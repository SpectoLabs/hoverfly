var path = require("path");
var webpack = require('webpack');


module.exports = {
    context: __dirname,
    entry: {
        home: "./src/home.jsx"
    },
    output: {
        path: path.resolve("./dist"),
        filename: "[name]-bundle.js"
    },

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

