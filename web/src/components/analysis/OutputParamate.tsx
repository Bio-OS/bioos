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
import { Table, Typography } from '@arco-design/web-react';

import Copy from 'components/Copy';
import Empty from 'components/Empty';

import styles from './analysis.less';

const tableStyle = { width: '100%', height: '100%', overflow: 'auto' };

const OutputParamate: React.FC<{
  outputs?: string;
  style?: React.CSSProperties;
  copyMaxWidth?: number;
  scrollY?: React.CSSProperties['height'];
}> = ({ outputs, style, copyMaxWidth, scrollY = 400 }) => {
  const columnsOutput = [
    {
      title: '变量',
      dataIndex: 'name',
      width: '35%',
      render: (value: string) => (
        <Typography.Paragraph
          ellipsis={{
            rows: 1,
            showTooltip: {
              type: 'tooltip',
              props: {
                color: '#ffffff',
                content: <div className={styles.inputTip}>{value}</div>,
              },
            },
          }}
        >
          {value}
        </Typography.Paragraph>
      ),
    },
    {
      title: '属性值',
      dataIndex: 'Analysis',
      width: '65%',
      render: (value: string) => {
        let copyValue = value;
        if (value.startsWith('"') && value.endsWith('"')) {
          copyValue = value.slice(1, -1);
        }
        return <Copy text={copyValue} maxWidth={copyMaxWidth || 255} />;
      },
    },
  ];
  const onDataChange = outputs => {
    if (!outputs) return;
    const outputsData = JSON.parse(outputs);
    return Object.keys(outputsData).map((key, index) => {
      return {
        key: index,
        name: key,
        Analysis:
          outputsData[key].toString().startsWith('"') &&
          outputsData[key].toString().endsWith('"')
            ? outputsData[key]
            : JSON.stringify(outputsData[key]),
        operator: [outputsData[key], key],
      };
    });
  };

  return (
    <Table
      className={styles.ellipsisText}
      scroll={{ y: scrollY }}
      style={{ ...tableStyle, ...style }}
      border={true}
      columns={columnsOutput}
      data={onDataChange(outputs)}
      pagination={false}
      noDataElement={<Empty desc="暂无数据" />}
    />
  );
};

export default OutputParamate;
