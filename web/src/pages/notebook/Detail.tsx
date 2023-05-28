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
import { useParams, useRouteMatch } from 'react-router-dom';

import Breadcrumbs from 'components/Breadcrumbs';
import DetailPage from 'components/DetailPage';
import ConfigureRuntime from 'components/notebook/ConfigureRuntime';
import EditNotebookButton from 'components/notebook/EditNotebookButton';
import NotebookViewer from 'components/notebook/notebook-viewer';
import Api from 'api/client';

export default function NotebookDetail() {
  const { workspaceId, notebookName } = useParams<{
    workspaceId: string;
    notebookName: string;
  }>();

  const match = useRouteMatch<{ workspaceId: string }>();

  const [preview, setPreview] = useState<any>();
  useEffect(() => {
    Api.workspaceIdNotebookDetail(workspaceId, notebookName, {
      format: 'json',
    }).then(({ data }) => {
      setPreview(data);
    });
  }, []);
  return (
    <DetailPage
      breadcrumbs={
        <Breadcrumbs
          className="fs12 lh20"
          breadcrumbs={[
            {
              text: 'Notebook',
              path: `/workspace/${match.params.workspaceId}/notebook`,
            },
            {
              text: 'Notebook详情',
            },
          ]}
        ></Breadcrumbs>
      }
      title={notebookName}
      rightArea={
        <ConfigureRuntime
          render={({ setVisibleConfigureRuntime }) => {
            return (
              <EditNotebookButton
                usage="notebook"
                notebookName={notebookName}
                setVisibleConfigureRuntime={setVisibleConfigureRuntime}
              />
            );
          }}
        />
      }
    >
      <NotebookViewer notebook={preview} loading={preview === undefined} />
    </DetailPage>
  );
}
