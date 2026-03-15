import { useState, useCallback, useEffect } from 'react';

import type { Todo } from '../generated/todo/v1/todo_object_gen';
import type { Spec } from '../generated/todo/v1/types.spec.gen';
import * as todoApi from '../api/todoApi';

interface UseTodosState {
  todos: Todo[];
  isLoading: boolean;
  error: string | null;
}

interface UseTodosResult extends UseTodosState {
  reload: () => Promise<void>;
  create: (spec: Spec) => Promise<void>;
  update: (name: string, todo: Todo) => Promise<void>;
  remove: (name: string) => Promise<void>;
}

export function useTodos(namespace: string): UseTodosResult {
  const [state, setState] = useState<UseTodosState>({
    todos: [],
    isLoading: true,
    error: null,
  });

  const reload = useCallback(async () => {
    setState((prev) => ({ ...prev, isLoading: true, error: null }));
    try {
      const result = await todoApi.listTodos(namespace);
      setState({ todos: result.items ?? [], isLoading: false, error: null });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load todos';
      setState((prev) => ({ ...prev, isLoading: false, error: message }));
    }
  }, [namespace]);

  const create = useCallback(
    async (spec: Spec) => {
      try {
        await todoApi.createTodo(namespace, spec);
        await reload();
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to create todo';
        setState((prev) => ({ ...prev, error: message }));
        throw err;
      }
    },
    [namespace, reload]
  );

  const update = useCallback(
    async (name: string, todo: Todo) => {
      try {
        await todoApi.updateTodo(namespace, name, todo);
        await reload();
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to update todo';
        setState((prev) => ({ ...prev, error: message }));
        throw err;
      }
    },
    [namespace, reload]
  );

  const remove = useCallback(
    async (name: string) => {
      try {
        await todoApi.deleteTodo(namespace, name);
        await reload();
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to delete todo';
        setState((prev) => ({ ...prev, error: message }));
        throw err;
      }
    },
    [namespace, reload]
  );

  useEffect(() => {
    reload();
  }, [reload]);

  return { ...state, reload, create, update, remove };
}
