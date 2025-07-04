import { Stack, Link, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import Avatar from '../avatar';

const User = ({
  id,
  username = '',
  email = '',
}: {
  id?: string;
  username?: string;
  email?: string;
}) => {
  return (
    <Stack>
      <Link
        component={RouterLink}
        to={id ? `/dashboard/member/${id}` : ''}
        sx={{
          '&:hover': {
            color: id ? 'info.main' : 'text.primary',
          },
          cursor: 'pointer',
        }}
      >
        <Stack direction={'row'} alignItems={'center'} gap={1}>
          <Avatar
            name={username}
            sx={{ width: 22, height: 22, fontSize: 14 }}
          />
          <Typography>{username}</Typography>
        </Stack>
      </Link>
      {email && (
        <Typography color='text.auxiliary' sx={{ fontSize: 14 }}>
          {email}
        </Typography>
      )}
    </Stack>
  );
};

export default User;
