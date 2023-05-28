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

import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useRequest } from 'ahooks';
import { Button, Modal, Popover, Space } from '@arco-design/web-react';
import { IconPauseCircleFill } from '@arco-design/web-react/icon';

import Api from 'api/client';
import { HertzGetResponse } from 'api/index';

import Loading from '../Loading';

import styles from './style.less';

interface JupyterEditProps {
  title?: string;
  url?: string;
  onCancel: () => void;
  onRefresh: () => Promise<string>;
  getStatus: () => Promise<boolean>;
}

export default function JupyterEdit({
  title,
  url,
  onCancel,
  onRefresh,
  getStatus,
}: JupyterEditProps) {
  const hiddenStartTimeRef = useRef(0);
  const ref = useRef<HTMLIFrameElement | null>(null);
  const { workspaceId, notebookName } = useParams<{
    workspaceId: string;
    notebookName: string;
  }>();

  const [flagFinishJupyterhubLogout, setFlagFinishJupyterhubLogout] =
    useState(false);
  const [notebookEditInfo, setNotebookEditInfo] = useState<HertzGetResponse>();

  const getNotebookServerDetail = async () => {
    const { data: notebookServerList } =
      await Api.workspaceIdNotebookserverList(workspaceId);
    const { data: serverDetail } = await Api.workspaceIdNotebookserverDetail(
      workspaceId,
      notebookServerList[0].id,
      {
        notebook: notebookName,
      },
    );
    setNotebookEditInfo(serverDetail);
    return serverDetail;
  };

  const stopNotebookServer = async () => {
    await Api.workspaceIdNotebookserverCreate2(
      workspaceId,
      notebookEditInfo.id,
      { off: true },
    );
  };

  const { run: start, cancel: stop } = useRequest(getNotebookServerDetail, {
    pollingInterval: 1000,
  });

  const flagReady = notebookEditInfo?.status === 'Running';

  useEffect(() => {
    if (flagReady) {
      stop();
    }
  }, [flagReady]);

  // 离开超过 2 分钟时探活
  useEffect(() => {
    start();
    const fn = async () => {
      if (!flagFinishJupyterhubLogout) return;

      if (!document.hidden) {
        const interval = Date.now() - hiddenStartTimeRef.current;
        if (interval >= 2 * 60 * 1000) {
          // console.info('离开超过2分钟，检测 server 状态', interval);
          if (!(await getStatus())) {
            // console.info('server 离线，进行刷新');
            Modal.destroyAll();
            Modal.confirm({
              title: '当前 Server 已离线，请选择您的操作',
              okText: '重启 Server',
              cancelText: '返回',
              onOk: async () => {
                if (ref.current) {
                  ref.current.src = await onRefresh();
                }
              },
              onCancel,
            });
          } else {
            // console.info('server 状态正常');
          }
        }
      } else {
        hiddenStartTimeRef.current = Date.now();
      }
    };
    document.addEventListener('visibilitychange', fn);
    return () => {
      document.removeEventListener('visibilitychange', fn);
    };
  }, []);

  // 本地调试时更改为本地origin，支持本地调试notebook
  const newUrl = useMemo(() => {
    if (!url) return;
    return url;
  }, [url]);

  const handleLoad = () => {
    const iframeDocument = ref.current?.contentWindow?.document;
    if (!iframeDocument) return;
    const targetNode = iframeDocument.body;
    const config = { childList: true, subtree: true };
    const observer = new MutationObserver(() => {
      iframeDocument.getElementById('login_widget')?.remove();
      iframeDocument
        .getElementsByClassName(
          'btn btn-default btn-sm navbar-btn pull-right',
        )?.[0]
        ?.remove();
    });
    observer.observe(targetNode, config);
  };

  useEffect(() => {
    if (!newUrl) return;
    // const urlObj = new URL(window.location.origin + newUrl);
    const iframe = document.createElement('iframe');

    iframe.style.width = '0px';
    // 适配后端notebook URL前缀更改为动态变化
    iframe.src = newUrl.replace(/login$/, 'logout');
    iframe.onload = () => {
      setFlagFinishJupyterhubLogout(true);

      document.body.removeChild(iframe);
    };

    document.body.appendChild(iframe);
  }, [newUrl]);

  if (!flagFinishJupyterhubLogout) {
    return <Loading />;
  }

  return (
    <div className={styles.jupyterEditContainer}>
      <div className={styles.header}>
        <span>{title}</span>
        <Space>
          {flagReady ? (
            <Button
              onClick={() => {
                Modal.confirm({
                  title: '确定要停止 Jupyter 实例吗？ ',
                  content: (
                    <div style={{ fontSize: 12 }}>
                      <div>
                        停止 Jupyter
                        实例后，/home/jovyan目录下的数据会保存，其余目录下内容会被清理。
                      </div>
                    </div>
                  ),
                  okButtonProps: {
                    status: 'danger',
                  },
                  escToExit: false,
                  maskClosable: false,
                  okText: '停止',
                  async onOk() {
                    await stopNotebookServer();

                    onCancel();
                  },
                });
              }}
              icon={<IconPauseCircleFill className="fs13 colorText2" />}
            >
              停止
            </Button>
          ) : (
            <Popover content="Notebook 启动中，无法停止" position="left">
              <Button
                disabled={true}
                icon={<IconPauseCircleFill className="fs13 colorF4" />}
              >
                停止
              </Button>
            </Popover>
          )}

          <Button onClick={onCancel}>取消编辑</Button>
        </Space>
      </div>
      <iframe
        id="iframe"
        title="jupyterEdit"
        className={styles.iframeWrap}
        src={newUrl}
        ref={ref}
        onLoad={handleLoad}
      />
    </div>
  );
}
