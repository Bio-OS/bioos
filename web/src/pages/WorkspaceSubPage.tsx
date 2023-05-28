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

import { lazy, memo, Suspense, useEffect, useRef, useState } from 'react';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router-dom';
import classNames from 'classnames';
import { Layout, Menu, Typography } from '@arco-design/web-react';
import { IconLeft } from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import Api from 'api/client';

import styles from './index.less';

const WorkflowList = lazy(() => import('./workflow'));
const WorkflowRun = lazy(() => import('./workflow/Run'));
const Data = lazy(() => import('./data'));
const Notebook = lazy(() => import('./notebook'));
const NotebooksDetail = lazy(() => import('./notebook/Detail'));
const NotebookEdit = lazy(() => import('./notebook/Edit'));

const AnalysisList = lazy(() => import('./analysis'));
const AnalysisDetail = lazy(() => import('./analysis/AnalysisDetail'));
const AnalysisTaskDetail = lazy(() => import('./analysis/AnalysisTaskDetail'));
const MenuItem = Menu.Item;

const MENUS = [
  {
    key: 'data',
    display: '数据',
  },
  {
    key: 'notebook',
    display: 'Notebooks',
  },
  {
    key: 'workflow',
    display: '工作流',
  },
  {
    key: 'analysis',
    display: '分析历史',
  },
];

export default function WorkspaceSubPage() {
  const match = useRouteMatch<{ workspaceId: string }>();
  const dataRef = useRef();
  const hideSide = /\/(add|edit|bind)$/.test(location.pathname);

  return (
    <Layout id="rootContainer" className="flexRow">
      {!hideSide && <SideMenu />}
      {/** 设置width: 0 使右侧宽度完全由flex: 1这个属性来分配 */}
      <Layout.Content style={{ width: 0 }}>
        <Suspense fallback={null}>
          <Switch>
            <Route
              exact
              path={`${match.path}/data`}
              render={() => <Data ref={dataRef} />}
            />
            <Route
              exact
              path={`${match.path}/workflow`}
              render={() => <WorkflowList />}
            />

            <Route
              exact
              path={`${match.path}/notebook`}
              render={() => <Notebook />}
            />
            <Route
              path={`${match.path}/notebook/:notebookName`}
              exact
              component={NotebooksDetail}
            />
            <Route
              path={`${match.path}/notebook/:notebookName/edit`}
              component={NotebookEdit}
            />
            <Route
              exact
              path={`${match.path}/workflow/:workflowId/run`}
              render={() => <WorkflowRun />}
            />
            <Route
              exact
              path={`${match.path}/analysis`}
              render={() => <AnalysisList />}
            />
            <Route
              path={`${match.path}/analysis/detail/taskDetail`}
              component={AnalysisTaskDetail}
            />
            <Route
              path={`${match.path}/analysis/detail`}
              component={AnalysisDetail}
            />
          </Switch>
        </Suspense>
      </Layout.Content>
    </Layout>
  );
}

const SideMenu = memo(({}) => {
  const history = useHistory();
  const match = useRouteMatch<{ workspaceId: string }>();
  const { workspaceId } = match.params;
  const checkedMenu = window.location.pathname
    .replace(match.url, '')
    .split('/')?.[1];
  const [workspaceInfo, setWorkspaceInfo] = useState(null);

  function handleChangeMenu(key: string) {
    history.push(`/workspace/${workspaceId}/${key}`);
  }

  function handleBack() {
    history.push('/workspace');
  }

  async function getWorkspaceInfo() {
    const res = await Api.workspaceDetail(workspaceId);
    setWorkspaceInfo(res.data);
  }

  useEffect(() => {
    getWorkspaceInfo();
  }, [workspaceId]);

  return (
    <Menu
      className={styles.sideMenu}
      style={{ width: 200 }}
      hasCollapseButton
      icons={{
        collapseDefault: <Icon glyph="menu" />,
        collapseActive: <Icon glyph="menu" className="rotate180" />,
      }}
      selectedKeys={[checkedMenu]}
      onClickMenuItem={handleChangeMenu}
    >
      <div
        className={classNames(['flexAlignCenter pt20 pb8', styles.menuTitle])}
      >
        <div
          className="flexCenter cursorPointer noShrink mr8 icon"
          onClick={handleBack}
        >
          <IconLeft />
        </div>
        <Typography.Text
          className="fs16 fw500"
          ellipsis={{ cssEllipsis: true }}
        >
          {workspaceInfo?.name}
        </Typography.Text>
      </div>
      {MENUS.map(item => {
        return <MenuItem key={item.key}>{item.display}</MenuItem>;
      })}
    </Menu>
  );
});
