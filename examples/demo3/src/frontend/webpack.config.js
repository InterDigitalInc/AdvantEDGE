const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const dotenv = require('dotenv');
const webpack = require('webpack');

var extractPlugin = new ExtractTextPlugin({
  filename: 'bundle.css'
});

var htmlPlugin = new HtmlWebpackPlugin({
  template: 'src/index.html'
});

module.exports = (env) => {
  // call dotenv and it will return an Object with a parsed key
  const environmentalVariable = dotenv.config().parsed;

  // reduce it to a nice object, the same as before
  const envKeys = Object.keys(environmentalVariable).reduce((prev, next) => {
    prev[`process.env.${next}`] = JSON.stringify(environmentalVariable[next]);
    return prev;
  }, {});

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
      }),
      new webpack.DefinePlugin(envKeys)
    ],
    devServer: {
      proxy: {
        '/': {
          target: 'http://' + (env ? env.MEEP_HOST : ''),
          secure: false
        }
      }
    }
  };
};
