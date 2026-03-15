import React, { useState } from 'react';
import { Button, Field, Input, Select, TextArea } from '@grafana/ui';

import type { Spec } from '../generated/todo/v1/types.spec.gen';
import { PRIORITY_OPTIONS, getPriority } from './priorityUtils';
import type { Priority } from './priorityUtils';

const STATUS_OPTIONS = [
  { label: 'Open', value: 'open' as const },
  { label: 'In Progress', value: 'in_progress' as const },
  { label: 'Done', value: 'done' as const },
];

interface TodoFormProps {
  initialValues?: Spec;
  onSubmit: (spec: Spec) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
}

export function TodoForm({ initialValues, onSubmit, onCancel, isSubmitting = false }: TodoFormProps): React.ReactElement {
  const [title, setTitle] = useState(initialValues?.title ?? '');
  const [description, setDescription] = useState(initialValues?.description ?? '');
  const [status, setStatus] = useState<Spec['status']>(initialValues?.status ?? 'open');
  const [priority, setPriority] = useState<Priority>(initialValues ? getPriority(initialValues) : 'medium');
  const [titleError, setTitleError] = useState<string | undefined>(undefined);

  const handleSubmit = (e: React.FormEvent): void => {
    e.preventDefault();

    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setTitleError('Title is required');
      return;
    }

    setTitleError(undefined);
    onSubmit({
      title: trimmedTitle,
      description: description.trim() || undefined,
      status,
      priority,
    });
  };

  return (
    <form onSubmit={handleSubmit} data-testid="todo-form">
      <Field label="Title" required invalid={!!titleError} error={titleError}>
        <Input
          data-testid="todo-title-input"
          value={title}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setTitle(e.target.value);
            if (titleError) {
              setTitleError(undefined);
            }
          }}
          placeholder="Enter todo title"
          disabled={isSubmitting}
        />
      </Field>

      <Field label="Description">
        <TextArea
          data-testid="todo-description-input"
          value={description}
          onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setDescription(e.target.value)}
          placeholder="Enter description (optional)"
          rows={3}
          disabled={isSubmitting}
        />
      </Field>

      <Field label="Priority">
        <Select
          data-testid="todo-priority-select"
          options={[...PRIORITY_OPTIONS]}
          value={priority}
          onChange={(v) => {
            if (v.value) {
              setPriority(v.value as Priority);
            }
          }}
          disabled={isSubmitting}
        />
      </Field>

      <Field label="Status">
        <Select
          data-testid="todo-status-select"
          options={STATUS_OPTIONS}
          value={status}
          onChange={(v) => {
            if (v.value) {
              setStatus(v.value);
            }
          }}
          disabled={isSubmitting}
        />
      </Field>

      <div style={{ display: 'flex', gap: '8px', marginTop: '16px' }}>
        <Button type="submit" disabled={isSubmitting} data-testid="todo-submit-btn">
          {initialValues ? 'Update' : 'Create'}
        </Button>
        <Button variant="secondary" onClick={onCancel} disabled={isSubmitting} data-testid="todo-cancel-btn">
          Cancel
        </Button>
      </div>
    </form>
  );
}
