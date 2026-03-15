import { AppPlugin } from '@grafana/data';

import { TodoPage } from './pages/TodoPage';

export const plugin = new AppPlugin().addPage({
  title: 'TODO List',
  body: TodoPage,
  id: 'todos',
});
