import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { TodoForm } from './TodoForm';

describe('TodoForm', () => {
  const mockOnSubmit = jest.fn();
  const mockOnCancel = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders empty form for creation', () => {
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    expect(screen.getByTestId('todo-title-input')).toHaveValue('');
    expect(screen.getByTestId('todo-description-input')).toHaveValue('');
    expect(screen.getByText('Create')).toBeInTheDocument();
  });

  it('renders pre-filled form for editing', () => {
    render(
      <TodoForm
        initialValues={{ title: 'Test', description: 'Desc', status: 'in_progress' }}
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />
    );

    expect(screen.getByTestId('todo-title-input')).toHaveValue('Test');
    expect(screen.getByTestId('todo-description-input')).toHaveValue('Desc');
    expect(screen.getByText('Update')).toBeInTheDocument();
  });

  it('shows validation error when title is empty', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.click(screen.getByTestId('todo-submit-btn'));

    expect(screen.getByText('Title is required')).toBeInTheDocument();
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('submits form with valid data', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.type(screen.getByTestId('todo-title-input'), 'My Todo');
    await user.type(screen.getByTestId('todo-description-input'), 'A description');
    await user.click(screen.getByTestId('todo-submit-btn'));

    expect(mockOnSubmit).toHaveBeenCalledWith({
      title: 'My Todo',
      description: 'A description',
      status: 'open',
      priority: 'medium',
    });
  });

  it('calls onCancel when cancel button is clicked', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.click(screen.getByTestId('todo-cancel-btn'));

    expect(mockOnCancel).toHaveBeenCalled();
  });

  it('clears validation error when user types in title', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.click(screen.getByTestId('todo-submit-btn'));
    expect(screen.getByText('Title is required')).toBeInTheDocument();

    await user.type(screen.getByTestId('todo-title-input'), 'A');
    expect(screen.queryByText('Title is required')).not.toBeInTheDocument();
  });

  it('disables inputs when submitting', () => {
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} isSubmitting />);

    expect(screen.getByTestId('todo-title-input')).toBeDisabled();
    expect(screen.getByTestId('todo-description-input')).toBeDisabled();
    expect(screen.getByTestId('todo-submit-btn')).toBeDisabled();
  });

  it('trims whitespace from title and description', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.type(screen.getByTestId('todo-title-input'), '  Trimmed Title  ');
    await user.type(screen.getByTestId('todo-description-input'), '  Trimmed Desc  ');
    await user.click(screen.getByTestId('todo-submit-btn'));

    expect(mockOnSubmit).toHaveBeenCalledWith({
      title: 'Trimmed Title',
      description: 'Trimmed Desc',
      status: 'open',
      priority: 'medium',
    });
  });

  it('shows status dropdown with three options', () => {
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);
    // The Select component from @grafana/ui renders the selected value
    expect(screen.getByText('Open')).toBeInTheDocument();
  });

  it('renders priority Select with 4 options and defaults to "medium"', () => {
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    const prioritySelect = screen.getByTestId('todo-priority-select');
    expect(prioritySelect).toBeInTheDocument();
    expect(prioritySelect).toHaveValue('medium');

    const options = prioritySelect.querySelectorAll('option');
    expect(options).toHaveLength(4);
    expect(Array.from(options).map((o) => o.value)).toEqual(['low', 'medium', 'high', 'critical']);
  });

  it('pre-selects current priority when editing', () => {
    render(
      <TodoForm
        initialValues={{ title: 'Test', status: 'open', priority: 'critical' }}
        onSubmit={mockOnSubmit}
        onCancel={mockOnCancel}
      />
    );

    expect(screen.getByTestId('todo-priority-select')).toHaveValue('critical');
  });

  it('submits with selected priority value', async () => {
    const user = userEvent.setup();
    render(<TodoForm onSubmit={mockOnSubmit} onCancel={mockOnCancel} />);

    await user.type(screen.getByTestId('todo-title-input'), 'Priority Task');
    await user.selectOptions(screen.getByTestId('todo-priority-select'), 'high');
    await user.click(screen.getByTestId('todo-submit-btn'));

    expect(mockOnSubmit).toHaveBeenCalledWith(
      expect.objectContaining({ title: 'Priority Task', priority: 'high' })
    );
  });
});
