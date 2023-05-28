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

import { isText } from 'istextorbinary';
import JSZip from 'jszip';
import { isRegExp } from 'lodash-es';

import { fileSuffixWhiteList, languages } from './constants';
import { TreeNode } from './types';
import { FileNode, UnzipOption } from './types';

export const runtimeDefaults = {
  golang: {
    extension: '.go',
    portal: 'main.go',
  },
  node: {
    extension: '.js',
    portal: 'index.js',
  },
  python: {
    extension: '.py',
    portal: 'index.py',
  },
  rust: {
    extension: '.rs',
    portal: 'main.rs',
  },
  java: {
    extension: '.java',
    portal: 'Handler.java',
  },
  wasm: {
    extension: /rs|ts/,
    portal: /index|lib/,
  },
};

// https://github.com/feross/clipboard-copy/blob/master/index.js
export function clipboard(text) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    return navigator.clipboard.writeText(text).catch(function (err) {
      throw err !== undefined
        ? err
        : new DOMException('The request is not allowed', 'NotAllowedError');
    });
  }

  const span = document.createElement('span');
  span.textContent = text;

  span.style.whiteSpace = 'pre';

  document.body.appendChild(span);

  const selection = window.getSelection();
  const range = window.document.createRange();
  selection.removeAllRanges();
  range.selectNode(span);
  selection.addRange(range);

  let success = false;
  try {
    success = window.document.execCommand('copy');
  } catch (err) {
    // eslint-disable-next-line
    console.log('error', err);
  }

  selection.removeAllRanges();
  window.document.body.removeChild(span);

  return success
    ? Promise.resolve()
    : Promise.reject(
        new DOMException('The request is not allowed', 'NotAllowedError'),
      );
}

// 解析出精确的 index 文件
export const resolveIndexFile = (
  files: FileNode[] = [],
  runtime?: string,
): FileNode => {
  const rootFiles = files.filter(a => !a.dir && a.name.split('/').length === 1);

  let indexFile: FileNode | undefined = undefined;

  if (!runtime) return indexFile;

  const conf = Object.entries(runtimeDefaults).find(([language]) => {
    if (runtime.includes(language)) return true;
  })?.[1];

  if (!conf) return indexFile;

  const regExt = !isRegExp(conf.extension)
    ? new RegExp(`${conf.extension}$`)
    : conf.extension;
  const regFileName = !isRegExp(conf.portal)
    ? new RegExp(conf.portal)
    : conf.portal;

  indexFile =
    rootFiles.find(file => {
      const isExt = regExt.test(file.fileName);
      if (!isExt) return false;
      return regFileName.test(file.fileName);
    }) || indexFile;

  return indexFile;
};

/**
 * 生成 tree 组件需要的树结构
 * @param all 所有子集 node
 * @param parents 父级 node
 * @returns 返回树结构
 */
export function generateTreeData(
  all: TreeNode[],
  parents: TreeNode[],
): TreeNode[] {
  parents.forEach(parent => {
    if (parent.isLeaf) return;

    const nextAll = all.filter(
      a => a.name.startsWith(parent.name) && a.name !== parent.name,
    );

    const nextParents = nextAll.filter(item => {
      const paths = item.name.replace(parent.name, '').split('/');
      return (paths.length === 2 && paths[1] === '') || paths.length === 1;
    });

    parent.children = generateTreeData(nextAll, nextParents);
    // 排序
    parent.children = parent.children.sort((a, b) => {
      if (a.isLeaf && b.isLeaf) {
        return a.name.localeCompare(b.name);
      }
      if (!a.isLeaf) return -1;
      return 1;
    });
  });
  return parents;
}

/**
 * 生成文件(夹)名，复制和添加时，可能文件名会重复，需要加数字区分
 * @param files 所有文件
 * @param defaultFileName 默认的文件名称
 * @param path 文件目录
 * @param dir 是否文件夹
 * @param suffix 文件名称后缀，例如 copy
 * @returns
 */
export function generateFileName(
  files: FileNode[],
  defaultFileName: string,
  path = '',
  dir = false,
  suffix = '',
): { name: string; fileName: string } {
  const { name: oldfileName, suffix: fileSuffix } =
    splitFileName(defaultFileName);
  let fileName = oldfileName;
  let name = `${path}${fileName}${dir ? '/' : ''}`;

  let allName = fileSuffix ? `${name}.${fileSuffix}` : name;
  let item = files.find(a => a.name === allName);
  for (let i = 1; item; i++) {
    fileName = `${oldfileName}${suffix}${i}`;
    name = `${path}${fileName}${dir ? '/' : ''}`;

    allName = fileSuffix ? `${name}.${fileSuffix}` : name;
    item = files.find(a => a.name === allName);
  }
  if (fileSuffix) {
    fileName = `${fileName}.${fileSuffix}`;
    name = `${name}.${fileSuffix}`;
  }

  return { name, fileName };
}

/**
 * 逃逸正则
 * 参考 https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions#escaping
 * @param string 字符串
 * @returns 字符串
 */
export function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); // $& means the whole matched string
}

const JWTTOKENKEY = '__cloud_jwt';

const commonHeader = {
  headers: {
    ...(window.localStorage.getItem(JWTTOKENKEY) && {
      // 如果是跨域请求，记得配置一下这个 header，否则跨域请求会报错
      'X-Jwt-Token': window.localStorage.getItem(JWTTOKENKEY),
    }),
  },
};

export const S4 = (): string => {
  return ((1 + Math.random()) * 0x10000 || 0).toString(16).substring(1);
};

export const guid = (): string => {
  return `${S4()}-${S4()}-${S4()}}-${S4()}${S4()}`;
};

export const getNameByPath = (name: string) => {
  const paths = name.split('/');
  // 文件夹最后一个是空字符串
  return paths[paths.length - 1] || paths[paths.length - 2];
};

export const splitFileName = (name: string) => {
  const names = name.split('.');
  if (names.length === 1) {
    return { name: names[0] };
  }

  const suffix = names.pop();
  return { name: names.join('.'), suffix, language: languages[suffix] };
};

async function unzip(
  { files }: JSZip,
  options: UnzipOption = { binary: true },
) {
  const nextFiles = await Promise.all(
    Object.values(files).map(async item => {
      let data: FileNode = {
        ...item,
        id: guid(),
        fileName: getNameByPath(item.name),
      };

      if (!item.dir) {
        const { suffix, language } = splitFileName(data.fileName);

        // 如果开启非binary解析模式，直接判断后缀，否则按照原逻辑进行判断
        const isTxt = !options?.binary
          ? fileSuffixWhiteList.includes(suffix)
          : isText('', await item.async('nodebuffer'));

        data = {
          ...data,
          content: isTxt
            ? await item.async('string')
            : !options?.binary
            ? new Uint8Array()
            : await item.async('uint8array'),
          language: language || (isTxt ? 'txt' : undefined),
          suffix,
        };
      }
      return data;
    }),
  );
  return nextFiles;
}

export const getFile = async (url: string, options?: UnzipOption) => {
  const response = await fetch(`${url}`, {
    ...commonHeader,
  });
  const responseFile = await response.blob();
  const zip = await JSZip.loadAsync(responseFile, { createFolders: true });
  const files = await unzip(zip, options);
  return { files, zip };
};
