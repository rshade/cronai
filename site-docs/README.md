# CronAI Documentation Site

This directory contains the Docusaurus documentation site for the CronAI project. The site is built with [Docusaurus](https://docusaurus.io/), a modern static website generator.

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/en/download/) version 16.14 or above
- [npm](https://www.npmjs.com/) (comes with Node.js)

### Installation

```bash
# Navigate to the site-docs directory
cd site-docs

# Install dependencies
npm install
```

### Local Development

```bash
# Start the development server
npm start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

### Build

```bash
# Build the website for production
npm run build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

### Serve Built Website

```bash
# Serve the built website locally
npm run serve
```

This command serves the built website locally, which is useful for testing the build output.

## Deployment

The documentation site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by a GitHub Actions workflow defined in `.github/workflows/deploy-docs.yml`.

If you want to manually deploy the site, you can use the following command:

```bash
# Deploy to GitHub Pages
npm run gh-pages
```

This command uses the `GIT_USER` environment variable defined in the script to authenticate with GitHub and deploy the site.

## Project Structure

```text
site-docs/
├── docs/            # Documentation files in markdown
├── src/             # React components and pages
│   ├── components/  # React components
│   ├── css/         # CSS files
│   └── pages/       # Custom React pages
├── static/          # Static files like images
│   └── img/         # Image files
├── docusaurus.config.js  # Docusaurus configuration
├── sidebars.js      # Sidebar configuration
└── package.json     # npm package configuration
```

## Adding Content

### Adding a New Document

1. Create a new markdown file in the `docs` directory
2. Add frontmatter at the top of the file:

   ```yaml
   ---
   id: unique-id
   title: Document Title
   sidebar_label: Sidebar Label
   ---
   ```

3. Add the document to the appropriate category in `sidebars.js`

### Adding a New Page

1. Create a new React component file in `src/pages`
2. The page will be available at the URL path based on the file name

## Customization

### Theme

The site's appearance can be customized in `src/css/custom.css`.

### Configuration

The site configuration is managed in `docusaurus.config.js`, where you can modify settings like:

- Site metadata (title, tagline, URL)
- Navigation bar structure
- Footer content
- Theme settings

For more information about configuring Docusaurus, refer to the [official documentation](https://docusaurus.io/docs/configuration).
