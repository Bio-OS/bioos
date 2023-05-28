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

import { useCallback, useRef } from 'react';
import { Form as FinalForm } from 'react-final-form';
import { FormApi } from 'final-form';
import {
  Alert,
  Form as ArcoForm,
  Input,
  Modal,
  Table,
  Typography,
} from '@arco-design/web-react';

import { currentTime } from 'helpers/utils';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';

interface FormValues {
  name: string;
}

export interface ParamModalProps {
  originName?: string;
  names: string[];
  type?: 'add' | 'rename' | 'copy';
  onChange: (values: FormValues) => void;
  onHide: () => void;
}

export function ParamModal({
  originName,
  names,
  type,
  onChange,
  onHide,
}: ParamModalProps) {
  const formRef = useRef<undefined | FormApi<FormValues>>();

  function hide() {
    onHide();
  }

  function handleSubmit(values: FormValues) {
    hide();
    onChange(values);
  }

  const getInitialValue = useCallback(() => {
    let title = '';
    let initialValue = '';
    switch (type) {
      case 'add':
        title = '新增实体属性';
        initialValue = currentTime();
        break;
      case 'rename':
        title = '重命名实体属性';
        initialValue = originName;
        break;
      case 'copy':
        title = '复制实体属性';
        initialValue = `${originName}-copy`;
        break;
      default:
        break;
    }

    return {
      title,
      initialValue,
    };
  }, [type]);

  const { title, initialValue } = getInitialValue();
  const vaildateNames =
    type === 'rename' ? names.filter(_ => _ !== originName) : names;

  return (
    <Modal
      visible={Boolean(type)}
      title={<div style={{ textAlign: 'left' }}>{title}</div>}
      onOk={() => formRef.current?.submit()}
      unmountOnExit={true}
      onCancel={hide}
      maskClosable={false}
      escToExit={false}
      style={{ width: 440 }}
    >
      {type === 'copy' && (
        <Alert
          content="将复制当前实体属性内容作为新的实体属性。"
          className="flexAlignCenter"
          style={{ width: 440, height: 40, margin: '0 0 24px -24px' }}
        />
      )}

      <FinalForm onSubmit={handleSubmit}>
        {({ form }) => {
          formRef.current = form;
          return (
            <ArcoForm
              labelCol={{ span: 4 }}
              wrapperCol={{ span: 20 }}
              labelAlign="left"
            >
              <FieldItem
                name="name"
                label="名称"
                required={true}
                initialValue={initialValue}
                rules={
                  type
                    ? [
                        ...getNameRules(1, 200),
                        {
                          asyncValidate: (value: string) => {
                            return vaildateNames.includes(value)
                              ? '属性名称不可与已有的重复'
                              : undefined;
                          },
                        },
                      ]
                    : []
                }
              >
                <Input placeholder="请输入" />
              </FieldItem>
            </ArcoForm>
          );
        }}
      </FinalForm>
    </Modal>
  );
}

const COLUMUNS = [
  {
    dataIndex: 'name',
    title: '实体属性',
    width: 380,
    render(value: string) {
      return (
        <Typography.Text
          ellipsis={{
            cssEllipsis: true,
            showTooltip: {
              type: 'popover',
            },
          }}
        >
          {value}
        </Typography.Text>
      );
    },
  },
];

export function DeleteModal({
  data,
  onHide,
  onDelete,
}: {
  data: string[] | undefined;
  onHide: () => void;
  onDelete: () => void;
}) {
  function handleDelete() {
    onDelete();
    hide();
  }

  function hide() {
    onHide();
  }

  return (
    <Modal
      visible={Boolean(data)}
      title={<div style={{ textAlign: 'left' }}>确定删除实体属性吗？</div>}
      okText="删除"
      okButtonProps={{ status: 'danger' }}
      onOk={handleDelete}
      onCancel={hide}
      style={{ width: 440 }}
      unmountOnExit={true}
      maskClosable={true}
      escToExit={true}
    >
      <Alert
        type="warning"
        content="删除后，以下实体属性将不再显示，请谨慎操作。"
        style={{ width: 440, height: 40, margin: '-24px -20px 24px -20px' }}
      />
      <Table
        className="tableOverflowXHidden"
        data={data?.map(item => ({ name: item }))}
        rowKey="name"
        columns={COLUMUNS}
        pagination={false}
        border={true}
        scroll={{ y: 368, x: false }}
      />
    </Modal>
  );
}
