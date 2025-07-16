import React from 'react';
import LoginHistory from './loginHistory';
import AdminTable from './adminTable';
import { Grid2 as Grid, Stack } from '@mui/material';

const Admin = () => {
  return (
    <Grid container  spacing={2} sx={{ height: '100%' }}>
        <Grid size={6}>
          <AdminTable />
        </Grid>
        <Grid size={6}>
          <LoginHistory />
        </Grid>
    </Grid>
  );
};

export default Admin;
