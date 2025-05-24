"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[891],{

/***/ 7697:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_prompt_management_md_60c_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-prompt-management-md-60c.json
const site_docs_prompt_management_md_60c_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"prompt-management","title":"Prompt Management","description":"CronAI includes a simple file-based prompt management system that helps you organize and use prompts efficiently.","source":"@site/docs/prompt-management.md","sourceDirName":".","slug":"/prompt-management","permalink":"/cronai/docs/prompt-management","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/prompt-management.md","tags":[],"version":"current","frontMatter":{"id":"prompt-management","title":"Prompt Management","sidebar_label":"Prompt Management"},"sidebar":"tutorialSidebar","previous":{"title":"systemd Service","permalink":"/cronai/docs/systemd"},"next":{"title":"Model Parameters","permalink":"/cronai/docs/model-parameters"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/prompt-management.md


const frontMatter = {
	id: 'prompt-management',
	title: 'Prompt Management',
	sidebar_label: 'Prompt Management'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "Directory Structure",
  "id": "directory-structure",
  "level": 2
}, {
  "value": "Prompt Files",
  "id": "prompt-files",
  "level": 2
}, {
  "value": "Variables in Prompts",
  "id": "variables-in-prompts",
  "level": 2
}, {
  "value": "CLI Commands",
  "id": "cli-commands",
  "level": 2
}, {
  "value": "List Prompts",
  "id": "list-prompts",
  "level": 3
}, {
  "value": "Search Prompts",
  "id": "search-prompts",
  "level": 3
}, {
  "value": "Show Prompt Details",
  "id": "show-prompt-details",
  "level": 3
}, {
  "value": "Preview Prompt",
  "id": "preview-prompt",
  "level": 3
}, {
  "value": "Using Prompts in CronAI Configuration",
  "id": "using-prompts-in-cronai-configuration",
  "level": 2
}, {
  "value": "Example Prompt Files",
  "id": "example-prompt-files",
  "level": 2
}, {
  "value": "Basic Prompt",
  "id": "basic-prompt",
  "level": 3
}, {
  "value": "Prompt with System Metrics",
  "id": "prompt-with-system-metrics",
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
    li: "li",
    p: "p",
    pre: "pre",
    ul: "ul",
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return (0,jsx_runtime.jsxs)(jsx_runtime.Fragment, {
    children: [(0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI includes a simple file-based prompt management system that helps you organize and use prompts efficiently."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "directory-structure",
      children: "Directory Structure"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["Prompts are stored as markdown files in the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "cron_prompts/"
      }), " directory:"]
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "cron_prompts/\n├── README.md           # Documentation for prompt structure\n├── monitoring/         # Prompts for monitoring purposes\n├── reports/            # Prompts for report generation\n├── system/             # Prompts for system operations\n└── [other_categories]/ # Custom categories\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "prompt-files",
      children: "Prompt Files"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["Prompts are standard markdown files with a ", (0,jsx_runtime.jsx)(_components.code, {
        children: ".md"
      }), " extension. During the MVP phase, prompts are simple text files that can contain variables:"]
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-markdown",
        children: "# System Health Check\n\nAnalyze the following system metrics and provide recommendations:\n\n- CPU Usage: {{cpu_usage}}%\n- Memory Usage: {{memory_usage}}%\n- Disk Usage: {{disk_usage}}%\n\nPlease provide:\n1. Assessment of current system health\n2. Potential issues identified\n3. Recommended actions\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "variables-in-prompts",
      children: "Variables in Prompts"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["Variables in prompts use the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "{{variable_name}}"
      }), " syntax. CronAI automatically provides the following special variables:"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "{{CURRENT_DATE}}"
        }), ": Current date in YYYY-MM-DD format"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "{{CURRENT_TIME}}"
        }), ": Current time in HH:MM", ":SS", " format"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "{{CURRENT_DATETIME}}"
        }), ": Current date and time in YYYY-MM-DD HH:MM", ":SS", " format"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Custom variables can be provided in the configuration file or command line using a comma-separated list of key=value pairs."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "cli-commands",
      children: "CLI Commands"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI provides several commands to help you manage your prompts:"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "list-prompts",
      children: "List Prompts"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "List all available prompts, optionally filtered by category:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "# List all prompts\ncronai prompt list\n\n# List prompts in a specific category\ncronai prompt list --category system\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "search-prompts",
      children: "Search Prompts"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Search for prompts by name or content:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "# Search by name or description\ncronai prompt search \"health check\"\n\n# Search in a specific category\ncronai prompt search --query \"monitoring\" --category system\n\n# Search in prompt content\ncronai prompt search --content --query \"CPU usage\"\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "show-prompt-details",
      children: "Show Prompt Details"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Show detailed information about a specific prompt:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "cronai prompt show system/system_health\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "preview-prompt",
      children: "Preview Prompt"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Preview a prompt with variables substituted:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "cronai prompt preview system/system_health --vars \"cpu_usage=85,memory_usage=70,disk_usage=50\"\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "using-prompts-in-cronai-configuration",
      children: "Using Prompts in CronAI Configuration"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Reference prompts in your cronai.config file using either the full path or category/name format:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Using a prompt from a category\n0 8 * * * openai system/system_health file-/var/log/cronai/health.log\n\n# Using a prompt with variables\n0 9 * * 1 claude reports/weekly_report github-issue:owner/repo date={{CURRENT_DATE}},team=Engineering\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "example-prompt-files",
      children: "Example Prompt Files"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "basic-prompt",
      children: "Basic Prompt"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-markdown",
        children: "# Daily Status Report\n\nGenerate a daily status report for {{project}} on {{CURRENT_DATE}}.\n\n1. Current status overview\n2. Progress since yesterday\n3. Planned tasks for today\n4. Blocking issues, if any\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "prompt-with-system-metrics",
      children: "Prompt with System Metrics"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-markdown",
        children: "# System Health Check\n\nCPU Usage: {{cpu_usage}}%\nMemory Usage: {{memory_usage}}%\nDisk Space: {{disk_usage}}%\n\nPlease analyze these metrics and provide:\n1. Current system health assessment\n2. Potential issues or warning signs\n3. Recommended actions\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "post-mvp-features",
      children: "Post-MVP Features"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following prompt management features are planned for future releases:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Prompt metadata (YAML frontmatter)"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Template inheritance and composition"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Includes for reusing common prompt components"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Conditional logic in prompts"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Advanced variable validation and defaults"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "For more information on these upcoming features, see the project roadmap."
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