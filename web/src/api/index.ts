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

/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ApiserverClientConfig {
  notebook?: {
    officialImages?: NotebookImage[];
    resourceOptions?: NotebookResourceSize[];
  };
  storage?: {
    fsPath?: string[];
  };
}

export interface ErrorsAppError {
  code?: number;
  message?: string;
}

export interface GithubComBioOSBioosInternalContextSubmissionInterfaceHertzHandlersWorkflowVersion {
  id?: string;
  versionID?: string;
}

export interface GithubComBioOSBioosInternalContextWorkspaceInterfaceHertzHandlersWorkflowVersion {
  createdAt?: string;
  files?: HandlersWorkflowFileInfo[];
  graph?: string;
  id?: string;
  inputs?: HandlersWorkflowParam[];
  language?: string;
  languageVersion?: string;
  mainWorkflowPath?: string;
  message?: string;
  metadata?: Record<string, string>;
  outputs?: HandlersWorkflowParam[];
  source?: string;
  status?: string;
  updatedAt?: string;
}

export interface HandlersCreateSubmissionRequest {
  description?: string;
  entity?: HandlersEntity;
  exposedOptions?: HandlersExposedOptions;
  inOutMaterial?: HandlersInOutMaterial;
  name?: string;
  type?: string;
  workflowID?: string;
  workspaceID?: string;
}

export interface HandlersCreateSubmissionResponse {
  id?: string;
}

export interface HandlersCreateWorkspaceRequest {
  description?: string;
  name?: string;
  storage?: HandlersWorkspaceStorage;
}

export interface HandlersCreateWorkspaceResponse {
  id?: string;
}

export interface HandlersDataModel {
  id?: string;
  name?: string;
  rowCount?: number;
  type?: string;
}

export interface HandlersEntity {
  dataModelID?: string;
  dataModelRowIDs?: string[];
  /**
   * * 输入配置，json 序列化后的 string
   * 	  采用 json 序列化原因基于以下两点考虑：
   * 	  - thrift/接口设计层面不允许 `Value` 类型不确定
   * 	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
   */
  inputsTemplate?: string;
  /**
   * * 输出配置，json 序列化后的 string
   * 	  采用 json 序列化原因基于以下两点考虑：
   * 	  - thrift/接口设计层面不允许 `Value` 类型不确定
   * 	  - 在 inputs/outputs 层级进行序列化可使得 `bioos-server` 不处理 `Inputs`/`Outputs`(非 `this.xxx` 索引的输入) 就入库/提交给计算引擎，达到透传效果
   */
  outputsTemplate?: string;
}

export interface HandlersExposedOptions {
  readFromCache?: boolean;
}

export interface HandlersGetDataModelResponse {
  dataModel?: HandlersDataModel;
  headers?: string[];
}

export interface HandlersGetWorkspaceByIdResponse {
  createTime?: number;
  description?: string;
  id?: string;
  name?: string;
  storage?: HandlersWorkspaceStorage;
  updateTime?: number;
}

export interface HandlersImportWorkspaceResponse {
  id?: string;
}

export interface HandlersInOutMaterial {
  inputsMaterial?: string;
  outputsMaterial?: string;
}

export interface HandlersListAllDataModelRowIDsResponse {
  rowIDs?: string[];
}

export interface HandlersListDataModelRowsResponse {
  headers?: string[];
  page?: number;
  rows?: string[][];
  size?: number;
  total?: number;
}

export interface HandlersListDataModelsResponse {
  Items?: HandlersDataModel[];
}

export interface HandlersListRunsResponse {
  items?: HandlersRunItem[];
  page?: number;
  size?: number;
  total?: number;
}

export interface HandlersListSubmissionsResponse {
  items?: HandlersSubmissionItem[];
  page?: number;
  size?: number;
  total?: number;
}

export interface HandlersListTasksResponse {
  items?: HandlersTaskItem[];
  page?: number;
  size?: number;
  total?: number;
}

