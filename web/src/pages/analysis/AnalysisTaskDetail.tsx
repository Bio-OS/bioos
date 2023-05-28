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

import { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useRequest } from 'ahooks';
import { Tag } from '@arco-design/web-react';

import AnalysisStatus from 'components/analysis/AnalysisStatus';
import DetailRightContainer from 'components/analysis/DetailRightContainer';
import DetailWorkflowRun from 'components/analysis/DetailWorkflowRun';
import WorkflowOverview from 'components/analysis/WorkflowOverview';
import Breadcrumbs from 'components/Breadcrumbs';
import DetailPage from 'components/DetailPage';
import { RUN_STATUS_TAG } from 'helpers/constants';
import { useQuery } from 'helpers/hooks';
import { genTime } from 'helpers/utils';
import Api from 'api/client';

export default function AnalysisTaskDetail() {
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const { submissionID, runID } = useQuery();
  const [runSubmissionData, setRunSubmissionData] = useState<any>({});
  const [runData, setRunData] = useState<any>({});
  const isFinished = ['Finished', 'Cancelled'].includes(
    runSubmissionData.status,
  );
  const { run: startAnalysisDetailPolling, cancel: stopAnalysisDetailPolling } =
    useRequest(fetchAnalysisDetail, {
      pollingInterval: 4000,
    });
  const {
    run: startAnalysisRunDetailPolling,
    cancel: stopAnalysisRunDetailPolling,
  } = useRequest(fetchAnalysisRunDetail, {
    pollingInterval: 4000,
  });

  // 获取分析历史某一条数据信息
  async function fetchAnalysisDetail() {
    const res = await Api.submissionDetail(workspaceId, {
      ids: [submissionID],
    });

    if (res?.ok) {
      const item = res.data?.items?.[0];
      setRunSubmissionData(item);
      if (['Finished', 'Cancelled'].includes(item.status)) {
        stopAnalysisDetailPolling();
        stopAnalysisRunDetailPolling();
      }
    }
  }

  // 获取分析历史详情 workflow 某一条数据详情
  async function fetchAnalysisRunDetail() {
    const res = await Api.submissionRunDetail(workspaceId, submissionID, {
      ids: [runID],
    });

    if (res?.ok) {
      const item = res.data?.items?.[0];
      setRunData(item);
    }
  }

  const statusData = useMemo(
    () => [
      {
        name: '运行成功',
        value: runData?.taskStatus?.succeeded,
        style: { color: '#1cb267' },
      },
      {
        name: '运行中',
        value: runData?.taskStatus?.running,
        style: { color: PRIMARY_COLOR },
      },
      {
        name: '运行失败',
        value: runData?.taskStatus?.failed,
        style: { color: '#db373f' },
      },
      { name: '启动中', value: runData?.taskStatus?.pending },
      { name: '排队中', value: runData?.taskStatus?.queued },
      { name: '已终止', value: runData?.taskStatus?.cancelled },
    ],
    [runData],
  );

  const logs = useMemo(
    () => [
      { name: '运行日志', info: runData?.log || '-' },
      { name: '错误信息', info: runData?.message || '-' },
    ],
    [runData],
  );

  const status = RUN_STATUS_TAG.find(item => item.value === runData.status);

  useEffect(() => {
    startAnalysisDetailPolling();
    return () => stopAnalysisDetailPolling();
  }, [workspaceId, submissionID]);

  useEffect(() => {
    startAnalysisRunDetailPolling();
    return () => stopAnalysisRunDetailPolling();
  }, [workspaceId, submissionID, runID]);
  return (
    <DetailPage
      breadcrumbs={
        <Breadcrumbs
          breadcrumbs={[
            { text: '分析历史', path: `/workspace/${workspaceId}/analysis` },
            {
              text: '分析历史详情',
              path: `/workspace/${workspaceId}/analysis/detail`,
              query: { submissionID },
            },
            { text: '工作流运行详情' },
          ]}
        />
      }
      title={runData?.name || 'default'}
      statusTag={<Tag color={status?.color}>{status?.text || ''}</Tag>}
      contentStyle={{ paddingRight: 0, paddingBottom: 0 }}
    >
      <div className="flex h100">
        <DetailWorkflowRun
          inputs={runData.inputs}
          outputs={runData.outputs}
          logs={logs}
          callCache={runSubmissionData?.exposedOptions?.readFromCache}
          isFinished={isFinished}
        />
        <DetailRightContainer>
          <AnalysisStatus
            statusData={statusData}
            title="工作流运行状态"
            tip="以下数据指代 task 运行状态"
          />
          <WorkflowOverview
            name={runSubmissionData?.name || '-'}
            id={runData?.engineRunID || '-'}
            runName={runData?.name || 'default'}
            startTime={genTime(runData?.startTime) || '-'}
            endTime={genTime(runData?.finishTime) || '-'}
          />
        </DetailRightContainer>
      </div>
    </DetailPage>
  );
}
