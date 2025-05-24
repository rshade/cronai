"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[419],{

/***/ 7637:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_model_parameters_md_8bc_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-model-parameters-md-8bc.json
const site_docs_model_parameters_md_8bc_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"model-parameters","title":"Model Parameters Configuration","description":"CronAI supports model-specific parameters that allow you to fine-tune AI model behavior for each prompt. This document explains how to configure and use these parameters.","source":"@site/docs/model-parameters.md","sourceDirName":".","slug":"/model-parameters","permalink":"/cronai/docs/model-parameters","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/model-parameters.md","tags":[],"version":"current","frontMatter":{"id":"model-parameters","title":"Model Parameters Configuration","sidebar_label":"Model Parameters"},"sidebar":"tutorialSidebar","previous":{"title":"Prompt Management","permalink":"/cronai/docs/prompt-management"},"next":{"title":"Logging in CronAI","permalink":"/cronai/docs/logging"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/model-parameters.md


const frontMatter = {
	id: 'model-parameters',
	title: 'Model Parameters Configuration',
	sidebar_label: 'Model Parameters'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "Supported Parameters",
  "id": "supported-parameters",
  "level": 2
}, {
  "value": "Common Parameters",
  "id": "common-parameters",
  "level": 3
}, {
  "value": "Model-Specific Parameters",
  "id": "model-specific-parameters",
  "level": 3
}, {
  "value": "OpenAI",
  "id": "openai",
  "level": 4
}, {
  "value": "Claude",
  "id": "claude",
  "level": 4
}, {
  "value": "Gemini",
  "id": "gemini",
  "level": 4
}, {
  "value": "Model-Specific Default Values",
  "id": "model-specific-default-values",
  "level": 2
}, {
  "value": "OpenAI Default Settings",
  "id": "openai-default-settings",
  "level": 3
}, {
  "value": "Claude Default Settings",
  "id": "claude-default-settings",
  "level": 3
}, {
  "value": "Gemini Default Settings",
  "id": "gemini-default-settings",
  "level": 3
}, {
  "value": "SDK Implementation",
  "id": "sdk-implementation",
  "level": 2
}, {
  "value": "Configuration Methods",
  "id": "configuration-methods",
  "level": 2
}, {
  "value": "1. Task-specific Configuration",
  "id": "1-task-specific-configuration",
  "level": 3
}, {
  "value": "Using Model-Specific Parameters",
  "id": "using-model-specific-parameters",
  "level": 4
}, {
  "value": "2. Environment Variables",
  "id": "2-environment-variables",
  "level": 3
}, {
  "value": "3. Command Line Parameters",
  "id": "3-command-line-parameters",
  "level": 3
}, {
  "value": "Examples",
  "id": "examples",
  "level": 2
}, {
  "value": "Low Temperature for Consistent Output",
  "id": "low-temperature-for-consistent-output",
  "level": 3
}, {
  "value": "Specific Model Version",
  "id": "specific-model-version",
  "level": 3
}, {
  "value": "GitHub Processor Example",
  "id": "github-processor-example",
  "level": 3
}, {
  "value": "Console Output Example",
  "id": "console-output-example",
  "level": 3
}, {
  "value": "Integration with Variables",
  "id": "integration-with-variables",
  "level": 3
}, {
  "value": "Advanced Configuration",
  "id": "advanced-configuration",
  "level": 2
}, {
  "value": "Custom API Endpoints",
  "id": "custom-api-endpoints",
  "level": 3
}, {
  "value": "Timeout Configuration",
  "id": "timeout-configuration",
  "level": 3
}, {
  "value": "Post-MVP Features",
  "id": "post-mvp-features",
  "level": 2
}];
function _createMdxContent(props) {
  const _components = {
    code: "code",
    h2: "h2",
    h3: "h3",
    h4: "h4",
    li: "li",
    ol: "ol",
    p: "p",
    pre: "pre",
    strong: "strong",
    table: "table",
    tbody: "tbody",
    td: "td",
    th: "th",
    thead: "thead",
    tr: "tr",
    ul: "ul",
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return (0,jsx_runtime.jsxs)(jsx_runtime.Fragment, {
    children: [(0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI supports model-specific parameters that allow you to fine-tune AI model behavior for each prompt. This document explains how to configure and use these parameters."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "supported-parameters",
      children: "Supported Parameters"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "common-parameters",
      children: "Common Parameters"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following parameters are supported across all models:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.table, {
      children: [(0,jsx_runtime.jsx)(_components.thead, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.th, {
            children: "Parameter"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Type"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Range"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Description"
          })]
        })
      }), (0,jsx_runtime.jsxs)(_components.tbody, {
        children: [(0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "temperature"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "float"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "0.0 - 1.0"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Controls response randomness (higher = more random)"
          })]
        }), (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "max_tokens"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "int"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "> 0"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Maximum number of tokens to generate"
          })]
        }), (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "model"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "string"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "-"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Specific model version to use"
          })]
        })]
      })]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "model-specific-parameters",
      children: "Model-Specific Parameters"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["Each model can also be configured with specific parameters using the prefix notation ", (0,jsx_runtime.jsx)(_components.code, {
        children: "model_name.parameter"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h4, {
      id: "openai",
      children: "OpenAI"
    }), "\n", (0,jsx_runtime.jsxs)(_components.table, {
      children: [(0,jsx_runtime.jsx)(_components.thead, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.th, {
            children: "Parameter"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Type"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Description"
          })]
        })
      }), (0,jsx_runtime.jsx)(_components.tbody, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "openai.model"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "string"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Specific OpenAI model to use"
          })]
        })
      })]
    }), "\n", (0,jsx_runtime.jsx)(_components.h4, {
      id: "claude",
      children: "Claude"
    }), "\n", (0,jsx_runtime.jsxs)(_components.table, {
      children: [(0,jsx_runtime.jsx)(_components.thead, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.th, {
            children: "Parameter"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Type"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Description"
          })]
        })
      }), (0,jsx_runtime.jsx)(_components.tbody, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "claude.model"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "string"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Specific Claude model to use"
          })]
        })
      })]
    }), "\n", (0,jsx_runtime.jsx)(_components.h4, {
      id: "gemini",
      children: "Gemini"
    }), "\n", (0,jsx_runtime.jsxs)(_components.table, {
      children: [(0,jsx_runtime.jsx)(_components.thead, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.th, {
            children: "Parameter"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Type"
          }), (0,jsx_runtime.jsx)(_components.th, {
            children: "Description"
          })]
        })
      }), (0,jsx_runtime.jsx)(_components.tbody, {
        children: (0,jsx_runtime.jsxs)(_components.tr, {
          children: [(0,jsx_runtime.jsx)(_components.td, {
            children: "gemini.model"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "string"
          }), (0,jsx_runtime.jsx)(_components.td, {
            children: "Specific Gemini model to use"
          })]
        })
      })]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "model-specific-default-values",
      children: "Model-Specific Default Values"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "openai-default-settings",
      children: "OpenAI Default Settings"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Default Model"
        }), ": ", (0,jsx_runtime.jsx)(_components.code, {
          children: "gpt-3.5-turbo"
        })]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Supported Models"
        }), ":", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "gpt-3.5-turbo"
            }), " - Fast and cost-effective for most tasks"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "gpt-4"
            }), " - Strong reasoning and instruction following"]
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "claude-default-settings",
      children: "Claude Default Settings"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Default Model"
        }), ": ", (0,jsx_runtime.jsx)(_components.code, {
          children: "claude-3-sonnet-20240229"
        })]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Supported Models"
        }), ":", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "claude-3-opus-20240229"
            }), " - Most powerful Claude model for complex tasks"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "claude-3-sonnet-20240229"
            }), " - Balanced performance and speed"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "claude-3-haiku-20240307"
            }), " - Fast and economical"]
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "gemini-default-settings",
      children: "Gemini Default Settings"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Default Model"
        }), ": ", (0,jsx_runtime.jsx)(_components.code, {
          children: "gemini-pro"
        })]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Supported Models"
        }), ":", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "gemini-pro"
            }), " - Original Gemini model"]
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "sdk-implementation",
      children: "SDK Implementation"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI uses official client SDKs for all supported AI models:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "OpenAI"
        }), ": Uses the official ", (0,jsx_runtime.jsx)(_components.code, {
          children: "github.com/sashabaranov/go-openai"
        }), " SDK"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Claude"
        }), ": Uses the official ", (0,jsx_runtime.jsx)(_components.code, {
          children: "github.com/anthropics/anthropic-sdk-go"
        }), " SDK"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Gemini"
        }), ": Uses the official ", (0,jsx_runtime.jsx)(_components.code, {
          children: "github.com/google/generative-ai-go"
        }), " SDK"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "configuration-methods",
      children: "Configuration Methods"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "You can configure model parameters in three ways, listed in order of precedence:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: (0,jsx_runtime.jsx)(_components.strong, {
          children: "Task-specific parameters in the config file"
        })
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: (0,jsx_runtime.jsx)(_components.strong, {
          children: "Environment variables"
        })
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: (0,jsx_runtime.jsx)(_components.strong, {
          children: "Default values"
        })
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "1-task-specific-configuration",
      children: "1. Task-specific Configuration"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["In the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "cronai.config"
      }), " file, you can specify model parameters using the prefix ", (0,jsx_runtime.jsx)(_components.code, {
        children: "model_params:"
      }), ":"]
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Format: timestamp model prompt response_processor [variables] [model_params:...]\n0 8 * * * claude product_manager file-output.txt model_params:temperature=0.8,model=claude-3-opus-20240229\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "You can also include both variables and model parameters:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "0 9 * * 1 openai report_template github-issue:owner/repo reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,model=gpt-4\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h4, {
      id: "using-model-specific-parameters",
      children: "Using Model-Specific Parameters"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "For model-specific configuration, use the prefix notation:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Use OpenAI-specific parameters\n0 9 * * 1 openai report_template file-output.txt model_params:openai.model=gpt-4\n\n# Use Claude-specific parameters\n0 8 * * * claude product_manager file-output.txt model_params:claude.model=claude-3-opus-20240229\n\n# Use Gemini-specific parameters\n*/15 * * * * gemini system_health file-output.txt model_params:gemini.model=gemini-pro\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "2-environment-variables",
      children: "2. Environment Variables"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "You can set global defaults for all tasks using environment variables:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "# Common parameters\nMODEL_TEMPERATURE=0.7\nMODEL_MAX_TOKENS=2048\n\n# Model-specific parameters\nOPENAI_MODEL=gpt-4\nCLAUDE_MODEL=claude-3-opus-20240229\nGEMINI_MODEL=gemini-pro\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "3-command-line-parameters",
      children: "3. Command Line Parameters"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["When using the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "run"
      }), " command, you can specify model parameters with the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "--model-params"
      }), " flag:"]
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "# Using common parameters\ncronai run --model openai --prompt weekly_report --processor file-output.txt --model-params \"temperature=0.5,max_tokens=4000,model=gpt-4\"\n\n# Using model-specific parameters\ncronai run --model gemini --prompt system_health --processor file-output.txt --model-params \"gemini.model=gemini-pro\"\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "examples",
      children: "Examples"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "low-temperature-for-consistent-output",
      children: "Low Temperature for Consistent Output"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Run system health check with very precise (low temperature) settings\n*/15 * * * * claude system_health file-health.log cluster=Primary model_params:temperature=0.1,max_tokens=1000\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "specific-model-version",
      children: "Specific Model Version"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Run weekly with OpenAI using a specific model\n0 9 * * 1 openai report_template file-report.log model_params:openai.model=gpt-4\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "github-processor-example",
      children: "GitHub Processor Example"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Create a GitHub issue with the weekly report\n0 9 * * 1 claude weekly_report github-owner/repo reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.7\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "console-output-example",
      children: "Console Output Example"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Output system health check to console (useful for testing)\n*/30 * * * * gemini system_health console model_params:temperature=0.3,max_tokens=2000\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "integration-with-variables",
      children: "Integration with Variables"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Model parameters can be used alongside variables:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "0 9 * * 1 openai report_template file-report.log reportType=Weekly,date={{CURRENT_DATE}} model_params:temperature=0.5,max_tokens=4000,openai.model=gpt-4\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "advanced-configuration",
      children: "Advanced Configuration"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "custom-api-endpoints",
      children: "Custom API Endpoints"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "For organizations using proxy services or custom endpoints for AI models, CronAI supports custom base URLs through environment variables:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "OpenAI"
        }), ": Set ", (0,jsx_runtime.jsx)(_components.code, {
          children: "OPENAI_BASE_URL"
        }), " environment variable to point to a custom endpoint"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Claude"
        }), ": Set ", (0,jsx_runtime.jsx)(_components.code, {
          children: "ANTHROPIC_BASE_URL"
        }), " environment variable to point to a custom endpoint"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "timeout-configuration",
      children: "Timeout Configuration"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "All model clients have a default timeout of 120 seconds (2 minutes) for API requests. This can be adjusted by setting the appropriate environment variables:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "OpenAI"
        }), ": Set ", (0,jsx_runtime.jsx)(_components.code, {
          children: "OPENAI_TIMEOUT"
        }), " to the desired timeout in seconds"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Claude"
        }), ": Set ", (0,jsx_runtime.jsx)(_components.code, {
          children: "ANTHROPIC_TIMEOUT"
        }), " to the desired timeout in seconds"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Gemini"
        }), ": Set ", (0,jsx_runtime.jsx)(_components.code, {
          children: "GEMINI_TIMEOUT"
        }), " to the desired timeout in seconds"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "post-mvp-features",
      children: "Post-MVP Features"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following features are planned for post-MVP releases:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Advanced parameters like top_p, frequency_penalty, and presence_penalty"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "System message customization"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Model fallback mechanism for automatic model switching on failure"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Safety setting configurations for Gemini"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Detailed error handling with retry mechanisms"
      }), "\n"]
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