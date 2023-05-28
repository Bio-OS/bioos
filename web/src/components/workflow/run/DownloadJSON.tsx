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

import { useForm } from 'react-final-form';
import { get } from 'lodash-es';
import { Link } from '@arco-design/web-react';

import { downloadFile } from 'helpers/utils';
import { HandlersWorkflowParam } from 'api/index';

import { renderKey } from './utils';

export default function DownloadJSON({
  disabled,
  checkData,
  isInput,
  data,
}: {
  disabled: boolean;
  checkData: () => string[];
  isInput?: boolean;
  data: HandlersWorkflowParam[];
}) {
  const form = useForm();

  function download() {
    const jsonObj: { [key: string]: string } = {};
    const downLoadData = checkData();
    downLoadData?.forEach(key => {
      data?.forEach(item => {
        let value = get(
          form.getState().values,
          renderKey(isInput, key, item.name),
        );
        // String和File类型 去掉多余的双引号
        if (
          (item?.type.startsWith('String') || item?.type.startsWith('File')) &&
          value?.startsWith('"') &&
          value?.endsWith('"') &&
          value.length >= 2
        ) {
          value = value.slice(1, -1);
        } else {
          // 除了string file 其他类型的都去掉双引号
          try {
            value = JSON.parse(value);
          } catch (error) {}
        }

        // 未输入下载下来显示null
        jsonObj[item.name] = typeof value === 'undefined' ? null : value;
      });
      const href = URL.createObjectURL(
        new Blob([JSON.stringify(jsonObj, null, 2)]),
      );
      downloadFile(href, isInput ? 'inputs.json' : 'outputs.json');
      URL.revokeObjectURL(href);
    });
  }
  return (
    <Link disabled={disabled} onClick={download}>
      下载JSON文件
    </Link>
  );
}
