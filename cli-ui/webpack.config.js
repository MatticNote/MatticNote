const path = require('path');
const mode = process.env.NODE_ENV === 'production' ? 'production' : 'development';
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

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
                ],
            },
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'matticnote-ui.css',
        })
    ],
}
