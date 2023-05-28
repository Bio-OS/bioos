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
import { useHistory } from 'react-router-dom';
import { Descriptions, Link, Typography } from '@arco-design/web-react';

import Copy from 'components/Copy';

import styles from './analysis.less';

interface AnalysisOverviewProps {
  name?: string;
  id?: string;
  description?: string;
  workspaceId?: string;
  workflowId?: string;
  flagReadFromCache?: boolean;
}

const AnalysisOverview: React.FC<AnalysisOverviewProps> = ({
  name,
  id,
  description,
  workspaceId,
  workflowId,
  flagReadFromCache,
}) => {
  const history = useHistory();

  const data = [
    {
      label: '投递名称',
      value: (
        <Typography.Text
          style={{ maxWidth: 195, marginBottom: 0 }}
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
      label: '投递 ID',
      value: <Copy text={id} maxWidth={175} style={{ color: '#020814' }} />,
    },
    {
      label: '工作流配置',
      value: (
        <Link
          style={{ maxWidth: 195 }}
          className="ellipsis"
          onClick={() => {
            history.push(
              `/workspace/${workspaceId}/workflow/${workflowId}/run?submissionId=${id}`,
            );
          }}
        >
          {name}
        </Link>
      ),
    },
    {
      label: 'CallCaching',
      value: flagReadFromCache ? '启用' : '未启用',
    },
    {
      label: '简短描述',
      value: (
        <Typography.Text
          style={{ width: 195, marginBottom: 0 }}
          ellipsis={{
            showTooltip: {
              type: 'popover',
              props: {
                style: { maxWidth: 350 },
                content: (
                  <div
                    style={{
                      fontSize: 12,
                      maxHeight: 300,
                      overflowY: 'auto',
                    }}
                  >
                    {description}
                  </div>
                ),
              },
            },
          }}
        >
          {description}
        </Typography.Text>
      ),
    },
  ];
  return (
    <Descriptions
      className={styles.workflowInfo}
      column={1}
      title="投递概览"
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
        width: '195px',
        fontSize: 13,
        padding: 0,
        height: 38,
        display: 'flex',
        alignItems: 'center',
      }}
    />
  );
};

export default AnalysisOverview;
