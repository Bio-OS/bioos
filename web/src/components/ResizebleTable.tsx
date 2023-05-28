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

import { forwardRef, RefObject } from 'react';
import { Resizable, ResizableProps } from 'react-resizable';
import classNames from 'classnames';
import { Table, TableColumnProps, TableProps } from '@arco-design/web-react';

import styles from './style.less';

export default function ResizebleTable({
  columns = [],
  minSize = 100,
  className,
  onResize,
  ...rest
}: TableProps & {
  minSize?: number;
  onResize: (index: number, width: number) => void;
}) {
  function handleResize(index: number) {
    return (_e: any, { size }: { size: { width: number } }) => {
      document.getElementsByTagName('html')[0].className = 'cursorResize';
      onResize(index, size.width < minSize ? minSize : size.width);
    };
  }

  function handleResizeEnd() {
    document.getElementsByTagName('html')[0].className = '';
  }

  const resizebleColumns = columns.map(
    (column: TableColumnProps, index: number) => {
      const last = index === columns.length - 1;
      if (last) return column;

      return {
        ...column,
        onHeaderCell: (col: TableColumnProps) => ({
          width: col.width,
          onResize: handleResize(index),
          onResizeStop: handleResizeEnd,
        }),
      };
    },
  );

  const components = {
    header: {
      th: ResizableTitle,
    },
  };

  return (
    <Table
      className={classNames([styles.resizeTable, className])}
      components={components}
      border={true}
      borderCell={true}
      scroll={{ x: true }}
      columns={resizebleColumns}
      {...rest}
    />
  );
}

const ResizableTitle = (props: ResizableProps) => {
  const { onResize, onResizeStop, width, ...restProps } = props;

  if (!width) {
    return <th {...restProps} />;
  }

  return (
    <Resizable
      width={width}
      height={0}
      handle={<CustomResizeHandle />}
      onResize={onResize}
      onResizeStop={onResizeStop}
      draggableOpts={{
        enableUserSelectHack: false,
      }}
    >
      <th {...restProps} />
    </Resizable>
  );
};

const CustomResizeHandle = forwardRef(
  (props: { handleAxis?: string }, ref: RefObject<HTMLSpanElement>) => {
    const { handleAxis, ...restProps } = props;

    return (
      <span
        ref={ref}
        className="react-resizable-handle"
        {...restProps}
        onClick={e => {
          e.stopPropagation();
        }}
      />
    );
  },
);
