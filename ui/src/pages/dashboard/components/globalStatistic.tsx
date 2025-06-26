import React, { useMemo } from 'react';
import { Grid2 as Grid, styled } from '@mui/material';
import dayjs from 'dayjs';
import {
  getStatisticsDashboard,
  getUserCodeRankDashboard,
  getTimeStatDashboard,
  getCategoryStatDashboard,
} from '@/api/Dashboard';
import { useRequest } from 'ahooks';
import { UserCard } from './statisticCard';
import LineCharts from './lineCharts';
import PieCharts from './pieCharts';
import BarCharts from './barCharts';
import { ContributionCard } from './statisticCard';
import { DomainTimeStat } from '@/api/types';
import { getRecent90DaysData, getRecent60MinutesData } from '@/utils';

export const StyledHighlight = styled('span')(({ theme }) => ({
  fontSize: 12,
  color: theme.palette.text.primary,
  padding: '0 4px',
}));

const GlobalStatistic = () => {
  const { data: statisticsData } = useRequest(getStatisticsDashboard);
  const { data: userCodeRankData } = useRequest(getUserCodeRankDashboard);
  const { data: timeStatData } = useRequest(getTimeStatDashboard);
  const {
    data: categoryStatData = {
      program_language: [],
      work_mode: [],
    },
  } = useRequest(getCategoryStatDashboard);

  const {
    userActiveChartData,
    chatChartData,
    codeCompletionChartData,
    codeLineChartData,
    realTimeTokenChartData,
    acceptedPerChartData,
  } = useMemo(() => {
    const {
      active_users = [],
      chats = [],
      code_completions = [],
      lines_of_code = [],
      real_time_tokens = [],
      accepted_per = [],
    } = timeStatData || {};
    const userActiveChartData = getRecent90DaysData(active_users, {
      valueLabel: 'value',
    });
    const chatChartData = getRecent90DaysData(chats, {
      valueLabel: 'value',
    });
    const codeCompletionChartData = getRecent90DaysData(code_completions, {
      valueLabel: 'value',
    });
    const codeLineChartData = getRecent90DaysData(lines_of_code, {
      valueLabel: 'value',
    });
    const realTimeTokenChartData = getRecent60MinutesData(real_time_tokens, {
      valueLabel: 'value',
    });
    const acceptedPerChartData = getRecent90DaysData(accepted_per, {
      valueLabel: 'value',
    });
    return {
      userActiveChartData,
      chatChartData,
      codeCompletionChartData,
      codeLineChartData,
      realTimeTokenChartData,
      acceptedPerChartData,
    };
  }, [timeStatData]);

  return (
    <Grid
      container
      spacing={2}
      sx={{
        height: '100%',
        overflow: 'auto',
        borderRadius: 2.5,
      }}
    >
      {/* <Grid container size={9} spacing={2}> */}
      <Grid size={3}>
        <UserCard data={statisticsData} />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='活跃用户'
          data={userActiveChartData}
          extra={
            <>
              最近 90 天共
              <StyledHighlight>
                {timeStatData?.total_users || 0}
              </StyledHighlight>
              个活跃用户数
            </>
          }
        />
      </Grid>
      <Grid size={3}>
        <ContributionCard data={userCodeRankData} />
      </Grid>
      <Grid size={4}>
        <PieCharts
          title='工作模式-对话任务'
          data={categoryStatData.work_mode || []}
          extra='最近 90 天'
        />
      </Grid>
      <Grid size={4}>
        <PieCharts
          title='编程语言'
          data={categoryStatData.program_language || []}
          extra='最近 90 天'
        />
      </Grid>
      <Grid size={4}>
        <BarCharts
          title='实时使用情况'
          data={realTimeTokenChartData}
          extra='近 60 分钟'
        />
      </Grid>
      <Grid size={6}>
        <LineCharts
          title='对话任务'
          data={chatChartData}
          extra={
            <>
              最近 90 天共
              <StyledHighlight>
                {timeStatData?.total_chats || 0}
              </StyledHighlight>
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
                {timeStatData?.total_completions || 0}
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
                {timeStatData?.total_lines_of_code || 0}
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
                {(timeStatData?.total_accepted_per || 0).toFixed(2)}
              </StyledHighlight>
              %
            </>
          }
        />
      </Grid>
      {/* </Grid> */}
    </Grid>
  );
};

export default GlobalStatistic;
