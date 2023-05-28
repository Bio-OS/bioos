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

import React, { CSSProperties } from 'react';
import { Field, FieldProps, FieldRenderProps, useForm } from 'react-final-form';
import { FieldState } from 'final-form';
import {
  get,
  isEmpty,
  isFunction,
  isNumber,
  isPlainObject,
  isUndefined,
  omit,
  pick,
  set,
} from 'lodash-es';
import { Form, PopoverProps } from '@arco-design/web-react';
import { FormItemProps } from '@arco-design/web-react/lib/Form';

import Icon from 'components/Icon';

import MultiRowPopover from '../MultiRowPopover';

import combineValidators from './utils/combineValidators';
import rulesCombineValidators, {
  Rule,
  RuleAsync,
  RuleBoolean,
} from './utils/rulesCombineValidators';

const FormItem = Form.Item;

/**
 * 只把 props 上的如下纯 ui 属性赋给 arco 的 Form.Item,
 * 抛弃 arco 自身数据绑定和校验等表单相关的 props,
 * 被抛弃的这部分 props 会赋给 final-form 的 Field 组件
 */
const formItemPropKeys = [
  'noStyle',
  'style',
  'className',
  'label',
  'labelCol',
  'wrapperCol',
  'colon',
  'disabled',
  'required',
  'extra',
  'hasFeedback',
  'help',
  'validateStatus',
];

// final-form 的 Field 组件的 数据绑定和校验等表单相关的 props: https://final-form.org/docs/react-final-form/types/FieldProps
const finalFormFieldPropKeys = [
  'afterSubmit',
  'allowNull',
  'beforeSubmit',
  'data',
  'defaultValue',
  'format',
  'formatOnBlur',
  'initialValue',
  'isEqual',
  'name',
  'parse',
  'subscription',
  'validateFields',
];

type formItemType = Pick<
  FormItemProps,
  | 'style'
  | 'className'
  | 'label'
  | 'labelCol'
  | 'wrapperCol'
  | 'colon'
  | 'disabled'
  | 'required'
  | 'extra'
  | 'hasFeedback'
  | 'noStyle'
  | 'help'
  | 'validateStatus'
  | 'rules'
>;

type fieldType<FieldValue> = Pick<
  FieldProps<FieldValue, FieldRenderProps<FieldValue, HTMLElement>>,
  | 'afterSubmit'
  | 'allowNull'
  | 'beforeSubmit'
  | 'data'
  | 'defaultValue'
  | 'format'
  | 'formatOnBlur'
  | 'initialValue'
  | 'isEqual'
  | 'name'
  | 'parse'
  | 'subscription'
  | 'validateFields'
>;

type FieldValidator<FieldValue> = (
  value: FieldValue,
  allValues: any,
  meta?: FieldState<FieldValue>,
) => any | Promise<any>;

// 容器组件可以承接 Form.Item 注册进来的属性(函数式的不能), 如 Form.Item 内部计算好的 disabled 结果
export const FormItemChildContainer = (props: {
  children: (disabled: boolean) => JSX.Element;
  disabled?: boolean;
}) => {
  const { children, disabled } = props;
  return children(disabled ?? false);
};

/**
 * @title FieldItemProps
 */
type FieldItemProps<FieldValue = any> = Omit<formItemType, 'rules'> &
  Omit<fieldType<FieldValue>, 'validate'> & {
    children:
      | React.ReactElement
      | ((
          props: FieldRenderProps<FieldValue, HTMLElement> & {
            disabled?: boolean;
          },
        ) => JSX.Element);
    /**  当失焦的时候, 恢复组件的值为 initialValue, 当同时设置 initialValue 的时候生效 */
    forceInitializeWhenClear?: boolean;
    /** 错误总是展示在 popover 内, 行内编辑场景下 */
    allErrorPopover?: boolean;
    initialValue?: FieldValue;
  } & (
    | {
        popoverProps?: PopoverProps;
        validate?: FieldValidator<FieldValue> | FieldValidator<FieldValue>[];
        rules?: never;
      }
    | {
        popoverProps?: Omit<PopoverProps, 'content'>;
        rules?: Array<Rule<FieldValue>>;
        validate?: never;
      }
  );

