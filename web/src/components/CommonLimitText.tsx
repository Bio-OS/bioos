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

import { CSSProperties } from 'react';
import { Typography } from '@arco-design/web-react';

const CommonLimitText: React.FC<{
  name: string;
  value?: string;
  style?: CSSProperties;
}> = ({ name, value = '', style }) => {
  return (
    <Typography.Text
      className="fs12 colorGrey"
      style={{ ...style }}
      ellipsis={{
        showTooltip: {
          type: 'popover',
          props: {
            content: value,
            position: 'right',
            style: { whiteSpace: 'break-spaces' },
            getPopupContainer() {
              return document.body;
            },
          },
        },
      }}
    >
      {`当前${name}：${value}`}
    </Typography.Text>
  );
};

export default CommonLimitText;
