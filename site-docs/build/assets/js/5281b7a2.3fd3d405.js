"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[443],{

/***/ 5874:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_architecture_md_528_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-architecture-md-528.json
const site_docs_architecture_md_528_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"architecture","title":"CronAI Architecture","description":"This document provides an overview of the CronAI system architecture, explaining the key components, their interactions, and design principles.","source":"@site/docs/architecture.md","sourceDirName":".","slug":"/architecture","permalink":"/cronai/docs/architecture","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/architecture.md","tags":[],"version":"current","frontMatter":{"id":"architecture","title":"CronAI Architecture","sidebar_label":"Architecture"},"sidebar":"tutorialSidebar","previous":{"title":"Introduction","permalink":"/cronai/docs/"},"next":{"title":"CronAI: Known Limitations and Future Improvements","permalink":"/cronai/docs/limitations-and-improvements"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/architecture.md


const frontMatter = {
	id: 'architecture',
	title: 'CronAI Architecture',
	sidebar_label: 'Architecture'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "System Overview",
  "id": "system-overview",
  "level": 2
}, {
  "value": "Components Diagram",
  "id": "components-diagram",
  "level": 2
}, {
  "value": "Key Components",
  "id": "key-components",
  "level": 2
}, {
  "value": "1. Command-Line Interface (CLI)",
  "id": "1-command-line-interface-cli",
  "level": 3
}, {
  "value": "2. Configuration Management",
  "id": "2-configuration-management",
  "level": 3
}, {
  "value": "3. Cron Scheduling Service",
  "id": "3-cron-scheduling-service",
  "level": 3
}, {
  "value": "4. Prompt Management",
  "id": "4-prompt-management",
  "level": 3
}, {
  "value": "5. Model Execution",
  "id": "5-model-execution",
  "level": 3
}, {
  "value": "6. Response Processing",
  "id": "6-response-processing",
  "level": 3
}, {
  "value": "7. Templating System",
  "id": "7-templating-system",
  "level": 3
}, {
  "value": "Interfaces and Design Patterns",
  "id": "interfaces-and-design-patterns",
  "level": 2
}, {
  "value": "1. Interface-Based Design",
  "id": "1-interface-based-design",
  "level": 3
}, {
  "value": "2. Factory Pattern",
  "id": "2-factory-pattern",
  "level": 3
}, {
  "value": "3. Singleton Pattern",
  "id": "3-singleton-pattern",
  "level": 3
}, {
  "value": "4. Registry Pattern",
  "id": "4-registry-pattern",
  "level": 3
}, {
  "value": "Data Flow",
  "id": "data-flow",
  "level": 2
}, {
  "value": "Configuration to Execution Flow",
  "id": "configuration-to-execution-flow",
  "level": 3
}, {
  "value": "Error Handling",
  "id": "error-handling",
  "level": 2
}, {
  "value": "Extension Points",
  "id": "extension-points",
  "level": 2
}, {
  "value": "Testing Strategy",
  "id": "testing-strategy",
  "level": 2
}, {
  "value": "Current Limitations",
  "id": "current-limitations",
  "level": 2
}];
function _createMdxContent(props) {
  const _components = {
    br: "br",
    code: "code",
    h2: "h2",
    h3: "h3",
    li: "li",
    ol: "ol",
    p: "p",
    pre: "pre",
    strong: "strong",
    ul: "ul",
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return (0,jsx_runtime.jsxs)(jsx_runtime.Fragment, {
    children: [(0,jsx_runtime.jsx)(_components.p, {
      children: "This document provides an overview of the CronAI system architecture, explaining the key components, their interactions, and design principles."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "system-overview",
      children: "System Overview"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI is built as a modular Go application that connects scheduled tasks to AI models and processes their responses. The system follows a pipeline architecture:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "Configuration → Scheduling → Prompt Management → Model Execution → Response Processing\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "components-diagram",
      children: "Components Diagram"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-ascii",
        children: "┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐     ┌────────────┐\n│            │     │            │     │            │     │            │     │            │\n│    CLI     │────▶│    Cron    │────▶│   Prompt   │────▶│   Models   │────▶│ Processors │\n│  Commands  │     │  Service   │     │  Manager   │     │   Client   │     │            │\n│            │     │            │     │            │     │            │     │            │\n└────────────┘     └────────────┘     └────────────┘     └────────────┘     └────────────┘\n       │                 ▲                  ▲                  ▲                  ▲\n       │                 │                  │                  │                  │\n       │                 │                  │                  │                  │\n       ▼                 │                  │                  │                  │\n┌────────────┐     ┌─────┴──────┐     ┌─────┴──────┐     ┌─────┴──────┐     ┌─────┴──────┐\n│            │     │            │     │            │     │            │     │            │\n│   Config   │────▶│Environment │     │  Template  │     │   Model    │     │  Template  │\n│  Manager   │     │ Variables  │     │  System    │     │   Config   │     │  System    │\n│            │     │            │     │            │     │            │     │            │\n└────────────┘     └────────────┘     └────────────┘     └────────────┘     └────────────┘\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "key-components",
      children: "Key Components"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "1-command-line-interface-cli",
      children: "1. Command-Line Interface (CLI)"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "cmd/cronai/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Provides the user interface to interact with the system", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Command parsing using Cobra framework"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Start/stop/run commands for service control"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Prompt management commands"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Configuration validation"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "2-configuration-management",
      children: "2. Configuration Management"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "pkg/config/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Loading and validating configuration", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Model parameter configuration"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Environment variable management"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Configuration file parsing"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Parameter validation"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "3-cron-scheduling-service",
      children: "3. Cron Scheduling Service"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "internal/cron/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Manages scheduled task execution", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Parses cron expressions"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Schedules tasks based on configuration"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Manages task lifecycle"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Handles service start/stop"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "4-prompt-management",
      children: "4. Prompt Management"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "internal/prompt/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Manages prompt loading and preprocessing", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "File-based prompt loading"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Variable substitution"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Prompt metadata parsing"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Prompt searching and listing"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "5-model-execution",
      children: "5. Model Execution"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "internal/models/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Communicates with AI model APIs", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Model client interface abstraction"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Multiple model support (OpenAI, Claude, Gemini)"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Standard response format"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Fallback mechanism"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Error handling and retries"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "6-response-processing",
      children: "6. Response Processing"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "internal/processor/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Processes model responses into output formats", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Processor interface for consistent handling"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Multiple output channels (File, GitHub, Console)"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Registry pattern for processor management"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Configuration validation"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "7-templating-system",
      children: "7. Templating System"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: [(0,jsx_runtime.jsx)(_components.strong, {
        children: "Location"
      }), ": ", (0,jsx_runtime.jsx)(_components.code, {
        children: "internal/processor/template/"
      }), (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Responsibility"
      }), ": Formats output based on templates", (0,jsx_runtime.jsx)(_components.br, {}), "\n", (0,jsx_runtime.jsx)(_components.strong, {
        children: "Key Features"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Template loading and management"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Standard template variables"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Output formatting"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "interfaces-and-design-patterns",
      children: "Interfaces and Design Patterns"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI employs several design patterns to maintain clean architecture:"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "1-interface-based-design",
      children: "1. Interface-Based Design"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The system uses interfaces to define clear boundaries between components:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-go",
        children: "// ModelClient defines the interface for AI model clients\ntype ModelClient interface {\n    Execute(promptContent string) (*ModelResponse, error)\n}\n\n// Processor defines the interface for response processors\ntype Processor interface {\n    Process(response *models.ModelResponse, templateName string) error\n    Validate() error\n    GetType() string\n    GetConfig() ProcessorConfig\n}\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "2-factory-pattern",
      children: "2. Factory Pattern"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Used in the processor system to create processor instances dynamically:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-go",
        children: "// RegisterProcessor adds a processor factory to the registry\nfunc RegisterProcessor(processorType string, factory ProcessorFactory)\n\n// GetProcessor creates a processor of the specified type\nfunc GetProcessor(processorType string, config ProcessorConfig) (Processor, error)\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "3-singleton-pattern",
      children: "3. Singleton Pattern"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Used for managers that need to maintain global state:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-go",
        children: "// GetInstance returns the singleton template manager instance\nfunc GetInstance() *Manager\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "4-registry-pattern",
      children: "4. Registry Pattern"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Used to register and manage available processors:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-go",
        children: "// In registry.go, processors register themselves:\nfunc init() {\n    RegisterProcessor(\"file\", NewFileProcessor)\n    RegisterProcessor(\"github\", NewGithubProcessor)\n    RegisterProcessor(\"console\", NewConsoleProcessor)\n}\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "data-flow",
      children: "Data Flow"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "configuration-to-execution-flow",
      children: "Configuration to Execution Flow"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Configuration Loading"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Parse configuration file with cron schedule, model, prompt, and processor"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Load environment variables for API keys and other settings"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Scheduling"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Cron service parses schedule and creates tasks"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Tasks are scheduled using cron library"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Task Execution"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "When triggered, task loads the specified prompt"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Variables are replaced in the prompt content"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Prompt is passed to the specified model"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Model Execution"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "ModelClient for the specified model is created"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Model parameters are applied"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Prompt is sent to the AI API"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Response is received and standardized"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Response Processing"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Appropriate processor is created based on configuration"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Processor formats and delivers the response"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Output is sent to the configured destination"
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "error-handling",
      children: "Error Handling"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The system uses a consistent error handling pattern:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Error Categorization"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Errors are categorized (Configuration, Validation, Application, IO)"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Structured logging with error context"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Graceful Degradation"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Model fallback mechanism when primary model fails"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Retry logic with configurable attempts"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Validation Hierarchy"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Configuration validation before execution"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Input validation at each processing stage"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Clear error messages for troubleshooting"
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "extension-points",
      children: "Extension Points"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI is designed to be extended in several ways:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "New Models"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: ["Implement the ", (0,jsx_runtime.jsx)(_components.code, {
              children: "ModelClient"
            }), " interface"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: ["Register in the ", (0,jsx_runtime.jsx)(_components.code, {
              children: "defaultCreateModelClient"
            }), " function"]
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "New Processors"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: ["Implement the ", (0,jsx_runtime.jsx)(_components.code, {
              children: "Processor"
            }), " interface"]
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Register with the processor registry"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "New CLI Commands"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Add new commands to the Cobra command structure"
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: ["Follow the existing pattern in ", (0,jsx_runtime.jsx)(_components.code, {
              children: "cmd/cronai/cmd/"
            })]
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "testing-strategy",
      children: "Testing Strategy"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The architecture supports comprehensive testing:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Unit Testing"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Each component can be tested in isolation"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Mock implementations of interfaces"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Table-driven tests for different scenarios"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
          children: [(0,jsx_runtime.jsx)(_components.strong, {
            children: "Integration Testing"
          }), ":"]
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "End-to-end workflow tests"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Configuration validation tests"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Real external services can be mocked"
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "current-limitations",
      children: "Current Limitations"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The MVP architecture has some known limitations:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "No automatic handling of API rate limits"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "No persistent storage for response history"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Limited response processor options"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "No web UI for management"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "No response templating capabilities yet"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "These limitations are planned to be addressed in post-MVP releases."
    })]
  });
}
function MDXContent(props = {}) {
  const {wrapper: MDXLayout} = {
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return MDXLayout ? (0,jsx_runtime.jsx)(MDXLayout, {
    ...props,
    children: (0,jsx_runtime.jsx)(_createMdxContent, {
      ...props
    })
  }) : _createMdxContent(props);
}



/***/ }),

/***/ 8453:
/***/ ((__unused_webpack___webpack_module__, __webpack_exports__, __webpack_require__) => {

/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   R: () => (/* binding */ useMDXComponents),
/* harmony export */   x: () => (/* binding */ MDXProvider)
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(6540);
/**
 * @import {MDXComponents} from 'mdx/types.js'
 * @import {Component, ReactElement, ReactNode} from 'react'
 */

/**
 * @callback MergeComponents
 *   Custom merge function.
 * @param {Readonly<MDXComponents>} currentComponents
 *   Current components from the context.
 * @returns {MDXComponents}
 *   Additional components.
 *
 * @typedef Props
 *   Configuration for `MDXProvider`.
 * @property {ReactNode | null | undefined} [children]
 *   Children (optional).
 * @property {Readonly<MDXComponents> | MergeComponents | null | undefined} [components]
 *   Additional components to use or a function that creates them (optional).
 * @property {boolean | null | undefined} [disableParentContext=false]
 *   Turn off outer component context (default: `false`).
 */



/** @type {Readonly<MDXComponents>} */
const emptyComponents = {}

const MDXContext = react__WEBPACK_IMPORTED_MODULE_0__.createContext(emptyComponents)

/**
 * Get current components from the MDX Context.
 *
 * @param {Readonly<MDXComponents> | MergeComponents | null | undefined} [components]
 *   Additional components to use or a function that creates them (optional).
 * @returns {MDXComponents}
 *   Current components.
 */
function useMDXComponents(components) {
  const contextComponents = react__WEBPACK_IMPORTED_MODULE_0__.useContext(MDXContext)

  // Memoize to avoid unnecessary top-level context changes
  return react__WEBPACK_IMPORTED_MODULE_0__.useMemo(
    function () {
      // Custom merge via a function prop
      if (typeof components === 'function') {
        return components(contextComponents)
      }

      return {...contextComponents, ...components}
    },
    [contextComponents, components]
  )
}

/**
 * Provider for MDX context.
 *
 * @param {Readonly<Props>} properties
 *   Properties.
 * @returns {ReactElement}
 *   Element.
 * @satisfies {Component}
 */
function MDXProvider(properties) {
  /** @type {Readonly<MDXComponents>} */
  let allComponents

  if (properties.disableParentContext) {
    allComponents =
      typeof properties.components === 'function'
        ? properties.components(emptyComponents)
        : properties.components || emptyComponents
  } else {
    allComponents = useMDXComponents(properties.components)
  }

  return react__WEBPACK_IMPORTED_MODULE_0__.createElement(
    MDXContext.Provider,
    {value: allComponents},
    properties.children
  )
}


/***/ })

}]);