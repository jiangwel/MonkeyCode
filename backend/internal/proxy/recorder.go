package proxy

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"

	"github.com/rokku-c/go-openai"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/promptparser"
	"github.com/chaitin/MonkeyCode/backend/pkg/tee"
)

type ChatRecorder struct {
	*tee.Tee
	ctx        context.Context
	cx         *Ctx
	usecase    domain.ProxyUsecase
	req        *openai.ChatCompletionRequest
	logger     *slog.Logger
	model      *domain.Model
	completion strings.Builder // 累积完整的响应内容
	usage      *openai.Usage   // 最终的使用统计
	buffer     strings.Builder // 缓存不完整的行
	recorded   bool            // 标记是否已记录完整对话
}

func NewChatRecorder(
	ctx context.Context,
	cx *Ctx,
	usecase domain.ProxyUsecase,
	model *domain.Model,
	req *openai.ChatCompletionRequest,
	r io.Reader,
	w io.Writer,
	logger *slog.Logger,
) *ChatRecorder {
	c := &ChatRecorder{
		ctx:     ctx,
		cx:      cx,
		usecase: usecase,
		model:   model,
		req:     req,
		logger:  logger,
	}
	c.Tee = tee.NewTee(ctx, logger, r, w, c.handle)
	return c
}

func (c *ChatRecorder) handle(ctx context.Context, data []byte) error {
	c.buffer.Write(data)
	bufferContent := c.buffer.String()

	lines := strings.Split(bufferContent, "\n")
	if len(lines) > 0 {
		lastLine := lines[len(lines)-1]
		lines = lines[:len(lines)-1]
		c.buffer.Reset()
		c.buffer.WriteString(lastLine)
	}

	for _, line := range lines {
		if err := c.processSSELine(ctx, line); err != nil {
			return err
		}
	}

	return nil
}

func (c *ChatRecorder) processSSELine(ctx context.Context, line string) error {
	line = strings.TrimSpace(line)

	if line == "" || !strings.HasPrefix(line, "data:") {
		return nil
	}

	dataContent := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if dataContent == "" {
		return nil
	}

	if dataContent == "[DONE]" {
		c.processCompletedChat(ctx)
		return nil
	}

	var resp openai.ChatCompletionStreamResponse
	if err := json.Unmarshal([]byte(dataContent), &resp); err != nil {
		c.logger.With("data", dataContent).WarnContext(ctx, "解析流式响应失败", "error", err)
		return nil
	}

	prompt := c.getPrompt(ctx, c.req)
	mode := c.req.Metadata["mode"]
	taskID := c.req.Metadata["task_id"]

	rc := &domain.RecordParam{
		RequestID: c.cx.RequestID,
		TaskID:    taskID,
		UserID:    c.cx.UserID,
		ModelID:   c.model.ID,
		ModelType: c.model.ModelType,
		WorkMode:  mode,
		Prompt:    prompt,
		Role:      consts.ChatRoleAssistant,
	}
	if resp.Usage != nil {
		rc.InputTokens = int64(resp.Usage.PromptTokens)
	}

	if rc.Prompt != "" {
		urc := rc.Clone()
		urc.Role = consts.ChatRoleUser
		urc.Completion = urc.Prompt
		if err := c.usecase.Record(context.Background(), urc); err != nil {
			c.logger.With("modelID", c.model.ID, "modelName", c.model.ModelName, "modelType", consts.ModelTypeLLM).
				WarnContext(ctx, "插入流式记录失败", "error", err)
		}
	}

	if len(resp.Choices) > 0 {
		content := resp.Choices[0].Delta.Content
		if content != "" {
			c.completion.WriteString(content)
		}
	}

	if resp.Usage != nil {
		c.usage = resp.Usage
	}

	return nil
}

func (c *ChatRecorder) processCompletedChat(ctx context.Context) {
	if c.recorded {
		return // 避免重复记录
	}

	mode := c.req.Metadata["mode"]
	taskID := c.req.Metadata["task_id"]

	rc := &domain.RecordParam{
		RequestID:    c.cx.RequestID,
		TaskID:       taskID,
		UserID:       c.cx.UserID,
		ModelID:      c.model.ID,
		ModelType:    c.model.ModelType,
		WorkMode:     mode,
		Role:         consts.ChatRoleAssistant,
		Completion:   c.completion.String(),
		InputTokens:  int64(c.usage.PromptTokens),
		OutputTokens: int64(c.usage.CompletionTokens),
	}

	if err := c.usecase.Record(context.Background(), rc); err != nil {
		c.logger.With("modelID", c.model.ID, "modelName", c.model.ModelName, "modelType", consts.ModelTypeLLM).
			WarnContext(ctx, "插入流式记录失败", "error", err)
	} else {
		c.recorded = true
		c.logger.With("requestID", c.cx.RequestID, "completion_length", len(c.completion.String())).
			InfoContext(ctx, "流式对话记录已保存")
	}
}

// Close 关闭 recorder 并确保数据被保存
func (c *ChatRecorder) Close() {
	// 如果有累积的内容但还没有记录，强制保存
	if !c.recorded && c.completion.Len() > 0 {
		c.logger.With("requestID", c.cx.RequestID).
			WarnContext(c.ctx, "数据流异常中断，强制保存已累积的内容")
		c.processCompletedChat(c.ctx)
	}

	if c.Tee != nil {
		c.Tee.Close()
	}
}

func (c *ChatRecorder) Reset() {
	c.completion.Reset()
	c.buffer.Reset()
	c.usage = nil
	c.recorded = false
}

func (c *ChatRecorder) getPrompt(ctx context.Context, req *openai.ChatCompletionRequest) string {
	prompt := ""
	parse := promptparser.New(promptparser.KindTask)
	for _, message := range req.Messages {
		if message.Role == "system" {
			continue
		}

		if strings.Contains(message.Content, "<task>") ||
			strings.Contains(message.Content, "<feedback>") ||
			strings.Contains(message.Content, "<user_message>") {
			if info, err := parse.Parse(message.Content); err == nil {
				prompt = info.Prompt
			} else {
				c.logger.With("message", message.Content).WarnContext(ctx, "解析Prompt失败", "error", err)
			}
		}

		for _, m := range message.MultiContent {
			if strings.Contains(m.Text, "<task>") ||
				strings.Contains(m.Text, "<feedback>") ||
				strings.Contains(m.Text, "<user_message>") {
				if info, err := parse.Parse(m.Text); err == nil {
					prompt = info.Prompt
				} else {
					c.logger.With("message", m.Text).WarnContext(ctx, "解析Prompt失败", "error", err)
				}
			}
		}
	}
	return prompt
}
