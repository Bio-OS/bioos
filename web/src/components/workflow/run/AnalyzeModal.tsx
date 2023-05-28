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

import { ReactNode, useMemo, useRef, useState } from 'react';
import { Form as FinalForm } from 'react-final-form';
import { useRouteMatch } from 'react-router-dom';
import { FormApi } from 'final-form';
import { noop } from 'lodash-es';
import { Form as ArcoForm, Input, Modal } from '@arco-design/web-react';

import { ToastLimitText } from 'components/index';
import TimerModal from 'components/TimerModal';
import { useQueryHistory } from 'helpers/hooks';
import { currentTime, maxLength } from 'helpers/utils';
import { getNameRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';

interface FormValues {
  name: string;
  description?: string;
}

export default function AnalyzeModal({
  title = '分析工作流',
  onOk,
  children,
  workflowName,
  onClickStay,
}: {
  title?: string;
  onOk: (data: FormValues) => Promise<string>;
  children: (open: () => void) => ReactNode;
  workflowName: string;
  onClickStay?: () => void;
}) {
  const [modalStep, setModalStep] = useState<'step-1' | 'step-2' | undefined>();

  const formRef = useRef<FormApi<FormValues>>();
  const submissionId = useRef<string>('');
  const navigate = useQueryHistory();

  const match = useRouteMatch<{ workspaceId: string; workflowId: string }>();
  const { workspaceId, workflowId } = match.params;
  const namePrefix = `${workflowName}-history-`;

  const initialName = useMemo(() => {
    if (modalStep !== 'step-1') return;

    return currentTime();
  }, [modalStep]);

  async function validateSubmissionName(value: string) {
    const res = await Api.submissionDetail(match.params.workspaceId, {
      searchWord: namePrefix + value,
      exact: true,
    });
    const exist = res?.data?.items?.some(
      item => item.name === namePrefix + value,
    );

    return exist ? '名称不可与已有的重复' : undefined;
  }

  function handleTimeout() {
    if (!submissionId.current) return;
    navigate(`/workspace/${workspaceId}/analysis/detail`, {
      workflowID: workflowId,
      submissionID: submissionId.current,
    });
  }

  function open() {
    setModalStep('step-1');
  }

  return (
    <>
      {children(open)}

      <Modal
        title={title}
        visible={modalStep === 'step-1'}
        onOk={async () => {
          if (!formRef.current) return;

          formRef.current.submit();

          const formState = formRef.current.getState();
          if (formState.errors && Object.keys(formState.errors).length) return;

          submissionId.current = await onOk({
            ...formState.values,
            name: namePrefix + formState.values.name,
          });

          setModalStep('step-2');
        }}
        onCancel={() => setModalStep(undefined)}
        unmountOnExit={true}
        autoFocus={false}
        maskClosable={false}
        escToExit={false}
      >
        <FinalForm<FormValues> onSubmit={noop}>
          {({ form }) => {
            formRef.current = form;

            return (
              <ArcoForm labelAlign="left">
                <FieldItem
                  name="name"
                  initialValue={initialName}
                  label="投递名称"
                  required={true}
                  rules={[
                    ...getNameRules(1, 200),
                    {
                      asyncValidate: validateSubmissionName,
                    },
                  ]}
                >
                  <Input
                    placeholder="请输入"
                    addBefore={
                      <div
                        className="ellipsis colorBlack3"
                        style={{ maxWidth: 150 }}
                      >
                        {namePrefix}
                      </div>
                    }
                  />
                </FieldItem>
                <FieldItem
                  name="description"
                  label="简短描述"
                  validate={maxLength(1000)}
                >
                  <Input.TextArea
                    placeholder="请输入"
                    showWordLimit={true}
                    maxLength={{ length: 1000, errorOnly: true }}
                  />
                </FieldItem>
              </ArcoForm>
            );
          }}
        </FinalForm>
      </Modal>

      <TimerModal
        visible={modalStep === 'step-2'}
        title={
          <ToastLimitText
            prefix="工作流"
            suffix="进入分析"
            name={workflowName}
          />
        }
        onClickStay={() => {
          setModalStep(undefined);
          onClickStay?.();
        }}
        onTimeout={handleTimeout}
        renderDesc={countdown => {
          return (
            <div className="flexCenter">
              即将在
              <span className="colorPrimary mr4 ml4">{countdown}</span>
              秒后进入工作流运行详情
            </div>
          );
        }}
      />
    </>
  );
}
