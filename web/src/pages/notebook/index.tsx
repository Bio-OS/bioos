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

import { useEffect, useState } from 'react';
import { useHistory, useLocation, useRouteMatch } from 'react-router-dom';
import { orderBy } from 'lodash-es';
import { Button } from '@arco-design/web-react';

import CardList from 'components/CardList';
import ListPage from 'components/ListPage';
import AddModal from 'components/notebook/AddModal';
import Card from 'components/notebook/Card';
import ConfigureRuntime from 'components/notebook/ConfigureRuntime';
import PaginationPage from 'components/PaginationPage';
import TopAction from 'components/TopAction';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import Api from 'api/client';

export interface Notebook {
  name: string;
  updateTime: number;
  contentLength: number;
}

export interface NotebookListResult {
  ListBucketResult: {
    Contents?: Notebook[] | Notebook;
    Prefix?: string;
  };
}

const SORT_OPTIONS = [
  {
    label: '名称',
    value: 'name',
  },
  {
    label: '更新时间',
    value: 'updateTime',
  },
];

let dataList = [];

export default function () {
  const match = useRouteMatch<{ workspaceId: string }>();
  const [list, setList] = useState<Notebook[]>([]);
  const [addModalVisible, setAddModalVisible] = useState(false);
  const { pathname } = useLocation();
  const navigate = useQueryHistory();
  const [loading, toggleLoading] = useState(false);
  const { search, sortBy, sort, page, size } = useQuery();

  const refreshList = async () => {
    toggleLoading(true);
    Api.workspaceIdNotebookList(match.params.workspaceId)
      .then(({ data }) => {
        dataList = data.items;
        updatePageList();
        toggleLoading(false);
      })
      .catch(error => {
        toggleLoading(false);
      });
  };

  const updatePageList = function () {
    let resultList = dataList;
    if (search) {
      resultList = dataList.filter(({ name }) => {
        return name?.includes(search);
      });
    }
    if (sort && sortBy) {
      resultList = orderBy(resultList, [sortBy], [sort as 'asc' | 'desc']);
    }
    if (page && size) {
      resultList = resultList.slice((page - 1) * size, page * size);
    }
    // console.log(resultList, 'resultList');
    setList(resultList as Notebook[]);
  };

  useEffect(() => {
    refreshList();
    return () => {
      dataList = [];
    };
  }, []);

  useEffect(() => {
    updatePageList();
  }, [search, sortBy, sort, page, size]);
  return (
    <ListPage title="Notebooks">
      <TopAction
        createText="新建notebook"
        onCreate={() => {
          setAddModalVisible(true);
        }}
        onRefresh={() => {
          refreshList();
        }}
        sortOptions={SORT_OPTIONS}
        afterCreateContent={
          <ConfigureRuntime
            render={({ setVisibleConfigureRuntime }) => {
              return (
                <>
                  <Button
                    onClick={() => {
                      setVisibleConfigureRuntime(true);
                    }}
                  >
                    运行资源配置
                  </Button>
                </>
              );
            }}
          ></ConfigureRuntime>
        }
      ></TopAction>

      <PaginationPage
        loading={loading}
        total={search ? list.length : dataList.length}
      >
        <CardList
          data={list}
          loading={loading}
          uniqKey="name"
          renderItem={(notebookProps: Notebook) => (
            <Card refetch={() => refreshList()} {...notebookProps}></Card>
          )}
        />
      </PaginationPage>

      <AddModal
        visible={addModalVisible}
        onClose={() => setAddModalVisible(false)}
        refetch={() => refreshList()}
      />
    </ListPage>
  );
}
/**
 * @param status Server当前状态
 *  Terminated 已停止
    Terminating 停止中
    Pending 启动中
    Running 运行中
    Unknown 未知，可以当成Terminated处理
 */
export function isNotebookServerOk(status?: string) {
  if (status === 'Pending' || status === 'Running') return true;

  return false;
}
