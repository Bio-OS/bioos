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

import { FC, PropsWithChildren } from 'react';
import cs from 'classnames';
import { Popover, PopoverProps } from '@arco-design/web-react';

export interface ComponentProps extends Omit<PopoverProps, 'content'> {
  /** popover content */
  content?: Array<React.ReactNode | string> | React.ReactNode;
}

const MultiRowPopover: FC<PropsWithChildren<ComponentProps>> = ({
  className,
  content,
  disabled,
  getPopupContainer,
  ...resetProps
}) => {
  const getContent = () => {
    if (Array.isArray(content) && content?.length) {
      return (
        <ul>
          {content.map((o, i) => (
            <li key={i}>{o}</li>
          ))}
        </ul>
      );
    }
    return content;
  };

  const cls = cs('filed-popover', className);
  return (
    <Popover
      className={cls}
      content={getContent()}
      {...resetProps}
      // 内容是空的时候不显示 Popover
      disabled={[undefined, null, ''].includes(content as string) || disabled}
      // 默认使用全局 ConfigProvider 的 getPopupContainer, 这里配置的话优先级会更高
      getPopupContainer={getPopupContainer}
    />
  );
};

export default MultiRowPopover;
