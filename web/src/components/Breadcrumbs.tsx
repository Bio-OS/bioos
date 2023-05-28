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

import React from 'react';
import { useHistory } from 'react-router-dom';
import classNames from 'classnames';
import { Breadcrumb, Typography } from '@arco-design/web-react';
import { IconRight } from '@arco-design/web-react/icon';

import { useQueryHistory } from 'helpers/hooks';

const BreadcrumbItem = Breadcrumb.Item;

/**
 * @param {Array} breadcrumbs
 * @param {Object} style
 * @param {string} className
 */
interface BreadcrumbspProps {
  breadcrumbs?: {
    text: React.ReactNode | string;
    path?: string;
    query?: { [key: string]: unknown };
  }[];
  style?: object;
  className?: string;
}
const Breadcrumbs: React.FC<BreadcrumbspProps> = ({
  breadcrumbs,
  style,
  className,
}) => {
  const navigate = useQueryHistory();
  return (
    <Breadcrumb separator={<IconRight />} style={style} className={className}>
      {breadcrumbs?.map(({ path, text, query }, index) => {
        return (
          <BreadcrumbItem
            key={index}
            className={classNames({
              cursorPointer: index < breadcrumbs?.length - 1,
            })}
            onClick={() => path && navigate(path, query)}
          >
            <Typography.Text style={{ maxWidth: 300 }} ellipsis={true}>
              {text}
            </Typography.Text>
          </BreadcrumbItem>
        );
      })}
    </Breadcrumb>
  );
};
export default Breadcrumbs;
