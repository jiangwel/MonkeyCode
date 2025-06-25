import KeyboardArrowRightRoundedIcon from '@mui/icons-material/KeyboardArrowRightRounded';
import { Box, Stack, useTheme } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';

const HomeBread = { title: '工作台', to: '/' };
const OtherBread = {
  dashboard: { title: '仪表盘', to: '/' },
  chat: { title: '对话记录', to: '/chat' },
  completion: { title: '补全记录', to: '/completion' },
  model: { title: '模型管理', to: '/model' },
  user: { title: '用户管理', to: '/user' },
  admin: { title: '管理员', to: '/admin' },
};

const Bread = () => {
  const theme = useTheme();
  const { pathname } = useLocation();
  const [breads, setBreads] = useState<
    { title: React.ReactNode; to: string }[]
  >([]);

  useEffect(() => {
    const curBreads: { title: React.ReactNode; to: string }[] = [
      {
        title: (
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
            }}
          >
            MonkeyCode
          </Box>
        ),
        to: '/dashboard',
      },
    ];
    if (pathname === '/') {
      curBreads.push(HomeBread);
    } else {
      const pieces = pathname.split('/').filter((it) => it !== '');
      pieces.forEach((it) => {
        const bread = OtherBread[it as keyof typeof OtherBread];
        if (bread) {
          curBreads.push(bread);
        }
      });
    }
    // if (pageName) {
    //   curBreads.push({ title: pageName, to: 'custom' })
    // }
    setBreads(curBreads);
  }, [pathname]);

  return (
    <Stack
      direction={'row'}
      alignItems={'center'}
      gap={1}
      sx={{
        flexGrow: 1,
        color: 'text.auxiliary',
        fontSize: '14px',
        a: { color: 'text.auxiliary' },
        lineHeight: '22px',
      }}
    >
      {/* <KBSelect /> */}
      {breads.map((it, idx) => {
        return (
          <Stack
            direction={'row'}
            alignItems={'center'}
            gap={1}
            key={idx}
            sx={{
              color:
                idx === breads.length - 1
                  ? `${theme.palette.text.primary} !important`
                  : 'text.disabled',

              ...(idx === breads.length - 1 && { fontWeight: 'bold' }),
            }}
          >
            {idx !== 0 && (
              <KeyboardArrowRightRoundedIcon sx={{ fontSize: 14 }} />
            )}
            {it.to === 'custom' ? (
              <Box
                sx={{ cursor: 'pointer', ':hover': { color: 'primary.main' } }}
              >
                {it.title}
              </Box>
            ) : (
              <NavLink to={it.to}>
                <Box
                  sx={{
                    cursor: 'pointer',
                    ':hover': { color: 'primary.main' },
                  }}
                >
                  {it.title}
                </Box>
              </NavLink>
            )}
          </Stack>
        );
      })}
    </Stack>
  );
};

export default Bread;
