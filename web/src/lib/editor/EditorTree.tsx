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

import React, {
  CSSProperties,
  FC,
  ReactNode,
  useEffect,
  useMemo,
  useState,
} from 'react';
import {
  Button,
  Notification,
  Popover,
  Space,
  Tree,
  Typography,
} from '@arco-design/web-react';
import { AvailableVirtualListProps } from '@arco-design/web-react/es/_class/VirtualList';
import { NodeInstance } from '@arco-design/web-react/es/Tree/interface';
import {
  IconCaretDown,
  IconExclamationCircle,
} from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import MultiRowPopover from 'lib/MultiRowPopover';

import EditorFileIcon from './EditorFileIcon';
import EditorInput from './EditorInput';
import EditorMenu from './EditorMenu';
import { Clipboard, TreeNode } from './types';
import { useBasicEditContext } from './useBasicEditState';
import { generateTreeData } from './utils';

interface Props {
  className?: string;
  style?: CSSProperties;
  virtualListProps?: AvailableVirtualListProps;
  readOnly?: boolean;
  /**
   * 左侧菜单顶部 icon 提示信息
   */
  isHeaderTip: boolean;
  /**
   * 左侧菜单顶部 icon 提示信息
   */
  headerTip?: React.ReactNode;
  onLayout: () => void;
  onSelectTreeNode?: (key: string) => void;
  editorTitle?: ReactNode;
}
const EditorTree: FC<Props> = ({
  className,
  style,
  onLayout,
  virtualListProps,
  readOnly,
  isHeaderTip,
  headerTip,
  onSelectTreeNode,
  editorTitle,
}) => {
  const {
    updateFileName,
    paste,
    active,
    expandedKeys,
    handleSelectTreeNode,
    handleExpandTree,
    files = [],
  } = useBasicEditContext();
  const [clipboard, setClipboard] = useState<Clipboard>();
  const [dragNode, setDragNode] = useState<NodeInstance>();
  const [treeFold, setTreeFold] = useState<boolean>(false);

  const handleSelect = (selectedKeys: string[]) => {
    handleSelectTreeNode(selectedKeys[0]);
    onSelectTreeNode?.(selectedKeys[0]);
  };

  const handleExpand = (keys: string[]) => {
    handleExpandTree(keys);
  };

  const handleChange = (id: string, fileName: string) => {
    const err = updateFileName(id, fileName);
    if (err) {
      Notification.error({ title: '重命名失败', content: err });
    }
  };

  const handleCopy = (id: string) => {
    setClipboard({
      id,
      action: 'copy',
    });
  };

  const handleCut = (id: string) => {
    setClipboard({
      id,
      action: 'cut',
    });
  };

  const handlePaste = (path: string) => {
    const err = paste(clipboard, path);

    if (err) {
      Notification.error({ title: '粘贴失败', content: err });
    }
  };

  const handleDrop = ({
    dropNode,
    dropPosition,
  }: {
    dropNode: NodeInstance | null;
    dropPosition: number;
  }) => {
    if (!dragNode || !dropNode) return;

    let path = (dropNode.props as { dataRef: TreeNode }).dataRef.name;
    // 拖拽到根目录上
    if (dropPosition === -1) path = '';

    paste({ id: dragNode.props._key, action: 'cut' }, path);
    setDragNode(undefined);
  };

  const handleDragStart = (_e: any, node: NodeInstance) => {
    setDragNode(node);
  };

  const handleAllowDrop = ({
    dropNode,
    dropPosition,
  }: {
    dropNode: NodeInstance;
    /**
     * -1 上面
     * 0 重叠
     * 1 下面
     */
    dropPosition: number;
  }) => {
    const drag = (dragNode.props as { dataRef: TreeNode }).dataRef;
    const dragtNames = drag.name.split('/').filter(n => n !== '');
    const target = (dropNode.props as { dataRef: TreeNode }).dataRef;
    const targetNames = target.name.split('/').filter(n => n !== '');

    // 拖拽到根目录，位于本目录下随便哪个文件（夹）上面就可以
    if (dropPosition === -1) {
      return targetNames.length === 1 && dragtNames.length !== 1;
    }

    // 下面无效
    if (dropPosition === 1) return false;

    // 只能拖拽到目录上
    if (target.isLeaf) return false;

    // 不允许拖拽到自己目录或者下级目录，逻辑上不通（只有文件夹有这个逻辑）
    if (!drag.isLeaf && target.name.startsWith(drag.name)) {
      return false;
    }

    // 不允许拖拽到当前自己所在目录下，没有意义
    if (
      `${target.name}${drag.fileName}${drag.isLeaf ? '' : '/'}` === drag.name
    ) {
      return false;
    }

    return true;
  };

  const handleFold = (bl: boolean) => {
    setTreeFold(bl);
    onLayout();
  };

  let fmtWidth = parseInt(style.width as any) || 200;
  fmtWidth = fmtWidth > 200 ? fmtWidth : 200;
  const [widthOptions, setWidthOptions] = useState<{
    width: number;
    currentWidth: number;
    startPageX?: number;
  }>({
    width: fmtWidth,
    currentWidth: fmtWidth,
  });

  const handleWidthStart = (
    e: React.MouseEvent<HTMLDivElement, MouseEvent>,
  ) => {
    setWidthOptions({ ...widthOptions, startPageX: e.pageX });
  };

  useEffect(() => {
    if (widthOptions.startPageX === undefined) return;

    const handleMousemove = (e: MouseEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setWidthOptions({
        ...widthOptions,
        width: widthOptions.currentWidth + e.pageX - widthOptions.startPageX,
      });
    };
    const handleWidthEnd = () => {
      setWidthOptions({
        ...widthOptions,
        currentWidth: widthOptions.width,
        startPageX: undefined,
      });
    };
    document.addEventListener('mousemove', handleMousemove);
    document.addEventListener('mouseup', handleWidthEnd);
    return () => {
      document.removeEventListener('mousemove', handleMousemove);
      document.removeEventListener('mouseup', handleWidthEnd);
    };
  }, [widthOptions]);

  /**
   * 生成 tree 数据
   */
  const treeNodes = useMemo(() => {
    if (files?.length === 0) return;
    const datas: Array<TreeNode> = files.map(
      ({ id, name, dir, rename, fileName, suffix }) => ({
        id,
        name,
        fileName,
        isLeaf: !dir,
        rename,
        suffix,
        draggable: !readOnly && !rename && !widthOptions.startPageX,
      }),
    );

    return generateTreeData(datas, [
      { id: '', name: '', isLeaf: false, fileName: '', draggable: true },
    ])[0].children;
  }, [files, readOnly, widthOptions.startPageX]);

  if (treeFold) {
    return (
      <div className="basic-editor-tree fold">
        <div className="basic-editor-tree-header">
          <Button
            type="text"
            size="small"
            onClick={() => handleFold(false)}
            className="basic-editor-tree-header-right-fold"
            icon={<Icon glyph="collapse" />}
          />
        </div>
      </div>
    );
  }

  const renderTitle = ({ dataRef }: { dataRef: TreeNode }) => {
    const isEdit = dataRef.rename;

    return (
      <EditorMenu
        data={dataRef}
        disabled={readOnly || isEdit}
        onCopy={handleCopy}
        onCut={handleCut}
        onPaste={handlePaste}
        dir={!dataRef.isLeaf}
        isPaste={!!clipboard}
      >
        <div
          className="basic-editor-tree-title"
          onContextMenu={e => e.stopPropagation()}
        >
          <EditorFileIcon suffix={dataRef.suffix} dir={!dataRef.isLeaf} />
          {isEdit ? (
            <EditorInput
              defaultValue={dataRef.fileName}
              onPressEnter={value => handleChange(dataRef.id, value)}
            />
          ) : (
            <Popover
              title={dataRef.fileName}
              position="right"
              triggerProps={{ popupAlign: { right: 30 } }}
            >
              <div className="basic-editor-tree-title-filename">
                {dataRef.fileName}
              </div>
            </Popover>
          )}
        </div>
      </EditorMenu>
    );
  };

  return (
    <div
      className={`basic-editor-tree ${className}`}
      style={{ ...style, width: widthOptions.width }}
      onClick={e => e.stopPropagation()}
    >
      <div className="basic-editor-tree-header">
        <div className="basic-editor-tree-header-left">
          <Space size={4}>
            <Typography.Text style={{ margin: 0 }}>
              {editorTitle || '资源管理器'}
            </Typography.Text>
            {isHeaderTip && (
              <MultiRowPopover
                content={
                  headerTip || (
                    <div style={{ width: '224px' }}>
                      Java、Golang 仅支持在线预览，暂不支持在线编辑代码
                    </div>
                  )
                }
                position="top"
              >
                <IconExclamationCircle className="basic-editor-tree-header-tip" />
              </MultiRowPopover>
            )}
          </Space>
        </div>
        <Button
          type="text"
          size="small"
          onClick={() => handleFold(true)}
          className="basic-editor-tree-header-right"
          icon={<Icon glyph="collapse" />}
        />
      </div>
      <div className="basic-editor-tree-content">
        <Tree
          blockNode
          draggable={!readOnly}
          showLine
          treeData={treeNodes}
          icons={({ isLeaf }) => ({
            switcherIcon: isLeaf ? null : <IconCaretDown />,
          })}
          allowDrop={handleAllowDrop}
          selectedKeys={[active?.id || '']}
          onSelect={handleSelect}
          expandedKeys={expandedKeys}
          onExpand={handleExpand}
          onDragStart={handleDragStart}
          onDrop={handleDrop}
          virtualListProps={virtualListProps}
          fieldNames={{
            key: 'id',
            title: 'fileName',
          }}
          renderTitle={renderTitle as any}
        />
        <EditorMenu
          isPaste={!!clipboard}
          onPaste={handlePaste}
          disabled={readOnly}
        >
          <div className="basic-editor-tree-global" />
        </EditorMenu>
      </div>
      <div className="basic-editor-tree-width" onMouseDown={handleWidthStart} />
    </div>
  );
};
export default EditorTree;
