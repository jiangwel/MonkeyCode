package repo

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/db/task"
	"github.com/chaitin/MonkeyCode/backend/db/taskrecord"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
)

type BillingRepo struct {
	db *db.Client
}

func NewBillingRepo(db *db.Client) domain.BillingRepo {
	return &BillingRepo{db: db}
}

// ChatInfo implements domain.BillingRepo.
func (b *BillingRepo) ChatInfo(ctx context.Context, id string) (*domain.ChatInfo, error) {
	record, err := b.db.Task.Query().
		WithTaskRecords(func(trq *db.TaskRecordQuery) {
			trq.Order(taskrecord.ByCreatedAt(sql.OrderAsc()))
		}).
		Where(task.TaskID(id)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return cvt.From(record, &domain.ChatInfo{}), nil
}

// CompletionInfo implements domain.BillingRepo.
func (b *BillingRepo) CompletionInfo(ctx context.Context, id string) (*domain.CompletionInfo, error) {
	record, err := b.db.Task.Query().Where(task.TaskID(id)).First(ctx)
	if err != nil {
		return nil, err
	}

	return cvt.From(record, &domain.CompletionInfo{}), nil
}

// ListChatRecord implements domain.BillingRepo.
func (b *BillingRepo) ListChatRecord(ctx context.Context, page *web.Pagination) (*domain.ListChatRecordResp, error) {
	q := b.db.Task.Query().
		WithUser().
		WithModel().
		WithTaskRecords().
		Where(task.ModelType(consts.ModelTypeLLM)).
		Order(task.ByCreatedAt(sql.OrderDesc()))

	records, p, err := q.Page(ctx, page.Page, page.Size)
	if err != nil {
		return nil, err
	}

	return &domain.ListChatRecordResp{
		PageInfo: p,
		Records: cvt.Iter(records, func(_ int, r *db.Task) *domain.ChatRecord {
			return cvt.From(r, &domain.ChatRecord{})
		}),
	}, nil
}

// ListCompletionRecord implements domain.BillingRepo.
func (b *BillingRepo) ListCompletionRecord(ctx context.Context, page *web.Pagination) (*domain.ListCompletionRecordResp, error) {
	q := b.db.Task.Query().
		WithUser().
		WithModel().
		Where(task.ModelType(consts.ModelTypeCoder)).
		Order(task.ByCreatedAt(sql.OrderDesc()))
	records, p, err := q.Page(ctx, page.Page, page.Size)
	if err != nil {
		return nil, err
	}

	return &domain.ListCompletionRecordResp{
		PageInfo: p,
		Records: cvt.Iter(records, func(_ int, r *db.Task) *domain.CompletionRecord {
			return cvt.From(r, &domain.CompletionRecord{})
		}),
	}, nil
}
