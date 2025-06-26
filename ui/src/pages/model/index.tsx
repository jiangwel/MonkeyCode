import React from 'react';
import TokenUsage from './components/tokenUsage';
import ModelCard from './components/modelCard';
import { Stack } from '@mui/material';

const Model = () => {
  return (
    <Stack gap={2}>
      <TokenUsage />
      <ModelCard title='对话模型' modelType='llm' />
      <ModelCard title='代码补全模型' modelType='coder' />
    </Stack>
  );
};

export default Model;
