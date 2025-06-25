// import {
//   addModel,
//   getModelByProviderBrand,
//   ModelItem,
//   testModel,
//   updateModel,
// } from '@/api';
import {
  postCreateModel,
  putUpdateModel,
  postCheckModel,
  getListModel,
} from '@/api/Model';
import Card from '@/components/card';
import { ModelProvider } from '../constant';
import { Icon, message, Modal } from '@c-x/ui';
import {
  Box,
  Button,
  MenuItem,
  Stack,
  TextField,
  useTheme,
  alpha,
} from '@mui/material';
import { useEffect, useMemo, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import {
  DomainModelBasic,
  ConstsModelType,
  DomainUpdateModelReq,
  DomainCreateModelReq,
  DomainProviderModel,
} from '@/api/types';

type ModelItem = any;

interface AddModelProps {
  open: boolean;
  data?: ModelItem;
  onClose: () => void;
  refresh: () => void;
  modelType: 'llm' | 'coder';
}

const ModelModal = ({
  open,
  onClose,
  refresh,
  data,
  modelType,
}: AddModelProps) => {
  const theme = useTheme();
  const spaceId = 1;
  const {
    formState: { errors },
    handleSubmit,
    control,
    reset,
    setValue,
    watch,
  } = useForm<Required<DomainCreateModelReq>>({
    defaultValues: {
      provider: data?.provider || 'DeepSeek',
      api_base: data?.api_base || ModelProvider.DeepSeek.defaultBaseUrl,
      model_name: data?.model_name || '',
      // api_version: data?.api_version || '',
      api_key: data?.api_key || '',
      // api_header_key: data?.api_header?.split('=')[0] || '',
      // api_header_value: data?.api_header?.split('=')[1] || '',
    },
  });

  const providerBrand = watch('provider') as keyof typeof ModelProvider;

  const [modelUserList, setModelUserList] = useState<DomainProviderModel[]>([]);

  const [loading, setLoading] = useState(false);
  const [modelLoading, setModelLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  // const getModel = (value: AddModelForm) => {
  //   let header = '';
  //   if (value.api_header_key && value.api_header_value) {
  //     header = value.api_header_key + '=' + value.api_header_value;
  //   }
  //   setModelLoading(true);
  //   getModelByProviderBrand({
  //     space_id: spaceId,
  //     api_key: value.api_key,
  //     base_url: value.base_url,
  //     provider_brand: value.provider_brand,
  //     api_header: header,
  //   })
  //     .then((res) => {
  //       setModelUserList(res.models || []);
  //       if (data && (res.models || []).find((it) => it.model === data.model)) {
  //         setValue('model', data.model);
  //       } else {
  //         setValue('model', res.models?.[0]?.model || '');
  //       }
  //       setSuccess(true);
  //     })
  //     .finally(() => {
  //       setModelLoading(false);
  //     });
  // };

  const onCreateModel = (value: DomainCreateModelReq) => {
    return postCreateModel({
      ...value,
      model_type: modelType as ConstsModelType,
    }).then(() => {
      message.success('添加成功');
      reset();
      onClose();
      refresh();
    });
  };

  const onUpdateModel = (value: DomainUpdateModelReq) => {
    return putUpdateModel({
      ...value,
    }).then(() => {
      message.success('修改成功');
      reset();
      onClose();
      refresh();
    });
  };

  const onSubmit = (value: Required<DomainCreateModelReq>) => {
    const header = '';
    // if (value.api_header_key && value.api_header_value) {
    //   header = value.api_header_key + '=' + value.api_header_value;
    // }
    setError('');
    setLoading(true);
    postCheckModel({
      ...value,
    })
      .then((res) => {
        if (data) {
          onUpdateModel({
            ...value,
            id: data.id,
          }).finally(() => {
            setLoading(false);
          });
        } else {
          onCreateModel(value).finally(() => {
            setLoading(false);
          });
        }
      })
      .catch(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    if (open) {
      if (data) {
        if (data.provider_brand && data.provider_brand !== 'Other') {
          // getModel({
          //   api_key: data.api_key || '',
          //   base_url: data.base_url || '',
          //   model: data.model || '',
          //   provider_brand: data.provider_brand || '',
          //   api_version: data.api_version || '',
          //   api_header_key: data.api_header?.split('=')[0] || '',
          //   api_header_value: data.api_header?.split('=')[1] || '',
          // });
        }
        reset(
          {
            provider: data.provider || 'Other',
            model_name: data.model_name || '',
            api_base: data.api_base || '',
            api_key: data.api_key || '',
            // api_version: data.api_version || '',
            // api_header_key: data.api_header?.split('=')[0] || '',
            // api_header_value: data.api_header?.split('=')[1] || '',
          },
          {
            keepDefaultValues: true,
          }
        );
      } else {
        reset();
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, open]);

  useEffect(() => {
    if (open) {
      getListModel().then((res) => {
        setModelUserList(res.providers || []);
      });
    }
  }, [open]);

  const currentModelList = useMemo(() => {
    return (
      modelUserList.find((it) => it.provider === providerBrand)?.models || []
    );
  }, [modelUserList, providerBrand]);

  useEffect(() => {
    if (currentModelList.length > 0) {
      setValue('api_base', currentModelList[0].api_base || '');
      setValue('model_name', currentModelList[0].name || '');
    }
  }, [currentModelList]);

  return (
    <Modal
      title={data ? '修改第三方模型' : '添加第三方模型'}
      open={open}
      width={800}
      onCancel={() => {
        reset();
        setModelUserList([]);
        setSuccess(false);
        setLoading(false);
        setError('');
        onClose();
      }}
      okText='保存'
      onOk={handleSubmit(onSubmit)}
      okButtonProps={{
        loading,
      }}
    >
      <Stack direction={'row'} alignItems={'stretch'} gap={3}>
        <Stack
          gap={1}
          sx={{
            width: 200,
            flexShrink: 0,
            bgcolor: 'rgb(248, 249, 250)',
            borderRadius: '10px',
            p: 1,
          }}
        >
          <Box
            sx={{ fontSize: 14, lineHeight: '24px', fontWeight: 'bold', p: 1 }}
          >
            模型供应商
          </Box>
          {modelUserList.map((it) => (
            <Stack
              direction={'row'}
              alignItems={'center'}
              gap={1.5}
              key={it.provider}
              sx={{
                cursor: 'pointer',
                fontSize: 14,
                lineHeight: '24px',
                p: 1,
                borderRadius: '10px',
                fontWeight: 'bold',
                fontFamily: 'Gbold',
                ...(providerBrand === it.provider && {
                  bgcolor: alpha(theme.palette.primary.main, 0.1),
                  color: 'primary.main',
                }),
                '&:hover': {
                  color: 'primary.main',
                },
              }}
              onClick={() => {
                if (data) return;
                // setModelUserList([]);
                setError('');
                reset(
                  {
                    provider: it.provider as keyof typeof ModelProvider,
                    api_base: '',
                    model_name: '',
                    // api_version: '',
                    api_key: '',
                    // api_header_key: '',
                    // api_header_value: '',
                  },
                  {
                    keepDefaultValues: true,
                  }
                );
              }}
            >
              <Icon
                type={
                  ModelProvider[it.provider as keyof typeof ModelProvider].icon
                }
                sx={{ fontSize: 18 }}
              />
              {it.provider}
            </Stack>
          ))}
        </Stack>
        <Box sx={{ flex: 1 }}>
          <Box sx={{ fontSize: 14, lineHeight: '32px' }}>
            API 地址{' '}
            <Box component={'span'} sx={{ color: 'red' }}>
              *
            </Box>
          </Box>
          <Controller
            control={control}
            name='api_base'
            rules={{
              required: {
                value: true,
                message: 'URL 不能为空',
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                disabled={!ModelProvider[providerBrand].urlWrite}
                size='small'
                placeholder={ModelProvider[providerBrand].defaultBaseUrl}
                error={!!errors.api_base}
                helperText={errors.api_base?.message}
                onChange={(e) => {
                  field.onChange(e.target.value);
                  setModelUserList([]);
                  setValue('model_name', '');
                }}
              />
            )}
          />
          <Stack
            direction={'row'}
            alignItems={'center'}
            justifyContent={'space-between'}
            sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}
          >
            <Box>
              API Secret
              {ModelProvider[providerBrand].secretRequired && (
                <Box component={'span'} sx={{ color: 'red' }}>
                  {' '}
                  *
                </Box>
              )}
            </Box>
            {ModelProvider[providerBrand].modelDocumentUrl && (
              <Box
                component={'span'}
                sx={{
                  color: 'info.main',
                  cursor: 'pointer',
                  ml: 1,
                  textAlign: 'right',
                }}
                onClick={() =>
                  window.open(
                    ModelProvider[providerBrand].modelDocumentUrl,
                    '_blank'
                  )
                }
              >
                查看文档
              </Box>
            )}
          </Stack>
          <Controller
            control={control}
            name='api_key'
            rules={{
              required: {
                value: ModelProvider[providerBrand].secretRequired,
                message: 'API Secret 不能为空',
              },
            }}
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                size='small'
                placeholder=''
                error={!!errors.api_key}
                helperText={errors.api_key?.message}
                onChange={(e) => {
                  field.onChange(e.target.value);
                }}
              />
            )}
          />
          {/* {providerBrand === 'AzureOpenAI' && (
            <>
              <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                API Version
              </Box>
              <Controller
                control={control}
                name='api_version'
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    size='small'
                    placeholder='2024-10-21'
                    error={!!errors.api_version}
                    helperText={errors.api_version?.message}
                    onChange={(e) => {
                      field.onChange(e.target.value);
                      setModelUserList([]);
                      setValue('model', '');
                      setSuccess(false);
                    }}
                  />
                )}
              />
            </>
          )} */}
          {providerBrand === 'Other' ? (
            <>
              <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                模型名称{' '}
                <Box component={'span'} sx={{ color: 'red' }}>
                  *
                </Box>
              </Box>
              <Controller
                control={control}
                name='model_name'
                rules={{
                  required: '模型名称不能为空',
                }}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    size='small'
                    placeholder=''
                    error={!!errors.model_name}
                    helperText={errors.model_name?.message}
                  />
                )}
              />
              <Box sx={{ fontSize: 12, color: 'error.main', mt: 1 }}>
                需要与模型供应商提供的名称完全一致，不要随便填写
              </Box>
            </>
          ) : (
            <>
              <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                模型名称{' '}
                <Box component={'span'} sx={{ color: 'red' }}>
                  *
                </Box>
              </Box>
              <Controller
                control={control}
                name='model_name'
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    select
                    size='small'
                    placeholder=''
                    error={!!errors.model_name}
                    helperText={errors.model_name?.message}
                  >
                    {currentModelList.map((it) => (
                      <MenuItem key={it.name} value={it.name}>
                        {it.name}
                      </MenuItem>
                    ))}
                  </TextField>
                )}
              />
              {/* {ModelProvider[providerBrand].customHeader && (
                <>
                  <Box sx={{ fontSize: 14, lineHeight: '32px', mt: 2 }}>
                    Header
                  </Box>
                  <Stack direction={'row'} gap={1}>
                    <Controller
                      control={control}
                      name='api_header_key'
                      render={({ field }) => (
                        <TextField
                          {...field}
                          fullWidth
                          size='small'
                          placeholder='key'
                          error={!!errors.api_header_key}
                          helperText={errors.api_header_key?.message}
                        />
                      )}
                    />
                    <Box sx={{ fontSize: 14, lineHeight: '36px' }}>=</Box>
                    <Controller
                      control={control}
                      name='api_header_value'
                      render={({ field }) => (
                        <TextField
                          {...field}
                          fullWidth
                          size='small'
                          placeholder='value'
                          error={!!errors.api_header_value}
                          helperText={errors.api_header_value?.message}
                        />
                      )}
                    />
                  </Stack>
                </>
              )} */}
            </>
          )}
          {error && (
            <Card
              sx={{
                color: 'error.main',
                mt: 2,
                fontSize: 12,
                p: 2,
                bgcolor: 'background.paper2',
                border: '1px solid',
                borderColor: 'error.main',
              }}
            >
              {error}
            </Card>
          )}
        </Box>
      </Stack>
    </Modal>
  );
};

export default ModelModal;
