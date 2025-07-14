import { Button, IconButton, Stack } from '@mui/material';
import { Icon, message } from '@c-x/ui';
import { useNavigate, useLocation } from 'react-router-dom';
import LogoutIcon from '@mui/icons-material/Logout';
import DownloadIcon from '@mui/icons-material/Download';
import { Box } from '@mui/material';
import Bread from './Bread';

const Header = () => {
  const navigate = useNavigate();
  const { pathname } = useLocation();
  return (
    <Stack
      direction={'row'}
      alignItems={'center'}
      justifyContent={'space-between'}
      sx={{
        mt: 2,
        pr: 2,
        width: '100%',
      }}
    >
      <Bread />

      <Stack direction={'row'} alignItems={'center'} gap={2}>
        <Button
          variant='contained'
          size='small'
          startIcon={<DownloadIcon />}
          href='/api/v1/static/vsix'
          download='monkeycode-client.vsix'
          rel='noopener noreferrer'
        >
          下载客户端
        </Button>
        <IconButton
          title='退出登录'
          size='small'
          sx={{
            bgcolor: '#fff',

            '&:hover': {
              color: 'primary.main',
            },
          }}
          onClick={() => {
            message.success('退出登录成功');
            navigate(pathname.startsWith('/user/') ? '/user/login' : '/login', {
              replace: true,
            });
          }}
        >
          <LogoutIcon sx={{ fontSize: 16 }} />
        </IconButton>
      </Stack>
    </Stack>
  );
};

export default Header;
