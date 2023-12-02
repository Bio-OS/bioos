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

import { memo, useEffect, useRef, useState } from 'react';
import classNames from 'classnames';
import { Popover } from '@arco-design/web-react';

import Icon from 'components/Icon';
import NextflowGraph from 'components/workflow/graph/NextflowGraph';
import WDLGraph from 'components/workflow/graph/WDLGraph';
import { Z_INDEX } from 'helpers/constants';

import styles from './style.less';

const ZOOM_ACTIONS = [
  {
    action: 'zoomIn',
    title: '放大',
    icon: <Icon glyph="zoom-in" />,
  },
  {
    action: 'zoomOut',
    title: '缩小',
    icon: <Icon glyph="zoom-out" />,
  },
  {
    action: 'full',
    title: '全屏',
    icon: <Icon glyph="full-screen" />,
  },
  {
    action: 'reset',
    title: '重置',
    icon: <Icon glyph="reset" />,
  },
];

function DAGToGraph({ data, language }: { data: string; language: string }) {
  const graphRef = useRef<HTMLDivElement>(null);
  const zoomRef = useRef<SvgPanZoom.Instance>(null);
  const [showFullScreen, setShowFullScreen] = useState(false);

  if (!data) return null;

  function resetSvg() {
    zoomRef.current.resetZoom();
    zoomRef.current.resetPan();
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      resetSvg();
      setShowFullScreen(!showFullScreen);
    }
  }

  useEffect(() => {
    if (showFullScreen) {
      document.addEventListener('keydown', handleKeyDown);
      return () => document.removeEventListener('keydown', handleKeyDown);
    }
  }, [showFullScreen]);

  function handleZoom(action: string) {
    switch (action) {
      case 'zoomIn':
        zoomRef.current.zoomIn();
        break;
      case 'zoomOut':
        zoomRef.current.zoomOut();
        break;
      case 'reset':
        resetSvg();
        break;
      case 'full':
        resetSvg();
        setShowFullScreen(!showFullScreen);
        break;
      default:
        break;
    }
  }

  ZOOM_ACTIONS[2].icon = showFullScreen ? (
    <Icon glyph="exit-full-screen" />
  ) : (
    <Icon glyph="full-screen" />
  );
  ZOOM_ACTIONS[2].title = showFullScreen ? '退出全屏' : '全屏';

  return (
    <div
      className={classNames([
        styles.graphWrap,
        { [styles.fullScreen]: showFullScreen },
      ])}
      style={{ zIndex: Z_INDEX.modal, cursor: 'pointer' }}
    >
      <div ref={graphRef}>
        <div className={styles.actions}>
          {ZOOM_ACTIONS.map(item => (
            <Popover key={item.title} title={item.title} position="left">
              <div onClick={() => handleZoom(item.action)}>{item.icon}</div>
            </Popover>
          ))}
        </div>

        {language === 'WDL' ? (
          <WDLGraph data={data} container={graphRef} zoom={zoomRef} />
        ) : (
          <NextflowGraph data={data} zoom={zoomRef} />
        )}
      </div>
    </div>
  );
}

export default memo(DAGToGraph);
