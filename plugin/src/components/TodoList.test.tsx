import React from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { TodoList } from './TodoList';
import type { Todo } from '../generated/todo/v1/todo_object_gen';

const makeTodo = (overrides: Partial<Todo> = {}): Todo => ({
  kind: 'Todo',
  apiVersion: 'todo.grafana.app/v1',
  metadata: {
    name: 'test-todo',
    namespace: 'default',
    uid: 'uid-1',
    creationTimestamp: '2026-03-15T10:00:00Z',
  },
  spec: {
    title: 'Test Todo',
    description: 'A test todo',
    status: 'open',
  },
  status: {},
  ...overrides,
});

describe('TodoList', () => {
  const mockOnEdit = jest.fn();
  const mockOnDelete = jest.fn();
  const mockOnStatusChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('shows empty message when no todos', () => {
    render(
      <TodoList todos={[]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    expect(screen.getByTestId('todo-empty')).toHaveTextContent('No todos found');
  });

  it('renders list of todos with title, status, and creation date', () => {
    const todos = [
      makeTodo(),
      makeTodo({
        metadata: { name: 'todo-2', namespace: 'default', uid: 'uid-2', creationTimestamp: '2026-03-14T10:00:00Z' },
        spec: { title: 'Second Todo', status: 'done' },
      }),
    ];

    render(
      <TodoList todos={todos} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    expect(screen.getByText('Test Todo')).toBeInTheDocument();
    expect(screen.getByText('Second Todo')).toBeInTheDocument();
    expect(screen.getByText('Open')).toBeInTheDocument();
    expect(screen.getByText('Done')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', async () => {
    const user = userEvent.setup();
    const todo = makeTodo();

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    await user.click(screen.getByTestId('todo-edit-test-todo'));
    expect(mockOnEdit).toHaveBeenCalledWith(todo);
  });

  it('shows delete confirmation modal and calls onDelete on confirm', async () => {
    const user = userEvent.setup();
    const todo = makeTodo();

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    await user.click(screen.getByTestId('todo-delete-test-todo'));
    // Confirm modal should appear
    expect(screen.getByText('Are you sure you want to delete this TODO? This action cannot be undone.')).toBeInTheDocument();

    const dialog = screen.getByRole('dialog');
    const confirmBtn = within(dialog).getByText('Delete');
    await user.click(confirmBtn);
    expect(mockOnDelete).toHaveBeenCalledWith('test-todo');
  });

  it('calls onStatusChange with next status when advance button clicked', async () => {
    const user = userEvent.setup();
    const todo = makeTodo({ spec: { title: 'Open Todo', status: 'open' } });

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    await user.click(screen.getByTestId('todo-advance-test-todo'));
    expect(mockOnStatusChange).toHaveBeenCalledWith(todo, 'in_progress');
  });

  it('shows Complete button for in_progress todos', () => {
    const todo = makeTodo({ spec: { title: 'IP Todo', status: 'in_progress' } });

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    expect(screen.getByText('Complete')).toBeInTheDocument();
  });

  it('does not show advance button for done todos', () => {
    const todo = makeTodo({ spec: { title: 'Done Todo', status: 'done' } });

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    expect(screen.queryByTestId('todo-advance-test-todo')).not.toBeInTheDocument();
  });

  it.each([
    ['low', 'blue', 'Low'],
    ['medium', 'yellow', 'Medium'],
    ['high', 'orange', 'High'],
    ['critical', 'red', 'Critical'],
  ] as const)('renders %s priority badge with color=%s', (priority, expectedColor, expectedLabel) => {
    const todo = makeTodo({
      spec: { title: `${priority} todo`, status: 'open', priority },
    });

    render(
      <TodoList todos={[todo]} onEdit={mockOnEdit} onDelete={mockOnDelete} onStatusChange={mockOnStatusChange} />
    );

    const badge = screen.getByTestId('todo-priority-badge-test-todo');
    expect(badge).toHaveTextContent(expectedLabel);
    expect(badge).toHaveAttribute('data-color', expectedColor);
  });
});
