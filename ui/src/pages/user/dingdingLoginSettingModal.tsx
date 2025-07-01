import React, { useState, useEffect } from 'react';
import { Modal, message, Loading } from '@c-x/ui';
import { useForm, Controller } from 'react-hook-form';
import { StyledFormLabel } from '@/components/form';
import { putUpdateSetting } from '@/api/User';
import {
  Box,
  Typography,
  IconButton,
  Paper,
  TextField,
  Stack,
} from '@mui/material';

const DingingLoginSettingModal = ({
  open,
  onClose,
  onOk,
}: {
  open: boolean;
  onClose: () => void;
  onOk: () => void;
}) => {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    defaultValues: {
      dingtalk_client_id: '',
      dingtalk_client_secret: '',
      // title: '',
    },
  });

  useEffect(() => {
    if (open) {
      reset();
    }
  }, [open]);

  const onSubmit = handleSubmit((data) => {
    putUpdateSetting({ ...data, enable_dingtalk_oauth: true }).then(() => {
      message.success('设置成功');
      onClose();
      onOk();
    });
  });

  return (
    <Modal
      title='钉钉登录设置'
      width={800}
      open={open}
      onOk={onSubmit}
      onCancel={onClose}
    >
      <Stack gap={2}>
        <Box>
          <StyledFormLabel required>Client ID</StyledFormLabel>
          <Controller
            control={control}
            name='dingtalk_client_id'
            rules={{
              required: {
                value: true,
                message: 'Client Id 不能为空',
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                size='small'
                placeholder='请输入'
                error={!!errors.dingtalk_client_id}
                helperText={errors.dingtalk_client_id?.message}
              />
            )}
          />
        </Box>
        <Box>
          <StyledFormLabel required>Client Secret</StyledFormLabel>
          <Controller
            control={control}
            name='dingtalk_client_secret'
            rules={{
              required: {
                value: true,
                message: 'Client Secret 不能为空',
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                size='small'
                placeholder='请输入'
                error={!!errors.dingtalk_client_secret}
                helperText={errors.dingtalk_client_secret?.message}
              />
            )}
          />
        </Box>
        {/* <Box>
          <StyledFormLabel>标题名称，默认为 身份认证-钉钉登录</StyledFormLabel>
          <Controller
            control={control}
            name='title'
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                size='small'
                placeholder='请输入'
                error={!!errors.title}
                helperText={errors.title?.message}
                onChange={(e) => {
                  field.onChange(e.target.value);
                }}
              />
            )}
          />
        </Box> */}
      </Stack>
    </Modal>
  );
};

export default DingingLoginSettingModal;
