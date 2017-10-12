import {getTypeTag} from './index';

test('getTypeTag', () => {
  expect(getTypeTag('type:timer')).toBe('timer');
  expect(getTypeTag('type:a')).toBe('a');
  expect(getTypeTag('this:test type')).toBe(null);
  expect(getTypeTag('type:counts ')).toBe('counts');
});
