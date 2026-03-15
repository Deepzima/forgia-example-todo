import React from 'react';

const mockBackendSrv = {
  get: jest.fn(),
  post: jest.fn(),
  put: jest.fn(),
  delete: jest.fn(),
};

export const getBackendSrv = jest.fn(() => mockBackendSrv);

export function PluginPage({ children }: { children: React.ReactNode }): React.ReactElement {
  return React.createElement('div', null, children);
}
