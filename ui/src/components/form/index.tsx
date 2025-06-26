'use client';
import { styled, FormLabel } from '@mui/material';

export const StyledFormLabel = styled(FormLabel)(({ theme }) => ({
  display: 'block',
  color: theme.vars.palette.text.primary,
  fontSize: 16,
  fontWeight: 500,
  marginBottom: theme.spacing(1),
  [theme.breakpoints.down('sm')]: {
    fontSize: 14,
  },
}));
