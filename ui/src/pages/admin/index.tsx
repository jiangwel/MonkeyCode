import React from 'react';
import LoginHistory from './loginHistory';
import AdminTable from './adminTable';
import { Stack } from '@mui/material';

const Admin = () => {
  return (
    <Stack gap={2} sx={{ height: '100%' }}>
      <AdminTable />
      <LoginHistory />
    </Stack>
  );
};

export default Admin;
