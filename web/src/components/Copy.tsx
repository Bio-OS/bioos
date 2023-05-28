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

import React, { CSSProperties, useState } from 'react';
import copy from 'clipboard-copy';
import { Popover, Typography } from '@arco-design/web-react';
import { IconCheckCircleFill, IconCopy } from '@arco-design/web-react/icon';

interface Props {
  text?: string;
  maxWidth?: number;
  copyValue?: string;
  style?: CSSProperties;
}

const Copy: React.FC<Props> = props => {
  const { text, maxWidth = 150, copyValue, style } = props;
  const [copyed, executeCopy] = useState(false);
  return (
    <div className="flexAlignCenter colorBlack">
      <Typography.Text
        style={{ maxWidth, marginBottom: 0, ...style }}
        className="colorBlack3 fs13"
        ellipsis={{
          cssEllipsis: true,
          showTooltip: {
            type: 'popover',
            props: {
              getPopupContainer: () => document.body,
              content: (
                <div
                  style={{
                    wordBreak: 'break-all',
                    maxHeight: 400,
                    overflow: 'auto',
                    fontSize: 13,
                  }}
                >
                  {text}
                </div>
              ),
            },
          },
        }}
      >
        {text}
      </Typography.Text>

      <Popover
        title={
          copyed ? (
            <div className="colorSuccess fs14">
              <IconCheckCircleFill />
              <span className="ml4">已复制</span>
            </div>
          ) : (
            <span className="fs14">复制</span>
          )
        }
        onVisibleChange={() => {
          executeCopy(false);
        }}
      >
        <IconCopy
          onClick={() => {
            copy(copyValue || text || '');
            executeCopy(true);
          }}
          className="cursorPointer ml4 colorGrey"
        />
      </Popover>
    </div>
  );
};
export default Copy;
