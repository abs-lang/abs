const { description } = require('../../package')

module.exports = {
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#title
   */
  title: ' ',
  /**
   * Ref：https://v1.vuepress.vuejs.org/config/#description
   */
  description: description,

  /**
   * Extra tags to be injected to the page HTML `<head>`
   *
   * ref：https://v1.vuepress.vuejs.org/config/#head
   */
  head: [
    ['meta', { name: 'theme-color', content: '#ec5975' }],
    ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
    ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;500&display=swap' }]
  ],

  /**
   * Theme configuration, here is the default theme configuration for VuePress.
   *
   * ref：https://v1.vuepress.vuejs.org/theme/default-theme-config.html
   */
  themeConfig: {
    repo: 'https://github.com/abs-lang',
    logo: '/abs-horizontal.png',
    search: false,
    editLinks: false,
    docsDir: '',
    editLinkText: '',
    lastUpdated: false,
    sidebarDepth: 0,
    nav: [
      {
        text: 'Introduction',
        link: '/introduction/',
      },
      {
        text: 'Docs',
        link: '/docs/'
      },
      {
        text: 'Playground',
        link: '/playground/'
      }
    ],
    sidebar: {
      '/introduction/': [
        {
          title: 'Introduction',
          collapsable: false,
          children: [
            '',
            'how-to-run-abs-code',
          ]
        }
      ],
      '/docs/': [
        {
          title: 'Getting Started',
          path: '/docs/',
          collapsable: false,
        },
        {
          title: 'Syntax',
          collapsable: false,
          children: [
            'syntax/assignments',
            'syntax/return',
            'syntax/if',
            'syntax/for',
            'syntax/while',
            'syntax/system-commands',
            'syntax/operators',
            'syntax/comments',
          ]
        },
        {
          title: 'Types and Functions',
          collapsable: false,
          children: [
            'types/string',
            'types/number',
            'types/array',
            'types/hash',
            'types/function',
            'types/builtin-function',
            'types/decorator',
          ]
        },
        {
          title: 'Standard Library',
          collapsable: false,
          children: [
            'standard-lib/intro',
            'standard-lib/runtime',
            'standard-lib/cli',
            'standard-lib/util',
          ]
        },
        {
          title: 'Miscellaneous',
          collapsable: false,
          children: [
            'misc/3pl',
            'misc/error',
            'misc/configuring-the-repl',
            'misc/runtime',
            'misc/technical-details',
            'misc/upgrade-from-abs-1-to-2',
            'misc/credits',
          ]
        }
      ],
    }
  },

  /**
   * Apply plugins，ref：https://v1.vuepress.vuejs.org/zh/plugin/
   */
  plugins: [
    '@vuepress/plugin-back-to-top',
    '@vuepress/plugin-medium-zoom',
  ]
}
