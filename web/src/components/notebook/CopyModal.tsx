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

import { useEffect, useRef, useState } from 'react';
import { Form as FinalForm, FormSpy } from 'react-final-form';
import { useHistory, useRouteMatch } from 'react-router-dom';
import { FormApi } from 'final-form';
import {
  Form as ArcoForm,
  Input,
  Message,
  Modal,
  Select,
} from '@arco-design/web-react';

import Icon from 'components/Icon';
import TimerModal from 'components/TimerModal';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';
interface Props {
  name: string;
  refetch: () => void;
  onClose: () => void;
}

interface FormCopyNotebook {
  workspace: string;
  name: string;
}

export default function Copy({ name, refetch, onClose }: Props) {
  const refForm = useRef(null);
  const history = useHistory();
  const [timerModalInfo, setTimerModalInfo] = useState<{
    notebookName: string;
  }>();
  const match = useRouteMatch<{ workspaceId: string }>();
  const workspaceId = match.params.workspaceId;
  const [workspaceList, setWorkspaceList] = useState([]);
  const [purposeUrl, setPurposeUrl] = useState<string>(workspaceId);
  const handleSubmit = async function (values) {
    const { data } = await Api.workspaceIdNotebookDetail(workspaceId, name, {
      format: 'json',
    });
    // console.log(data, 'data');
    const res = await Api.workspaceIdNotebookUpdate(
      values.workspace,
      values.name,
      data,
    );
    if (res?.data === null) {
      onClose();
      setTimerModalInfo({ notebookName: values.name });
    } else {
      Message.error('复制 Notebook 失败');
    }
  };

  async function validateNotebookName(value, allValues: FormCopyNotebook) {
    const WorkspaceID = allValues.workspace;
    //  查询当前选中的workspace下的所有notebook
    const { data } = await Api.workspaceIdNotebookList(WorkspaceID);
    // console.log(data, 'werwerw');
    const exist = data.items.some(i => i.name === value);
    if (exist) {
      return 'Notebook 名称不可与目的 Workspace 已有的重复';
    }
  }

  useEffect(() => {
    if (Boolean(name)) {
      Api.workspaceList({ orderBy: 'Name:asc' }).then(({ data }) => {
        setWorkspaceList(data.items);
      });
    }
  }, [name]);
  return (
    <>
      <Modal
        visible={Boolean(name)}
        onCancel={onClose}
        title="复制 Notebook"
        onOk={() => refForm.current?.submit()}
        unmountOnExit={true}
        autoFocus={false}
        maskClosable={false}
        escToExit={false}
      >
        <FinalForm
          subscription={{
            initialValues: true,
          }}
          initialValues={{
            workspace: workspaceId,
            name: name + '-copy',
          }}
          onSubmit={handleSubmit}
        >
          {({ form }) => {
            refForm.current = form;
            return (
              <ArcoForm
                labelAlign="left"
                layout="vertical"
                labelCol={{ span: 5 }}
                wrapperCol={{ span: 19 }}
              >
                <FieldItem
                  name="workspace"
                  label="目的 Workspace"
                  required={true}
                >
                  <Select
                    placeholder="请选择 Workspace"
                    showSearch={true}
                    filterOption={(key, value) =>
                      value.props.children.includes(key)
                    }
                    onChange={value => {
                      setPurposeUrl(value);
                    }}
                  >
                    {workspaceList.map(v => (
                      <Select.Option value={v.id} key={v.id}>
                        {v.name}
                      </Select.Option>
                    ))}
                  </Select>
                </FieldItem>
                <FieldItem
                  key={purposeUrl}
                  name="name"
                  label="Notebook 名称"
                  required={true}
                  rules={[
                    ...getNameRules(1, 60),
                    { asyncValidate: validateNotebookName },
                  ]}
                >
                  <Input placeholder="请输入" />
                </FieldItem>
              </ArcoForm>
            );
          }}
        </FinalForm>
      </Modal>
      <TimerModal
        visible={Boolean(timerModalInfo)}
        title="复制 Notebook 成功"
        onClickStay={() => {
          refetch();
          setTimerModalInfo(undefined);
        }}
        onTimeout={() => {
          history.push(
            `/workspace/${purposeUrl}/notebook/${timerModalInfo?.notebookName}`,
          );
        }}
        renderDesc={countdown => {
          return (
            <div className="flexJustifyCenter flexAlignCenter">
              即将在
              <span style={{ margin: '0 0.5em' }}>{countdown}</span>秒 后进入
              <span
                className="ellipsis"
                style={{ maxWidth: 100, margin: '0 0.5em' }}
              >
                {timerModalInfo?.notebookName}
              </span>
              详情
            </div>
          );
        }}
      />
    </>
  );
}
