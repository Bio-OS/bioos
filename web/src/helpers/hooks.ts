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

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import queryString from 'query-string';

import Api from 'api/client';
import { ApiserverClientConfig } from 'api/index';

import { GLOBAL_CONFIG_STORAGE_KEY } from './constants';

export function useQuery(size = 12) {
  const { search } = useLocation();
  const parsed = queryString.parse(search);
  const queryNum: { page: number; size: number } = {
    page: parsed.page ? Number(parsed.page) : 1,
    size: parsed.size ? Number(parsed.size) : size,
  };
  return { ...parsed, ...queryNum } as {
    [key: string]: string;
  } & { page: number; size: number };
}

export const useQueryHistory = () => {
  const history = useHistory();
  return (pathname: string, params?: { [key: string]: unknown }) =>
    history.push({
      pathname,
      search: params && queryString.stringify(params, { skipNull: true }),
    });
};

export function useDestroyed() {
  const refDestroyed = useRef(false);

  useEffect(() => {
    return () => {
      refDestroyed.current = true;
    };
  }, []);

  return refDestroyed;
}

export function useRefCallback<T extends (...args: any[]) => any>(callback: T) {
  const callbackRef = useRef(callback);
  callbackRef.current = callback;

  return useCallback((...args: any[]) => callbackRef.current(...args), []) as T;
}

export function useResize<T extends (...args: any[]) => any>(callback: T) {
  const callbackRef = useRefCallback(callback);
  useEffect(() => {
    window.addEventListener('resize', callbackRef);
    return () => {
      window.removeEventListener('resize', callbackRef);
    };
  }, []);
}

export function useResizeMemo<T extends (...args: any[]) => any>(
  callback: T,
  deps: any[],
) {
  const [tick, setTick] = useState(0);
  useResize(() => setTick(v => v + 1));

  // 防止 dom ref 从无到有的时候没有触发 useMemo 的问题
  useEffect(() => setTick(v => v + 1), deps);

  return useMemo(callback, [tick, ...deps]);
}

export function useForceUpdate() {
  const [rerender, setRerender] = useState<boolean>(false);

  return () => {
    setRerender(!rerender);
  };
}

const env = {
  notebook: {},
  storage: {},
};

export function useGetEnvQuery() {
  const [data, setData] = useState<ApiserverClientConfig>(env);
  const sessionStorage = JSON.parse(
    window.sessionStorage.getItem(GLOBAL_CONFIG_STORAGE_KEY),
  );

  useEffect(() => {
    if (sessionStorage) {
      setData(sessionStorage);
    } else {
      Api.configurationList()
        .then(({ data }) => {
          if (data) {
            window.sessionStorage.setItem(
              GLOBAL_CONFIG_STORAGE_KEY,
              JSON.stringify(data),
            );
            setData(data);
          }
        })
        .catch(() => {
          setData(env);
        });
    }
  }, []);

  return data;
}
