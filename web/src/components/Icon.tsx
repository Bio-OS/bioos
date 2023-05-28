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

import { noop } from 'lodash-es';

export interface IconProps {
  className?: string;
  glyph: string;
  size?: number | string;
  width?: number | string;
  height?: number | string;
  style?: React.CSSProperties;
  onClick?: () => void;
}

export default function Icon({
  className,
  glyph,
  size = 24,
  width,
  height,
  style,
  onClick = noop,
}: IconProps) {
  return (
    <svg
      className={className}
      width={width || size}
      height={height || size}
      onClick={onClick}
      style={style}
    >
      <use xlinkHref={`#icon-${glyph}`} />
    </svg>
  );
}
