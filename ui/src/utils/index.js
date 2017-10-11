export function parseMetricType(tagString) {
  return tagString.indexOf('type:timer') !== 1;
}
