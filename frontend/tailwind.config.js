/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,vue}",
  ],
  darkMode: "class",
  theme: {
    extend: {},
  },
  plugins: [],
  safelist: [
    { pattern: /^text-/ },
    { pattern: /^bg-/, variants: ['hover'] },
  ]
}

