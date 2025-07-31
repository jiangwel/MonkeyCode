import React from 'react';
import BaseURLSettings from './baseURLSettings';
import { Grid2 as Grid, Stack } from '@mui/material';

const GeneralSetting = () => {
  return (
    <Grid container  spacing={2} sx={{ height: '100%' }}>
        <Grid size={6}>
          <BaseURLSettings />
        </Grid>
        <Grid size={6}>
        </Grid>
    </Grid>
  );
};

export default GeneralSetting;
