"use strict";
(self["webpackChunkcronai_docs"] = self["webpackChunkcronai_docs"] || []).push([[634],{

/***/ 5:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  "default": () => (/* binding */ Home)
});

// EXTERNAL MODULE: ./node_modules/react/index.js
var react = __webpack_require__(6540);
// EXTERNAL MODULE: ./node_modules/clsx/dist/clsx.mjs
var clsx = __webpack_require__(4164);
// EXTERNAL MODULE: ./node_modules/@docusaurus/core/lib/client/exports/Link.js
var Link = __webpack_require__(6289);
// EXTERNAL MODULE: ./node_modules/@docusaurus/core/lib/client/exports/useDocusaurusContext.js
var useDocusaurusContext = __webpack_require__(797);
// EXTERNAL MODULE: ./node_modules/@docusaurus/theme-classic/lib/theme/Layout/index.js + 67 modules
var Layout = __webpack_require__(1410);
;// ./src/components/HomepageFeatures/styles.module.css
// extracted by mini-css-extract-plugin
/* harmony default export */ const styles_module = ({"features":"features_t9lD","featureSvg":"featureSvg_GfXr"});
// EXTERNAL MODULE: ./node_modules/react/jsx-runtime.js
var jsx_runtime = __webpack_require__(4848);
;// ./src/components/HomepageFeatures/index.js
const FeatureList=[{title:'Schedule AI Prompts',Svg:(__webpack_require__(9132)/* ["default"] */ .A),description:/*#__PURE__*/(0,jsx_runtime.jsx)(jsx_runtime.Fragment,{children:"CronAI allows you to schedule AI model prompts to run on a cron-type schedule, automating your AI workflows and keeping you updated with regular insights."})},{title:'Multiple Model Support',Svg:(__webpack_require__(7709)/* ["default"] */ .A),description:/*#__PURE__*/(0,jsx_runtime.jsx)(jsx_runtime.Fragment,{children:"Work with multiple AI models including OpenAI, Claude, and Gemini. Configure model-specific parameters and fallbacks for reliability."})},{title:'Flexible Output Processing',Svg:(__webpack_require__(9066)/* ["default"] */ .A),description:/*#__PURE__*/(0,jsx_runtime.jsx)(jsx_runtime.Fragment,{children:"Process AI responses through various channels including email, Slack, webhooks, GitHub, and file output with templated formatting."})}];function Feature({Svg,title,description}){return/*#__PURE__*/(0,jsx_runtime.jsxs)("div",{className:(0,clsx/* default */.A)('col col--4'),children:[/*#__PURE__*/(0,jsx_runtime.jsx)("div",{className:"text--center",children:/*#__PURE__*/(0,jsx_runtime.jsx)(Svg,{className:styles_module.featureSvg,role:"img"})}),/*#__PURE__*/(0,jsx_runtime.jsxs)("div",{className:"text--center padding-horiz--md",children:[/*#__PURE__*/(0,jsx_runtime.jsx)("h3",{children:title}),/*#__PURE__*/(0,jsx_runtime.jsx)("p",{children:description})]})]});}function HomepageFeatures(){return/*#__PURE__*/(0,jsx_runtime.jsx)("section",{className:styles_module.features,children:/*#__PURE__*/(0,jsx_runtime.jsx)("div",{className:"container",children:/*#__PURE__*/(0,jsx_runtime.jsx)("div",{className:"row",children:FeatureList.map((props,idx)=>/*#__PURE__*/(0,jsx_runtime.jsx)(Feature,{...props},idx))})})});}
;// ./src/pages/index.module.css
// extracted by mini-css-extract-plugin
/* harmony default export */ const index_module = ({"heroBanner":"heroBanner_qdFl","buttons":"buttons_AeoN"});
;// ./src/pages/index.js
function HomepageHeader(){const{siteConfig}=(0,useDocusaurusContext/* default */.A)();return/*#__PURE__*/(0,jsx_runtime.jsx)("header",{className:(0,clsx/* default */.A)('hero hero--primary',index_module.heroBanner),children:/*#__PURE__*/(0,jsx_runtime.jsxs)("div",{className:"container",children:[/*#__PURE__*/(0,jsx_runtime.jsx)("h1",{className:"hero__title",children:siteConfig.title}),/*#__PURE__*/(0,jsx_runtime.jsx)("p",{className:"hero__subtitle",children:siteConfig.tagline}),/*#__PURE__*/(0,jsx_runtime.jsx)("div",{className:index_module.buttons,children:/*#__PURE__*/(0,jsx_runtime.jsx)(Link/* default */.A,{className:"button button--secondary button--lg",to:"/docs/intro",children:"Get Started"})})]})});}function Home(){const{siteConfig}=(0,useDocusaurusContext/* default */.A)();return/*#__PURE__*/(0,jsx_runtime.jsxs)(Layout/* default */.A,{title:`${siteConfig.title}`,description:"Schedule AI prompts with cron-like precision",children:[/*#__PURE__*/(0,jsx_runtime.jsx)(HomepageHeader,{}),/*#__PURE__*/(0,jsx_runtime.jsx)("main",{children:/*#__PURE__*/(0,jsx_runtime.jsx)(HomepageFeatures,{})})]});}

/***/ }),

