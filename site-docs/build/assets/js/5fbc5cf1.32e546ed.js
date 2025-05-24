"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[624],{

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


/***/ }),

/***/ 9362:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_api_md_5fb_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-api-md-5fb.json
const site_docs_api_md_5fb_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"api","title":"CronAI API Documentation","description":"Note: This document is a placeholder for future API documentation. The external API is planned for post-MVP releases.","source":"@site/docs/api.md","sourceDirName":".","slug":"/api","permalink":"/cronai/docs/api","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/api.md","tags":[],"version":"current","frontMatter":{"id":"api","title":"CronAI API Documentation","sidebar_label":"API Reference"},"sidebar":"tutorialSidebar","previous":{"title":"Troubleshooting Guide (Coming Soon)","permalink":"/cronai/docs/troubleshooting"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/api.md


const frontMatter = {
	id: 'api',
	title: 'CronAI API Documentation',
	sidebar_label: 'API Reference'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "Current Status",
  "id": "current-status",
  "level": 2
}, {
  "value": "Planned API Features",
  "id": "planned-api-features",
  "level": 2
}, {
  "value": "API Design Principles",
  "id": "api-design-principles",
  "level": 2
}, {
  "value": "Planned Endpoints",
  "id": "planned-endpoints",
  "level": 2
}, {
  "value": "Authentication",
  "id": "authentication",
  "level": 2
}, {
  "value": "Status Codes",
  "id": "status-codes",
  "level": 2
}, {
  "value": "Stay Tuned",
  "id": "stay-tuned",
  "level": 2
}];
function _createMdxContent(props) {
  const _components = {
    blockquote: "blockquote",
    code: "code",
    h2: "h2",
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
    children: [(0,jsx_runtime.jsxs)(_components.blockquote, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.p, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "Note"
        }), ": This document is a placeholder for future API documentation. The external API is planned for post-MVP releases."]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "This document will describe the API endpoints, request/response formats, and authentication mechanisms for CronAI's external API."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "current-status",
      children: "Current Status"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The external API is currently in development and planned for post-MVP releases. The MVP version does not include external API access."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "planned-api-features",
      children: "Planned API Features"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following API features are planned for future releases:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ol, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsx)(_components.p, {
          children: (0,jsx_runtime.jsx)(_components.strong, {
            children: "Task Management API"
          })
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Create, read, update, and delete scheduled tasks"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Control task execution (start, stop, pause)"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Query task execution history and status"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsx)(_components.p, {
          children: (0,jsx_runtime.jsx)(_components.strong, {
            children: "Prompt Management API"
          })
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Create, read, update, and delete prompts"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Organize prompts into categories"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Test prompts with variables"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsx)(_components.p, {
          children: (0,jsx_runtime.jsx)(_components.strong, {
            children: "Model Management API"
          })
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Configure model parameters"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "View model usage and statistics"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Test model execution"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["\n", (0,jsx_runtime.jsx)(_components.p, {
          children: (0,jsx_runtime.jsx)(_components.strong, {
            children: "Response Processing API"
          })
        }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Configure response processors"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "View response history"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Query processed outputs"
          }), "\n"]
        }), "\n"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "api-design-principles",
      children: "API Design Principles"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The future API will follow these design principles:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "RESTful architecture"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "JSON request/response format"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Token-based authentication"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Comprehensive error responses"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Versioned endpoints"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Rate limiting"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Pagination for list endpoints"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "planned-endpoints",
      children: "Planned Endpoints"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "A preview of planned endpoints:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-http",
        children: "# Task Management\nGET    /api/v1/tasks\nPOST   /api/v1/tasks\nGET    /api/v1/tasks/:id\nPUT    /api/v1/tasks/:id\nDELETE /api/v1/tasks/:id\nPOST   /api/v1/tasks/:id/execute\n\n# Prompt Management\nGET    /api/v1/prompts\nPOST   /api/v1/prompts\nGET    /api/v1/prompts/:id\nPUT    /api/v1/prompts/:id\nDELETE /api/v1/prompts/:id\nPOST   /api/v1/prompts/:id/test\n\n# Model Management\nGET    /api/v1/models\nGET    /api/v1/models/:id/config\nPUT    /api/v1/models/:id/config\nPOST   /api/v1/models/:id/test\n\n# Response Processing\nGET    /api/v1/processors\nGET    /api/v1/processors/:type/config\nPUT    /api/v1/processors/:type/config\nGET    /api/v1/responses\nGET    /api/v1/responses/:id\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "authentication",
      children: "Authentication"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Future API authentication will likely use:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Bearer token authentication"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "OAuth 2.0 integration (for third-party applications)"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Role-based access control"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "status-codes",
      children: "Status Codes"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The API will use standard HTTP status codes:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "200 OK"
        }), ": Successful request"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "201 Created"
        }), ": Resource created successfully"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "400 Bad Request"
        }), ": Invalid request parameters"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "401 Unauthorized"
        }), ": Missing or invalid authentication"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "403 Forbidden"
        }), ": Insufficient permissions"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "404 Not Found"
        }), ": Resource not found"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "500 Internal Server Error"
        }), ": Server-side error"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "stay-tuned",
      children: "Stay Tuned"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The API documentation will be expanded as the external API is implemented. Check back in future releases for comprehensive API documentation."
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



/***/ })

}]);