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

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from 'react';

import { Clipboard, FileNode } from './types';
import {
  escapeRegExp,
  generateFileName,
  getNameByPath,
  guid,
  resolveIndexFile,
  splitFileName,
} from './utils';

export default function useBasicEditState(props: {
  initialActiveFile?: FileNode;
  files: FileNode[];
  maxCountTab: number;
  runtime?: string;
  updateFilesData?: (files: FileNode[]) => void;
}) {
  const { maxCountTab, runtime, updateFilesData, initialActiveFile } = props;
  const [files, setFiles] = useState<FileNode[]>([]);
  const [active, setActive] = useState<FileNode>();
  const [tabs, setTabs] = useState<FileNode[]>([]);
  const [expandedKeys, setExpandedKeys] = useState<string[]>([]);

  useEffect(() => {
    if (props.files === files) return;
    resetState(props.files);
  }, [props.files, runtime]);

  /**
   * 同步删除/重命名/剪切/目录调整等变动
   * @param nextFiles 最新的 files
   */
  function resetTabs(nextFiles: FileNode[]) {
    let activeIndex: number;
    const nextTabs = tabs
      .map((tab, i) => {
        const file = nextFiles.find(item => item.id === tab.id);
        // 如果删除了活跃的 tab，那么需要重新选择活跃的 tab
        if (!file && tab.id === active.id) {
          activeIndex = i;
        }
        return file;
      })
      .filter(item => !!item);

    setTabs(nextTabs);

    if (activeIndex === undefined) {
      setActive(nextTabs.find(item => item.id === active.id));
      return;
    }

    if (nextTabs.length === 0) {
      setActive(undefined);
    } else if (activeIndex + 1 > nextTabs.length) {
      setActive(nextTabs[nextTabs.length - 1]);
    } else {
      setActive(nextTabs[activeIndex]);
    }
  }

  /**
   * 重新上传或者重新拉取线上代码，需要重置所有状态
   * @param newFiles 新的文件
   */
  async function resetState(nextFiles: FileNode[]) {
    const nextExpandedKeys =
      nextFiles?.length && nextFiles[0].dir ? [nextFiles[0].id] : [];
    const indexFile = initialActiveFile || resolveIndexFile(nextFiles, runtime);

    if (indexFile) {
      const names = indexFile.name.split('/');

      setActive(indexFile);
      setTabs([indexFile]);

      names.slice(0, names.length - 1).reduce((agg, curr) => {
        agg = `${agg}${curr}/`;
        nextExpandedKeys.push(nextFiles.find(file => file.name === agg)?.id);
        return agg;
      }, '');
    } else {
      setActive(undefined);
      setTabs([]);
    }

    // TODO: 加上原来老的，否则 tree 组件会报错，找个时间复现一下，看看是不是基础组件 bug
    setExpandedKeys([...expandedKeys, ...nextExpandedKeys]);
    setFiles(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 选中树节点
   * @param id 选中的树节点
   * @param nextFiles 添加文件时全局状态还没有更新，需要传入 files
   * @returns
   */
  function handleSelectTreeNode(id: string, nextFiles: FileNode[] = files) {
    const file = nextFiles.find(item => item.id === id);
    if (!file) return;

    if (!file.dir) {
      setActive(file);
      if (tabs.find(item => item.id === id)) return;
      // 只允许 maxCountTab 个 tab，超出替换最后一个
      setTabs([...tabs.slice(0, maxCountTab - 1), file]);
      return;
    }

    // 选择文字或者 icon 都可以展开或者关闭
    if (expandedKeys.find(a => a === file.id)) {
      setExpandedKeys([...expandedKeys.filter(a => a !== file.id)]);
    } else {
      setExpandedKeys([...expandedKeys, file.id]);
    }
  }

  /**
   * 关闭 tab
   * 需要同步设置活跃的 tab
   * @param id 关闭的 tab id
   */
  function handleCloseTab(id: string) {
    const nextTabs = tabs.filter(item => item.id !== id);
    setTabs(nextTabs);

    // 删除的不是活跃的 tab
    if (id !== active?.id) return;

    const activeIndex = tabs.findIndex(item => item.id === id);

    if (nextTabs.length === 0) {
      setActive(undefined);
    } else if (activeIndex + 1 > nextTabs.length) {
      setActive(nextTabs[nextTabs.length - 1]);
    } else {
      setActive(nextTabs[activeIndex]);
    }
  }

  /**
   * 选中 tab
   * @param id
   */
  function handleSelectTab(id: string) {
    const file = files.find(item => item.id === id);
    setActive(file);
  }

  /**
   * 删除文件(夹)
   * @param id
   */
  function deleteFile(id: string) {
    const file = files.find(item => item.id === id);
    if (!file) return;

    let nextFiles = files;
    if (file.dir) {
      nextFiles = files.filter(item => !item.name.startsWith(file.name));
    } else {
      nextFiles = files.filter(item => item.name !== file.name);
    }

    setFiles(nextFiles);
    resetTabs(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 更新文件或者文件夹名称
   * @param id 文件 id
   * @param fileName 新文件名称
   */
  function updateFileName(id: string, fileName: string) {
    const file = files.find(item => item.id === id);
    if (!file) return;

    // 文件名里面可能包括正则字符，需要逃逸一下
    const newReg = new RegExp(`${escapeRegExp(file.fileName)}/?$`);
    // 新的文件名称
    const name = file.name.replace(newReg, `${fileName}${file.dir ? '/' : ''}`);

    if (files.find(item => item.id !== id && item.name === name)) {
      // 编辑状态改回正常状态
      const nextFiles: FileNode[] = [
        ...files.filter(item => item.id !== id),
        { ...file, rename: false },
      ];
      setFiles(nextFiles);
      updateFilesData?.(nextFiles);
      return `文件名“${fileName}”已被占用，请选取其他名称`;
    }

    const reg = new RegExp(`^${escapeRegExp(file.name)}`);
    const nextFiles = files.map<FileNode>(item => {
      if (item.name === file.name) {
        const { suffix, language } = splitFileName(fileName);
        return {
          ...item,
          name,
          rename: false,
          fileName,
          language:
            language || (typeof item.content === 'string' ? 'txt' : undefined),
          suffix,
        };
      }

      // 文件夹名称修改需要同步修改叶子节点
      if (file.dir && item.name.startsWith(file.name)) {
        const newName = item.name.replace(reg, name);
        return {
          ...item,
          name: newName,
          rename: false,
          fileName: getNameByPath(newName),
        };
      }
      return item;
    });

    setFiles(nextFiles);
    resetTabs(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 更新 active 文件内容
   * TODO: 全量更新性能太差了，键盘快速输入卡顿严重
   * @param value 编辑器内容
   */
  function updateFileContent(value: string) {
    if (!active) return;
    const nextActive: FileNode = { ...active, content: value };
    const nextFiles = [
      ...files.filter(item => item.id !== active.id),
      nextActive,
    ];

    setFiles(nextFiles);
    resetTabs(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 文件名编辑状态，编辑状态
   * @param name 路径
   */
  function renameState(id: string) {
    const nextFiles = files.map<FileNode>(item => {
      if (item.id === id) {
        return { ...item, rename: true };
      }
      return item;
    });

    setFiles(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 添加文件（夹）
   * @param id 文件夹 id
   * @param path 文件夹路径
   * @param dir 是否文件夹
   */
  function addFile(id?: string, path?: string, dir = false) {
    let file: FileNode;

    if (dir) {
      file = {
        id: guid(),
        // 生成文件夹名称，文件夹名称可能存在重复
        ...generateFileName(files, '未命名文件夹', path, true),
        dir: true,
        rename: true,
      };
    } else {
      file = {
        id: guid(),
        // 生成文件名称，文件名称可能存在重复
        ...generateFileName(files, '未命名文件', path),
        dir: false,
        content: '',
        language: 'txt',
        rename: true,
      };
    }

    if (id && !expandedKeys.includes(id)) {
      setExpandedKeys([...expandedKeys, id]);
    }

    const nextFiles = [...files, file];
    !file.dir && handleSelectTreeNode(file.id, nextFiles);
    setFiles(nextFiles);
    updateFilesData?.(nextFiles);
  }

  /**
   * 粘贴文件或者文件夹
   * @param clipboard 粘贴内容
   * @param path 粘贴到哪里
   * @returns
   */
  function paste(clipboard: Clipboard, path: string) {
    const file = files.find(item => item.id === clipboard.id);
    if (!file) return;

    if (clipboard.action === 'cut' && file.dir && path.startsWith(file.name)) {
      return '不允许剪切到下级目录';
    }

    // 剪切到同一个目录，那么不需要操作
    if (
      clipboard.action === 'cut' &&
      file.name === `${path}${file.fileName}${file.dir ? '/' : ''}`
    ) {
      return;
    }

    // 生成文件名称，文件名称可能存在重复
    const { fileName, name } = generateFileName(
      files,
      getNameByPath(file.name),
      path,
      file.dir,
      'Copy',
    );

    const reg = new RegExp(`^${escapeRegExp(file.name)}`);

    // 修改文件名称
    const clipFiles = files
      .filter(
        item =>
          (file.dir && item.name.startsWith(file.name)) ||
          item.name === file.name,
      )
      .map(item => {
        const newName = item.name.replace(reg, name);
        return {
          ...item,
          name: newName,
          fileName: newName === name ? fileName : item.fileName,
          ...(clipboard.action === 'copy' && { id: guid() }), // 复制属于新增文件，需要生成新的 guid
        };
      });

    let nextFiles = [...files, ...clipFiles];

    if (clipboard.action === 'cut') {
      // 删除被剪切的文件
      nextFiles = [
        ...clipFiles,
        ...files.filter(
          item =>
            !(
              (file.dir && item.name.startsWith(file.name)) ||
              item.name === file.name
            ),
        ),
      ];
    }

    setFiles(nextFiles);
    resetTabs(nextFiles);
    updateFilesData?.(nextFiles);
  }

  return {
    handleSelectTreeNode,
    handleExpandTree: setExpandedKeys,
    handleCloseTab,
    handleSelectTab,
    deleteFile,
    updateFileName,
    renameState,
    addFile,
    paste,
    tabs,
    expandedKeys,
    active,
    resetState,
    files,
    updateFileContent,
  };
}

const BasicEditContext = createContext<
  ReturnType<typeof useBasicEditState> | undefined
>(undefined);

export function BasicEditProvider({
  value,
  children,
}: {
  value: ReturnType<typeof useBasicEditState>;
  children: ReactNode;
}): JSX.Element {
  return (
    <BasicEditContext.Provider value={value}>
      {children}
    </BasicEditContext.Provider>
  );
}

export const useBasicEditContext = (): ReturnType<typeof useBasicEditState> => {
  const context = useContext(BasicEditContext);

  // 这个一定会有的，只是为了规避 typescript 类型定义初始值问题
  if (context === undefined)
    throw new Error('Basic Edit Context Value does not exist');

  return context;
};
