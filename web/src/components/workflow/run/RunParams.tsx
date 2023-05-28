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

import {
  forwardRef,
  memo,
  useEffect,
  useImperativeHandle,
  useMemo,
  useState,
} from 'react';
import { FormSpy, useForm } from 'react-final-form';
import classNames from 'classnames';
import {
  Button,
  Input,
  Message,
  Popover,
  Space,
  Switch,
  Table,
  TableColumnProps,
} from '@arco-design/web-react';
import { IconInfoCircle, IconPlus } from '@arco-design/web-react/icon';

import PageEmpty from 'components/Empty';
import ResizebleTable from 'components/ResizebleTable';
import { useForceUpdate } from 'helpers/hooks';
import { HandlersWorkflowParam } from 'api/index';

import { Line } from '../..';

import { ParamsTitle, renderValidateForm } from './components';
import DownloadJSON from './DownloadJSON';
import ParamFieldItem from './ParamFieldItem';
import { DeleteModal, ParamModal, ParamModalProps } from './ParamModal';
import UploadJSON from './UploadJSON';
import {
  isStringType,
  renderKey,
  replaceDotToHyphen,
  stringKey,
} from './utils';

import styles from './style.less';

interface Props {
  type: 'input' | 'output';
  isPath: boolean;
  isHidden: boolean;
  rowIds?: string[];
  data: HandlersWorkflowParam[];
  selectOptions?: string[];
  invalidHeaderArr?: string[];
}

const BASE_COLUMNS: TableColumnProps[] = [
  {
    dataIndex: 'name',
    width: 200,
    title: '变量',
  },
  {
    dataIndex: 'type',
    width: 122,
    title: (
      <>
        类型
        <Popover content="file、Array[File] 类型下属性值仅支持“s3://”格式">
          <IconInfoCircle className="ml4" style={{ color: '#86909c' }} />
        </Popover>
      </>
    ),
  },
];

