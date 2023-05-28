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

import React, {
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef,
  useState,
} from 'react';
import { useParams } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { last } from 'lodash-es';
import Papa from 'papaparse';
import { v4 } from 'uuid';
import {
  Button,
  Divider,
  Dropdown,
  Input,
  Link,
  Menu,
  Message,
  Popconfirm,
  Popover,
  Space,
  Spin,
  Table,
  Typography,
} from '@arco-design/web-react';
import { IconMore, IconMoreVertical } from '@arco-design/web-react/icon';
import { useIntersectionObserver } from '@asyarb/use-intersection-observer';

import Copy from 'components/Copy';
import PaginationPage from 'components/PaginationPage';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import { downloadFile } from 'helpers/utils';
import Api, { apiInstance } from 'api/client';
import {
  HandlersDataModel,
  HandlersListAllDataModelRowIDsResponse,
  HandlersListDataModelRowsResponse,
} from 'api/index';

import PageEmpty from '../Empty';

import CreateSetModal from './modal/CreateSetModal';
import DeleteSetModal from './modal/DeleteSetModal';
import SetDetailModal from './modal/SetDetailModal';

import styles from './ModelContent.less';

interface Props {
  model?: HandlersDataModel;
  includeSet?: HandlersDataModel;
  refresh?: () => void;
  entityList?: HandlersDataModel[];
}
interface Column {
  dataIndex: string;
  title: JSX.Element;
  fixed?: 'left' | 'right';
  width: number;
  sorter: boolean;
  filterIcon: JSX.Element;
  filterDropdown: ({ confirm }: { confirm: any }) => JSX.Element;
  render: (value: string) => JSX.Element | '-';
}
const defaultModel = {
  headers: [],
  total: 0,
  rows: [],
};
const ModelContent: React.ForwardRefRenderFunction<
  {
    updateModel: () => void;
  },
  Props
