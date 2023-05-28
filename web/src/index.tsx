/**
 *
 * Copyright 2023 Beijing Volcano Engine Technology Ltd.
 * Copyright 2023 Guangzhou Laboratory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { ConfigProvider } from '@arco-design/web-react';

import App from './App';

import '@arco-themes/react-bioos/index.less';
import 'assets/styles/arco.less';
import 'assets/styles/global.less';

// render sprite
const request = require.context('./assets/svg', false, /\.svg$/);
request.keys().forEach(request);

const ArcoComponentConfig = {
  Pagination: {
    selectProps: {
      getPopupContainer: () => document.body,
    },
  },
};

const container = document.getElementById('root');
const root = createRoot(container);
root.render(
  <ConfigProvider componentConfig={ArcoComponentConfig}>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </ConfigProvider>,
);
