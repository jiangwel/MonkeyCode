package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/db/model"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
	"github.com/chaitin/MonkeyCode/backend/pkg/request"
)

type ModelUsecase struct {
	logger *slog.Logger
	repo   domain.ModelRepo
	cfg    *config.Config
	client *http.Client
}

func NewModelUsecase(
	logger *slog.Logger,
	repo domain.ModelRepo,
	cfg *config.Config,
) domain.ModelUsecase {
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 30,
		},
	}
	return &ModelUsecase{repo: repo, cfg: cfg, logger: logger, client: client}
}

func (m *ModelUsecase) Check(ctx context.Context, req *domain.CheckModelReq) (*domain.Model, error) {
	if req.Type == consts.ModelTypeEmbedding || req.Type == consts.ModelTypeReranker {
		url := req.APIBase
		reqBody := map[string]any{}
		if req.Type == consts.ModelTypeEmbedding {
			reqBody = map[string]any{
				"model":           req.ModelName,
				"input":           "MonkeyCode 是一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				"encoding_format": "float",
			}
			url = req.APIBase + "/embeddings"
		}
		if req.Type == consts.ModelTypeReranker {
			reqBody = map[string]any{
				"model": req.ModelName,
				"documents": []string{
					"MonkeyCode 是一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"MonkeyCode 是一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
					"MonkeyCode 是一个基于大模型的代码生成器，它可以根据用户的需求生成代码。",
				},
				"query": "MonkeyCode",
			}
			url = req.APIBase + "/rerank"
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
		request.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("request failed: %s", resp.Status)
		}
		return &domain.Model{}, nil
	}
	config := &openai.ChatModelConfig{
		APIKey:  req.APIKey,
		BaseURL: req.APIBase,
		Model:   string(req.ModelName),
	}
	// for azure openai
	if req.Provider == consts.ModelProviderAzureOpenAI {
		config.ByAzure = true
		config.APIVersion = req.APIVersion
		if config.APIVersion == "" {
			config.APIVersion = "2024-10-21"
		}
	}
	if req.APIHeader != "" {
		client := getHttpClientWithAPIHeaderMap(req.APIHeader)
		if client != nil {
			config.HTTPClient = client
		}
	}
	chatModel, err := openai.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	resp, err := chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are a helpful assistant."),
		schema.UserMessage("hi"),
	})
	if err != nil {
		return nil, err
	}
	content := resp.Content
	if content == "" {
		return nil, fmt.Errorf("generate failed")
	}
	return &domain.Model{
		ModelType: req.Type,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		APIBase:   req.APIBase,
	}, nil
}

type headerTransport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.base.RoundTrip(req)
}

func getHttpClientWithAPIHeaderMap(header string) *http.Client {
	headerMap := request.GetHeaderMap(header)
	if len(headerMap) > 0 {
		// create http client with custom transport for headers
		client := &http.Client{
			Timeout: 0,
		}
		// Wrap the transport to add headers
		client.Transport = &headerTransport{
			headers: headerMap,
			base:    http.DefaultTransport,
		}
		return client
	}
	return nil
}

func (m *ModelUsecase) MyModelList(ctx context.Context, req *domain.MyModelListReq) ([]*domain.Model, error) {
	models, err := m.repo.MyModelList(ctx, req)
	if err != nil {
		return nil, err
	}
	ids := cvt.Iter(models, func(_ int, e *db.Model) uuid.UUID {
		return e.ID
	})
	usages, err := m.repo.ModelUsage(ctx, ids)
	if err != nil {
		return nil, err
	}
	return cvt.Iter(models, func(_ int, e *db.Model) *domain.Model {
		tmp := cvt.From(e, &domain.Model{}).From(e)
		if usage, ok := usages[e.ID]; ok {
			tmp.Input = usage.Input
			tmp.Output = usage.Output
		}
		return tmp
	}), nil
}

func (m *ModelUsecase) List(ctx context.Context) (*domain.AllModelResp, error) {
	return m.repo.List(ctx)
}

// Create implements domain.ModelUsecase.
func (m *ModelUsecase) Create(ctx context.Context, req *domain.CreateModelReq) (*domain.Model, error) {
	model, err := m.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return cvt.From(model, &domain.Model{}), nil
}

// GetTokenUsage implements domain.ModelUsecase.
func (m *ModelUsecase) GetTokenUsage(ctx context.Context, modelType consts.ModelType) (*domain.ModelTokenUsageResp, error) {
	return m.repo.GetTokenUsage(ctx, modelType)
}

