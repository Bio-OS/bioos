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

import { HandlersWorkflowParam } from 'api/index';

export function replaceDotToHyphen(value?: string) {
  if (!value) return '';
  return value.replace(/\./g, '-');
}

export const renderKey = (isInput: boolean, key: string, name: string) => {
  const formPrefix = isInput ? 'input' : 'output';
  return `${formPrefix}.${stringKey(key)}.${replaceDotToHyphen(name)}`;
};

export const stringKey = (key: string) => {
  if (!isNaN(Number(key))) return `key${key}`;
  return key;
};

export const isStringType = (value: string) => {
  return value.startsWith('String') || value.startsWith('File');
};

export const isContextValue = (value: string | undefined) => {
  return value?.startsWith('this.') || value?.startsWith('workspace.');
};

export function getValue(value: string) {
  return isContextValue(value) ? value : JSON.parse(value);
}

const validateQuota = (value?: string) => {
  if (!value) return false;
  return value.length >= 2 && value?.startsWith('"') && value?.endsWith('"');
};

const validateCommon = (
  parame: HandlersWorkflowParam | undefined,
  value: string | undefined,
) => {
  if (isContextValue(value)) {
    return undefined;
  }
  if (parame?.optional && !value) return undefined;

  if (isStringType(parame?.type) && !validateQuota(value)) {
    return '解析错误，除this\\workspace规则索引外，其余输入内容均需在""内。';
  }
  if (
    isStringType(parame?.type) &&
    validateQuota(value) &&
    value?.replace(/^"/, '').replace(/"$/, '').includes('"')
  ) {
    return '解析错误，除this\\workspace规则索引外，其余输入内容均需在""内, 且""内不能再带有"。';
  }
  if (parame?.type.startsWith('Float') && !/^[0-9.]+$/.test(value || '')) {
    return '解析错误，除this\\workspace规则索引外，其值必须为数值，且不允许有""';
  }
  if (parame?.type.startsWith('Int') && !/^[0-9]+$/.test(value || '')) {
    return '解析错误，除this\\workspace规则索引外，其值必须为整数数值，且不允许有""';
  }
  // 对于数组中有些数据无法解析，前端给出提示，例如[a]
  if (
    parame?.type.startsWith('Array') &&
    value?.startsWith('[') &&
    value?.endsWith(']')
  ) {
    try {
      JSON.parse(value);
    } catch (error) {
      if (error) {
        return '解析错误，除this\\workspace规则索引外，其值必须为数组, 且Array内元素必须大于等于1, 数组中的元素为统一类型, 且元素格式要正确';
      }
    }
  }
  // 解决下Array[String]的参数填入[666]没校验元素是否带""，可投递成功，但后端解析的时候就会报错、工作流运行失败
  if (
    parame?.type.startsWith('Array') &&
    value?.startsWith('[') &&
    value?.endsWith(']')
  ) {
    if (
      value?.replace(/\[|]/g, '').length === 0 &&
      !(parame.type.endsWith('+') || parame.type.endsWith('+?'))
    )
      return;

    try {
      const isString = JSON.parse(value).some(
        (item: any) => typeof item !== 'string',
      );
      if (isString) {
        return '解析错误，除this\\workspace规则索引外，其值必须为数组, 且Array内元素必须带""';
      }
    } catch (error) {
      if (error) {
        return '解析错误，除this\\workspace规则索引外，其值必须为数组, 且Array内元素必须带""';
      }
    }
  }
  if (
    (parame?.type.startsWith('Array') &&
      !(value?.startsWith('[') && value?.endsWith(']'))) ||
    (parame?.type.startsWith('Array') &&
      value?.startsWith('[') &&
      value?.endsWith(']') &&
      (parame.type.endsWith('+') || parame.type.endsWith('+?')) &&
      value?.replace(/\[|]/g, '').length === 0)
  ) {
    if (
      parame?.type.startsWith('Array') &&
      !(value?.startsWith('[') && value?.endsWith(']')) &&
      (!parame.type.endsWith('+') || !parame.type.endsWith('+?'))
    ) {
      return '解析错误，除this\\workspace规则索引外，其值必须为数组';
    }

    return '解析错误，除this\\workspace规则索引外，其值必须为数组，且Array内元素必须大于等于1';
  }
  if (
    parame?.type.startsWith('Boolean') &&
    value !== 'true' &&
    parame?.type.startsWith('Boolean') &&
    value !== 'false'
  ) {
    return '解析错误，除this\\workspace规则索引外，其值必须为true或者false，且不允许有""';
  }
};

export const validateInput = (
  isHidden: boolean,
  isPath: boolean,
  item: HandlersWorkflowParam,
  value: string | undefined,
) => {
  if (isHidden) return;
  if (isPath && value?.startsWith('this.')) {
    return '仅采用数据模型开始工作流分析时允许this.规则索引';
  }

  return validateCommon(item, value);
};

export const validateOutput = (
  isHidden: boolean,
  isPath: boolean,
  item: HandlersWorkflowParam,
  invalidHeaderArr: string[],
  value: string,
  allValues,
) => {
  if (isHidden) return;
  if (!value || isPath) return;
  if (!/^this\.[0-9a-zA-Z][0-9a-zA-Z-_]{0,99}$/.test(value)) {
    return '请输入 this.columnName 样式属性值，其中 columnName 仅支持字母、数字、连接符(-、_)，且不能以连接符开头，长度 1 ~ 100 个字符';
  }

  const outputValues = allValues as {
    output?: { [key: string]: string };
  };
  if (
    Object.values(outputValues?.output || {}).filter(item => item === value)
      .length > 1
  ) {
    return '存在相同属性值，请检查修改';
  }

  if (invalidHeaderArr && invalidHeaderArr.includes(value)) {
    return '此列不支持回写';
  }

  return validateCommon(item, value);
};

type ParamObj = {
  [key: string]: string;
};

export function getFormValue(obj: ParamObj) {
  return Object.keys(obj).reduce<ParamObj>((acc, key) => {
    const val = obj[key];
    if (!val && typeof val !== 'boolean') return acc;

    // final form中点符号，有特别意义，这里替换下
    const keyDashed = replaceDotToHyphen(key);

    if (typeof val === 'string' && isContextValue(val)) {
      acc[keyDashed] = val;
    } else {
      acc[keyDashed] = JSON.stringify(val);
    }

    return acc;
  }, {});
}

export function getFormValueDefault(obj?: HandlersWorkflowParam[] | null) {
  if (!obj) return;

  return obj.reduce<ParamObj>((acc, item) => {
    const val = item.default;
    if (!val) return acc;

    // final form中点符号，有特别意义，这里替换下
    const keyDashed = replaceDotToHyphen(item.name);
    acc[keyDashed] = val;

    return acc;
  }, {});
}
