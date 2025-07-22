package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/db/admin"
	"github.com/chaitin/MonkeyCode/backend/db/adminloginhistory"
	"github.com/chaitin/MonkeyCode/backend/db/apikey"
	"github.com/chaitin/MonkeyCode/backend/db/invitecode"
	"github.com/chaitin/MonkeyCode/backend/db/user"
	"github.com/chaitin/MonkeyCode/backend/db/useridentity"
	"github.com/chaitin/MonkeyCode/backend/db/userloginhistory"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/errcode"
	"github.com/chaitin/MonkeyCode/backend/pkg/entx"
	"github.com/chaitin/MonkeyCode/backend/pkg/ipdb"
)

type UserRepo struct {
	db    *db.Client
	ipdb  *ipdb.IPDB
	redis *redis.Client
}

func NewUserRepo(db *db.Client, ipdb *ipdb.IPDB, redis *redis.Client) domain.UserRepo {
	return &UserRepo{db: db, ipdb: ipdb, redis: redis}
}

func (r *UserRepo) InitAdmin(ctx context.Context, username, password string) error {
	_, err := r.AdminByName(ctx, username)
	if db.IsNotFound(err) {
		_, err = r.CreateAdmin(ctx, &db.Admin{
			Username: username,
			Password: password,
			Status:   consts.AdminStatusActive,
		})
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) CreateAdmin(ctx context.Context, admin *db.Admin) (*db.Admin, error) {
	return r.db.Admin.Create().
		SetUsername(admin.Username).
		SetPassword(admin.Password).
		SetStatus(admin.Status).
		Save(ctx)

}

func (r *UserRepo) AdminByName(ctx context.Context, username string) (*db.Admin, error) {
	return r.db.Admin.Query().Where(admin.Username(username)).Only(ctx)
}

func (r *UserRepo) GetByName(ctx context.Context, username string) (*db.User, error) {
	return r.db.User.Query().Where(
		user.Or(
			user.Username(username),
			user.Email(username),
		),
	).Only(ctx)
}

func (r *UserRepo) ValidateInviteCode(ctx context.Context, code string) (*db.InviteCode, error) {
	var res *db.InviteCode
	err := entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		ic, err := r.innerValidateInviteCode(ctx, tx, code)
		if err != nil {
			return err
		}
		res = ic
		return nil
	})
	return res, err
}

func (r *UserRepo) innerValidateInviteCode(ctx context.Context, tx *db.Tx, code string) (*db.InviteCode, error) {
	ic, err := tx.InviteCode.Query().Where(invitecode.Code(code)).Only(ctx)
	if err != nil {
		return nil, err
	}

	if ic.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("invite code has expired")
	}
	if ic.Status == consts.InviteCodeStatusUsed {
		return nil, errors.New("invite code has been used")
	}

	ic, err = tx.InviteCode.UpdateOneID(ic.ID).
		SetStatus(consts.InviteCodeStatusUsed).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return ic, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *db.User) (*db.User, error) {
	return r.db.User.Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetStatus(user.Status).
		SetPlatform(user.Platform).
		Save(ctx)
}

func (r *UserRepo) UserLoginHistory(ctx context.Context, page *web.Pagination) ([]*db.UserLoginHistory, *db.PageInfo, error) {
	ctx = entx.SkipSoftDelete(ctx)
	q := r.db.UserLoginHistory.Query().WithOwner().Order(userloginhistory.ByCreatedAt(sql.OrderDesc()))
	return q.Page(ctx, page.Page, page.Size)
}

func (r *UserRepo) AdminLoginHistory(ctx context.Context, page *web.Pagination) ([]*db.AdminLoginHistory, *db.PageInfo, error) {
	q := r.db.AdminLoginHistory.Query().WithOwner().Order(adminloginhistory.ByCreatedAt(sql.OrderDesc()))
	return q.Page(ctx, page.Page, page.Size)
}

func (r *UserRepo) CreateInviteCode(ctx context.Context, userID string, code string) (*db.InviteCode, error) {
	adminID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return r.db.InviteCode.Create().
		SetAdminID(adminID).
		SetCode(code).
		SetStatus(consts.InviteCodeStatusPending).
		SetExpiredAt(time.Now().Add(15 * time.Minute)).
		Save(ctx)
}

func (r *UserRepo) AdminList(ctx context.Context, page *web.Pagination) ([]*db.Admin, *db.PageInfo, error) {
	q := r.db.Admin.Query()
	return q.Page(ctx, page.Page, page.Size)
}