> = (props, ref) => {
  const { model, includeSet, refresh, entityList } = props;
  const navigate = useQueryHistory();
  const { pathname } = useLocation();
  const query = useQuery();
  const { search, page, size, modelId } = query;
  const tokenRef = useRef(null);
  const [innerSearch, setInnerSearch] = useState('');
  const { workspaceId } = useParams<{ workspaceId: string }>();
  // 内部状态
  const [sorter, setSorter] = useState('');
  const [loading, setLoading] = useState(false);
  const [modelData, setModelData] =
    useState<HandlersListDataModelRowsResponse>(defaultModel);
  const [modelRowsId, setModelRowsId] =
    useState<HandlersListAllDataModelRowIDsResponse>();
  const [selectedRowKeys, setSelectedRowKeys] = useState<(string | number)[]>(
    [],
  );
  const isSelectAll = modelRowsId?.rowIDs?.length === selectedRowKeys?.length;
  const [setValue, setSetValue] = useState<string[] | null>();
  const [collectionVisible, setCollectionVisible] = useState(false);
  const [deleteVisible, setDeleteVisible] = useState(false);
  const { headers = [], total = 0, rows = [] } = modelData;
  useImperativeHandle(ref, () => {
    return {
      updateModel: () => {
        setSorter('');
        setModelRowsId(null);
        setModelData(defaultModel);
      },
    };
  });
  const fetchAllDataModelRowIDs = id => {
    Api.dataModelRowsIdsDetail(workspaceId, id)
      .then(res => {
        if (res?.ok) {
          setModelRowsId(res.data);
        }
      })
      .catch(e => {
        const info = e?.error?.message || e?.statusText;
        info && Message.error(info);
      });
  };
  const fetchDataModelRows = (id, page, size, search, sorter) => {
    setLoading(true);
    if (tokenRef.current) {
      apiInstance.abortRequest(tokenRef.current);
    }
    tokenRef.current = v4();
    Api.dataModelRowsDetail(
      workspaceId,
      id,
      {
        page,
        size,
        searchWord: search,
        orderBy: sorter,
      },
      {
        cancelToken: tokenRef.current,
      },
    )
      .then(res => {
        res?.ok && setModelData(res.data || defaultModel);
      })
      .catch(e => {
        const info = e?.error?.message || e?.statusText;
        info && Message.error(info);
      })
      .finally(() => {
        setLoading(false);
        tokenRef.current = null;
      });
  };
  useEffect(() => {
    setSelectedRowKeys([]);
  }, [model]);
  useEffect(() => {
    if (!modelId || !model) return;
    setInnerSearch(search);
    if (modelId === '-1') {
      setModelData({
        headers: ['Key', 'Value'],
        total: 0,
        rows: [],
      });
      return;
    }
    fetchDataModelRows(modelId, page, size, search, sorter);
    fetchAllDataModelRowIDs(modelId);
    return () => {
      if (tokenRef.current) {
        apiInstance.abortRequest(tokenRef.current);
      }
    };
  }, [model, modelId, page, size, search, sorter]);
  const handleSearch = (value?: string) => {
    const queryMap = { ...query, search: value, page: 1, size: 10 };
    navigate(pathname, queryMap);
  };
  const handleDownLoadCsv = () => {
    Api.dataModelRowsDetail(workspaceId, model.id, {
      size: selectedRowKeys.length,
      page: 1,
      rowIDs: selectedRowKeys as string[],
    })
      .then(res => {
        if (res?.ok && res?.data) {
          const csv = Papa.unparse([res?.data?.headers, ...res?.data?.rows]);
          const blob = new Blob([`\ufeff${csv}`]);
          downloadFile(window.URL.createObjectURL(blob), `${model.name}.csv`);
        }
      })
      .catch(e => {
        Message.error(e?.error?.message || e?.statusText || '下载失败');
      });
  };
  const handleDeleteRow = () => {
    if (isSelectAll) {
      setDeleteVisible(false);
      refresh?.();
    } else {
      setDeleteVisible(false);
      refresh?.();
      fetchDataModelRows(model.id, page, size, search, sorter);
      fetchAllDataModelRowIDs(model.id);
    }
    setSelectedRowKeys([]);
  };
  const handleDeleteColumn = async (header: string) => {
    try {
      const res = await Api.dataModelDelete(workspaceId, model.id, {
        headers: [header],
      });
      if (res?.ok) {
        Message.success(`删除数据${header}成功`);
        fetchDataModelRows(model.id, page, size, search, sorter);
        fetchAllDataModelRowIDs(model.id);
      }
    } catch (_) {
      Message.error(`删除数据${header}失败`);
    }
  };
  const { columns, data } = useMemo(() => {
    if (!model) {
      return {
        columns: [],
        data: [],
      };
    }
    const columns: Column[] = headers?.map((item, index) => {
      const isFilter = model.type === 'entity' && index !== 0;
      const isWorkSpace = model.type === 'workspace';
      const isSet = model.type === 'entity_set';
      const basicColumn = {
        dataIndex: item,
        title: <TextEllipsis text={item} />,
        width: 200,
        sorter:
          model.type === 'entity' ||
          (model.type === 'entity_set' && index === 0),
        filterIcon: isFilter ? <IconMoreVertical /> : null,
        filterDropdown: isFilter ? filterDown(item, handleDeleteColumn) : null,
        render: (value: string) => {
          if (!value) return '-';
          if (isWorkSpace) return <TextEllipsis text={value} />;
          return isSet ? (
            <LinkRender
              value={value}
              index={index}
              search={search}
              name={model?.name}
              onClick={val => {
                setSetValue(val);
              }}
            />
          ) : (
            <TextRender value={value} index={index} />
          );
        },
      };
      return {
        ...basicColumn,
        ...(index === 0
          ? {
              fixed: 'left',
            }
          : {}),
      };
    });
    const data = rows?.map(row => {
      return row?.reduce((cur, next, index) => {
        return {
          ...cur,
          [headers[index]]: next,
        };
      }, {});
    });
    return {
      columns,
      data,
    };
  }, [headers, rows, model, search]);
  if (!model) {
    return (
      <PageEmpty
        desc="暂无实体数据模型，快去添加一个吧"
        style={{
          marginTop: 229,
        }}
      />
    );
  }
  return (
    <div className={styles.modelTable}>
      <Spin loading={loading} block>
        <div className="colorBlack mb16 fs16 fw500">{model?.name}</div>
        <div className="flexBetween mb16">
          <Space>
            <Button
              onClick={() => {
                if (isSelectAll) {
                  setSelectedRowKeys([]);
                } else {
                  setSelectedRowKeys(modelRowsId?.rowIDs);
                }
              }}
              disabled={!data?.length || !modelRowsId?.rowIDs.length}
            >
              {isSelectAll ? '取消全部' : '选择全部'}
            </Button>
            <Button
              disabled={!selectedRowKeys?.length}
              onClick={handleDownLoadCsv}
            >
              下载
            </Button>
            {model?.type !== 'workspace' && (
              <Button
                disabled={selectedRowKeys?.length < 2}
                onClick={() => {
                  setCollectionVisible(true);
                }}
              >
                生成实体集合
              </Button>
            )}
            <Dropdown
              droplist={
                <Menu>
                  <Menu.Item
                    key="delete"
                    onClick={e => {
                      e.stopPropagation();
                      setDeleteVisible(true);
                    }}
                    disabled={!selectedRowKeys?.length}
                  >
                    <span className="colorDanger">删除</span>
                  </Menu.Item>
                </Menu>
              }
              position="bl"
            >
              <Button icon={<IconMore className="fs14" />} />
            </Dropdown>
            <Divider type="vertical" style={{ margin: '0 4px' }} />
            <div className="colorGrey fs12 lh20">
              <div>已选{selectedRowKeys.length}项</div>
            </div>
          </Space>
          <Space>
            <Input.Search
              allowClear
              placeholder="请输入数据搜索"
              style={{ width: 240, marginRight: 4 }}
              value={innerSearch}
              onChange={value => {
                setInnerSearch(value);
              }}
              onSearch={handleSearch}
              onClear={() => {
                handleSearch('');
              }}
            />
          </Space>
        </div>
        <PaginationPage
          total={total}
          sizeOptions={[10, 20, 30, 40, 50, 100]}
          defaultSize={10}
        >
          <Table
            components={{
              body: {
                cell: Cell,
              },
            }}
            showSorterTooltip={false}
            tableLayoutFixed
            style={{ width: '100%' }}
            columns={columns}
            data={data || []}
            rowKey={headers[0]}
            scroll={{
              x: model.type !== 'workspace',
              y: 'calc(100vh - 352px)',
            }}
            rowSelection={{
              type: 'checkbox',
              selectedRowKeys,
              onChange(keys) {
                const hiddenKeys =
                  selectedRowKeys?.filter(
                    key => !data.find(item => item[headers[0]] === key),
                  ) || [];
                setSelectedRowKeys(keys.concat(hiddenKeys));
              },
            }}
            pagination={false}
            onChange={(_, sorter, filter, extra) => {
              if (extra.action !== 'sort') return;
              setSorter(
                `${sorter.field}:${
                  sorter.direction === 'descend' ? 'desc' : 'asc'
                }`,
              );
            }}
            noDataElement={<PageEmpty search={search} />}
          />
        </PaginationPage>
        <CreateSetModal
          model={model}
          visible={collectionVisible}
          rowKeys={selectedRowKeys}
          includeSet={includeSet}
          onConfirm={() => {
            refresh();
          }}
          onClose={() => {
            setCollectionVisible(false);
          }}
        />
        <DeleteSetModal
          visible={deleteVisible}
          model={model}
          isAll={isSelectAll}
          includeSet={includeSet}
          rowKeys={selectedRowKeys}
          onClose={() => {
            setSelectedRowKeys([]);
            setDeleteVisible(false);
          }}
          onConfirm={handleDeleteRow}
          entityList={entityList}
        />
        <SetDetailModal
          data={setValue || []}
          keyword={search}
          name={model?.name?.replace('_set', '')}
          visible={Boolean(setValue)}
          onOK={() => setSetValue(undefined)}
          onCancel={() => setSetValue(undefined)}
          title={!model?.name?.includes('_set_set') ? '实体' : '实体集合'}
        />
      </Spin>
    </div>
  );
};

