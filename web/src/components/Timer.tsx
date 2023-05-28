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

import { ReactElement, ReactNode, useEffect, useRef, useState } from 'react';

import { useDestroyed } from 'helpers/hooks';

export default function Timer({
  initialSeconds,
  render,
  onTimeout,
}: {
  /**
   * 倒计时初始值（单位：秒）
   */
  initialSeconds: number;
  render: (seconds: number) => ReactNode;
  onTimeout?: () => void;
}) {
  const [seconds, setSeconds] = useState(initialSeconds);
  const timer = useRef<number>();
  const refDestroyed = useDestroyed();

  function startTimer() {
    timer.current = window.setTimeout(() => {
      if (seconds === 0 || refDestroyed.current) return;

      const next = seconds - 1;
      setSeconds(next);

      if (next === 0) {
        onTimeout?.();

        return;
      }

      startTimer();
    }, 1000);
  }

  useEffect(() => {
    startTimer();
    return () => window.clearTimeout(timer.current);
  }, [seconds]);

  // https://github.com/DefinitelyTyped/DefinitelyTyped/issues/18051
  // open issue
  // 实际react可以render ReactNode，但这里如果返回ReactNode类型，ts会报错，因此使用as转换为ReactElement，
  return render(seconds) as ReactElement;
}
