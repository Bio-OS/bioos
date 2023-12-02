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

import React, { useEffect } from 'react';
import { Select, Space, Typography } from '@arco-design/web-react';

import { ANALYSIS_STATUS, RUN_STATUS_TAG, LANGUAGE_TYPES } from 'helpers/constants';

import styles from './analysis.less';

interface ListSelectProps {
  showStatusFlag?: boolean;
  showWorkflowFlag?: boolean;
  showLanguageFlag?: boolean;
  statusID?: string;
  workflowID?: string;
  language?: string;
  listWorkFlowItems?: any;
  onChangeStatus?: (value: string) => void;
  onChangeWorkflow?: (value: string) => void;
  onChangeLanguage?: (value: string) => void;

}

const ListSelects: React.FC<ListSelectProps> = ({
  showStatusFlag = true,
  showWorkflowFlag = false,
  showLanguageFlag= false,
  statusID,
  workflowID,
  language,
  listWorkFlowItems,
  onChangeStatus,
  onChangeWorkflow,
  onChangeLanguage,
}) => (
  <Space size={12}>
    {showStatusFlag && (
      <Select
        placeholder={'请选择'}
        addBefore={`${showWorkflowFlag ? '分析' : '运行'}状态`}
        value={statusID}
        className={styles.historySelect}
        onChange={value => {
          onChangeStatus && onChangeStatus(value);
        }}
        triggerProps={{
          autoFitPosition: false,
        }}
        notFoundContent={'暂无数据'}
      >
        <Select.Option value="All">全部</Select.Option>
        {(showWorkflowFlag ? ANALYSIS_STATUS : RUN_STATUS_TAG).map(
          ({ value, text }) => (
            <Select.Option key={value} value={value}>
              {text}
            </Select.Option>
          ),
        )}
      </Select>
    )}
    {showWorkflowFlag && (
      <Select
        placeholder={'请选择'}
        addBefore="所属工作流"
        value={workflowID}
        className={styles.historySelect}
        renderFormat={(_, value) =>
          listWorkFlowItems?.find(_ => _?.id === value)?.name || '全部'
        }
        onChange={value => {
          onChangeWorkflow && onChangeWorkflow(value);
        }}
        triggerProps={{
          autoFitPosition: false,
        }}
        notFoundContent={'暂无数据'}
      >
        {<Select.Option value="">全部</Select.Option>}
        {listWorkFlowItems?.map(({ id, name }, index) => (
          <Select.Option
            key={index}
            value={id}
            className="flexAlignCenter"
            style={{ height: 36 }}
          >
            <Typography.Text
              ellipsis={{
                showTooltip: {
                  type: 'popover',
                  props: {
                    position: 'right',
                    style: {
                      transform: 'translate(8px, 0px)',
                    },
                    content: <div style={{ maxWidth: 200 }}>{name}</div>,
                  },
                },
              }}
            >
              {name}
            </Typography.Text>
          </Select.Option>
        ))}
      </Select>
    )}
    {showLanguageFlag && (
        <Select
            placeholder={'请选择'}
            addBefore={'规范'}
            value={language}
            className={styles.historySelect}
            onChange={value => {
              onChangeLanguage && onChangeLanguage(value);
            }}
            triggerProps={{
              autoFitPosition: false,
            }}
            notFoundContent={'暂无数据'}
        >
          <Select.Option value="All">全部</Select.Option>
          {LANGUAGE_TYPES.map(({value, text}) => (
            <Select.Option value={value} key={value}>{text}</Select.Option>
          ))}
        </Select>
    )}
  </Space>
);

export default ListSelects;
