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

import { ReactNode, useEffect, useMemo, useRef, useState } from 'react';
import { useRouteMatch } from 'react-router-dom';
import {
  Input,
  Link,
  Modal,
  Spin,
  Table,
  TableColumnProps,
} from '@arco-design/web-react';

import SetDetailModal from 'components/data/modal/SetDetailModal';
import {
  LinkRender,
  TextEllipsis,
  TextRender,
} from 'components/data/ModelContent';
import PageEmpty from 'components/Empty';
import { Line } from 'components/index';
import Api from 'api/client';
import {
  HandlersDataModel,
  HandlersListDataModelRowsResponse,
} from 'api/index';

export default function SelectDataModelModal({
  dataModelInfo = {},
  selectRows = [],
  children,
  onChange,
}: {
  dataModelInfo: HandlersDataModel;
  selectRows: string[];
  children: (open: () => void) => ReactNode;
  onChange: (selectedRowKeys: string[]) => void;
}) {
  const [visible, setVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [sorter, setSorter] = useState('');
  const [setValue, setSetValue] = useState<string[] | null>();
  const [data, setData] = useState<HandlersListDataModelRowsResponse>({});
  const [selectedRowKeys, setSelectedRowKeys] = useState(selectRows);
  const [pagination, setPagination] = useState<{ page: number; size: number }>({
    page: 1,
    size: 10,
  });
  const { headers = [], rows = [], total = 0 } = data;
  const match = useRouteMatch<{ workspaceId: string }>();
  const { workspaceId } = match.params;
  const allRowKeys = useRef([]);
  const isSelectAll = allRowKeys.current.length === selectedRowKeys.length;

  function open() {
    setVisible(true);
  }

  function close() {
    setSearch('');
    setSetValue(null);
    setSorter('');
    setPagination({
      page: 1,
      size: 10,
    });
    setVisible(false);
  }

  function confirm() {
    close();
    onChange(selectedRowKeys);
  }

  const columns: TableColumnProps[] = useMemo(() => {
    return headers.map((item, index) => {
      const isSet = dataModelInfo.type === 'entity_set';
      return {
        dataIndex: item,
        title: <TextEllipsis text={item} />,
        fixed: index === 0 ? 'left' : undefined,
        width: 200,
        sorter: !isSet || (isSet && index === 0),
        render: (value: string) => {
          if (value === '-') return value;
          return isSet ? (
            <LinkRender
              value={value}
              index={index}
              search={search}
              name={dataModelInfo.name}
              onClick={setSetValue}
            />
          ) : (
            <TextRender value={value} index={index} />
          );
        },
      };
    });
  }, [headers]);

  const rowsData = useMemo(() => {
    return rows?.map(row => {
      return row.reduce((cur, next, index) => {
        return {
          ...cur,
          [headers[index]]: next || '-',
        };
      }, {});
    });
  }, [rows]);

  async function fetchDataModel() {
    setLoading(true);
    const res = await Api.dataModelRowsDetail(workspaceId, dataModelInfo.id, {
      page: pagination.page,
      size: pagination.size,
      searchWord: search || undefined,
      orderBy: sorter || undefined,
    });
    setLoading(false);
    setData(res.data);
  }

  async function fetchAllRowsId() {
    const res = await Api.dataModelRowsIdsDetail(workspaceId, dataModelInfo.id);
    allRowKeys.current = res.data.rowIDs;
  }

  function handleSelectAll() {
    isSelectAll
      ? setSelectedRowKeys([])
      : setSelectedRowKeys(allRowKeys.current);
  }

  useEffect(() => {
    if (!dataModelInfo.id) return;
    if (visible) {
      fetchAllRowsId();
    }
  }, [dataModelInfo, visible]);

  useEffect(() => {
    if (!dataModelInfo.id) return;
    if (visible) {
      fetchDataModel();
    }
  }, [dataModelInfo, visible, search, pagination, sorter]);

  useEffect(() => {
    if (!dataModelInfo.id) return;
    if (visible) {
      setSelectedRowKeys(selectRows);
    }
  }, [visible]);

  return (
    <>
      {children(open)}
      <Modal
        visible={visible}
        title="选择数据"
        style={{ width: 800 }}
        onOk={confirm}
        onCancel={close}
        autoFocus={false}
        unmountOnExit={true}
        maskClosable={false}
        escToExit={false}
      >
        <div className="flexBetween mb16">
          <div className="flexAlignCenter">
            <Link onClick={handleSelectAll}>
              {isSelectAll ? '取消选择全部' : '选择全部实体'}
            </Link>
            <Line />
            <div className="colorGrey fs12 lh20">
              <div>已选{selectedRowKeys.length}项</div>
            </div>
          </div>
          <Input.Search
            allowClear
            placeholder="输入数据搜索"
            style={{ width: 240 }}
            onSearch={value => {
              setSearch(value);
              setPagination({ page: 1, size: pagination.size });
            }}
            onClear={() => setSearch('')}
          />
        </div>
        {!headers.length ? (
          <div className="w100 textAlignCenter" style={{ padding: '80px 0' }}>
            <Spin />
          </div>
        ) : (
          <Table
            loading={loading}
            showSorterTooltip={false}
            tableLayoutFixed
            style={{ width: '100%' }}
            columns={columns}
            data={rowsData}
            rowKey={headers[0]}
            scroll={{
              x: true,
              y: 430,
            }}
            rowSelection={{
              type: 'checkbox',
              selectedRowKeys,
              checkCrossPage: true,
              onChange(keys) {
                const hiddenKeys =
                  selectedRowKeys?.filter(
                    key => !rowsData.find(item => item[headers[0]] === key),
                  ) || [];
                setSelectedRowKeys((keys as string[]).concat(hiddenKeys));
              },
            }}
            pagination={{
              size: 'small',
              hideOnSinglePage: pagination.size === 10,
              sizeCanChange: true,
              current: pagination.page,
              pageSize: pagination.size,
              total,
              showTotal: true,
              sizeOptions: [10, 20, 30, 40, 50, 100],
              onChange(pageNumber, sizeNumber) {
                setPagination({ page: pageNumber, size: sizeNumber });
              },
            }}
            onChange={(_pagination, sorter, _filter, extra) => {
              if (extra.action !== 'sort') return;
              setSorter(
                `${sorter.field}:${
                  sorter.direction === 'descend' ? 'desc' : 'asc'
                }`,
              );
            }}
            noDataElement={<PageEmpty search={search} />}
          />
        )}
      </Modal>
      <SetDetailModal
        data={setValue || []}
        keyword={search}
        name={dataModelInfo.name?.replace('_set', '')}
        visible={Boolean(setValue)}
        onOK={() => setSetValue(undefined)}
        onCancel={() => setSetValue(undefined)}
        title={!dataModelInfo.name?.includes('_set_set') ? '实体' : '实体集合'}
      />
    </>
  );
}