export interface HandlersListWorkspacesResponse {
  items?: HandlersWorkspaceItem[];
  page?: number;
  size?: number;
  total?: number;
}

export interface HandlersNFSWorkspaceStorage {
  mountPath?: string;
}

export interface HandlersPatchDataModelRequest {
  async?: boolean;
  headers?: string[];
  name?: string;
  rows?: string[][];
  workspaceID?: string;
}

export interface HandlersPatchDataModelResponse {
  id?: string;
}

export interface HandlersRunItem {
  duration?: number;
  engineRunID?: string;
  finishTime?: number;
  id?: string;
  inputs?: string;
  log?: string;
  message?: string;
  name?: string;
  outputs?: string;
  startTime?: number;
  status?: string;
  taskStatus?: HandlersStatus;
}

export interface HandlersStatus {
  cancelled?: number;
  cancelling?: number;
  count?: number;
  failed?: number;
  initializing?: number;
  pending?: number;
  queued?: number;
  running?: number;
  succeeded?: number;
}

export interface HandlersSubmissionItem {
  description?: string;
  duration?: number;
  entity?: HandlersEntity;
  exposedOptions?: HandlersExposedOptions;
  finishTime?: number;
  id?: string;
  inOutMaterial?: HandlersInOutMaterial;
  name?: string;
  runStatus?: HandlersStatus;
  startTime?: number;
  status?: string;
  type?: string;
  workflowVersion?: GithubComBioOSBioosInternalContextSubmissionInterfaceHertzHandlersWorkflowVersion;
}

export interface HandlersTaskItem {
  duration?: number;
  finishTime?: number;
  name?: string;
  runID?: string;
  startTime?: number;
  status?: string;
  stderr?: string;
  stdout?: string;
}

export interface HandlersUpdateWorkspaceRequest {
  description?: string;
  id?: string;
  name?: string;
}

export interface HandlersWorkflowFile {
  content?: string;
  createdAt?: string;
  id?: string;
  path?: string;
  updatedAt?: string;
  workflowVersionID?: string;
}

export interface HandlersWorkflowFileInfo {
  id?: string;
  path?: string;
}

export interface HandlersWorkflowItem {
  createdAt?: string;
  description?: string;
  id?: string;
  latestVersion?: GithubComBioOSBioosInternalContextWorkspaceInterfaceHertzHandlersWorkflowVersion;
  name?: string;
  updatedAt?: string;
}

export interface HandlersWorkflowParam {
  default?: string;
  name?: string;
  optional?: boolean;
  type?: string;
}

export interface HandlersWorkspaceItem {
  createTime?: number;
  description?: string;
  id?: string;
  name?: string;
  storage?: HandlersWorkspaceStorage;
  updateTime?: number;
}

export interface HandlersWorkspaceStorage {
  nfs?: HandlersNFSWorkspaceStorage;
}

export interface HandlersCreateWorkflowRequest {
  description?: string;
  id?: string;
  language: string;
  mainWorkflowPath: string;
  name: string;
  source: 'git';
  tag: string;
  token?: string;
  url: string;
  workspaceID?: string;
}

export interface HandlersCreateWorkflowResponse {
  id?: string;
}

export interface HandlersGetWorkflowFileResponse {
  file?: HandlersWorkflowFile;
}

export interface HandlersGetWorkflowResponse {
  workflow?: HandlersWorkflowItem;
}

export interface HandlersGetWorkflowVersionResponse {
  version?: GithubComBioOSBioosInternalContextWorkspaceInterfaceHertzHandlersWorkflowVersion;
}

export interface HandlersListNotebooksResponse {
  items?: HandlersNotebookItem[];
}

export interface HandlersListWorkflowFilesRequest {
  ids?: string;
  orderBy?: string;
  page?: number;
  size?: number;
  workflowID?: string;
  workflowVersionID?: string;
  workspaceID?: string;
}