func (r *UserRepo) List(ctx context.Context, page *web.Pagination) ([]*db.User, *db.PageInfo, error) {
	q := r.db.User.Query()
	return q.Page(ctx, page.Page, page.Size)
}

func (r *UserRepo) GetOrCreateApiKey(ctx context.Context, userID string) (*db.ApiKey, error) {
	i, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var apiKey *db.ApiKey
	err = entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		k, err := tx.ApiKey.Query().Where(apikey.UserID(i)).First(ctx)
		if db.IsNotFound(err) {
			n, err := tx.ApiKey.Create().
				SetUserID(i).
				SetKey(uuid.NewString()).
				SetName("default").
				SetStatus(consts.ApiKeyStatusActive).
				Save(ctx)
			if err != nil {
				return err
			}
			apiKey = n
			return nil
		}
		if err != nil {
			return err
		}
		apiKey = k
		return nil
	})
	return apiKey, err
}

func (r *UserRepo) GetSetting(ctx context.Context) (*db.Setting, error) {
	s, err := r.db.Setting.Query().First(ctx)
	if db.IsNotFound(err) {
		return r.db.Setting.Create().
			SetEnableSSO(false).
			SetForceTwoFactorAuth(false).
			SetDisablePasswordLogin(false).
			Save(ctx)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *UserRepo) UpdateSetting(ctx context.Context, fn func(*db.Setting, *db.SettingUpdateOne)) (*db.Setting, error) {
	var res *db.Setting
	err := entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		s, err := tx.Setting.Query().First(ctx)
		if err != nil {
			return err
		}
		up := tx.Setting.UpdateOneID(s.ID)
		fn(s, up)
		s, err = up.Save(ctx)
		if err != nil {
			return err
		}
		res = s
		return nil
	})
	return res, err
}

func (r *UserRepo) Update(ctx context.Context, id string, fn func(*db.User, *db.UserUpdateOne) error) (*db.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var u *db.User
	err = entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		u, err = tx.User.Query().Where(user.ID(uid)).Only(ctx)
		if err != nil {
			return err
		}
		if err := fn(u, u.Update()); err != nil {
			return err
		}
		return u.Update().Exec(ctx)
	})
	return u, err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		user, err := tx.User.Query().
			WithIdentities().
			Where(user.ID(uid)).
			Only(ctx)
		if err != nil {
			return err
		}

		keys, err := tx.ApiKey.Query().Where(apikey.UserID(user.ID)).All(ctx)
		if err != nil {
			return err
		}

		for _, v := range keys {
			if _, err := tx.ApiKey.Delete().Where(apikey.ID(v.ID)).Exec(ctx); err != nil {
				return err
			}
			if err := r.redis.Del(ctx, fmt.Sprintf("sk-%s", v.Key)).Err(); err != nil {
				return err
			}
		}

		for _, v := range user.Edges.Identities {
			if _, err := tx.UserIdentity.Delete().Where(useridentity.ID(v.ID)).Exec(ctx); err != nil {
				return err
			}
		}
		return tx.User.DeleteOneID(uid).Exec(ctx)
	})
}

func (r *UserRepo) DeleteAdmin(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	admin, err := r.db.Admin.Get(ctx, uid)
	if err != nil {
		return err
	}
	if admin.Username == "admin" {
		return errors.New("admin cannot be deleted")
	}
	return r.db.Admin.DeleteOne(admin).Exec(ctx)
}

func (r *UserRepo) OAuthRegister(ctx context.Context, platform consts.UserPlatform, inviteCode string, req *domain.OAuthUserInfo) (*db.User, error) {
	var u *db.User
	err := entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		if _, err := r.innerValidateInviteCode(ctx, tx, inviteCode); err != nil {
			return errcode.ErrInviteCodeInvalid.Wrap(err)
		}

		_, err := tx.UserIdentity.Query().
			WithUser().
			Where(useridentity.Platform(platform), useridentity.IdentityID(req.ID)).
			First(ctx)
		if err == nil {
			return fmt.Errorf("user already exists for platform %s and identity ID %s", platform, req.ID)
		}
		if !db.IsNotFound(err) {
			return err
		}
		user, err := tx.User.Create().
			SetUsername(req.Name).
			SetEmail(req.Email).
			SetAvatarURL(req.AvatarURL).
			SetPlatform(platform).
			SetStatus(consts.UserStatusActive).
			Save(ctx)
		if err != nil {
			return err
		}
		_, err = tx.UserIdentity.Create().
			SetUserID(user.ID).
			SetPlatform(platform).
			SetIdentityID(req.ID).
			SetUnionID(req.UnionID).
			SetNickname(req.Name).
			SetAvatarURL(req.AvatarURL).
			SetEmail(req.Email).
			Save(ctx)
		if err != nil {
			return err
		}
		u = user
		return nil
	})
	return u, err
}

