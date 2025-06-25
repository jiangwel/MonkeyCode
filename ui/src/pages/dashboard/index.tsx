import React, { useEffect, useMemo, useState } from 'react';
import { getListUser } from '@/api/User';
import {
  Stack,
  TextField,
  MenuItem,
  InputAdornment,
  IconButton,
} from '@mui/material';
import dayjs from 'dayjs';
import { CusTabs } from '@c-x/ui';
import GlobalStatistic from './components/globalStatistic';
import CancelRoundedIcon from '@mui/icons-material/CancelRounded';
import { useRequest } from 'ahooks';
import MemberStatistic from './components/memberStatistic';
import { useParams } from 'react-router-dom';
import { useNavigate } from 'react-router-dom';
import { DomainUser } from '@/api/types';

const Dashboard = () => {
  const navigate = useNavigate();
  const { tab, id } = useParams();
  const [tabValue, setTabValue] = useState(tab || 'global');
  const [memberId, setMemberId] = useState(id || '');
  const [memberData, setMemberData] = useState<DomainUser | null>(null);

  const { data: userData, refresh } = useRequest(
    () =>
      getListUser({
        page: 1,
        size: 99999,
      }),
    {
      manual: true,
      onSuccess: (res) => {
        if (id) {
          setMemberId(id);
          setMemberData(res.users?.find((item) => item.id === id) || null);
        } else {
          setMemberData(res.users?.[0] || null);
          setMemberId(res.users?.[0]?.id || '');
          navigate(`/dashboard/member/${res.users?.[0]?.id}`);
        }
      },
    }
  );
  const userList = useMemo(() => {
    return userData?.users || [];
  }, [userData]);
  useEffect(() => {
    if (tabValue === 'member') {
      refresh();
    }
  }, [tabValue]);

  const onMemberChange = (data: DomainUser) => {
    setMemberId(data.id!);
    setMemberData(data);
    navigate(`/dashboard/member/${data.id}`);
  };

  const onTabChange = (value: string) => {
    setTabValue(value);
    navigate(`/dashboard/${value}`);
  };

  return (
    <Stack gap={2} sx={{ height: '100%' }}>
      <Stack direction='row' justifyContent='space-between' alignItems='center'>
        <CusTabs
          value={tabValue}
          onChange={onTabChange}
          list={[
            { label: '全局统计', value: 'global' },
            { label: '成员统计', value: 'member' },
          ]}
          sx={{
            borderColor: 'transparent',
            backgroundColor: '#fff',
            height: 'auto',
            borderRadius: 2.5,
            '.MuiTabs-indicator': {
              top: '0px',
              bottom: '0px',
              height: 'auto',
              borderRadius: '5px',
            },
            '.MuiTab-root': {
              lineHeight: '1.5 !important',
            },
          }}
        />
      </Stack>

      {tabValue === 'global' && <GlobalStatistic />}
      {tabValue === 'member' && (
        <MemberStatistic
          memberData={memberData}
          userList={userList}
          onMemberChange={onMemberChange}
        />
      )}
    </Stack>
  );
};

export default Dashboard;
