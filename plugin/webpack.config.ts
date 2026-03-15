import path from 'path';
import type { Configuration } from 'webpack';
import CopyWebpackPlugin from 'copy-webpack-plugin';

const config = (env: Record<string, unknown>): Configuration => {
  const isProduction = env.production === true;

  return {
    mode: isProduction ? 'production' : 'development',
    entry: './src/module.ts',
    devtool: isProduction ? 'source-map' : 'eval-source-map',
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: 'module.js',
      library: {
        type: 'amd',
      },
    },
    externals: [
      'react',
      'react-dom',
      '@grafana/data',
      '@grafana/runtime',
      '@grafana/ui',
    ],
    resolve: {
      extensions: ['.ts', '.tsx', '.js'],
    },
    module: {
      rules: [
        {
          test: /\.tsx?$/,
          use: 'swc-loader',
          exclude: /node_modules/,
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
      ],
    },
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: 'src/plugin.json', to: '.' },
          { from: 'src/img', to: 'img', noErrorOnMissing: true },
        ],
      }),
    ],
  };
};

export default config;
