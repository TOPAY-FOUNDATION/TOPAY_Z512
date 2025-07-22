const path = require('path');

module.exports = {
  entry: './src/index.ts',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
  },
  output: {
    filename: 'topayz512.js',
    path: path.resolve(__dirname, 'dist'),
    library: 'TOPAYZ512',
    libraryTarget: 'umd',
    globalObject: 'this',
  },
  target: 'web',
  mode: 'production',
  optimization: {
    minimize: true,
  },
  externals: {
    crypto: 'crypto',
  },
  node: {
    crypto: true,
  },
};