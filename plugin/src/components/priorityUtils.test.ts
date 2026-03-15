import { sortByPriority, filterByPriority, getPriority } from './priorityUtils';
import type { Spec } from '../generated/todo/v1/types.spec.gen';

const makeItem = (priority?: Spec['priority']) => ({
  spec: { title: 'T', status: 'open' as const, priority },
});

describe('priorityUtils', () => {
  describe('getPriority', () => {
    it('returns the priority when present', () => {
      expect(getPriority({ title: 'T', status: 'open', priority: 'critical' })).toBe('critical');
    });

    it('defaults to "medium" when priority is missing and logs warning', () => {
      const warnSpy = jest.spyOn(console, 'warn').mockImplementation(() => {});
      expect(getPriority({ title: 'T', status: 'open' })).toBe('medium');
      expect(warnSpy).toHaveBeenCalledWith(expect.stringContaining('missing priority'));
      warnSpy.mockRestore();
    });
  });

  describe('sortByPriority', () => {
    const items = [
      makeItem('low'),
      makeItem('critical'),
      makeItem('medium'),
      makeItem('high'),
    ];

    it('sorts descending (critical first) by default', () => {
      const sorted = sortByPriority(items);
      expect(sorted.map((i) => i.spec.priority)).toEqual(['critical', 'high', 'medium', 'low']);
    });

    it('sorts ascending (low first) when descending=false', () => {
      const sorted = sortByPriority(items, false);
      expect(sorted.map((i) => i.spec.priority)).toEqual(['low', 'medium', 'high', 'critical']);
    });

    it('does not mutate the original array', () => {
      const original = [...items];
      sortByPriority(items);
      expect(items).toEqual(original);
    });
  });

  describe('filterByPriority', () => {
    const items = [
      makeItem('low'),
      makeItem('critical'),
      makeItem('medium'),
      makeItem('high'),
    ];

    it('returns all items when filter is empty', () => {
      expect(filterByPriority(items, [])).toEqual(items);
    });

    it('filters to only matching priorities', () => {
      const filtered = filterByPriority(items, ['high', 'critical']);
      expect(filtered.map((i) => i.spec.priority)).toEqual(['critical', 'high']);
    });

    it('returns empty array when no items match', () => {
      const filtered = filterByPriority([makeItem('low')], ['critical']);
      expect(filtered).toEqual([]);
    });
  });
});
