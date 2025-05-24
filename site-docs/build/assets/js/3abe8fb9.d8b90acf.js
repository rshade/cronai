"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[485],{

/***/ 6180:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_logging_md_3ab_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-logging-md-3ab.json
const site_docs_logging_md_3ab_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"logging","title":"Logging in CronAI","description":"CronAI implements structured logging with configurable log levels to help with troubleshooting and monitoring.","source":"@site/docs/logging.md","sourceDirName":".","slug":"/logging","permalink":"/cronai/docs/logging","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/logging.md","tags":[],"version":"current","frontMatter":{},"sidebar":"tutorialSidebar","previous":{"title":"Model Parameters","permalink":"/cronai/docs/model-parameters"},"next":{"title":"Troubleshooting Guide (Coming Soon)","permalink":"/cronai/docs/troubleshooting"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/logging.md


const frontMatter = {};
const contentTitle = 'Logging in CronAI';

const assets = {

};



const toc = [{
  "value": "Log Levels",
  "id": "log-levels",
  "level": 2
}, {
  "value": "Configuring Log Level",
  "id": "configuring-log-level",
  "level": 2
}];
function _createMdxContent(props) {
  const _components = {
    code: "code",
    h1: "h1",
    h2: "h2",
    header: "header",
    li: "li",
    p: "p",
    pre: "pre",
    strong: "strong",
    ul: "ul",
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return (0,jsx_runtime.jsxs)(jsx_runtime.Fragment, {
    children: [(0,jsx_runtime.jsx)(_components.header, {
      children: (0,jsx_runtime.jsx)(_components.h1, {
        id: "logging-in-cronai",
        children: "Logging in CronAI"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI implements structured logging with configurable log levels to help with troubleshooting and monitoring."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "log-levels",
      children: "Log Levels"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following log levels are supported, in order of increasing severity:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "DEBUG"
        }), ": Detailed information, typically only useful when troubleshooting issues"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "INFO"
        }), ": General information about the normal operation of the application"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "WARN"
        }), ": Warnings that don't affect application function but should be addressed"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "ERROR"
        }), ": Errors that affect application function but don't cause termination"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "FATAL"
        }), ": Fatal errors that require application termination"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "configuring-log-level",
      children: "Configuring Log Level"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["The log level can be configured through the ", (0,jsx_runtime.jsx)(_components.code, {
        children: "LOG_LEVEL"
      }), " environment variable:"]
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "export LOG_LEVEL=DEBUG\n./cronai start\n```text\n\nValid values for `LOG_LEVEL` are: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`.\n\nIf not specified, the default log level is `INFO`.\n\n## Structured Logging\n\nCronAI uses structured logging to provide context for log messages. Each log message includes:\n\n- Timestamp in RFC3339 format\n- Log level\n- Message\n- File and line number (for debugging)\n- Additional metadata specific to the log message\n\nExample log output:\n\n```text\n[2025-05-18T14:30:45Z] [INFO] (service.go:40) Starting CronAI service | config_path=/etc/cronai/cronai.config\n[2025-05-18T14:30:45Z] [INFO] (service.go:207) Successfully parsed configuration file | path=/etc/cronai/cronai.config, task_count=3\n[2025-05-18T14:30:45Z] [INFO] (service.go:70) Scheduled task | task_index=0, schedule=0 9 * * 1-5, model=claude, prompt=weekly_report, processor=email-team@company.com\n```text\n\n## JSON Logging\n\nFor integration with log management systems, CronAI supports JSON-formatted logs. To enable JSON logging, set the `LOG_FORMAT` environment variable to `JSON`:\n\n```bash\nexport LOG_FORMAT=JSON\n./cronai start\n```text\n\nExample JSON log output:\n\n```json\n{\"time\":\"2025-05-18T14:30:45Z\",\"level\":\"INFO\",\"message\":\"Starting CronAI service\",\"file\":\"service.go\",\"line\":40,\"metadata\":{\"config_path\":\"/etc/cronai/cronai.config\"}}\n{\"time\":\"2025-05-18T14:30:45Z\",\"level\":\"INFO\",\"message\":\"Successfully parsed configuration file\",\"file\":\"service.go\",\"line\":207,\"metadata\":{\"path\":\"/etc/cronai/cronai.config\",\"task_count\":3}}\n```text\n\n## Error Handling\n\nCronAI implements categorized error handling through the `errors` package. Errors are categorized as:\n\n- **CONFIGURATION**: Errors related to configuration files and parameters\n- **VALIDATION**: Errors related to input validation\n- **EXTERNAL**: Errors from external services (APIs, etc.)\n- **SYSTEM**: System-level errors (file I/O, etc.)\n- **APPLICATION**: Application-level errors\n\nError logs include the error category and additional context information to aid in troubleshooting.\n\n## Log File\n\nBy default, logs are written to STDOUT. To direct logs to a file, use the `LOG_FILE` environment variable:\n\n```bash\nexport LOG_FILE=/var/log/cronai.log\n./cronai start\n```text\n\nIf not specified, logs are written to STDOUT.\n\n## Troubleshooting\n\nFor troubleshooting issues, set the log level to DEBUG:\n\n```bash\nexport LOG_LEVEL=DEBUG\n./cronai start\n```text\n\nThis will provide detailed logs of all operations, including prompt loading, model execution, and response processing.\n"
      })
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