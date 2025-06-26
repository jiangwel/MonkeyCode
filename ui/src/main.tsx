import ReactDOM from 'react-dom/client';
import 'dayjs/locale/zh-cn';
import { RouterProvider } from 'react-router-dom';
import '@/assets/fonts/iconfont';
import './index.css';
import '@/assets/styles/markdown.css';
import { ThemeProvider } from '@c-x/ui';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';

import { lightTheme } from './theme';

import router from './router';

dayjs.locale('zh-cn');
dayjs.extend(duration);
dayjs.extend(relativeTime);

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <ThemeProvider theme={lightTheme}>
    <RouterProvider router={router} />
  </ThemeProvider>
);
