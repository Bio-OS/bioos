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

import { useEffect, useMemo, useState } from 'react';
import { Spin } from '@arco-design/web-react';
import { useMonaco } from '@monaco-editor/react';

import PageEmpty from 'components/Empty';
import { useForceUpdate } from 'helpers/hooks';
import BasicEditor, { splitFileName } from 'lib/editor';
import { FileNode } from 'lib/editor/types';
import Api from 'api/client';
import { HandlersWorkflowFileInfo } from 'api/index';

export default function WDLFileViewer({
  workspaceId,
  workflowId,
  files,
  initialFile,
}: {
  workspaceId: string;
  workflowId: string;
  files: HandlersWorkflowFileInfo[];
  initialFile?: string;
}) {
  const forceUpdate = useForceUpdate();
  const [fileAndDirArr, setFileAndDirArr] = useState<FileNode[]>([]);

  async function fetchWDLFiles() {
    const fileAndDirArr = files.reduce<FileNode[]>((acc, item) => {
      const partArr = item.path.split('/');

      if (!partArr.length) {
        throw new Error(`文件格式错误：${item.path}`);
      }

      partArr.forEach((item, index) => {
        if (index === partArr.length - 1) return;

        const dirName = `${partArr.slice(0, index + 1).join('/')}/`;

        const flagDirExist = acc.find(
          fileObj => fileObj.dir && fileObj.name === dirName,
        );

        if (!flagDirExist) {
          acc.unshift({
            name: dirName,
            fileName: item,
            dir: true,
            id: dirName,
          });
        }
      });

      const fileName = partArr[partArr.length - 1];

      acc.push({
        name: item.path,
        fileName,
        dir: false,
        id: item.id,
        content: '',
        language: splitFileName(fileName).language,
      });

      return acc;
    }, []);
    setFileAndDirArr(fileAndDirArr);
  }
  useEffect(() => {
    fetchWDLFiles();
  }, []);

  const monaco = useMonaco();
  useEffect(() => {
    if (!monaco) return;

    monaco.languages.register({ id: 'wdl' });
    monaco.languages.setMonarchTokensProvider('wdl', {
      keywords: [
        'alias',
        'as',
        'call',
        'command',
        'else',
        'false',
        'if',
        'in',
        'import',
        'input',
        'left',
        'meta',
        'object',
        'output',
        'parameter_meta',
        'right',
        'runtime',
        'scatter',
        'struct',
        'task',
        'then',
        'true',
        'workflow',
        'hints',
      ],
      typeKeywords: [
        'Array',
        'Boolean',
        'Float',
        'Int',
        'Map',
        'None',
        'Object',
        'Pair',
        'String',
        'Directory',
        'File',
      ],

      tokenizer: {
        root: [
          // whitespace
          [/[ \t\r\n]+/, ''],

          // identifiers, keywords, type
          [
            /[a-zA-Z][a-zA-Z0-9_]*/,
            {
              cases: {
                '@keywords': 'keyword',
                '@typeKeywords': 'type',
                '@default': 'identifier',
              },
            },
          ],

          // comment
          [/#.*$/, 'comment'],

          // string
          [/"/, 'string', '@string_double'],
          [/'/, 'string', '@string_single'],
        ],

        string_double: [
          [/[^\\"]+/, 'string'],
          [/\\./, 'string.escape'],
          [/"/, 'string', '@pop'],
        ],
        string_single: [
          [/[^\\']+/, 'string'],
          [/\\./, 'string.escape'],
          [/'/, 'string', '@pop'],
        ],
      },
    });
  }, [monaco]);

  async function fetchFile(fileObj?: FileNode) {
    // 已经加载文件了，或者正在加载，那么不再重复加载
    if (!fileObj || fileObj.contentLoading !== undefined) return;

    fileObj.contentLoading = true;

    const res = await Api.workspaceIdWorkflowWorkflowIdFileDetail(
      fileObj.id,
      workspaceId,
      workflowId,
    );
    if (!res.ok) return;
    fileObj.content = window.atob(res.data.file.content);
    fileObj.contentLoading = false;
    forceUpdate();
  }

  const initialActiveFile = useMemo(() => {
    return fileAndDirArr?.find(item => item.name === initialFile);
  }, [fileAndDirArr]);

  useEffect(() => {
    fetchFile(initialActiveFile);
  }, [initialActiveFile]);

  if (!fileAndDirArr) {
    return (
      <Spin className="w100 textAlignCenter" style={{ padding: '100px 0' }} />
    );
  }

  if (!fileAndDirArr.length) {
    return <PageEmpty desc="没有可以查看的描述文件" />;
  }

  return (
    <BasicEditor
      isReadOnly={true}
      editorTitle="描述文件"
      fullscreenDisabled={true}
      options={{ minimap: { enabled: true }, renderLineHighlight: 'none' }}
      files={fileAndDirArr}
      wrapperStyle={{ minHeight: 600 }}
      initialActiveFile={initialActiveFile}
      onSelectTreeNode={async key => {
        // 点击目录不做操作
        if (key.endsWith('/')) return;

        const fileObj = fileAndDirArr.find(item => item.id === key);
        await fetchFile(fileObj);
      }}
    />
  );
}
