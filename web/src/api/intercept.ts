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

import fetchIntercept from 'fetch-intercept';
import { Message } from '@arco-design/web-react';

export default function interceptFetch() {
  const unregister = fetchIntercept.register({
    request: function (url, config) {
      config = config || {};
      config.headers = {
        ...config.headers,
        authorization: 'Basic YWRtaW46YWRtaW4=',
      };
      return [url, config];
    },

    requestError: function (error) {
      return Promise.reject(error);
    },

    response: function (response) {
      if (response.status === 401) {
        alert('请登录');
      }
      if (response.status === 500) {
        Message.error('系统内部错误');
      }
      return response;
    },

    responseError: function (error) {
      return Promise.reject(error);
    },
  });

  return unregister;
}
