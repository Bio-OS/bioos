/* eslint-disable */
const path = require('path');
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const FriendlyErrorsWebpackPlugin = require('@nuxt/friendly-errors-webpack-plugin');
const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
const portfinder = require('portfinder');

module.exports = (env, args) =>
  new Promise((resolve, reject) => {
    portfinder.getPort(
      {
        port: 8901,
      },
      (error, port) => {
        if (error) {
          reject(error);
        } else {
          resolve(
            merge(common(env, args), {
              mode: 'development',
              stats: 'errors-only',
              devServer: {
                historyApiFallback: {
                  disableDotRule: true, //  解决路径中包含.ipynb刷新出现404的问题
                },
                proxy: {
                  '/api': {
                    target: 'http://localhost:8888',
                    changeOrigin: true,
                    pathRewrite: { '^/api': '' },
                  },
                },
                static: {
                  directory: path.join(__dirname, 'build'),
                },
                port,
                hot: true,
              },
              plugins: [
                new ReactRefreshWebpackPlugin(),
                new ForkTsCheckerWebpackPlugin(),
                new FriendlyErrorsWebpackPlugin({
                  compilationSuccessInfo: {
                    messages: [
                      'Project is running at \033[1;36mhttp://localhost:' +
                        port,
                    ],
                  },
                }),
              ],
            }),
          );
        }
      },
    );
  });
