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

import { useEffect, useState } from 'react';
import { useRouteMatch } from 'react-router-dom';
import { useRequest } from 'ahooks';
import { Button, Message, Popover } from '@arco-design/web-react';
import {
  IconEdit,
  IconExclamationCircleFill,
  IconLoading,
} from '@arco-design/web-react/icon';

import { useDestroyed } from 'helpers/hooks';
import Api from 'api/client';
import { HertzGetResponse } from 'api/index';

import { isNotebookServerOk } from '../../pages/notebook/index';

import styles from './style.less';

interface Props {
  notebookName: string;
  type?: 'button' | 'icon';
  setVisibleConfigureRuntime?: React.Dispatch<React.SetStateAction<boolean>>;
  usage: 'dashboard' | 'notebook';
}

export default function EditNotebookButton({
  notebookName,
  type = 'button',
  setVisibleConfigureRuntime,
  usage,
}: Props) {
  const match = useRouteMatch<{ workspaceId: string }>();
  const workspaceId = match.params.workspaceId;
  const [btnStatus, setBtnStatus] = useState<
    'init' | 'disabled' | 'loading' | 'starting' | 'configServer'
  >('init');

  const [visibleAlertEnvRestart, setVisibleAlertEnvRestart] = useState(false);

  const [serverDetail, setServerDetail] = useState<HertzGetResponse>();

  const flagDestroyed = useDestroyed();

  const getNotebookServerDetail = async () => {
    const { data: notebookServerList } =
      await Api.workspaceIdNotebookserverList(workspaceId);
    if (!notebookServerList.length) {
      return;
    }
    const { data: serverDetail } = await Api.workspaceIdNotebookserverDetail(
      workspaceId,
      notebookServerList[0]?.id,
      {
        notebook: notebookName,
      },
    );
    setServerDetail(serverDetail);
    return serverDetail;
  };

  const turnOnServer = async id => {
    const data = Api.workspaceIdNotebookserverCreate2(workspaceId, id, {
      on: true,
    });
  };

  const { run: start, cancel: stop } = useRequest(getNotebookServerDetail, {
    pollingInterval: 1000,
  });

  useEffect(() => {
    getNotebookServerDetail();
    return () => {
      stop();
    };
  }, []);

  const notebookServerOk = isNotebookServerOk(serverDetail?.status);

  useEffect(() => {
    if (notebookServerOk && btnStatus === 'starting') {
      Message.success('环境已启动');
      setBtnStatus('init');
      stop();
      window.open(`${match.url}/edit`);
    }
  }, [btnStatus, notebookServerOk]);

  async function handleClick() {
    if (btnStatus === 'loading') return;

    setBtnStatus('loading');

    const currentServerDetail = await getNotebookServerDetail();

    // 检查是否配置server
    if (!(currentServerDetail?.image || currentServerDetail?.resourceSize)) {
      setBtnStatus('configServer');
      setVisibleAlertEnvRestart(true);
      return;
    }

    if (flagDestroyed.current) return;

    setVisibleAlertEnvRestart(true);
    turnOnServer(currentServerDetail.id);
    start();
    setBtnStatus('starting');
    getNotebookServerDetail();
  }

  const renderIcon = () => {
    if (btnStatus === 'loading') {
      return <IconLoading />;
    }

    if (btnStatus === 'disabled') {
      return (
        <IconEdit style={{ color: '#94c2ff' }} className="disabledPointer" />
      );
    }

    return (
      <IconEdit
        className="cursorPointer"
        onClick={() => {
          handleClick();
        }}
      />
    );
  };

  const renderButton = () => {
    return (
      <Button
        loading={btnStatus === 'loading'}
        disabled={btnStatus === 'disabled'}
        type="primary"
        icon={<IconEdit />}
        onClick={() => {
          handleClick();
        }}
      >
        编辑
      </Button>
    );
  };

  const btn = type === 'button' ? renderButton() : renderIcon();

  if (
    btnStatus === 'init' ||
    btnStatus === 'loading' ||
    btnStatus === 'disabled'
  ) {
    return btn;
  } else if (btnStatus === 'configServer') {
    // 未进行 Server 配置
    return (
      <Popover
        className={styles.popoverAlert}
        style={{ maxWidth: 500 }}
        popupVisible={visibleAlertEnvRestart}
        onVisibleChange={() => {
          setVisibleAlertEnvRestart(false);
          setBtnStatus('init');
        }}
        content={
          <div className="dpfx alignCenter" style={{ color: 'white' }}>
            <IconExclamationCircleFill className="mr8 fs20" />
            <span className="noWrap">
              首次编辑 Notebooks 请先进行
              <span
                className="colorPrimary cursorPointer ml4"
                onClick={() => setVisibleConfigureRuntime?.(true)}
              >
                运行资源配置
              </span>
            </span>
          </div>
        }
      >
        {btn}
      </Popover>
    );
  } else {
    return (
      <Popover
        popupVisible={visibleAlertEnvRestart}
        content={
          <>
            <div className="flexAlignCenter" style={{ marginBottom: 6 }}>
              <IconExclamationCircleFill
                className="mr8 fs20"
                style={{ color: '#ff7d00' }}
              />
              <span className="fs14 fw500" style={{ color: '#1a2233' }}>
                环境启动中，请耐心等待
              </span>
            </div>
            <div>
              首次启动或更新运行资源配置后，存在一段环境启动时间，请耐心等待。
            </div>
          </>
        }
      >
        {btn}
      </Popover>
    );
  }
}
