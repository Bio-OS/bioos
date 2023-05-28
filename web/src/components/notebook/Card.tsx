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

import { memo, ReactNode, useRef, useState } from 'react';
import { useHistory, useRouteMatch } from 'react-router-dom';
import classNames from 'classnames';
import {
  Button,
  Card,
  CardProps,
  Dropdown,
  Menu,
  Message,
  Typography,
} from '@arco-design/web-react';
import { IconMore } from '@arco-design/web-react/icon';
import { IconClockCircle } from '@arco-design/web-react/icon';

import DeleteModal from 'components/DeleteModal';
import { downloadFileByBlob, genTime } from 'helpers/utils';
import Api from 'api/client';

import CopyModal from './CopyModal';

interface Notebook {
  name: string;
  updateTime: number;
  contentLength: number;
  refetch: () => void;
}

function NotebookCard({ updateTime, contentLength, name, refetch }: Notebook) {
  const history = useHistory();
  const match = useRouteMatch<{ workspaceId: string }>();
  const [copyName, setCopyName] = useState<string>();
  const [deleteItem, setDeleteItem] = useState({
    visible: false,
    name: '',
  });
  const handleDownload = () => {
    Api.workspaceIdNotebookDetail(match.params.workspaceId, name, {
      format: 'blob',
    })
      .then(({ data }) => {
        downloadFileByBlob(data, `${name}.ipynb`);
      })
      .catch(() => {
        Message.error('下载失败，文件不存在。');
      });
  };

  const handleDelete = () => {
    Api.workspaceIdNotebookDelete(
      match.params.workspaceId,
      deleteItem.name,
    ).then(res => {
      if (res.data === null) {
        Message.success(`删除${deleteItem.name}成功`);
        setDeleteItem({ visible: false, name: null });
        refetch();
      }
    });
  };

  const renderAction = () => {
    return (
      <>
        <Menu.Item
          key="edit"
          onClick={() => {
            setCopyName(name);
          }}
        >
          复制到workspace
        </Menu.Item>
        <Menu.Item key="download" onClick={() => handleDownload()}>
          下载ipynb文件
        </Menu.Item>
        <Menu.Item
          key="delete"
          onClick={() => setDeleteItem({ visible: true, name: name })}
        >
          删除
        </Menu.Item>
      </>
    );
  };
  const renderCardAction = () => {
    return (
      <span onClick={e => e.stopPropagation()}>
        <Dropdown
          unmountOnExit={false}
          droplist={<Menu style={{ minWidth: 64 }}>{renderAction()}</Menu>}
          position="br"
        >
          <span className="flexCenter iconBox">
            <IconMore />
          </span>
        </Dropdown>
      </span>
    );
  };

  return (
    <>
      <Card
        hoverable={true}
        bodyStyle={{ padding: '16px 20px' }}
        onClick={() => {
          history.push(`${match.url}/${name}`);
        }}
        className="cursorPointer br8"
      >
        <>
          <div className="flexJustifyBetween mb8">
            <div className="mr16 flexAlignCenter">
              <Typography.Paragraph
                className={classNames('fw500 fs14')}
                ellipsis={{
                  showTooltip: { type: 'popover' },
                }}
              >
                {name}
              </Typography.Paragraph>
            </div>

            {renderCardAction()}
          </div>
          <div className="flexAlignCenter mt24 fs12 colorGrey">
            <IconClockCircle fontSize={14} className="mr8" />
            更新时间: {genTime(updateTime)}
          </div>
        </>
      </Card>
      <CopyModal
        name={copyName}
        onClose={() => setCopyName(undefined)}
        refetch={refetch}
      />
      <DeleteModal
        type=" Notebook "
        tips="删除后，将同步删除 Notebook 存储的 .ipynb 后缀文件，请谨慎操作。"
        name={deleteItem.name}
        visible={deleteItem.visible}
        title="确定删除 Notebook 吗？"
        onClose={() =>
          setDeleteItem({
            visible: false,
            name: '',
          })
        }
        onDelete={handleDelete}
      />
    </>
  );
}

export default memo(NotebookCard);
