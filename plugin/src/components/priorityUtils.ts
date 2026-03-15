import type { Spec } from '../generated/todo/v1/types.spec.gen';

export type Priority = NonNullable<Spec['priority']>;

export const PRIORITY_OPTIONS = [
  { label: 'Low', value: 'low' as const },
  { label: 'Medium', value: 'medium' as const },
  { label: 'High', value: 'high' as const },
  { label: 'Critical', value: 'critical' as const },
] as const;

export const priorityColorMap: Record<Priority, 'blue' | 'yellow' | 'orange' | 'red'> = {
  low: 'blue',
  medium: 'yellow',
  high: 'orange',
  critical: 'red',
};

const prioritySortWeight: Record<Priority, number> = {
  low: 1,
  medium: 2,
  high: 3,
  critical: 4,
};

export function getPriority(spec: Spec): Priority {
  if (spec.priority) {
    return spec.priority;
  }
  console.warn('Todo missing priority field, defaulting to "medium"');
  return 'medium';
}

/** Sort todos by priority. `descending = true` means critical first. */
export function sortByPriority<T extends { spec: Spec }>(todos: T[], descending = true): T[] {
  return [...todos].sort((a, b) => {
    const wA = prioritySortWeight[getPriority(a.spec)];
    const wB = prioritySortWeight[getPriority(b.spec)];
    return descending ? wB - wA : wA - wB;
  });
}

/** Filter todos to only those matching the given priority levels. */
export function filterByPriority<T extends { spec: Spec }>(todos: T[], priorities: Priority[]): T[] {
  if (priorities.length === 0) {
    return todos;
  }
  return todos.filter((t) => priorities.includes(getPriority(t.spec)));
}
