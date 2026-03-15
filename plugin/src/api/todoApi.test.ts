import { getBackendSrv } from '@grafana/runtime';

import { listTodos, getTodo, createTodo, updateTodo, deleteTodo } from './todoApi';
import type { Todo } from '../generated/todo/v1/todo_object_gen';

// The @grafana/runtime mock is loaded via moduleNameMapper.
// getBackendSrv is already a jest.fn() returning {get, post, put, delete} as jest.fn()s.

const mockBackendSrv = (getBackendSrv as jest.Mock)();

const makeTodo = (name: string): Todo => ({
  kind: 'Todo',
  apiVersion: 'todo.grafana.app/v1',
  metadata: { name, namespace: 'default', uid: `uid-${name}`, creationTimestamp: '2026-03-15T10:00:00Z' },
  spec: { title: `Todo ${name}`, status: 'open' },
  status: {},
});

describe('todoApi', () => {
  const ns = 'default';
  const basePath = `/apis/todo.grafana.app/v1/namespaces/${ns}/todos`;

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('listTodos', () => {
    it('calls GET on the correct path', async () => {
      const response = { items: [makeTodo('t1')], metadata: {} };
      mockBackendSrv.get.mockResolvedValue(response);

      const result = await listTodos(ns);

      expect(mockBackendSrv.get).toHaveBeenCalledWith(basePath);
      expect(result.items).toHaveLength(1);
    });

    it('propagates errors', async () => {
      mockBackendSrv.get.mockRejectedValue(new Error('Server error'));

      await expect(listTodos(ns)).rejects.toThrow('Server error');
    });
  });

  describe('getTodo', () => {
    it('calls GET with name in path', async () => {
      const todo = makeTodo('t1');
      mockBackendSrv.get.mockResolvedValue(todo);

      const result = await getTodo(ns, 't1');

      expect(mockBackendSrv.get).toHaveBeenCalledWith(`${basePath}/t1`);
      expect(result.metadata.name).toBe('t1');
    });
  });

  describe('createTodo', () => {
    it('calls POST with spec in body', async () => {
      const todo = makeTodo('new');
      mockBackendSrv.post.mockResolvedValue(todo);

      const spec = { title: 'New Todo', status: 'open' as const };
      const result = await createTodo(ns, spec);

      expect(mockBackendSrv.post).toHaveBeenCalledWith(basePath, {
        metadata: { namespace: ns },
        spec,
      });
      expect(result.kind).toBe('Todo');
    });

    it('propagates 400 errors', async () => {
      mockBackendSrv.post.mockRejectedValue(new Error('Bad Request'));

      await expect(createTodo(ns, { title: '', status: 'open' })).rejects.toThrow('Bad Request');
    });
  });

  describe('updateTodo', () => {
    it('calls PUT with full todo in body', async () => {
      const todo = makeTodo('t1');
      todo.spec.title = 'Updated';
      mockBackendSrv.put.mockResolvedValue(todo);

      const result = await updateTodo(ns, 't1', todo);

      expect(mockBackendSrv.put).toHaveBeenCalledWith(`${basePath}/t1`, todo);
      expect(result.spec.title).toBe('Updated');
    });
  });

  describe('deleteTodo', () => {
    it('calls DELETE with name in path', async () => {
      mockBackendSrv.delete.mockResolvedValue(undefined);

      await deleteTodo(ns, 't1');

      expect(mockBackendSrv.delete).toHaveBeenCalledWith(`${basePath}/t1`);
    });

    it('propagates 404 errors', async () => {
      mockBackendSrv.delete.mockRejectedValue(new Error('Not Found'));

      await expect(deleteTodo(ns, 'nonexistent')).rejects.toThrow('Not Found');
    });
  });
});
