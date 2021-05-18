const colors = require('tailwindcss/colors')

module.exports = {
  purge: [
      '../client/src/**/*.html',
      '../mn_template/**/*.django',
  ],
  darkMode: 'class', // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        'light-blue': colors.lightBlue,
        cyan: colors.cyan,
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
