import { defineConfig, HeadConfig } from 'vitepress'

const umamiScript: HeadConfig = ["script", {
  defer: "true",
  src: "https://analytics.amberbyte.dev/clackor",
  "data-website-id": "ec93a5fa-a1e4-48bd-8f1a-3640f2f3b879"
}]

const baseHeaders: HeadConfig[] = []

const headers = process.env.NODE_ENV === "production" ?
  [...baseHeaders, umamiScript] :
  baseHeaders

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Flamigo Docs",
  description: "Documentation of Flamigo Framework",
  head: headers,
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
          { text: 'Error Handling', link: '/errors' },
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
          { text: 'Config', link: '/config' },
        ]
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/amberbyte/flamigo' }
    ]
  }
})
