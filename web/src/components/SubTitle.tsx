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

import { memo } from 'react';
import classNames from 'classnames';

import style from './style.less';

export interface Props {
  title: string;
  bg?: string;
  className?: string;
}

function SubTitle({ title, bg, className }: Props) {
  return (
    <div
      className={classNames([style.subTitle, className])}
      style={{ backgroundColor: bg }}
    >
      {title}
    </div>
  );
}

export default memo(SubTitle);
