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

import { lazy, Suspense, useEffect } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { Layout } from '@arco-design/web-react';

import Header from 'components/Header';
import { Z_INDEX } from 'helpers/constants';
import interceptFetch from 'api/intercept';

const WorkspaceList = lazy(() => import('./pages/workspace'));
const WorkspaceSubPage = lazy(() => import('./pages/WorkspaceSubPage'));
const Icons = lazy(() => import('./pages/Icons'));

const HIDE_HEADER_PATH = ['/icons'];

function App() {
  useEffect(() => {
    const unregister = interceptFetch();
    return () => unregister();
  }, []);

  const hideHeader = HIDE_HEADER_PATH.includes(window.location.pathname);
  return (
    <Layout>
      {!hideHeader && (
        <Layout.Header style={{ zIndex: Z_INDEX.header }}>
          <Header />
        </Layout.Header>
      )}
      <Layout.Content>
        <Suspense fallback={null}>
          <Switch>
            <Route path="/workspace" exact render={() => <WorkspaceList />} />
            <Route
              path="/workspace/:workspaceId"
              render={() => <WorkspaceSubPage />}
            />
            <Route path="/icons" render={() => <Icons />} />
            <Redirect to="/workspace" />
          </Switch>
        </Suspense>
      </Layout.Content>
    </Layout>
  );
}

export default App;
