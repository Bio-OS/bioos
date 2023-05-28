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

import {
  CSSProperties,
  ReactNode,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import { Form as FinalForm, FormSpy } from 'react-final-form';
import { Link, useRouteMatch } from 'react-router-dom';
import classNames from 'classnames';
import { FormApi } from 'final-form';
import {
  Alert,
  Button,
  Drawer,
  Form as ArcoForm,
  Input,
  Link as ArcoLinkCom,
  Message,
  Popover,
  Select,
  Space,
  Table,
  TableProps,
} from '@arco-design/web-react';
import {
  IconExclamationCircleFill,
  IconLeft,
  IconQuestionCircle,
} from '@arco-design/web-react/icon';

import SubTitle from 'components/SubTitle';
import { useGetEnvQuery } from 'helpers/hooks';
import { genTime, getSize, required } from 'helpers/utils';
import { getLengthRules } from 'helpers/validate';
import FieldItem from 'lib/field-item/FieldFormItem';
import Api from 'api/client';
import { HertzListResponseItem } from 'api/index';

import styles from './style.less';

const Option = Select.Option;

const IMAGE_TYPE = [
  { name: '预设镜像', value: 'pre' },
  { name: '自定义', value: 'custom' },
];

interface FormConfigureRuntime {
  image: string;
  resourceSize: number;
}

enum CalculationType {
  CPU = 'CPU',
  GPU = 'GPU',
}

const CALCULATIONTYPE = [CalculationType.CPU, CalculationType.GPU];

export default function ConfigureRuntime({
  render,
}: {
  render: (param: {
    setVisibleConfigureRuntime: React.Dispatch<React.SetStateAction<boolean>>;
  }) => ReactNode;
}) {
  const match = useRouteMatch<{ workspaceId: string }>();
  const workspaceId = match.params.workspaceId;
  const [visibleConfigureRuntime, setVisibleConfigureRuntime] = useState(false);
  const [visibleAlertUpdateEnv, setVisibleAlertUpdateEnv] = useState(false);
  const [formLoading, setFormLoading] = useState(false);
  const [currentImageID, setCurrentImageID] = useState('');
  const originNotebookServerConfig = useRef<FormConfigureRuntime>({
    image: '',
    resourceSize: 0,
  }); //  记录初始配置
  const [calculationType, setCalculationType] = useState<CalculationType>(
    CalculationType.CPU,
  );
  const [currentNotebookServer, setCurrentNotebookServer] =
    useState<HertzListResponseItem>();
  const [gpuDisabled, setGpuDisabled] = useState(false);
  const [loading, setLoading] = useState(false);

  const { notebook: notebookEnv } = useGetEnvQuery();

  // console.log(notebookEnv, 'notebookEnv');
  const imageOptions = notebookEnv?.officialImages;

  const currentImageInfo =
    imageOptions &&
    (imageOptions.find(item => item.image === currentImageID) ||
      imageOptions[0]);

  const getNotebookServer = async () => {
    setFormLoading(true);
    const { data: notebookServerList } =
      await Api.workspaceIdNotebookserverList(workspaceId);
    setCurrentNotebookServer(notebookServerList[0]);
    setFormLoading(false);
  };

  const isConfigBefore =
    currentNotebookServer?.image || currentNotebookServer?.resourceSize;

  useEffect(() => {
    if (!visibleConfigureRuntime) {
      return;
    }

    getNotebookServer();
  }, [visibleConfigureRuntime]);

  //  计算规格可能出现不是存量的，需要手动加入到select的options的第一条，并展示出来
  useEffect(() => {
    if (!currentNotebookServer || !notebookEnv) {
      return;
    }
    //  找到匹配的值
    const resourceSize = currentNotebookServer.resourceSize;
    let resourceSizeIndex = notebookEnv.resourceOptions.findIndex(
      a =>
        a.cpu === resourceSize.cpu &&
        a.memory === resourceSize.memory &&
        a.disk === resourceSize.disk,
    );

    if (resourceSizeIndex === -1) {
      optionsResource.unshift(resourceSize);
      resourceSizeIndex = 0;
    }

    // 记录原始类型和size
    originNotebookServerConfig.current = {
      image: currentNotebookServer.image,
      resourceSize: resourceSizeIndex,
    };

    refForm.current?.reset({
      ...refForm.current.getState().values,
      resourceSize: resourceSizeIndex,
    });
    //  设置初始值
  }, [currentNotebookServer, notebookEnv]);

  const optionsResource = useMemo(() => {
    if (!notebookEnv) {
      return [];
    }
    const options = notebookEnv.resourceOptions;
    const cpuOptions = options?.filter(item => !item.gpu);
    const gpuOptions = options?.filter(item => item.gpu);
    if (gpuOptions?.length === 0) {
      setGpuDisabled(true);
    }
    return calculationType === CalculationType.CPU ? cpuOptions : gpuOptions;
  }, [notebookEnv, calculationType]);

  useEffect(() => {
    if (!currentNotebookServer && imageOptions) {
      refForm.current?.reset({
        ...refForm.current.getState().values,
        image: imageOptions[0].image,
      });
    }
    if (
      !visibleConfigureRuntime ||
      !currentNotebookServer ||
      !imageOptions ||
      formLoading
    )
      return;

    setCurrentImageID(currentNotebookServer.image as string);
    //  判断是否为预设镜像
    const isPre = imageOptions.find(
      i => i.image === currentNotebookServer.image,
    );
    if (isPre) {
      refForm.current?.reset({
        ...refForm.current.getState().values,
        image: currentNotebookServer.image,
        type: 'pre',
      });
    } else {
      refForm.current?.reset({
        ...refForm.current.getState().values,
        imageUrl: currentNotebookServer.image,
        type: 'custom',
      });
    }
  }, [
    visibleConfigureRuntime,
    currentNotebookServer,
    imageOptions,
    formLoading,
  ]);

  useEffect(() => {
    if (!optionsResource) {
      return;
    }
    refForm.current?.reset({
      ...refForm.current.getState().values,
      resourceSize: 0,
      store: getSize(optionsResource[0].disk),
    });
  }, [calculationType, optionsResource]);

  const refForm = useRef<FormApi>();

  function checkChanged(values: FormConfigureRuntime) {
    const originConfig = originNotebookServerConfig.current;
    return !(
      originConfig.image === values?.image &&
      originConfig.resourceSize === values?.resourceSize
    );
  }

  const handleResourceSizeChange = value => {
    refForm.current?.reset({
      ...refForm.current.getState().values,
      store: getSize(optionsResource[value].disk),
    });
  };

  const handleSubmit = async values => {
    if (
      refForm.current?.getState().pristine &&
      isConfigBefore &&
      !checkChanged(values)
    ) {
      setVisibleConfigureRuntime(false);
      setVisibleAlertUpdateEnv(false);
      return;
    }

    setLoading(true);
    const hertzParam = {
      image: values.type === 'pre' ? values.image : values.imageUrl,
      resourceSize: optionsResource[values.resourceSize],
    };
    if (currentNotebookServer) {
      //  没有修改的参数不传入
      if (hertzParam.image === originNotebookServerConfig.current.image) {
        delete hertzParam.image;
      }
      if (
        values.resourceSize === originNotebookServerConfig.current.resourceSize
      ) {
        delete hertzParam.resourceSize;
      }
      const { data: updateRes } = await Api.workspaceIdNotebookserverUpdate(
        workspaceId,
        currentNotebookServer.id,
        hertzParam,
      );
    } else {
      const { data: createRes } = await Api.workspaceIdNotebookserverCreate(
        workspaceId,
        hertzParam,
      );
    }
    setLoading(false);
    Message.success('更新运行资源配置成功');
    setVisibleConfigureRuntime(false);
  };

  return (
    <>
      {render({
        setVisibleConfigureRuntime,
      })}

      <FinalForm
        onSubmit={handleSubmit}
        initialValues={{ type: 'pre' }}
        subscription={{ initialValues: true, pristine: true }}
        render={({ handleSubmit, form, pristine, errors }) => {
          refForm.current = form;
          return (
            <Drawer
              width={560}
              title="运行资源配置"
              className={styles.drawPadding}
              okText="更新环境"
              visible={visibleConfigureRuntime}
              unmountOnExit={true}
              onCancel={() => {
                setVisibleConfigureRuntime(false);
              }}
              footer={
                formLoading ? null : (
                  <FormSpy>
                    {({ values }) => (
                      <div className={styles.footerWrap}>
                        <Button
                          className="mr12"
                          onClick={() => {
                            setVisibleConfigureRuntime(false);
                            setVisibleAlertUpdateEnv(false);
                          }}
                        >
                          取消
                        </Button>

                        {isConfigBefore &&
                        (!errors || !Object.keys(errors).length) ? (
                          <Popover
                            trigger={'click'}
                            popupVisible={visibleAlertUpdateEnv}
                            onVisibleChange={val =>
                              setVisibleAlertUpdateEnv(val)
                            }
                            content={
                              <>
                                <section className="flexAlignCenter mb8">
                                  <IconExclamationCircleFill
                                    className="mr8 fs20"
                                    style={{ color: '#FF7D00' }}
                                  />
                                  <span className="fs14 fw500 colorBlack">
                                    {pristine &&
                                    !checkChanged(
                                      values as FormConfigureRuntime,
                                    )
                                      ? '当前无配置修改'
                                      : '确定更新环境吗？'}
                                  </span>
                                </section>
                                <section className="mb16">
                                  {pristine &&
                                  !checkChanged(values as FormConfigureRuntime)
                                    ? '注意：当前无配置修改，点击确定后不会更新环境。'
                                    : '修改环境参数将终止运行中Notebook，若未保存可能导致内容丢失，请谨慎操作。'}
                                </section>
                                <section className="textAlignRight">
                                  <Button
                                    size="mini"
                                    className="mr8"
                                    onClick={() =>
                                      setVisibleAlertUpdateEnv(false)
                                    }
                                  >
                                    放弃更新
                                  </Button>
                                  <Button
                                    type="primary"
                                    size="mini"
                                    onClick={() => {
                                      handleSubmit();
                                      setVisibleAlertUpdateEnv(false);
                                    }}
                                    loading={loading}
                                  >
                                    确定更新
                                  </Button>
                                </section>
                              </>
                            }
                          >
                            <Button type="primary">更新环境</Button>
                          </Popover>
                        ) : (
                          <Button
                            type="primary"
                            loading={loading}
                            onClick={() => {
                              handleSubmit();
                              setVisibleAlertUpdateEnv(false);
                            }}
                          >
                            更新环境
                          </Button>
                        )}
                      </div>
                    )}
                  </FormSpy>
                )
              }
            >
              <ArcoForm
                labelCol={{ span: 5 }}
                wrapperCol={{ span: 19 }}
                onSubmit={handleSubmit}
                labelAlign="left"
              >
                <SubTitle className="mb24" title="应用配置"></SubTitle>
                <FormSpy subscription={{ values: true }}>
                  {({ values }) => (
                    <FieldItem name="type" label="镜像来源" required={true}>
                      <>
                        <Space>
                          {IMAGE_TYPE.map(item => (
                            <Button
                              key={item.value}
                              className={
                                values.type === item.value
                                  ? styles.selectedType
                                  : styles.unselectedType
                              }
                              onClick={() => {
                                form.change('type', item.value);
                              }}
                            >
                              {item.name}
                            </Button>
                          ))}
                        </Space>
                        {values.type === 'custom' && (
                          <Alert
                            className={styles.customAlert}
                            type="warning"
                            content={
                              <div className="colorText3">
                                <div className="mb4">
                                  1、使用自定义镜像会增加启动时间
                                </div>
                                <div>
                                  2、较大镜像可能会导致 Notebook 启动超时
                                </div>
                              </div>
                            }
                          />
                        )}
                      </>
                    </FieldItem>
                  )}
                </FormSpy>

                <FormSpy subscription={{ values: true }}>
                  {({ values }) =>
                    values.type === 'custom' ? (
                      <FieldItem
                        name="imageUrl"
                        label="镜像地址"
                        required={true}
                        extra={<>自定义镜像必须基于 Bio-OS 基础镜像</>}
                        popoverProps={{ position: 'top' }}
                        rules={getLengthRules(7, 188)}
                      >
                        <Input placeholder="请输入" autoComplete="on" />
                      </FieldItem>
                    ) : (
                      <FieldItem
                        name="image"
                        label="容器镜像"
                        required={true}
                        validate={required}
                        validateFields={[]}
                        extra={
                          <Desc
                            itemArr={[
                              {
                                text: '版本号',
                                value: currentImageInfo?.version,
                              },
                              {
                                text: '更新时间',
                                value: genTime(currentImageInfo?.updateTime),
                              },
                              {
                                text: '描述',
                                value: currentImageInfo?.description,
                              },
                            ]}
                          />
                        }
                      >
                        <Select
                          placeholder="请选择"
                          onChange={val => setCurrentImageID(val)}
                        >
                          {imageOptions?.map(item => (
                            <Option key={item.image} value={item.image}>
                              {item.name}
                            </Option>
                          ))}
                        </Select>
                      </FieldItem>
                    )
                  }
                </FormSpy>

                <SubTitle className="mb24" title="资源配置"></SubTitle>

                <FieldItem
                  required={true}
                  name="calculation_type"
                  label="计算类型"
                >
                  <>
                    <Space>
                      {CALCULATIONTYPE.map((item, i) => {
                        return (
                          <Button
                            disabled={
                              gpuDisabled && item === CalculationType.GPU
                            }
                            key={i}
                            onClick={() => {
                              if (item !== calculationType) {
                                setCalculationType(item);
                              }
                            }}
                            className={
                              calculationType === item
                                ? styles.selectedType
                                : styles.unselectedType
                            }
                          >
                            {item}
                          </Button>
                        );
                      })}
                    </Space>
                  </>
                </FieldItem>

                <FieldItem name="resourceSize" required={true} label="计算规格">
                  <Select
                    placeholder="请选择"
                    onChange={val => handleResourceSizeChange(val)}
                  >
                    {optionsResource?.map((item, index) => {
                      return calculationType === CalculationType.GPU ? (
                        <Select.Option value={index} key={index}>
                          <span
                            style={{
                              padding: 2,
                              background: '#FDEDD9',
                              borderRadius: 4,
                            }}
                            className="mr8"
                          >{`${item.gpu?.model} (${getSize(
                            item.gpu?.memory,
                          )})`}</span>
                          <span>{`${item.cpu} Core ${getSize(
                            item.memory,
                          )}`}</span>
                        </Select.Option>
                      ) : (
                        <Select.Option value={index} key={index}>
                          <span className="mr8">{`${item.cpu} Core`}</span>
                          <span>{getSize(item.memory)}</span>
                        </Select.Option>
                      );
                    })}
                  </Select>
                </FieldItem>

                <FieldItem
                  name="store"
                  label={
                    <span className="flexAlignCenter">
                      <span className="ml16">存储规格</span>
                    </span>
                  }
                  disabled={true}
                >
                  <Input className="mb12" />
                </FieldItem>
              </ArcoForm>
            </Drawer>
          );
        }}
      ></FinalForm>
    </>
  );
}

function Desc({
  itemArr,
  className,
  style,
}: {
  itemArr: { text: string; value: ReactNode }[];
  className?: string;
  style?: CSSProperties;
}) {
  return (
    <section
      className={classNames(['fs12', styles.descSection, { className }])}
      style={style}
    >
      {itemArr.map(({ text, value }) => (
        <div className="flexRow pt8 pb8" key={text}>
          <div className={classNames(['colorBlack2', styles.descLeft])}>
            {text}
          </div>
          <div className="colorBlack3">{value}</div>
        </div>
      ))}
    </section>
  );
}
