package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/db/apikey"
	"github.com/chaitin/MonkeyCode/backend/db/model"
	"github.com/chaitin/MonkeyCode/backend/db/task"
	"github.com/chaitin/MonkeyCode/backend/db/taskrecord"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/entx"
)

type ProxyRepo struct {
	db *db.Client
}

func NewProxyRepo(db *db.Client) domain.ProxyRepo {
	return &ProxyRepo{db: db}
}

func (r *ProxyRepo) SelectModelWithLoadBalancing(modelName string, modelType consts.ModelType) (*db.Model, error) {
	ctx := context.Background()
	m, err := r.db.Model.Query().
		Where(model.ModelName(modelName), model.ModelType(modelType)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *ProxyRepo) ValidateApiKey(ctx context.Context, key string) (*db.ApiKey, error) {
	a, err := r.db.ApiKey.Query().
		Where(apikey.Key(key), apikey.Status(consts.ApiKeyStatusActive)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *ProxyRepo) Record(ctx context.Context, record *domain.RecordParam) error {
	userID, err := uuid.Parse(record.UserID)
	if err != nil {
		return err
	}
	modelID, err := uuid.Parse(record.ModelID)
	if err != nil {
		return err
	}

	return entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		isNew := false
		t, err := tx.Task.Query().ForUpdate().Where(task.TaskID(record.TaskID)).Only(ctx)
		if db.IsNotFound(err) {
			t, err = tx.Task.Create().
				SetTaskID(record.TaskID).
				SetRequestID(record.RequestID).
				SetUserID(userID).
				SetModelID(modelID).
				SetProgramLanguage(record.ProgramLanguage).
				SetInputTokens(record.InputTokens).
				SetOutputTokens(record.OutputTokens).
				SetIsAccept(record.IsAccept).
				SetModelType(record.ModelType).
				SetWorkMode(record.WorkMode).
				SetCodeLines(record.CodeLines).
				Save(ctx)
			isNew = true
		}
		if err != nil {
			return err
		}
		if !isNew {
			up := tx.Task.UpdateOneID(t.ID)
			if record.OutputTokens > 0 {
				up.AddOutputTokens(record.OutputTokens)
			}
			if t.InputTokens == 0 && record.InputTokens > 0 {
				up.SetInputTokens(record.InputTokens)
			}
			if t.RequestID != record.RequestID {
				up.SetRequestID(record.RequestID)
				up.AddInputTokens(record.InputTokens)
			}
			if err := up.Exec(ctx); err != nil {
				return err
			}
		}

		if record.Role == consts.ChatRoleUser {
			count, err := tx.TaskRecord.Query().
				Where(
					taskrecord.TaskID(t.ID),
					taskrecord.Role(consts.ChatRoleUser),
					taskrecord.Prompt(record.Prompt),
				).
				Count(ctx)
			if err != nil {
				return err
			}
			if count > 0 {
				return nil
			}
		}

		_, err = tx.TaskRecord.Create().
			SetTaskID(t.ID).
			SetRole(record.Role).
			SetPrompt(record.Prompt).
			SetCompletion(record.Completion).
			SetOutputTokens(record.OutputTokens).
			Save(ctx)

		return err
	})
}

func (r *ProxyRepo) UpdateByTaskID(ctx context.Context, taskID string, fn func(*db.TaskUpdateOne)) error {
	rc, err := r.db.Task.Query().Where(task.TaskID(taskID)).Only(ctx)
	if err != nil {
		return err
	}

	return entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		up := tx.Task.UpdateOneID(rc.ID)
		fn(up)
		return up.Exec(ctx)
	})
}
