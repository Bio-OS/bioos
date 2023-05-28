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

import React, { ReactNode } from 'react';
import classNames from 'classnames';
import { Alert, Typography } from '@arco-design/web-react';

import styles from './style.less';

interface DetailPageProps {
  showBorderBottom?: boolean;
  title: string;
  breadcrumbs: ReactNode;
  rightArea?: React.ReactNode;
  statusTag?: ReactNode;
  description?: ReactNode;
  tipInfo?:
    | {
        text: React.ReactNode;
        type?: 'info' | 'success' | 'warning' | 'error';
        closable?: boolean;
        onClose?: () => void;
      }
    | undefined;
  children: React.ReactNode;
  className?: string;
  contentStyle?: React.CSSProperties;
}

export default function DetailPage({
  title,
  breadcrumbs,
  rightArea,
  statusTag,
  children,
  tipInfo,
  description,
  showBorderBottom = true,
  className,
  contentStyle,
}: DetailPageProps) {
  return (
    <div className={classNames([styles.detailPageContainer, className])}>
      <div
        className={styles.header}
        style={{
          borderBottom: showBorderBottom ? '1px solid #eaedf1' : 0,
        }}
      >
        <div className={styles.leftWrap}>
          <div>{breadcrumbs}</div>
          <div className="mt10 lh26 flexAlignCenter">
            <Typography.Paragraph
              className="fs18 fw500 mr12"
              style={{ maxWidth: 400 }}
              ellipsis={{
                cssEllipsis: true,
                showTooltip: {
                  type: 'popover',
                },
              }}
            >
              {title}
            </Typography.Paragraph>
            {statusTag}
          </div>
          {description && (
            <div className="flexAlignCenter fs12 colorGrey mt8">
              {description}
            </div>
          )}
        </div>
        <div>{rightArea}</div>
      </div>
      {tipInfo && (
        <Alert
          content={tipInfo.text}
          type={tipInfo.type}
          closable={tipInfo.closable}
          onClose={tipInfo.onClose}
        />
      )}
      <div className={styles.content} style={contentStyle}>
        {children}
      </div>
    </div>
  );
}
