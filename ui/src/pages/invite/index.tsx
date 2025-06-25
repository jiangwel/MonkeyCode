import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Stepper,
  Step,
  StepLabel,
  TextField,
  Button,
  Stack,
  Typography,
  Container,
  CardContent,
} from '@mui/material';
import { useRequest } from 'ahooks';
import { postRegister } from '@/api/User';

import Card from '@/components/card';
import DownloadIcon from '@mui/icons-material/Download';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import { Controller, useForm } from 'react-hook-form';
import { StyledFormLabel } from '@/components/form';

const Invite = () => {
  const { id } = useParams();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      email: '',
      password: '',
    },
  });
  const { run: register, loading } = useRequest(postRegister, {
    manual: true,
  });
  const [activeStep, setActiveStep] = useState(0);

  const onNext = () => {
    setActiveStep(activeStep + 1);
  };

  const onRegister = handleSubmit((data) => {
    register({ ...data, code: id });
    onNext();
  });

  return (
    <Box sx={{ bgcolor: 'background.paper', height: '100vh' }}>
      <Container maxWidth='md' sx={{ py: 4 }}>
        <Card sx={{ p: 4 }}>
          <Typography variant='h4' component='h1' align='center' gutterBottom>
            欢迎加入
          </Typography>
          <Typography
            variant='body1'
            color='text.secondary'
            align='center'
            sx={{ mb: 4 }}
          >
            请完成以下步骤开始使用我们的服务
          </Typography>

          <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 4 }}>
            <Step>
              <StepLabel>注册账号</StepLabel>
            </Step>
            <Step>
              <StepLabel>下载客户端</StepLabel>
            </Step>
            <Step>
              <StepLabel>使用教程</StepLabel>
            </Step>
          </Stepper>

          {activeStep === 0 && (
            <Stack gap={3}>
              <Box>
                <StyledFormLabel required>邮箱</StyledFormLabel>
                <Controller
                  control={control}
                  rules={{
                    required: '请输入邮箱',
                    pattern: {
                      value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                      message: '请输入正确的邮箱地址',
                    },
                  }}
                  name='email'
                  render={({ field }) => (
                    <TextField
                      {...field}
                      size='small'
                      fullWidth
                      placeholder='请输入您的邮箱地址'
                      error={!!errors.email}
                      helperText={errors.email?.message as string}
                    />
                  )}
                />
              </Box>
              <Box>
                <StyledFormLabel required>密码</StyledFormLabel>
                <Controller
                  control={control}
                  rules={{
                    required: '请输入密码',
                    minLength: {
                      value: 8,
                      message: '密码至少需要8位',
                    },
                  }}
                  name='password'
                  render={({ field }) => (
                    <TextField
                      {...field}
                      size='small'
                      fullWidth
                      type='password'
                      placeholder='请设置您的密码'
                      error={!!errors.password}
                      helperText={errors.password?.message as string}
                    />
                  )}
                />
              </Box>

              <Button
                onClick={onRegister}
                variant='contained'
                loading={loading}
                size='medium'
                sx={{
                  mt: 2,
                  px: 4,
                  py: 1,
                  fontSize: '1rem',
                  textTransform: 'none',
                  borderRadius: 1,
                  width: 'fit-content',
                  alignSelf: 'center',
                }}
              >
                立即注册
              </Button>
            </Stack>
          )}

          {activeStep === 1 && (
            <Stack gap={3} alignItems='center'>
              <Card sx={{ width: '100%', maxWidth: 600, textAlign: 'center' }}>
                <CardContent>
                  <DownloadIcon
                    sx={{ fontSize: 60, color: 'primary.main', mb: 2 }}
                  />
                  <Typography variant='h6' gutterBottom>
                    下载 MonkeyCode 客户端
                  </Typography>
                  <Typography
                    variant='body2'
                    color='text.secondary'
                    sx={{ mb: 3 }}
                  >
                    请下载并安装 MonkeyCode 客户端，这是使用我们服务的必要步骤
                  </Typography>
                  <Button
                    variant='contained'
                    startIcon={<DownloadIcon />}
                    href='/api/v1/static/vsix'
                    download='monkeycode-client.vsix'
                    rel='noopener noreferrer'
                    sx={{
                      textTransform: 'none',
                      px: 4,
                      py: 1,
                      borderRadius: 1,
                    }}
                    onClick={onNext}
                  >
                    下载客户端
                  </Button>
                </CardContent>
              </Card>
            </Stack>
          )}

          {activeStep === 2 && (
            <Stack gap={3} alignItems='center'>
              <Card sx={{ width: '100%', maxWidth: 600, textAlign: 'center' }}>
                <CardContent>
                  <MenuBookIcon
                    sx={{ fontSize: 60, color: 'primary.main', mb: 2 }}
                  />
                  <Typography variant='h6' gutterBottom>
                    使用教程
                  </Typography>
                  <Typography
                    variant='body2'
                    color='text.secondary'
                    sx={{ mb: 3 }}
                  >
                    查看详细的使用教程，快速上手 MonkeyCode 客户端
                  </Typography>
                  <Button
                    variant='outlined'
                    startIcon={<MenuBookIcon />}
                    href='https://www.monkeycode.cn/docs/client/quick-start'
                    target='_blank'
                    rel='noopener noreferrer'
                    sx={{
                      textTransform: 'none',
                      px: 4,
                      py: 1,
                      borderRadius: 1,
                    }}
                  >
                    查看教程
                  </Button>
                </CardContent>
              </Card>
            </Stack>
          )}
        </Card>
      </Container>
    </Box>
  );
};

export default Invite;
