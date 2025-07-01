import { styled, FormLabel } from '@mui/material';

export const StyledFormLabel = styled(FormLabel)(({ theme }) => ({
  display: 'block',
  color: theme.vars.palette.text.primary,
  fontSize: 14,
  fontWeight: 400,
  marginBottom: theme.spacing(1),
  [theme.breakpoints.down('sm')]: {
    fontSize: 14,
  },
}));
