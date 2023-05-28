/* eslint-disable */
const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const ReactRefreshTypeScript = require('react-refresh-typescript');
const WebpackBar = require('webpackbar');

module.exports = (_, args) => {
  const isDev = args.mode === 'development';
  const cssLoader = isDev ? 'style-loader' : MiniCssExtractPlugin.loader;

  return {
    entry: path.resolve(__dirname, '/src/index.tsx'),
    output: {
      path: path.resolve(__dirname, 'build'),
      filename: '[name].[chunkhash:8].js',
      chunkFilename: '[name].chunk.[chunkhash:8].js',
      publicPath: '/',
    },
    resolve: {
      extensions: ['.js', '.ts', '.tsx'],
      alias: {
        components: path.join(__dirname, 'src/components'),
        helpers: path.join(__dirname, 'src/helpers'),
        assets: path.join(__dirname, 'src/assets'),
        lib: path.join(__dirname, 'src/lib'),
        api: path.join(__dirname, 'src/api'),
      },
      fallback: { path: false },
    },
    cache: {
      type: 'filesystem',
    },
    plugins: [
      new CleanWebpackPlugin(),
      new HtmlWebpackPlugin({
        template: path.join(__dirname, 'public', 'index.html'),
        favicon: path.join(__dirname, 'public', 'favicon.png'),
      }),
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify(args.mode),
        PRIMARY_COLOR: JSON.stringify('#1664ff'),
      }),
      new WebpackBar({
        name: 'Bio-OS',
        color: 'cyan',
      }),
    ],
    module: {
      rules: [
        {
          test: /\.less$/i,
          include: [/src/],
          exclude: [/src\/lib/],
          use: [
            cssLoader,
            {
              loader: 'css-loader',
              options: {
                modules: {
                  localIdentName: '[name]__[local]___[hash:base64:5]',
                },
              },
            },
            'postcss-loader',
            {
              loader: 'less-loader',
              options: {
                lessOptions: {
                  globalVars: {
                    primary: '#1664ff',
                  },
                  javascriptEnabled: true,
                },
              },
            },
          ],
        },
        {
          test: /\.less$/i,
          include: [/node_modules/, /src\/lib/],
          use: [cssLoader, 'css-loader', 'less-loader'],
        },
        {
          test: /\.css$/i,
          include: [/node_modules/],
          use: [cssLoader, 'css-loader'],
        },
        {
          test: /\.tsx?$/,
          include: /src/,
          exclude: /node_modules/,
          use: [
            {
              loader: require.resolve('ts-loader'),
              options: {
                getCustomTransformers: () => ({
                  before: [isDev && ReactRefreshTypeScript()].filter(Boolean),
                }),
                transpileOnly: isDev,
              },
            },
          ],
        },
        {
          test: /\.svg$/,
          loader: 'svg-sprite-loader',
          include: path.resolve(__dirname, 'src/assets/svg'),
        },
        {
          test: /\.(jpg|png)$/,
          type: 'asset',
          parser: {
            dataUrlCondition: {
              maxSize: 10 * 1024,
            },
          },
          generator: {
            filename: 'img/[name].[hash:6][ext]',
          },
        },
      ],
    },
  };
};
