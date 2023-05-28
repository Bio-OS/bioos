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

import { memo } from 'react';
import { Field } from 'react-final-form';
import { Popover, Select, SelectProps } from '@arco-design/web-react';
import { IconExclamationCircleFill } from '@arco-design/web-react/icon';

import { HandlersWorkflowParam } from 'api/index';

interface Props {
  name: string;
  item: HandlersWorkflowParam;
  disabled: boolean;
  selectOptions: SelectProps['options'];
}

export default memo(function FieldItem({
  name,
  item,
  disabled,
  selectOptions,
}: Props) {
  function getPlaceholder() {
    if (disabled) {
      return '请选择';
    }
    return `${item.optional ? '可选项' : '必选项'}，请选择或输入`;
  }

  return (
    <Field<string | undefined> name={name}>
      {({ input, meta }) => {
        return (
          <div className="flexAlignCenter posRelative">
            <Select
              disabled={disabled}
              className="mr8"
              {...input}
              showSearch={{ retainInputValue: true }}
              allowCreate={true}
              allowClear={true}
              value={disabled ? undefined : input.value || undefined}
              error={meta.touched && meta.error}
              placeholder={getPlaceholder()}
              options={selectOptions}
              triggerProps={{
                containerScrollToClose: true,
              }}
              onClear={() => {
                if (!meta.touched) {
                  input.onFocus();
                  input.onBlur();
                }
              }}
            />
            <Popover
              content={
                <div style={{ wordBreak: 'break-word' }}>
                  {meta.touched && meta.error}
                </div>
              }
              position="tr"
            >
              <IconExclamationCircleFill
                className="fs16"
                style={{
                  color: '#e63f3f',
                  visibility: meta.touched && meta.error ? 'visible' : 'hidden',
                }}
              />
            </Popover>
          </div>
        );
      }}
    </Field>
  );
});
