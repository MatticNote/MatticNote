/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'media',
  content: [
    './src/**/*.{html,js}',
    '../template/**/*.html',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}
