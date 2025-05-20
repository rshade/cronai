# CronAI Architecture

This document provides an overview of the CronAI system architecture, explaining the key components, their interactions, and design principles.

## System Overview

CronAI is built as a modular Go application that connects scheduled tasks to AI models and processes their responses. The system follows a pipeline architecture:

```text
Configuration → Scheduling → Prompt Management → Model Execution → Response Processing
```text

## Components Diagram

```text
┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐
│            │     │            │     │            │     │            │     │            │
│    CLI     │────▶│    Cron    │────▶│   Prompt   │────▶│   Models   │────▶│ Processors │
│  Commands  │     │  Service   │     │  Manager   │     │   Client   │     │            │
│            │     │            │     │            │     │            │     │            │
└────────────┘     └────────────┘     └────────────┘     └────────────┘     └────────────┘
       │                 ▲                  ▲                  ▲                  ▲
       │                 │                  │                  │                  │
       │                 │                  │                  │                  │
       ▼                 │                  │                  │                  │
┌────────────┐     ┌─────┴──────┐     ┌─────┴──────┐     ┌─────┴──────┐     ┌─────┴──────┐
│            │     │            │     │            │     │            │     │            │
│   Config   │────▶│Environment │     │  Template  │     │   Model    │     │  Template  │
│  Manager   │     │ Variables  │     │  System    │     │   Config   │     │  System    │
│            │     │            │     │            │     │            │     │            │
└────────────┘     └────────────┘     └────────────┘     └────────────┘     └────────────┘
```text

## Key Components

### 1. Command-Line Interface (CLI)

**Location**: `cmd/cronai/`
**Responsibility**: Provides the user interface to interact with the system
**Key Features**:

- Command parsing using Cobra framework
- Start/stop/run commands for service control
- Prompt management commands
- Configuration validation

### 2. Configuration Management

**Location**: `pkg/config/`
**Responsibility**: Loading and validating configuration
**Key Features**:

- Model parameter configuration
- Environment variable management
- Configuration file parsing
- Parameter validation

### 3. Cron Scheduling Service

**Location**: `internal/cron/`
**Responsibility**: Manages scheduled task execution
**Key Features**:

- Parses cron expressions
- Schedules tasks based on configuration
- Manages task lifecycle
- Handles service start/stop

### 4. Prompt Management

**Location**: `internal/prompt/`
**Responsibility**: Manages prompt loading and preprocessing
**Key Features**:

- File-based prompt loading
- Variable substitution
- Prompt metadata parsing
- Prompt searching and listing

### 5. Model Execution

**Location**: `internal/models/`
**Responsibility**: Communicates with AI model APIs
**Key Features**:

- Model client interface abstraction
- Multiple model support (OpenAI, Claude, Gemini)
- Standard response format
- Fallback mechanism
- Error handling and retries

### 6. Response Processing

**Location**: `internal/processor/`
**Responsibility**: Processes model responses into output formats
**Key Features**:

- Processor interface for consistent handling
- Multiple output channels (File, GitHub, Console)
- Registry pattern for processor management
- Configuration validation

### 7. Templating System

**Location**: `internal/processor/template/`
**Responsibility**: Formats output based on templates
**Key Features**:

- Template loading and management
- Standard template variables
- Output formatting

## Interfaces and Design Patterns

CronAI employs several design patterns to maintain clean architecture:

### 1. Interface-Based Design

The system uses interfaces to define clear boundaries between components:

```go
// ModelClient defines the interface for AI model clients
type ModelClient interface {
    Execute(promptContent string) (*ModelResponse, error)
}

// Processor defines the interface for response processors
type Processor interface {
    Process(response *models.ModelResponse, templateName string) error
    Validate() error
    GetType() string
    GetConfig() ProcessorConfig
}
```text

### 2. Factory Pattern

Used in the processor system to create processor instances dynamically:

```go
// RegisterProcessor adds a processor factory to the registry
func RegisterProcessor(processorType string, factory ProcessorFactory)

// GetProcessor creates a processor of the specified type
func GetProcessor(processorType string, config ProcessorConfig) (Processor, error)
```text

### 3. Singleton Pattern

Used for managers that need to maintain global state:

```go
// GetInstance returns the singleton template manager instance
func GetInstance() *Manager
```text

### 4. Registry Pattern

Used to register and manage available processors:

```go
// In registry.go, processors register themselves:
func init() {
    RegisterProcessor("file", NewFileProcessor)
    RegisterProcessor("github", NewGithubProcessor)
    RegisterProcessor("console", NewConsoleProcessor)
}
```text

## Data Flow

### Configuration to Execution Flow

1. **Configuration Loading**:
   - Parse configuration file with cron schedule, model, prompt, and processor
   - Load environment variables for API keys and other settings

2. **Scheduling**:
   - Cron service parses schedule and creates tasks
   - Tasks are scheduled using cron library

3. **Task Execution**:
   - When triggered, task loads the specified prompt
   - Variables are replaced in the prompt content
   - Prompt is passed to the specified model

4. **Model Execution**:
   - ModelClient for the specified model is created
   - Model parameters are applied
   - Prompt is sent to the AI API
   - Response is received and standardized

5. **Response Processing**:
   - Appropriate processor is created based on configuration
   - Processor formats and delivers the response
   - Output is sent to the configured destination

## Error Handling

The system uses a consistent error handling pattern:

1. **Error Categorization**:
   - Errors are categorized (Configuration, Validation, Application, IO)
   - Structured logging with error context

2. **Graceful Degradation**:
   - Model fallback mechanism when primary model fails
   - Retry logic with configurable attempts

3. **Validation Hierarchy**:
   - Configuration validation before execution
   - Input validation at each processing stage
   - Clear error messages for troubleshooting

## Extension Points

CronAI is designed to be extended in several ways:

1. **New Models**:
   - Implement the `ModelClient` interface
   - Register in the `defaultCreateModelClient` function

2. **New Processors**:
   - Implement the `Processor` interface
   - Register with the processor registry

3. **New CLI Commands**:
   - Add new commands to the Cobra command structure
   - Follow the existing pattern in `cmd/cronai/cmd/`

## Testing Strategy

The architecture supports comprehensive testing:

1. **Unit Testing**:
   - Each component can be tested in isolation
   - Mock implementations of interfaces
   - Table-driven tests for different scenarios

2. **Integration Testing**:
   - End-to-end workflow tests
   - Configuration validation tests
   - Real external services can be mocked

## Current Limitations

The MVP architecture has some known limitations:

1. No automatic handling of API rate limits
2. No persistent storage for response history
3. Limited response processor options
4. No web UI for management
5. No response templating capabilities yet

These limitations are planned to be addressed in post-MVP releases.
