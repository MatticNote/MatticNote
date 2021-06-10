const colors = require('tailwindcss/colors')

module.exports = {
  purge: [
      './src/**',
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
  plugins: [
      require('@tailwindcss/forms'),
  ],
}
