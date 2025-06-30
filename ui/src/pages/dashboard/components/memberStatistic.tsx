import React, { useMemo } from 'react';
import { Grid2 as Grid } from '@mui/material';
import { useParams } from 'react-router-dom';
import MemberInfo from './memberInfo';
import PieCharts from './pieCharts';
import LineCharts from './lineCharts';
import { RecentActivityCard } from './statisticCard';
import { useRequest } from 'ahooks';
import {
  getUserEventsDashboard,
  getUserStatDashboard,
  getUserHeatmapDashboard,
} from '@/api/Dashboard';
import { StyledHighlight } from './globalStatistic';
import { getRecent90DaysData } from '@/utils';
import { DomainUser } from '@/api/types';

const MemberStatistic = ({
  memberData,
  userList,
  onMemberChange,
}: {
  memberData: DomainUser | null;
  userList: DomainUser[];
  onMemberChange: (data: DomainUser) => void;
}) => {
  const { id } = useParams();
  const { data: userEvents } = useRequest(
    () =>
      getUserEventsDashboard({
        user_id: id || '',
      }),
    {
      refreshDeps: [id],
      manual: false,
      ready: !!id,
    }
  );
  const { data: userStat } = useRequest(
    () =>
      getUserStatDashboard({
        user_id: id || '',
      }),
    {
      refreshDeps: [id],
      manual: false,
      ready: !!id,
    }
  );
  const { data: userHeatmap } = useRequest(
    () =>
      getUserHeatmapDashboard({
        user_id: id || '',
      }),
    {
      refreshDeps: [id],
      manual: false,
      ready: !!id,
    }
  );

  const {
    chatChartData,
    codeCompletionChartData,
    codeLineChartData,
    acceptedPerChartData,
  } = useMemo(() => {
    const {
      accepted_per = [],
      chats = [],
      code_completions = [],
      lines_of_code = [],
    } = userStat || {};
    const chatChartData = getRecent90DaysData(chats, {
      valueLabel: 'value',
    });
    const codeCompletionChartData = getRecent90DaysData(code_completions, {
      valueLabel: 'value',
    });
    const codeLineChartData = getRecent90DaysData(lines_of_code, {
      valueLabel: 'value',
    });
    const acceptedPerChartData = getRecent90DaysData(accepted_per, {
      valueLabel: 'value',
    });
    return {
      chatChartData,
      codeCompletionChartData,
      codeLineChartData,
      acceptedPerChartData,
    };
  }, [userStat]);
  return (
    <Grid
      container
      spacing={2}
      sx={{
        height: '100%',
        overflow: 'auto',
        borderRadius: '10px',
      }}
    >
      <Grid container size={9}>
        <Grid size={12}>
          <MemberInfo
            data={userHeatmap || {}}
            memberData={memberData}
            userList={userList}
            onMemberChange={onMemberChange}
          />
        </Grid>
        <Grid size={6}>
          <PieCharts
            title='工作模式-对话任务'
            extra='最近 90 天'
            data={userStat?.work_mode || []}
          />
        </Grid>
        <Grid size={6}>
          <PieCharts
            title='编程语言'
            extra='最近 90 天'
            data={userStat?.program_language || []}
          />
        </Grid>
      </Grid>
      <Grid size={3}>
        <RecentActivityCard data={userEvents || []} />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='对话任务'
          data={chatChartData}
          extra={
            <>
              最近 90 天共
              <StyledHighlight>{userStat?.total_chats || 0}</StyledHighlight>
              个对话任务
            </>
          }
        />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='补全任务'
          data={codeCompletionChartData}
          extra={
            <>
              最近 90 天共
              <StyledHighlight>
                {userStat?.total_completions || 0}
              </StyledHighlight>
              个补全任务
            </>
          }
        />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='代码量'
          data={codeLineChartData}
          extra={
            <>
              最近 90 天共修改
              <StyledHighlight>
                {userStat?.total_lines_of_code || 0}
              </StyledHighlight>
              行代码
            </>
          }
        />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='补全任务采纳率'
          data={acceptedPerChartData}
          extra={
            <>
              最近 90 天平均采纳率为
              <StyledHighlight>
                {(userStat?.total_accepted_per || 0).toFixed(2)}
              </StyledHighlight>
              %
            </>
          }
        />
      </Grid>
    </Grid>
  );
};

export default MemberStatistic;
