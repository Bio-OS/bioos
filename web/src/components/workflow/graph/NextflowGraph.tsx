/**
 *
 * Copyright 2023 Beijing Volcano Engine Technology Ltd.
 * Copyright 2023 Guangzhou Laboratory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useRef } from 'react';
import mermaid from 'mermaid';
import svgPanZoom from 'svg-pan-zoom';

export default function ({
  data,
  zoom,
}: {
  data: string;
  zoom: React.MutableRefObject<SvgPanZoom.Instance>;
}) {
  const svgWrapper = useRef<HTMLDivElement>(null);

  mermaid.initialize({ startOnLoad: false });
  mermaid
    .parse(data, {
      suppressErrors: true,
    })
    .then(parsed => {
      if (parsed) {
        mermaid.render('theGraph', data).then(result => {
          svgWrapper.current.innerHTML = result.svg;
          svgWrapper.current.querySelector('svg').style.maxWidth = null;
          zoom.current = svgPanZoom(svgWrapper.current.querySelector('svg'), {
            zoomEnabled: true,
            fit: true,
            center: true,
          });
        });
      }
    });

  return (
    <div
      ref={svgWrapper}
      id="svgContainer"
      style={{
        width: '100%',
        height: '100%',
      }}
    />
  );
}
