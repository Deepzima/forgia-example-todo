// Grafana-compatible Jest config
module.exports = {
  testEnvironment: 'jsdom',
  testMatch: ['<rootDir>/src/**/*.test.{ts,tsx}'],
  transform: {
    '^.+\\.(t|j)sx?$': [
      '@swc/jest',
      {
        sourceMaps: true,
        jsc: {
          parser: {
            syntax: 'typescript',
            tsx: true,
            decorators: false,
            dynamicImport: true,
          },
          transform: {
            react: {
              runtime: 'automatic',
            },
          },
        },
      },
    ],
  },
  transformIgnorePatterns: [
    'node_modules/(?!(@grafana|ol|react-colorful|uuid|d3|d3-color|d3-interpolate|d3-force|d3-scale|rxjs))',
  ],
  moduleNameMapper: {
    '\\.(css|scss|sass)$': '<rootDir>/src/__mocks__/styleMock.ts',
    '\\.(svg|png|jpg|gif)$': '<rootDir>/src/__mocks__/fileMock.ts',
    '^@/(.*)$': '<rootDir>/src/$1',
    '^@grafana/ui$': '<rootDir>/src/__mocks__/@grafana/ui.tsx',
    '^@grafana/runtime$': '<rootDir>/src/__mocks__/@grafana/runtime.ts',
    '^@grafana/data$': '<rootDir>/src/__mocks__/@grafana/data.ts',
  },
  setupFiles: ['<rootDir>/src/setupTests.ts'],
  setupFilesAfterEnv: ['<rootDir>/src/setupAfterEnv.ts'],
  collectCoverageFrom: [
    'src/**/*.{ts,tsx}',
    '!src/**/*.test.{ts,tsx}',
    '!src/generated/**',
    '!src/__mocks__/**',
    '!src/module.ts',
    '!src/setupTests.ts',
  ],
  coverageThreshold: {
    global: {
      branches: 60,
      functions: 80,
      lines: 78,
      statements: 78,
    },
  },
};
