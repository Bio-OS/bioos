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

import { memo, useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useRequest } from 'ahooks';
import {
  Button,
  Link,
  Message,
  Modal,
  Popover,
  Table,
  Tag,
} from '@arco-design/web-react';

import InputOutputParamete from 'components/analysis/InputOutputModal';
import ListSelects from 'components/analysis/ListSelects';
import CommonLimitText from 'components/CommonLimitText';
import Copy from 'components/Copy';
import PageEmpty from 'components/Empty';
import PaginationPage from 'components/PaginationPage';
import TopAction from 'components/TopAction';
import { RUN_STATUS_TAG, WORKFLOW_STATUS_TO_NAME } from 'helpers/constants';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import { convertDuration, genTime } from 'helpers/utils';
import Api from 'api/client';
import { HandlersListRunsResponse } from 'api/index';

import styles from './analysis.less';

export default memo(function DetailWorkflowList({
  isFinished,
}: {
  isFinished: boolean;
}) {
  const [loading, setLoading] = useState(true);
  const navigate = useQueryHistory();
  const { page, size, search, workflowID, submissionID } = useQuery(10);
  const [status, setStatus] = useState('All');
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [data, setData] = useState<HandlersListRunsResponse>({});
  const { items = [], total } = data;
  const { run: startPolling, cancel: stopPolling } = useRequest(
    getAnalysisRunList,
    {
      pollingInterval: 4000,
    },
  );

  async function getAnalysisRunList() {
    const res = await Api.submissionRunDetail(workspaceId, submissionID, {
      page,
      size,
      searchWord: search ?? '',
      status:
        status === 'All'
          ? [
              'Pending',
              'Running',
              'Succeeded',
              'Failed',
              'Cancelling',
              'Cancelled',
            ]
          : [status],
    });
    if (res.ok) {
      setData(res?.data);
    }
  }

  // 获取分析历史详情 workflow list 数据
  const fetchAnalysisRunList = async () => {
    setLoading(true);
    await getAnalysisRunList();
    setLoading(false);
  };

  // 终止 workflow
  const handleCancelRun = (runId, name) => {
    Api.submissionRunCancelCreate(workspaceId, submissionID, runId)
      .then(() => {
        Message.success(`终止 workflow ${name} 成功`);
        fetchAnalysisRunList();
      })
      .catch(err => {
        Message.error(err?.statusText || '终止 workflow 失败');
        console.error(err);
      });
  };

  const runListData = useMemo(
    () =>
      items?.map(
        ({
          id,
          name,
          status,
          taskStatus,
          startTime,
          finishTime,
          duration,
          inputs,
          outputs,
        } = {}) => ({
          key: id,
          dataEntityRowID: [name, id],
          status: [status, taskStatus],
          startTime: genTime(startTime) || '-',
          endTime: genTime(finishTime) || '-',
          duration: convertDuration(duration),
          input: inputs,
          output: outputs,
          inputOutput: [inputs, outputs],
          operation: [status, id, name],
        }),
      ),
    [items],
  );

  const columns = [
    {
      title: '数据实体',
      dataIndex: 'dataEntityRowID',
      render: ([name, runID]: string[]) => (
        <>
          <Button
            type="text"
            className="ellipsis"
            style={{
              maxWidth: 160,
              padding: 0,
              height: 'auto',
            }}
            onClick={() => {
              navigate(`/workspace/${workspaceId}/analysis/detail/taskDetail`, {
                submissionID,
                runID,
                workflowStatus: status,
                workflowID,
              });
            }}
          >
            {name || 'default'}
          </Button>
          {runID && (
            <div className="colorBlack2 flexAlignCenter">
              <span className="noShrink">运行ID：</span>
              <Copy text={runID} maxWidth={100} />
            </div>
          )}
        </>
      ),
      with: '18%',
    },
    {
      title: '运行状态',
      dataIndex: 'status',
      width: '12%',
      render: ([value, runStatus]: any[]) => {
        const status = RUN_STATUS_TAG.find(item => item.value === value);
        return (
          <Popover
            content={
              <>
                <div>{`启动中：${runStatus?.pending}`}</div>
                <div>{`运行中：${runStatus?.running}`}</div>
                <div>{`运行成功：${runStatus?.succeeded}`}</div>
                <div>{`运行失败：${runStatus?.failed}`}</div>
                <div>{`排队中：${runStatus?.queued}`}</div>
                <div>{`已终止：${runStatus?.cancelled}`}</div>
              </>
            }
            position="right"
          >
            {<Tag color={status?.color}>{status?.text || ''}</Tag>}
          </Popover>
        );
      },
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      width: '18%',
      render: (value: string) => value || '-',
    },
    {
      title: '结束时间',
      dataIndex: 'endTime',
      width: '18%',
      render: (value: string) => value || '-',
    },
    {
      title: '耗时',
      dataIndex: 'duration',
      width: '12%',
      render: (value: string) => value || '-',
    },
    {
      title: '参数',
      dataIndex: 'inputOutput',
      width: '8%',
      render: ([input, output]) => (
        <InputOutputParamete inputs={input} outputs={output} />
      ),
    },
    {
      title: '操作',
      dataIndex: 'operation',
      width: '8%',
      render: ([value, runId, name]: string[]) => {
        let content;
        switch (value) {
          case value:
            content = WORKFLOW_STATUS_TO_NAME[value];
            break;
          default:
            break;
        }
        return value === 'Running' || value === 'Pending' ? (
          <Link
            style={{ padding: 0 }}
            onClick={() => {
              Modal.confirm({
                title: '确定要终止吗？',
                content: (
                  <div style={{ fontSize: 12 }}>
                    终止运行会终止所有运行中的 Task，请谨慎操作
                    <CommonLimitText
                      name="投递"
                      value={runId || 'default'}
                      style={{ marginTop: 12 }}
                    />
                  </div>
                ),
                okButtonProps: {
                  status: 'danger',
                },
                style: {
                  width: 374,
                },
                escToExit: false,
                maskClosable: false,
                okText: '终止',
                onConfirm: () => {
                  handleCancelRun(runId, name);
                },
              });
            }}
          >
            终止
          </Link>
        ) : (
          <Popover content={content} position="left">
            <Link style={{ padding: 0 }} disabled={true}>
              终止
            </Link>
          </Popover>
        );
      },
    },
  ];

  useEffect(() => {
    if (isFinished) {
      return stopPolling();
    }
    startPolling();
    return () => stopPolling();
  }, [page, size, search, status, workspaceId, submissionID, isFinished]);

  useEffect(() => {
    fetchAnalysisRunList();
  }, [page, size, search, status, workspaceId, submissionID]);

  return (
    <div className={styles.detailWorklowlistBox}>
      <div className="fs14 fw500 mb20">工作流运行列表</div>
      <TopAction
        onRefresh={fetchAnalysisRunList}
        afterCreateContent={
          <ListSelects
            statusID={status}
            onChangeStatus={status => setStatus(status)}
          />
        }
      />
      <PaginationPage
        total={total}
        loading={loading}
        sizeOptions={[10, 20, 30, 40, 50, 100]}
        defaultSize={10}
      >
        <Table
          key="0"
          loading={loading}
          columns={columns}
          data={runListData || []}
          noDataElement={<PageEmpty search={search} />}
          pagination={false}
          scroll={{ y: 'calc(100vh - 380px)' }}
        />
      </PaginationPage>
    </div>
  );
});
