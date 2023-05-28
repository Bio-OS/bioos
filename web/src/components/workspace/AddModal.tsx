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

import { memo, useEffect, useRef } from 'react';
import { Form as FinalForm, FormSpy } from 'react-final-form';
import { FormApi } from 'final-form';
import { omit, pick } from 'lodash-es';
import { Form as ArcoForm, Input, Modal, Select } from '@arco-design/web-react';

import { useGetEnvQuery } from 'helpers/hooks';
import { maxLength, required } from 'helpers/utils';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import API from 'api/client';

const Option = Select.Option;

interface Props {
  type: 'add' | 'edit';
  visible: boolean;
  workspaceId?: string;
  refetch: () => void;
  onClose: () => void;
}

const options = [
  {
    name: 'NAS 文件存储',
    value: 'nfs',
  },
];

const FORM_WIDTH = 432;

function AddModal({
  type = 'add',
  workspaceId,
  visible,
  onClose,
  refetch,
}: Props) {
  const refForm = useRef<FormApi>();
  const { storage: storageConfig } = useGetEnvQuery();
  const isAdd = type === 'add';

  function action(params) {
    if (isAdd) {
      API.workspaceCreate(params).then(() => {
        refetch();
        onClose();
      });
      return;
    }

    API.workspacePartialUpdate(workspaceId, {
      ...params,
      id: workspaceId,
    }).then(() => {
      refetch();
      onClose();
    });
  }
  function handleSubmit(values) {
    const parmas = omit(values, ['storage', 'mountPath', 'mountDir']);
    if (isAdd) {
      parmas.storage = {};
      parmas.storage[values.storage] = {
        mountPath: values.mountDir + values.mountPath,
      };
    }
    action(parmas);
  }

  function getWorkspaceDetail(id: string) {
    API.workspaceDetail(id)
      .then(({ data }) => {
        const values = pick(data, ['name', 'description']);
        const storage = Object.keys(data.storage)?.[0];
        let mountPath: string = data.storage[storage].mountPath;
        let mountDir = storageConfig?.fsPath?.[0];
        storageConfig?.fsPath?.forEach(path => {
          if (mountPath.startsWith(path)) {
            mountDir = path;
            mountPath = mountPath.slice(path.length);
          }
        });

        refForm.current.reset({
          ...values,
          storage,
          mountPath,
          mountDir,
        });
      })
      .catch(err => {
        console.error(err);
      });
  }

  async function validateWorkspaceName(value: string) {
    if (!isAdd) {
      return;
    }
    const res = await API.workspaceList({
      searchWord: value,
      exact: true,
    });
    const exist = res?.data?.items?.some(item => item.name === value);

    return exist ? '名称不可与已有的重复' : undefined;
  }

  useEffect(() => {
    if (isAdd) return;
    if (workspaceId) {
      getWorkspaceDetail(workspaceId);
    }
  }, [type, workspaceId]);

  return (
    <Modal
      title={`${isAdd ? '新建' : '编辑'} Workspace`}
      visible={visible}
      focusLock={false}
      style={{ width: 480 }}
      maskClosable={false}
      escToExit={false}
      onConfirm={() => refForm.current.submit()}
      onCancel={onClose}
      unmountOnExit={true}
    >
      <FinalForm
        initialValues={{
          storage: 'nfs',
          mountDir: storageConfig?.fsPath?.[0],
        }}
        subscription={{ initialValues: true }}
        onSubmit={handleSubmit}
      >
        {({ form }) => {
          refForm.current = form;
          return (
            <ArcoForm layout="vertical">
              <FieldItem
                name="name"
                label="名称"
                required={true}
                rules={[
                  ...getNameRules(1, 60),
                  { asyncValidate: validateWorkspaceName },
                ]}
              >
                <Input.TextArea
                  placeholder="请输入"
                  style={{
                    width: FORM_WIDTH,
                    overflowY: 'hidden',
                  }}
                  allowClear={true}
                  autoSize={{ minRows: 1, maxRows: 6 }}
                />
              </FieldItem>
              <FieldItem
                name="storage"
                label="存储类型"
                required={true}
                disabled={!isAdd}
              >
                <Select style={{ width: FORM_WIDTH }}>
                  {options.map(item => (
                    <Option key={item.value} value={item.value}>
                      {item.name}
                    </Option>
                  ))}
                </Select>
              </FieldItem>
              <FormSpy>
                {({ values }) => (
                  <FieldItem
                    name="mountPath"
                    label="挂载目录"
                    required={true}
                    disabled={!isAdd}
                    validate={[
                      required,
                      value => {
                        if (!value.startsWith('/')) return '挂载目录需以/开头';
                        return undefined;
                      },
                    ]}
                  >
                    <Input
                      style={{ width: FORM_WIDTH }}
                      addBefore={
                        <Select
                          value={values.mountDir}
                          style={{ width: 160 }}
                          disabled={!isAdd}
                        >
                          {storageConfig?.fsPath?.map(item => (
                            <Select.Option key={item} value={item}>
                              {item}
                            </Select.Option>
                          ))}
                        </Select>
                      }
                      allowClear={true}
                      placeholder="请输入目录"
                    />
                  </FieldItem>
                )}
              </FormSpy>
              <FieldItem
                name="description"
                label="简短描述"
                required
                validate={[required, maxLength(1000)]}
              >
                <Input.TextArea
                  placeholder="请输入"
                  style={{ height: 74, width: FORM_WIDTH }}
                  maxLength={{ length: 1000, errorOnly: true }}
                  showWordLimit={true}
                  allowClear={true}
                />
              </FieldItem>
            </ArcoForm>
          );
        }}
      </FinalForm>
    </Modal>
  );
}

export default memo(AddModal);
