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
import { Link, Popover, Table, Tag, Typography } from '@arco-design/web-react';
import { IconQuestionCircle } from '@arco-design/web-react/icon';

import LogModal from 'components/analysis/LogModal';
import PageEmpty from 'components/Empty';
import PaginationPage from 'components/PaginationPage';
import { RUN_STATUS_TAG } from 'helpers/constants';
import { useQuery } from 'helpers/hooks';
import { convertDuration, genTime } from 'helpers/utils';
import Api from 'api/client';
import { HandlersListTasksResponse } from 'api/index';

export default function DetailTaskList({
  callCache,
  isFinished,
}: {
  callCache: boolean;
  isFinished: boolean;
}) {
  const { page, size, submissionID, runID } = useQuery(10);
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [data, setData] = useState<HandlersListTasksResponse>({});
  const { items = [], total } = data;
  const { run: startPolling, cancel: stopPolling } = useRequest(
    fetchAnalysisTaskList,
    {
      pollingInterval: 4000,
    },
  );

  // 获取分析历史详情 task list 数据
  async function fetchAnalysisTaskList() {
    const res = await Api.submissionRunTaskDetail(
      workspaceId,
      submissionID,
      runID,
      {
        page,
        size,
      },
    );
    if (res.ok) {
      setData(res?.data);
    }
  }

  const taskData = useMemo(
    () =>
      items?.map(item => {
        const { name, startTime, finishTime, duration, stdout, stderr } = item;
        return {
          ...item,
          name: name || '-',
          startTime: genTime(startTime) || '-',
          endTime: genTime(finishTime) || '-',
          duration: convertDuration(duration),
          opration: { stdout, stderr },
        };
      }),
    [items],
  );
  const taskColumns = [
    {
      title: 'Task 名称',
      dataIndex: 'name',
      width: '25%',
      render: value => {
        return (
          <Typography.Text
            style={{ maxWidth: 180 }}
            ellipsis={{
              cssEllipsis: true,
              showTooltip: {
                type: 'popover',
              },
            }}
          >
            {value}
          </Typography.Text>
        );
      },
    },
    {
      title: '运行状态',
      dataIndex: 'status',
      width: '20%',
      render: value => {
        const status = RUN_STATUS_TAG.find(item => item.value === value);
        return <Tag color={status?.color}>{status?.text || ''}</Tag>;
      },
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      width: '25%',
      render: value => value || '-',
    },
    {
      title: (
        <span>
          总耗时
          <Popover
            position="top"
            content={
              <div style={{ maxWidth: 170 }}>
                当前task所有耗时之和，包含启动耗时及分析耗时
              </div>
            }
          >
            <IconQuestionCircle className="ml4" />
          </Popover>
        </span>
      ),
      dataIndex: 'duration',
      width: '15%',
      render: value => value || '-',
    },
    {
      title: '运行详情',
      dataIndex: 'opration',
      width: '15%',
      render: value => {
        if (callCache) {
          return <Link disabled={callCache}>运行日志</Link>;
        }
        return <LogModal {...value} />;
      },
    },
  ];
  useEffect(() => {
    if (isFinished) {
      return stopPolling();
    }
    startPolling();
    return () => stopPolling();
  }, [page, size, submissionID, runID, workspaceId, isFinished]);

  useEffect(() => {
    fetchAnalysisTaskList();
  }, [page, size, submissionID, runID, workspaceId]);

  return (
    <PaginationPage
      total={total}
      sizeOptions={[10, 20, 30, 40, 50, 100]}
      defaultSize={10}
    >
      <Table
        rowKey="name"
        columns={taskColumns}
        data={taskData || []}
        noDataElement={<PageEmpty />}
        pagination={false}
        scroll={{ y: 'calc(100vh - 360px)' }}
      />
    </PaginationPage>
  );
}
