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

import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import classNames from 'classnames';
import { Alert, Message, Modal, Table } from '@arco-design/web-react';

import CommonLimitText from 'components/CommonLimitText';
import { genHighlightText } from 'components/getHighLight';
import Api from 'api/client';
import { HandlersDataModel } from 'api/index';

import styles from './DeleteSetModal.less';

interface Props {
  model?: HandlersDataModel;
  visible: boolean;
  isAll: boolean;
  includeSet?: HandlersDataModel;
  rowKeys: (string | number)[];
  onClose?: () => void;
  onConfirm?: () => void;
  entityList?: HandlersDataModel[];
}
const DeleteSetModal: React.FC<Props> = props => {
  const {
    visible,
    model,
    includeSet,
    isAll,
    rowKeys,
    onClose,
    onConfirm,
    entityList,
  } = props;
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    page: 1,
    size: 10,
  });
  const [total, setTotal] = useState(0);
  const [includeSetData, setIncludeSetData] = useState({
    columns: [],
    data: [],
  });
  const isSet = model?.type === 'entity_set';
  const collectionTips = isSet ? '集合' : '';
  useEffect(() => {
    if (includeSet && !isAll && visible) {
      setLoading(true);
      Api.dataModelRowsDetail(workspaceId, includeSet.id, {
        size: pagination.size,
        page: pagination.page,
        inSetIDs: rowKeys as string[],
      })
        .then(res => {
          if (res.ok && res?.data) {
            const { rows, total } = res.data;
            const deleteIDs = rows?.map(item => item[0]);
            const data = deleteIDs?.map((item, index) => {
              return { key: index, ID: item, name: `${model.name}_set` };
            });
            const columns = [
              {
                title: '实体集合ID',
                dataIndex: 'ID',
              },
              {
                title: '实体集合表名称',
                dataIndex: 'name',
                width: '42%',
              },
            ];
            setIncludeSetData({
              columns,
              data,
            });
            setTotal(total);
            setLoading(false);
          }
        })
        .catch(e => {
          Message.error(
            e?.error?.message || e?.statusText || '获取实体数据失败',
          );
        });
    }
  }, [includeSet, model, isAll, visible, pagination]);
  const handleOk = () => {
    setLoading(true);
    Api.dataModelDelete(
      workspaceId,
      model.id,
      isAll
        ? undefined
        : {
            rowIDs: rowKeys as string[],
          },
      { format: 'json' },
    )
      .then(res => {
        if (res?.ok) {
          Message.success('删除数据成功');
          onConfirm();
          setLoading(false);
        }
      })
      .catch(e => {
        setLoading(false);
        Message.error(e?.error?.message || e?.statusText || '删除数据失败');
      });
  };
  return (
    <Modal
      visible={visible}
      style={{ width: includeSet && !isAll && !!total ? 580 : 400 }}
      title={`确定删除${
        model?.type !== 'entity_set' ? '所选实体ID' : '所选实体集合ID'
      } 吗？`}
      okText="删除"
      okButtonProps={{
        status: 'danger',
        disabled: loading,
      }}
      onOk={handleOk}
      onCancel={() => {
        setLoading(false);
        onClose?.();
      }}
      maskClosable={false}
      className={classNames({
        [styles.deleteSet]: includeSet && !isAll && !!total,
      })}
    >
      <div className="colorBlack3 fs13">
        {isAll && (
          <>
            <span>{`已选择删除实体${collectionTips}表中所有的实体${collectionTips}
            ID，删除后实体${collectionTips}表${
              includeSet ? '、引用此数据的实体集合表' : ''
            }也将被同时删除，请谨慎操作。`}</span>
            {!isSet && (
              <CommonLimitText
                name="实体表"
                value={model?.name}
                style={{ marginTop: 12 }}
              />
            )}
            {(includeSet || isSet) && (
              <CommonLimitText
                name="实体集合表"
                value={entityList
                  .filter(
                    item =>
                      item.name.replace(/_set/g, '') ===
                        `${model?.name.replace(/_set/g, '')}` &&
                      item.name.includes(isSet ? model?.name : '_set'),
                  )
                  .map(_ => _.name)
                  .join('、')}
                style={{ marginTop: 12 }}
              />
            )}
          </>
        )}
        {((!includeSet && !isAll) || (includeSet && !total)) && (
          <span>
            {genHighlightText(
              `当前已选择删除${rowKeys?.length}个实体${collectionTips}ID，删除后无法通过平台恢复，请谨慎操作。`,
              String(rowKeys?.length || ''),
            )}
          </span>
        )}
      </div>
      {includeSet && !isAll && !!total && (
        <>
          <Alert
            className="fs12"
            type="warning"
            content={
              <span>
                {genHighlightText(
                  `已选择删除${rowKeys?.length}个实体${collectionTips}ID，同时也会同步从以下引用此数据的实体集合中删除，请谨慎操作。`,
                  String(rowKeys?.length || ''),
                )}
              </span>
            }
          />
          <div className={styles.deleteSetTable}>
            <Table
              data={includeSetData.data || []}
              columns={includeSetData.columns || []}
              scroll={{
                y: 318,
              }}
              loading={loading}
              pagination={{
                total,
                current: pagination.page,
                pageSize: pagination.size,
                hideOnSinglePage: total <= pagination.size,
                showTotal: true,
                onChange(page, size) {
                  setPagination({
                    page,
                    size,
                  });
                },
              }}
            />
          </div>
        </>
      )}
    </Modal>
  );
};

export default DeleteSetModal;
