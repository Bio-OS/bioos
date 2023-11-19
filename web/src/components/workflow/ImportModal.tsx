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

import {memo, useMemo, useRef, useState} from 'react';
import { Form as FinalForm } from 'react-final-form';
import { useRouteMatch } from 'react-router-dom';
import { FormApi } from 'final-form';
import {
  Alert,
  Form as ArcoForm,
  Input,
  Message,
  Modal,
  Popover,
  Radio,
} from '@arco-design/web-react';
import { IconQuestionCircle } from '@arco-design/web-react/icon';

import FoldString from 'components/FoldString';
import { maxLength, required } from 'helpers/utils';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';
import { HandlersWorkflowItem } from 'api/index';

import { ToastLimitText } from '..';

import style from './style.less';

const DEFAULT_LANGUAGE = 'WDL'

interface Props {
  visible: boolean;
  workflowInfo?: HandlersWorkflowItem;
  refetch?: () => void;
  onClose?: () => void;
}

function ImportModal({ visible, workflowInfo, onClose, refetch }: Props) {
  const refForm = useRef<FormApi>();
  const match = useRouteMatch<{ workspaceId: string }>();
  const isEdit = workflowInfo?.latestVersion.status === 'Success';
  const isReimport = workflowInfo?.latestVersion.status === 'Failed';
  const [language, _] = useState(DEFAULT_LANGUAGE)

  const initialValues = useMemo(
    () =>
      workflowInfo
        ? {
            name: workflowInfo.name,
            language: workflowInfo.latestVersion.language,
            url: workflowInfo.latestVersion.metadata.gitURL,
            tag: workflowInfo.latestVersion.metadata.gitTag,
            mainWorkflowPath: workflowInfo.latestVersion.mainWorkflowPath,
            description: workflowInfo.description,
          }
        : {
            language,
          },
    [workflowInfo],
  );

  async function handleSubmit(values) {
    let res;
    if (isEdit || isReimport) {
      res = await Api.workspaceIdWorkflowPartialUpdate(
        match.params.workspaceId,
        workflowInfo.id,
        {
          ...values,
          token: values.token || '',
          description: values.description || '',
        },
      );
    } else {
      res = await Api.workspaceIdWorkflowCreate(match.params.workspaceId, {
        ...values,
        source: 'git',
      });
    }

    if (res.ok) {
      Message.success({
        content: (
          <ToastLimitText
            name={values.name}
            prefix="工作流"
            suffix="开始导入至当前Workspace"
          />
        ),
      });
      refetch();
      onClose();
    }
  }

  async function validateWorkflowName(value: string) {
    if ((isEdit || isReimport) && value === workflowInfo.name) {
      return;
    }
    const res = await Api.workspaceIdWorkflowList(match.params.workspaceId, {
      searchWord: value,
      exact: true,
    });
    const exist = res?.data?.items?.some(item => item.name === value);

    return exist ? '名称不可与已有的重复' : undefined;
  }

  function getTitle() {
    if (isReimport) return '重新导入';
    if (isEdit) return '更新工作流';
    return '导入工作流';
  }

  return (
    <Modal
      className={style.importModal}
      title={getTitle()}
      visible={visible}
      focusLock={false}
      style={{ width: 480 }}
      unmountOnExit={true}
      maskClosable={false}
      escToExit={false}
      onConfirm={() => refForm.current.submit()}
      onCancel={onClose}
    >
      {isReimport && workflowInfo.latestVersion.message && (
        <Alert
          type="warning"
          title="请检查以下错误信息，修改完成后重新导入"
          style={{
            margin: '0 -24px 24px',
            width: 'calc(100% + 48px)',
            padding: '8px 24px',
          }}
          content={
            <FoldString
              lineClamp={1}
              className="colorBlack4 fs12"
              style={{ whiteSpace: 'pre-wrap', lineHeight: '22px' }}
            >
              {workflowInfo.latestVersion.message}
            </FoldString>
          }
        />
      )}

      <FinalForm
        subscription={{ initialValues: true }}
        initialValues={initialValues}
        onSubmit={handleSubmit}
      >
        {({ form }) => {
          refForm.current = form;

          return (
            <ArcoForm layout="vertical">
              <FieldItem
                name="name"
                label="工作流名称"
                required={true}
                rules={[
                  ...getNameRules(1, 60),
                  { asyncValidate: validateWorkflowName },
                ]}
              >
                <Input placeholder="请输入" allowClear={true} />
              </FieldItem>
              <FieldItem
                  name="language"
                  label="规范"
                  validate={[
                    (val: string | undefined) => {
                      if (val === undefined) {
                        return '请选择 workflow 规范';
                      }
                      return;
                    }
                  ]}
              >
                <Radio.Group type="button" defaultValue={language}>
                  <Radio value="WDL">WDL</Radio>
                  <Radio value="Nextflow">Nextflow</Radio>
                </Radio.Group>
              </FieldItem>
              <FieldItem
                name="url"
                label="Git 地址"
                required={true}
                validate={[
                  required,
                  (val: string) => {
                    if (val.startsWith('http://') || val.startsWith('https://'))
                      return;
                    return '需要以 http:// 或者 https:// 开头';
                  },
                ]}
              >
                <Input
                  allowClear={true}
                  placeholder="请输入 http 或 https 协议的地址"
                />
              </FieldItem>
              <FieldItem
                name="tag"
                label="Branch/Tag"
                required={true}
                validate={[required]}
              >
                <Input placeholder="请输入" allowClear={true} />
              </FieldItem>
              <FieldItem
                name="token"
                label={
                  <>
                    Token
                    <Popover content=" 私有仓库需指定token" position="right">
                      <IconQuestionCircle
                        className="ml4"
                        style={{ color: '#86909c' }}
                      />
                    </Popover>
                  </>
                }
              >
                <Input placeholder="请输入" allowClear={true} />
              </FieldItem>
              <FieldItem
                name="mainWorkflowPath"
                label="主工作流路径"
                required={true}
                validate={[required]}
              >
                <Input.TextArea placeholder="请输入仓库内的指定文件的具体路径，不包含仓库地址信息" />
              </FieldItem>
              <FieldItem
                name="description"
                label="简短描述"
                validate={[maxLength(1000)]}
              >
                <Input.TextArea
                  placeholder="请输入"
                  style={{ height: 74 }}
                  maxLength={{ length: 1000, errorOnly: true }}
                  showWordLimit={true}
                />
              </FieldItem>
            </ArcoForm>
          );
        }}
      </FinalForm>
    </Modal>
  );
}

export default memo(ImportModal, (pre, next) => {
  return !(pre.visible !== next.visible || pre.workflowInfo !== next.workflowInfo);
});
