import ReactDOM from 'react-dom/client';
import { useState, useEffect } from 'react';
import 'dayjs/locale/zh-cn';
import { RouterProvider } from 'react-router-dom';
import '@/assets/fonts/font.css';
import '@/assets/fonts/iconfont';
import './index.css';
import '@/assets/styles/markdown.css';
import { ThemeProvider } from '@c-x/ui';
import { getUserProfile } from '@/api/UserManage';
import { getAdminProfile } from '@/api/Admin';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import { AuthContext } from './context';
import { DomainUser, DomainAdminUser } from './api/types';
import { lightTheme } from './theme';
import router from './router';
import { getRedirectUrl } from './utils';
import { Loading } from '@c-x/ui';

dayjs.locale('zh-cn');
dayjs.extend(duration);
dayjs.extend(relativeTime);

const App = () => {
  const [user, setUser] = useState<DomainUser | DomainAdminUser | null>(null);
  const [loading, setLoading] = useState(true);

  const onGotoRedirect = (source: 'user' | 'admin') => {
    const redirectUrl = getRedirectUrl(source);
    window.location.href = redirectUrl.href;
  };

  const getUser = () => {
    setLoading(true);
    if (location.pathname.startsWith('/user')) {
      return getUserProfile()
        .then((res) => {
          setUser(res);
          if (location.pathname.startsWith('/user/login')) {
            onGotoRedirect('user');
          }
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      return getAdminProfile()
        .then((res) => {
          console.log(res);
          setUser(res);
          if (location.pathname.startsWith('/login')) {
            onGotoRedirect('admin');
          }
        })
        .finally(() => {
          setLoading(false);
        });
    }
  };

  useEffect(() => {
    getUser();
  }, []);

  if (loading) {
    return null;
  }

  return (
    <ThemeProvider theme={lightTheme}>
      <AuthContext.Provider
        value={[
          user,
          {
            loading,
            setUser,
            refreshUser: getUser,
          },
        ]}
      >
        <RouterProvider router={router} />
      </AuthContext.Provider>
    </ThemeProvider>
  );
};

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(<App />);
