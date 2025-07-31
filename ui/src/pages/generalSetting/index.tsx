import CardServiceSettings from './components/cardServiceSettings';
import { Grid2 as Grid } from '@mui/material';
import CardAdminUser from './components/cardAdminUser';

const GeneralSetting = () => {
  return (
    <Grid container  spacing={2} sx={{ height: '100%' }}>
        <Grid size={6}>
          <CardServiceSettings />
        </Grid>
        <Grid size={6}>
          <CardAdminUser />
        </Grid>
    </Grid>
  );
};

export default GeneralSetting;
