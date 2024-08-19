import { defaultTheme } from '@vuepress/theme-default'
import { defineUserConfig } from 'vuepress/cli'
import { viteBundler } from '@vuepress/bundler-vite'
import { shikiPlugin } from '@vuepress/plugin-shiki'

export default defineUserConfig({
  lang: 'en-US',

  title: 'DaVinci Script',
  description: 'DaVinci Script is a fast, simple and powerful scripting language for creating dynamic web pages.',

  theme: defaultTheme({
    logo: 'davinci.png',
    locales: {
      '/': {
        navbar: [
          '/guide/get-started.md',
          '/guide/tutorials.md',
          '/guide/functions.md',
        ],
        sidebar: [
          {
            text: 'Guide',
            children: [
              '/guide/get-started.md',
              '/guide/tutorials.md',
              '/guide/functions.md',
            ]
          }
        ]
      }
    }
  }),

  bundler: viteBundler(),
  plugins: [
     shikiPlugin({
        langs: ['php', 'bash', 'diff', 'json', 'md', 'ts', 'vue'],
        theme: 'dark-plus',
      })
    ]
})

