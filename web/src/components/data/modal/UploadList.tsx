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

import React, { memo, useEffect, useMemo, useState } from 'react';
import {
  IconCheck,
  IconClose,
  IconCloseCircleFill,
} from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import { getSize } from 'helpers/utils';

import { FileInfo } from './ImportDataModal';

import styles from './UploadList.less';
interface Props {
  list: FileInfo[];
  onClose?: () => void;
}
const UploadList: React.FC<Props> = props => {
  const { list, onClose } = props;
  const [uploadVisible, setUploadVisible] = useState(false);
  useEffect(() => {
    setUploadVisible(list.some(item => item?.status !== 'success'));
  }, [list]);
  return (
    <div
      className={styles.list}
      style={{
        display: uploadVisible ? 'block' : 'none',
      }}
    >
      <div className="flexJustifyBetween">
        <span className="colorBlack fs14 fw500 mb12">传输列表</span>
        <IconClose
          className="cursorPointer"
          strokeWidth={8}
          onClick={() => {
            setUploadVisible(false);
            onClose?.();
          }}
        />
      </div>
      <div className={styles.content}>
        {list
          ?.filter(item => item?.status !== 'success')
          .map(item => {
            return (
              <div key={item?.id} className="flexBetween pt12 pb12">
                <div className="flexBetween">
                  <Icon glyph="file" style={{ width: 16, height: 20 }} />
                  <span className="ml12 colorGrey1">{item?.name}</span>
                  <span className="ml12 colorGrey">
                    {getSize(item?.size || 0)}
                  </span>
                </div>
                {item?.status === 'init' && <Icon glyph="loading" />}
                {item?.status === 'success' && (
                  <IconCheck className="colorSuccess" />
                )}
                {item?.status === 'error' && (
                  <IconCloseCircleFill className="colorDanger" />
                )}
              </div>
            );
          })}
      </div>
    </div>
  );
};
export default memo(UploadList);
