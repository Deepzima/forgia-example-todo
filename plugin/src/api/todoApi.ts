import { getBackendSrv } from '@grafana/runtime';

import type { Todo, Metadata } from '../generated/todo/v1/todo_object_gen';
import type { Spec } from '../generated/todo/v1/types.spec.gen';

export interface TodoList {
  items: Todo[];
  metadata: {
    continue?: string;
    remainingItemCount?: number;
  };
}

const BASE_PATH = '/apis/todo.grafana.app/v1/namespaces';

function namespacePath(namespace: string): string {
  return `${BASE_PATH}/${namespace}/todos`;
}

export async function listTodos(namespace: string): Promise<TodoList> {
  return getBackendSrv().get<TodoList>(namespacePath(namespace));
}

export async function getTodo(namespace: string, name: string): Promise<Todo> {
  return getBackendSrv().get<Todo>(`${namespacePath(namespace)}/${name}`);
}

export async function createTodo(
  namespace: string,
  spec: Spec
): Promise<Todo> {
  return getBackendSrv().post<Todo>(namespacePath(namespace), {
    metadata: { namespace } as Partial<Metadata>,
    spec,
  });
}

export async function updateTodo(
  namespace: string,
  name: string,
  todo: Todo
): Promise<Todo> {
  return getBackendSrv().put<Todo>(`${namespacePath(namespace)}/${name}`, todo);
}

export async function deleteTodo(
  namespace: string,
  name: string
): Promise<void> {
  return getBackendSrv().delete(`${namespacePath(namespace)}/${name}`);
}
