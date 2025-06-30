import Avatar from '@/components/avatar';
import Card from '@/components/card';
import { getChatInfo } from '@/api/Billing';
import MarkDown from '@/components/markDown';
import { addCommasToNumber, processText } from '@/utils';
import { Ellipsis, Icon, Modal } from '@c-x/ui';
import { Box, Stack, Tooltip, useTheme } from '@mui/material';
import dayjs from 'dayjs';
import { useEffect, useState } from 'react';
import { DomainChatRecord } from '@/api/types';

type ConversationItem = any;
type ToolInfo = any;

const ChatDetailModal = ({
  data,
  open,
  onClose,
}: {
  data?: DomainChatRecord;
  open: boolean;
  onClose: () => void;
}) => {
  const theme = useTheme();
  const [ChatDetailModal, setChatDetailModal] =
    useState<ConversationItem | null>(null);
  const [content, setContent] = useState<string>('');
  const [showToolInfo, setShowToolInfo] = useState<{ [key: string]: ToolInfo }>(
    {}
  );

  const getChatDetailModal = () => {
    if (!data) return;
    getChatInfo({ id: data.id! }).then((res) => {
      setContent(res.content || '');
    });
    // getConversationChatDetailModal({ id }).then((res) => {
    //   const newAnswer = res.answer
    //   const toolWrapsIds = newAnswer.match(/<tools id="([^"]+)">/g)?.map(match => {
    //     const idMatch = match.match(/<tools id="([^"]+)">/);
    //     return idMatch ? idMatch[1] : null;
    //   }).filter(Boolean) || [];
    //   const toolIds = newAnswer.match(/<tool id="([^"]+)">/g)?.map(match => {
    //     const idMatch = match.match(/<tool id="([^"]+)">/);
    //     return idMatch ? idMatch[1] : null;
    //   }).filter(Boolean) || [];
    //   const obj: { [key: string]: ToolInfo } = {}
    //   toolWrapsIds.forEach(id => {
    //     obj[id!] = {
    //       done: true,
    //     }
    //   })
    //   toolIds.forEach(id => {
    //     obj[id!] = {
    //       args: false,
    //       result: false,
    //       done: true,
    //     }
    //   })
    //   setShowToolInfo(obj)
    //   setChatDetailModal({ ...res, answer: processText(res.answer) })
    // })
  };

  useEffect(() => {
    if (open) getChatDetailModal();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, open]);

  return (
    <Modal
      title={
        <Ellipsis
          sx={{
            fontWeight: 'bold',
            fontSize: 20,
            lineHeight: '22px',
            width: 700,
          }}
        >
          {data?.question?.replace(/^<task>|<\/task>$/g, '') || '-'}
        </Ellipsis>
      }
      width={800}
      open={open}
      onCancel={onClose}
      footer={null}
    >
      {ChatDetailModal ? (
        <Box sx={{ fontSize: 14 }}>
          <Stack
            direction={'row'}
            alignItems={'center'}
            gap={3}
            sx={{
              fontSize: 14,
              color: 'text.auxiliary',
            }}
          >
            {ChatDetailModal.created_at && (
              <Stack direction={'row'} alignItems={'center'} gap={1}>
                <Icon type='icon-a-shijian2' />
                {dayjs(ChatDetailModal.created_at).format(
                  'YYYY-MM-DD HH:mm:ss'
                )}
              </Stack>
            )}
            {ChatDetailModal.remote_ip && (
              <Stack direction={'row'} alignItems={'center'} gap={1}>
                <Icon type='icon-IPdizhijiancha' />
                {ChatDetailModal.remote_ip}
              </Stack>
            )}
            {ChatDetailModal.model && (
              <Stack direction={'row'} alignItems={'center'} gap={1}>
                <Icon type='icon-moxing' />
                使用模型
                <Box>{ChatDetailModal.model}</Box>
              </Stack>
            )}
            {data?.input_tokens && data?.output_tokens && (
              <Tooltip
                title={
                  <Stack gap={1} sx={{ minWidth: 100, py: 1 }}>
                    <Box>
                      输入 Token 使用： {addCommasToNumber(data?.input_tokens)}
                    </Box>
                    <Box>
                      输出 Token 使用： {addCommasToNumber(data?.output_tokens)}
                    </Box>
                  </Stack>
                }
              >
                <Stack
                  direction={'row'}
                  alignItems={'center'}
                  gap={1}
                  sx={{ cursor: 'pointer' }}
                >
                  <Icon type='icon-moxing' />
                  Token 统计
                  <Box>
                    {addCommasToNumber(
                      data?.input_tokens + data?.output_tokens
                    )}
                  </Box>
                  <Icon type='icon-a-wenhao8' />
                </Stack>
              </Tooltip>
            )}
          </Stack>
          {ChatDetailModal.references?.length > 0 && (
            <>
              <Stack
                direction={'row'}
                alignItems={'center'}
                gap={1}
                sx={{
                  fontWeight: 'bold',
                  mt: 2,
                  mb: 1,
                  '&::before': {
                    content: '""',
                    display: 'inline-block',
                    width: '4px',
                    height: '12px',
                    borderRadius: '2px',
                    backgroundColor: theme.palette.primary.main,
                  },
                }}
              >
                内容来源
              </Stack>
              <Card sx={{ p: 2, bgcolor: 'background.paper2' }}>
                {ChatDetailModal.references.map((item: any, index: number) => (
                  <Stack
                    direction={'row'}
                    alignItems={'center'}
                    gap={1}
                    key={index}
                  >
                    <Avatar
                      src={item.favicon}
                      sx={{ width: 18, height: 18 }}
                      errorIcon={
                        <Icon
                          type='icon-ditu_diqiu'
                          sx={{ fontSize: 18, color: 'text.auxiliary' }}
                        />
                      }
                    />
                    <Ellipsis>
                      <Box
                        component={'a'}
                        href={item.url}
                        target='_blank'
                        sx={{
                          color: 'text.primary',
                          '&:hover': { color: 'primary.main' },
                        }}
                      >
                        {item.title}
                      </Box>
                    </Ellipsis>
                  </Stack>
                ))}
              </Card>
            </>
          )}
          <Stack
            direction={'row'}
            alignItems={'center'}
            gap={1}
            sx={{
              fontWeight: 'bold',
              mt: 2,
              mb: 1,
              '&::before': {
                content: '""',
                display: 'inline-block',
                width: '4px',
                height: '12px',
                borderRadius: '2px',
                backgroundColor: theme.palette.primary.main,
              },
            }}
          >
            回答
          </Stack>
        </Box>
      ) : (
        <Box></Box>
      )}
      <Card
        sx={{
          '.markdown-body': {
            background: 'transparent',
          },
          p: 0,
        }}
      >
        <MarkDown
          showToolInfo={showToolInfo}
          setShowToolInfo={setShowToolInfo}
          content={content}
        />
      </Card>
    </Modal>
  );
};

export default ChatDetailModal;