// Update implements domain.ModelUsecase.
func (m *ModelUsecase) Update(ctx context.Context, req *domain.UpdateModelReq) (*domain.Model, error) {
	model, err := m.repo.Update(ctx, req.ID, func(tx *db.Tx, old *db.Model, up *db.ModelUpdateOne) error {
		if req.ModelName != nil {
			up.SetModelName(*req.ModelName)
		}
		if req.Provider != nil {
			up.SetProvider(*req.Provider)
		}
		if req.APIBase != nil {
			up.SetAPIBase(*req.APIBase)
		}
		if req.APIKey != nil {
			up.SetAPIKey(*req.APIKey)
		}
		if req.APIVersion != nil {
			up.SetAPIVersion(*req.APIVersion)
		}
		if req.APIHeader != nil {
			up.SetAPIHeader(*req.APIHeader)
		}
		if req.ShowName != nil {
			up.SetShowName(*req.ShowName)
		}
		if req.Status != nil {
			if *req.Status == consts.ModelStatusActive {
				if err := tx.Model.Update().
					Where(model.ModelType(old.ModelType)).
					SetStatus(consts.ModelStatusInactive).
					Exec(ctx); err != nil {
					return err
				}
			}
			up.SetStatus(*req.Status)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cvt.From(model, &domain.Model{}), nil
}

func (m *ModelUsecase) InitModel(ctx context.Context) error {
	m.logger.With("init_model", m.cfg.InitModel).Debug("init model")
	if m.cfg.InitModel.Name == "" {
		return nil
	}
	return m.repo.InitModel(ctx, m.cfg.InitModel.Name, m.cfg.InitModel.Key, m.cfg.InitModel.URL)
}

func (m *ModelUsecase) GetProviderModelList(ctx context.Context, req *domain.GetProviderModelListReq) (*domain.GetProviderModelListResp, error) {
	switch req.Provider {
	case consts.ModelProviderAzureOpenAI,
		consts.ModelProviderVolcengine:
		return &domain.GetProviderModelListResp{
			Models: domain.ModelProviderBrandModelsList[req.Provider],
		}, nil
	case consts.ModelProviderOpenAI,
		consts.ModelProviderHunyuan,
		consts.ModelProviderMoonshot,
		consts.ModelProviderDeepSeek,
		consts.ModelProviderBaiLian:
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "/models")
		client := request.NewClient(u.Scheme, u.Host, m.client.Timeout, request.WithClient(m.client))
		resp, err := request.Get[domain.OpenAIResp](client, u.Path, request.WithHeader(
			request.Header{
				"Authorization": fmt.Sprintf("Bearer %s", req.APIKey),
			},
		))
		if err != nil {
			return nil, err
		}

		return &domain.GetProviderModelListResp{
			Models: cvt.Iter(resp.Data, func(_ int, e *domain.OpenAIData) domain.ProviderModelListItem {
				return domain.ProviderModelListItem{
					Model: e.ID,
				}
			}),
		}, nil
	case consts.ModelProviderOllama:
		// get from ollama http://10.10.16.24:11434/api/tags
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = "/api/tags"
		client := request.NewClient(u.Scheme, u.Host, m.client.Timeout, request.WithClient(m.client))

		h := request.Header{}
		if req.APIHeader != "" {
			headers := request.GetHeaderMap(req.APIHeader)
			maps.Copy(h, headers)
		}

		return request.Get[domain.GetProviderModelListResp](client, u.Path, request.WithHeader(h))

	case consts.ModelProviderSiliconFlow, consts.ModelProviderBaiZhiCloud:
		if req.Type == consts.ModelTypeEmbedding || req.Type == consts.ModelTypeReranker {
			if req.Provider == consts.ModelProviderBaiZhiCloud {
				if req.Type == consts.ModelTypeEmbedding {
					return &domain.GetProviderModelListResp{
						Models: []domain.ProviderModelListItem{
							{
								Model: "bge-m3",
							},
						},
					}, nil
				} else {
					return &domain.GetProviderModelListResp{
						Models: []domain.ProviderModelListItem{
							{
								Model: "bge-reranker-v2-m3",
							},
						},
					}, nil
				}
			}
		}
		u, err := url.Parse(req.BaseURL)
		if err != nil {
			return nil, err
		}
		st := string(req.Type)
		if req.Type == consts.ModelTypeLLM {
			st = "chat"
		}
		client := request.NewClient(u.Scheme, u.Host, m.client.Timeout, request.WithClient(m.client))
		resp, err := request.Get[domain.OpenAIResp](client, "/v1/models", request.WithHeader(
			request.Header{
				"Authorization": fmt.Sprintf("Bearer %s", req.APIKey),
			},
		), request.WithQuery(request.Query{
			"type":     "text",
			"sub_type": st,
		}))
		if err != nil {
			return nil, err
		}

		return &domain.GetProviderModelListResp{
			Models: cvt.Iter(resp.Data, func(_ int, e *domain.OpenAIData) domain.ProviderModelListItem {
				return domain.ProviderModelListItem{
					Model: e.ID,
				}
			}),
		}, nil
	default:
		return nil, fmt.Errorf("invalid provider: %s", req.Provider)
	}
}

func (m *ModelUsecase) Delete(ctx context.Context, id string) error {
	return m.repo.Delete(ctx, id)
}
