var ExtractTextPlugin = require("extract-text-webpack-plugin");

module.exports = {
    entry: {
        main: [
            './assets/javascripts/schedule.js',
            './assets/stylesheets/main.scss',
            './assets/javascripts/main.js',
            './assets/stylesheets/schedule.scss',
        ],
        schedule: [
            './assets/javascripts/schedule.js',
            './assets/stylesheets/schedule.scss',
        ],
        fields: [
            './assets/javascripts/fields.js',
        ],
    },
    output: {
        path: 'public/assets/',
        filename: '[name].js'
    },
    module: {
        loaders: [
            { test: /\.js$/, loader: 'babel-loader?stage=0' },
            { test: /\.s?css$/, loader: ExtractTextPlugin.extract("style-loader", "css-loader!sass-loader") },
            { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: "file" },
            { test: /\.(woff|woff2)$/, loader:"url?prefix=font/&limit=5000" },
            { test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: "url?limit=10000&mimetype=application/octet-stream" },
            { test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: "url?limit=10000&mimetype=image/svg+xml" },
        ]
    },
    plugins: [
        new ExtractTextPlugin("[name].css"),
    ],
};
