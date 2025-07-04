import { Box, Stack } from '@mui/material';
import { DomainUserCodeRank } from '@/api/types';
import { Modal } from '@c-x/ui';
import { StyledItem, StyledSerialNumber, StyledText } from './statisticCard';

const ContributionModal = ({
  open,
  onCancel,
  data,
}: {
  open: boolean;
  onCancel: () => void;
  data: DomainUserCodeRank[];
}) => {
  return (
    <Modal
      open={open}
      onCancel={onCancel}
      title='用户贡献榜'
      showCancel={false}
      okText='关闭'
    >
      <Box
        gap={0.5}
        sx={{
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        {data.map((item, index) => (
          <Stack
            direction='row'
            alignItems='center'
            justifyContent='space-between'
            key={item.username}
          >
            <StyledItem>
              <StyledSerialNumber num={index + 1}>
                {index + 1}
              </StyledSerialNumber>
              <Stack
                direction='row'
                alignItems='center'
                gap={1.5}
                sx={{
                  flex: 1,
                  minWidth: 0,
                }}
              >
                {/* <Avatar size={24} src={item.avatar} /> */}
                <StyledText className='active-user-name'>
                  {item.username}
                </StyledText>
              </Stack>
            </StyledItem>
            <Box sx={{ fontSize: 14 }}>{item.lines}</Box>
          </Stack>
        ))}
      </Box>
    </Modal>
  );
};

export default ContributionModal;
