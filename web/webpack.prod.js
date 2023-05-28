/* eslint-disable */
const { merge } = require('webpack-merge');
const common = require('./webpack.common.js');
const TerserPlugin = require('terser-webpack-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

module.exports = (env, args) => {
  const plugins = [
    new MiniCssExtractPlugin({
      filename: 'main.[contenthash:8].css',
    }),
  ];

  if (process.env.BUNDLE_ANALYZE === 'true') {
    plugins.push(new BundleAnalyzerPlugin());
  }

  return merge(common(env, args), {
    mode: 'production',
    optimization: {
      minimizer: [
        new TerserPlugin({
          extractComments: false,
        }),
        new CssMinimizerPlugin(),
      ],
      splitChunks: {
        chunks: 'all',
        minSize: 0,
        minChunks: 1,
        cacheGroups: {
          main: {
            test: /[\\/]node_modules[\\/]/,
            name: 'react',
            priority: -20,
          },
          react: {
            test: /[\\/]node_modules[\\/](react|react-dom)[\\/]/,
            name: 'react',
            priority: -10,
          },
          arco: {
            test: /@arco-design[\\/]web-react/,
            name: 'arco',
            priority: -5,
          },
        },
      },
      runtimeChunk: true,
    },
    plugins,
  });
};
