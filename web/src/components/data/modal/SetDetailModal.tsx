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
import { Modal, Table } from '@arco-design/web-react';

import { genHighlightText } from 'components/getHighLight';

interface tableDataProps {
  key: number;
  name: string | JSX.Element[];
}

const SetDetailModal: React.FC<{
  data: string[];
  name?: string;
  title?: string;
  visible: boolean;
  keyword?: string;
  onOK: () => void;
  onCancel: () => void;
}> = ({ data = [], title, name, visible, keyword, onOK, onCancel }) => {
  const dataSort = data?.sort?.();
  const [sorter, setSorter] = useState('ascend');
  const [paginationInfo, setPaginationInfo] = useState({
    current: 1,
    pageSize: 10,
  });

  const { columns, tableData } = useMemo(() => {
    const columns = [
      {
        title: `${name}_id`,
        dataIndex: 'name',
        sorter: true,
      },
    ];
    let tableData: tableDataProps[];
    if (keyword) {
      const preData = [];
      const lastData = [];
      dataSort.forEach(item => {
        if (item.includes(keyword)) {
          preData.push(item);
        } else {
          lastData.push(item);
        }
      });
      tableData = [...preData, ...lastData].map((item, index) => {
        return {
          key: index,
          name: genHighlightText(item, keyword, '#1D2129', '#BDDCFF', false),
          sortName: item,
        };
      });
    } else {
      tableData = dataSort?.map((item, index) => {
        return { key: index, name: item, sortName: item };
      });
    }
    return {
      columns,
      tableData: sorter === 'ascend' ? tableData : tableData.reverse(),
    };
  }, [keyword, dataSort, sorter]);
  const getCurrentData = () => {
    const { current, pageSize } = paginationInfo;
    return tableData?.slice((current - 1) * pageSize, current * pageSize);
  };
  return (
    <Modal
      unmountOnExit={true}
      key={dataSort?.[0]}
      style={{ width: 800 }}
      onOk={() => {
        setPaginationInfo({ current: 1, pageSize: 10 });
        onOK();
      }}
      onCancel={() => {
        setPaginationInfo({ current: 1, pageSize: 10 });
        onCancel();
      }}
      visible={visible}
      title={<div style={{ textAlign: 'left' }}>{title}</div>}
    >
      {keyword && (
        <div className="flexJustifyStart flexAlignCenter mb12">
          与关键词匹配的
          <div className="colorB1 ml4 mr4">
            {dataSort?.filter(item => item.includes(keyword))?.length}
          </div>
          条数据，已用
          <div
            className="ml4 mr4"
            style={{
              width: 14,
              height: 14,
              background: '#BDDCFF',
              marginTop: 2,
            }}
          ></div>
          高亮
        </div>
      )}
      <Table
        showSorterTooltip={false}
        key={dataSort?.[0]}
        data={getCurrentData()}
        columns={columns}
        border={true}
        scroll={{ y: 450 }}
        pagination={{
          current: paginationInfo.current,
          total: dataSort?.length,
          hideOnSinglePage: dataSort?.length <= 10,
          pageSize: paginationInfo.pageSize,
          sizeOptions: [10, 20, 30, 40, 50, 100],
          showTotal: true,
          sizeCanChange: true,
          onChange: (current, pageSize) => {
            setPaginationInfo({ current, pageSize });
          },
        }}
        onChange={(_, sorter) => {
          sorter.direction && setSorter(sorter.direction);
        }}
      />
    </Modal>
  );
};
export default SetDetailModal;
