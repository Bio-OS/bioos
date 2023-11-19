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
import {
  Button,
  Dropdown,
  Link,
  Menu,
  Message,
  Modal,
  Popover,
  Space,
  Table,
  Tag,
} from '@arco-design/web-react';
import {
  IconCheckCircleFill,
  IconMinusCircleFill,
  IconMore,
  IconSunFill,
} from '@arco-design/web-react/icon';

import { AnalyzeResult } from 'components/analysis/AnalysisStatus';
import ListSelects from 'components/analysis/ListSelects';
import CommonLimitText from 'components/CommonLimitText';
import DeleteModal from 'components/DeleteModal';
import PageEmpty from 'components/Empty';
import { ToastLimitText } from 'components/index';
import ListPage from 'components/ListPage';
import PaginationPage from 'components/PaginationPage';
import TopAction from 'components/TopAction';
import AnalyzeModal from 'components/workflow/run/AnalyzeModal';
import { ANALYSIS_STATUS } from 'helpers/constants';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import { convertDuration, genTime } from 'helpers/utils';
import Api from 'api/client';
import {
  HandlersListSubmissionsResponse,
  HandlersListWorkflowsResponse,
  HandlersSubmissionItem,
} from 'api/index';

export default function AnalysisList() {
  const navigate = useQueryHistory();
  const { page, size, search } = useQuery(10);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState('All');
  const [language, setLanguage] = useState('All');
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [workflowID, setWorkflowID] = useState(''); //空值为全部
  const [deleteItem, setDeleteItem] = useState({
    visible: false,
    id: '',
    name: '',
  });
  const [resetData, setResetData] = useState<HandlersSubmissionItem>({}); // 重新投递 item data
  const [data, setData] = useState<HandlersListSubmissionsResponse>({});
  const [workflowData, setWorkflowData] =
    useState<HandlersListWorkflowsResponse>({});

  const { items = [], total } = data;
  const { run: startPolling, cancel: stopPolling } = useRequest(
    getAnalysisList,
    {
      pollingInterval: 4000,
    },
  );

  async function getAnalysisList() {
    const res = await Api.submissionDetail(workspaceId, {
      orderBy: 'StartTime:desc',
      page,
      size,
      searchWord: search ?? '',
      workflowID: workflowID,
      status:
        status === 'All'
          ? ['Pending', 'Running', 'Finished', 'Cancelling', 'Cancelled']
          : status === 'Running'
          ? ['Pending', 'Running']
          : [status],
      language: language === 'All' ? ['WDL', 'Nextflow'] : [language],
    });

    if (res.ok) {
      setData(res.data);
    }
  }

  // 获取分析历史列表数据
  const fetchAnalysisList = async () => {
    setLoading(true);
    await getAnalysisList();
    setLoading(false);
  };

  // 获取 workflow 列表数据
  const fetchWorkFlowList = () => {
    Api.workspaceIdWorkflowList(workspaceId, { page: 1, size: 100 })
      .then(res => {
        setWorkflowData(res?.data);
      })
      .catch(err => {
        Message.error(err?.statusText || '获取 workflow list 失败');
      });
  };
  // 终止分析历史
  const handleCancelSubmission = (id, name) => {
    Api.submissionCancelCreate(workspaceId, id)
      .then(() => {
        Message.success(`终止分析历史${name}成功`);
        fetchAnalysisList();
      })
      .catch(err => {
        Message.error(err?.statusText || '终止 analysis 失败');
        console.error(err);
      });
  };

  //删除分析历史
  const handleDeleteSubmission = () => {
    Api.submissionDelete(workspaceId, deleteItem.id)
      .then(() => {
        Message.success({
          content: (
            <ToastLimitText
              name={deleteItem.name}
              prefix="删除分析历史"
              suffix="成功"
            />
          ),
        });
        fetchAnalysisList();
        setDeleteItem({
          visible: false,
          id: '',
          name: '',
        });
      })
      .catch(err => {
        console.error(err);
        Message.error(err?.statusText || '删除 analysis 失败');
        setDeleteItem({
          visible: false,
          id: '',
          name: '',
        });
      });
  };

  // 终止
  const stopSubmission = (submissionName: string, submissionID: string) => {
    return (
      <Link
        onClick={() => {
          Modal.confirm({
            title: '确定要终止吗？',
            content: (
              <div style={{ fontSize: 12 }}>
                终止投递将会终止所有启动中和运行中的工作流，请谨慎操作
                <CommonLimitText
                  name={'投递'}
                  value={submissionName}
                  style={{ marginTop: '12px' }}
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
              handleCancelSubmission(submissionID, submissionName);
            },
          });
        }}
      >
        终止
      </Link>
    );
  };

  // 终止中
  const stoppingSubmission = () => {
    return (
      <Popover content={'工作流终止中'} position="left">
        <Link disabled={true}>终止</Link>
      </Popover>
    );
  };

  // 重新投递
  async function createSubmission(data) {
    const { name, description } = data || {};
    const body = {
      ...resetData,
      language: resetData.language,
      workflowID: resetData?.workflowVersion?.id,
      workspaceID: workspaceId,
      name,
      description,
    };

    const res = await Api.submissionCreate(workspaceId, {
      ...body,
    });
    if (res.ok) {
      setResetData({});
      await fetchAnalysisList();
      return res?.data?.id;
    }
    return '';
  }

  // 处理 table data
  const ListSubmissionsData = useMemo(() => {
    return items?.map(item => ({
      key: item?.id,
      Name: [item?.name || '-', item?.description || '-', item?.id || ''],
      Language: item.language,
      Status: [item?.status, item?.runStatus],
      Count: item?.runStatus?.count,
      Time: genTime(item?.startTime) || '-',
      Duration: convertDuration(item?.duration),
      Operator: item,
    }));
  }, [items]);

  const columns = [
    {
      title: '投递名称',
      dataIndex: 'Name',
      width: 280,
      render: ([name, description, submissionID]: string[]) => (
        <Popover
          position="tl"
          content={
            <div
              style={{
                color: 'black',
                fontSize: 12,
                maxHeight: 300,
                overflowY: 'auto',
              }}
            >
              <div>{name}</div>
              <div className="colorBlack2">描述：{description}</div>
            </div>
          }
        >
          <Button
            className="ellipsis"
            type="text"
            style={{
              padding: 0,
              height: 'auto',
              maxWidth: 260,
            }}
            onClick={() =>
              navigate(`/workspace/${workspaceId}/analysis/detail`, {
                submissionID,
              })
            }
          >
            {name}
          </Button>
        </Popover>
      ),
    },
    {
      title: '分析状态',
      dataIndex: 'Status',
      width: 100,
      render: ([value]) => {
        const v = value === 'Pending' ? 'Running' : value;
        const analysisStatus = ANALYSIS_STATUS.find(item => item.value === v);
        const getIcon = () => {
          switch (analysisStatus?.icon) {
            case 'IconCheckCircleFill':
              return <IconCheckCircleFill />;
            case 'IconSunFill':
              return <IconSunFill />;
            case 'IconMinusCircleFill':
              return <IconMinusCircleFill />;
            default:
              break;
          }
        };
        return (
          <Tag color={analysisStatus?.color} icon={getIcon && getIcon()}>
            {analysisStatus?.text || ''}
          </Tag>
        );
      },
    },
    {
      title: '规范',
      width: 100,
      dataIndex: 'Language',
      render: (language: string) => {
        return <Tag color='arcoblue'>{language}</Tag>
      }
    },
    {
      title: '执行结果',
      width: 234,
      dataIndex: 'StatusResult',
      render: (_: undefined, { Status }) => {
        const [value, runStatus] = Status;
        if (['Finished'].includes(value)) {
          return (
            <AnalyzeResult
              data={[
                {
                  text: '成功',
                  count: runStatus.succeeded,
                  color: '#5bcf78',
                },
                {
                  text: '失败',
                  count: runStatus.failed,
                  color: '#ff6b72',
                },
              ]}
              total={runStatus.succeeded + runStatus.failed}
            />
          );
        }

        if (['Cancelled'].includes(value)) {
          return (
            <AnalyzeResult
              data={[
                {
                  text: '成功',
                  count: runStatus.succeeded,
                  color: '#5bcf78',
                },
                {
                  text: '失败',
                  count: runStatus.failed,
                  color: '#ff6b72',
                },
                {
                  text: '已终止',
                  count: runStatus.cancelled,
                  color: '#dde2e9',
                },
              ]}
              total={
                runStatus.succeeded + runStatus.failed + runStatus.cancelled
              }
            />
          );
        }

        return '-';
      },
    },
    {
      title: '数据实体个数',
      width: 110,
      dataIndex: 'Count',
    },
    {
      title: '开始时间',
      width: 170,
      dataIndex: 'Time',
    },
    {
      title: '分析耗时',
      width: 126,
      dataIndex: 'Duration',
    },
    {
      title: '操作',
      dataIndex: 'Operator',
      width: 124,
      render: items => {
        const { status, name, id, workflowVersion } = items;
        return (
          <Space size="mini">
            {(status === 'Finished' || status === 'Cancelled') && (
              <AnalyzeModal
                title="工作流重新投递"
                workflowName={
                  workflowData.items?.find(_ => _.id === workflowVersion.id)
                    ?.name || ''
                }
                onOk={createSubmission}
              >
                {open => (
                  <Link
                    onClick={() => {
                      setResetData(items);
                      open();
                    }}
                  >
                    重新投递
                  </Link>
                )}
              </AnalyzeModal>
            )}
            {(status === 'Pending' || status === 'Running') &&
              stopSubmission(name, id)}
            {status === 'Cancelling' && stoppingSubmission()}
            <Dropdown
              droplist={
                <Menu style={{ minWidth: 64 }}>
                  <Menu.Item
                    key="delete"
                    onClick={() =>
                      setDeleteItem({
                        visible: true,
                        id,
                        name,
                      })
                    }
                  >
                    删除
                  </Menu.Item>
                </Menu>
              }
              position="br"
            >
              <Button
                size="small"
                icon={<IconMore className="fs18 colorBlack4" />}
                iconOnly={true}
                type="text"
              />
            </Dropdown>
          </Space>
        );
      },
    },
  ];

  useEffect(() => {
    fetchWorkFlowList();
  }, []);

  useEffect(() => {
    fetchAnalysisList();
    startPolling();
    return () => stopPolling();
  }, [page, size, search, status, workflowID, language]);

  return (
    <ListPage title="分析历史">
      <TopAction
        onRefresh={fetchAnalysisList}
        afterCreateContent={
          <ListSelects
            statusID={status}
            workflowID={workflowID}
            language={language}
            listWorkFlowItems={workflowData?.items}
            showWorkflowFlag={true}
            showLanguageFlag={true}
            onChangeStatus={status => setStatus(status)}
            onChangeWorkflow={id => setWorkflowID(id)}
            onChangeLanguage={language => setLanguage(language)}
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
          data={ListSubmissionsData}
          noDataElement={<PageEmpty search={search} />}
          pagination={false}
          scroll={{ y: 'calc(100vh - 300px)' }}
        />
      </PaginationPage>
      <DeleteModal
        title="确定删除分析历史吗？"
        type="分析历史"
        name={deleteItem.name}
        tips="删除分析历史将会删除此次分析所产生的中间及最终输出结果数据，并且数据无法恢复。"
        visible={deleteItem.visible}
        onClose={() => {
          setDeleteItem({
            visible: false,
            id: '',
            name: '',
          });
        }}
        onDelete={handleDeleteSubmission}
      />
    </ListPage>
  );
}
