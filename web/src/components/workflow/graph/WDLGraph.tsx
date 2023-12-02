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
import svgPanZoom from 'svg-pan-zoom';
import Viz from 'viz.js';
import { Module, render } from 'viz.js/full.render';

let viz = new Viz({ Module, render });

export default function ({
  data,
  container,
  zoom,
}: {
  data: string;
  container: React.MutableRefObject<HTMLDivElement>;
  zoom: React.MutableRefObject<SvgPanZoom.Instance>;
}) {
  viz
    ?.renderSVGElement(data, { yInvert: false })
    .then(result => {
      container?.current?.append(result);
      zoom.current = svgPanZoom(result, {
        zoomEnabled: true,
        fit: true,
        center: true,
      });
    })
    .catch(error => {
      viz = new Viz({ Module, render });
      console.error(error);
    });

  return <></>;
}
