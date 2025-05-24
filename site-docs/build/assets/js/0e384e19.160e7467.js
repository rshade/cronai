"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[976],{

/***/ 7879:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_intro_md_0e3_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-intro-md-0e3.json
const site_docs_intro_md_0e3_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"intro","title":"CronAI","description":"AI agent for scheduled prompt execution - Your automated AI assistant.","source":"@site/docs/intro.md","sourceDirName":".","slug":"/","permalink":"/cronai/docs/","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/intro.md","tags":[],"version":"current","frontMatter":{"id":"intro","title":"CronAI","sidebar_label":"Introduction","slug":"/"},"sidebar":"tutorialSidebar","next":{"title":"Architecture","permalink":"/cronai/docs/architecture"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/intro.md


const frontMatter = {
	id: 'intro',
	title: 'CronAI',
	sidebar_label: 'Introduction',
	slug: '/'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "Overview",
  "id": "overview",
  "level": 2
}, {
  "value": "MVP Features",
  "id": "mvp-features",
  "level": 2
}, {
  "value": "Planned Features (Not Yet Implemented)",
  "id": "planned-features-not-yet-implemented",
  "level": 3
}, {
  "value": "Installation",
  "id": "installation",
  "level": 2
}, {
  "value": "Configuration",
  "id": "configuration",
  "level": 2
}, {
  "value": "Format",
  "id": "format",
  "level": 3
}, {
  "value": "Example Configuration",
  "id": "example-configuration",
  "level": 3
}];
function _createMdxContent(props) {
  const _components = {
    a: "a",
    code: "code",
    h2: "h2",
    h3: "h3",
    li: "li",
    p: "p",
    pre: "pre",
    strong: "strong",
    ul: "ul",
    ...(0,lib/* useMDXComponents */.R)(),
    ...props.components
  };
  return (0,jsx_runtime.jsxs)(jsx_runtime.Fragment, {
    children: [(0,jsx_runtime.jsx)(_components.p, {
      children: "AI agent for scheduled prompt execution - Your automated AI assistant."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "overview",
      children: "Overview"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "CronAI is an intelligent agent that schedules and executes AI model prompts automatically. It acts as your personal AI automation system, running tasks on schedule and delivering results through your preferred channels."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "mvp-features",
      children: "MVP Features"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The current MVP release includes:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "✅ Cron-style scheduling for automated execution"
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["✅ Support for multiple AI models:", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "OpenAI (gpt-3.5-turbo, gpt-4)"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Claude (claude-3-sonnet, claude-3-opus)"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Gemini"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "✅ Customizable prompts stored as markdown files"
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: ["✅ Implemented response processors:", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "File output - Save responses to files"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "GitHub integration - Create issues and add comments"
          }), "\n", (0,jsx_runtime.jsx)(_components.li, {
            children: "Console output - Display responses in terminal"
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "✅ Variable substitution in prompts"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "✅ Systemd service for deployment"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "planned-features-not-yet-implemented",
      children: "Planned Features (Not Yet Implemented)"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "The following processors are planned but not yet implemented in the current release:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "⚠️ Email processor - Currently logs actions only"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "⚠️ Slack processor - Currently logs actions only"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "⚠️ Webhook processor - Currently logs actions only"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Additional planned features:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Enhanced templating capabilities"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Web UI for prompt management"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Model fallback mechanisms"
      }), "\n", (0,jsx_runtime.jsx)(_components.li, {
        children: "Advanced scheduling options"
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["See ", (0,jsx_runtime.jsx)(_components.a, {
        href: "https://github.com/rshade/cronai/blob/main/docs/limitations-and-improvements.md",
        children: "Limitations and Improvements"
      }), " for a detailed breakdown of current limitations and planned improvements."]
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "installation",
      children: "Installation"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "# Install directly\ngo install github.com/rshade/cronai/cmd/cronai@latest\n\n# Or clone and build\ngit clone https://github.com/rshade/cronai.git\ncd cronai\ngo build -o cronai ./cmd/cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "configuration",
      children: "Configuration"
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["Create a configuration file called ", (0,jsx_runtime.jsx)(_components.code, {
        children: "cronai.config"
      }), " with your scheduled tasks."]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "format",
      children: "Format"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "timestamp model prompt response_processor [variables] [model_params:...]\n"
      })
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "timestamp"
        }), ": Standard cron format (minute hour day-of-month month day-of-week)"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "model"
        }), ": AI model to use (openai, claude, gemini)"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "prompt"
        }), ": Name of prompt file in cron_prompts directory (with or without .md extension)"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "response_processor"
        }), ": How to process the response:", "\n", (0,jsx_runtime.jsxs)(_components.ul, {
          children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "file-path/to/output.txt"
            }), ": Save to file"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "github-issue:owner/repo"
            }), ": Create GitHub issue"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "github-comment:owner/repo#123"
            }), ": Add comment to GitHub issue"]
          }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
            children: [(0,jsx_runtime.jsx)(_components.code, {
              children: "console"
            }), ": Display in console"]
          }), "\n"]
        }), "\n"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "variables"
        }), " (optional): Variables to replace in the prompt file, in the format ", (0,jsx_runtime.jsx)(_components.code, {
          children: "key1=value1,key2=value2,..."
        })]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.strong, {
          children: "model_params"
        }), " (optional): Model-specific parameters in the format ", (0,jsx_runtime.jsx)(_components.code, {
          children: "model_params:param1=value1,param2=value2,..."
        })]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "example-configuration",
      children: "Example Configuration"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-text",
        children: "# Run daily at 8 AM using OpenAI, saving to file\n0 8 * * * openai product_manager file-/var/log/cronai/product_manager.log\n\n# Run weekly on Monday at 9 AM using Claude, creating GitHub issue\n0 9 * * 1 claude weekly_report github-issue:your-org/your-repo\n\n# Run daily health check with variables\n0 6 * * * openai system_check file-/var/log/cronai/health.log system=production,check_level=detailed\n\n# Run with custom model parameters (temperature and specific model version)\n0 9 * * 1 openai weekly_report file-/var/log/cronai/report.log model_params:temperature=0.5,model=gpt-4\n"
      })
    }), "\n", (0,jsx_runtime.jsxs)(_components.p, {
      children: ["See the ", (0,jsx_runtime.jsx)(_components.a, {
        href: "https://github.com/rshade/cronai/blob/main/cronai.config.example",
        children: "Example Configuration Files"
      }), " in the repository for more examples."]
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