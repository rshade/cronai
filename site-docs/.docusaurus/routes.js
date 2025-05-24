import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/cronai/__docusaurus/debug',
    component: ComponentCreator('/cronai/__docusaurus/debug', 'c0e'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/config',
    component: ComponentCreator('/cronai/__docusaurus/debug/config', 'a47'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/content',
    component: ComponentCreator('/cronai/__docusaurus/debug/content', '6c7'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/globalData',
    component: ComponentCreator('/cronai/__docusaurus/debug/globalData', 'b01'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/metadata',
    component: ComponentCreator('/cronai/__docusaurus/debug/metadata', 'e89'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/registry',
    component: ComponentCreator('/cronai/__docusaurus/debug/registry', '19b'),
    exact: true
  },
  {
    path: '/cronai/__docusaurus/debug/routes',
    component: ComponentCreator('/cronai/__docusaurus/debug/routes', '52b'),
    exact: true
  },
  {
    path: '/cronai/docs',
    component: ComponentCreator('/cronai/docs', '40c'),
    routes: [
      {
        path: '/cronai/docs',
        component: ComponentCreator('/cronai/docs', 'd1a'),
        routes: [
          {
            path: '/cronai/docs',
            component: ComponentCreator('/cronai/docs', '5bd'),
            routes: [
              {
                path: '/cronai/docs',
                component: ComponentCreator('/cronai/docs', 'a7b'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/api',
                component: ComponentCreator('/cronai/docs/api', '736'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/architecture',
                component: ComponentCreator('/cronai/docs/architecture', '6ce'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/limitations-and-improvements',
                component: ComponentCreator('/cronai/docs/limitations-and-improvements', '348'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/logging',
                component: ComponentCreator('/cronai/docs/logging', '591'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/model-parameters',
                component: ComponentCreator('/cronai/docs/model-parameters', 'f00'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/prompt-management',
                component: ComponentCreator('/cronai/docs/prompt-management', 'd70'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/systemd',
                component: ComponentCreator('/cronai/docs/systemd', '630'),
                exact: true,
                sidebar: "tutorialSidebar"
              },
              {
                path: '/cronai/docs/troubleshooting',
                component: ComponentCreator('/cronai/docs/troubleshooting', 'bb1'),
                exact: true,
                sidebar: "tutorialSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/cronai/',
    component: ComponentCreator('/cronai/', '4a7'),
    exact: true
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
