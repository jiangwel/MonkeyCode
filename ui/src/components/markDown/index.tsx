// import { ToolInfo } from '@/api';
import { Icon, message } from '@c-x/ui';
import { Box, Button, IconButton, Stack, useTheme, alpha } from '@mui/material';
import React, { useState } from 'react';
import ReactMarkdown, { Components } from 'react-markdown';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import rehypeRaw from 'rehype-raw';
import rehypeSanitize, { defaultSchema } from 'rehype-sanitize';
import remarkBreaks from 'remark-breaks';
import remarkGfm from 'remark-gfm';
import ExpandMoreRoundedIcon from '@mui/icons-material/ExpandMoreRounded';

interface ExtendedComponents extends Components {
  tools?: React.ComponentType<any>;
}

export const toolNames = [
  'execute_command',
  'read_file',
  'write_to_file',
  'apply_diff',
  'insert_content',
  'search_and_replace',
  'search_files',
  'list_files',
  'list_code_definition_names',
  'browser_action',
  'use_mcp_tool',
  'access_mcp_resource',
  'ask_followup_question',
  'attempt_completion',
  'switch_mode',
  'new_task',
  'fetch_instructions',
];

// 去掉下划线的标签名，用于Markdown渲染
export const toolTagNames = toolNames.map((name) => name.replace(/_/g, ''));

type ToolInfo = any;

