import React, { useState } from 'react';
import { Alert, Button, Spinner, Stack } from '@grafana/ui';
import { PluginPage } from '@grafana/runtime';

import type { Todo } from '../generated/todo/v1/todo_object_gen';
import type { Spec } from '../generated/todo/v1/types.spec.gen';
import { useTodos } from '../hooks/useTodos';
import { TodoList } from '../components/TodoList';
import { TodoForm } from '../components/TodoForm';

const NAMESPACE = 'default';

type ViewMode = 'list' | 'create' | 'edit';

export function TodoPage(): React.ReactElement {
  const { todos, isLoading, error, create, update, remove } = useTodos(NAMESPACE);
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

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
            <div>
              <Button onClick={() => setViewMode('create')} data-testid="todo-create-btn">
                Create TODO
              </Button>
            </div>

            {isLoading ? (
              <Spinner data-testid="todo-spinner" />
            ) : (
              <TodoList
                todos={todos}
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
