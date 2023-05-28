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

import React, { useEffect, useRef, useState } from 'react';
import { useHistory } from 'react-router-dom';
import classNames from 'classnames';
import dayjs from 'dayjs';
import { Dropdown, Menu, Message, Typography } from '@arco-design/web-react';
import { IconMore } from '@arco-design/web-react/icon';

import CardList from 'components/CardList';
import DeleteModal from 'components/DeleteModal';
import PaginationPage from 'components/PaginationPage';
import TopAction from 'components/TopAction';
import AddModal from 'components/workspace/AddModal';
import { useQuery } from 'helpers/hooks';
import Api from 'api/client';
import { HandlersWorkspaceItem } from 'api/index';
import { HandlersListWorkspacesResponse } from 'api/index';

import styles from './style.less';

interface ModalInfo {
  type?: 'add' | 'edit';
  id?: string;
}
const SORT_OPTIONS = [
  { label: '名称', value: 'Name' },
  { label: '创建时间', value: 'CreatedAt' },
];

const List: React.FC = () => {
  const history = useHistory();
  const [visible, setVisible] = useState(false);
  const [deleteItem, setDeleteItem] = useState({
    visible: false,
    id: '',
    name: '',
  });
  const [data, setData] = useState<HandlersListWorkspacesResponse>({});
  const [loading, setLoading] = useState(false);
  const { page, size, sortBy, sort, search } = useQuery();
  const modalInfo = useRef<ModalInfo>({});

  const { items = [], total } = data;

  function setModalType({ type, id }: ModalInfo) {
    modalInfo.current = { type, id };
    setVisible(true);
  }

  function fetchWorkSpaceList() {
    setLoading(true);
    Api.workspaceList({
      page,
      size,
      searchWord: search ?? undefined,
      orderBy: sortBy ? `${sortBy}:${sort}` : 'Name:asc',
    })
      .then(res => {
        setData(res.data);
        setLoading(false);
      })
      .catch(err => {
        setLoading(false);
        Message.error(err?.statusText || '获取workspace失败');
      });
  }

  function handleDelete() {
    Api.workspaceDelete(deleteItem.id)
      .then(() => {
        Message.success(`删除workspace${deleteItem.name}成功`);
        fetchWorkSpaceList();
        setDeleteItem({
          visible: false,
          id: '',
          name: '',
        });
      })
      .catch(err => {
        console.error(err);
      });
  }

  useEffect(() => {
    fetchWorkSpaceList();
  }, [page, size, search, sortBy, sort]);

  return (
    <div className={styles.workspaceList}>
      <TopAction
        createText="新建 workspace"
        onCreate={() => setModalType({ type: 'add' })}
        onRefresh={fetchWorkSpaceList}
        sortOptions={SORT_OPTIONS}
      />
      <PaginationPage total={total} loading={loading}>
        <CardList
          data={items}
          gutter={[32, 32]}
          renderItem={({
            id,
            name,
            description,
            createTime,
          }: HandlersWorkspaceItem) => {
            return (
              <div
                className={styles.workspaceItem}
                onClick={() => {
                  history.push(`/workspace/${id}/data`);
                }}
              >
                <Typography.Text
                  className={classNames([
                    'fs14 colorPrimary lh22 fw500',
                    styles.title,
                  ])}
                  ellipsis={{
                    rows: 4,
                  }}
                >
                  {name}
                </Typography.Text>
                <Typography.Paragraph
                  className="fs12 lh20 fw400 flex1"
                  style={{
                    pointerEvents: 'none',
                  }}
                  ellipsis={{
                    rows: 5,
                  }}
                >
                  {description}
                </Typography.Paragraph>
                <div className="flexBetween">
                  <span className="colorGrey">
                    {dayjs(createTime * 1000).format('YYYY-MM-DD')}
                  </span>
                  <Dropdown
                    droplist={
                      <Menu>
                        <Menu.Item
                          key="edit"
                          onClick={e => {
                            e.stopPropagation();
                            setModalType({ type: 'edit', id });
                          }}
                        >
                          编辑
                        </Menu.Item>
                        <Menu.Item
                          key="delete"
                          onClick={e => {
                            e.stopPropagation();
                            setDeleteItem({
                              visible: true,
                              id,
                              name,
                            });
                          }}
                        >
                          删除
                        </Menu.Item>
                      </Menu>
                    }
                    position="bl"
                  >
                    <IconMore
                      fontSize={20}
                      className={styles.svgMore}
                      onClick={e => e.stopPropagation()}
                    />
                  </Dropdown>
                </div>
              </div>
            );
          }}
        />
      </PaginationPage>
      <AddModal
        type={modalInfo.current.type}
        workspaceId={modalInfo.current.id}
        visible={visible}
        onClose={() => {
          modalInfo.current = {};
          setVisible(false);
        }}
        refetch={fetchWorkSpaceList}
      />
      <DeleteModal
        type=" Workspace "
        width={600}
        name={deleteItem.name}
        visible={deleteItem.visible}
        verify={true}
        useTextArea={true}
        onClose={() =>
          setDeleteItem({
            visible: false,
            id: '',
            name: '',
          })
        }
        onDelete={handleDelete}
      />
    </div>
  );
};

export default List;
