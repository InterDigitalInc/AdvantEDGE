const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

var extractPlugin = new ExtractTextPlugin({
  filename: 'bundle.css'
});

var htmlPlugin = new HtmlWebpackPlugin({
  template: 'src/index.html'
});

module.exports = () => {
 
  return {
    mode: 'development',
    entry: ['./src/js/demo-controller.js'],
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: 'bundle.js'
    },
    resolve: {
      extensions: ['.js', '.json'],
      alias: {
        '@': path.resolve('src')
      }
    },
    module: {
      rules: [
        {
          test: /\.(js|jsx)$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader',
            options: {
              presets: ['@babel/preset-env', '@babel/preset-react']
            }
          }
        },
        {
          parser: {
            amd: false
          }
        },
        {
          test: /\.css$/,
          use: extractPlugin.extract({
            use: ['css-loader']
          })
        },
        {
          test: /\.(png|svg|jpg|gif)$/,
          use: [
            {
              loader: 'file-loader',
              options: {
                name: '[name].[ext]',
                outputPath: 'img',
                publicPath: 'img'
              }
            }
          ]
        },
        {
          test: /\.scss$/,
          use: extractPlugin.extract({
            use: [
              {
                loader: 'css-loader'
              },
              {
                loader: 'sass-loader',
                options: {
                  includePaths: [
                    './node_modules/material-components-web/node_modules',
                    './node_modules'
                  ]
                }
              }
            ]
          })
        }
      ]
    },
    plugins: [
      htmlPlugin,
      extractPlugin,
      new webpack.DefinePlugin({
        __VERSION__: JSON.stringify('v0.0.0')
      })
    ]
   
  };
};
