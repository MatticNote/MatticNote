const path = require('path');
const mode = process.env.NODE_ENV === 'production' ? 'production' : 'development';

const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');

const webpackPlugins = [
    new MiniCssExtractPlugin({
        filename: 'matticnote-ui.css',
    })
]

if (mode === 'development') {
    webpackPlugins.push(new HtmlWebpackPlugin({
        template: 'src/ui-scratch.html',
        filename: 'index.html',
    }))
}

module.exports = {
    mode: mode,
    entry: './src/index.js',
    output: {
        path: path.resolve(__dirname, '..', 'static', 'ui'),
        filename: 'matticnote-ui.js',
    },
    module: {
        rules: [
            {
                test: /\.js$/i,
                include: path.resolve(__dirname, 'src'),
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                    },
                },
            },
            {
                test: /\.(sc|c|sa)ss$/i,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'postcss-loader',
                    'sass-loader',
                    'import-glob-loader',
                ],
            },
        ]
    },
    plugins: webpackPlugins,
    optimization: {
        minimizer: [
            new TerserPlugin(),
            new OptimizeCSSAssetsPlugin(),
        ]
    },
    devServer: {
        contentBase: path.resolve(__dirname, 'devDist'),
        writeToDisk: true,
        watchContentBase: true,
        compress: true,
        port: 9000,
    }
}
