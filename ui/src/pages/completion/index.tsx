import React, { useState, useEffect } from 'react';
import { DomainCompletionRecord } from '@/api/types';
import { getListCompletionRecord } from '@/api/Billing';
import { useRequest } from 'ahooks';
import { Table } from '@c-x/ui';
import Card from '@/components/card';
import { Box, Button, ButtonBase, Chip, Stack, alpha } from '@mui/material';
import dayjs from 'dayjs';
import { ColumnsType } from '@c-x/ui/dist/Table';
import { addCommasToNumber } from '@/utils';
import CompletionDetailModal from './completionDetailModal';
import StyledLabel from '@/components/label';

const Completion = () => {
  const [page, setPage] = useState(1);
  const [size, setSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [dataSource, setDataSource] = useState<DomainCompletionRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [completionDetailModal, setCompletionDetailModal] = useState<
    DomainCompletionRecord | undefined
  >();
  const fetchData = async () => {
    setLoading(true);
    const res = await getListCompletionRecord({
      page: page,
      size: size,
    });
    setLoading(false);
    setTotal(res?.total_count || 0);
    setDataSource(res.records || []);
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, size]);
  const columns: ColumnsType<DomainCompletionRecord> = [
    {
      dataIndex: 'user',
      title: '成员',

      render(value: DomainCompletionRecord['user']) {
        return value?.username;
      },
    },
    {
      dataIndex: 'task',
      title: '补全内容',
      width: 150,
      render(_, record) {
        return (
          <ButtonBase
            disableRipple
            onClick={() => setCompletionDetailModal(record)}
            sx={{ color: 'info.main' }}
          >
            点击查看
          </ButtonBase>
        );
      },
    },

    {
      dataIndex: 'is_accept',
      title: '是否采纳',
      width: 130,
      render(value: boolean) {
        const color = value ? 'success' : 'default';
        return (
          <StyledLabel color={color}>{value ? '已采纳' : '未采纳'}</StyledLabel>
        );
      },
    },

    {
      dataIndex: 'program_language',
      title: '编程语言',
      width: 140,
    },
    {
      dataIndex: 'input_tokens',
      title: '输入 Token',
      render(value: number) {
        return addCommasToNumber(value);
      },
    },
    {
      dataIndex: 'output_tokens',
      title: '输出 Token',

      render(value: number) {
        return addCommasToNumber(value);
      },
    },
    {
      dataIndex: 'created_at',
      title: '时间',
      width: 200,
      render(value: number) {
        return dayjs.unix(value).format('YYYY-MM-DD HH:mm:ss');
      },
    },
  ];
  return (
    <Card sx={{ flex: 1, height: '100%' }}>
      <Table
        size='large'
        height='100%'
        loading={loading}
        columns={columns}
        dataSource={dataSource || []}
        sx={{ mx: -2 }}
        PaginationProps={{
          sx: {
            pt: 2,
            mx: 2,
          },
        }}
        rowKey='id'
        pagination={{
          pageSize: size,
          total: total,
          page: page,
          onChange: (page, size) => {
            setPage(page);
            setSize(size);
          },
        }}
      />
      <CompletionDetailModal
        open={!!completionDetailModal}
        onClose={() => setCompletionDetailModal(undefined)}
        data={completionDetailModal}
      />
    </Card>
  );
};

export default Completion;
