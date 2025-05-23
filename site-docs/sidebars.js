/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  tutorialSidebar: [
    {
      type: 'category',
      label: 'Introduction',
      items: ['intro', 'architecture', 'limitations-and-improvements'],
    },
    {
      type: 'category',
      label: 'Guides',
      items: [
        'systemd',
        'prompt-management',
        'model-parameters',
        'logging',
        'troubleshooting',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: ['api'],
    },
  ],
};

module.exports = sidebars;