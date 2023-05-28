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
import classNames from 'classnames';
import { IconDown, IconUp } from '@arco-design/web-react/icon';

import { HandlersDataModel } from 'api/index';

import Item from './CategoryItem';

import styles from './Category.less';

interface Props {
  title: string;
  suffix?: React.ReactNode;
  list: HandlersDataModel[];
  activeItem?: string;
  onSelect?: (val: HandlersDataModel) => void;
}

const Category = (props: Props) => {
  const { title, suffix, list, activeItem = list[0], onSelect } = props;
  const [renderChild, setRenderChild] = useState(true);
  const handleSelectItem = item => {
    onSelect?.(item);
  };
  return (
    <div className={styles.category}>
      <div className={classNames(['colorGrey', styles.header])}>
        <span
          className="cursorPointer mr8"
          onClick={() => {
            setRenderChild(!renderChild);
          }}
        >
          {renderChild ? <IconUp /> : <IconDown />}
        </span>
        <span className="flex1">{title}</span>
        {suffix}
      </div>
      <div>
        {renderChild &&
          list.map(item => {
            const active = activeItem === item.id;
            return (
              <Item
                item={item}
                active={active}
                onSelect={handleSelectItem}
                key={item.id}
              />
            );
          })}
        {!list.length && (
          <div style={{ height: 48, padding: '13px 0px 13px 20px' }}>
            暂无数据
          </div>
        )}
      </div>
    </div>
  );
};

Category.Item = Item;
export default Category;
