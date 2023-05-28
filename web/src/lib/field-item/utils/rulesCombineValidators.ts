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

import { FieldState } from 'final-form';

/**
 * 布尔型规则
 * 返回 true 代表校验失败，false 代表通过，这里有点反直觉，但是要与 validate 对齐
 */
export type RuleBoolean<FieldValue> = {
  validate: (
    value: FieldValue,
    allValues: any,
    meta?: FieldState<FieldValue>,
  ) => boolean;
  message:
    | string
    | ((
        value: FieldValue,
        allValues: any,
        meta?: FieldState<FieldValue>,
      ) => string);
  children?: Array<{
    validate: (value: FieldValue) => boolean;
    message: string;
  }>;
};

export type RuleAsync<FieldValue> = {
  asyncValidate: (
    value: FieldValue,
    allValues: any,
    meta?: FieldState<FieldValue>,
  ) => string | undefined | Promise<string | undefined>;
};

export type Rule<FieldValue> = RuleBoolean<FieldValue> | RuleAsync<FieldValue>;

export default function tempCombineValidators<FieldValue = any>(
  rules: Rule<FieldValue>[],
): (
  value: FieldValue,
  allValues: any,
  meta?: FieldState<FieldValue>,
) =>
  | Array<string | boolean | undefined>
  | undefined
  | string
  | Promise<Array<string | boolean | undefined> | undefined | string> {
  return (value, allValues, meta) => {
    // 直接用 async/await 会导致，整个方法都是异步的，校验函数异步会导致抖动
    // 表单的其他字段也会触发校验函数，异步校验函数可以做缓存，如果 value 没有变化，
    // 可以同步返回（同步函数不会有抖动）
    function runAsyncRule(
      ruleArr: RuleAsync<FieldValue>[],
    ): Promise<string | undefined> | string | undefined {
      const rule = ruleArr.shift();
      if (!rule) return undefined;
      const errMsg = rule.asyncValidate(value, allValues, meta);
      if (!errMsg) return runAsyncRule(ruleArr);
      if (typeof errMsg === 'string') return errMsg;

      return (errMsg as Promise<string | undefined>).then(e => {
        if (!e) {
          return runAsyncRule(ruleArr);
        }
        return e;
      });
    }

    const errArr: Array<boolean | string | undefined> = rules
      .filter(rule => !!(rule as RuleBoolean<FieldValue>).message)
      .map(rule =>
        (rule as RuleBoolean<FieldValue>).validate(value, allValues, meta),
      );

    if (errArr.some(e => e)) return errArr;

    const errAsync = runAsyncRule(
      rules.filter(
        rule => (rule as RuleAsync<FieldValue>).asyncValidate,
      ) as RuleAsync<FieldValue>[],
    );

    if (errAsync && (errAsync as Promise<any>).then) {
      return (errAsync as Promise<string | undefined>).then(e => {
        if (errArr.length === 0) return e;
        errArr.push(e);
        if (errArr.every(e => !e)) return undefined;
        return errArr;
      });
    }

    if (errArr.length === 0) return errAsync as string;

    errArr.push(errAsync as string);
    if (errArr.every(e => !e)) return undefined;
    return errArr;
  };
}