/***/ 7709:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   A: () => (__WEBPACK_DEFAULT_EXPORT__)
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(6540);
var _circle, _path, _circle2;
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }

const SvgUndrawDocusaurusTree = ({
  title,
  titleId,
  ...props
}) => /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("svg", _extends({
  xmlns: "http://www.w3.org/2000/svg",
  width: 200,
  height: 200,
  viewBox: "0 0 200 200",
  "aria-labelledby": titleId
}, props), title ? /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("title", {
  id: titleId
}, title) : null, _circle || (_circle = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("circle", {
  cx: 100,
  cy: 100,
  r: 80,
  fill: "#2e8555"
})), _path || (_path = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("path", {
  fill: "#fff",
  d: "M90 100h20v60H90z"
})), _circle2 || (_circle2 = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("circle", {
  cx: 100,
  cy: 80,
  r: 40,
  fill: "#fff"
})));
/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (SvgUndrawDocusaurusTree);

/***/ }),

/***/ 9066:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   A: () => (__WEBPACK_DEFAULT_EXPORT__)
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(6540);
var _circle, _circle2, _ellipse, _ellipse2;
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }

const SvgUndrawDocusaurusReact = ({
  title,
  titleId,
  ...props
}) => /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("svg", _extends({
  xmlns: "http://www.w3.org/2000/svg",
  width: 200,
  height: 200,
  viewBox: "0 0 200 200",
  "aria-labelledby": titleId
}, props), title ? /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("title", {
  id: titleId
}, title) : null, _circle || (_circle = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("circle", {
  cx: 100,
  cy: 100,
  r: 80,
  fill: "#2e8555"
})), _circle2 || (_circle2 = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("circle", {
  cx: 100,
  cy: 100,
  r: 20,
  fill: "#fff"
})), _ellipse || (_ellipse = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("ellipse", {
  cx: 100,
  cy: 100,
  fill: "none",
  stroke: "#fff",
  strokeWidth: 4,
  rx: 60,
  ry: 20
})), _ellipse2 || (_ellipse2 = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("ellipse", {
  cx: 100,
  cy: 100,
  fill: "none",
  stroke: "#fff",
  strokeWidth: 4,
  rx: 20,
  ry: 60
})));
/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (SvgUndrawDocusaurusReact);

/***/ }),

/***/ 9132:
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   A: () => (__WEBPACK_DEFAULT_EXPORT__)
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(6540);
var _circle, _path;
function _extends() { return _extends = Object.assign ? Object.assign.bind() : function (n) { for (var e = 1; e < arguments.length; e++) { var t = arguments[e]; for (var r in t) ({}).hasOwnProperty.call(t, r) && (n[r] = t[r]); } return n; }, _extends.apply(null, arguments); }

const SvgUndrawDocusaurusMountain = ({
  title,
  titleId,
  ...props
}) => /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("svg", _extends({
  xmlns: "http://www.w3.org/2000/svg",
  width: 200,
  height: 200,
  viewBox: "0 0 200 200",
  "aria-labelledby": titleId
}, props), title ? /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("title", {
  id: titleId
}, title) : null, _circle || (_circle = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("circle", {
  cx: 100,
  cy: 100,
  r: 80,
  fill: "#2e8555"
})), _path || (_path = /*#__PURE__*/react__WEBPACK_IMPORTED_MODULE_0__.createElement("path", {
  fill: "#fff",
  d: "m70 140 30-80 30 80Z"
})));
/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (SvgUndrawDocusaurusMountain);

/***/ })

}]);