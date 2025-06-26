import React from 'react';
import Card from '@/components/card';
import { Grid2 as Grid, Stack, styled, Switch } from '@mui/material';
import { useRequest } from 'ahooks';
import { getGetSetting, putUpdateSetting } from '@/api/User';
import MemberManage from './memberManage';
import LoginHistory from './loginHistory';
import { message } from '@c-x/ui';

const StyledCard = styled(Card)({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
});

const StyledLabel = styled('div')(({ theme }) => ({
  fontWeight: 700,
  fontSize: 14,
  color: theme.vars.palette.text.primary,
}));

const User = () => {
  const {
    data = {
      enable_sso: false,
      force_two_factor_auth: false,
      disable_password_login: false,
    },
    refresh,
  } = useRequest(() => getGetSetting());
  const { run: updateSetting } = useRequest(putUpdateSetting, {
    manual: true,
    onSuccess: () => {
      refresh();
      message.success('设置更新成功');
    },
  });
  return (
    <Stack gap={2}>
      <Grid container spacing={2}>
        <Grid size={4}>
          <StyledCard>
            <StyledLabel>单点登录</StyledLabel>
            <Switch
              checked={data?.enable_sso}
              onChange={(e) => updateSetting({ enable_sso: e.target.checked })}
            />
          </StyledCard>
        </Grid>
        <Grid size={4}>
          <StyledCard>
            <StyledLabel>强制成员启用两步认证</StyledLabel>
            <Switch
              checked={data?.force_two_factor_auth}
              onChange={(e) => {
                console.log(e.target.checked);
                updateSetting({ force_two_factor_auth: e.target.checked });
              }}
            />
          </StyledCard>
        </Grid>
        <Grid size={4}>
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

      <MemberManage />
      <LoginHistory />
    </Stack>
  );
};

export default User;
