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

import { memo, ReactElement, useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { Button, Input, Popover, Select, Space } from '@arco-design/web-react';
import { IconPlus, IconRefresh } from '@arco-design/web-react/icon';

import { useQuery, useQueryHistory } from 'helpers/hooks';

import Icon from './Icon';

import style from './style.less';

interface Props {
  createText?: string;
  onCreate?: () => void;
  sortOptions?: { label: string; value: string }[];
  onRefresh: () => void;
  afterCreateContent?: ReactElement;
}

type Sort = 'asc' | 'desc';

function TopAction({
  createText,
  onCreate,
  sortOptions,
  onRefresh,
  afterCreateContent,
}: Props) {
  const navigate = useQueryHistory();
  const { pathname } = useLocation();
  const query = useQuery();
  const { sort = 'asc', sortBy } = query;
  const defaultSortBy = sortOptions && sortOptions[0]?.value;

  const [search, setSearch] = useState(query.search);

  useEffect(() => {
    setSearch(query.search);
  }, [query.search]);

  function handleSort(sortValue: Sort, value: string) {
    const queryMap = { ...query, sort: sortValue, sortBy: value };
    navigate(pathname, queryMap);
  }

  function handleSearch(value?: string) {
    const queryMap = { ...query };
    queryMap.search = value;
    queryMap.page = 1;
    navigate(pathname, queryMap);
  }

  return (
    <div className="flexBetween mb16 mt4">
      <Space>
        {createText && (
          <Button type="primary" onClick={onCreate}>
            <IconPlus /> {createText}
          </Button>
        )}
        {afterCreateContent}
        {sortOptions && (
          <div className={`flex ${style.topActionSelect}`}>
            <Select
              style={{ width: 140, borderRadius: '4px 0 0 4px' }}
              defaultValue={sortBy || defaultSortBy}
              onChange={value => {
                handleSort(sort as Sort, value);
              }}
            >
              {sortOptions.map(option => (
                <Select.Option key={option.value} value={option.value}>
                  按 {option.label} 排序
                </Select.Option>
              ))}
            </Select>
            <Popover content={sort === 'asc' ? '点击降序' : '点击升序'}>
              <Button
                className="flexCenter"
                style={{ background: 'white' }}
                icon={
                  sort === 'asc' ? (
                    <Icon glyph="asc" size={14} className="colorPrimary" />
                  ) : (
                    <Icon glyph="desc" size={14} className="colorPrimary" />
                  )
                }
                onClick={() => {
                  handleSort(
                    sort === 'asc' ? 'desc' : 'asc',
                    sortBy || defaultSortBy,
                  );
                }}
              />
            </Popover>
          </div>
        )}
      </Space>
      <Space>
        <Input.Search
          allowClear
          placeholder="请输入名称或描述搜索"
          value={search}
          style={{ width: 264 }}
          onChange={setSearch}
          onSearch={handleSearch}
          onClear={() => handleSearch()}
        />
        <Button
          style={{ background: 'white', padding: '0 10px' }}
          onClick={onRefresh}
        >
          <IconRefresh />
        </Button>
      </Space>
    </div>
  );
}

export default memo(TopAction);
