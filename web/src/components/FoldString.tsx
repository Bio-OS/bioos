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

import { CSSProperties, ReactNode, useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import { Link } from '@arco-design/web-react';

export default function FoldString({
  children,
  lineClamp = 2,
  textFold = '展开',
  className,
  style = {},
}: {
  children?: ReactNode;
  lineClamp?: number;
  textFold?: string;
  className?: string;
  style?: CSSProperties;
}) {
  const [flagShowFold, setflagShowFold] = useState(true);
  const [flagFold, setFlagFold] = useState(true);

  const refDiv = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = refDiv.current;
    if (!el) return;

    if (el.scrollHeight > el.offsetHeight) {
      setflagShowFold(true);
    } else {
      setflagShowFold(false);
    }
  }, [children]);

  if (!children) return <>-</>;

  if (!flagShowFold) {
    return <>{children}</>;
  }

  return (
    <div className="w100 flexAlignItemsEnd">
      <div
        className={classNames('flex1', className)}
        ref={refDiv}
        style={
          flagFold
            ? {
                display: '-webkit-box',
                WebkitBoxOrient: 'vertical',
                WebkitLineClamp: lineClamp,
                overflow: 'hidden',
                marginRight: 8,
                ...style,
              }
            : {
                marginRight: 8,
                minWidth: 0,
                ...style,
              }
        }
      >
        {children}
      </div>

      <div style={{ marginBottom: 2, overflowWrap: 'normal' }}>
        <Link
          className="fs12"
          onClick={() => {
            setFlagFold(!flagFold);
          }}
        >
          {flagFold ? textFold : '收起'}
        </Link>
      </div>
    </div>
  );
}
