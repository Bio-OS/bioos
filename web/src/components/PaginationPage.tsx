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

import { ReactNode, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Pagination, Spin } from '@arco-design/web-react';

import { useQuery, useQueryHistory } from 'helpers/hooks';

export interface Props {
  children: ReactNode;
  total: number;
  loading?: boolean;
  sizeOptions?: number[];
  defaultSize?: number;
}

export default function PaginationPage({
  children,
  loading,
  total = 0,
  sizeOptions = [12, 24, 36, 72],
  defaultSize = 12,
}: Props) {
  const { pathname } = useLocation();
  const navigate = useQueryHistory();
  const query = useQuery(defaultSize);
  const { page, size } = query;
  function handleChangePage(current: number, pageSize: number) {
    const queryMap = { ...query, page: current, size: pageSize };
    navigate(pathname, queryMap);
  }
  useEffect(() => {
    // 当前页无数据自动跳转前一页
    if (total !== 0 && page > 1 && total <= size * (page - 1)) {
      const lastPage = Math.ceil(total / size);
      const queryMap = { ...query, page: lastPage };
      navigate(pathname, queryMap);
    }
  }, [total, page]);

  return (
    <>
      <Spin loading={loading} block>
        {children}
      </Spin>
      <Pagination
        className="flexJustifyEnd mt16"
        total={total}
        hideOnSinglePage={size === defaultSize}
        sizeOptions={sizeOptions}
        showTotal
        sizeCanChange
        pageSizeChangeResetCurrent
        pageSize={size}
        current={page}
        onChange={handleChangePage}
      />
    </>
  );
}