export interface HandlersListWorkflowFilesResponse {
  items?: HandlersWorkflowFile[];
  page?: number;
  size?: number;
  total?: number;
  workflowID?: string;
  workspaceID?: string;
}

export interface HandlersListWorkflowVersionsResponse {
  items?: GithubComBioOSBioosInternalContextWorkspaceInterfaceHertzHandlersWorkflowVersion[];
  page?: number;
  size?: number;
  total?: number;
  workflowID?: string;
  workspaceID?: string;
}

export interface HandlersListWorkflowsResponse {
  items?: HandlersWorkflowItem[];
  page?: number;
  size?: number;
  total?: number;
}

export interface HandlersNotebookItem {
  contentLength?: number;
  name?: string;
  updateTime?: number;
}

export interface HandlersUpdateWorkflowRequest {
  description?: string;
  id?: string;
  language: 'WDL';
  mainWorkflowPath?: string;
  name?: string;
  source: 'git';
  tag?: string;
  token?: string;
  url?: string;
  workspaceID?: string;
}

export interface HertzCreateRequest {
  image?: string;
  resourceSize?: NotebookResourceSize;
}

export interface HertzCreateResponse {
  id?: string;
}

export interface HertzGetResponse {
  accessURL?: string;
  createTime?: number;
  id?: string;
  image?: string;
  resourceSize?: NotebookResourceSize;
  status?: string;
  updateTime?: number;
}

export interface HertzListResponseItem {
  createTime?: number;
  id?: string;
  image?: string;
  resourceSize?: NotebookResourceSize;
  status?: string;
  updateTime?: number;
}

export interface HertzUpdateSettingsRequest {
  image?: string;
  resourceSize?: NotebookResourceSize;
}

export interface NotebookGPU {
  /** float is for mgpu */
  card?: number;
  memory?: number;
  model?: string;
}

export interface NotebookIPythonNotebook {
  cells?: NotebookIPythonNotebookCell[];
  metadata?: NotebookIPythonNotebookMeta;
  /** @min 4 */
  nbformat: number;
  nbformat_minor?: number;
}

export interface NotebookIPythonNotebookCell {
  attachments?: Record<string, any>;
  cell_type: 'code' | 'markdown' | 'raw';
  execution_count?: number;
  id?: string;
  metadata?: Record<string, any>;
  outputs?: NotebookIPythonNotebookCellOutput[];
  source?: string[];
}

export interface NotebookIPythonNotebookCellOutput {
  data?: Record<string, any>;
  /** In errors */
  ename?: string;
  evalue?: string;
  /** in execute result; used to be pyout / prompt_number */
  execute_count?: number;
  /** in display data */
  metadata?: Record<string, any>;
  name: string;
  output_type: 'stream' | 'display_data' | 'execute_result' | 'error';
  text?: string[];
  traceback?: string[];
}

export interface NotebookIPythonNotebookKernelSpec {
  display_name?: string;
  language: 'python' | 'R';
  name: string;
}

export interface NotebookIPythonNotebookLanguageInfo {
  /**
   * codemirror_mode in docs is string but jupyter created is {"name":"ipython",...}
   * CodeMirrorMode string `json:"codemirror_mode"`
   */
  file_extension: string;
  mimetype: string;
  name: 'python' | 'R';
  pygments_lexer?: string;
  version: string;
}

export interface NotebookIPythonNotebookMeta {
  kernelspec?: NotebookIPythonNotebookKernelSpec;
  language_info?: NotebookIPythonNotebookLanguageInfo;
  max_cell_id?: number;
}

export interface NotebookImage {
  description?: string;
  image?: string;
  name?: string;
  updateTime?: string;
  version?: string;
}

