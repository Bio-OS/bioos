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

import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import classNames from 'classnames';
import { chunk } from 'lodash-es';
import Papa from 'papaparse';
import { v4 } from 'uuid';
import { Message, Modal, Notification, Upload } from '@arco-design/web-react';
import {
  IconCaretDown,
  IconCaretUp,
  IconPlus,
} from '@arco-design/web-react/icon';

import Icon from 'components/Icon';
import {
  entityRules,
  ParseResult,
  transformName,
  workspaceRules,
} from 'helpers/validate';
import Api from 'api/client';
import { HandlersDataModel } from 'api/index';

import styles from '../Category.less';
export interface Props {
  visible: boolean;
  type?: 'entity' | 'workspace';
  title?: string;
  itemTips?: string;
  csvTips?: (string | React.ReactNode)[];
  onClose?: () => void;
  onConfirm?: (name: string) => void;
  entityList?: HandlersDataModel[];
  startUpload?: (file) => void;
  completeUpload?: (file) => void;
}
export interface FileInfo {
  name: string;
  csvName: string;
  size: number;
  status: 'init' | 'success' | 'error';
  id: string;
}
const ImportEntityModal: React.FC<Props> = props => {
  const {
    visible,
    title,
    type,
    onClose,
    onConfirm,
    itemTips,
    csvTips,
    entityList,
    startUpload,
    completeUpload,
  } = props;
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const [display, setDisplay] = useState(true);
  const [disabled, setDisabled] = useState(true);
  const [data, setData] = useState({
    header: [],
    rows: [],
  });
  const [fileReader, setFileReader] = useState<FileInfo | null>({
    name: '',
    csvName: '',
    size: 0,
    status: 'init',
    id: '',
  });
  const handleUploadCsv = async () => {
    const isEntity = type === 'entity';
    const { header, rows } = data;
    const perChunk = Math.ceil((fileReader?.size || 0) / (4 * 1024 * 1024));
    const chunks = chunk(rows, Math.floor(rows.length / perChunk));
    const promiseArr = chunks.map(chunk => {
      return () =>
        Api.dataModelPartialUpdate(workspaceId, {
          headers: isEntity ? header : ['Key', 'Value'],
          name: fileReader.csvName,
          rows: chunk,
          workspaceID: workspaceId,
        });
    });
    onClose?.();
    setFileReader(null);
    setDisabled(true);
    try {
      const res = await promiseArr[0]();
      if (res.ok) {
        if (promiseArr.length > 1) {
          await Promise.all(
            promiseArr.slice(1, promiseArr.length).map(req => req()),
          );
        }
        Message.success(
          isEntity
            ? `上传数据模型${fileReader.csvName}成功`
            : '上传Workspace Data成功',
        );
        completeUpload({ ...fileReader, status: 'success' });
        onConfirm(isEntity ? fileReader.csvName : 'Workspace Data');
      }
    } catch (err) {
      Notification.error({
        title: isEntity
          ? `上传数据模型${fileReader.csvName}失败`
          : '上传Workspace Data失败',
        content: err?.error?.Message || err?.statusText || '',
      });
      completeUpload({ ...fileReader, status: 'error' });
    }
  };
  return (
    <Modal
      visible={visible}
      title={title}
      style={{ width: 580 }}
      onConfirm={() => {
        startUpload({ ...fileReader });
        handleUploadCsv();
      }}
      onCancel={() => {
        setFileReader(null);
        onClose?.();
        setDisabled(true);
      }}
      okText="导入数据表"
      okButtonProps={{
        disabled,
      }}
      focusLock={false}
      maskClosable={false}
      unmountOnExit
    >
      <div className="fs13">{itemTips}</div>
      <Upload
        className={styles.trigger}
        onDrop={e => {
          if (e.dataTransfer?.files?.[0]?.type !== 'text/csv') {
            Message.error('文件格式错误，仅支持 .csv文件');
          }
        }}
        drag={true}
        accept=".csv"
        beforeUpload={file => {
          return file.size >= 0;
        }}
        showUploadList={false}
        onChange={(_, file) => {
          if (file.status === 'init') {
            const readerBase64 = new FileReader();
            readerBase64.readAsText(file.originFile as Blob, 'gb2312');
            readerBase64.onload = () => {
              const result: ParseResult = Papa.parse(
                readerBase64.result as string,
                { delimiter: '' },
              );
              const header = result?.data?.shift();
              result.data = result?.data.filter(
                row => !row.every(item => item === ''),
              );
              const rows = result?.data || [];
              const csvName = transformName(header);
              setFileReader({
                name: file?.originFile?.name,
                csvName: type === 'entity' ? csvName : 'workspace_data',
                size: file.originFile.size,
                status: 'init',
                id: v4(),
              });
              const validateRules =
                type === 'entity' ? entityRules : workspaceRules;
              for (const rule of validateRules) {
                if (rule.validate(header, rows)) {
                  const message =
                    typeof rule?.message === 'function' &&
                    rule?.message(csvName);
                  setDisabled(true);
                  return Message.error({
                    content: message || (rule?.message as string),
                    duration: 10000,
                  });
                }
              }
              setData({
                header,
                rows,
              });
              setDisabled(false);
            };
          }
        }}
        tip={
          <div className="flexCenter flexCol">
            {fileReader?.name && !disabled ? (
              <Icon glyph="csv" size={40} />
            ) : (
              <IconPlus
                fontSize={20}
                style={{ display: 'block', marginBottom: 0 }}
              />
            )}
            <span className="colorBlack mt8 fw500">
              {fileReader?.name && !disabled
                ? fileReader?.name
                : '点击或拖拽选择 .csv 文件到此处导入'}
            </span>
            <span className="fs12">
              {fileReader?.name && !disabled
                ? type === 'entity'
                  ? `默认实体表名称：${fileReader?.name}`
                  : '导入数据表后将增量更新，若存在重复的key，则覆盖原数据'
                : '将文件拖拽至框内上传'}
            </span>
            {fileReader?.name && !fileReader?.csvName && (
              <span className="colorWarning">
                数据模型格式无效，请按指定格式制作
              </span>
            )}
            {type === 'entity' &&
              entityList?.some(item => item.name === fileReader?.csvName) && (
                <span>
                  当前 Workspace
                  已经存在同名数据表，继续导入可能会覆盖部分原数据
                </span>
              )}
          </div>
        }
      />
      <>
        <span className="colorBlack mr8 fs13">制作csv文件指南</span>
        {display ? (
          <IconCaretUp
            style={{ color: '#80838a' }}
            onClick={() => {
              setDisplay(!display);
            }}
          />
        ) : (
          <IconCaretDown
            style={{ color: '#80838a' }}
            onClick={() => {
              setDisplay(!display);
            }}
          />
        )}
      </>
      <div className={classNames(['colorBlack3 fs13 mt4', styles.uploadLi])}>
        {display &&
          csvTips?.map((item, index) => {
            return <li key={index}>{item}</li>;
          })}
      </div>
    </Modal>
  );
};
export default ImportEntityModal;
