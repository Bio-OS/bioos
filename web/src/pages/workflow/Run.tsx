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

import { ReactNode, useEffect, useMemo, useRef, useState } from 'react';
import { Form as FinalForm } from 'react-final-form';
import { useHistory, useRouteMatch } from 'react-router-dom';
import { FormApi } from 'final-form';
import { get, noop } from 'lodash-es';
import omitDeep from 'omit-deep';
import {
  Button,
  Link,
  Message,
  Popover,
  Select,
  Space,
  Switch,
  Tabs,
  Typography,
} from '@arco-design/web-react';
import {
  IconCaretRight,
  IconQuestionCircle,
} from '@arco-design/web-react/icon';

import BeforeUnloadModal from 'components/BeforeUnloadModal';
import Breadcrumbs from 'components/Breadcrumbs';
import DetailPage from 'components/DetailPage';
import PageEmpty from 'components/Empty';
import { Line } from 'components/index';
import Loading from 'components/Loading';
import SubTitle from 'components/SubTitle';
import TimerAlert from 'components/TimerAlert';
import AnalysisMode, {
  AnalysisModeType,
} from 'components/workflow/AnalysisMode';
import DAGToGraph from 'components/workflow/DAGToGraph';
import AnalyzeModal from 'components/workflow/run/AnalyzeModal';
import RunParams from 'components/workflow/run/RunParams';
import SelectDataModelModal from 'components/workflow/run/SelectDataModelModal';
import {
  getFormValue,
  getFormValueDefault,
  getValue,
  replaceDotToHyphen,
  stringKey,
} from 'components/workflow/run/utils';
import WDLFileViewer from 'components/workflow/WDLFileViewer';
import {
  INPUT_MODEL_DEFAULT_KEY,
  INPUT_PATH_MODEL_DEFAULT_KEY,
  OUTPUT_MODEL_DEFAULT_KEY,
  OUTPUT_PATH_MODEL_DEFAULT_KEY,
  WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY,
} from 'helpers/constants';
import { useQuery } from 'helpers/hooks';
import Api from 'api/client';
import { HandlersSubmissionItem } from 'api/index';

import style from './style.less';

