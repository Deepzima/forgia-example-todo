import React from 'react';

// Minimal mocks for @grafana/ui components used in the plugin

export function Button({
  children,
  onClick,
  disabled,
  type,
  variant,
  size,
  ...rest
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: string;
  size?: string;
}): React.ReactElement {
  return (
    <button onClick={onClick} disabled={disabled} type={type ?? 'button'} {...rest}>
      {children}
    </button>
  );
}

export function Input({
  value,
  onChange,
  placeholder,
  disabled,
  ...rest
}: React.InputHTMLAttributes<HTMLInputElement>): React.ReactElement {
  return (
    <input value={value} onChange={onChange} placeholder={placeholder} disabled={disabled} {...rest} />
  );
}

export function TextArea({
  value,
  onChange,
  placeholder,
  rows,
  disabled,
  ...rest
}: React.TextareaHTMLAttributes<HTMLTextAreaElement>): React.ReactElement {
  return (
    <textarea value={value} onChange={onChange} placeholder={placeholder} rows={rows} disabled={disabled} {...rest} />
  );
}

export function Select({
  options,
  value,
  onChange,
  disabled,
  ...rest
}: {
  options: Array<{ label: string; value: string }>;
  value: string;
  onChange: (v: { value?: string }) => void;
  disabled?: boolean;
  'data-testid'?: string;
}): React.ReactElement {
  return (
    <select
      data-testid={rest['data-testid']}
      value={value}
      onChange={(e) => onChange({ value: e.target.value })}
      disabled={disabled}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  );
}

export function Field({
  label,
  children,
  required,
  invalid,
  error,
}: {
  label: string;
  children: React.ReactNode;
  required?: boolean;
  invalid?: boolean;
  error?: string;
}): React.ReactElement {
  return (
    <div>
      <label>
        {label}
        {required && ' *'}
      </label>
      {children}
      {invalid && error && <span>{error}</span>}
    </div>
  );
}

export function Card({
  children,
  ...rest
}: {
  children: React.ReactNode;
  'data-testid'?: string;
}): React.ReactElement {
  return <div data-testid={rest['data-testid']}>{children}</div>;
}

Card.Heading = function Heading({ children }: { children: React.ReactNode }): React.ReactElement {
  return <h4>{children}</h4>;
};

Card.Meta = function Meta({ children }: { children: string[] }): React.ReactElement {
  return (
    <div>
      {children.map((c, i) => (
        <span key={i}>{c}</span>
      ))}
    </div>
  );
};

Card.Description = function Description({ children }: { children: React.ReactNode }): React.ReactElement {
  return <p>{children}</p>;
};

Card.Tags = function Tags({ children }: { children: React.ReactNode }): React.ReactElement {
  return <div>{children}</div>;
};

Card.Actions = function Actions({ children }: { children: React.ReactNode }): React.ReactElement {
  return <div>{children}</div>;
};

export function Badge({ text, color, ...rest }: { text: string; color: string; 'data-testid'?: string }): React.ReactElement {
  return <span data-color={color} data-testid={rest['data-testid']}>{text}</span>;
}

export function RadioButtonGroup({
  options,
  value,
  onChange,
}: {
  options: Array<{ label: string; value: string }>;
  value: string;
  onChange: (value: string) => void;
}): React.ReactElement {
  return (
    <div role="radiogroup">
      {options.map((opt) => (
        <label key={opt.value}>
          <input
            type="radio"
            name="radio-group"
            value={opt.value}
            checked={value === opt.value}
            onChange={() => onChange(opt.value)}
          />
          {opt.label}
        </label>
      ))}
    </div>
  );
}

export function Alert({
  severity,
  title,
  children,
  ...rest
}: {
  severity: string;
  title: string;
  children?: React.ReactNode;
  'data-testid'?: string;
}): React.ReactElement {
  return (
    <div role="alert" data-testid={rest['data-testid']}>
      <strong>{title}</strong>
      {children}
    </div>
  );
}

export function Spinner({ ...rest }: { 'data-testid'?: string }): React.ReactElement {
  return <div data-testid={rest['data-testid']}>Loading...</div>;
}

export function ConfirmModal({
  isOpen,
  title,
  body,
  confirmText,
  onConfirm,
  onDismiss,
}: {
  isOpen: boolean;
  title: string;
  body: string;
  confirmText: string;
  onConfirm: () => void;
  onDismiss: () => void;
}): React.ReactElement | null {
  if (!isOpen) {
    return null;
  }
  return (
    <div role="dialog">
      <h3>{title}</h3>
      <p>{body}</p>
      <button onClick={onConfirm}>{confirmText}</button>
      <button onClick={onDismiss}>Cancel</button>
    </div>
  );
}

export function Stack({
  children,
  direction,
  gap,
}: {
  children: React.ReactNode;
  direction?: string;
  gap?: number;
}): React.ReactElement {
  return <div>{children}</div>;
}
