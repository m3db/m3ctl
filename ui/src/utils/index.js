import _ from 'lodash';

export function getTypeTag(tagString) {
  const regex = /type:([^\s]*)/g;
  const results = regex.exec(tagString);
  return _.nth(results, 1) || null;
}

export function filterPoliciesBasedOnTag(policies, typeTag) {
  if (_.includes(['gauge', 'counter'], typeTag)) {
    return policies;
  }
  return _.filter(policies, p => !p.startsWith('10m'));
}