export default function WorkflowRun() {
  const refForm = useRef<FormApi>();
  const dataModelHeaders = useRef([]);
  const runParamPathRef = useRef<{
    pathRowKeys: string[];
    checkedPathRowKeys: string[];
  }>({ pathRowKeys: [], checkedPathRowKeys: [] });

  const history = useHistory();
  const { submissionId } = useQuery();
  const match = useRouteMatch<{ workspaceId: string; workflowId: string }>();
  const { workspaceId, workflowId } = match.params;

  const [callCaching, setCallCaching] = useState(undefined);
  const [pathRowIds, setPathRowIds] = useState<string[]>([
    INPUT_PATH_MODEL_DEFAULT_KEY,
  ]);
  const [mode, setMode] = useState<AnalysisModeType>('dataModel');
  const [workflow, setWorkflow] = useState(null);
  const [submission, setSubmission] = useState<HandlersSubmissionItem>(null);
  const [dataModels, setDataModels] = useState(null);
  const [modelId, setModelId] = useState('');
  const [modelValidate, setModelValidate] = useState(null);
  const [options, setOptions] = useState<{
    dataModel: string[];
    workspaceModel: string[];
  }>({
    dataModel: [],
    workspaceModel: [],
  });
  const [selectDataModelRowKeys, setSelectDataModelRowKeys] = useState([]);
  const isPath = mode === 'filePath';

  const currentDataModel = useMemo(
    () => getDateModelInfo({ id: modelId }),
    [modelId],
  );

  const workspaceDataModel = useMemo(() => {
    return getDateModelInfo({ name: 'workspace_data' });
  }, [dataModels]);

  const dataModelList = useMemo(() => {
    return dataModels?.filter(_ => _.type !== 'workspace');
  }, [dataModels]);

  const outputSelectOptions = useMemo(() => {
    if (currentDataModel?.Type === 'set') {
      return {
        invalidHeaderArr: options.dataModel.slice(0, 2),
        validHeaderArr: options.dataModel.slice(2),
      };
    }

    return {
      invalidHeaderArr: options.dataModel.slice(0, 1),
      validHeaderArr: options.dataModel.slice(1),
    };
  }, [options.dataModel, currentDataModel]);

  function getDateModelInfo({ name, id }: { name?: string; id?: string }) {
    return dataModels?.find(_ => _.name === name || _.id === id);
  }

  async function getModelHeaders(prefix, modelInfo) {
    const res = await Api.dataModelDetail2(workspaceId, modelInfo.id);
    if (res.ok) {
      dataModelHeaders.current.push(
        ...res.data.headers.map(item => `${prefix}.${item}`),
      );

      if (modelInfo.type !== 'entity_set') {
        setOptions(pre => ({
          dataModel: dataModelHeaders.current,
          workspaceModel: pre.workspaceModel,
        }));
        return;
      }

      const nextModelInfo = getDateModelInfo({ name: res.data.headers[1] });
      if (!nextModelInfo) return;

      const nextPrefix = `${prefix}.${res.data.headers[1]}`;
      getModelHeaders(nextPrefix, nextModelInfo);
    }
  }

  function renderAlert() {
    if (!modelValidate) return;
    let content: ReactNode;

    if (modelValidate?.status === 'model-deleted') {
      content =
        '当前分析历史的工作流配置所对应的实体数据表、指定实体已被删除。';
    }

    if (modelValidate?.status === 'rows-deleted') {
      content = (
        <span>
          当前分析历史的工作流配置所对应的指定实体已被删除
          <span className="mr4 ml4">
            {modelValidate.originalRows.length - modelValidate.validRows.length}
          </span>
          项。
        </span>
      );
    }
    return (
      <TimerAlert content={content} onClose={() => setModelValidate(null)} />
    );
  }

  function handleChangeSelectRowKeys(keys: string[]) {
    setSelectDataModelRowKeys(keys);
  }

  async function validateDataModelRows() {
    const dataEntityId = submission.entity.dataModelID;
    const submissionRowIds = submission.entity.dataModelRowIDs;
    const res = await Api.dataModelRowsDetail(workspaceId, dataEntityId, {
      rowIDs: submissionRowIds,
    });
    const dataModelRowIds = res.data.rows?.map(_ => _[0]);
    if (submissionRowIds.length === dataModelRowIds.length) {
      setSelectDataModelRowKeys(submissionRowIds);
      return;
    }

    const realRowIds = submissionRowIds.filter(rowID =>
      dataModelRowIds?.includes(rowID),
    );

    setSelectDataModelRowKeys(realRowIds);

    setModelValidate({
      status: 'rows-deleted',
      validRows: realRowIds,
      originalRows: submissionRowIds,
    });
  }

  function handleSaveLocalStorage() {
    const { errors, values } = refForm.current.getState();
    let valuesData = values;
    if (errors) {
      // 过滤掉报错项
      const errKeys: string[] = ['dirty'];
      Object.keys(errors).forEach(item => {
        Object.keys(errors[item]).forEach(key => {
          Object.keys(errors[item][key]).forEach(subKey => {
            errKeys.push(`${key}.${subKey}`);
          });
        });
      });
      valuesData = omitDeep(values, errKeys);
    }

    const paramInfo = {
      workspaceId,
      workflowId,
      paramValues: valuesData,
      pathRowKeys: runParamPathRef.current?.pathRowKeys,
    };

    if (!window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY]) {
      window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY] =
        JSON.stringify([paramInfo]);
      return;
    }

    const paramInfoArr = JSON.parse(
      window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY],
    );

    const currentIndex = paramInfoArr.findIndex(
      item =>
        item.workspaceId === workspaceId && item.workflowId === workflowId,
    );

    if (currentIndex < 0) {
      paramInfoArr.push(paramInfo);
    } else {
      paramInfoArr[currentIndex] = paramInfo;
    }

    window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY] =
      JSON.stringify(paramInfoArr);
  }

  async function createSubmission(data) {
    const { values } = refForm.current.getState();
    const { pathRowKeys, checkedPathRowKeys } = runParamPathRef.current;

    const inputJsonObj: { [key: string]: any } = {};
    const outputJsonObj: { [key: string]: any } = {};

    if (isPath) {
      const rowKeys =
        pathRowKeys.length === 1 ? pathRowKeys : checkedPathRowKeys;
      // 只有一列属性 直接投递
      rowKeys.forEach(key => {
        const inputObj = get(values, `input.${stringKey(key)}`, {});
        inputJsonObj[key] = {};
        workflow.latestVersion.inputs.forEach(item => {
          const v = inputObj[replaceDotToHyphen(item.name)];
          if (v) {
            inputJsonObj[key][item.name] = getValue(v);
          }
        });
      });
    } else {
      const inputObj = get(values, `input.${INPUT_MODEL_DEFAULT_KEY}`, {});
      const outputObj = get(values, `output.${OUTPUT_MODEL_DEFAULT_KEY}`, {});

      workflow.latestVersion.inputs.forEach(item => {
        const v = inputObj[replaceDotToHyphen(item.name)];
        if (v) {
          inputJsonObj[item.name] = getValue(v);
        }
      });

      workflow.latestVersion.outputs.forEach(item => {
        const v = outputObj[replaceDotToHyphen(item.name)];
        if (v) {
          outputJsonObj[item.name] = getValue(v);
        }
      });
    }

    const body = {
      ...data,
      workspaceID: workspaceId,
      workflowID: workflowId,
      type: isPath ? 'filePath' : 'dataModel',
      exposedOptions: {
        readFromCache: callCaching,
      },
    };

    if (isPath) {
      body.inOutMaterial = {
        inputsMaterial: JSON.stringify(inputJsonObj),
        outputsMaterial: '',
      };
    } else {
      body.entity = {
        dataModelID: isPath ? undefined : modelId,
        dataModelRowIDs: isPath ? undefined : selectDataModelRowKeys,
        inputsTemplate: JSON.stringify(inputJsonObj),
        outputsTemplate: Object.keys(outputJsonObj).length
          ? JSON.stringify(outputJsonObj)
          : undefined,
      };
    }

    const res = await Api.submissionCreate(workspaceId, body);
    afterCreateSubmission();
    return res.data.id;
  }

  function afterCreateSubmission() {
    refForm.current.reset(refForm.current.getState().values);
    if (window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY]) {
      const { index, list } = getCurrentLocalStorage();

      if (index >= 0) {
        list.splice(index, 1);
        window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY] =
          JSON.stringify(list);
      }
    }
  }

  function handleOpenAnalyze(open) {
    if (mode === 'dataModel') {
      if (!modelId) {
        Message.error('请选择实体名称');
        return;
      }
      if (!selectDataModelRowKeys?.length) {
        Message.error('请选择实体数据');
        return;
      }
    } else if (
      // 路径分析 多列input 需要至少勾选一个
      runParamPathRef.current.pathRowKeys.length > 1 &&
      !runParamPathRef.current.checkedPathRowKeys.length
    ) {
      Message.error('请选择实体属性');
      return;
    }
    refForm.current.submit();
    const { errors } = refForm.current.getState();

    let hasErr: boolean | undefined = false;
    if (isPath && runParamPathRef.current.checkedPathRowKeys.length) {
      hasErr = Boolean(
        runParamPathRef.current.checkedPathRowKeys.some(
          key => errors?.input?.[stringKey(key)],
        ),
      );
    } else if (isPath) {
      hasErr = Boolean(
        Object.keys(
          errors?.input?.[runParamPathRef.current.pathRowKeys[0]] || {},
        ).length,
      );
    } else {
      hasErr = Boolean(
        Object.keys(errors?.input?.[INPUT_MODEL_DEFAULT_KEY] || {}).length,
      );
    }

    if (hasErr) {
      Message.error('请检查输入输出参数配置，填写必选参数');
      return;
    }
    open();
  }

  function getCurrentLocalStorage() {
    const paramInfoArr = JSON.parse(
      window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY],
    );

    const index = paramInfoArr.findIndex(
      item =>
        item.workspaceId === workspaceId && item.workflowId === workflowId,
    );

    return { index, current: paramInfoArr[index], list: paramInfoArr };
  }

  function initFromSubmission() {
    if (!submission) return;
    const usePath = submission.type === 'filePath';
    const inputData = usePath
      ? submission.inOutMaterial.inputsMaterial
      : submission.entity.inputsTemplate;
    const outputData = usePath
      ? submission.inOutMaterial.outputsMaterial
      : submission.entity.outputsTemplate;
    const submissionInput = JSON.parse(inputData);
    const submissionOutput = JSON.parse(outputData);
    const outputValue = getFormValue(submissionOutput);
    const input: { [key: string]: { [key: string]: string } } = {};
    const output: { [key: string]: { [key: string]: string } } = {};

    // 使用数据模型分析
    if (!usePath) {
      const inputValue = getFormValue(submissionInput);
      input[INPUT_MODEL_DEFAULT_KEY] = inputValue;
      input[INPUT_PATH_MODEL_DEFAULT_KEY] = {};
      output[OUTPUT_MODEL_DEFAULT_KEY] = outputValue;
    } else {
      const keys = Object.keys(submissionInput);
      setPathRowIds(keys);
      keys.forEach(key => {
        input[stringKey(key)] = getFormValue(submissionInput[key]);
      });
      input[INPUT_MODEL_DEFAULT_KEY] = {};
      output[OUTPUT_MODEL_DEFAULT_KEY] = {};
    }

    refForm.current?.reset({
      input,
      output,
    });
  }

  function initFromLocalStorage() {
    const { current } = getCurrentLocalStorage();
    if (!current) return initDefault();
    refForm.current?.reset(current.paramValues);

    setPathRowIds(current.pathRowKeys || [INPUT_PATH_MODEL_DEFAULT_KEY]);
  }

  function initDefault() {
    if (!workflow) return;
    const inputValue = getFormValueDefault(workflow.latestVersion.inputs);
    const outputValue = getFormValueDefault(workflow.latestVersion.outputs);
    const input: { [key: string]: { [key: string]: string } } = {};
    const output: { [key: string]: { [key: string]: string } } = {};

    input[INPUT_MODEL_DEFAULT_KEY] = inputValue;
    input[INPUT_PATH_MODEL_DEFAULT_KEY] = inputValue;
    output[OUTPUT_MODEL_DEFAULT_KEY] = outputValue;

    refForm.current?.reset({
      input,
      output,
    });
  }

  async function getDataModelRows(id: string) {
    const res = await Api.dataModelRowsDetail(workspaceId, id);
    if (res.ok) {
      const rows = res.data.rows.map(item => `workspace.${item[0]}`) || [];
      setOptions(pre => ({
        dataModel: pre.dataModel,
        workspaceModel: rows,
      }));
    }
  }

  async function getWorkflow() {
    const res = await Api.workspaceIdWorkflowDetail(workflowId, workspaceId);
    if (res.ok) {
      setWorkflow(res.data.workflow);
    }
  }

  async function getDataModels() {
    const res = await Api.dataModelDetail(workspaceId);
    if (res.ok) {
      setDataModels(res.data.Items);
    }
  }

  async function getSubmission() {
    const res = await Api.submissionDetail(workspaceId, {
      ids: [submissionId],
    });
    if (res.ok) {
      const data = res.data.items[0];
      setSubmission(data);
    }
  }

  // 获取初始数据
  useEffect(() => {
    getWorkflow();
    getDataModels();
    if (submissionId) {
      getSubmission();
    }
  }, []);

  useEffect(() => {
    if (!dataModelList) return;
    setCallCaching(submission?.exposedOptions?.readFromCache ?? true);
    if (!submissionId || submission?.type === 'filePath') {
      setModelId(dataModelList?.[0]?.id);
      setMode(submission?.type || 'dataModel');
      return;
    }

    if (submission?.entity) {
      const dataEntityId = submission.entity.dataModelID;
      if (dataEntityId) {
        const index = dataModelList?.findIndex(
          item => item.id === dataEntityId,
        );
        if (index > -1) {
          validateDataModelRows();
          setModelId(dataEntityId);
        } else {
          setModelValidate({ status: 'model-deleted' });
        }
      }
    }
  }, [submissionId, submission, dataModelList]);

  // 递归获取实体/集合表 headers 作为下拉菜单
  useEffect(() => {
    if (currentDataModel) {
      dataModelHeaders.current = [];
      getModelHeaders('this', currentDataModel);
    }
  }, [currentDataModel]);

  // 获取workspace data rows 作为下拉菜单
  useEffect(() => {
    if (workspaceDataModel) {
      getDataModelRows(workspaceDataModel.id);
    }
  }, [workspaceDataModel]);

  useEffect(() => {
    if (!workflow || (submissionId && !submission)) return;

    if (submissionId) {
      return initFromSubmission();
    }

    if (window.localStorage[WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY]) {
      return initFromLocalStorage();
    }

    initDefault();
  }, [workflow, submission]);

  if (!workflow) {
    return <Loading />;
  }

  return (
    <FinalForm
      onSubmit={noop}
      subscription={{ initialValues: true, pristine: true }}
    >
      {({ form, pristine }) => {
        refForm.current = form;
        return (
          <DetailPage
            showBorderBottom={false}
            rightArea={
              <Space>
                <Button
                  onClick={() => {
                    history.push(`/workspace/${workspaceId}/workflow`);
                  }}
                >
                  取消
                </Button>
                <AnalyzeModal
                  workflowName={workflow?.name}
                  onOk={createSubmission}
                >
                  {open => (
                    <Button
                      type="primary"
                      status="success"
                      onClick={() => handleOpenAnalyze(open)}
                    >
                      <IconCaretRight />
                      开始分析
                    </Button>
                  )}
                </AnalyzeModal>
              </Space>
            }
            breadcrumbs={
              <Breadcrumbs
                className="fs12 lh20"
                breadcrumbs={[
                  {
                    text: '工作流',
                    path: `/workspace/${match.params.workspaceId}/workflow`,
                  },
                  {
                    text: workflow?.name,
                  },
                ]}
              ></Breadcrumbs>
            }
            description={
              <>
                <span className="">描述：</span>
                <Typography.Paragraph
                  className="colorGrey mr20"
                  style={{ maxWidth: 100 }}
                  ellipsis={{
                    showTooltip: {
                      type: 'popover',
                    },
                  }}
                >
                  {workflow?.description}
                </Typography.Paragraph>
                <span>来源：</span>
                <Link>{workflow?.latestVersion?.metadata?.gitURL}</Link>
                <span>工作流类型：</span>
                {workflow?.latestVersion?.language}
              </>
            }
            title={workflow?.name}
          >
            <div className={style.workflowRunWrap}>
              <div className={style.blockBox}>
                <SubTitle title="运行选项" bg="#f6f8fa" className="mb20" />
                {renderAlert()}
                <div className={style.row}>
                  <div className={style.col}>
                    分析方式
                    <Popover content="切换分析方式，配置的内容将不做保留">
                      <IconQuestionCircle className="ml4" />
                    </Popover>
                  </div>
                  <div className={style.col}>
                    <AnalysisMode value={mode} onChange={setMode} />
                  </div>
                </div>
                {!isPath && (
                  <>
                    <div className={`${style.row} flexAlignCenter`}>
                      <div className={style.col}>实体名称</div>
                      <div className={style.col}>
                        <Select
                          value={modelId}
                          style={{ width: 492 }}
                          dropdownRender={
                            !dataModelList?.length
                              ? () => <PageEmpty />
                              : undefined
                          }
                          onChange={id => {
                            setModelId(id);
                            setSelectDataModelRowKeys([]);
                          }}
                        >
                          {dataModelList?.map(option => (
                            <Select.Option key={option.id} value={option.id}>
                              {option.name}
                            </Select.Option>
                          ))}
                        </Select>
                      </div>
                    </div>
                    <div className={style.row}>
                      <div className={style.col}>指定实体</div>
                      <div className={style.col}>
                        <div className="flexAlignCenter">
                          <SelectDataModelModal
                            dataModelInfo={currentDataModel}
                            selectRows={selectDataModelRowKeys}
                            onChange={handleChangeSelectRowKeys}
                          >
                            {open => <Link onClick={open}>选择数据</Link>}
                          </SelectDataModelModal>
                          <Line height={10} />
                          <span className="colorGrey noShrink">
                            当前已选
                            <span className="fw600 ml4 mr4">
                              {selectDataModelRowKeys.length || 0}
                            </span>
                            项
                          </span>
                        </div>
                      </div>
                    </div>
                  </>
                )}

                <div className={style.row}>
                  <div className={style.col}>
                    CallCaching
                    <Popover content="CallCaching会在之前运行的任务的缓存中搜索具有完全相同的命令和完全相同的输入的任务。 如果缓存命中，将使用前一个任务的结果而不是重新运行，从而节省时间和资源。">
                      <IconQuestionCircle className="ml4" />
                    </Popover>
                  </div>
                  <div className={style.col}>
                    <Switch
                      checkedText="开"
                      uncheckedText="关"
                      checked={callCaching}
                      onChange={setCallCaching}
                    />
                  </div>
                </div>
              </div>
            </div>
            <div className={style.blockBox}>
              <SubTitle title="运行参数" bg="#f6f8fa" className="mb20" />
              <Tabs type="card-gutter" defaultActiveTab="input">
                <Tabs.TabPane key="wdl" title="描述文件">
                  <WDLFileViewer
                    workspaceId={workspaceId}
                    workflowId={workflowId}
                    files={workflow?.latestVersion?.files}
                    initialFile={workflow?.latestVersion?.mainWorkflowPath}
                  />
                </Tabs.TabPane>
                <Tabs.TabPane key="input" title="输入参数">
                  <RunParams
                    type="input"
                    isPath={false}
                    isHidden={isPath}
                    data={workflow?.latestVersion?.inputs}
                    rowIds={[INPUT_MODEL_DEFAULT_KEY]}
                    selectOptions={options.dataModel.concat(
                      options.workspaceModel,
                    )}
                  />
                  <RunParams
                    type="input"
                    isPath={true}
                    isHidden={!isPath}
                    data={workflow?.latestVersion?.inputs}
                    ref={runParamPathRef}
                    rowIds={pathRowIds}
                    selectOptions={options.workspaceModel}
                  />
                </Tabs.TabPane>
                <Tabs.TabPane key="ouput" title="输出参数">
                  <RunParams
                    type="output"
                    isPath={false}
                    isHidden={isPath}
                    data={workflow?.latestVersion?.outputs}
                    rowIds={[OUTPUT_MODEL_DEFAULT_KEY]}
                    selectOptions={outputSelectOptions.validHeaderArr}
                    invalidHeaderArr={outputSelectOptions.invalidHeaderArr}
                  />
                  <RunParams
                    type="output"
                    isPath={true}
                    isHidden={!isPath}
                    data={workflow?.latestVersion?.outputs}
                    rowIds={[OUTPUT_PATH_MODEL_DEFAULT_KEY]}
                  />
                </Tabs.TabPane>
                <Tabs.TabPane key="graph" title="Graph" style={{ height: 700 }}>
                  <DAGToGraph data={workflow?.latestVersion?.graph} />
                </Tabs.TabPane>
              </Tabs>
            </div>
            {!pristine && (
              <BeforeUnloadModal
                title="退出是否保存当前运行参数"
                content={
                  <span>
                    建议保存当前配置,便于记录当前工作流运行参数的详细信息。
                  </span>
                }
                onSave={handleSaveLocalStorage}
              />
            )}
          </DetailPage>
        );
      }}
    </FinalForm>
  );
}
