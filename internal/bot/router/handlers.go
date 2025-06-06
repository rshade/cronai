// Package router provides event routing functionality for processing GitHub webhook events.
package router

import (
	"encoding/json"
	"fmt"

	"github.com/rshade/cronai/internal/logger"
	"github.com/rshade/cronai/internal/models"
	"github.com/rshade/cronai/internal/processor"
)

// baseHandler provides common functionality for all event handlers
type baseHandler struct {
	model          models.ModelClient
	processor      processor.Processor
	promptTemplate string
	logger         *logger.Logger
}

// generatePrompt creates a prompt from event data
func (h *baseHandler) generatePrompt(eventType string, data interface{}) string {
	// Format event data as JSON for the prompt
	eventJSON, _ := json.MarshalIndent(data, "", "  ") //nolint:errcheck // Fallback handled below

	prompt := fmt.Sprintf(`Analyze this GitHub %s event and provide insights:

Event Data:
%s

Please provide:
1. A summary of what happened
2. Any notable patterns or concerns
3. Recommended actions (if any)`, eventType, string(eventJSON))

	if h.promptTemplate != "" {
		prompt = h.promptTemplate
	}

	return prompt
}

// processWithAI sends the prompt to the AI model and processes the response
func (h *baseHandler) processWithAI(prompt string, metadata map[string]string) error {
	// Execute the model
	response, err := h.model.Execute(prompt)
	if err != nil {
		return fmt.Errorf("model execution failed: %w", err)
	}

	// Add metadata to response
	if response.Variables == nil {
		response.Variables = make(map[string]string)
	}
	for k, v := range metadata {
		response.Variables[k] = v
	}

	// Process the response if a processor is configured
	if h.processor != nil {
		// Use a default template name based on event type
		templateName := metadata["event_type"] + "_template"
		if err := h.processor.Process(response, templateName); err != nil {
			return fmt.Errorf("processor failed: %w", err)
		}
	}

	return nil
}

// IssuesHandler handles GitHub issues events
type IssuesHandler struct {
	baseHandler
}

// NewIssuesHandler creates a new issues event handler
func NewIssuesHandler(model models.ModelClient, processor processor.Processor) *IssuesHandler {
	return &IssuesHandler{
		baseHandler: baseHandler{
			model:     model,
			processor: processor,
			logger:    logger.GetLogger(),
		},
	}
}

// Handle processes an issues event
func (h *IssuesHandler) Handle(event Event) error {
	h.logger.Info("Handling issues event", logger.Fields{"action": event.Action})

	// Parse the payload
	var data struct {
		Action string `json:"action"`
		Issue  struct {
			Number int    `json:"number"`
			Title  string `json:"title"`
			Body   string `json:"body"`
			State  string `json:"state"`
			User   struct {
				Login string `json:"login"`
			} `json:"user"`
		} `json:"issue"`
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return fmt.Errorf("failed to parse issues event: %w", err)
	}

	// Generate prompt
	prompt := h.generatePrompt("issues", data)

	// Add metadata
	metadata := map[string]string{
		"event_type":   "issues",
		"action":       event.Action,
		"issue_number": fmt.Sprintf("%d", data.Issue.Number),
		"repository":   fmt.Sprintf("%s/%s", data.Repository.Owner.Login, data.Repository.Name),
	}

	return h.processWithAI(prompt, metadata)
}

// PullRequestHandler handles GitHub pull request events
type PullRequestHandler struct {
	baseHandler
}

// NewPullRequestHandler creates a new pull request event handler
func NewPullRequestHandler(model models.ModelClient, processor processor.Processor) *PullRequestHandler {
	return &PullRequestHandler{
		baseHandler: baseHandler{
			model:     model,
			processor: processor,
			logger:    logger.GetLogger(),
		},
	}
}

// Handle processes a pull request event
func (h *PullRequestHandler) Handle(event Event) error {
	h.logger.Info("Handling pull request event", logger.Fields{"action": event.Action})

	// Parse the payload
	var data struct {
		Action      string `json:"action"`
		Number      int    `json:"number"`
		PullRequest struct {
			Title string `json:"title"`
			Body  string `json:"body"`
			State string `json:"state"`
			Head  struct {
				Ref string `json:"ref"`
			} `json:"head"`
			Base struct {
				Ref string `json:"ref"`
			} `json:"base"`
			User struct {
				Login string `json:"login"`
			} `json:"user"`
		} `json:"pull_request"`
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return fmt.Errorf("failed to parse pull request event: %w", err)
	}

	// Generate prompt
	prompt := h.generatePrompt("pull_request", data)

	// Add metadata
	metadata := map[string]string{
		"event_type": "pull_request",
		"action":     event.Action,
		"pr_number":  fmt.Sprintf("%d", data.Number),
		"repository": fmt.Sprintf("%s/%s", data.Repository.Owner.Login, data.Repository.Name),
	}

	return h.processWithAI(prompt, metadata)
}

// PushHandler handles GitHub push events
type PushHandler struct {
	baseHandler
}

// NewPushHandler creates a new push event handler
func NewPushHandler(model models.ModelClient, processor processor.Processor) *PushHandler {
	return &PushHandler{
		baseHandler: baseHandler{
			model:     model,
			processor: processor,
			logger:    logger.GetLogger(),
		},
	}
}

// Handle processes a push event
func (h *PushHandler) Handle(event Event) error {
	h.logger.Info("Handling push event", logger.Fields{})

	// Parse the payload
	var data struct {
		Ref     string `json:"ref"`
		Before  string `json:"before"`
		After   string `json:"after"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"commits"`
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return fmt.Errorf("failed to parse push event: %w", err)
	}

	// Generate prompt
	prompt := h.generatePrompt("push", data)

	// Add metadata
	metadata := map[string]string{
		"event_type":   "push",
		"ref":          data.Ref,
		"commit_count": fmt.Sprintf("%d", len(data.Commits)),
		"repository":   fmt.Sprintf("%s/%s", data.Repository.Owner.Login, data.Repository.Name),
	}

	return h.processWithAI(prompt, metadata)
}

// ReleaseHandler handles GitHub release events
type ReleaseHandler struct {
	baseHandler
}

// NewReleaseHandler creates a new release event handler
func NewReleaseHandler(model models.ModelClient, processor processor.Processor) *ReleaseHandler {
	return &ReleaseHandler{
		baseHandler: baseHandler{
			model:     model,
			processor: processor,
			logger:    logger.GetLogger(),
		},
	}
}

// Handle processes a release event
func (h *ReleaseHandler) Handle(event Event) error {
	h.logger.Info("Handling release event", logger.Fields{"action": event.Action})

	// Parse the payload
	var data struct {
		Action  string `json:"action"`
		Release struct {
			TagName    string `json:"tag_name"`
			Name       string `json:"name"`
			Body       string `json:"body"`
			Draft      bool   `json:"draft"`
			Prerelease bool   `json:"prerelease"`
			Author     struct {
				Login string `json:"login"`
			} `json:"author"`
		} `json:"release"`
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return fmt.Errorf("failed to parse release event: %w", err)
	}

	// Generate prompt
	prompt := h.generatePrompt("release", data)

	// Add metadata
	metadata := map[string]string{
		"event_type": "release",
		"action":     event.Action,
		"tag":        data.Release.TagName,
		"repository": fmt.Sprintf("%s/%s", data.Repository.Owner.Login, data.Repository.Name),
	}

	return h.processWithAI(prompt, metadata)
}
