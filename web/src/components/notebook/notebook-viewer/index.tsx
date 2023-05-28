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

import { useEffect, useRef } from 'react';
import katex from 'katex';
import { marked } from 'marked';
import { Skeleton } from '@arco-design/web-react';

import prism from './prism/prism.js';
import nbv_constructor from './nbv';

import prismStyles from './prism/prism.less';
import styles from './style.less';
// 处理js中类型不正确
(prism.plugins as any).customClass.map(prismStyles);

export default function NotebookViewer({
  loading,
  notebook,
}: {
  loading?: boolean;
  notebook: object;
}) {
  const refDiv = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (notebook) {
      const nbv = nbv_constructor(document, { katex, prism, marked });
      nbv.render(notebook, refDiv.current);
    }
  }, [notebook]);

  return (
    <>
      {loading && <Skeleton animation={true} loading={loading} />}
      <div ref={refDiv} className={styles.nbvWrap} />
    </>
  );
}
