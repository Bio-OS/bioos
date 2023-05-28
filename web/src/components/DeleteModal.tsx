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

import { memo, ReactNode, useState } from 'react';
import { Input, Modal } from '@arco-design/web-react';

import CommonLimitText from 'components/CommonLimitText';

export interface Props {
  type: string;
  title?: string | ReactNode;
  visible: boolean;
  tips?: string | string[];
  name?: string;
  width?: number;
  verify?: boolean;
  useTextArea?: boolean;
  showLimitText?: boolean;
  onClose: () => void;
  onDelete: () => void;
}

function DeleteModal({
  type,
  title,
  name,
  visible,
  tips,
  width = 400,
  verify = false,
  useTextArea = false,
  showLimitText = true,
  onClose,
  onDelete,
}: Props) {
  const [verifyValue, setVerifyValue] = useState('');
  const [error, setError] = useState('');

  function handleVerify() {
    if (!verifyValue) {
      setError(`请输入${type}名称，以确认删除`);
    } else if (verifyValue !== name) {
      setError(`${type} 名称错误，请重新输入`);
    } else {
      setError('');
    }
  }

  const InputCom = useTextArea ? Input.TextArea : Input;

  return (
    <Modal
      title={title || `确定删除所选${type}吗？`}
      visible={visible}
      style={{ width }}
      unmountOnExit
      onCancel={onClose}
      onOk={onDelete}
      okText="删除"
      okButtonProps={{
        status: 'danger',
        disabled: verify && verifyValue !== name,
      }}
    >
      {!!tips && (
        <div className="fs13 colorBlack3 mb12">
          {tips instanceof Array
            ? tips.map((tip, index) => {
                return (
                  <div key={tip}>
                    {index + 1}、{tip}
                  </div>
                );
              })
            : tips}
        </div>
      )}

      {verify && (
        <div className="mb20" style={{ position: 'relative' }}>
          <InputCom
            error={!!error}
            autoSize={{ minRows: 1, maxRows: 6 }}
            allowClear={true}
            onClear={() => setVerifyValue('')}
            placeholder={`请输入${type}名称，以确认删除`}
            onChange={setVerifyValue}
            onBlur={handleVerify}
          />
          {error && (
            <span
              className="fs12 colorDanger"
              style={{ position: 'absolute', bottom: -20, left: 0 }}
            >
              {error}
            </span>
          )}
        </div>
      )}
      {showLimitText && <CommonLimitText name={type} value={name} />}
    </Modal>
  );
}

export default memo(DeleteModal);
