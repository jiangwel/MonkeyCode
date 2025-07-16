import Card from '@/components/card';
import { Stack, Box } from '@mui/material';
import { Table } from '@c-x/ui';
import dayjs from 'dayjs';
import { useRequest } from 'ahooks';
import { getAdminLoginHistory } from '@/api/Admin';
import { ColumnsType } from '@c-x/ui/dist/Table';
import { DomainListAdminLoginHistoryResp } from '@/api/types';
import User from '@/components/user';

type LoginHistory = NonNullable<
  DomainListAdminLoginHistoryResp['login_histories']
>[number];

const LoginHistory = () => {
  const { data, loading } = useRequest(() => getAdminLoginHistory({}));
  const columns: ColumnsType<LoginHistory> = [
    {
      title: '账号',
      dataIndex: 'user',
      render: (user, record) => {
        return <User username={record!.user!.username!} />;
      },
    },
    {
      title: '源 IP 地址',
      dataIndex: 'ip',
      render: (ip, record) => {
        return (
          <Stack direction='column'>
            <Box>{record?.ip_info?.ip}</Box>
            <Box sx={{ color: 'text.secondary' }}>
              {record?.ip_info?.country === '中国' ? ('' + record?.ip_info?.province + '-' + record?.ip_info?.city) : (record?.ip_info?.country || '未知')}
            </Box>
          </Stack>
        );
      },
    },
    {
      title: '登录时间',
      dataIndex: 'created_at',
      render: (text) => {
        return dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss');
      },
    },
  ];
  return (
    <Card sx={{ height: '100%' }}>
      <Stack
        direction='row'
        justifyContent='space-between'
        alignItems='center'
        sx={{ mb: 2 }}
      >
        <Box sx={{ fontWeight: 700, lineHeight: '36px' }}>登录记录</Box>
      </Stack>
      <Table
        columns={columns}
        dataSource={data?.login_histories || []}
        pagination={false}
        rowKey='id'
        loading={loading}
        sx={{ mx: -2 }}
      />
    </Card>
  );
};

export default LoginHistory;
