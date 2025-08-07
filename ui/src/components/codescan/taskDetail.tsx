import Card from '@/components/card';
import { Ellipsis, Modal } from '@c-x/ui';

import { useEffect, useState } from 'react';
import { DomainSecurityScanningResult, DomainSecurityScanningRiskDetail } from '@/api/types';
import { getSecurityScanningDetail } from '@/api';
import { Box, CircularProgress, List, ListItem, ListItemButton, Stack } from '@mui/material';

interface RiskLevelBoxProps {
  level: 'severe' | 'critical' | 'suggest';
}

const RiskLevelBox = ({ level }: RiskLevelBoxProps) => {
  const riskConfig = {
    severe: {
      text: '严重',
      color: 'risk.severe',
    },
    critical: {
      text: '高风险',
      color: 'risk.critical',
    },
    suggest: {
      text: '低风险',
      color: 'risk.suggest',
    },
  };

  const config = riskConfig[level];

  if (!config) return null;

  return (
    <Box sx={{
      backgroundColor: config.color,
      color: '#fff',
      borderRadius: '4px',
      textAlign: 'center',
      width: '80px',
      minWidth: '80px',
      fontSize: '12px',
      lineHeight: '20px'
    }}>
      {config.text}
    </Box>
  );
};

const TaskDetail = ({
  task,
  open,
  onClose,
}: {
  task?: DomainSecurityScanningResult;
  open: boolean;
  onClose: () => void;
}) => {
  const [loading, setLoading] = useState(true);
  const [vulns, setVulns] = useState<DomainSecurityScanningRiskDetail[]>([]);

  const fetchData = async () => {
    setLoading(true);
    const resp = await getSecurityScanningDetail({
      id: task?.id as string
    });
    setVulns(resp);
    setLoading(false);
  };

  useEffect(() => {
    console.log(task)
    if (open) {
      fetchData();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [task, open]);

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
          {task?.name} / {task?.project_name}
        </Ellipsis>
      }
      sx={{
        '.MuiDialog-paper': {
          maxWidth: 1300,
        },
      }}
      width={1200}
      open={open}
      onCancel={onClose}
      footer={null}
    >
      <Card sx={{ p: 0, background: 'transparent', boxShadow: 'none' }}>
        <List>
          {loading ? (
            <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
              <CircularProgress />
            </div>
          ) : (
            vulns.map((vuln) => (
              <ListItem key={vuln.id} sx={{
                padding: 0,
                width: '100%'
              }}>
                <ListItemButton sx={{
                  borderBottomWidth: '1px',
                  borderBottomStyle: 'solid',
                  borderBottomColor: 'background.paper',
                  fontSize: '14px',
                  width: '100%'
                }}>
                  <Stack direction={"column"} sx={{
                    width: '100%'
                  }}>
                    <Stack direction={"row"}>
                      <RiskLevelBox level={vuln.level as 'severe' | 'critical' | 'suggest'} />
                      <Box sx={{
                        fontSize: '14px',
                        ml: '20px',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        lineHeight: '20px'
                      }}>{vuln.desc}</Box>
                    </Stack>
                    <Box sx={{
                      color: 'text.tertiary',
                      fontSize: '14px',
                      mt: '6px',
                      width: '100%',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}>
                      {vuln.filename}:{vuln?.start?.line}
                    </Box>
                  </Stack>
                </ListItemButton>
              </ListItem>
            ))
          )}
        </List>
      </Card>
    </Modal>
  );
};

export default TaskDetail;
