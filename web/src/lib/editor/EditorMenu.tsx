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

import { FC, ReactNode } from 'react';
import { Divider, Dropdown, Menu, Modal } from '@arco-design/web-react';

import { TreeNode } from './types';
import { useBasicEditContext } from './useBasicEditState';
import { clipboard } from './utils';

const EditorMenu: FC<{
  data?: TreeNode;
  disabled?: boolean;
  dir?: boolean;
  /**
   * 粘贴版上有内容才允许粘贴
   */
  isPaste: boolean;
  onCopy?: (path: string) => void;
  onCut?: (path: string) => void;
  onPaste: (path: string) => void;
  children: ReactNode;
}> = ({ data, disabled, dir, isPaste, onCopy, onCut, onPaste, children }) => {
  const { addFile, renameState, deleteFile } = useBasicEditContext();

  const handleAddFile = () => {
    addFile(data?.id, data?.name);
  };

  const handleAddFolder = () => {
    addFile(data?.id, data?.name, true);
  };

  const handleRename = () => {
    renameState(data?.id);
  };

  const handleDelete = () => {
    Modal.confirm({
      title: `确定删除所选文件${data?.isLeaf ? '' : '夹'} ${
        data?.fileName
      } 吗？`,
      okButtonProps: { status: 'danger' },
      onOk: () => {
        deleteFile(data?.id);
      },
    });
  };

  const handleCopyPath = () => {
    clipboard(data?.name);
  };

  const handleCopy = () => {
    onCopy?.(data?.id);
  };

  const handleCut = () => {
    onCut?.(data?.id);
  };

  const handlePaste = () => {
    onPaste(data?.name);
  };

  return (
    <Dropdown
      trigger="contextMenu"
      position="bl"
      disabled={disabled}
      droplist={
        <Menu
          className="basic-editor-tree-menu"
          onClickMenuItem={(_key, e) => e.stopPropagation()}
        >
          {(dir || !data) && (
            <>
              <Menu.Item key="1" onClick={handleAddFile}>
                新增文件
              </Menu.Item>
              <Menu.Item key="2" onClick={handleAddFolder}>
                新增文件夹
              </Menu.Item>
              <Divider />
            </>
          )}
          <Menu.Item key="4" onClick={handleCopyPath}>
            复制相对路径
          </Menu.Item>
          {(data || isPaste) && <Divider />}
          {data && (
            <>
              <Menu.Item key="5" onClick={handleCut}>
                剪切
              </Menu.Item>
              <Menu.Item key="6" onClick={handleCopy}>
                复制
              </Menu.Item>
            </>
          )}
          {(dir || !data) && isPaste && (
            <Menu.Item key="7" onClick={handlePaste}>
              粘贴
            </Menu.Item>
          )}

          {data && (
            <>
              <Divider />
              <Menu.Item key="8" onClick={handleRename}>
                重命名
              </Menu.Item>
              <Menu.Item key="9" onClick={handleDelete}>
                delete
              </Menu.Item>
            </>
          )}
        </Menu>
      }
    >
      {children}
    </Dropdown>
  );
};

export default EditorMenu;
