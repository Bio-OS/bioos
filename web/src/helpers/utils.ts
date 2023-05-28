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

import dayjs from 'dayjs';
import { isBoolean, isEmpty, isNumber } from 'lodash-es';

export const required = (
  val?: string | number | boolean | null,
): string | undefined => {
  let empty = true;
  if (typeof val === 'string') {
    empty = !val.trim();
  } else if (isNumber(val)) {
    empty = false;
  } else if (isBoolean(val)) {
    empty = !val;
  } else {
    empty = isEmpty(val);
  }
  return empty ? '不能为空' : undefined;
};

export const maxLength = (length: number) => (value: string) => {
  if (!value) return;
  if (value.length > length) {
    return '长度超过限制';
  }
};

// 当前时间 年-月-日-时-分-秒
export const currentTime = () => {
  const now = new Date();

  function padString(val: number | string) {
    return String(val).padStart(2, '0');
  }

  return `${now.getFullYear()}-${padString(now.getMonth() + 1)}-${padString(
    now.getDate(),
  )}-${padString(now.getHours())}-${padString(now.getMinutes())}-${padString(
    now.getSeconds(),
  )}`;
};

export function genTime(
  value?: number | null | string,
  format = 'YYYY-MM-DD HH:mm:ss',
) {
  if (!value) return '-';
  if (typeof value === 'number') {
    const result = value / 10 ** 10;
    if (result > 1) {
      return dayjs(value).format(format);
    } else {
      return dayjs.unix(value).format(format);
    }
  } else {
    return dayjs(value).format(format);
  }
}

/**
 * 转换时间为持续时间例如
 * 60000 ms => 1min
 *
 * @param seconds - 间隔多少秒
 * @returns 几天几小时几分钟几秒
 */
export function convertDuration(seconds: number): string {
  const days = Math.floor(seconds / (24 * 3600));
  const hours = Math.floor((seconds - days * 24 * 3600) / 3600);
  const minutes = Math.floor((seconds - days * 24 * 3600 - hours * 3600) / 60);
  const sec = Math.floor(
    seconds - days * 24 * 3600 - hours * 3600 - minutes * 60,
  );
  let str = '';
  if (days) {
    str += `${days}d`;
  }
  if (hours) {
    str += `${hours}h`;
  }
  if (minutes) {
    str += `${minutes}min`;
  }
  if (sec) {
    str += `${sec}s`;
  }
  return str;
}

export const downloadFile = (url: string, name?: string) => {
  const a = document.createElement('a');
  a.href = url;
  a.download = name;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
};

export const getSize = (value: number) => {
  const unit = ['B', 'KB', 'MB', 'GB'];
  const index = Math.floor(Math.log(value) / Math.log(1024));
  const sizeValue = value / Math.pow(1024, index);
  const result = sizeValue.toFixed(0);
  return result + unit[index];
};

export const downloadFileByBlob = (blob: Blob, name: string) => {
  // 生成url对象
  const urlObject = window.URL || window.webkitURL || (window as any);
  const url = urlObject.createObjectURL(blob);
  downloadFile(url, name);
};
