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

import Papa from 'papaparse';

// 长度校验
export function validateLen(value: string, min: number, max: number) {
  const regExp = new RegExp(`\^.{${min},${max}}\$`);
  return !regExp.test(value ?? '');
}

// 字符校验
export function validateChart(value: string) {
  return !/^[\p{Script=Han}a-zA-Z0-9-_]+$/u.test(value ?? '');
}

// 字符校验
export function validateBegin(value: string) {
  return !/^[^-_]/.test(value ?? '');
}

export function getNameRules(min: number, max: number) {
  return [
    {
      validate: validateChart,
      message: '仅支持中文、数字、字母、“-”、“_”',
    },
    {
      validate: validateBegin,
      message: '不能以“-”、“_”开头',
    },
    {
      validate: (value?: string) => validateLen(value, min, max),
      message: `长度 ${min}~${max} 个字符`,
    },
  ];
}

export function getLengthRules(min: number, max: number) {
  return [
    {
      validate: (value?: string) => validateLen(value, min, max),
      message: `长度 ${min}~${max} 个字符`,
    },
  ];
}

// 数据解析校验

export type ParseResult = Papa.ParseResult<string[]>;
export const transformName = (header: ParseResult) => {
  if (!header) return '';
  if (!header[0].endsWith('_id')) return '';
  return header[0].replace('_id', '');
};

export const entityRules = [
  {
    validate: (header: ParseResult) => {
      return !transformName(header);
    },
    message: '数据模型格式无效，请按指定格式制作',
  },
  {
    validate: (header: ParseResult) => {
      return transformName(header).length > 30;
    },
    message: 'ID列标题不可超过30个字符',
  },
  {
    validate: (header: ParseResult) => {
      return transformName(header).endsWith('_set');
    },
    message: '实体集合暂不支持上传',
  },
  {
    validate: (header: ParseResult) => {
      return header?.some(item => item?.length > 100);
    },
    message: '行名称超过100个字符',
  },
  {
    validate: (header: ParseResult) => {
      return new Set(header).size !== header?.length;
    },
    message: '表头名称不允许重复',
  },
  {
    validate: (header: ParseResult) => {
      return header?.some(
        item => !new RegExp(`^[0-9a-zA-Z][0-9a-zA-Z-_]*$`).test(item),
      );
    },
    message: `行名称仅支持字母、数字及连接符(-或_)，且不能以连接符开头`,
  },
  {
    validate: (_, rows: ParseResult) => {
      return !rows?.length;
    },
    message: (name: string) =>
      `上传数据模型${name}失败，数据模型格式无效，请按指定格式制作`,
  },
  {
    validate: (header: ParseResult, rows: ParseResult) => {
      const nonStringRows = rows?.filter(row => {
        return !row.every(item => item === '');
      });
      return header?.length > 51 || nonStringRows?.length > 10000;
    },
    message: '表格超过50列或10000行',
  },
  {
    validate: (_, rows: ParseResult) => {
      return rows?.some(row => row[0].length > 500);
    },
    message: '实体ID超过500个字符',
  },
  {
    validate: (_, rows: ParseResult) => {
      return rows?.some(row => row[0] === '');
    },
    message: 'id为空, 请检查相应的id',
  },
];
export const workspaceRules = [
  {
    validate: (header, rows: ParseResult) => {
      return (
        header?.length !== 2 ||
        header?.[0] !== 'Key' ||
        header?.[1] !== 'Value' ||
        !rows?.length ||
        rows?.[0]?.length < 2
      );
    },
    message: '上传Workspace Data失败, 格式无效, 请按指定格式制作',
  },
  {
    validate: (_, rows: ParseResult) => {
      return rows?.length > 3000;
    },
    message: '表格超过了3000行',
  },
  {
    validate: (_, rows: ParseResult) => {
      return rows.some(row => {
        return row.some(item => item.length > 500);
      });
    },
    message: '单元格中超过500个字符',
  },
  {
    validate: (_, rows: ParseResult) => {
      const nonStringRows = rows?.filter(row => {
        return !row.every(item => item === '');
      });
      return nonStringRows.some(row => row[0] === '');
    },
    message: 'key为空, 请检查相应的key',
  },
];
