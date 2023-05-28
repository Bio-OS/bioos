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

import { memo } from 'react';
import classNames from 'classnames';
import { Radio } from '@arco-design/web-react';

const RADIO_ARRAY = [
  {
    title: '以数据模型作为输入',
    description: '使用已有数据模型，批量分析更便捷',
    value: 'dataModel',
  },
  {
    title: '以文件路径作为输入',
    description: '直接选择文件作为输入，分析更简单',
    value: 'filePath',
  },
];

export type AnalysisModeType = (typeof RADIO_ARRAY)[number]['value'];

function AnalysisMode({
  value,
  onChange,
  disabled,
}: {
  value: string;
  onChange: (val: string) => void;
  disabled?: boolean;
}) {
  return (
    <Radio.Group value={value} className="flex" disabled={disabled}>
      {RADIO_ARRAY.map(item => {
        return (
          <div
            key={item.value}
            className={classNames([
              'fs12 br4 mr12',
              { cursorPointer: !disabled },
            ])}
            onClick={() => {
              if (disabled) return;
              onChange(item.value);
            }}
            style={{
              padding: '12px 16px',
              width: 240,
              border:
                item.value === value
                  ? '1px solid #94c2ff'
                  : '1px solid #e4e8ff',
              background: item.value === value ? '#e8f4ff' : 'white',
            }}
          >
            <div className="flexJustifyBetween">
              <div className="flexAlignCenter">
                <span className="colorBlack fw500 mr4">{item.title}</span>
              </div>

              <Radio value={item.value} style={{ marginRight: 0 }} />
            </div>

            <div className="colorGrey" style={{ marginTop: 2 }}>
              {item.description}
            </div>
          </div>
        );
      })}
    </Radio.Group>
  );
}

export default memo(AnalysisMode);