const MarkDown = ({
  loading = false,
  content,
  showToolInfo = {},
  setShowToolInfo,
  setCurrentToolId,
  handleSearchAbort,
}: {
  loading?: boolean;
  content: string;
  showToolInfo: Record<string, ToolInfo>;
  setShowToolInfo: (value: Record<string, ToolInfo>) => void;
  setCurrentToolId?: (value: string) => void;
  handleSearchAbort?: () => void;
}) => {
  const theme = useTheme();
  const [showThink, setShowThink] = useState(false);

  // 将content中的下划线标签替换为无下划线版本
  const processContent = (content: string) => {
    let processedContent = content;

    toolNames.forEach((toolName) => {
      const withUnderscore = toolName;
      const withoutUnderscore = toolName.replace(/_/g, '');

      // 替换开始标签
      const openTagPattern = new RegExp(`<${withUnderscore}>`, 'g');
      processedContent = processedContent.replace(
        openTagPattern,
        `<${withoutUnderscore}>`
      );

      // 替换结束标签
      const closeTagPattern = new RegExp(`</${withUnderscore}>`, 'g');
      processedContent = processedContent.replace(
        closeTagPattern,
        `</${withoutUnderscore}>`
      );
    });

    return processedContent;
  };

  // const answer = processContent(content);

  const answer = content;

  if (content.length === 0) return null;

  return (
    <Box
      className='markdown-body'
      id='markdown-body'
      sx={{
        fontSize: '14px',
        background: 'transparent',
        '#chat-thinking': {
          display: 'flex',
          alignItems: 'flex-end',
          gap: '16px',
          fontSize: '12px',
          color: 'text.secondary',
          marginBottom: '40px',
          lineHeight: '20px',
          backgroundColor: 'background.paper',
          padding: '16px',
          cursor: 'pointer',
          borderRadius: '10px',
          div: {
            transition: 'height 0.3s',
            overflow: 'hidden',
            height: showThink ? 'auto' : '60px',
          },
        },
      }}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm, remarkBreaks]}
        rehypePlugins={[
          rehypeRaw,
          [
            rehypeSanitize,
            {
              tagNames: [
                ...(defaultSchema.tagNames! as string[]),
                // 'thinking',
                // 'tools',
                // 'tool',
                // 'toolname',
                // 'toolargs',
                // 'toolresult',
                // 'error',
                // 'attemptcompletion',
                ...toolTagNames,
              ],
            },
          ],
        ]}
        components={
          {
            error: ({
              children,
              ...rest
            }: React.HTMLAttributes<HTMLElement>) => {
              return (
                <div className='chat-error' {...rest}>
                  {children}
                </div>
              );
            },
            tools: ({
              id = '',
              ...rest
            }: React.HTMLAttributes<HTMLElement>) => {
              const _id = id.replace('user-content-', '');
              return (
                <div className='chat-tools' id={_id} {...rest}>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      marginBottom: 8,
                    }}
                  >
                    <div className='chat-tool-name-text'>
                      <Icon type='icon-gongju-tool' />
                      工具调用
                    </div>
                  </div>
                  <div {...rest} style={{ zIndex: 1 }}></div>
                  {!showToolInfo[_id].done && (
                    <Stack direction='row' alignItems='center' gap={2}>
                      <Button
                        variant='contained'
                        onClick={() => setCurrentToolId?.(_id)}
                        className='chat-tool-run'
                        sx={{
                          width: 80,
                          bgcolor: 'primary.main',
                          '.MuiButton-startIcon': {
                            mr: 0,
                          },
                        }}
                        startIcon={
                          <Icon
                            type='icon-yunhang'
                            sx={{ fontSize: 12, mr: 0 }}
                          />
                        }
                      >
                        运行
                      </Button>
                      <Button
                        variant='outlined'
                        sx={{ width: 80 }}
                        onClick={handleSearchAbort}
                      >
                        结束
                      </Button>
                    </Stack>
                  )}
                </div>
              );
            },
            tool: ({
              id = '',
              ...rest
            }: React.HTMLAttributes<HTMLElement> & { id?: string }) => {
              const _id = id.replace('user-content-', '');
              const className = showToolInfo[_id]
                ? showToolInfo[_id].args
                  ? 'chat-tool chat-tool-expend-args'
                  : showToolInfo[_id].result
                  ? 'chat-tool chat-tool-expend-result'
                  : 'chat-tool'
                : 'chat-tool';
              return (
                <div className={className} id={_id}>
                  <div {...rest} style={{ zIndex: 1 }}></div>
                  {!!showToolInfo[_id] && (
                    <div className='chat-tool-expend-btn'>
                      <div
                        className={
                          showToolInfo[_id].args
                            ? 'chat-tool-expend-text chat-tool-expend-text-active'
                            : 'chat-tool-expend-text'
                        }
                        onClick={() => {
                          setShowToolInfo({
                            ...showToolInfo,
                            [_id]: {
                              args: !showToolInfo[_id].args,
                              result: false,
                              done: showToolInfo[_id].done,
                            },
                          });
                        }}
                      >
                        参数
                        <ExpandMoreRoundedIcon
                          sx={{
                            fontSize: 16,
                            ml: 0,
                            transform: showToolInfo[_id].args
                              ? 'rotate(-180deg)'
                              : 'rotate(0deg)',
                          }}
                        />
                      </div>
                      {showToolInfo[_id].done && (
                        <div
                          className={
                            showToolInfo[_id].result
                              ? 'chat-tool-expend-text chat-tool-expend-text-active'
                              : 'chat-tool-expend-text'
                          }
                          onClick={() => {
                            setShowToolInfo({
                              ...showToolInfo,
                              [_id]: {
                                args: false,
                                result: !showToolInfo[_id].result,
                                done: showToolInfo[_id].done,
                              },
                            });
                          }}
                        >
                          结果
                          <ExpandMoreRoundedIcon
                            sx={{
                              fontSize: 16,
                              ml: 0,
                              transform: showToolInfo[_id].result
                                ? 'rotate(-180deg)'
                                : 'rotate(0deg)',
                            }}
                          />
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            },
            toolname: (props: React.HTMLAttributes<HTMLElement>) => {
              return <div className='chat-tool-name'>{props.children}</div>;
            },
            toolargs: ({
              children,
              ...rest
            }: React.HTMLAttributes<HTMLElement>) => {
              const safeChildren = React.Children.toArray(children).filter(
                (child) => child !== '\n'
              );
              let innerText: React.ReactNode = '';
              try {
                if (
                  safeChildren.length > 1 &&
                  React.isValidElement(safeChildren[1])
                ) {
                  const secondChild = safeChildren[1] as React.ReactElement<{
                    children?: React.ReactNode;
                  }>;
                  if (secondChild.props && secondChild.props.children) {
                    const jsonString = String(secondChild.props.children);
                    innerText = JSON.stringify(JSON.parse(jsonString), null, 2);
                  }
                } else {
                  innerText = safeChildren;
                }
              } catch (err) {
                console.error(err);
                innerText = safeChildren;
              }
              return (
                <div className='chat-tool-args'>
                  <pre {...rest}>{innerText}</pre>
                </div>
              );
            },
            toolresult: ({
              children,
              ...rest
            }: React.HTMLAttributes<HTMLElement>) => {
              const safeChildren = React.Children.toArray(
                children || []
              ).filter((child) => child !== '\n');
              const hasPreTag = safeChildren.some(
                (child) => React.isValidElement(child) && child.type === 'pre'
              );
              return hasPreTag ? (
                <div
                  className='chat-tool-result'
                  {...rest}
                  children={safeChildren}
                />
              ) : (
                <div className='chat-tool-result'>
                  <pre {...rest} children={safeChildren} />
                </div>
              );
            },
            // thinking: (props: React.HTMLAttributes<HTMLElement>) => {
            //   return (
            //     <div id='chat-thinking'>
            //       <div
            //         className={!showThink ? 'three-ellipsis' : ''}
            //         {...props}
            //       ></div>
            //       {!loading && (
            //         <IconButton
            //           size='small'
            //           onClick={() => setShowThink(!showThink)}
            //           sx={{
            //             bgcolor: 'background.paper',
            //             ':hover': {
            //               bgcolor: alpha(theme.palette.primary.main, 0.1),
            //               color: theme.palette.primary.main,
            //             },
            //           }}
            //         >
            //           <ExpandMoreRoundedIcon
            //             sx={{
            //               fontSize: 18,
            //               flexShrink: 0,
            //               transform: showThink
            //                 ? 'rotate(-180deg)'
            //                 : 'rotate(0deg)',
            //             }}
            //           />
            //         </IconButton>
            //       )}
            //     </div>
            //   );
            // },
            h1: (props: React.HTMLAttributes<HTMLHeadingElement>) => (
              <h2 {...props} />
            ),
            a: ({
              children,
              style,
              ...rest
            }: React.HTMLAttributes<HTMLAnchorElement>) => (
              <a
                {...rest}
                target='_blank'
                rel='noopener noreferrer'
                style={{
                  color: theme.palette.primary.main,
                  textDecoration: 'underline',
                  ...style,
                }}
              >
                {children}
              </a>
            ),
            img: (props: React.ImgHTMLAttributes<HTMLImageElement>) => {
              const { style, alt, ...rest } = props;
              return (
                <img
                  alt={alt || 'markdown-img'}
                  {...rest}
                  style={{
                    ...style,
                    borderRadius: '10px',
                    marginLeft: '5px',
                    boxShadow: '0px 0px 3px 1px rgba(0,0,5,0.15)',
                    cursor: 'pointer',
                  }}
                  referrerPolicy='no-referrer'
                />
              );
            },
            attemptcompletion: (props: React.HTMLAttributes<HTMLElement>) => {
              return (
                <div
                  style={{
                    backgroundColor: '#f0f9ff',
                    border: '1px solid #0ea5e9',
                    borderRadius: '8px',
                    padding: '12px',
                    margin: '8px 0',
                    color: '#0c4a6e',
                    whiteSpace: 'pre-line',
                  }}
                >
                  {props.children}
                </div>
              );
            },
            code({
              children,
              className,
              ...rest
            }: React.HTMLAttributes<HTMLElement>) {
              const match = /language-(\w+)/.exec(className || '');
              return match ? (
                <SyntaxHighlighter
                  showLineNumbers
                  {...rest}
                  language={match[1] || 'bash'}
                  style={github}
                  onClick={() => {
                    if (navigator.clipboard) {
                      navigator.clipboard.writeText(
                        String(children).replace(/\n$/, '')
                      );
                      message.success('复制成功');
                    }
                  }}
                >
                  {String(children).replace(/\n$/, '')}
                </SyntaxHighlighter>
              ) : (
                <code
                  {...rest}
                  className={className}
                  onClick={() => {
                    if (navigator.clipboard) {
                      navigator.clipboard.writeText(String(children));
                      message.success('复制成功');
                    }
                  }}
                >
                  {children}
                </code>
              );
            },
          } as ExtendedComponents
        }
      >
        {answer}
      </ReactMarkdown>
    </Box>
  );
};

export default MarkDown;