/** 错误提示 */
interface ErrorTipProps {
  message: string;
  style?: CSSProperties;
  status?: 'normal' | 'error' | 'success';
  icon?: 'normal' | 'error' | 'success';
}
function ValidateTip({
  message,
  status = 'normal',
  icon = 'normal',
  style,
}: ErrorTipProps) {
  return (
    <div
      className={`flexAlignCenter field-form-item-verify-${status}`}
      key={message}
      style={style}
    >
      <Icon glyph={icon} size={12} className="mr4" />
      <div style={{ whiteSpace: 'break-spaces', wordBreak: 'break-word' }}>
        {message}
      </div>
    </div>
  );
}

/**
 * - 一个 FieldItem 只能包裹一个子组件, 且该子组件必须是表单组件
 */
export default function FieldItem<FieldValue = any>(
  props: FieldItemProps<FieldValue>,
): React.ReactElement {
  const { className, validateStatus, help, ...formItemProps } = pick(
    props,
    formItemPropKeys,
  ) as Record<string, any>;

  const finalFormFieldProps = pick(props, finalFormFieldPropKeys) as FieldProps<
    FieldValue,
    FieldRenderProps<FieldValue, HTMLElement>
  >;

  const {
    children,
    forceInitializeWhenClear,
    initialValue,
    rules,
    validate,
    allErrorPopover,
  } = props;

  const formInstance = useForm();

  if (rules) {
    finalFormFieldProps.validate = rulesCombineValidators(rules);
  } else {
    finalFormFieldProps.validate = validate
      ? combineValidators<FieldValue>(validate)
      : undefined;
  }

  if (!children || (Array.isArray(children) && children.length > 1)) {
    throw Error('FieldFormItem children.length must be 1');
  }

  type InputProps = FieldRenderProps<FieldValue, HTMLElement> & {
    disabled?: boolean;
  };

  let renderInput: (inputProps: InputProps) => JSX.Element;

  if (typeof children === 'function') {
    renderInput = ({ meta, input, disabled }) =>
      React.cloneElement(
        children({
          meta,
          input,
        }),
        { disabled },
      );
  } else {
    const { onChange, onBlur, onFocus } = children.props;
    if (
      children.props.disabled !== undefined &&
      // InputNumber AutoComplete 内置了 defaultProps.disabled, 肯定不是 undefined
      !['InputNumber', 'AutoComplete'].includes(
        get(children, 'type.displayName'),
      )
    ) {
      console.warn(
        'disabled 属性需要设置在外层 Item 上, 因为 Arco 的 FormItem 源码会强行将 Item 上的 disabled || Form 全局 disabled 覆盖到 children 上',
      );
    }

    renderInput = ({ input, meta, disabled }: InputProps) => {
      const defaultProps = {
        ...input,
        // 手动传递 error API，用于子组件判断当前错误是否由外部(form)触发，会被用户的 error 覆盖
        error: meta.touched ? meta.error || meta.submitError : undefined,
        // 开发者不要在表单组件上写默认值属性, 应该把该属性写在 FieldFormItem 上去
        ...omit(children.props, ['defaultValue', 'initialValue']),
        // input 对象初始值默认是 '', 导致 select 组件的 placehoder 会失效, 需要改为 undefiend, 但这样会使 AutoComplete 组件删除内容时报错, 因为 AutoComplete 内部调用 String 的 toLowerCase 方法, 此时要给 AutoComplete 加 tofilterOption={(inputValue: string, option: any) => !!inputValue} 来修复
        value: String(input.value) === '' ? undefined : input.value,
        // 不同组件的 onChange 抛出来的参数, 数量和类型都是不定的
        onChange: (...args: [any]) => {
          input.onChange(args[0]); // Upload 组件最好自己封装, 别用这个公共组件, 因为它抛出来的参数是 (fileList: UploadItem[], file: UploadItem);
          onChange && onChange(...args);
        },
        onBlur: (e: React.FocusEvent<HTMLElement>) => {
          // 当失焦的时候, 恢复组件的值为 initialValue
          if (
            forceInitializeWhenClear &&
            initialValue !== undefined &&
            (String(input.value) === '' || input.value === undefined)
          ) {
            input.onChange(initialValue);
          }
          input.onBlur(e);
          onBlur && onBlur(e);
        },
        onFocus: (e: React.FocusEvent<HTMLElement>) => {
          input.onFocus(e);
          onFocus && onFocus(e);
        },
        disabled,
      };
      const displayName = get(children, 'type.displayName');
      if (displayName === 'Input' || displayName === 'InputNumber') {
        set(
          defaultProps,
          'autoComplete',
          children?.props?.autoComplete || 'off',
        );
      }
      if (displayName === 'Password') {
        set(
          defaultProps,
          'autoComplete',
          children?.props?.autoComplete || 'new-password',
        );
      }
      return React.cloneElement(children, defaultProps);
    };
  }

  return (
    <Field
      {...finalFormFieldProps}
      validateFields={[]}
      initialValue={initialValue}
    >
      {({ meta, input }) => {
        /** submit 校验时，会把 touched 改为 true，所以不影响滚动聚焦逻辑 */
        const modified = meta.touched || meta.modified || input.value;

        /** popover 的 rules 校验规则提示只在active 和 hover 的时候才可见 */
        /** 用于控制规则 pop 是否固定可见，而不需要用户 focus 或 hover 才可见 */
        const hadValidateError =
          rules &&
          modified &&
          Array.isArray(meta.error) &&
          meta.error.some(e => typeof e === 'string' || e === true);

        return (
          <FormItem
            {...formItemProps}
            validateStatus={
              validateStatus ||
              (modified &&
              ((meta.error && !isPlainObject(meta.error)) ||
                (meta.submitError && !meta.dirtySinceLastSubmit)) &&
              !meta.submitting
                ? 'error'
                : undefined)
            }
            help={
              allErrorPopover
                ? undefined
                : help ??
                  (modified &&
                  (meta.error ||
                    (meta.submitError && !meta.dirtySinceLastSubmit)) &&
                  !meta.submitting
                    ? getError(
                        meta.error,
                        input.value,
                        formItemProps.required,
                        meta.active,
                      ) || meta.submitError
                    : undefined)
            }
            // 规则导致的校验错误，需要添加一个额外的类名，辅助滚动聚焦
            className={`field-form-item-${props.name} ${className || ''} ${
              hadValidateError ? 'form-item-rule-error' : ''
            }`}
          >
            <FormItemChildContainer>
              {(disabled: boolean | undefined) => {
                // FormItem 和 输入组件 之间有其他元素, 会阻隔属性传递, 目前发现的有 disabled; 故手动将 disabled 传递到输入组件上;
                if (rules) {
                  // 异步校验信息不展示在 popover 中
                  const popupRules = rules.filter(
                    item => !(item as RuleAsync<FieldValue>).asyncValidate,
                  );

                  const content = (
                    <>
                      {popupRules.map((rule, index) => {
                        const invalid =
                          Array.isArray(meta.error) && meta.error[index];
                        const ruleWithChildren =
                          rule as RuleBoolean<FieldValue>;
                        let { message: ruleMessage } = ruleWithChildren;
                        if (
                          !ruleMessage &&
                          (rule as RuleAsync<FieldValue>).asyncValidate
                        ) {
                          ruleMessage = meta.error[index];
                        }
                        // 如果输入框不是 required 并且当前 value 为空，则校验状态维持 normal
                        const notRequiredAndEmptyVal =
                          !formItemProps.required && !input.value;

                        const status =
                          !modified || notRequiredAndEmptyVal
                            ? 'normal'
                            : invalid ||
                              (!Array.isArray(meta.error) && meta.error)
                            ? 'error'
                            : 'success';
                        const iconStatus =
                          !modified || notRequiredAndEmptyVal
                            ? 'normal'
                            : invalid
                            ? 'error'
                            : 'success';

                        // 支持根据输入值动态调整message
                        if (isFunction(ruleMessage)) {
                          ruleMessage = ruleMessage(
                            input.value,
                            formInstance.getState().values,
                            formInstance.getFieldState(input.name),
                          );
                        }

                        return (
                          <React.Fragment key={index}>
                            <ValidateTip
                              message={ruleMessage}
                              status={status}
                              icon={iconStatus}
                            />
                            {/* children 不参与校验，只用于子项状态展示 */}
                            {(ruleWithChildren.children || []).map(sub => {
                              const isOk = !sub.validate(input.value);
                              return (
                                <ValidateTip
                                  key={sub.message || ''}
                                  style={{ marginLeft: 24 }}
                                  message={sub.message}
                                  icon={isOk ? 'success' : 'normal'}
                                  status={isOk ? 'success' : 'normal'}
                                />
                              );
                            })}
                          </React.Fragment>
                        );
                      })}
                    </>
                  );

                  const popupVisible =
                    allErrorPopover && hadValidateError
                      ? { popupVisible: true }
                      : {};

                  return (
                    <MultiRowPopover
                      trigger="hover"
                      position="right"
                      getPopupContainer={dom => {
                        let containerDom = dom.parentElement!;
                        // 暂时只考虑悬浮在水平位置的情况（纵向一般是 bottom，且纵向空间相对灵活），最近的挂载点需大于 input 宽度的两倍。
                        if (
                          !['top', 'bottom'].includes(
                            props.popoverProps?.position as string,
                          )
                        ) {
                          while (
                            containerDom &&
                            containerDom !== document.body
                          ) {
                            if (
                              containerDom?.offsetWidth >=
                              dom.offsetWidth * 2
                            ) {
                              break;
                            }

                            containerDom = containerDom.parentElement!;
                          }
                        }
                        return containerDom!;
                      }}
                      {...props.popoverProps}
                      {...popupVisible}
                      content={content}
                    >
                      {renderInput({ meta, input, disabled })}
                    </MultiRowPopover>
                  );
                }

                if (props.popoverProps) {
                  return (
                    <MultiRowPopover
                      trigger={['focus', 'hover']}
                      position="right"
                      {...props.popoverProps}
                    >
                      {renderInput({ meta, input, disabled })}
                    </MultiRowPopover>
                  );
                }
                return renderInput({ meta, input, disabled });
              }}
            </FormItemChildContainer>
          </FormItem>
        );
      }}
    </Field>
  );
}

const getError = (
  error: string | undefined | Array<string | undefined | boolean>,
  value: any,
  required?: boolean,
  active?: boolean,
) => {
  // active 有 rule 的情况下不展示下方的错误
  if (active && Array.isArray(error)) {
    return undefined;
  }
  if (isPlainObject(error)) {
    return undefined;
  }

  if (!Array.isArray(error)) {
    return error;
  }

  // 数组里面有对象会直接报错
  if (error.some(e => isPlainObject(e))) return undefined;

  // TODO: active
  if (required) {
    let empty = true;
    if (typeof value === 'string') {
      empty = !value.trim();
    } else if (isNumber(value)) {
      empty = false;
    } else {
      empty = isEmpty(value);
    }
    if (empty) {
      return '不能为空';
    }
  }
  // 气泡提示有错误条目，并且文本框下没有错误信息，那么文本框下显示固定文案
  if (error.every(a => typeof a === 'boolean' || isUndefined(a)))
    return '不符合输入规则';

  return error;
};
