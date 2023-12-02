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

import { memo } from 'react';
import { useRouteMatch } from 'react-router-dom';
import classNames from 'classnames';
import {
  Card,
  Dropdown,
  Menu,
  Popover, Tag,
  Typography,
} from '@arco-design/web-react';
import {
  IconExclamationCircleFill,
  IconLoading,
  IconMore,
  IconPlayCircleFill,
} from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import { useQueryHistory } from 'helpers/hooks';

interface Props {
  id: string;
  name: string;
  language: string;
  status: string;
  description: string;
  originUrl: string;
  onEdit: () => void;
  onDelete: () => void;
  onReImport: () => void;
}

function WorkflowCard({
  id,
  name,
  language,
  status,
  description,
  originUrl,
  onEdit,
  onDelete,
  onReImport,
}: Props) {
  const navigate = useQueryHistory();
  const match = useRouteMatch<{ workspaceId: string }>();

  const importing = status === 'Pending';
  const failed = status === 'Failed';
  const succeeded = status === 'Success';

  function renderStatus() {
    if (importing) {
      return <IconLoading className="colorPrimary" />;
    }

    if (failed) {
      return (
        <Popover content="导入失败，建议重新导入">
          <IconExclamationCircleFill style={{ color: '#fa9600' }} />
        </Popover>
      );
    }
    return null;
  }

  function renderAction() {
    return (
      <Dropdown
        disabled={importing}
        droplist={
          <Menu style={{ width: 64 }}>
            <Menu.Item
              key="edit"
              onClick={e => {
                e.stopPropagation();
                onEdit();
              }}
            >
              更新
            </Menu.Item>

            <Menu.Item
              key="delete"
              onClick={e => {
                e.stopPropagation();
                onDelete();
              }}
            >
              删除
            </Menu.Item>
          </Menu>
        }
        position="br"
      >
        <span
          className={classNames([
            'flexCenter',
            importing ? 'notAllowed' : 'iconBox',
          ])}
          onClick={e => e.stopPropagation()}
        >
          <IconMore />
        </span>
      </Dropdown>
    );
  }

  function handleClick() {
    if (!succeeded) return;
    navigate(`${match.url}/${id}/run`);
  }

  return (
    <Card
      hoverable
      className={classNames(['br8', { hoverableArea: succeeded }])}
      onClick={handleClick}
    >
      <div className="flexBetween mb8">
        <div
          className={classNames([
            'flexAlignCenter',
            {
              notAllowed: !succeeded,
            },
          ])}
        >
          <Typography.Paragraph
            className="fs16 fw500 mr8 hoverActive"
            style={{ maxWidth: 180 }}
            ellipsis={{
              cssEllipsis: true,
              showTooltip: { type: 'popover' },
            }}
          >
            {name}
          </Typography.Paragraph>
          <span className="noShrink">{renderStatus()}</span>
        </div>
        {renderAction()}
      </div>

      <Typography.Paragraph
        className="colorBlack2 fs12 inlineBlock lh20 w100"
        style={{ height: 20 }}
        ellipsis={{
          cssEllipsis: true,
          showTooltip: { type: 'popover' },
        }}
      >
        {description}
      </Typography.Paragraph>
      <div className="flexBetween mt24 w100 fs12">
        <span className="noShrink colorGrey">规范：</span>
        <Typography.Paragraph
            className="colorBlack2 fs12 inlineBlock w100"
            // style={{ height: 20 }}
        >
          <Tag color="arcoblue">{language}</Tag>
        </Typography.Paragraph>
        <span className="noShrink colorGrey">来源：</span>
        <a
          className="colorGrey ellipsis hoverableText mr16"
          target="_blank"
          href={originUrl}
          style={{ width: 'calc(100% - 60px)' }}
        >
          {originUrl}
        </a>
        {failed ? (
          <span
            className="noShrink colorPrimary cursorPointer flexAlignCenter"
            onClick={onReImport}
          >
            <span className="mr4">重新导入</span>
            <Icon glyph="linkto" className="colorPrimary" size={12} />
          </span>
        ) : (
          <IconPlayCircleFill className="fs20 colorBg noShrink" />
        )}
      </div>
    </Card>
  );
}
export default memo(WorkflowCard);
