import { Box, Stack } from '@mui/material';
import { DomainUserCodeRank } from '@/api/types';
import { Modal } from '@c-x/ui';
import { StyledItem, StyledSerialNumber } from './statisticCard';
import User from '@/components/user';

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
      footer={false}
      onOk={onCancel}
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
                <User
                  id={item.user?.id}
                  username={item.user?.username}
                  email={item.user?.email}
                  avatar={item.user?.avatar_url}
                  deleted={item.user?.is_deleted} />
              </Stack>
            </StyledItem>
            <Box sx={{ fontSize: 14 }}>{item.lines} 行</Box>
          </Stack>
        ))}
      </Box>
    </Modal>
  );
};

export default ContributionModal;
