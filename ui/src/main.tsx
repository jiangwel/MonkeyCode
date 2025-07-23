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
import { getMyModelList } from '@/api/Model';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import { AuthContext, CommonContext } from './context';
import { DomainUser, DomainAdminUser, DomainModel } from './api/types';
import { lightTheme } from './theme';
import router from './router';
import { getRedirectUrl } from './utils';

dayjs.locale('zh-cn');
dayjs.extend(duration);
dayjs.extend(relativeTime);

const App = () => {
  const [user, setUser] = useState<DomainUser | DomainAdminUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [coderModel, setCoderModel] = useState<DomainModel[]>([]);
  const [llmModel, setLlmModel] = useState<DomainModel[]>([]);
  const [isConfigModel, setIsConfigModel] = useState(false);

  const onGotoRedirect = (source: 'user' | 'admin') => {
    const redirectUrl = getRedirectUrl(source);
    window.location.href = redirectUrl.href;
  };

  const getModelList = () => {
    return Promise.all([
      getMyModelList({
        model_type: 'coder',
      }),
      getMyModelList({
        model_type: 'llm',
      }),
    ])
      .then((res) => {
        setCoderModel(res[0] || []);
        setLlmModel(res[1] || []);
        return res;
      })
      .then(handleModelConfig);
  };

  const handleModelConfig = (res: [DomainModel[], DomainModel[]]) => {
    if ((res[0] || [])?.length == 0 || (res[1] || [])?.length == 0) {
      if (location.pathname !== '/model') {
        window.location.href = '/model';
      }
      setIsConfigModel(false);
      return false;
    } else {
      const isActive =
        res[0].every((item) => item.is_active) &&
        res[1].every((item) => item.is_active);
      if (isActive) {
        setIsConfigModel(true);
        return true;
      } else {
        if (location.pathname !== '/model') {
          window.location.href = '/model';
        }
        setIsConfigModel(false);
        return false;
      }
    }
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
    } else if (!location.pathname.startsWith('/auth')) {
      return getAdminProfile()
        .then((res) => {
          setUser(res);
          getModelList().then((res) => {
            if (res) {
              if (location.pathname.startsWith('/login')) {
                onGotoRedirect('admin');
              }
            }
          });
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      setLoading(false);
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
      <CommonContext.Provider
        value={{
          coderModel,
          llmModel,
          isConfigModel,
          refreshModel: getModelList,
        }}
      >
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
      </CommonContext.Provider>
    </ThemeProvider>
  );
};

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(<App />);
