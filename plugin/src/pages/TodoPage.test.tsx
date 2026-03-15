import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { TodoPage } from './TodoPage';
import * as todoApi from '../api/todoApi';
import type { Todo } from '../generated/todo/v1/todo_object_gen';

// @grafana/runtime and @grafana/ui are auto-mocked via moduleNameMapper.
// Mock the API module for controlling responses.
jest.mock('../api/todoApi');

const mockedApi = todoApi as jest.Mocked<typeof todoApi>;

const makeTodo = (name: string, status: 'open' | 'in_progress' | 'done' = 'open'): Todo => ({
  kind: 'Todo',
  apiVersion: 'todo.grafana.app/v1',
  metadata: { name, namespace: 'default', uid: `uid-${name}`, creationTimestamp: '2026-03-15T10:00:00Z' },
  spec: { title: `Todo ${name}`, status },
  status: {},
});

describe('TodoPage — E2E flow', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders list, creates, edits, and deletes a todo', async () => {
    const user = userEvent.setup();

    // Initial load: empty list
    mockedApi.listTodos.mockResolvedValueOnce({ items: [], metadata: {} });

    render(<TodoPage />);

    // Wait for loading to finish — empty state
    await waitFor(() => expect(screen.getByTestId('todo-empty')).toBeInTheDocument());

    // Click Create
    await user.click(screen.getByTestId('todo-create-btn'));

    // Fill in form
    mockedApi.createTodo.mockResolvedValue(makeTodo('new-todo'));
    mockedApi.listTodos.mockResolvedValueOnce({ items: [makeTodo('new-todo')], metadata: {} });

    await user.type(screen.getByTestId('todo-title-input'), 'New Task');
    await user.click(screen.getByTestId('todo-submit-btn'));

    // Should return to list with the new todo
    await waitFor(() => expect(screen.getByText('Todo new-todo')).toBeInTheDocument());

    // Edit the todo
    await user.click(screen.getByTestId('todo-edit-new-todo'));

    const titleInput = screen.getByTestId('todo-title-input');
    expect(titleInput).toHaveValue('Todo new-todo');

    const updatedTodo = makeTodo('new-todo');
    updatedTodo.spec.title = 'Updated Task';
    mockedApi.updateTodo.mockResolvedValue(updatedTodo);
    mockedApi.listTodos.mockResolvedValueOnce({ items: [updatedTodo], metadata: {} });

    await user.clear(titleInput);
    await user.type(titleInput, 'Updated Task');
    await user.click(screen.getByTestId('todo-submit-btn'));

    await waitFor(() => expect(screen.getByText('Updated Task')).toBeInTheDocument());

    // Delete the todo
    mockedApi.deleteTodo.mockResolvedValue(undefined);
    mockedApi.listTodos.mockResolvedValueOnce({ items: [], metadata: {} });

    await user.click(screen.getByTestId('todo-delete-new-todo'));
    // Confirm deletion in the dialog
    const dialog = screen.getByRole('dialog');
    await user.click(within(dialog).getByText('Delete'));

    await waitFor(() => expect(screen.getByTestId('todo-empty')).toBeInTheDocument());
  });

  it('shows error alert when API fails', async () => {
    mockedApi.listTodos.mockRejectedValueOnce(new Error('Network failure'));

    render(<TodoPage />);

    await waitFor(() => {
      expect(screen.getByTestId('todo-error')).toBeInTheDocument();
    }, { timeout: 3000 });
    expect(screen.getByText('Network failure')).toBeInTheDocument();
  });
});
