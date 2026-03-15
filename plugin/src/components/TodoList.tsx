import React from 'react';
import { Badge, Button, Card, ConfirmModal, Stack } from '@grafana/ui';

import type { Todo } from '../generated/todo/v1/todo_object_gen';
import type { Spec } from '../generated/todo/v1/types.spec.gen';

interface TodoListProps {
  todos: Todo[];
  onEdit: (todo: Todo) => void;
  onDelete: (name: string) => void;
  onStatusChange: (todo: Todo, status: Spec['status']) => void;
}

function statusColor(status: Spec['status']): 'blue' | 'orange' | 'green' {
  switch (status) {
    case 'open':
      return 'blue';
    case 'in_progress':
      return 'orange';
    case 'done':
      return 'green';
  }
}

function statusLabel(status: Spec['status']): string {
  switch (status) {
    case 'open':
      return 'Open';
    case 'in_progress':
      return 'In Progress';
    case 'done':
      return 'Done';
  }
}

function formatDate(timestamp?: string): string {
  if (!timestamp) {
    return '-';
  }
  return new Date(timestamp).toLocaleString();
}

export function TodoList({ todos, onEdit, onDelete, onStatusChange }: TodoListProps): React.ReactElement {
  const [deleteTarget, setDeleteTarget] = React.useState<string | null>(null);

  if (todos.length === 0) {
    return <p data-testid="todo-empty">No todos found. Create one to get started!</p>;
  }

  return (
    <>
      <div data-testid="todo-list">
        {todos.map((todo) => (
          <Card key={todo.metadata.uid ?? todo.metadata.name} data-testid={`todo-item-${todo.metadata.name}`}>
            <Card.Heading>{todo.spec.title}</Card.Heading>
            <Card.Meta>
              {[
                `Created: ${formatDate(todo.metadata.creationTimestamp)}`,
              ]}
            </Card.Meta>
            <Card.Description>{todo.spec.description ?? ''}</Card.Description>
            <Card.Tags>
              <Badge text={statusLabel(todo.spec.status)} color={statusColor(todo.spec.status)} />
            </Card.Tags>
            <Card.Actions>
              <Stack direction="row" gap={1}>
                {todo.spec.status !== 'done' && (
                  <Button
                    size="sm"
                    variant="secondary"
                    data-testid={`todo-advance-${todo.metadata.name}`}
                    onClick={() => {
                      const nextStatus: Spec['status'] = todo.spec.status === 'open' ? 'in_progress' : 'done';
                      onStatusChange(todo, nextStatus);
                    }}
                  >
                    {todo.spec.status === 'open' ? 'Start' : 'Complete'}
                  </Button>
                )}
                <Button size="sm" variant="secondary" onClick={() => onEdit(todo)} data-testid={`todo-edit-${todo.metadata.name}`}>
                  Edit
                </Button>
                <Button
                  size="sm"
                  variant="destructive"
                  onClick={() => setDeleteTarget(todo.metadata.name)}
                  data-testid={`todo-delete-${todo.metadata.name}`}
                >
                  Delete
                </Button>
              </Stack>
            </Card.Actions>
          </Card>
        ))}
      </div>

      <ConfirmModal
        isOpen={deleteTarget !== null}
        title="Delete TODO"
        body="Are you sure you want to delete this TODO? This action cannot be undone."
        confirmText="Delete"
        onConfirm={() => {
          if (deleteTarget) {
            onDelete(deleteTarget);
            setDeleteTarget(null);
          }
        }}
        onDismiss={() => setDeleteTarget(null)}
      />
    </>
  );
}
