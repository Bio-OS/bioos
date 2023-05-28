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

import React, { useMemo, useState } from 'react';
import { Link, Modal, Table, Typography } from '@arco-design/web-react';

import Empty from 'components/Empty';

import styles from './analysis.less';

const LogModal: React.FC<{
  log?: string;
  stdout?: string;
  stderr?: string;
}> = ({ stdout, stderr }) => {
  const [visible, setVisible] = useState(false);
  const logData = useMemo(
    () => [
      { type: 'Stdout', path: stdout },
      { type: 'Stderr', path: stderr },
    ],
    [stdout, stderr],
  );
  const columnsLog = [
    {
      title: '日志类型',
      dataIndex: 'type',
      width: '25%',
      render: (value: string) => (
        <Typography.Paragraph
          ellipsis={{
            rows: 1,
          }}
        >
          {value}
        </Typography.Paragraph>
      ),
    },
    {
      title: '日志路径',
      dataIndex: 'path',
      width: '75%',
      render: (value: any) => (
        <Typography.Paragraph
          ellipsis={{
            rows: 5,
            showTooltip: {
              type: 'tooltip',
              props: {
                color: '#ffffff',
                content: <div className={styles.inputTip}>{value}</div>,
              },
            },
          }}
        >
          {value}
        </Typography.Paragraph>
      ),
    },
  ];
  return (
    <>
      <Link onClick={() => setVisible(true)}>运行日志</Link>
      <Modal
        className="pb20"
        visible={visible}
        footer={null}
        title={
          <div className="fw500 fs14" style={{ textAlign: 'left' }}>
            运行日志
          </div>
        }
        maskClosable={true}
        onCancel={() => {
          setVisible(false);
        }}
      >
        <Table
          scroll={{ y: true }}
          border={true}
          columns={columnsLog}
          data={logData}
          pagination={false}
          noDataElement={<Empty desc="暂无数据" />}
        />
      </Modal>
    </>
  );
};

export default LogModal;
