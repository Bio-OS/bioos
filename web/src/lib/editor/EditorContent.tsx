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

import MonacoEditor, { EditorProps, OnMount } from '@monaco-editor/react';

import { useForceUpdate } from 'helpers/hooks';

import { MonacoOptions } from './constants';
import { FileNode } from './types';
import { useBasicEditContext } from './useBasicEditState';

export default function EditorContent({
  theme,
  options,
  isReadOnly,
  handleEditorDidMount,
  updateFilesData,
}: {
  theme: string;
  options: EditorProps['options'];
  isReadOnly: boolean;
  handleEditorDidMount: OnMount;
  updateFilesData?: (files: FileNode[]) => void;
}) {
  const forceUpdate = useForceUpdate();
  const { files, active } = useBasicEditContext();

  const handleEditorChange = (value?: string) => {
    if (active) {
      // 这种更新方式 react 里面是不允许的，但是全量更新（useBasicEditState.updateFileContent）性能太差了，hack 一下
      active.content = value;
    }

    // 强制组件重新渲染，使 MonacoEditor 组件能获取到最新值
    // MonacoEditor 获取不到最新值，会导致如果修改到原始值，就不会触发 onChange 事件
    // 例如 value 初始 'aaa'，'aaa' > 'aaaa' > 'aaa' 修改成 aaa 时不会触发 onChange 事件，
    // 这时编辑器值和 value 值一致，所以不会出发 onChange
    updateFilesData?.(files);
    forceUpdate();
  };

  return (
    <MonacoEditor
      theme={theme}
      options={{
        ...MonacoOptions,
        ...options,
        readOnly: isReadOnly || typeof active?.content !== 'string',
      }}
      language={active?.language}
      value={
        typeof active?.content === 'string'
          ? active?.content
          : '此文件使用了不支持的文件类型或者是二进制文件，不支持展示'
      }
      path={active?.name}
      onChange={handleEditorChange}
      onMount={handleEditorDidMount}
    />
  );
}
