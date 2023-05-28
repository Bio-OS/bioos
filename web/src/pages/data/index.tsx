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

import { forwardRef, useEffect, useMemo, useRef, useState } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import { v4 } from 'uuid';
import { Layout, Link, Message, Popover } from '@arco-design/web-react';
import { IconDownload, IconPlus } from '@arco-design/web-react/icon';

import Category from 'components/data/Category';
import ImportEntityModal, {
  FileInfo,
  Props as CategoryInfo,
} from 'components/data/modal/ImportDataModal';
import UploadList from 'components/data/modal/UploadList';
import ModelContent from 'components/data/ModelContent';
import Icon from 'components/Icon';
import ListPage from 'components/ListPage';
import {
  ENTITY_CSV,
  ENTITY_ITEM,
  WORKSPACE_CSV,
  WORKSPACE_ITEM,
} from 'helpers/constants';
import { useQuery, useQueryHistory } from 'helpers/hooks';
import { downloadFile } from 'helpers/utils';
import Api, { apiInstance } from 'api/client';
import { HandlersDataModel } from 'api/index';

import style from './index.less';
const { Sider, Content } = Layout;
interface CategoryList {
  entity: HandlersDataModel[];
  workspace: HandlersDataModel[];
}
const Data = (_, ref) => {
  const query = useQuery();
  const { pathname } = useLocation();
  const navigate = useQueryHistory();
  const tokenModelRef = useRef(null);
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [importInfo, setImportInfo] = useState<CategoryInfo>({
    visible: false,
    type: 'entity',
    title: '',
  });
  const [uploadFileList, setUploadFileList] = useState<FileInfo[]>([]);
  const [categoryList, setCategoryList] = useState<CategoryList>({
    entity: [],
    workspace: [],
  });

  const [activeItem, setActiveItem] = useState<null | HandlersDataModel>();
  const modelRef = useRef(null);
  const fetchDataModels = (name?: string) => {
    const entityList = [];
    const workspaceList = [];
    tokenModelRef.current = v4();
    Api.dataModelDetail(
      workspaceId,
      {},
      {
        cancelToken: tokenModelRef.current,
      },
    )
      .then(res => {
        if (res?.ok) {
          (res?.data?.Items || []).map(item => {
            if (item.type === 'workspace') {
              workspaceList.push({ ...item, name: 'Workspace Data' });
            } else {
              entityList.push(item);
            }
          });
          setCategoryList({
            entity: entityList,
            workspace: workspaceList.length
              ? workspaceList
              : [
                  {
                    id: '-1',
                    name: 'Workspace Data',
                    rowCount: 0,
                    type: 'workspace',
                  },
                ],
          });
          // name 上传 query 刷新 ref 路由
          const item = [...entityList, ...workspaceList].find(item =>
            name
              ? item.name === name
              : item.id === query?.modelId || item.id === ref.current?.id,
          );
          const activeItem = item || entityList[0] || workspaceList[0];
          name && modelRef?.current?.updateModel();
          setActiveItem(activeItem);
          navigate(pathname, {
            ...query,
            modelId: activeItem?.id,
            search: name ? '' : query.search,
            size: 10,
            page: 1,
          });
          ref.current = activeItem;
        }
      })
      .catch(e => {
        const info = e?.error?.message || e?.statusText;
        info && Message.error(info);
      })
      .finally(() => {
        tokenModelRef.current = null;
      });
  };

  useEffect(() => {
    fetchDataModels();
    return () => {
      if (tokenModelRef.current) {
        apiInstance.abortRequest(tokenModelRef.current);
      }
    };
  }, []);
  const includeSet = useMemo(() => {
    return categoryList?.entity?.find(
      item => item.name === `${activeItem?.name}_set`,
    );
  }, [activeItem, categoryList]);

  const LinkToDownLoad = (props: { name: string; link: string }) => {
    const { name, link } = props;
    return (
      <>
        <span>模板参考：</span>
        <Link
          className="colorPrimary"
          onClick={() => {
            downloadFile(link);
          }}
        >
          {name} <IconDownload />
        </Link>
      </>
    );
  };

  const updateImportInfo = (type: 'entity' | 'workspace') => {
    setImportInfo({
      visible: true,
      type,
      title: type === 'entity' ? '导入实体表' : '导入 Workspace Data',
      itemTips: type === 'entity' ? ENTITY_ITEM : WORKSPACE_ITEM,
      csvTips: [
        ...(type === 'entity' ? ENTITY_CSV : WORKSPACE_CSV),
        <LinkToDownLoad
          name={type === 'entity' ? '下载 csv 文件模板' : 'workspace_data.csv'}
          link={`https://bioos.tos-cn-beijing.volces.com/${
            type === 'entity' ? 'sample.csv' : 'workspace_data.csv'
          }`}
        />,
      ],
    });
  };

  const completeUpload = (file: FileInfo) => {
    setUploadFileList(preState => {
      return preState
        .map(item => {
          if (item.csvName === file.csvName && item.id === file.id) {
            item.status = file.status;
          }
          return item;
        })
        .filter(item => item.status !== 'success');
    });
  };
  const handleSelect = val => {
    setActiveItem(val);
    ref.current = val;
    modelRef?.current?.updateModel();
    navigate(pathname, {
      ...query,
      modelId: val.id,
      page: 1,
      size: 10,
      search: '',
    });
  };
  const { entity, workspace } = categoryList;
  return (
    <ListPage title="数据" showTitleBorder>
      <Layout className={style.dataModel}>
        <Sider width={244}>
          <Category
            title="实体数据模型"
            list={entity}
            activeItem={activeItem?.id}
            onSelect={handleSelect}
            suffix={
              <>
                <Popover
                  content={
                    <div className="colorBlack fw500">按 名称 正序排列</div>
                  }
                  position="bottom"
                >
                  <div style={{ display: 'flex' }}>
                    <Icon glyph="datasort" size={22} className="colorPrimary" />
                  </div>
                </Popover>
                <IconPlus
                  style={{ marginLeft: 12 }}
                  fontSize={20}
                  onClick={() => {
                    updateImportInfo('entity');
                  }}
                />
              </>
            }
          />
          <Category
            title="Workspace 数据模型"
            list={workspace}
            activeItem={activeItem?.id}
            onSelect={handleSelect}
            suffix={
              <IconPlus
                fontSize={20}
                onClick={() => {
                  updateImportInfo('workspace');
                }}
              />
            }
          />
        </Sider>
        <Content style={{ width: '976px', overflow: 'hidden' }}>
          <ModelContent
            model={activeItem}
            includeSet={includeSet}
            refresh={fetchDataModels}
            entityList={entity}
            ref={modelRef}
          />
        </Content>
      </Layout>
      <ImportEntityModal
        {...importInfo}
        entityList={entity}
        startUpload={file => {
          setUploadFileList([...uploadFileList, { ...file }]);
        }}
        completeUpload={completeUpload}
        onConfirm={fetchDataModels}
        onClose={() => {
          setImportInfo({ visible: false });
        }}
      />
      <UploadList
        list={uploadFileList}
        onClose={() => {
          setUploadFileList(preState => {
            return preState.filter(item => item.status === 'init');
          });
        }}
      />
    </ListPage>
  );
};
export default forwardRef(Data);
