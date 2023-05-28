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

import JSZip from 'jszip';

export type TreeNode = {
  id: string;
  name: string;
  fileName: string;
  isLeaf: boolean;
  children?: TreeNode[];
  rename?: boolean;
  suffix?: string;
  draggable: boolean;
};

export type Clipboard = { id: string; action: 'copy' | 'cut' };

export type FileNode = Partial<JSZip.JSZipObject> & {
  id: string;
  /**
   * 全局唯一，路径+文件名
   */
  name: string;
  dir: boolean;
  fileName: string;
  content?: string | Uint8Array;
  contentLoading?: boolean;
  language?: string;
  suffix?: string;
  rename?: boolean;
};

export type UnzipOption = {
  binary?: boolean;
};