export default memo(
  forwardRef(function RunParams(
    {
      type,
      isPath,
      isHidden,
      data,
      rowIds,
      selectOptions = [],
      invalidHeaderArr = [],
    }: Props,
    ref,
  ) {
    const [rowKeys, setRowKeys] = useState<string[]>(rowIds);
    const [checkedRowKeys, setCheckedRowKeys] = useState<string[]>([]);
    const [onlyShowRequired, setOnlyShowRequired] = useState<boolean>(false);
    const [search, setSearch] = useState<string>();
    const [deleteData, setDeleteData] = useState<string[] | undefined>(
      undefined,
    );
    const [modalParams, setModalParams] = useState<
      Pick<ParamModalProps, 'originName' | 'type'>
    >({ type: undefined, originName: '' });

    const forceUpdate = useForceUpdate();
    const form = useForm();
    const isInput = type === 'input';
    const multiInput = isInput && isPath;

    const allChecked = rowKeys.length === checkedRowKeys?.length;
    const noChecked = !checkedRowKeys?.length;
    const onlyOneRow = rowKeys.length === 1;
    const deleteButtonDisabled = noChecked || onlyOneRow || allChecked;

    useEffect(() => {
      setRowKeys(rowIds);
    }, [rowIds]);

    const columns = useMemo(() => {
      return BASE_COLUMNS.concat(rowKeys?.map(key => renderField(key)) || []);
    }, [selectOptions, rowKeys, checkedRowKeys]);

    const numRequiredFields = useMemo(() => {
      if (!data || !form.getState().errors) return 0;

      const requiredData = data.filter(item => !item.optional);

      return requiredData.reduce((acc, item) => {
        rowKeys.forEach((key: string) => {
          const errorFieldArr = Object.keys(
            form.getState().errors?.[type]?.[key] || {},
          );

          if (errorFieldArr.includes(replaceDotToHyphen(item.name))) {
            acc++;
          }
        });

        return acc;
      }, 0);
    }, [data, columns]);

    useImperativeHandle(ref, () => {
      if (!(isPath && isInput)) {
        return null;
      }
      return {
        pathRowKeys: rowKeys,
        checkedPathRowKeys: rowKeys.filter(_ => checkedRowKeys.includes(_)),
      };
    });

    function handleCheck(key: string) {
      const index = checkedRowKeys.indexOf(key);
      if (index > -1) {
        checkedRowKeys.splice(index, 1);
        setCheckedRowKeys([...checkedRowKeys]);
      } else {
        setCheckedRowKeys(checkedRowKeys.concat(key));
      }
    }

    const autoTouched = (key: string, value: string) => {
      // 这里仅对有值field进行校验
      if (!value) return;

      // 设置field touched 状态
      form.focus(key);
      // 触发field校验，因为change为相同值不会触发校验，所以这里先设为undefined，再修改为原来值
      form.change(key, undefined);
      form.change(key, value);
      form.blur(key);
    };

    // 触发保存
    const makeFormDirty = () => {
      autoTouched('dirty', new Date().toString());
    };

    function handleAdd(name: string) {
      setRowKeys([...rowKeys, name]);
      makeFormDirty();
      Message.success('新增实体属性成功');
      setTimeout(() => {
        // 滚动到最右端
        const container: HTMLDivElement = document
          .querySelector('.inputPathParams')
          ?.querySelector('.arco-table-body');

        const scrollLeft = container?.offsetWidth;
        container?.scrollTo({ left: scrollLeft });
      }, 0);
    }

    function handleCopy(newName: string, preName: string) {
      const preIndex = rowKeys.findIndex(item => item === preName);
      rowKeys.splice(preIndex + 1, 0, newName);
      setRowKeys([...rowKeys]);
      const preValues = form.getState().values.input?.[stringKey(preName)];
      form.change(`input.${stringKey(newName)}`, {});

      makeFormDirty();
      Message.success('复制实体属性成功');

      if (preValues) {
        setTimeout(() => {
          Object.keys(preValues).forEach(item =>
            autoTouched(renderKey(isInput, newName, item), preValues[item]),
          );
        }, 0);
      }
    }

    function handleRename(newName: string, preName: string) {
      if (newName === preName) return;
      const i = rowKeys.indexOf(preName);
      const j = checkedRowKeys.indexOf(preName);
      if (i > -1) {
        rowKeys[i] = newName;
        setRowKeys([...rowKeys]);
      }
      if (j > -1) {
        checkedRowKeys[j] = newName;
        setCheckedRowKeys([...checkedRowKeys]);
      }

      const preValues = form.getState().values.input?.[stringKey(preName)];
      form.change(`input.${stringKey(newName)}`, {});

      makeFormDirty();
      Message.success('重命名实体属性成功');

      if (preValues) {
        setTimeout(() => {
          Object.keys(preValues).forEach(item =>
            autoTouched(renderKey(isInput, newName, item), preValues[item]),
          );
          form.change(`input.${stringKey(preName)}`, undefined);
        }, 0);
      }
    }

    function handleDelete() {
      const newColumns = rowKeys.filter(key => !deleteData?.includes(key));
      const newCheckedColumns = checkedRowKeys.filter(
        key => !deleteData?.includes(key),
      );

      setRowKeys(newColumns);
      setCheckedRowKeys(newCheckedColumns);

      form.batch(() => {
        deleteData?.forEach(key =>
          form.change(`input.${stringKey(key)}`, undefined),
        );
      });

      makeFormDirty();
      Message.success('删除实体属性成功');
    }

    function handleChangeColumn({ name }: { name: string }) {
      switch (modalParams.type) {
        case 'add':
          handleAdd(name);
          break;
        case 'rename':
          handleRename(name, modalParams.originName);
          break;
        case 'copy':
          handleCopy(name, modalParams.originName);
          break;
        default:
          break;
      }
    }

    function renderField(key: string) {
      return {
        dataIndex: key,
        width: 300,
        title: (
          <ParamsTitle
            name={key}
            deleteDisabled={multiInput && rowKeys?.length === 1}
            checked={checkedRowKeys.includes(key)}
            onCheck={() => handleCheck(key)}
            onCopy={() => setModalParams({ type: 'copy', originName: key })}
            onRename={() => setModalParams({ type: 'rename', originName: key })}
            onDelete={() => setDeleteData([key])}
          />
        ),
        bodyCellStyle: {
          padding: '8px 16px',
          backgroundColor: checkedRowKeys.includes(key) ? '#fafbfc' : '#fff',
        },
        headerCellStyle: {
          backgroundColor: checkedRowKeys.includes(key) ? '#f1f3f5' : '#f6f8fa',
        },
        render: (_val: unknown, item: HandlersWorkflowParam) => (
          <ParamFieldItem
            name={renderKey(isInput, key, item.name)}
            item={item}
            disabled={!isInput && isPath}
            selectOptions={selectOptions}
          />
        ),
      };
    }

    const tableData = useMemo(() => {
      let dataFiltered = onlyShowRequired
        ? data?.filter(item => !item.optional)
        : data?.map(item => {
            item.default =
              item.default &&
              isStringType(item.type) &&
              !(item.default.startsWith('"') && item.default.endsWith('"'))
                ? `"${item.default}"`
                : item.default;
            return item;
          });

      if (search) {
        dataFiltered = dataFiltered?.filter(item =>
          item.name.toLocaleLowerCase().includes(search.toLocaleLowerCase()),
        );
      }

      return dataFiltered;
    }, [onlyShowRequired, data, search]);

    const getUsedData = () => {
      let downLoadData: string[] = [];

      if (rowKeys.length === 1) {
        downLoadData = [rowKeys[0]];
      } else if (!checkedRowKeys.length) {
        Message.error('请选择实体属性');
        return;
      } else {
        downLoadData = checkedRowKeys;
      }

      return downLoadData;
    };

    const handleResizeTable = (index: number, width: number) => {
      columns[index].width = width;
      forceUpdate();
    };

    const TableComponent = isInput && isPath ? ResizebleTable : Table;
    // const className =
    return (
      <div
        className={classNames([
          'inputPathParams',
          styles.runWorkflowWrap,
          { hidden: isHidden },
        ])}
      >
        <div className="flexBetween mb16 mt4 mr4">
          <div>
            {isInput && isPath && (
              <Space className="mr12">
                <Button
                  type="primary"
                  onClick={() => setModalParams({ type: 'add' })}
                >
                  <IconPlus />
                  新增实体属性
                </Button>
                <Button
                  onClick={() => {
                    allChecked
                      ? setCheckedRowKeys([])
                      : setCheckedRowKeys([...rowKeys]);
                  }}
                >
                  {allChecked ? '取消全选实体属性' : '勾选全部实体属性'}
                </Button>
                <Popover
                  disabled={!deleteButtonDisabled || (!onlyOneRow && noChecked)}
                  trigger="hover"
                  content="请至少保留一个属性值"
                >
                  <Button
                    onClick={() => setDeleteData(checkedRowKeys)}
                    disabled={deleteButtonDisabled}
                  >
                    删除
                  </Button>
                </Popover>
              </Space>
            )}
            <Switch
              className="mr8"
              checked={onlyShowRequired}
              onChange={setOnlyShowRequired}
            />
            <span>只显示必选参数</span>

            <Line />

            <span className="colorGrey">
              输入存在
              <FormSpy subscription={{ errors: true }}>
                {() => (
                  <span className="mr4 ml4 colorPrimary">
                    {numRequiredFields}
                  </span>
                )}
              </FormSpy>
              项必选参数未填写成功
            </span>
          </div>

          <div className="flexAlignCenter">
            <>
              <DownloadJSON
                disabled={isPath && !isInput}
                checkData={getUsedData}
                isInput={isInput}
                data={data}
              />
              <Line />
              <UploadJSON
                disabled={isPath && !isInput}
                checkData={getUsedData}
                isInput={isInput}
                data={data}
                clickDisabled={!checkedRowKeys.length && rowKeys.length > 1}
                onUploadComplete={() => setSearch('')}
              />
            </>

            <Input.Search
              className="w200"
              placeholder="请输入变量搜索"
              onSearch={setSearch}
              allowClear={true}
              onClear={() => setSearch('')}
            />
          </div>
        </div>
        {renderValidateForm(
          rowKeys,
          data,
          isHidden,
          isInput,
          isPath,
          invalidHeaderArr,
        )}
        <TableComponent
          border={true}
          borderCell={true}
          noDataElement={<PageEmpty search={search} />}
          scroll={{ y: 600 }}
          rowKey="name"
          columns={columns}
          data={tableData}
          pagination={false}
          onResize={handleResizeTable}
        />
        {isInput && isPath && (
          <>
            <ParamModal
              type={modalParams.type}
              originName={modalParams.originName}
              names={rowKeys}
              onChange={handleChangeColumn}
              onHide={() => setModalParams({})}
            />
            <DeleteModal
              data={deleteData}
              onHide={() => setDeleteData(undefined)}
              onDelete={handleDelete}
            />
          </>
        )}
      </div>
    );
  }),
);