export interface NotebookResourceSize {
  cpu?: number;
  disk?: number;
  gpu?: NotebookGPU;
  memory?: number;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, 'body' | 'bodyUsed'>;

export interface FullRequestParams extends Omit<RequestInit, 'body'> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  'body' | 'method' | 'query' | 'path'
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, 'baseUrl' | 'cancelToken' | 'signal'>;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = 'application/json',
  FormData = 'multipart/form-data',
  UrlEncoded = 'application/x-www-form-urlencoded',
  Text = 'text/plain',
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = '/';
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>['securityWorker'];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: 'same-origin',
    headers: {},
    redirect: 'follow',
    referrerPolicy: 'no-referrer',
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(
      typeof value === 'number' ? value : `${value}`,
    )}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join('&');
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      key => 'undefined' !== typeof query[key],
    );
    return keys
      .map(key =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join('&');
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : '';
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === 'object' || typeof input === 'string')
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== 'string'
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === 'object' && property !== null
            ? JSON.stringify(property)
            : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === 'boolean' ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ''}${path}${
        queryString ? `?${queryString}` : ''
      }`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { 'Content-Type': type }
            : {}),
        },
        signal: cancelToken
          ? this.createAbortSignal(cancelToken)
          : requestParams.signal,
        body:
          typeof body === 'undefined' || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async response => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then(data => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch(e => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title BioOS Apiserver
 * @version 1.0
 * @license Apache 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html)
 * @baseUrl /
 * @contact hertz-contrib (https://github.com/hertz-contrib)
 *
 * This is bioos apiserver using Hertz.
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  wellKnown = {
    /**
     * @description get client configuration
     *
     * @name ConfigurationList
     * @summary use to get client configuration
     * @request GET:/.well-known/configuration
     */
    configurationList: (params: RequestParams = {}) =>
      this.request<ApiserverClientConfig, ErrorsAppError>({
        path: `/.well-known/configuration`,
        method: 'GET',
        format: 'json',
        ...params,
      }),
  };
  ping = {
    /**
     * @description ping
     *
     * @name PingList
     * @summary ping
     * @request GET:/ping
     */
    pingList: (params: RequestParams = {}) =>
      this.request<void, any>({
        path: `/ping`,
        method: 'GET',
        type: ContentType.Json,
        ...params,
      }),
  };
  version = {
    /**
     * @description version Description
     *
     * @name VersionList
     * @summary version Summary
     * @request GET:/version
     */
    versionList: (params: RequestParams = {}) =>
      this.request<void, any>({
        path: `/version`,
        method: 'GET',
        type: ContentType.Json,
        ...params,
      }),
  };
  workspace = {
    /**
     * @description list workspaces
     *
     * @tags workspace
     * @name WorkspaceList
     * @summary use to list workspaces
     * @request GET:/workspace
     * @secure
     */
    workspaceList: (
      query?: {
        /** query page */
        page?: number;
        /** query size */
        size?: number;
        /** query order, just like field1,field2:desc */
        orderBy?: string;
        /** query searchWord */
        searchWord?: string;
        /** query exact */
        exact?: boolean;
        /** query ids, split by comma */
        ids?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListWorkspacesResponse, ErrorsAppError>({
        path: `/workspace`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description import workspace
     *
     * @tags workspace
     * @name WorkspaceUpdate
     * @summary use to import workspace
     * @request PUT:/workspace
     * @secure
     */
    workspaceUpdate: (
      query: {
        /** workspace mount path */
        mountType: string;
        /** workspace mount type, only support nfs */
        mountPath: string;
      },
      data: {
        /** file */
        file: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersImportWorkspaceResponse, ErrorsAppError>({
        path: `/workspace`,
        method: 'PUT',
        query: query,
        body: data,
        secure: true,
        type: ContentType.FormData,
        format: 'json',
        ...params,
      }),

    /**
     * @description create workspace
     *
     * @tags workspace
     * @name WorkspaceCreate
     * @summary use to create workspace
     * @request POST:/workspace
     * @secure
     */
    workspaceCreate: (
      request: HandlersCreateWorkspaceRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersCreateWorkspaceResponse, ErrorsAppError>({
        path: `/workspace`,
        method: 'POST',
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get workspace
     *
     * @tags workspace
     * @name WorkspaceDetail
     * @summary use to get workspace by id
     * @request GET:/workspace/{id}
     * @secure
     */
    workspaceDetail: (id: string, params: RequestParams = {}) =>
      this.request<HandlersGetWorkspaceByIdResponse, ErrorsAppError>({
        path: `/workspace/${id}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description delete workspace
     *
     * @tags workspace
     * @name WorkspaceDelete
     * @summary use to delete workspace
     * @request DELETE:/workspace/{id}
     * @secure
     */
    workspaceDelete: (id: string, params: RequestParams = {}) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${id}`,
        method: 'DELETE',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description update workspace
     *
     * @tags workspace
     * @name WorkspacePartialUpdate
     * @summary use to update workspace
     * @request PATCH:/workspace/{id}
     * @secure
     */
    workspacePartialUpdate: (
      id: string,
      request: HandlersUpdateWorkspaceRequest,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${id}`,
        method: 'PATCH',
        body: request,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list notebook of workspace
     *
     * @tags notebook
     * @name WorkspaceIdNotebookList
     * @summary use to list notebook of workspace
     * @request GET:/workspace/{workspace-id}/notebook
     * @secure
     */
    workspaceIdNotebookList: (
      workspaceId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersListNotebooksResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebook`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get notebook content
     *
     * @tags notebook
     * @name WorkspaceIdNotebookDetail
     * @summary get notebook content
     * @request GET:/workspace/{workspace-id}/notebook/{name}
     * @secure
     */
    workspaceIdNotebookDetail: (
      workspaceId: string,
      name: string,
      params: RequestParams = {},
    ) =>
      this.request<any, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebook/${name}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description create notebook, update if name exist, set ipynb content in http body
     *
     * @tags notebook
     * @name WorkspaceIdNotebookUpdate
     * @summary use to create or update notebook
     * @request PUT:/workspace/{workspace-id}/notebook/{name}
     * @secure
     */
    workspaceIdNotebookUpdate: (
      workspaceId: string,
      name: string,
      request: NotebookIPythonNotebook,
      params: RequestParams = {},
    ) =>
      this.request<any, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebook/${name}`,
        method: 'PUT',
        body: request,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description delete notebook
     *
     * @tags notebook
     * @name WorkspaceIdNotebookDelete
     * @summary use to delete notebook
     * @request DELETE:/workspace/{workspace-id}/notebook/{name}
     * @secure
     */
    workspaceIdNotebookDelete: (
      workspaceId: string,
      name: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebook/${name}`,
        method: 'DELETE',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list notebook server
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverList
     * @summary use to list notebook server
     * @request GET:/workspace/{workspace-id}/notebookserver
     * @secure
     */
    workspaceIdNotebookserverList: (
      workspaceId: string,
      params: RequestParams = {},
    ) =>
      this.request<HertzListResponseItem[], ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver`,
        method: 'GET',
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
     * @description create notebook server
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverCreate
     * @summary use to create notebook server
     * @request POST:/workspace/{workspace-id}/notebookserver
     * @secure
     */
    workspaceIdNotebookserverCreate: (
      workspaceId: string,
      request: HertzCreateRequest,
      params: RequestParams = {},
    ) =>
      this.request<HertzCreateResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver`,
        method: 'POST',
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get notebook server
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverDetail
     * @summary use to get notebook server
     * @request GET:/workspace/{workspace-id}/notebookserver/{id}
     * @secure
     */
    workspaceIdNotebookserverDetail: (
      workspaceId: string,
      id: string,
      query?: {
        /** notebook object to edit */
        notebook?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HertzGetResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver/${id}`,
        method: 'GET',
        query: query,
        secure: true,
        format: 'json',
        ...params,
      }),

    /**
     * @description update notebook server settings
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverUpdate
     * @summary use to update notebook server settings
     * @request PUT:/workspace/{workspace-id}/notebookserver/{id}
     * @secure
     */
    workspaceIdNotebookserverUpdate: (
      workspaceId: string,
      id: string,
      request: HertzUpdateSettingsRequest,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver/${id}`,
        method: 'PUT',
        body: request,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description turn notebook server on or off
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverCreate2
     * @summary use to turn notebook server on or off
     * @request POST:/workspace/{workspace-id}/notebookserver/{id}
     * @originalName workspaceIdNotebookserverCreate
     * @duplicate
     * @secure
     */
    workspaceIdNotebookserverCreate2: (
      workspaceId: string,
      id: string,
      query?: {
        /** turn on notebook server */
        on?: boolean;
        /** turn off notebook server */
        off?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver/${id}`,
        method: 'POST',
        query: query,
        secure: true,
        ...params,
      }),

    /**
     * @description delete notebook server
     *
     * @tags notebook server
     * @name WorkspaceIdNotebookserverDelete
     * @summary use to delete notebook server
     * @request DELETE:/workspace/{workspace-id}/notebookserver/{id}
     * @secure
     */
    workspaceIdNotebookserverDelete: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/notebookserver/${id}`,
        method: 'DELETE',
        secure: true,
        ...params,
      }),

    /**
     * @description list workflow of workspace
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowList
     * @summary use to list workflows of workspace
     * @request GET:/workspace/{workspace-id}/workflow
     * @secure
     */
    workspaceIdWorkflowList: (
      workspaceId: string,
      query?: {
        /** page number */
        page?: number;
        /** page size */
        size?: number;
        /** support order field: name/createdAt, support order: asc/desc, seperated by comma, eg: createdAt:desc,name:asc */
        orderBy?: string;
        /** workflow name */
        searchWord?: string;
        /** exact */
        exact?: boolean;
        /** workspace ids seperated by comma */
        ids?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListWorkflowsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description create workflow, add workflow version if id is given
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowCreate
     * @summary use to create or update workflow
     * @request POST:/workspace/{workspace-id}/workflow
     * @secure
     */
    workspaceIdWorkflowCreate: (
      workspaceId: string,
      request: HandlersCreateWorkflowRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersCreateWorkflowResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow`,
        method: 'POST',
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get workflow
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowDetail
     * @summary get workflow
     * @request GET:/workspace/{workspace-id}/workflow/{id}
     * @secure
     */
    workspaceIdWorkflowDetail: (
      id: string,
      workspaceId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGetWorkflowResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${id}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description delete workflow
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowDelete
     * @summary use to delete workflow
     * @request DELETE:/workspace/{workspace-id}/workflow/{id}
     * @secure
     */
    workspaceIdWorkflowDelete: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${id}`,
        method: 'DELETE',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description update workspace
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowPartialUpdate
     * @summary use to update workflow
     * @request PATCH:/workspace/{workspace-id}/workflow/{id}
     * @secure
     */
    workspaceIdWorkflowPartialUpdate: (
      workspaceId: string,
      id: string,
      request: HandlersUpdateWorkflowRequest,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${id}`,
        method: 'PATCH',
        body: request,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list workflow of workspace
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowWorkflowIdFileList
     * @summary use to list workflow files
     * @request GET:/workspace/{workspace-id}/workflow/{workflow-id}/file
     * @secure
     */
    workspaceIdWorkflowWorkflowIdFileList: (
      workspaceId: string,
      workflowId: string,
      request: HandlersListWorkflowFilesRequest,
      query?: {
        /** page number */
        page?: number;
        /** page size */
        size?: number;
        /** support order field: version/path, support order: asc/desc, seperated by comma, eg: version:desc,path:asc */
        orderBy?: string;
        /** workflow name */
        searchWord?: string;
        /** workspace file ids seperated by comma */
        ids?: string;
        /** workspace version id */
        workflowVersionID?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListWorkflowFilesResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${workflowId}/file`,
        method: 'GET',
        query: query,
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get workflow file
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowWorkflowIdFileDetail
     * @summary get workflow file
     * @request GET:/workspace/{workspace-id}/workflow/{workflow-id}/file/{id}
     * @secure
     */
    workspaceIdWorkflowWorkflowIdFileDetail: (
      id: string,
      workspaceId: string,
      workflowId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGetWorkflowFileResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${workflowId}/file/${id}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get workflow version
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowWorkflowIdVersionDetail
     * @summary get workflow version
     * @request GET:/workspace/{workspace-id}/workflow/{workflow-id}/version/{id}
     * @secure
     */
    workspaceIdWorkflowWorkflowIdVersionDetail: (
      id: string,
      workspaceId: string,
      workflowId: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGetWorkflowVersionResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${workflowId}/version/${id}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description list workflow verions of workspace
     *
     * @tags workflow
     * @name WorkspaceIdWorkflowWorkflowIdVersionsList
     * @summary use to list workflow versions
     * @request GET:/workspace/{workspace-id}/workflow/{workflow-id}/versions
     * @secure
     */
    workspaceIdWorkflowWorkflowIdVersionsList: (
      workspaceId: string,
      workflowId: string,
      query?: {
        /** page number */
        page?: number;
        /** page size */
        size?: number;
        /** support order field: source/language/status, support order: asc/desc, seperated by comma, eg: status:desc,language:asc */
        orderBy?: string;
        /** workspace version ids seperated by comma */
        ids?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListWorkflowVersionsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/workflow/${workflowId}/versions`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description list data models
     *
     * @tags datamodel
     * @name DataModelDetail
     * @summary use to list data models
     * @request GET:/workspace/{workspace_id}/data_model
     * @secure
     */
    dataModelDetail: (
      workspaceId: string,
      query?: {
        /** data model types */
        types?: string[];
        /** query searchWord */
        searchWord?: string;
        /** data model ids */
        ids?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListDataModelsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description patch data model
     *
     * @tags datamodel
     * @name DataModelPartialUpdate
     * @summary use to patch data model
     * @request PATCH:/workspace/{workspace_id}/data_model
     * @secure
     */
    dataModelPartialUpdate: (
      workspaceId: string,
      request: HandlersPatchDataModelRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersPatchDataModelResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model`,
        method: 'PATCH',
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description get data model
     *
     * @tags datamodel
     * @name DataModelDetail2
     * @summary use to get data model
     * @request GET:/workspace/{workspace_id}/data_model/{id}
     * @originalName dataModelDetail
     * @duplicate
     * @secure
     */
    dataModelDetail2: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersGetDataModelResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model/${id}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description delete data model
     *
     * @tags datamodel
     * @name DataModelDelete
     * @summary use to delete data model,support delete with data model name/row ids/headers
     * @request DELETE:/workspace/{workspace_id}/data_model/{id}
     * @secure
     */
    dataModelDelete: (
      workspaceId: string,
      id: string,
      query?: {
        /** the data model headers should delete */
        headers?: string[];
        /** the data model row ids should delete */
        rowIDs?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model/${id}`,
        method: 'DELETE',
        query: query,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list data model rows
     *
     * @tags datamodel
     * @name DataModelRowsDetail
     * @summary use to list data model rows
     * @request GET:/workspace/{workspace_id}/data_model/{id}/rows
     * @secure
     */
    dataModelRowsDetail: (
      workspaceId: string,
      id: string,
      query?: {
        /** query page */
        page?: number;
        /** query size */
        size?: number;
        /** query order, just like field1,field2:desc */
        orderBy?: string;
        /** data model entity set reffed entity row ids */
        inSetIDs?: string[];
        /** query searchWord */
        searchWord?: string;
        /** data model row ids */
        rowIDs?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListDataModelRowsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model/${id}/rows`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description list all data model row ids
     *
     * @tags datamodel
     * @name DataModelRowsIdsDetail
     * @summary use to list all data model row ids
     * @request GET:/workspace/{workspace_id}/data_model/{id}/rows/ids
     * @secure
     */
    dataModelRowsIdsDetail: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<HandlersListAllDataModelRowIDsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/data_model/${id}/rows/ids`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description list submissions
     *
     * @tags submission
     * @name SubmissionDetail
     * @summary use to list submissions
     * @request GET:/workspace/{workspace_id}/submission
     * @secure
     */
    submissionDetail: (
      workspaceId: string,
      query?: {
        /** query page */
        page?: number;
        /** query size */
        size?: number;
        /** query order, just like field1,field2:desc */
        orderBy?: string;
        /** query searchWord */
        searchWord?: string;
        /** query exact */
        exact?: boolean;
        /** query ids */
        ids?: string[];
        /** workflow id */
        workflowID?: string;
        /** query status */
        status?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListSubmissionsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description create submission
     *
     * @tags submission
     * @name SubmissionCreate
     * @summary use to create submission
     * @request POST:/workspace/{workspace_id}/submission
     * @secure
     */
    submissionCreate: (
      workspaceId: string,
      request: HandlersCreateSubmissionRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlersCreateSubmissionResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission`,
        method: 'POST',
        body: request,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description delete submission
     *
     * @tags submission
     * @name SubmissionDelete
     * @summary use to delete submission
     * @request DELETE:/workspace/{workspace_id}/submission/{id}
     * @secure
     */
    submissionDelete: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${id}`,
        method: 'DELETE',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description cancel submission
     *
     * @tags submission
     * @name SubmissionCancelCreate
     * @summary use to cancel submission
     * @request POST:/workspace/{workspace_id}/submission/{id}/cancel
     * @secure
     */
    submissionCancelCreate: (
      workspaceId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${id}/cancel`,
        method: 'POST',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description check submission name unique
     *
     * @tags submission
     * @name SubmissionDetail2
     * @summary use to check submission name unique
     * @request GET:/workspace/{workspace_id}/submission/{name}
     * @originalName submissionDetail
     * @duplicate
     * @secure
     */
    submissionDetail2: (
      workspaceId: string,
      name: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${name}`,
        method: 'GET',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list runs
     *
     * @tags submission
     * @name SubmissionRunDetail
     * @summary use to list runs
     * @request GET:/workspace/{workspace_id}/submission/{submission_id}/run
     * @secure
     */
    submissionRunDetail: (
      workspaceId: string,
      submissionId: string,
      query?: {
        /** query page */
        page?: number;
        /** query size */
        size?: number;
        /** query order, just like field1,field2:desc */
        orderBy?: string;
        /** query searchWord */
        searchWord?: string;
        /** query ids */
        ids?: string[];
        /** query status */
        status?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListRunsResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${submissionId}/run`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),

    /**
     * @description cancel run
     *
     * @tags submission
     * @name SubmissionRunCancelCreate
     * @summary use to cancel run
     * @request POST:/workspace/{workspace_id}/submission/{submission_id}/run/{id}/cancel
     * @secure
     */
    submissionRunCancelCreate: (
      workspaceId: string,
      submissionId: string,
      id: string,
      params: RequestParams = {},
    ) =>
      this.request<void, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${submissionId}/run/${id}/cancel`,
        method: 'POST',
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description list tasks
     *
     * @tags submission
     * @name SubmissionRunTaskDetail
     * @summary use to list tasks
     * @request GET:/workspace/{workspace_id}/submission/{submission_id}/run/{run_id}/task
     * @secure
     */
    submissionRunTaskDetail: (
      workspaceId: string,
      submissionId: string,
      runId: string,
      query?: {
        /** query page */
        page?: number;
        /** query size */
        size?: number;
        /** query order, just like field1,field2:desc */
        orderBy?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<HandlersListTasksResponse, ErrorsAppError>({
        path: `/workspace/${workspaceId}/submission/${submissionId}/run/${runId}/task`,
        method: 'GET',
        query: query,
        secure: true,
        type: ContentType.Json,
        format: 'json',
        ...params,
      }),
  };
}
