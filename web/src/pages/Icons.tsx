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

import Icon from 'components/Icon';

const request = require.context('../assets/svg', false, /\.svg$/);

// 仅供开发预览使用
export default function Icons() {
  return (
    <div className="flexCenter">
      {request.keys().map((k, index) => {
        const child = [
          <Icon
            key={k}
            className="mr8 mt8"
            size={48}
            glyph={k.replace('./icon-', '').replace('.svg', '')}
          />,
        ];
        if (index > 0 && index % 10 === 0) {
          child.push(<br key={`${k}-br`} />);
        }
        return child;
      })}
    </div>
  );
}
