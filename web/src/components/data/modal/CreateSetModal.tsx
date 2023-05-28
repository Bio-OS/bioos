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
import { Form as FinalForm } from 'react-final-form';
import { useParams } from 'react-router-dom';
import dayjs from 'dayjs';
import { FormApi } from 'final-form';
import {
  Form as ArcoForm,
  Input,
  Message,
  Modal,
} from '@arco-design/web-react';

import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';
import { HandlersDataModel } from 'api/index';
interface Props {
  model?: HandlersDataModel;
  visible?: boolean;
  rowKeys?: (string | number)[];
  onClose?: () => void;
  onConfirm?: () => void;
  includeSet?: HandlersDataModel;
}
interface CreateSetFormValue {
  name: string;
  setName: string;
}
const CreateSetModal: React.FC<Props> = props => {
  const { visible, model, rowKeys, onClose, onConfirm, includeSet } = props;
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [setRows, setSetRows] = useState([]);
  const formRef = useRef<undefined | FormApi<CreateSetFormValue>>();
  const initialName = dayjs().format('YYYY-MM-DD-HH-mm-ss');
  const handleSubmit = (values: CreateSetFormValue) => {
    if (values.name.length > 50) {
      Message.error('实体集合表名称超过50字符限制，创建失败');
      return;
    }
    const headers = [
      `${values.name}_id`,
      values.name.endsWith('_set') ? values.name.replace('_set', '') : '',
    ];
    const rows = [[`${values.setName}`, JSON.stringify(rowKeys)]];
    Api.dataModelPartialUpdate(workspaceId, {
      headers,
      rows,
      name: values?.name,
      workspaceID: workspaceId,
    })
      .then(res => {
        if (res.ok) {
          Message.success({
            content: (
              <>
                <span>生成实体集合</span>
                <span
                  className="ellipsis ml4 mr4"
                  style={{
                    maxWidth: 120,
                    display: 'inline-block',
                  }}
                >
                  {values.name}
                </span>
                <span>成功</span>
              </>
            ),
          });
          onConfirm?.();
          onClose?.();
        } else {
          Message.error(res?.statusText || '生成实体集合失败');
        }
      })
      .catch(e => {
        Message.error(e?.error?.message || e?.statusText || '生成实体集合失败');
      });
  };
  useEffect(() => {
    if (!visible || !includeSet) return;
    Api.dataModelRowsDetail(workspaceId, includeSet.id).then(res => {
      if (res.ok) {
        setSetRows(res.data.rows);
      }
    });
  }, [includeSet, visible]);
  return (
    <Modal
      style={{ width: 500 }}
      onOk={() => {
        formRef.current?.submit();
      }}
      onCancel={() => {
        onClose?.();
      }}
      visible={visible}
      title="生成实体集合"
    >
      <FinalForm<CreateSetFormValue>
        onSubmit={handleSubmit}
        initialValues={{
          name: `${model?.name}_set`,
          setName: `${model?.name}_set-${initialName}`,
        }}
      >
        {({ form }) => {
          formRef.current = form;
          return (
            <ArcoForm
              labelCol={{ span: 8 }}
              wrapperCol={{ span: 16 }}
              labelAlign="left"
            >
              <FieldItem
                name="name"
                label="实体集合表名称"
                required={true}
                disabled={true}
              >
                <div>{`${model?.name}_set`}</div>
              </FieldItem>
              <FieldItem
                name="setName"
                label="实体集合ID"
                required={true}
                validateFields={[]}
                rules={[
                  {
                    validate: (value: string) =>
                      !/^.{1,100}$/.test(value ?? ''),
                    message: '请输入 1 ~ 100 个字符。',
                  },
                  {
                    asyncValidate: (value: string) => {
                      return setRows?.map(item => item[0]).includes(value)
                        ? '名称不可与已有的重复'
                        : undefined;
                    },
                  },
                ]}
              >
                <Input placeholder="请输入" />
              </FieldItem>
            </ArcoForm>
          );
        }}
      </FinalForm>
    </Modal>
  );
};
export default CreateSetModal;
