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

/* eslint-disable @typescript-eslint/no-non-null-assertion */
import React, {
  CSSProperties,
  forwardRef,
  ForwardRefRenderFunction,
  ReactNode,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from 'react';
import { Button, Spin } from '@arco-design/web-react';
import { AvailableVirtualListProps } from '@arco-design/web-react/es/_class/VirtualList';
import {
  IconFullscreen,
  IconFullscreenExit,
} from '@arco-design/web-react/icon';
import { EditorProps, OnMount } from '@monaco-editor/react';

import PageEmpty from 'components/Empty';
import MultiRowPopover from 'lib/MultiRowPopover';

import EditorContent from './EditorContent';
import EditorTabs from './EditorTabs';
import EditorTree from './EditorTree';
import { FileNode } from './types';
import useBasicEditState, { BasicEditProvider } from './useBasicEditState';
import { getFile } from './utils';

import './index.less';

export * from './utils';

type IStandaloneCodeEditor = Parameters<OnMount>[0];
export interface BasicEditorProps extends EditorProps {
  /**
   * 是否禁用（不显示）全屏按钮
   */
  fullscreenDisabled?: boolean;
  /**
   * 自定义设置编辑器标题，默认为 “资源管理器”
   */
  editorTitle?: ReactNode;
  /**
   * 由外部控制初始展示的文件
   */
  initialActiveFile?: FileNode;
  /**
   * zip 解压后的目录树数据
   */
  files?: FileNode[];
  /**
   * 容器样式类名
   */
  wrapperClassName?: string;
  /**
   * 容器树样式
   */
  wrapperStyle?: CSSProperties;
  /**
   * 目录树样式类名
   */
  treeClassName?: string;
  /**
   * 目录树样式
   */
  treeStyle?: CSSProperties;
  /**
   * 最大有多少个 tab
   */
  maxCountTab?: number;
  /**
   * 是否只读模式
   */
  isReadOnly?: boolean;
  /**
   * 只读提示信息
   */
  readOnlyTip?: string;
  /**
   * 是否显示左侧菜单顶部 icon 提示信息
   */
  isHeaderTip?: boolean;
  /**
   * 左侧菜单顶部 icon 提示信息
   */
  headerTip?: React.ReactNode;
  /**
   * 请求 zip 包接口 url
   */
  fetchZipUrl?: string;
  /**
   * loading 文案
   */
  loadText?: string;
  /**
   * loading
   */
  loading?: boolean;
  /**
   * 运行时，可选
   */
  runtime?: string;
  /**
   * 目录树虚拟滚动
   */
  treeVirtualListProps?: AvailableVirtualListProps;
  /**
   * 有目录树场景下，编辑器代码文件更新回调
   */
  updateFilesData?: (files: FileNode[]) => void;
  /**
   * 选择目录树节点
   */
  onSelectTreeNode?: (key: string) => void;
}
const BasicEditor: ForwardRefRenderFunction<
  IStandaloneCodeEditor,
  BasicEditorProps
> = (props, ref) => {
  const {
    fetchZipUrl = '',
    theme = 'light',
    options,
    wrapperClassName = '',
    wrapperStyle = {},
    treeClassName = '',
    treeStyle = {},
    maxCountTab = 5,
    runtime,
    loadText = 'loading...',
    loading = false,
    isReadOnly = false,
    readOnlyTip,
    isHeaderTip = false,
    headerTip,
    treeVirtualListProps,
    updateFilesData,
    onSelectTreeNode,
    initialActiveFile,
    editorTitle,
    fullscreenDisabled,
  } = props;
  // 获取当前语言环境

  // 外部需要调用编辑器接口重置布局
  const editorRef = useRef<IStandaloneCodeEditor>(null);
  useImperativeHandle(ref, () => editorRef.current);

  const [fetchLoading, setFetchLoading] = useState(false);
  const [fullscreen, setFullscreen] = useState<boolean>(false);

  const state = useBasicEditState({
    initialActiveFile,
    files: props.files,
    runtime,
    updateFilesData,
    maxCountTab,
  });

  const { resetState, active } = state;

  // 菜单展开或者关闭 width 会变动，重新布局一下
  const layout = () => {
    setTimeout(() => {
      editorRef.current?.layout({ glyphMarginWidth: 200 } as any);
    });
  };

  // 获取 zip 数据，并进行解压处理
  useEffect(() => {
    if (!fetchZipUrl) return;
    (async () => {
      // 先重置组件所有状态，否则表单可以直接用老的内容提交了
      resetState([]);
      setFetchLoading(true);
      try {
        const { files } = await getFile(fetchZipUrl);
        resetState(files);
      } finally {
        setFetchLoading(false);
        layout();
      }
    })();
  }, [fetchZipUrl, runtime]);

  const handleEditorDidMount: OnMount = editor => {
    editorRef.current = editor;
    editor.onDidContentSizeChange(() => {
      // 配合下面 handleTreeExpand 重新布局，菜单展开或者关闭 width 会变动
      editor.layout();
    });
  };

  return (
    <BasicEditProvider value={state}>
      <div
        className={`basic-editor-wrapper ${wrapperClassName}${
          fullscreen ? ' fullscreen' : ''
        }`}
      >
        <Spin loading={loading || fetchLoading} tip={loadText}>
          <div className="basic-editor" style={wrapperStyle}>
            {!fullscreenDisabled && (
              <Button
                type="text"
                size="mini"
                onClick={() => {
                  setFullscreen(!fullscreen);
                  layout();
                }}
                className="basic-editor-wrapper-fullscreen"
                icon={fullscreen ? <IconFullscreenExit /> : <IconFullscreen />}
              />
            )}
            <EditorTree
              className={treeClassName}
              style={treeStyle}
              onLayout={() => {
                layout();
              }}
              virtualListProps={treeVirtualListProps}
              readOnly={isReadOnly}
              isHeaderTip={isHeaderTip}
              headerTip={headerTip}
              editorTitle={editorTitle}
              onSelectTreeNode={onSelectTreeNode}
            />
            <div className="basic-editor-editor-wrapper">
              <EditorTabs />
              <div className="basic-editor-editor">
                {active ? (
                  <MultiRowPopover
                    content={readOnlyTip}
                    position="right"
                    triggerProps={{
                      alignPoint: true,
                      popupAlign: { bottom: 8, left: 20 },
                    }}
                    disabled={!isReadOnly || !readOnlyTip}
                  >
                    <div className="basic-editor-editor-tip">
                      {active.contentLoading ? (
                        <Spin
                          style={{
                            display: 'block',
                            textAlign: 'center',
                            padding: '30%',
                          }}
                        />
                      ) : (
                        <EditorContent
                          theme={theme}
                          options={options}
                          isReadOnly={isReadOnly}
                          handleEditorDidMount={handleEditorDidMount}
                          updateFilesData={updateFilesData}
                        />
                      )}
                    </div>
                  </MultiRowPopover>
                ) : (
                  <PageEmpty desc="未选中任何代码文件" />
                )}
              </div>
            </div>
          </div>
        </Spin>
      </div>
    </BasicEditProvider>
  );
};
export default forwardRef(BasicEditor);
