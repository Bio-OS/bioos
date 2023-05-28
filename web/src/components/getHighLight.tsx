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

import { escapeRegExp } from 'lodash-es';
export function genHighlightText(
  text: string,
  keyword: string | undefined,
  color?: string,
  backgroundColor?: string,
  space = true,
) {
  if (!keyword) return text;

  const re = new RegExp(`(${escapeRegExp(keyword)})`, 'i');
  const strArr = text.split(re);

  return strArr.map((item, index) =>
    re.test(item) ? (
      <span
        key={index}
        className={`${space ? 'ml4 mr4' : ''}`}
        style={{
          color: color || '#4086FF',
          background: backgroundColor || '#ffffff',
        }}
      >
        {item}
      </span>
    ) : (
      <span key={index}>{item}</span>
    ),
  );
}
