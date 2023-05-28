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

import React from 'react';
import { Popover, Space } from '@arco-design/web-react';
import { IconQuestionCircle } from '@arco-design/web-react/icon';

import StatusText from './StatusText';

import styles from './analysis.less';

export default function AnalysisStatus({
  title = '分析状态',
  tip = '以下数据指代工作流运行状态',
  statusData = [],
}) {
  return (
    <div className={styles.analysisStatus}>
      <Space className="pl20 pt20">
        <div className="fw500 fs14 colorBlack">{title}</div>
        <Popover position="right" content={tip}>
          <IconQuestionCircle className="colorGrey fs16" />
        </Popover>
      </Space>
      <div className="flex flexWrap pt20 pb12 ">
        {statusData.map(({ name, value, style }, index) => (
          <StatusText
            key={index}
            name={name}
            value={value}
            style={{ ...style }}
          />
        ))}
      </div>
    </div>
  );
}

export function AnalyzeResult({
  data,
  total,
}: {
  total: number;
  data: {
    color: string;
    text: string;
    count: number;
  }[];
}) {
  return (
    <div>
      <Space size="small">
        {data.map(item => (
          <div key={item.text} className="flexAlignCenter colorBlack fs13">
            <div
              style={{
                width: 4,
                height: 4,
                borderRadius: 1,
                background: item.color,
              }}
            ></div>
            <div className="ml4 mr4">{item.text}</div>
            <span>{item.count || 0}</span>
          </div>
        ))}
      </Space>
      <div className="mt4 flexAlignCenter">
        {data.map(
          item =>
            Boolean(item.count) && (
              <React.Fragment key={item.text}>
                <div
                  style={{
                    width: `${Math.floor((item.count / total) * 100)}%`,
                    height: 3,
                    backgroundColor: item.color,
                    marginRight: 2,
                  }}
                ></div>
              </React.Fragment>
            ),
        )}
      </div>
    </div>
  );
}
