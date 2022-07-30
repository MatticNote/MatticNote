const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const {WebpackManifestPlugin} = require("webpack-manifest-plugin");
const mode = (process.env.NODE_ENV === 'production') ? 'production' : 'development';

module.exports = [
  {
    entry: {
      ui: './src/ui/index.js',
    },
    output: {
      path: `${__dirname}/dist/ui`,
      filename: '[name].[hash].js',
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
      new CleanWebpackPlugin(),
      new MiniCssExtractPlugin({
        filename: 'ui.[hash].css',
      }),
      new WebpackManifestPlugin({
        publicPath: '/static/ui/'
      }),
    ]
  }
];
