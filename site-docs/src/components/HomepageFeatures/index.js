import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Schedule AI Prompts',
    Svg: require('@site/static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        CronAI allows you to schedule AI model prompts to run on a cron-type schedule,
        automating your AI workflows and keeping you updated with regular insights.
      </>
    ),
  },
  {
    title: 'Multiple Model Support',
    Svg: require('@site/static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Work with multiple AI models including OpenAI, Claude, and Gemini.
        Configure model-specific parameters and fallbacks for reliability.
      </>
    ),
  },
  {
    title: 'Flexible Output Processing',
    Svg: require('@site/static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
        Process AI responses through various channels including email, Slack,
        webhooks, GitHub, and file output with templated formatting.
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}