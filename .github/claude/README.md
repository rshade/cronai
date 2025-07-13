# Claude GitHub Assistant Configuration

This directory contains configuration and prompts for the Claude GitHub Assistant integration with multiple personas.

## Persona System

Claude can take on different roles based on GitHub labels:

### Available Personas

| Label | File | Role | Best For |
|-------|------|------|----------|
| `claude-reviewer` | `code-reviewer.md` | Code Reviewer | PR reviews, code quality, security analysis |
| `claude-engineer` | `software-engineer.md` | Software Engineer | Implementation, bug fixes, feature development |
| `claude-assistant` | `default.md` | General Assistant | Questions, documentation, general help |
| No label | `default.md` | General Assistant | Default behavior |

## How Persona Selection Works

1. **Label Priority**: The workflow checks for persona-specific labels on issues/PRs
2. **Default Fallback**: If no persona label is found, uses the default assistant
3. **Automatic Selection**: The system automatically loads the appropriate persona file

## Usage Examples

### Code Review
```
# Add this label to your PR:
claude-reviewer

# Then mention Claude in a comment:
@claude Please review this implementation for security and performance issues.
```

### Implementation Help
```
# Add this label to your issue:
claude-engineer

# Then describe your problem:
@claude I need help implementing a new processor for Microsoft Teams integration.
```

### General Questions
```
# Add this label (or no label):
claude-assistant

# Ask your question:
@claude How do I configure the Slack processor for webhooks?
```

## Persona Capabilities

### 🔍 Code Reviewer (`claude-reviewer`)
- **Focus**: Code quality, security, performance
- **Output**: Detailed review comments with severity levels
- **Best for**: PR reviews, security audits, architecture review
- **Style**: Structured, thorough, educational

### 🛠️ Software Engineer (`claude-engineer`)
- **Focus**: Implementation, problem-solving, system design
- **Output**: Working code, step-by-step implementation plans
- **Best for**: Feature development, bug fixes, technical guidance
- **Style**: Practical, solution-oriented, code-heavy

### 💬 General Assistant (`claude-assistant` or default)
- **Focus**: Questions, documentation, general guidance
- **Output**: Helpful explanations and guidance
- **Best for**: Documentation, troubleshooting, project questions
- **Style**: Adaptive, informative, user-friendly

## Triggering Claude

### Method 1: Labels + Mention
```bash
# Add appropriate label to issue/PR, then:
@claude [your request]
```

### Method 2: Issue Templates
Use the provided issue templates that automatically add the correct labels.

### Method 3: Direct Mention
```bash
# Mention Claude directly (uses default persona):
@claude [your request]
```

## File Structure

```
.github/claude/
├── README.md              # This documentation
├── default.md             # Default assistant persona
├── code-reviewer.md       # Code reviewer persona
└── software-engineer.md   # Software engineer persona
```

## Customization

### Adding New Personas
1. Create a new `.md` file in this directory
2. Update the workflow in `.github/workflows/claude.yml` to recognize new labels
3. Add the persona to this documentation

### Modifying Existing Personas
1. Edit the appropriate `.md` file
2. Changes take effect immediately on next interaction
3. Consider testing with a small issue first

## Best Practices

### Choosing the Right Persona
- **PR Reviews**: Always use `claude-reviewer` for comprehensive code review
- **Implementation**: Use `claude-engineer` when you need working code
- **Questions**: Use `claude-assistant` for general inquiries
- **Complex Issues**: You can change labels during the conversation

### Label Management
- Add persona labels when creating issues/PRs
- Change labels if you need a different type of assistance
- Remove old persona labels when switching to avoid conflicts

### Conversation Management
- Be specific about what you need
- Provide context and relevant code/error messages
- Use follow-up questions to dive deeper into topics

### Cost Management
- Conversations are limited to 10 turns to prevent runaway costs
- Each persona interaction counts as a separate conversation
- Be concise but complete in your requests

## Integration with CLAUDE.md

The persona system works alongside your `CLAUDE.md` file:
- **CLAUDE.md**: Provides project context and conventions
- **Persona files**: Define the assistant's role and behavior
- **Combined**: Claude gets both project knowledge and role-specific guidance

## Troubleshooting

### Claude Not Responding
1. Check that the `ANTHROPIC_API_KEY` secret is set
2. Verify the label is spelled correctly
3. Ensure you mentioned `@claude` in your comment

### Wrong Persona
1. Check the labels on your issue/PR
2. Remove unwanted persona labels
3. Add the correct persona label
4. Mention `@claude` again

### Need Different Assistance
1. Change the label on your issue/PR
2. Mention `@claude` with your new request
3. The system will use the new persona for subsequent responses