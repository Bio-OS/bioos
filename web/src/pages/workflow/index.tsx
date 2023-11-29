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

import { useEffect, useRef, useState } from 'react';
import { useRouteMatch } from 'react-router-dom';
import { useRequest } from 'ahooks';
import { Message } from '@arco-design/web-react';

import CardList from 'components/CardList';
import DeleteModal from 'components/DeleteModal';
import { ToastLimitText } from 'components/index';
import ListPage from 'components/ListPage';
import PaginationPage from 'components/PaginationPage';
import TopAction from 'components/TopAction';
import Card from 'components/workflow/Card';
import ImportModal from 'components/workflow/ImportModal';
import { useQuery } from 'helpers/hooks';
import Api from 'api/client';
import { HandlersListWorkflowsResponse, HandlersWorkflowItem } from 'api/index';

const SORT_OPTIONS = [
  { label: '名称', value: 'name' },
  { label: '创建时间', value: 'createdAt' },
];

export default function WorkflowList() {
  const currentInfo = useRef(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [data, setData] = useState<HandlersListWorkflowsResponse>(null);
  const [importVisible, setImportVisible] = useState(false);
  const [deleteVisible, setDeleteVisible] = useState(false);
  const { page, size, search, sort, sortBy } = useQuery();

  const { run: startPolling, cancel: stopPolling } = useRequest(
    getWorkflowList,
    {
      pollingInterval: 4000,
    },
  );
  const match = useRouteMatch<{ workspaceId: string }>();
  const { workspaceId } = match.params;

  function hideDeleteModal() {
    currentInfo.current = null;
    setDeleteVisible(false);
  }

  function hideImportModal() {
    currentInfo.current = null;
    setImportVisible(false);
  }

  async function getWorkflowList() {
    const res = await Api.workspaceIdWorkflowList(workspaceId, {
      page,
      size,
      searchWord: search || undefined,
      orderBy: sortBy ? `${sortBy}:${sort}` : 'name:asc',
    });

    if (res.ok) {
      setData(res.data);
    }
  }

  async function fetch() {
    setLoading(true);
    await getWorkflowList();
    setLoading(false);
  }

  async function handleDelete() {
    const res = await Api.workspaceIdWorkflowDelete(
      workspaceId,
      currentInfo.current.id,
    );
    if (res.ok) {
      Message.success({
        content: (
          <ToastLimitText
            name={currentInfo.current.name}
            prefix="删除"
            suffix="成功"
          />
        ),
      });
      hideDeleteModal();
      fetch();
    }
  }

  useEffect(() => {
    fetch();
    startPolling();
    return () => stopPolling();
  }, [page, size, search, sort, sortBy]);

  return (
    <ListPage title="工作流">
      <TopAction
        createText="导入工作流"
        onCreate={() => setImportVisible(true)}
        onRefresh={fetch}
        sortOptions={SORT_OPTIONS}
      />
      <PaginationPage total={data?.total} loading={loading}>
        <CardList
          data={data?.items}
          renderItem={(item: HandlersWorkflowItem) => {
            return (
              <Card
                id={item.id}
                name={item.name}
                status={item.latestVersion.status}
                description={item.description}
                language={item.latestVersion.language}
                originUrl={item.latestVersion.metadata.gitURL}
                onEdit={() => {
                  currentInfo.current = item;
                  setImportVisible(true);
                }}
                onDelete={() => {
                  currentInfo.current = item;
                  setDeleteVisible(true);
                }}
                onReImport={() => {
                  currentInfo.current = item;
                  setImportVisible(true);
                }}
              />
            );
          }}
        />
      </PaginationPage>
      <ImportModal
        visible={importVisible}
        workflowInfo={currentInfo.current}
        refetch={fetch}
        onClose={hideImportModal}
      />
      <DeleteModal
        type="工作流"
        name={currentInfo.current?.name}
        visible={deleteVisible}
        tips={[
          '若工作流正在运行中，删除时将会自动停止运行工作流；',
          '删除工作流将会删除此次工作流所对应的分析历史及其输出结果数据，并且数据无法恢复；',
        ]}
        onClose={hideDeleteModal}
        onDelete={handleDelete}
      />
    </ListPage>
  );
}
