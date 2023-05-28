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
  Alert,
  Form as ArcoForm,
  Input,
  Message,
  Modal,
  Select,
} from '@arco-design/web-react';
import { IconExclamationCircleFill } from '@arco-design/web-react/icon';

import TimerModal from 'components/TimerModal';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';

import styles from './style.less';
export interface NotebookCreateFormValue {
  name: string;
  language: string;
}
interface Props {
  visible: boolean;
  refetch: () => void;
  onClose: () => void;
}

const LanguageInfoMap = {
  'Python 3': {
    file_extension: '.py',
    mimetype: 'text/x-python',
    name: 'python',
    pygments_lexer: 'ipython3',
    version: '3.9.7',
  },
  'R 语言': {
    file_extension: '.r',
    mimetype: 'text/x-r-source',
    name: 'R',
    pygments_lexer: 'r',
    version: '4.0.3',
  },
};

export default function Add({ visible, onClose, refetch }: Props) {
  const history = useHistory();
  const match = useRouteMatch<{ workspaceId: string }>();
  const formRef = useRef<undefined | FormApi<NotebookCreateFormValue>>();
  const [loading, toggleLoading] = useState(false);
  const [timerModalInfo, setTimerModalInfo] = useState<{
    notebookName: string;
  }>();

  async function validateNotebookName(value) {
    const workspaceID = match.params.workspaceId;
    //  查询当前选中的workspace下的所有notebook
    const { data } = await Api.workspaceIdNotebookList(workspaceID);
    const exist = data.items.some(i => i.name === value);
    if (exist) {
      return 'Notebook 名称不可与已有的重复';
    }
  }

  const handleSubmit = async (values: NotebookCreateFormValue) => {
    toggleLoading(true);
    Api.workspaceIdNotebookUpdate(match.params.workspaceId, values.name, {
      nbformat: 4,
      nbformat_minor: 5,
      cells: [
        {
          cell_type: 'code',
          execution_count: null,
          id: 'ee2a969b',
          metadata: {},
          outputs: [],
          source: [],
        },
      ],
      metadata: {
        kernelspec: {
          display_name:
            values.language === 'Python 3' ? 'Python 3 (ipykernel)' : 'R',
          language: values.language === 'Python 3' ? 'python' : 'R',
          name: values.language === 'Python 3' ? 'python3' : 'ir',
        },
        language_info: LanguageInfoMap[values.language],
      },
    })
      .then(res => {
        // console.log(res.data);
        toggleLoading(false);
        onClose();
        refetch();
        setTimerModalInfo({ notebookName: values.name });
      })
      .catch(error => {
        Message.error('新建 NoteBook 失败！');
      });
  };
  return (
    <>
      <Modal
        className={styles.modalContent}
        title="新建 Notebook"
        visible={visible}
        focusLock={false}
        unmountOnExit={true}
        confirmLoading={loading}
        onConfirm={() => formRef.current.submit()}
        onCancel={onClose}
      >
        <Alert
          className={styles.tipWrap}
          type="warning"
          content={
            <div className="colorText3">
              <div className="mb4">1、Notebook 名称不可与已有的重复；</div>
              <div>2、创建完成Notebook后，会生成.ipynb后缀的文件。</div>
            </div>
          }
        />
        <FinalForm
          subscription={{ initialValues: true }}
          onSubmit={handleSubmit}
        >
          {({ form }) => {
            formRef.current = form;

            return (
              <ArcoForm
                labelAlign="left"
                layout="vertical"
                className="pl24 pr24"
                labelCol={{ span: 5 }}
                wrapperCol={{ span: 19 }}
              >
                <FieldItem
                  name="name"
                  label="名称"
                  required={true}
                  rules={[
                    ...getNameRules(1, 60),
                    {
                      asyncValidate: validateNotebookName,
                    },
                  ]}
                >
                  <Input allowClear={true} placeholder="请输入" />
                </FieldItem>
                <FieldItem name="language" label="选择语言" required={true}>
                  <Select
                    allowClear
                    placeholder="请选择"
                    options={['Python 3', 'R 语言']}
                  ></Select>
                </FieldItem>
              </ArcoForm>
            );
          }}
        </FinalForm>
      </Modal>
      <TimerModal
        visible={Boolean(timerModalInfo)}
        title="新建 Notebook 成功"
        onClickStay={() => {
          setTimerModalInfo(undefined);
        }}
        onTimeout={() => {
          history.push(
            `/workspace/${match.params.workspaceId}/notebook/${timerModalInfo?.notebookName}`,
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
