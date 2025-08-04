import React, { useEffect, useRef, useState } from 'react';
import Card from '@/components/card';
import dayjs from 'dayjs';
import {
  Box,
  Stack,
  Typography,
  Tooltip,
  useTheme,
  IconButton,
  Menu,
  MenuItem,
} from '@mui/material';
import { ActivityCalendar } from 'react-activity-calendar';
import { DomainUserHeatmap, DomainUserHeatmapResp } from '@/api/types';
import { DomainUser } from '@/api/types';
import Avatar from '@/components/avatar';
import { Icon } from '@c-x/ui';

const getRecent1YearData = (
  data: DomainUserHeatmap[] = [],
  max_count: number
) => {
  const today = dayjs();
  const lastYearToday = today.subtract(1, 'year');
  const diffInDays = today.diff(lastYearToday, 'day');
  const result: { count: number; date: string; level: number }[] = [];
  const dateMap: Record<string, number> = {};
  data.forEach((item) => {
    dateMap[dayjs.unix(item.date!).format('YYYY-MM-DD')] = item.count!;
  });

  const getLevel = (count: number) => {
    if (count === 0) return 0;
    if (count === 1) return 1;
    if (count <= Math.max(2, max_count * 0.25)) return 2;
    if (count <= Math.max(3, max_count * 0.6)) return 3;
    return 4;
  };

  for (let i = 0; i < diffInDays; i++) {
    const time = today.subtract(i, 'day').format('YYYY-MM-DD');
    const count = dateMap[time] || 0;
    result.unshift({
      count,
      date: time,
      level: getLevel(count),
    });
  }
  return result;
};

const MemberInfo = ({
  data,
  memberData,
  userList,
  onMemberChange,
}: {
  data: DomainUserHeatmapResp;
  memberData: DomainUser | null;
  userList?: DomainUser[];
  onMemberChange?: (data: DomainUser) => void;
}) => {
  const theme = useTheme();
  const [blockSize, setBlockSize] = useState(8);
  const ref = useRef<HTMLDivElement>(null);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };
  useEffect(() => {
    const innerWidth = ref.current?.offsetWidth;
    const dis = Math.max(0, Math.ceil((innerWidth! - 980) / 54));
    console.log(dis);
    setBlockSize(dis + 8);
  }, []);

  return (
    <Card
      ref={ref}
      sx={{
        display: 'flex',
        justifyContent: 'space-between',
        gap: 2,
        '.react-activity-calendar__scroll-container': {
          p: '1px',
        },
      }}
    >
      <Stack
        gap={1}
        justifyContent='space-between'
        sx={{
          bgcolor: '#f8f9fa',
          p: 2,
          borderRadius: 2.5,
          width: 240,
        }}
      >
        <Stack gap={2}>
          <Menu
            id='basic-menu'
            anchorEl={anchorEl}
            open={open}
            onClose={handleClose}
          >
            {userList?.map((item) => (
              <MenuItem
                key={item.id}
                selected={memberData?.id === item.id}
                onClick={() => {
                  onMemberChange?.(item);
                  handleClose();
                }}
                sx={{
                  fontSize: 14,
                }}
              >
                {item.username}
              </MenuItem>
            ))}
          </Menu>
          <Stack
            direction='row'
            alignItems='center'
            justifyContent='space-between'
            sx={{ mb: 1 }}
          >
            <Avatar name={memberData?.username || ''} />
            {userList && (
              <IconButton onClick={handleClick} size='small'>
                <Icon
                  type='icon-qiehuan'
                  sx={{ fontSize: 16, color: 'text.primary' }}
                />
              </IconButton>
            )}
          </Stack>
          <Box sx={{ fontSize: 16, fontWeight: 700 }}>
            {memberData?.username}
          </Box>
        </Stack>

        <Box sx={{ fontSize: 14, color: 'text.tertiary' }}>
          {dayjs.unix(memberData?.created_at || 0).fromNow()}加入
        </Box>
      </Stack>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          position: 'relative',
          minWidth: 0,
        }}
      >
        <ActivityCalendar
          data={getRecent1YearData(data.points || [], data?.max_count || 1)}
          colorScheme='light'
          blockSize={blockSize}
          theme={{
            light: [
              theme.palette.grey[200], // level 0: 无使用 - 浅灰色背景
              theme.palette.grey[400], // level 1: 轻度使用 - 浅灰色
              theme.palette.grey[600], // level 2: 中度使用 - 中灰色
              theme.palette.grey[800], // level 3: 高度使用 - 深灰色
              theme.palette.grey[900], // level 4: 重度使用 - 最深灰色
            ],
          }}
          labels={{
            totalCount: '{{count}} 次使用',
            months: [
              '一月',
              '二月',
              '三月',
              '四月',
              '五月',
              '六月',
              '七月',
              '八月',
              '九月',
              '十月',
              '十一月',
              '十二月',
            ],
            weekdays: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
            legend: {
              less: '较少',
              more: '较多',
            },
          }}
          // hideTotalCount
          // hideMonthLabels
          // hideColorLegend
          renderBlock={(block, activity) => (
            <Tooltip
              slotProps={{
                tooltip: {
                  sx: {
                    mb: '0 !important',
                  },
                },
              }}
              title={
                <Box
                  sx={{ height: 14 }}
                >{`${activity.date} ${activity.count} 次使用`}</Box>
              }
              placement='top'
            >
              {block}
            </Tooltip>
          )}
        />
      </Box>
    </Card>
  );
};

export default MemberInfo;
