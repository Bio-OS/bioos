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

export const Z_INDEX = {
  /** 100 */
  header: 100,
  /** 10 */
  modal: 99,
};
export const ENTITY_ITEM =
  '请导入 .csv 文件，实体表名称平台默认把它指定为 your-entity-name';
export const ENTITY_CSV = [
  '文件内需要至少包含1个实体行。',
  'ID 列标题为 your-entity-name_id，其中 your-entity-name 不可超过30个字符。',
  '表头仅支持字母、数字、连接符(-、_)，且不能以连接符开头，不可超过100个字符。',
  '表头名称不允许重复',
];
export const WORKSPACE_ITEM = '仅支持导入 Workspace 级别的 .csv 文件';
export const WORKSPACE_CSV = [
  '文件内需要至少包含 1 行 Workspace Data。',
  '属性列第一列为 Key 值，第二列为 Value 值，Value 一般为对应数据文件的路径信息。',
];

export const ANALYSIS_STATUS = [
  {
    text: '分析中',
    value: 'Running',
    color: 'arcoblue',
    icon: 'IconSunFill',
  },
  {
    text: '分析完成',
    value: 'Finished',
    color: 'green',
    icon: 'IconCheckCircleFill',
  },
  {
    text: '终止中',
    value: 'Cancelling',
    color: 'arcoblue',
    icon: 'IconSunFill',
  },
  {
    text: '已终止',
    value: 'Cancelled',
    color: 'gray',
    icon: 'IconMinusCircleFill',
  },
];

/** 数据模型输入默认属性列 */
export const INPUT_MODEL_DEFAULT_KEY = 'Default_DateModel_Input';
/** 路径分析输入默认属性列 */
export const INPUT_PATH_MODEL_DEFAULT_KEY = 'Default';
/** 数据模型输出默认属性列 */
export const OUTPUT_MODEL_DEFAULT_KEY = 'Default_DateModel_Output';
/** 路径分析输出默认属性列 */
export const OUTPUT_PATH_MODEL_DEFAULT_KEY = 'Default_Output';
/** 运行参数localstorage key */
export const WORKFLOW_RUN_PARAMS_LOCALSTORAGE_KEY = 'workflowRunParams';
export const WORKFLOW_STATUS_TO_NAME = {
  Succeeded: '工作流已运行成功',
  Failed: '工作流已运行失败',
  Cancelling: '工作流终止中',
  Cancelled: '工作流已终止',
};
export const RUN_STATUS_TAG = [
  {
    text: '启动中',
    value: 'Pending',
    color: 'arcoblue',
  },
  {
    text: '运行中',
    value: 'Running',
    color: 'arcoblue',
  },
  {
    text: '运行成功',
    value: 'Succeeded',
    color: 'green',
  },
  {
    text: '运行失败',
    value: 'Failed',
    color: 'red',
  },
  {
    text: '终止中',
    value: 'Cancelling',
    color: 'arcoblue',
  },
  {
    text: '已终止',
    value: 'Cancelled',
    color: 'orangered',
  },
];

export const GLOBAL_CONFIG_STORAGE_KEY = 'global-config-storage-key';
