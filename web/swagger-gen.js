/* eslint-disable @typescript-eslint/no-var-requires */
const { generateApi } = require('swagger-typescript-api');
const path = require('path');

generateApi({
  name: 'index.ts',
  input: path.join(__dirname, '../docs/swagger.json'),
  output: path.join(__dirname, './src/api'),
  generateClient: true,
  generateRouteTypes: false,
  silent: true,
  sortTypes: true,
});
