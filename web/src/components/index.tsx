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

import { Typography } from '@arco-design/web-react';

export function Line({ height = 14 }: { height?: number }) {
  return (
    <span
      className="mr12 ml12 inlineBlock"
      style={{
        borderRight: '1px solid #eaedf1',
        height,
      }}
    ></span>
  );
}

export function ToastLimitText({
  name,
  prefix,
  suffix,
  maxWidth = 200,
}: {
  name: string;
  prefix: string;
  suffix: string;
  maxWidth?: number;
}) {
  return (
    <div className="flexAlignCenter">
      {prefix}
      <Typography.Text
        className="ml4 mr4"
        style={{ maxWidth }}
        ellipsis={{
          cssEllipsis: true,
        }}
      >
        {name}
      </Typography.Text>
      {suffix}
    </div>
  );
}
