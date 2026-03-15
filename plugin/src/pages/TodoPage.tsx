import React, { useMemo, useState } from 'react';
import { Alert, Button, RadioButtonGroup, Select, Spinner, Stack } from '@grafana/ui';
import { PluginPage } from '@grafana/runtime';

import type { Todo } from '../generated/todo/v1/todo_object_gen';
import type { Spec } from '../generated/todo/v1/types.spec.gen';
import { useTodos } from '../hooks/useTodos';
import { TodoList } from '../components/TodoList';
import { TodoForm } from '../components/TodoForm';
import { PRIORITY_OPTIONS, sortByPriority, filterByPriority } from '../components/priorityUtils';
import type { Priority } from '../components/priorityUtils';

const NAMESPACE = 'default';

type ViewMode = 'list' | 'create' | 'edit';
type SortDirection = 'desc' | 'asc' | 'none';

const SORT_OPTIONS = [
  { label: 'No sort', value: 'none' as const },
  { label: 'Critical first', value: 'desc' as const },
  { label: 'Low first', value: 'asc' as const },
];

const FILTER_OPTIONS = [
  { label: 'All', value: '' },
  ...PRIORITY_OPTIONS.map((o) => ({ label: o.label, value: o.value })),
];

export function TodoPage(): React.ReactElement {
  const { todos, isLoading, error, create, update, remove } = useTodos(NAMESPACE);
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [sortDirection, setSortDirection] = useState<SortDirection>('none');
  const [priorityFilter, setPriorityFilter] = useState<Priority | ''>('');

  const displayedTodos = useMemo(() => {
    let result = todos;
    if (priorityFilter) {
      result = filterByPriority(result, [priorityFilter]);
    }
    if (sortDirection !== 'none') {
      result = sortByPriority(result, sortDirection === 'desc');
    }
    return result;
  }, [todos, sortDirection, priorityFilter]);

  const handleCreate = async (spec: Spec): Promise<void> => {
    setIsSubmitting(true);
    try {
      await create(spec);
      setViewMode('list');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUpdate = async (spec: Spec): Promise<void> => {
    if (!editingTodo) {
      return;
    }
    setIsSubmitting(true);
    try {
      await update(editingTodo.metadata.name, {
        ...editingTodo,
        spec,
      });
      setViewMode('list');
      setEditingTodo(null);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (todo: Todo): void => {
    setEditingTodo(todo);
    setViewMode('edit');
  };

  const handleStatusChange = async (todo: Todo, status: Spec['status']): Promise<void> => {
    await update(todo.metadata.name, {
      ...todo,
      spec: { ...todo.spec, status },
    });
  };

  const handleDelete = async (name: string): Promise<void> => {
    await remove(name);
  };

  const handleCancel = (): void => {
    setViewMode('list');
    setEditingTodo(null);
  };

  return (
    <PluginPage>
      <Stack direction="column" gap={2}>
        {error && (
          <Alert severity="error" title="Error" data-testid="todo-error">
            {error}
          </Alert>
        )}

        {viewMode === 'list' && (
          <>
            <Stack direction="row" gap={2}>
              <Button onClick={() => setViewMode('create')} data-testid="todo-create-btn">
                Create TODO
              </Button>
              <div data-testid="todo-sort-control">
                <RadioButtonGroup
                  options={SORT_OPTIONS}
                  value={sortDirection}
                  onChange={(v) => setSortDirection(v as SortDirection)}
                />
              </div>
              <div data-testid="todo-filter-control">
                <Select
                  data-testid="todo-priority-filter"
                  options={FILTER_OPTIONS}
                  value={priorityFilter}
                  onChange={(v) => setPriorityFilter((v.value ?? '') as Priority | '')}
                  placeholder="Filter by priority"
                />
              </div>
            </Stack>

            {isLoading ? (
              <Spinner data-testid="todo-spinner" />
            ) : (
              <TodoList
                todos={displayedTodos}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onStatusChange={handleStatusChange}
              />
            )}
          </>
        )}

        {viewMode === 'create' && (
          <TodoForm onSubmit={handleCreate} onCancel={handleCancel} isSubmitting={isSubmitting} />
        )}

        {viewMode === 'edit' && editingTodo && (
          <TodoForm
            initialValues={editingTodo.spec}
            onSubmit={handleUpdate}
            onCancel={handleCancel}
            isSubmitting={isSubmitting}
          />
        )}
      </Stack>
    </PluginPage>
  );
}
