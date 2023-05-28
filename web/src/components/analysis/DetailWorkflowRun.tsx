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

import { Tabs } from '@arco-design/web-react';

import DetailTaskList from './DetailTaskList';
import DetailWorkflowLog from './DetailWorkflowLog';
import InputParamate from './InputParamate';
import OutputParamate from './OutputParamate';

import styles from './analysis.less';

interface DetailWorkflowRunProps {
  inputs?: string;
  outputs?: string;
  logs?: { name: string; info: string }[];
  callCache?: boolean;
  isFinished: boolean;
}
const TabPane = Tabs.TabPane;

export default function DetailWorkflowRun({
  inputs,
  outputs,
  logs,
  callCache,
  isFinished,
}: DetailWorkflowRunProps) {
  return (
    <div className={styles.detailWorklowlistBox}>
      <div className="fs14 fw500 mb20">工作流运行详情</div>

      <Tabs defaultActiveTab="task" type="card-gutter">
        <TabPane key="task" title="Task 列表">
          <DetailTaskList callCache={callCache} isFinished={isFinished} />
        </TabPane>
        <TabPane key="input" title="输入">
          <InputParamate
            inputs={inputs}
            style={{ width: '100%' }}
            scrollY="calc(100vh - 280px)"
          />
        </TabPane>
        <TabPane key="output" title="输出">
          <OutputParamate
            outputs={outputs}
            style={{ width: '100%' }}
            copyMaxWidth={500}
            scrollY="calc(100vh - 280px)"
          />
        </TabPane>
        <TabPane key="log" title="工作流运行日志">
          <DetailWorkflowLog logs={logs} />
        </TabPane>
      </Tabs>
    </div>
  );
}
