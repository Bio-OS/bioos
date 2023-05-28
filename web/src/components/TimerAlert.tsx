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

import { ReactNode } from 'react';
import { Alert } from '@arco-design/web-react';

import Timer from './Timer';

function TimerAlert({
  content,
  time = 5,
  onClose,
}: {
  content: ReactNode;
  time?: number;
  onClose: () => void;
}) {
  return (
    <Alert
      closable={true}
      onClose={onClose}
      type="warning"
      className="mb20 flexAlignCenter"
      action={
        <span className="colorGrey">
          提示将在
          <Timer
            initialSeconds={time}
            onTimeout={onClose}
            render={seconds => (
              <span className="colorPrimary mr4 ml4">{seconds}s</span>
            )}
          />
          后自动关闭
        </span>
      }
      content={content}
    />
  );
}

export default TimerAlert;
