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
import classNames from 'classnames';
import { Typography } from '@arco-design/web-react';
import { IconCheck } from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import { HandlersDataModel } from 'api/index';

import styles from './Category.less';

interface ItemProps {
  item: HandlersDataModel;
  active: boolean;
  onSelect?: (val: HandlersDataModel) => void;
}
const Item: React.FC<ItemProps> = props => {
  const { item, active, onSelect } = props;
  return (
    <div
      className={classNames([
        styles.listItem,
        { [styles.listItemActive]: active },
      ])}
      onClick={() => {
        onSelect?.(item);
      }}
    >
      <Icon
        glyph="data"
        size={16}
        className={classNames([styles.dataIcon, { [styles.active]: active }])}
      />
      <Typography.Text
        ellipsis={{ rows: 1, showTooltip: { type: 'popover' } }}
        style={{
          color: active ? PRIMARY_COLOR : '#020814',
          fontSize: 13,
          maxWidth: 140,
          flex: 1,
        }}
      >
        {item.name}
      </Typography.Text>
      <div className="flexJustifyEnd" style={{ width: 50 }}>
        <span className={styles.count}>{`(${item.rowCount})`}</span>
        <IconCheck className={styles.pinCheck} />
      </div>
    </div>
  );
};
export default Item;
