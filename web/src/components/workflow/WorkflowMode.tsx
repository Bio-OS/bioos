import { memo} from 'react';
import classNames from 'classnames';
import { Radio } from '@arco-design/web-react';


const RADIO_ARRAY = [
  {
    title: 'WDL',
    value: 'WDL',
  },
  {
    title: 'CWL',
    value: 'CWL',
  },
  {
    title: 'Nextflow',
    description: '开发中',
    value: 'NFL',
    disabled: true,
  },
  {
    title: 'Snakemake',
    description: '开发中',
    value: 'SMK',
    disabled: true,
  },
];

export type WorkflowModeType = (typeof RADIO_ARRAY)[number]['value'];

function WorkflowMode({
  value,
  onChange,
  disabled,
}: {
  value: string;
  onChange: (val: string) => void;
  disabled?: boolean;
}) {
  return (
    <Radio.Group value={value} className="flex" disabled={disabled}>
      {RADIO_ARRAY.map(item => {
        return (
          <div
            key={item.value}
            className={classNames([
              'fs12 br4 mr12',
              { cursorPointer: !disabled && !item.disabled },
            ])}
            onClick={() => {
              if (disabled || item.disabled) return;
              onChange(item.value);
            }}
            style={{
              padding: '12px 16px',
              width: 114,
              border:
                item.value === value
                  ? '1px solid #94c2ff'
                  : '1px solid #e4e8ff',
              background: item.value === value ? '#e8f4ff' : 'white',
            }}
          >
            <div className="flexJustifyBetween">
              <div className="flexAlignCenter">
                <span className="colorBlack fw500 mr4">{item.title}</span>
              </div>

              <Radio value={item.value} style={{ marginRight: 0 }} />
            </div>

         
          </div>
        );
      })}
    </Radio.Group>
  );
}

export default memo(WorkflowMode);