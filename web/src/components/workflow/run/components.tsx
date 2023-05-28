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

import { Field } from 'react-final-form';
import classNames from 'classnames';
import {
  Checkbox,
  Dropdown,
  Input,
  Menu,
  Popover,
  Typography,
} from '@arco-design/web-react';
import { IconMoreVertical } from '@arco-design/web-react/icon';

import {
  INPUT_MODEL_DEFAULT_KEY,
  OUTPUT_MODEL_DEFAULT_KEY,
  OUTPUT_PATH_MODEL_DEFAULT_KEY,
} from 'helpers/constants';
import { HandlersWorkflowParam } from 'api/index';

import { renderKey, validateInput, validateOutput } from './utils';

import styles from './style.less';

interface ParamsTitleProps {
  name: string;
  deleteDisabled: boolean;
  checked: boolean;
  onCheck: (checked: boolean) => void;
  onCopy: () => void;
  onRename: () => void;
  onDelete: () => void;
}

export function ParamsTitle({
  name,
  deleteDisabled,
  checked,
  onCheck,
  onCopy,
  onRename,
  onDelete,
}: ParamsTitleProps): JSX.Element {
  if (
    [
      INPUT_MODEL_DEFAULT_KEY,
      OUTPUT_MODEL_DEFAULT_KEY,
      OUTPUT_PATH_MODEL_DEFAULT_KEY,
    ].includes(name)
  ) {
    return <span>属性值</span>;
  }
  return (
    <div className="flexAlignCenter">
      <div style={{ maxWidth: 'calc(100% - 46px)' }}>
        <Typography.Text
          className="lh20"
          ellipsis={{
            cssEllipsis: true,
            showTooltip: {
              type: 'popover',
            },
          }}
        >
          属性值-{name}
        </Typography.Text>
      </div>

      <Dropdown
        droplist={
          <Menu>
            <Menu.Item key="copy" onClick={onCopy}>
              复制
            </Menu.Item>
            <Menu.Item key="rename" onClick={onRename}>
              重命名
            </Menu.Item>
            <Menu.Item
              key="delete"
              disabled={deleteDisabled}
              onClick={onDelete}
            >
              <Popover
                disabled={!deleteDisabled}
                trigger="hover"
                content="请至少保留一个属性值"
                position="right"
              >
                <div>删除</div>
              </Popover>
            </Menu.Item>
          </Menu>
        }
        position="br"
      >
        <IconMoreVertical
          className={classNames(
            [styles.iconMoreVertical],
            'cursorPointer ml4 mr8',
          )}
        />
      </Dropdown>
      <Checkbox checked={checked} onChange={onCheck} />
    </div>
  );
}

// 用于所有表单校验 需要放在真实表单之前渲染 否则真实表单校验失效
export function renderValidateForm(
  columnsNames: string[],
  data: HandlersWorkflowParam[],
  isHidden: boolean,
  isInput: boolean,
  isPath: boolean,
  invalidHeaderArr: string[],
) {
  function getKey() {
    let key = isHidden ? 'hidden' : 'visible';
    if (invalidHeaderArr?.length) {
      key += invalidHeaderArr.join('-');
    }

    return key;
  }
  return (
    <div style={{ display: 'none' }}>
      <Field key="dirty" name="dirty">
        {() => null}
      </Field>
      {columnsNames.map((columnKey: string) => {
        return data?.map(item => (
          <Field
            key={`${getKey()}-${item.name}`} // key不同时才会触发validate
            name={renderKey(isInput, columnKey, item.name)}
            validateFields={[]}
            validate={(value: string, allValues: { [key: string]: string }) =>
              isInput
                ? validateInput(isHidden, isPath, item, value)
                : validateOutput(
                    isHidden,
                    isPath,
                    item,
                    invalidHeaderArr,
                    value,
                    allValues,
                  )
            }
          >
            {() => {
              return (
                <Input placeholder={item.optional ? '可选项' : '必选项'} />
              );
            }}
          </Field>
        ));
      })}
    </div>
  );
}
