import React, { useState } from 'react';
import {
  Box,
  Button,
  TextField,
  Typography,
  Container,
  Paper,
  Alert,
  CircularProgress,
  Grid2 as Grid,
  Divider,
} from '@mui/material';
import { postAdminLogin } from '@/api/User';
import { useForm, Controller } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { styled } from '@mui/material/styles';
import { Lock, Person } from '@mui/icons-material';
import { getRedirectUrl } from '@/utils';

// 样式化组件
const StyledContainer = styled(Container)(({ theme }) => ({
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  maxWidth: '100% !important',
  backgroundColor: theme.palette.background.paper,
}));

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  maxWidth: 400,
  width: '100%',
  borderRadius: 12,
  boxShadow: '0px 10px 20px 0px rgba(54,59,76,0.2)',
  backgroundColor: theme.palette.background.default,
}));

const LogoContainer = styled(Box)(({ theme }) => ({
  textAlign: 'center',
  marginBottom: theme.spacing(3),
}));

// 表单数据类型
interface LoginFormData {
  username: string;
  password: string;
}

const LoginPage = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>();

  const loginUser = async (data: LoginFormData) => {
    return postAdminLogin({
      username: data.username,
      password: data.password,
    });
  };

  // 处理登录表单提交
  const onSubmit = async (data: LoginFormData) => {
    setLoading(true);
    setError(null);

    try {
      await loginUser(data);
      const redirectUrl = getRedirectUrl();
      window.location.href = redirectUrl.href;
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <StyledContainer>
      <StyledPaper elevation={0}>
        <LogoContainer>
          <Box
            sx={{
              width: 60,
              height: 60,
              borderRadius: '50%',
              backgroundColor: 'primary.main',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              margin: '0 auto 16px',
            }}
          >
            <Lock sx={{ color: 'white', fontSize: 28 }} />
          </Box>
          <Typography
            variant='h4'
            component='h1'
            gutterBottom
            fontWeight='bold'
            sx={{
              color: 'primary.main',
              mb: 1,
            }}
          >
            欢迎登录
          </Typography>
          <Typography variant='body1' color='text.secondary' paragraph>
            请输入您的账号和密码
          </Typography>
        </LogoContainer>

        <Divider sx={{ mb: 3 }} />

        <Box component='form' onSubmit={handleSubmit(onSubmit)}>
          <Grid container spacing={3}>
            <Grid size={12}>
              <Controller
                name='username'
                control={control}
                defaultValue=''
                rules={{
                  required: '请输入用户名',
                  minLength: {
                    value: 2,
                    message: '用户名至少需要2个字符',
                  },
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label='用户名'
                    variant='outlined'
                    error={!!errors.username}
                    helperText={errors.username?.message}
                    disabled={loading}
                    slotProps={{
                      input: {
                        startAdornment: (
                          <Person sx={{ color: 'text.secondary', mr: 1 }} />
                        ),
                      },
                    }}
                  />
                )}
              />
            </Grid>

            <Grid size={12}>
              <Controller
                name='password'
                control={control}
                defaultValue=''
                rules={{
                  required: '请输入密码',
                  minLength: {
                    value: 3,
                    message: '密码至少需要3个字符',
                  },
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label='密码'
                    type='password'
                    variant='outlined'
                    error={!!errors.password}
                    helperText={errors.password?.message}
                    disabled={loading}
                    slotProps={{
                      input: {
                        startAdornment: (
                          <Lock sx={{ color: 'text.secondary', mr: 1 }} />
                        ),
                      },
                    }}
                  />
                )}
              />
            </Grid>

            {error && (
              <Grid size={12}>
                <Alert severity='error' sx={{ borderRadius: 2 }}>
                  {error}
                </Alert>
              </Grid>
            )}

            <Grid size={12}>
              <Button
                type='submit'
                fullWidth
                variant='contained'
                size='large'
                disabled={loading}
                sx={{ mt: 1, mb: 2 }}
              >
                {loading ? (
                  <>
                    <CircularProgress
                      size={20}
                      sx={{ mr: 1, color: 'inherit' }}
                    />
                    登录中...
                  </>
                ) : (
                  '登录'
                )}
              </Button>
            </Grid>
          </Grid>
        </Box>
      </StyledPaper>
    </StyledContainer>
  );
};

export default LoginPage;
