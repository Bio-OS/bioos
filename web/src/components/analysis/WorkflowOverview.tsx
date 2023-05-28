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
import { Descriptions, Typography } from '@arco-design/web-react';

import Copy from 'components/Copy';

import styles from './analysis.less';

interface WorkflowOverviewProps {
  name?: string;
  id?: string;
  runName?: string;
  startTime?: string;
  endTime?: string;
}

const WorkflowOverview: React.FC<WorkflowOverviewProps> = ({
  name,
  id,
  runName,
  startTime,
  endTime,
}) => {
  const data = [
    {
      label: '投递名称',
      value: (
        <Typography.Text
          style={{ width: 215, marginBottom: 0 }}
          ellipsis={{
            showTooltip: {
              type: 'popover',
              props: {
                style: { maxWidth: 350 },
              },
            },
          }}
        >
          {name}
        </Typography.Text>
      ),
    },
    {
      label: '数据实体',
      value: (
        <Typography.Text
          style={{ maxWidth: 215, marginBottom: 0 }}
          ellipsis={{
            cssEllipsis: true,
            showTooltip: {
              type: 'popover',
            },
          }}
        >
          {runName}
        </Typography.Text>
      ),
    },
    {
      label: '工作流ID',
      value: <Copy text={id} maxWidth={195} />,
    },
    {
      label: '开始时间',
      value: startTime,
    },
    {
      label: '结束时间',
      value: endTime,
    },
  ];
  return (
    <Descriptions
      className={styles.workflowInfo}
      column={1}
      title="工作流运行概览"
      data={data}
      labelStyle={{
        color: '#86909C',
        padding: '0 20px 0 0',
        lineHeight: '38px',
        fontSize: 13,
        fontWeight: 400,
      }}
      valueStyle={{
        lineHeight: '38px',
        width: '215px',
        fontSize: 13,
        padding: 0,
        height: 38,
        display: 'flex',
        alignItems: 'center',
      }}
    />
  );
};

export default WorkflowOverview;
