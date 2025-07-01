import React, { useState } from 'react';
import Card from '@/components/card';
import {
  Grid2 as Grid,
  Stack,
  styled,
  Switch,
  Button,
  Box,
  Select,
  MenuItem,
  Radio,
} from '@mui/material';
import { Icon, Modal } from '@c-x/ui';
import { useRequest } from 'ahooks';
import { getGetSetting, putUpdateSetting } from '@/api/User';
import MemberManage from './memberManage';
import LoginHistory from './loginHistory';
import { message } from '@c-x/ui';
import DingingLoginSettingModal from './dingdingLoginSettingModal';

const StyledCard = styled(Card)({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  boxShadow:
    '0px 0px 10px 0px rgba(68, 80, 91, 0.1), 0px 0px 2px 0px rgba(68, 80, 91, 0.1)',
});

const StyledLabel = styled('div')(({ theme }) => ({
  fontWeight: 700,
  fontSize: 14,
  color: theme.vars.palette.text.primary,
}));

const User = () => {
  const [dingdingLoginSettingModalOpen, setDingdingLoginSettingModalOpen] =
    useState(false);
  const [dingdingCheck, setDingdingCheck] = useState(false);
  const {
    data = {
      enable_sso: false,
      force_two_factor_auth: false,
      disable_password_login: false,
      enable_dingtalk_oauth: false,
    },
    refresh,
  } = useRequest(getGetSetting, {
    onSuccess: (data) => {
      setDingdingCheck(data.enable_dingtalk_oauth!);
    },
  });

  const { runAsync: updateSetting } = useRequest(putUpdateSetting, {
    manual: true,
    onSuccess: () => {
      refresh();
      message.success('设置更新成功');
    },
  });

  const onDisabledDingdingLogin = () => {
    Modal.confirm({
      title: '提示',
      content: '确定要关闭钉钉登录吗？',
      onOk: () => {
        updateSetting({ enable_dingtalk_oauth: false }).then(() => {
          refresh();
        });
      },
    });
  };

  return (
    <Stack gap={2}>
      <Grid container spacing={2}>
        <Grid size={6} container>
          <Grid size={12}>
            <StyledCard>
              <StyledLabel>强制成员启用两步认证</StyledLabel>
              <Switch
                checked={data?.force_two_factor_auth}
                onChange={(e) => {
                  updateSetting({ force_two_factor_auth: e.target.checked });
                }}
              />
            </StyledCard>
          </Grid>
          <Grid size={12}>
            <StyledCard>
              <StyledLabel>禁止成员使用密码登录</StyledLabel>
              <Switch
                checked={data?.disable_password_login}
                onChange={(e) =>
                  updateSetting({ disable_password_login: e.target.checked })
                }
              />
            </StyledCard>
          </Grid>
        </Grid>
        <Grid size={6} container>
          <Grid size={12}>
            <Card sx={{ height: '100%' }}>
              <StyledLabel>第三方登录</StyledLabel>
              <Stack
                direction='row'
                alignItems='center'
                spacing={2}
                sx={{ mt: 2, height: 'calc(100% - 40px)' }}
              >
                <Button
                  variant='outlined'
                  color='primary'
                  sx={{ gap: 3 }}
                  onClick={() => {
                    if (dingdingCheck) {
                      onDisabledDingdingLogin();
                    } else {
                      setDingdingLoginSettingModalOpen(true);
                    }
                  }}
                >
                  <Radio size='small' sx={{ p: 0.5 }} checked={dingdingCheck} />
                  <Stack direction='row' alignItems='center' gap={2}>
                    <Stack direction='row' alignItems='center' gap={1}>
                      <Icon type='icon-dingding' sx={{ fontSize: 18 }}></Icon>
                      钉钉
                    </Stack>
                  </Stack>
                </Button>
                <Button
                  variant='outlined'
                  color='primary'
                  sx={{ gap: 3 }}
                  disabled
                >
                  <Radio size='small' sx={{ p: 0.5 }} disabled />
                  <Stack direction='row' alignItems='center' gap={2}>
                    <Stack direction='row' alignItems='center' gap={1}>
                      <Icon type='icon-weixin' sx={{ fontSize: 18 }}></Icon>
                      微信
                    </Stack>
                  </Stack>
                </Button>
                <Button
                  variant='outlined'
                  color='primary'
                  sx={{ gap: 3 }}
                  disabled
                >
                  <Radio size='small' sx={{ p: 0.5 }} disabled />
                  <Stack direction='row' alignItems='center' gap={2}>
                    <Stack direction='row' alignItems='center' gap={1}>
                      <Icon type='icon-github' sx={{ fontSize: 18 }}></Icon>
                      GitHub
                    </Stack>
                  </Stack>
                </Button>
              </Stack>
            </Card>
          </Grid>
        </Grid>
      </Grid>

      <MemberManage />
      <LoginHistory />
      <DingingLoginSettingModal
        open={dingdingLoginSettingModalOpen}
        onClose={() => setDingdingLoginSettingModalOpen(false)}
        onOk={() => {
          refresh();
        }}
      />
    </Stack>
  );
};

export default User;