const filterDown = (item: string, deleteColumn: (header: string) => void) => {
  return ({ confirm }) => {
    return (
      <Popconfirm
        title={
          <>
            <span className="fw500 dpfx">确定删除列数据吗？</span>
            <div className="colorText2 mt12">
              删除后，当前数据列将不再显示，请谨慎操作
            </div>
          </>
        }
        okButtonProps={{ status: 'danger' }}
        onCancel={() => {
          confirm?.();
        }}
        onOk={async () => {
          await deleteColumn(item);
          confirm?.();
        }}
      >
        <div className={styles.deleteColumn}>删除</div>
      </Popconfirm>
    );
  };
};
export const TextRender = ({ value, index }) => {
  return index !== 0 && value?.startsWith('s3://') ? (
    <Copy text={last(value.split('/'))} copyValue={value} />
  ) : (
    <Copy text={value} />
  );
};
export const LinkRender = ({ value, index, onClick, name, search }) => {
  if (index == 1) {
    let parseValue = [];
    try {
      parseValue = JSON.parse(value);
    } catch (error) {
      parseValue = [];
    }
    return (
      <>
        <Link
          disabled={value && value?.length === 0}
          onClick={() => {
            onClick(parseValue);
          }}
          className="fs13"
        >
          {`包含${parseValue?.length || 0}个${
            name.includes('_set_set') ? '实体集合' : '实体'
          }`}
          {search && value && (
            <span className="ml8 colorGrey">
              {`存在${
                parseValue.filter(item => item.includes(search)).length
              }条与关键词匹配的数据`}
            </span>
          )}
        </Link>
      </>
    );
  }
  return (
    <Popover position="tl" content={<Copy text={value} />}>
      <div className="ellipsis fs13">
        {value?.startsWith('s3://') ? last(value.split('/')) : value}
      </div>
    </Popover>
  );
};
export const TextEllipsis = ({ text }) => {
  return (
    <Typography.Text
      style={{ margin: 0 }}
      className="colorBlack3 fs13"
      ellipsis={{
        rows: 1,
        showTooltip: {
          type: 'popover',
          props: {
            content: <div style={{ wordBreak: 'keep-all' }}>{text}</div>,
          },
        },
      }}
    >
      {text || '-'}
    </Typography.Text>
  );
};
export const Cell: React.FC<any> = ({ children }) => {
  const ref = useRef<HTMLTableDataCellElement | null>(null);
  const isVisible = useIntersectionObserver({
    ref,
    options: {
      threshold: 0.01,
      triggerOnce: true,
    },
  });

  return (
    <div ref={ref}>{isVisible ? children : <div style={{ height: 20 }} />}</div>
  );
};
export default React.forwardRef(ModelContent);
