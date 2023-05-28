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

import { memo, ReactNode, useEffect, useRef, useState } from 'react';
import { Beforeunload } from 'react-beforeunload';
import { useHistory } from 'react-router-dom';
import { Button, Modal, Space } from '@arco-design/web-react';

function BeforeUnloadModal({
  title,
  content,
  cancelText = '不保存',
  confirmText = '保存',
  onSave,
}: {
  title: ReactNode;
  content: ReactNode;
  cancelText?: string;
  confirmText?: string;
  onSave?: () => void;
}) {
  const [visible, setVisible] = useState(false);
  const history = useHistory();
  const unblockRef = useRef<() => void>();
  const nextPath = useRef<string>('');

  function handler(event) {
    event.preventDefault();
    return (event.returnValue = 'Are you sure want to exit?');
  }

  function leave() {
    setVisible(false);
    history.block(true);
    history.push(nextPath.current);
  }

  function cancel() {
    leave();
  }

  function confirm() {
    leave();
    onSave?.();
  }

  function prompt(location): string | false | void {
    nextPath.current = location.pathname;
    setVisible(true);
    return false;
  }

  useEffect(() => {
    unblockRef.current = history.block(prompt);

    return () => unblockRef.current?.();
  }, []);

  return (
    <>
      <Beforeunload onBeforeunload={handler} />
      <Modal
        style={{ width: 360 }}
        title={title}
        visible={visible}
        onCancel={() => setVisible(false)}
        footer={
          <Space>
            <Button onClick={cancel}>{cancelText}</Button>
            <Button type="primary" onClick={confirm}>
              {confirmText}
            </Button>
          </Space>
        }
      >
        <span className="colorBlack3">{content}</span>
      </Modal>
    </>
  );
}

export default memo(BeforeUnloadModal);
