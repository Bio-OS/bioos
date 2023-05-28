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
import { Button, Dropdown, Menu, Message, Tag } from '@arco-design/web-react';
import {
  IconCheckCircleFill,
  IconMinusCircleFill,
  IconMore,
  IconSunFill,
} from '@arco-design/web-react/icon';

import AnalysisOverview from 'components/analysis/AnalysisOverview';
import AnalysisStatus from 'components/analysis/AnalysisStatus';
import DetailRightContainer from 'components/analysis/DetailRightContainer';
import DetailWorkflowList from 'components/analysis/DetailWorkflowList';
import Breadcrumbs from 'components/Breadcrumbs';
import DeleteModal from 'components/DeleteModal';
import DetailPage from 'components/DetailPage';
import { ToastLimitText } from 'components/index';
import { ANALYSIS_STATUS } from 'helpers/constants';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import Api from 'api/client';

export default function AnalysisDetail() {
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const { submissionID } = useQuery();
  const navigate = useQueryHistory();
  const [deleteItem, setDeleteItem] = useState({
    visible: false,
    id: '',
    name: '',
  });
  const [detailData, setDetailData] = useState<any>({});
  const { run: startPolling, cancel: stopPolling } = useRequest(
    getAnalysisDetail,
    {
      pollingInterval: 4000,
    },
  );

  const isFinished = ['Finished', 'Cancelled'].includes(detailData.status);

  async function getAnalysisDetail() {
    const res = await Api.submissionDetail(workspaceId, {
      ids: [submissionID],
    });
    if (res.ok) {
      const item = res?.data?.items?.[0];
      setDetailData(item);
      if (['Finished', 'Cancelled'].includes(item.status)) stopPolling();
    }
  }

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
        setDeleteItem({
          visible: false,
          id: '',
          name: '',
        });
        navigate(`/workspace/${workspaceId}/analysis`);
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

  const statusData = useMemo(() => {
    return [
      {
        name: '运行成功',
        value: detailData?.runStatus?.succeeded,
        style: { color: '#1CB267' },
      },
      {
        name: '运行中',
        value: detailData?.runStatus?.running,
        style: { color: PRIMARY_COLOR },
      },
      {
        name: '运行失败',
        value: detailData?.runStatus?.failed,
        style: { color: '#DB373F' },
      },
      { name: '启动中', value: detailData?.runStatus?.pending },
      { name: '终止中', value: detailData?.runStatus?.cancelling },
      { name: '已终止', value: detailData?.runStatus?.cancelled },
    ];
  }, [detailData]);

  const status = ANALYSIS_STATUS.find(
    item =>
      item.value ===
      (detailData?.status === 'Pending' ? 'Running' : detailData?.status),
  );

  //获取 icon 组件
  const getIcon = () => {
    switch (status?.icon) {
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

  useEffect(() => {
    startPolling();
    return () => stopPolling();
  }, [workspaceId, submissionID]);

  return (
    <DetailPage
      breadcrumbs={
        <Breadcrumbs
          breadcrumbs={[
            { text: '分析历史', path: `/workspace/${workspaceId}/analysis` },
            { text: '分析历史详情' },
          ]}
        />
      }
      title={detailData?.name}
      statusTag={
        <Tag color={status?.color} icon={getIcon && getIcon()}>
          {status?.text || ''}
        </Tag>
      }
      contentStyle={{ paddingRight: 0, paddingBottom: 0 }}
      rightArea={
        <Dropdown
          droplist={
            <Menu style={{ minWidth: 64 }}>
              <Menu.Item
                key="delete"
                onClick={() =>
                  setDeleteItem({
                    visible: true,
                    id: detailData?.id,
                    name: detailData?.name,
                  })
                }
              >
                删除
              </Menu.Item>
            </Menu>
          }
          position="br"
        >
          <Button icon={<IconMore className="fs18 mt4" />} />
        </Dropdown>
      }
    >
      <div className="flex h100">
        <DetailWorkflowList isFinished={isFinished} />
        <DetailRightContainer>
          <AnalysisStatus statusData={statusData} />
          <AnalysisOverview
            name={detailData?.name}
            description={detailData?.description ?? '-'}
            id={detailData?.id}
            workspaceId={workspaceId}
            workflowId={detailData?.workflowVersion?.id}
            flagReadFromCache={detailData?.exposedOptions?.readFromCache}
          />
        </DetailRightContainer>
      </div>
      <DeleteModal
        title="确定删除分析历史吗？"
        type="分析历史"
        name={detailData?.name}
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
    </DetailPage>
  );
}