func (r *UserRepo) OAuthLogin(ctx context.Context, platform consts.UserPlatform, req *domain.OAuthUserInfo) (*db.User, error) {
	ui, err := r.db.UserIdentity.Query().
		WithUser().
		Where(useridentity.Platform(platform), useridentity.IdentityID(req.ID)).
		Where(useridentity.HasUser()).
		Only(ctx)
	if err != nil {
		return nil, errcode.ErrNotInvited.Wrap(err)
	}
	if ui.AvatarURL != req.AvatarURL {
		if err = entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
			return r.updateAvatar(ctx, tx, ui, req.AvatarURL)
		}); err != nil {
			return nil, err
		}
	}
	return ui.Edges.User, nil
}

func (r *UserRepo) updateAvatar(ctx context.Context, tx *db.Tx, ui *db.UserIdentity, avatar string) error {
	if err := tx.UserIdentity.UpdateOneID(ui.ID).SetAvatarURL(avatar).Exec(ctx); err != nil {
		return err
	}
	return tx.User.UpdateOneID(ui.UserID).SetAvatarURL(avatar).Exec(ctx)
}

func (r *UserRepo) SignUpOrIn(ctx context.Context, platform consts.UserPlatform, req *domain.OAuthUserInfo) (*db.User, error) {
	var u *db.User
	err := entx.WithTx(ctx, r.db, func(tx *db.Tx) error {
		ui, err := tx.UserIdentity.Query().
			WithUser().
			Where(useridentity.Platform(platform), useridentity.IdentityID(req.ID)).
			First(ctx)
		if err == nil {
			u = ui.Edges.User
			if ui.AvatarURL != req.AvatarURL {
				if err = r.updateAvatar(ctx, tx, ui, req.AvatarURL); err != nil {
					return err
				}
			}
			return nil
		}
		if !db.IsNotFound(err) {
			return err
		}
		user, err := tx.User.Create().
			SetUsername(req.Name).
			SetEmail(req.Email).
			SetAvatarURL(req.AvatarURL).
			SetPlatform(platform).
			SetStatus(consts.UserStatusActive).
			Save(ctx)
		if err != nil {
			return err
		}
		_, err = tx.UserIdentity.Create().
			SetUserID(user.ID).
			SetPlatform(platform).
			SetIdentityID(req.ID).
			SetUnionID(req.UnionID).
			SetNickname(req.Name).
			SetAvatarURL(req.AvatarURL).
			SetEmail(req.Email).
			Save(ctx)
		if err != nil {
			return err
		}
		u = user
		return nil
	})
	return u, err
}

func (r *UserRepo) SaveAdminLoginHistory(ctx context.Context, adminID string, ip string) error {
	uid, err := uuid.Parse(adminID)
	if err != nil {
		return err
	}
	addr, err := r.ipdb.Lookup(ip)
	if err != nil {
		return err
	}
	_, err = r.db.AdminLoginHistory.Create().
		SetAdminID(uid).
		SetIP(ip).
		SetCity(addr.City).
		SetCountry(addr.Country).
		SetProvince(addr.Province).
		Save(ctx)
	return err
}

func (r *UserRepo) SaveUserLoginHistory(ctx context.Context, userID string, ip string, session *domain.VSCodeSession) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	addr, err := r.ipdb.Lookup(ip)
	if err != nil {
		return err
	}
	c := r.db.UserLoginHistory.Create().
		SetUserID(uid).
		SetIP(ip).
		SetCity(addr.City).
		SetCountry(addr.Country).
		SetProvince(addr.Province)

	if session != nil {
		c.SetClientVersion(session.Version).
			SetOsType(session.OSType).
			SetOsRelease(session.OSRelease).
			SetClientID(session.ClientID).
			SetHostname(session.Hostname)
	}

	return c.Exec(ctx)
}
