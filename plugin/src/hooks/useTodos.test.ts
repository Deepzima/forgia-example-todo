import { renderHook, act, waitFor } from '@testing-library/react';

import { useTodos } from './useTodos';
import * as todoApi from '../api/todoApi';
import type { Todo } from '../generated/todo/v1/todo_object_gen';

jest.mock('../api/todoApi');

const mockedApi = todoApi as jest.Mocked<typeof todoApi>;

const makeTodo = (name: string): Todo => ({
  kind: 'Todo',
  apiVersion: 'todo.grafana.app/v1',
  metadata: { name, namespace: 'default', uid: `uid-${name}`, creationTimestamp: '2026-03-15T10:00:00Z' },
  spec: { title: `Todo ${name}`, status: 'open' },
  status: {},
});

describe('useTodos', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockedApi.listTodos.mockResolvedValue({ items: [makeTodo('t1')], metadata: {} });
  });

  it('loads todos on mount', async () => {
    const { result } = renderHook(() => useTodos('default'));

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.todos).toHaveLength(1);
    expect(result.current.todos[0].metadata.name).toBe('t1');
    expect(result.current.error).toBeNull();
  });

  it('sets error on load failure', async () => {
    mockedApi.listTodos.mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useTodos('default'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.error).toBe('Network error');
    expect(result.current.todos).toHaveLength(0);
  });

  it('creates a todo and reloads', async () => {
    mockedApi.createTodo.mockResolvedValue(makeTodo('t2'));
    mockedApi.listTodos
      .mockResolvedValueOnce({ items: [makeTodo('t1')], metadata: {} })
      .mockResolvedValueOnce({ items: [makeTodo('t1'), makeTodo('t2')], metadata: {} });

    const { result } = renderHook(() => useTodos('default'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await act(async () => {
      await result.current.create({ title: 'New', status: 'open' });
    });

    expect(mockedApi.createTodo).toHaveBeenCalledWith('default', { title: 'New', status: 'open' });
    expect(result.current.todos).toHaveLength(2);
  });

  it('updates a todo and reloads', async () => {
    const updated = makeTodo('t1');
    updated.spec.title = 'Updated';
    mockedApi.updateTodo.mockResolvedValue(updated);
    mockedApi.listTodos
      .mockResolvedValueOnce({ items: [makeTodo('t1')], metadata: {} })
      .mockResolvedValueOnce({ items: [updated], metadata: {} });

    const { result } = renderHook(() => useTodos('default'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await act(async () => {
      await result.current.update('t1', updated);
    });

    expect(mockedApi.updateTodo).toHaveBeenCalledWith('default', 't1', updated);
  });

  it('removes a todo and reloads', async () => {
    mockedApi.deleteTodo.mockResolvedValue(undefined);
    mockedApi.listTodos
      .mockResolvedValueOnce({ items: [makeTodo('t1')], metadata: {} })
      .mockResolvedValueOnce({ items: [], metadata: {} });

    const { result } = renderHook(() => useTodos('default'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await act(async () => {
      await result.current.remove('t1');
    });

    expect(mockedApi.deleteTodo).toHaveBeenCalledWith('default', 't1');
    expect(result.current.todos).toHaveLength(0);
  });

  it('sets error on create failure', async () => {
    mockedApi.createTodo.mockRejectedValue(new Error('Create failed'));

    const { result } = renderHook(() => useTodos('default'));

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    await act(async () => {
      try {
        await result.current.create({ title: 'Fail', status: 'open' });
      } catch {
        // expected
      }
    });

    expect(result.current.error).toBe('Create failed');
  });
});
