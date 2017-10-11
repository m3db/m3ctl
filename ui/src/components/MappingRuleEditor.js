// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import React from 'react';
import {Button, Input, Form, Select} from 'antd';
import {withFormik} from 'formik';
import _ from 'lodash';
const POLICIES = [
  '1s:1h',
  '10s:6h',
  '10s:12h',
  '10s:1d',
  '10s:2d',
  '1m:2d',
  '1m:7d',
  '1m:14d',
  '1m:40d',
  '10m:2d',
  '10m:7d',
  '10m:14d',
  '10m:40d',
];

const AGGREGATION_FUNCTIONS = [
  'Min',
  'Max',
  'Mean',
  'Median',
  'Count',
  'Sum',
  'Sumsq',
  'Stdev',
  'P10',
  'P20',
  'P30',
  'P40',
  'P50',
  'P60',
  'P70',
  'P80',
  'P90',
  'P95',
  'P99',
  'P999',
  'P9999',
];

const FormItem = Form.Item;
const Option = Select.Option;
const formItemLayout = {
  labelCol: {
    xs: {span: 24},
    sm: {span: 4},
  },
  wrapperCol: {
    xs: {span: 24},
    sm: {span: 18},
  },
};
function MappingRuleEditor({mappingRule, values, handleChange, setFieldValue}) {
  return (
    <Form layout="horizontal">
      <div className="clearfix">
        <div className="col col-12">
          <FormItem
            required="true"
            colon={false}
            label="Name"
            {...formItemLayout}>
            <Input
              name="name"
              autoComplete="off"
              value={values.name}
              onChange={handleChange}
            />
          </FormItem>
        </div>
        <div className="col col-12">
          <FormItem
            required="true"
            colon={false}
            label="Metric Filter"
            {...formItemLayout}
            help="eg. tag1:value1 tag2:value2">
            <Input
              name="filter"
              autoComplete="off"
              value={values.filter}
              onChange={handleChange}
            />
          </FormItem>
        </div>
        <div className="col col-12">
          <FormItem
            colon={false}
            label="Policies"
            {...formItemLayout}
            help="eg:1m:10d = 1 min granularity for 10 days">
            {_.map(values.policies, (p, i) => {
              const policyParsed = p.split('|');
              const time = policyParsed[0];
              let aggFunctions;
              if (policyParsed.length > 1) {
                aggFunctions = policyParsed[0].split(',');
              }
              console.log(time, aggFunctions);
              return (
                <div
                  className="clearfix mt1 p1"
                  style={{border: '1px dashed lightgray'}}>
                  <div className="col col-6 px1">
                    <label className="block">Granularity/Retention Time</label>
                    <Select
                      value={time}
                      autoComplete="off"
                      onChange={e => {
                        values.policies[i] = e;
                        setFieldValue('policies', values.policies);
                      }}
                      style={{width: '100%'}}
                      notFoundContent={false}>
                      {POLICIES.map(opt => <Option key={opt}>{opt}</Option>)}
                    </Select>
                  </div>
                  <div className="col col-6">
                    <label className="block">Aggregation Function</label>

                    <Select
                      values={aggFunctions}
                      autoComplete="off"
                      mode="tags"
                      style={{width: '100%'}}
                      tokenSeparators={[',']}
                      notFoundContent={false}>
                      {AGGREGATION_FUNCTIONS.map(opt => (
                        <Option key={opt}>{opt}</Option>
                      ))}
                    </Select>
                  </div>
                </div>
              );
            })}
          </FormItem>
        </div>

        <div className="clearfix my1">
          <Button type="primary" htmlType="submit" className="right">
            Add
          </Button>
        </div>
      </div>
    </Form>
  );
}

export default withFormik({
  mapPropsToValues: ({mappingRule}) => ({
    name: mappingRule.name,
    filter: mappingRule.filter,
    policies: mappingRule.policies,
  }),
  mapPropsToFields({mappingRule = {}}) {
    return {
      name: {
        value: mappingRule.name,
      },
      filter: {
        value: mappingRule.filter,
      },
      policies: {
        value: mappingRule.policies,
      },
    };
  },
})(MappingRuleEditor);
