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

import { memo, ReactNode } from 'react';
import { Grid } from '@arco-design/web-react';

import { useQuery } from 'helpers/hooks';

import PageEmpty from './Empty';

const { Row, Col } = Grid;

export interface Props<T> {
  uniqKey?: string;
  renderItem: (item: T) => ReactNode;
  data: T[];
  loading?: boolean;
  gutter?: number[];
}

function CardList<T>({
  data,
  uniqKey = 'id',
  renderItem,
  gutter = [16, 16],
  loading,
}: Props<T>) {
  const { search } = useQuery();
  if (loading) {
    return;
  }
  return (
    <Row gutter={gutter}>
      {data?.length ? (
        data.map(item => (
          <Col key={item[uniqKey]} xs={6} sm={6} xxxl={4}>
            {renderItem(item)}
          </Col>
        ))
      ) : (
        <PageEmpty search={search} />
      )}
    </Row>
  );
}

export default memo(CardList);
