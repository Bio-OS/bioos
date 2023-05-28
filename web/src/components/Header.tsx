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
import { Link, useHistory } from 'react-router-dom';
import classNames from 'classnames';

import Icon from './Icon';

import style from './style.less';

const TABS = [
  {
    name: 'Workspace',
    url: '/',
  },
];

function Header() {
  const history = useHistory();

  function goHome() {
    history.push('/workspace');
  }
  return (
    <div className={classNames('flexAlignCenter', [style.headerWrap])}>
      <div className="flexCenter cursorPointer" onClick={goHome}>
        <Icon glyph="logo" className="colorPrimary" />
        <span className="fs16 fw600 colorBlack">Bio-OS</span>
      </div>
      <div className={style.tabs}>
        {TABS.map(tab => (
          <Link
            key={tab.name}
            className={classNames('fw500 colorBlack', { [style.active]: true })}
            to={tab.url}
          >
            <span>{tab.name}</span>
          </Link>
        ))}
      </div>
    </div>
  );
}

export default memo(Header);
