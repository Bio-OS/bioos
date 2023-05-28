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

import React, { useState } from 'react';
import { Link, Modal, Tabs } from '@arco-design/web-react';

import InputParamate from './InputParamate';
import OutputParamate from './OutputParamate';

const TabPane = Tabs.TabPane;

const InputOutputModal: React.FC<{ inputs?: string; outputs?: string }> = ({
  inputs,
  outputs,
}) => {
  const [visible, setVisible] = useState(false);

  return (
    <>
      <Link onClick={() => setVisible(true)}>查看</Link>
      <Modal
        className="pb20"
        visible={visible}
        footer={null}
        title={
          <div className="fw500 fs14" style={{ textAlign: 'left' }}>
            输入输出
          </div>
        }
        maskClosable={true}
        onCancel={() => {
          setVisible(false);
        }}
      >
        <Tabs defaultActiveTab="input" type="card-gutter">
          <TabPane key="input" title="输入参数">
            <InputParamate inputs={inputs} />
          </TabPane>
          <TabPane key="output" title="输出参数">
            <OutputParamate outputs={outputs} />
          </TabPane>
        </Tabs>
      </Modal>
    </>
  );
};

export default InputOutputModal;
