const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

var extractPlugin = new ExtractTextPlugin({
  filename: 'bundle.css'
});
var htmlPlugin = new HtmlWebpackPlugin({
  template: 'src/index.html'
});

module.exports = {
  mode: 'development',
  entry: [
    './src/js/meep-controller.js'
  ],
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist'),
  },
  node: {
    fs: 'empty'
  },
  module: {
    rules: [
      {
        parser: {
          amd: false
        }
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              presets: [require.resolve('@babel/preset-react')]
            }
          }]
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
                includePaths: ['./node_modules/material-components-web/node_modules', './node_modules']
                // importer: function(url, prev) {
                //     if (url.indexOf('@material') === 0) {
                //         var filePath = url.split('@material')[1];
                //         var nodeModulePath = `./node_modules/material-components-web/node_modules/@material/${filePath}`;
                //         return { file: require('path').resolve(nodeModulePath) };
                //     }
                //     return { file: url };
                // }
              }
            }]
        })
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
            },
          }]
      },
      {
        test: /\.(ttf|woff|woff2|eot)$/,
        use: [
          {
            loader: 'file-loader',
            options: {
              name: '[name].[ext]',
              outputPath: 'icons',
              publicPath: 'icons'
            },
          }]
      },
      {
        test: /\.html$/,
        use: ['html-loader']
      }]
  },
  plugins: [
    extractPlugin,
    htmlPlugin
  ]
};
