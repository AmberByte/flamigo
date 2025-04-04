import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Flamigo Docs",
  description: "Documentation of Flamigo Framework",

  themeConfig: {
    logo: 'logo-simple.svg',
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/introduction.html' },
      { text: 'GoDoc', link: 'https://pkg.go.dev/github.com/amberbyte/flamigo' }
    ],

    sidebar: [{
        base: '/guide',
        text: 'Getting Started',
        items: [
          { text: 'Introduction', link: '/introduction' },
          { text: 'First Steps', link: '/first-steps' },
        ]
      },
      {
        base: '/guide',
        text: 'Concepts',
        items: [
          { text: 'Actors', link: '/actors' },
          { text: 'Dependency Injection', link: '/dependency-injection' },
          { text: 'Event Bus', link: '/realtime' },
          { text: 'Strategies', link: '/strategies' },
          { text: 'Mocking', link: '/mocking' },
        ]
      },
      {
        base: '/guide',
        text: 'Project Structure',
        items: [
          { text: 'Structure', link: '/structure' },
          { text: '/api', link: '/api' },
          { text: '/domains', link: '/domains' },
          { text: '/interfaces', link: '/interfaces' },
        ]
      },
      {
        base: '/guide',
        text: 'Features',
        items: [
          { text: 'Auth', link: '/auth' },
          { text: 'Websocket', link: '/websocket' },
        ]
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/amberbyte/flamigo' }
    ]
  }
})
