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

import React, { CSSProperties } from 'react';
import { Table, Typography } from '@arco-design/web-react';

import Empty from 'components/Empty';

import styles from './analysis.less';

const tableStyle = { width: '100%', overflow: 'auto' };

const InputParamate: React.FC<{
  inputs?: string;
  style?: React.CSSProperties;
  scrollY?: CSSProperties['height'];
}> = ({ inputs, style, scrollY = 400 }) => {
  const columnsInput = [
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
      render: (value: any) => (
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
  ];

  const onDataChange = inputs => {
    // deal inputs data
    if (!inputs) return;
    const inputsData = JSON.parse(inputs);
    return Object.keys(inputsData).map((key, index) => {
      return {
        key: index,
        name: key,
        Analysis:
          inputsData[key].toString().startsWith('"') &&
          inputsData[key].toString().endsWith('"')
            ? inputsData[key]
            : JSON.stringify(inputsData[key]),
      };
    });
  };

  return (
    <>
      <Table
        className={styles.ellipsisText}
        scroll={{ y: scrollY }}
        style={{ ...tableStyle, ...style }}
        border={true}
        columns={columnsInput}
        data={onDataChange(inputs)}
        pagination={false}
        noDataElement={<Empty desc="暂无数据" />}
      />
    </>
  );
};

export default InputParamate;
