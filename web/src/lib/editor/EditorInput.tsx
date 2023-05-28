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

import React, { useEffect, useRef, useState } from 'react';
import { Input } from '@arco-design/web-react';
import { RefInputType } from '@arco-design/web-react/es/Input/interface';

import Icon from 'components/Icon';
import MultiRowPopover from 'lib/MultiRowPopover';

export default function EditorInput({
  defaultValue,
  onPressEnter,
}: {
  defaultValue: string;
  onPressEnter: (value: string) => void;
}) {
  const [value, setValue] = useState<string>(defaultValue);
  const [valid, setValid] = useState<boolean>(true);
  const inputRef = useRef<RefInputType>(null);

  // focus 并且选中文本（排除后缀名称）
  useEffect(() => {
    if (!inputRef.current) return;

    inputRef.current.focus();

    if (inputRef.current.dom.setSelectionRange) {
      let index = defaultValue.lastIndexOf('.');
      if (index === -1) index = defaultValue.length;
      inputRef.current.dom.setSelectionRange(0, index);
    } else {
      inputRef.current.dom.select();
    }
  }, []);

  const handlePressEnter = () => {
    onPressEnter(valid ? value || defaultValue : defaultValue);
  };

  const handleChange = (val: string) => {
    setValue(val);
    setValid(!val.includes('/'));
  };

  return (
    <MultiRowPopover
      position="rt"
      content={
        <div className="basic-editor-tree-input-verify-error">
          <Icon glyph="editor-exclamation-circle" className="mr8 arco-icon" />
          文件名中不能包含正斜杠 (/)
        </div>
      }
      popupVisible={!valid}
    >
      <Input
        size="mini"
        value={value}
        onChange={handleChange}
        onClick={e => e.stopPropagation()}
        onPressEnter={handlePressEnter}
        onBlur={handlePressEnter}
        ref={inputRef}
        className={!valid ? 'basic-editor-tree-input-error' : ''}
      />
    </MultiRowPopover>
  );
}
