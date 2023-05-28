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

import { CSSProperties } from 'react';

import Icon from './Icon';

export default function PageEmpty({
  style,
  desc,
  search,
}: {
  style?: CSSProperties;
  search?: string;
  desc?: string;
}) {
  return (
    <div
      style={{ margin: '80px auto', ...style }}
      className="flexAlignCenter flexCol"
    >
      <Icon size={60} glyph={search ? 'search' : 'blank'} />
      <div className="mt12 fs12 colorGrey" style={{ lineHeight: '12px' }}>
        {search ? `未找到与“${search}”相关的数据` : desc || '暂无数据'}
      </div>
    </div>
  );
}
