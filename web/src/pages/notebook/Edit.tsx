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

import React, { useEffect, useState } from 'react';
import { useHistory, useParams } from 'react-router-dom';

import JupyterEdit from 'components/notebook/JupyterEdit';
import Api from 'api/client';
import { HertzGetResponse } from 'api/index';

import { isNotebookServerOk } from '.';

export default function NotebookEdit() {
  const { workspaceId, notebookName } = useParams<{
    workspaceId: string;
    notebookName: string;
  }>();
  const history = useHistory();
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
  useEffect(() => {
    getNotebookServerDetail();
  }, []);

  return (
    <JupyterEdit
      title={notebookName}
      url={notebookEditInfo?.accessURL}
      onCancel={() => window.close()}
      getStatus={async () => {
        const serverDetail = await getNotebookServerDetail();
        return isNotebookServerOk(serverDetail.status);
      }}
      onRefresh={async () => {
        const serverDetail = await getNotebookServerDetail();
        return serverDetail.accessURL;
      }}
    />
  );
}
