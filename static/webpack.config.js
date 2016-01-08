var path = require("path");
var webpack = require('webpack');
var CopyWebpackPlugin = require('copy-webpack-plugin');


module.exports = {
    context: __dirname,
    entry: {
        home: "./js/src/home.jsx"
    },
    output: {
        path: path.resolve("./dist"),
        filename: "js/[name]-bundle.js"
    },
    plugins: [
        new CopyWebpackPlugin([
            { from: 'index.html' },
            { from: 'css', to: 'css' },
            { from: 'images', to: 'images' }
        ])
    ],

    module: {
        loaders: [
            {
                //regex for file type supported by the loader
                test: /\.(jsx)$/,
                exclude: /node_modules/,
                //type of loader to be used
                //loaders can accept parameters as a query string
                loader: 'babel-loader',
                query:
                {
                    plugins: ['transform-runtime'],
                    presets:['react', 'es2015']
                }
            },
            {
                test: /\.js$/, loader: 'babel-loader'
            }
        ]
    }
};

