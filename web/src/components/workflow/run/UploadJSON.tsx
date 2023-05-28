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
import { Link, Message, Upload } from '@arco-design/web-react';

import { HandlersWorkflowParam } from 'api/index';

import { isContextValue, renderKey } from './utils';

export default function UploadJSON({
  disabled,
  checkData,
  isInput,
  data,
  clickDisabled,
  onUploadComplete,
}: {
  disabled: boolean;
  checkData: () => string[];
  isInput?: boolean;
  data: HandlersWorkflowParam[];
  clickDisabled: boolean;
  onUploadComplete?: () => void;
}) {
  const form = useForm();

  function upload(_fileList, file) {
    if (!file.originFile) return;

    const reader = new FileReader();
    reader.readAsText(file.originFile);

    reader.onerror = () => {
      Message.error('读取JSON文件失败');
    };

    reader.onload = () => {
      let jsonObj: { [key: string]: string } = {};

      try {
        jsonObj = JSON.parse(reader.result as string);
      } catch {
        Message.error('解析JSON文件失败，请检查文件格式');
        return;
      }

      if (Object.getPrototypeOf(jsonObj) !== Object.prototype) {
        Message.error('JSON内容格式不符合规范，请检查');
        return;
      }

      Object.keys(jsonObj).forEach(key => {
        const val = jsonObj[key];
        jsonObj[key] =
          typeof val === 'string' || val === null ? val : JSON.stringify(val);
      });

      let counter = 0;
      const uploadData = checkData();
      form.batch(() => {
        uploadData?.forEach(key => {
          data?.forEach(item => {
            const formKey = renderKey(isInput, key, item.name);
            let jsonVal = jsonObj[item.name];

            // String和File类型需要多加一层双引号
            if (
              jsonVal &&
              (item?.type.startsWith('String') ||
                item?.type.startsWith('File')) &&
              !(jsonVal?.startsWith('"') && jsonVal?.endsWith('"')) &&
              !isContextValue(jsonVal)
            ) {
              // 如果jsonVal可以被JSON.parse，则说明是[]或者{}等复杂类型，不添加引号，触发类型报错
              try {
                const parsedValue = JSON.parse(jsonVal);
                // "123" 应该保留双引号
                if (
                  typeof parsedValue === 'number' ||
                  typeof parsedValue === 'boolean' ||
                  parsedValue === null
                ) {
                  jsonVal = `"${jsonVal}"`;
                }
              } catch (error) {
                // 如果jsonVal是字符串且没有被双引号扩起来，JSON.parse会报错，则需要添加双引号
                jsonVal = `"${jsonVal}"`;
              }
            }

            // 空字符串显示为""
            if (typeof jsonVal === 'string' && jsonVal.trim() === '') {
              jsonVal = `"${jsonVal}"`;
            }

            const formValue = get(form.getState().values, formKey);

            if (
              !(item.name in jsonObj) ||
              formValue === jsonVal ||
              (typeof formValue === 'undefined' && jsonVal === null)
            )
              // 不存在该key，或者值相同时，不进行更新
              return;

            counter++;
            // 自动校验
            form.focus(formKey);
            form.change(formKey, jsonVal === null ? undefined : jsonVal);
            form.blur(formKey);
          });
        });
      });

      if (counter === 0) {
        Message.success('无参数属性值更新');
      } else {
        Message.success(
          `新参数属性值成功，已从JSON文件中更新 ${counter} 条参数`,
        );
      }

      onUploadComplete();
    };
  }

  function handleClick(e) {
    if (clickDisabled) {
      Message.error('请选择实体属性');
      e.stopPropagation();
    }
  }
  return (
    <Upload
      disabled={disabled}
      autoUpload={false}
      showUploadList={false}
      drag={false}
      accept=".json"
      onChange={upload}
    >
      <Link className="mr16" onClick={handleClick}>
        上传JSON文件
      </Link>
    </Upload>
  );
}
