import React, { useState, useEffect } from 'react';
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
  Link,
  Chip,
  Stack,
} from '@mui/material';
import CheckCircleRoundedIcon from '@mui/icons-material/CheckCircleRounded';
import { useSearchParams } from 'react-router-dom';
import { postLogin } from '@/api/User';
import Lock from '@mui/icons-material/Lock';
import Person from '@mui/icons-material/Person';
import { useForm, Controller } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { styled } from '@mui/material/styles';
import {
  generateVSCodeAuthToken,
  openVSCodeWithAuth,
  checkVSCodeInstalled,
  getVSCodeDownloadUrl,
  storeAuthToken,
  generatePluginConfig,
  type VSCodeAuthToken,
} from '@/utils/vscode';

// 样式化组件
const StyledContainer = styled(Container)(({ theme }) => ({
  minHeight: '100vh',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  maxWidth: '100% !important',
  background: theme.palette.background.paper,
}));

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  maxWidth: 450,
  width: '100%',
  borderRadius: theme.spacing(2),
  boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
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

const AuthPage = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [authStep, setAuthStep] = useState<'login' | 'authorize' | 'success'>(
    'login'
  );
  const [vscodeInstalled, setVscodeInstalled] = useState<boolean | null>(null);
  const [authToken, setAuthToken] = useState<VSCodeAuthToken | null>(null);
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>();

  // 检查 VSCode 是否已安装
  useEffect(() => {
    // checkVSCodeInstalled().then(setVscodeInstalled);
  }, []);

  // 打开 VSCode 并处理授权
  const handleOpenVSCode = (token: VSCodeAuthToken) => {
    try {
      openVSCodeWithAuth(token);

      // 生成插件配置
      const pluginConfig = generatePluginConfig(token);

      setSuccess(`
        授权成功！
        
        令牌: ${token.token}
        用户ID: ${token.userId}
        权限: ${token.permissions.join(', ')}
        
        插件配置已生成，VSCode 插件将自动获取这些信息。
        如果 VSCode 没有自动打开，请手动打开 VSCode。
      `);

      console.log('插件配置:', pluginConfig);
    } catch (error) {
      setError(error instanceof Error ? error.message : '打开 VSCode 失败');
    }
  };

  // 处理登录表单提交
  const onSubmit = async (data: LoginFormData) => {
    setLoading(true);
    setError(null);

    try {
      // 步骤1: 用户登录
      const loginResult = await postLogin({
        ...data,
        session_id: searchParams.get('session_id')!,
      });

      // setAuthStep('authorize');

      // 步骤2: 生成 VSCode 授权令牌
      // const vscodeAuth = generateVSCodeAuthToken(loginResult.id!, [
      //   'read',
      //   'write',
      //   'execute',
      // ]);

      // // 存储 VSCode 授权信息
      // storeAuthToken(vscodeAuth);
      // setAuthToken(vscodeAuth);

      setAuthStep('success');

      // 步骤3: 打开 VSCode
      // handleOpenVSCode(vscodeAuth);

      try {
        // 尝试打开 VSCode
        window.location.href = loginResult.redirect_url!;

        // 记录授权日志
        console.log('VSCode 授权信息已发送:', {
          uri: loginResult.redirect_url,
        });
      } catch (error) {
        console.error('打开 VSCode 失败:', error);
        setError('无法打开 VSCode，请确保已安装 VSCode');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败，请重试');
      setAuthStep('login');
    } finally {
      setLoading(false);
    }
  };

  // 重新授权
  const handleReauthorize = async () => {
    if (!authToken) return;

    setLoading(true);
    setError(null);

    try {
      const newToken = generateVSCodeAuthToken(
        authToken.userId,
        authToken.permissions
      );
      storeAuthToken(newToken);
      setAuthToken(newToken);
      handleOpenVSCode(newToken);
    } catch (err) {
      setError('重新授权失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  // 继续到主应用
  const handleContinue = () => {
    navigate('/dashboard');
  };

  // 渲染登录表单
  const renderLoginForm = () => (
    <>
      <LogoContainer>
        <Typography
          variant='h4'
          component='h1'
          gutterBottom
          fontWeight='bold'
          color='primary'
        >
          Monkey Code
        </Typography>
        <Typography variant='body1' color='textSecondary' paragraph>
          登录您的账户以授权 VSCode 插件
        </Typography>

        {/* {vscodeInstalled !== null && (
          <Box sx={{ mb: 2 }}>
            <Chip
              label={
                vscodeInstalled ? '✅ 已检测到 VSCode' : '⚠️ 未检测到 VSCode'
              }
              color={vscodeInstalled ? 'success' : 'warning'}
              variant='outlined'
              size='small'
            />
          </Box>
        )} */}
      </LogoContainer>

      <Box component='form' onSubmit={handleSubmit(onSubmit)}>
        <Grid container spacing={2}>
          <Grid size={12}>
            <Controller
              name='username'
              control={control}
              defaultValue=''
              rules={{ required: '请输入用户名' }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label='用户名'
                  placeholder='请输入用户名'
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
              rules={{ required: '请输入密码' }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  label='密码'
                  placeholder='请输入密码'
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
              <Alert severity='error'>{error}</Alert>
            </Grid>
          )}

          <Grid size={12}>
            <Button
              type='submit'
              fullWidth
              variant='contained'
              size='large'
              disabled={loading}
              sx={{ mt: 2, mb: 2, height: 48 }}
            >
              {loading ? <CircularProgress size={24} /> : '登录并授权 VSCode'}
            </Button>
          </Grid>

          {!vscodeInstalled && vscodeInstalled !== null && (
            <Grid size={12}>
              <Alert severity='info' sx={{ mt: 1 }}>
                <Typography variant='body2'>
                  未检测到 VSCode，建议先{' '}
                  <Link
                    href={getVSCodeDownloadUrl()}
                    target='_blank'
                    rel='noopener noreferrer'
                  >
                    下载并安装 Visual Studio Code
                  </Link>
                </Typography>
              </Alert>
            </Grid>
          )}
        </Grid>
      </Box>
    </>
  );

  // 渲染授权进度
  const renderAuthorizationStep = () => (
    <LogoContainer>
      <CircularProgress size={60} sx={{ mb: 2 }} />
      <Typography variant='h5' gutterBottom>
        正在授权 VSCode 插件
      </Typography>
      <Typography variant='body1' color='textSecondary'>
        请稍候，正在生成授权令牌并尝试打开 VSCode...
      </Typography>
    </LogoContainer>
  );

  // 渲染成功页面
  const renderSuccessStep = () => (
    <>
      <LogoContainer>
        <Stack
          direction='row'
          alignItems='center'
          justifyContent='center'
          gap={1}
          sx={{ fontSize: 40, color: 'success.main', mb: 2 }}
        >
          <CheckCircleRoundedIcon sx={{ fontSize: 40 }} />
          授权成功！
        </Stack>
        <Typography variant='body1' color='textSecondary' paragraph>
          VSCode 授权已完成，插件可以正常使用了。
        </Typography>
      </LogoContainer>

      {/* {success && (
        <Alert
          severity='success'
          sx={{ mb: 2, whiteSpace: 'pre-line', fontSize: '0.875rem' }}
        >
          {success}
        </Alert>
      )} */}

      {/* {authToken && (
        <Box sx={{ mb: 2 }}>
          <Typography variant='subtitle2' gutterBottom>
            授权信息:
          </Typography>
          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
            {authToken.permissions.map((permission) => (
              <Chip
                key={permission}
                label={permission}
                size='small'
                color='primary'
                variant='outlined'
              />
            ))}
          </Box>
        </Box>
      )} */}

      <Grid container spacing={2}>
        {/* <Grid size={12}>
          <Button
            fullWidth
            variant='contained'
            size='large'
            onClick={handleContinue}
            sx={{ mb: 1 }}
          >
            进入主应用
          </Button>
        </Grid> */}

        {/* <Grid size={12}>
          <Button
            fullWidth
            variant='outlined'
            size='large'
            onClick={handleReauthorize}
            disabled={loading}
          >
            {loading ? <CircularProgress size={20} /> : '重新授权'}
          </Button>
        </Grid> */}

        <Grid size={12}>
          <Box textAlign='center'>
            <Button
              href={getVSCodeDownloadUrl()}
              target='_blank'
              color='info'
              rel='noopener noreferrer'
            >
              下载 Visual Studio Code
            </Button>
          </Box>
        </Grid>
      </Grid>
    </>
  );

  return (
    <StyledContainer>
      <StyledPaper elevation={3}>
        {authStep === 'login' && renderLoginForm()}
        {authStep === 'authorize' && renderAuthorizationStep()}
        {authStep === 'success' && renderSuccessStep()}
      </StyledPaper>
    </StyledContainer>
  );
};

export default AuthPage;
