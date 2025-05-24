"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[174],{

/***/ 7106:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  assets: () => (/* binding */ assets),
  contentTitle: () => (/* binding */ contentTitle),
  "default": () => (/* binding */ MDXContent),
  frontMatter: () => (/* binding */ frontMatter),
  metadata: () => (/* reexport */ site_docs_systemd_md_e3a_namespaceObject),
  toc: () => (/* binding */ toc)
});

;// ./.docusaurus/docusaurus-plugin-content-docs/default/site-docs-systemd-md-e3a.json
const site_docs_systemd_md_e3a_namespaceObject = /*#__PURE__*/JSON.parse('{"id":"systemd","title":"Running CronAI as a systemd Service","description":"This document explains how to set up CronAI to run as a systemd service on Linux systems.","source":"@site/docs/systemd.md","sourceDirName":".","slug":"/systemd","permalink":"/cronai/docs/systemd","draft":false,"unlisted":false,"editUrl":"https://github.com/rshade/cronai/tree/main/site-docs/docs/systemd.md","tags":[],"version":"current","frontMatter":{"id":"systemd","title":"Running CronAI as a systemd Service","sidebar_label":"systemd Service"},"sidebar":"tutorialSidebar","previous":{"title":"CronAI: Known Limitations and Future Improvements","permalink":"/cronai/docs/limitations-and-improvements"},"next":{"title":"Prompt Management","permalink":"/cronai/docs/prompt-management"}}');
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
// EXTERNAL MODULE: ./node_modules/@mdx-js/react/lib/index.js
var lib = __webpack_require__(8453);
;// ./docs/systemd.md


const frontMatter = {
	id: 'systemd',
	title: 'Running CronAI as a systemd Service',
	sidebar_label: 'systemd Service'
};
const contentTitle = undefined;

const assets = {

};



const toc = [{
  "value": "Setup Steps",
  "id": "setup-steps",
  "level": 2
}, {
  "value": "1. Build and install the CronAI binary",
  "id": "1-build-and-install-the-cronai-binary",
  "level": 3
}, {
  "value": "2. Create your configuration and prompt files",
  "id": "2-create-your-configuration-and-prompt-files",
  "level": 3
}, {
  "value": "3. Set up your environment file",
  "id": "3-set-up-your-environment-file",
  "level": 3
}, {
  "value": "4. Create the systemd service file",
  "id": "4-create-the-systemd-service-file",
  "level": 3
}, {
  "value": "5. Enable and start the service",
  "id": "5-enable-and-start-the-service",
  "level": 3
}, {
  "value": "6. Check the service status",
  "id": "6-check-the-service-status",
  "level": 3
}, {
  "value": "7. View the logs",
  "id": "7-view-the-logs",
  "level": 3
}, {
  "value": "Managing the Service",
  "id": "managing-the-service",
  "level": 2
}, {
  "value": "Restart the service",
  "id": "restart-the-service",
  "level": 3
}, {
  "value": "Stop the service",
  "id": "stop-the-service",
  "level": 3
}, {
  "value": "Disable the service (prevents it from starting at boot)",
  "id": "disable-the-service-prevents-it-from-starting-at-boot",
  "level": 3
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
      children: "This document explains how to set up CronAI to run as a systemd service on Linux systems."
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "setup-steps",
      children: "Setup Steps"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "1-build-and-install-the-cronai-binary",
      children: "1. Build and install the CronAI binary"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "cd /path/to/cronai\ngo build -o cronai ./cmd/cronai\nsudo cp cronai /usr/local/bin/cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "2-create-your-configuration-and-prompt-files",
      children: "2. Create your configuration and prompt files"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "mkdir -p /etc/cronai/cron_prompts\ncp cronai.config.example /etc/cronai/cronai.config\ncp -r cron_prompts/* /etc/cronai/cron_prompts/\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "3-set-up-your-environment-file",
      children: "3. Set up your environment file"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "cp .env.example /etc/cronai/.env\n# Edit the .env file with your API keys and settings\nsudo nano /etc/cronai/.env\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "4-create-the-systemd-service-file",
      children: "4. Create the systemd service file"
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Copy the example service file and modify it for your system:"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo cp cronai.service /etc/systemd/system/cronai.service\nsudo nano /etc/systemd/system/cronai.service\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.p, {
      children: "Update the following fields in the service file:"
    }), "\n", (0,jsx_runtime.jsxs)(_components.ul, {
      children: ["\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "User"
        }), ": The user account that will run the service"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "WorkingDirectory"
        }), ": The directory where your configuration is located (e.g., ", (0,jsx_runtime.jsx)(_components.code, {
          children: "/etc/cronai"
        }), ")"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "ExecStart"
        }), ": The path to the CronAI binary (e.g., ", (0,jsx_runtime.jsx)(_components.code, {
          children: "/usr/local/bin/cronai start --config /etc/cronai/cronai.config"
        }), ")"]
      }), "\n", (0,jsx_runtime.jsxs)(_components.li, {
        children: [(0,jsx_runtime.jsx)(_components.code, {
          children: "EnvironmentFile"
        }), ": The path to your .env file (e.g., ", (0,jsx_runtime.jsx)(_components.code, {
          children: "/etc/cronai/.env"
        }), ")"]
      }), "\n"]
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "5-enable-and-start-the-service",
      children: "5. Enable and start the service"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo systemctl daemon-reload\nsudo systemctl enable cronai\nsudo systemctl start cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "6-check-the-service-status",
      children: "6. Check the service status"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo systemctl status cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "7-view-the-logs",
      children: "7. View the logs"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo journalctl -u cronai -f\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h2, {
      id: "managing-the-service",
      children: "Managing the Service"
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "restart-the-service",
      children: "Restart the service"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo systemctl restart cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "stop-the-service",
      children: "Stop the service"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo systemctl stop cronai\n"
      })
    }), "\n", (0,jsx_runtime.jsx)(_components.h3, {
      id: "disable-the-service-prevents-it-from-starting-at-boot",
      children: "Disable the service (prevents it from starting at boot)"
    }), "\n", (0,jsx_runtime.jsx)(_components.pre, {
      children: (0,jsx_runtime.jsx)(_components.code, {
        className: "language-bash",
        children: "sudo systemctl disable cronai\n"
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