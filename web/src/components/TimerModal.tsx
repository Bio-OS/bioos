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

import { ReactNode } from 'react';
import { Button, Modal } from '@arco-design/web-react';
import { IconCheckCircleFill } from '@arco-design/web-react/icon';

import ImgWorkspace from 'assets/imgs/ImgWorkspace.png';

import Timer from './Timer';

export default function TimerModal({
  visible,
  title,
  onClickStay,
  onTimeout,
  renderDesc,
}: {
  visible: boolean;
  title: ReactNode;
  onClickStay: () => void;
  onTimeout: () => void;
  renderDesc;
}) {
  return (
    <Modal
      visible={visible}
      simple={true}
      unmountOnExit={true}
      footer={null}
      maskClosable={false}
      escToExit={false}
      autoFocus={false}
      style={{ width: 440 }}
    >
      <div className="w100">
        <div className="flexAlignCenter flexJustifyCenter mb20">
          <IconCheckCircleFill
            className="mr8 fs20"
            style={{ color: '#00AA2A' }}
          />
          <span className="fw500 fs14 ellipsis">{title}</span>
        </div>

        <Timer
          initialSeconds={3}
          render={renderDesc}
          onTimeout={onTimeout}
        ></Timer>

        <img
          title="imgWorkspace"
          className="displayBlock marginAuto mb24 mt24"
          style={{ width: 120 }}
          src={ImgWorkspace}
        />

        <div className="textAlignCenter">
          <Button onClick={onClickStay}>我想留在当前页面</Button>
        </div>
      </div>
    </Modal>
  );
}
