// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Permify',
  url: 'https://permify.co',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'Permify', // Usually your GitHub org/user name.
  projectName: 'docs', // Usually your repo name.
  trailingSlash: false,

  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        googleAnalytics: {
          trackingID: 'G-BSRXWHBYR1',
          anonymizeIP: true,
        },
        gtag: {
          trackingID: 'G-BSRXWHBYR1',
          anonymizeIP: true,
        },
      },
    ],
  ],

  themeConfig:
      {
        navbar: {
          title: 'Permify',
          logo: {
            alt: 'Permify Logo',
            src: 'img/logo.svg',
            href: 'https://www.permify.co/',
          },
          items: [
            {
              type: 'doc',
              docId: 'permify-overview/intro',
              position: 'left',
              label: 'Docs',
            },
            {
              type: 'dropdown',
              label: 'API Reference',
              position: 'left',
              items: [
                {
                  label: 'gRPC API Reference',
                  href: 'https://buf.build/permify/permify/docs/main:base.v1',
                },
                {
                  label: 'REST API Reference',
                  href: 'https://app.swaggerhub.com/apis-docs/permify/permify/latest',
                },
              ],
            },
            {
              href: 'https://github.com/Permify/permify',
              position: 'right',
              className: 'header-github-link'
            },
            {
              href: 'https://discord.gg/MJbUjwskdH',
              position: 'right',
              className: 'header-discord-link'
            },
            {
              href: 'https://twitter.com/getPermify',
              position: 'right',
              className: 'header-twitter-link'
            }
          ],
        },
        footer: {
          style: 'dark',
          links: [
            {
              title: 'Docs',
              items: [
                {
                  label: 'Docs',
                  to: '/docs/permify-overview/intro',
                },
              ],
            },
            {
              title: 'Community',
              items: [
                {
                  label: 'Twitter',
                  href: 'https://twitter.com/getPermify',
                },
              ],
            },
            {
              title: 'More',
              items: [
                {
                  label: 'Blog',
                  to: 'https://www.permify.co/blog',
                },
                {
                  label: 'GitHub',
                  href: 'https://github.com/Permify',
                },
              ],
            },
          ],
          copyright: `Copyright © ${new Date().getFullYear()} Permify.`,
        },
        colorMode: {
          disableSwitch: false,
          respectPrefersColorScheme: true,
        },
        prism: {
          theme: darkCodeTheme,
          darkTheme: darkCodeTheme,
          additionalLanguages: ['php'],
        },
      },
};

module.exports = config;
