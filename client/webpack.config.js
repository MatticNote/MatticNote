const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const mode = (process.env.NODE_ENV === 'production') ? 'production' : 'development';

module.exports = [
  {
    entry: {
      ui: './src/ui/index.js',
    },
    output: {
      path: `${__dirname}/dist/ui`,
      filename: '[name].js',
    },
    mode: mode,
    module: {
      rules: [
        {
          test: /\.scss$/,
          use: [
            'sass-loader',
            MiniCssExtractPlugin.loader,
            'css-loader',
            'postcss-loader',
          ]
        }
      ]
    },
    plugins: [
      new MiniCssExtractPlugin()
    ]
  }
];
