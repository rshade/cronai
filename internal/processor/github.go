package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v72/github"
	"github.com/rshade/cronai/internal/errors"
	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor/template"
	"golang.org/x/oauth2"
)

// GitHubProcessor handles GitHub operations (issues, comments, pull requests)
type GitHubProcessor struct {
	config Config
	client *github.Client
}

// NewGitHubProcessor creates a new GitHub processor
func NewGitHubProcessor(config Config) (Processor, error) {
	// Validate the target format
	if config.Target == "" {
		return nil, errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("github target cannot be empty"),
			"invalid github processor configuration")
	}

	// Initialize GitHub client with token
	token := os.Getenv(EnvGitHubToken)
	if token != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		return &GitHubProcessor{
			config: config,
			client: github.NewClient(tc),
		}, nil
	}

	// Create processor without client (will validate in Validate())
	return &GitHubProcessor{
		config: config,
	}, nil
}

// Process handles the model response with optional template
func (g *GitHubProcessor) Process(response *models.ModelResponse, templateName string) error {
	// Create template data
	tmplData := template.Data{
		Content:     response.Content,
		Model:       response.Model,
		Timestamp:   response.Timestamp,
		PromptName:  response.PromptName,
		Variables:   response.Variables,
		ExecutionID: response.ExecutionID,
		Metadata:    make(map[string]string),
	}

	// Add standard metadata fields
	tmplData.Metadata["timestamp"] = response.Timestamp.Format(time.RFC3339)
	tmplData.Metadata["date"] = response.Timestamp.Format("2006-01-02")
	tmplData.Metadata["time"] = response.Timestamp.Format("15:04:05")
	tmplData.Metadata["execution_id"] = response.ExecutionID
	tmplData.Metadata["processor"] = g.GetType()
	if templateName != "" {
		tmplData.Metadata["template"] = templateName
	}

	return g.processGitHubWithTemplate(g.config.Target, tmplData, templateName)
}

// Validate checks if the processor is properly configured
func (g *GitHubProcessor) Validate() error {
	if g.config.Target == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("github target cannot be empty"),
			"invalid github processor configuration")
	}

	// Skip token validation in tests
	if os.Getenv("GO_TEST") == "1" {
		return nil
	}

	// Check for required environment variables
	githubToken := os.Getenv(EnvGitHubToken)
	if githubToken == "" {
		return errors.Wrap(errors.CategoryConfiguration,
			fmt.Errorf("GITHUB_TOKEN environment variable not set"),
			"missing GitHub configuration")
	}

	// Validate target format (action:repo or issue:repo#number or pr:repo)
	parts := strings.Split(g.config.Target, ":")
	if len(parts) != 2 {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid github target format: %s, expected 'action:repo' format", g.config.Target),
			"invalid github processor configuration")
	}

	action := parts[0]
	validActions := []string{"issue", "comment", "pr"}
	valid := false
	for _, a := range validActions {
		if action == a {
			valid = true
			break
		}
	}
	if !valid {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid github action: %s, expected one of: %v", action, validActions),
			"invalid github processor configuration")
	}

	return nil
}

// GetType returns the processor type identifier
func (g *GitHubProcessor) GetType() string {
	return "github"
}

// GetConfig returns the processor configuration
func (g *GitHubProcessor) GetConfig() Config {
	return g.config
}

// processGitHubWithTemplate processes GitHub operations with templates
func (g *GitHubProcessor) processGitHubWithTemplate(target string, data template.Data, templateName string) error {
	// Ensure client is initialized
	if g.client == nil {
		// Check if running in test mode
		if os.Getenv("GO_TEST") == "1" {
			// Test mode - use a stub client
			g.client = github.NewClient(nil)
		} else {
			// Normal mode - use authenticated client
			token := os.Getenv(EnvGitHubToken)
			if token == "" {
				log.Error("GitHub token not set", logger.Fields{
					"target": target,
				})
				return errors.Wrap(errors.CategoryConfiguration, fmt.Errorf("GITHUB_TOKEN environment variable not set"),
					"missing GitHub configuration")
			}

			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)
			g.client = github.NewClient(tc)
		}
	}

	// Parse target format (action:repo or comment:repo#number)
	parts := strings.Split(target, ":")
	if len(parts) != 2 {
		return errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("invalid target format: %s", target),
			"GitHub target parsing failed")
	}

	action := parts[0]
	repoInfo := parts[1]

	// Get template manager
	manager := template.GetManager()

	// Use default template if none specified
	if templateName == "" {
		templateName = fmt.Sprintf("default_github_%s", action)
	}

	// Execute template to get payload
	payload := manager.SafeExecute(templateName, data)
	if payload == "" {
		log.Error("Failed to generate GitHub payload", logger.Fields{
			"template": templateName,
			"target":   target,
		})
		return errors.Wrap(errors.CategoryApplication, fmt.Errorf("empty payload generated from template %s", templateName),
			"GitHub payload generation failed")
	}

	// Add to metadata for logging
	data.Metadata["github_action"] = action
	data.Metadata["github_repo"] = repoInfo
	data.Metadata["template_used"] = templateName

	// Parse the payload as JSON to extract fields
	var jsonPayload map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &jsonPayload); err != nil {
		log.Error("Invalid GitHub JSON payload", logger.Fields{
			"template": templateName,
			"error":    err.Error(),
		})
		return errors.Wrap(errors.CategoryApplication, err, "GitHub payload is not valid JSON")
	}

	// Process based on action type
	switch action {
	case "issue":
		return g.processGitHubIssue(repoInfo, jsonPayload)
	case "comment":
		return g.processGitHubComment(repoInfo, jsonPayload)
	case "pr":
		return g.processGitHubPR(repoInfo, jsonPayload)
	default:
		return errors.Wrap(errors.CategoryApplication,
			fmt.Errorf("unsupported github action: %s", action),
			"GitHub action not supported")
	}
}

