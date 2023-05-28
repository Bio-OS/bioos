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

import { FieldState, FieldValidator } from 'final-form';

export default function combineValidators<FieldValue = any>(
  validators: FieldValidator<FieldValue>[] | FieldValidator<FieldValue>,
): (
  value: FieldValue,
  allValues: any,
  meta?: FieldState<FieldValue>,
) => string | undefined | Promise<any> {
  if (!Array.isArray(validators)) {
    validators = [validators as FieldValidator<FieldValue>];
  }

  return (value: FieldValue, allValues: any, meta?: FieldState<FieldValue>) => {
    // 直接用 async/await 会导致，整个方法都是异步的，校验函数异步会导致抖动
    // 表单的其他字段也会触发校验函数，异步校验函数可以做缓存，如果 value 没有变化，
    // 可以同步返回（同步函数不会有抖动）
    function runValidator(
      err: string | undefined | Promise<any>,
      validatorArr: FieldValidator<FieldValue>[],
    ): string | undefined | Promise<any> {
      if (err) {
        return err;
      }
      const validator = validatorArr.shift();
      if (!validator) {
        return err;
      }

      const result = validator(value, allValues, meta);
      if (!result || !result.then) {
        return runValidator(result, validatorArr);
      }

      return result.then((e: any) => runValidator(e, validatorArr));
    }

    return runValidator('', [...(validators as FieldValidator<FieldValue>[])]);
  };
}
