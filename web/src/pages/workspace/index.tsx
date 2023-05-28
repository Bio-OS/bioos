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

import { Tabs } from '@arco-design/web-react';

import List from '../../components/workspace/list';

import styles from './index.less';

const TabPane = Tabs.TabPane;
const tab = [
  {
    title: 'My workspace',
    key: 'self',
    renderChildren: () => {
      return <List />;
    },
  },
  // {
  //   title: 'Public',
  //   key: 'public',
  // },
];
export default function Index() {
  return (
    <div className={styles.workspaceList}>
      <div className={styles.workspaceTitle}>Workspace</div>
      <Tabs type="card-gutter" size="large">
        {tab.map(tab => {
          return (
            <TabPane key={tab.key} title={tab.title}>
              {tab.renderChildren?.()}
            </TabPane>
          );
        })}
      </Tabs>
    </div>
  );
}
