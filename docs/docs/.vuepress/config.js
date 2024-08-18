import { defaultTheme } from '@vuepress/theme-default'
import { defineUserConfig } from 'vuepress/cli'
import { viteBundler } from '@vuepress/bundler-vite'

export default defineUserConfig({
  lang: 'en-US',

  title: 'Davinci Script',
  description: 'Davinci Script is a fast, simple and powerful scripting language for creating dynamic web pages.',

  theme: defaultTheme({
    logo: 'davinci.png',

    navbar: ['/', '/get-started'],
  }),

  bundler: viteBundler(),
})