// processGitHubIssue creates a new GitHub issue
func (g *GitHubProcessor) processGitHubIssue(repo string, payload map[string]interface{}) error {
	title, ok := payload["title"].(string)
	if !ok || title == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("missing or invalid 'title' in payload"),
			"GitHub issue creation failed")
	}

	body, ok := payload["body"].(string)
	if !ok {
		body = "" // Body can be empty
	}

	labels := []string{}
	if labelList, ok := payload["labels"].([]interface{}); ok {
		for _, label := range labelList {
			if labelStr, ok := label.(string); ok {
				labels = append(labels, labelStr)
			}
		}
	}

	// Parse owner and repo
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid repo format: %s, expected 'owner/repo'", repo),
			"GitHub issue creation failed")
	}
	owner := parts[0]
	repoName := parts[1]

	// Create issue request
	// issueRequest := &github.IssueRequest{
	// 	Title:  &title,
	// 	Body:   &body,
	// 	Labels: &labels,
	// }

	// In MVP, just log rather than actually creating
	log.Info("Would create GitHub issue", logger.Fields{
		"owner":  owner,
		"repo":   repoName,
		"title":  title,
		"body":   body,
		"labels": labels,
	})

	// For production implementation:
	/*
		ctx := context.Background()
		issue, _, err := g.client.Issues.Create(ctx, owner, repoName, issueRequest)
		if err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "failed to create GitHub issue")
		}

		log.Info("Created GitHub issue", logger.Fields{
			"owner":  owner,
			"repo":   repoName,
			"issue":  issue.GetNumber(),
			"url":    issue.GetHTMLURL(),
		})
	*/

	return nil
}

// processGitHubComment adds a comment to an existing issue
func (g *GitHubProcessor) processGitHubComment(repoInfo string, payload map[string]interface{}) error {
	// Parse repo and issue number (format: owner/repo#123)
	parts := strings.Split(repoInfo, "#")
	if len(parts) != 2 {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid comment target format: %s, expected 'owner/repo#123'", repoInfo),
			"GitHub comment creation failed")
	}

	repo := parts[0]
	issueNumberStr := parts[1]

	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid issue number: %s", issueNumberStr),
			"GitHub comment creation failed")
	}

	body, ok := payload["body"].(string)
	if !ok || body == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("missing or invalid 'body' in payload"),
			"GitHub comment creation failed")
	}

	// Parse owner and repo
	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid repo format: %s, expected 'owner/repo'", repo),
			"GitHub comment creation failed")
	}
	owner := repoParts[0]
	repoName := repoParts[1]

	// Create comment request
	// commentRequest := &github.IssueComment{
	// 	Body: &body,
	// }

	// In MVP, just log rather than actually creating
	log.Info("Would add GitHub comment", logger.Fields{
		"owner": owner,
		"repo":  repoName,
		"issue": issueNumber,
		"body":  body,
	})

	// For production implementation:
	/*
		ctx := context.Background()
		comment, _, err := g.client.Issues.CreateComment(ctx, owner, repoName, issueNumber, commentRequest)
		if err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "failed to add GitHub comment")
		}

		log.Info("Added GitHub comment", logger.Fields{
			"owner":  owner,
			"repo":   repoName,
			"issue":  issueNumber,
			"url":    comment.GetHTMLURL(),
		})
	*/

	return nil
}

// processGitHubPR creates a new pull request
func (g *GitHubProcessor) processGitHubPR(repo string, payload map[string]interface{}) error {
	title, ok := payload["title"].(string)
	if !ok || title == "" {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("missing or invalid 'title' in payload"),
			"GitHub PR creation failed")
	}

	body, ok := payload["body"].(string)
	if !ok {
		body = "" // Body can be empty
	}

	// Parse owner and repo
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("invalid repo format: %s, expected 'owner/repo'", repo),
			"GitHub PR creation failed")
	}
	owner := parts[0]
	repoName := parts[1]

	// Optional fields
	baseBranch := "main"
	if base, ok := payload["base"].(string); ok && base != "" {
		baseBranch = base
	}

	headBranch := ""
	if head, ok := payload["head"].(string); ok && head != "" {
		headBranch = head
	} else {
		return errors.Wrap(errors.CategoryValidation,
			fmt.Errorf("missing 'head' branch in payload"),
			"GitHub PR creation failed")
	}

	// Create PR request
	// prRequest := &github.NewPullRequest{
	// 	Title: &title,
	// 	Body:  &body,
	// 	Base:  &baseBranch,
	// 	Head:  &headBranch,
	// }

	// In MVP, just log rather than actually creating
	log.Info("Would create GitHub PR", logger.Fields{
		"owner": owner,
		"repo":  repoName,
		"title": title,
		"body":  body,
		"base":  baseBranch,
		"head":  headBranch,
	})

	// For production implementation:
	/*
		ctx := context.Background()
		pr, _, err := g.client.PullRequests.Create(ctx, owner, repoName, prRequest)
		if err != nil {
			return errors.Wrap(errors.CategoryExternal, err, "failed to create GitHub PR")
		}

		log.Info("Created GitHub PR", logger.Fields{
			"owner":  owner,
			"repo":   repoName,
			"pr":     pr.GetNumber(),
			"url":    pr.GetHTMLURL(),
		})
	*/

	return nil
}
