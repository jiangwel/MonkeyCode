package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/ent/types"
	"github.com/chaitin/MonkeyCode/backend/errcode"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
	"github.com/chaitin/MonkeyCode/backend/pkg/oauth"
)

type UserUsecase struct {
	cfg    *config.Config
	redis  *redis.Client
	repo   domain.UserRepo
	logger *slog.Logger
}

func NewUserUsecase(
	cfg *config.Config,
	redis *redis.Client,
	repo domain.UserRepo,
	logger *slog.Logger,
) domain.UserUsecase {
	u := &UserUsecase{
		cfg:    cfg,
		redis:  redis,
		repo:   repo,
		logger: logger,
	}
	return u
}

func (u *UserUsecase) InitAdmin(ctx context.Context) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("generate admin password", "error", err)
		return err
	}
	return u.repo.InitAdmin(ctx, u.cfg.Admin.User, string(hash))
}

func (u *UserUsecase) List(ctx context.Context, req domain.ListReq) (*domain.ListUserResp, error) {
	users, p, err := u.repo.List(ctx, &req.Pagination)
	if err != nil {
		return nil, err
	}

	ids := cvt.Iter(users, func(_ int, u *db.User) string { return u.ID.String() })
	m, err := u.getUserActive(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &domain.ListUserResp{
		PageInfo: p,
		Users: cvt.Iter(users, func(_ int, e *db.User) *domain.User {
			return cvt.From(e, &domain.User{
				LastActiveAt: m[e.ID.String()],
			})
		}),
	}, nil
}

func (u *UserUsecase) getUserActive(ctx context.Context, ids []string) (map[string]int64, error) {
	m := make(map[string]int64)
	for _, id := range ids {
		key := fmt.Sprintf(consts.UserActiveKeyFmt, id)
		if t, err := u.redis.Get(ctx, key).Int64(); err != nil && !errors.Is(err, redis.Nil) {
			u.logger.With("key", key).With("error", err).Warn("get user active time failed")
		} else {
			m[id] = t
		}
	}

	return m, nil
}

// AdminList implements domain.UserUsecase.
func (u *UserUsecase) AdminList(ctx context.Context, page *web.Pagination) (*domain.ListAdminUserResp, error) {
	admins, p, err := u.repo.AdminList(ctx, page)
	if err != nil {
		return nil, err
	}

	ids := cvt.Iter(admins, func(_ int, u *db.Admin) string { return u.ID.String() })
	m, err := u.getAdminActive(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &domain.ListAdminUserResp{
		PageInfo: p,
		Users: cvt.Iter(admins, func(_ int, e *db.Admin) *domain.AdminUser {
			return cvt.From(e, &domain.AdminUser{
				LastActiveAt: m[e.ID.String()],
			})
		}),
	}, nil
}

func (u *UserUsecase) getAdminActive(ctx context.Context, ids []string) (map[string]int64, error) {
	m := make(map[string]int64)
	for _, id := range ids {
		key := fmt.Sprintf(consts.AdminActiveKeyFmt, id)
		if t, err := u.redis.Get(ctx, key).Int64(); err != nil && !errors.Is(err, redis.Nil) {
			u.logger.With("key", key).With("error", err).Warn("get admin active time failed")
		} else {
			m[id] = t
		}
	}

	return m, nil
}

// AdminLoginHistory implements domain.UserUsecase.
func (u *UserUsecase) AdminLoginHistory(ctx context.Context, page *web.Pagination) (*domain.ListAdminLoginHistoryResp, error) {
	histories, p, err := u.repo.AdminLoginHistory(ctx, page)
	if err != nil {
		return nil, err
	}

	return &domain.ListAdminLoginHistoryResp{
		PageInfo: p,
		LoginHistories: cvt.Iter(histories, func(_ int, e *db.AdminLoginHistory) *domain.AdminLoginHistory {
			return cvt.From(e, &domain.AdminLoginHistory{})
		}),
	}, nil
}

// Invite implements domain.UserUsecase.
func (u *UserUsecase) Invite(ctx context.Context, userID string) (*domain.InviteResp, error) {
	icode := uuid.NewString()
	code, err := u.repo.CreateInviteCode(ctx, userID, icode)
	if err != nil {
		return nil, err
	}
	return &domain.InviteResp{
		Code: code.Code,
	}, nil
}

// LoginHistory implements domain.UserUsecase.
func (u *UserUsecase) LoginHistory(ctx context.Context, page *web.Pagination) (*domain.ListLoginHistoryResp, error) {
	histories, p, err := u.repo.UserLoginHistory(ctx, page)
	if err != nil {
		return nil, err
	}

	return &domain.ListLoginHistoryResp{
		PageInfo: p,
		LoginHistories: cvt.Iter(histories, func(_ int, e *db.UserLoginHistory) *domain.UserLoginHistory {
			return cvt.From(e, &domain.UserLoginHistory{}).From(e)
		}),
	}, nil
}

// Register implements domain.UserUsecase.
func (u *UserUsecase) Register(ctx context.Context, req *domain.RegisterReq) (*domain.User, error) {
	if _, err := u.repo.ValidateInviteCode(ctx, req.Code); err != nil {
		return nil, errcode.ErrInviteCodeInvalid.Wrap(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &db.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Status:   consts.UserStatusActive,
		Platform: consts.UserPlatformEmail,
	}
	n, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return cvt.From(n, &domain.User{}), nil
}

func (u *UserUsecase) Login(ctx context.Context, req *domain.LoginReq) (*domain.LoginResp, error) {
	user, err := u.repo.GetByName(ctx, req.Username)
	if err != nil {
		return nil, errcode.ErrUserNotFound.Wrap(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errcode.ErrPassword.Wrap(err)
	}

	apiKey, err := u.repo.GetOrCreateApiKey(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	r, session, err := u.getVSCodeURL(ctx, req.SessionID, apiKey.Key, user.Username)
	if err != nil {
		return nil, err
	}

	if err := u.repo.SaveUserLoginHistory(ctx, user.ID.String(), req.IP, session); err != nil {
		u.logger.With("error", err).Error("save user login history")
	}
	return &domain.LoginResp{
		RedirectURL: r,
	}, nil
}

func (u *UserUsecase) getVSCodeURL(ctx context.Context, sessionID, apiKey, username string) (string, *domain.VSCodeSession, error) {
	s, err := u.redis.Get(ctx, fmt.Sprintf("vscode:session:%s", sessionID)).Result()
	if err != nil {
		return "", nil, err
	}
	session := &domain.VSCodeSession{}
	if err := json.Unmarshal([]byte(s), session); err != nil {
		return "", nil, err
	}
	r := fmt.Sprintf("%s?state=%s&api_key=%s&expires_in=3600&username=%s", session.RedirectURI, session.State, apiKey, username)
	return r, session, nil
}

func (u *UserUsecase) AdminLogin(ctx context.Context, req *domain.LoginReq) (*domain.AdminUser, error) {
	admin, err := u.repo.AdminByName(ctx, req.Username)
	if err != nil {
		return nil, errcode.ErrUserNotFound.Wrap(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, errcode.ErrPassword.Wrap(err)
	}

	if err := u.repo.SaveAdminLoginHistory(ctx, admin.ID.String(), req.IP); err != nil {
		return nil, err
	}
	return cvt.From(admin, &domain.AdminUser{}), nil
}

func (u *UserUsecase) VSCodeAuthInit(ctx context.Context, req *domain.VSCodeAuthInitReq) (*domain.VSCodeAuthInitResp, error) {
	i := uuid.NewString()
	session := &domain.VSCodeSession{
		ID:          i,
		State:       req.State,
		RedirectURI: req.RedirectURI,
		Version:     req.Version,
		OSType:      req.OSType,
		OSRelease:   req.OSRelease,
		Hostname:    req.Hostname,
	}

	b, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	if err := u.redis.Set(ctx, fmt.Sprintf("vscode:session:%s", i), b, 15*time.Minute).Err(); err != nil {
		return nil, err
	}
	return &domain.VSCodeAuthInitResp{
		AuthURL: fmt.Sprintf("%s?session_id=%s", u.cfg.BaseUrl+"/auth", i),
	}, nil
}

func (u *UserUsecase) CreateAdmin(ctx context.Context, req *domain.CreateAdminReq) (*domain.AdminUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin := &db.Admin{
		Username: req.Username,
		Password: string(hash),
	}
	n, err := u.repo.CreateAdmin(ctx, admin)
	if err != nil {
		return nil, err
	}
	return &domain.AdminUser{
		ID:        n.ID.String(),
		Username:  n.Username,
		CreatedAt: n.CreatedAt.Unix(),
	}, nil
}

func (u *UserUsecase) GetSetting(ctx context.Context) (*domain.Setting, error) {
	s, err := u.repo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}
	return cvt.From(s, &domain.Setting{}), nil
}

func (u *UserUsecase) UpdateSetting(ctx context.Context, req *domain.UpdateSettingReq) (*domain.Setting, error) {
	s, err := u.repo.UpdateSetting(ctx, func(old *db.Setting, up *db.SettingUpdateOne) {
		if req.EnableSSO != nil {
			up.SetEnableSSO(*req.EnableSSO)
		}
		if req.ForceTwoFactorAuth != nil {
			up.SetForceTwoFactorAuth(*req.ForceTwoFactorAuth)
		}
		if req.DisablePasswordLogin != nil {
			up.SetDisablePasswordLogin(*req.DisablePasswordLogin)
		}
		if req.EnableAutoLogin != nil {
			up.SetEnableAutoLogin(*req.EnableAutoLogin)
		}
		if req.DingtalkOAuth != nil {
			dingtalk := cvt.NilWithDefault(old.DingtalkOauth, &types.DingtalkOAuth{})
			if req.DingtalkOAuth.Enable != nil {
				dingtalk.Enable = *req.DingtalkOAuth.Enable
			}
			if req.DingtalkOAuth.ClientID != nil {
				dingtalk.ClientID = *req.DingtalkOAuth.ClientID
			}
			if req.DingtalkOAuth.ClientSecret != nil {
				dingtalk.ClientSecret = *req.DingtalkOAuth.ClientSecret
			}
			up.SetDingtalkOauth(dingtalk)
		}
		if req.CustomOAuth != nil {
			custom := cvt.NilWithDefault(old.CustomOauth, &types.CustomOAuth{})
			if req.CustomOAuth.Enable != nil {
				custom.Enable = *req.CustomOAuth.Enable
			}
			if req.CustomOAuth.ClientID != nil {
				custom.ClientID = *req.CustomOAuth.ClientID
			}
			if req.CustomOAuth.ClientSecret != nil {
				custom.ClientSecret = *req.CustomOAuth.ClientSecret
			}
			if req.CustomOAuth.AuthorizeURL != nil {
				custom.AuthorizeURL = *req.CustomOAuth.AuthorizeURL
			}
			if req.CustomOAuth.AccessTokenURL != nil {
				custom.AccessTokenURL = *req.CustomOAuth.AccessTokenURL
			}
			if req.CustomOAuth.UserInfoURL != nil {
				custom.UserInfoURL = *req.CustomOAuth.UserInfoURL
			}
			if req.CustomOAuth.Scopes != nil {
				custom.Scopes = req.CustomOAuth.Scopes
			}
			if req.CustomOAuth.IDField != nil {
				custom.IDField = *req.CustomOAuth.IDField
			}
			if req.CustomOAuth.NameField != nil {
				custom.NameField = *req.CustomOAuth.NameField
			}
			if req.CustomOAuth.AvatarField != nil {
				custom.AvatarField = *req.CustomOAuth.AvatarField
			}
			if req.CustomOAuth.EmailField != nil {
				custom.EmailField = *req.CustomOAuth.EmailField
			}
			up.SetCustomOauth(custom)
		}
	})
	if err != nil {
		return nil, err
	}
	return cvt.From(s, &domain.Setting{}), nil
}

func (u *UserUsecase) Update(ctx context.Context, req *domain.UpdateUserReq) (*domain.User, error) {
	user, err := u.repo.Update(ctx, req.ID, func(u *db.UserUpdateOne) error {
		if req.Status != nil {
			u.SetStatus(*req.Status)
		}
		if req.Password != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.SetPassword(string(hash))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cvt.From(user, &domain.User{}), nil
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *UserUsecase) DeleteAdmin(ctx context.Context, id string) error {
	return u.repo.DeleteAdmin(ctx, id)
}

func (u *UserUsecase) getOAuthConfig(setting *db.Setting, platform consts.UserPlatform) (*domain.OAuthConfig, error) {
	cfg := domain.OAuthConfig{
		Debug:       u.cfg.Debug,
		Platform:    platform,
		RedirectURI: fmt.Sprintf("%s/api/v1/user/oauth/callback", u.cfg.BaseUrl),
	}

	switch platform {
	case consts.UserPlatformDingTalk:
		if setting.DingtalkOauth == nil || !setting.DingtalkOauth.Enable {
			return nil, errcode.ErrDingtalkNotEnabled
		}
		cfg.ClientID = setting.DingtalkOauth.ClientID
		cfg.ClientSecret = setting.DingtalkOauth.ClientSecret
	case consts.UserPlatformCustom:
		if setting.CustomOauth == nil || !setting.CustomOauth.Enable {
			return nil, errcode.ErrCustomNotEnabled
		}
		cfg.ClientID = setting.CustomOauth.ClientID
		cfg.ClientSecret = setting.CustomOauth.ClientSecret
		cfg.AuthorizeURL = setting.CustomOauth.AuthorizeURL
		cfg.Scopes = setting.CustomOauth.Scopes
		cfg.TokenURL = setting.CustomOauth.AccessTokenURL
		cfg.UserInfoURL = setting.CustomOauth.UserInfoURL
		cfg.IDField = setting.CustomOauth.IDField
		cfg.NameField = setting.CustomOauth.NameField
		cfg.AvatarField = setting.CustomOauth.AvatarField
		cfg.EmailField = setting.CustomOauth.EmailField
	default:
		return nil, errcode.ErrUnsupportedPlatform
	}

	return &cfg, nil
}

func (u *UserUsecase) OAuthSignUpOrIn(ctx context.Context, req *domain.OAuthSignUpOrInReq) (*domain.OAuthURLResp, error) {
	setting, err := u.repo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := u.getOAuthConfig(setting, req.Platform)
	if err != nil {
		return nil, err
	}

	oauth, err := oauth.NewOAuther(*cfg)
	if err != nil {
		return nil, err
	}
	state, url := oauth.GetAuthorizeURL()

	session := &domain.OAuthState{
		SessionID:   req.SessionID,
		Kind:        req.OAuthKind(),
		Platform:    req.Platform,
		RedirectURL: req.RedirectURL,
		InviteCode:  req.InviteCode,
	}
	b, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	if err := u.redis.Set(ctx, fmt.Sprintf("oauth:state:%s", state), b, 15*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &domain.OAuthURLResp{
		URL: url,
	}, nil
}

func (u *UserUsecase) OAuthCallback(ctx context.Context, req *domain.OAuthCallbackReq) (string, error) {
	b, err := u.redis.Get(ctx, fmt.Sprintf("oauth:state:%s", req.State)).Result()
	if err != nil {
		return "", err
	}
	var session domain.OAuthState
	if err := json.Unmarshal([]byte(b), &session); err != nil {
		return "", err
	}

	switch session.Kind {
	case consts.OAuthKindInvite:
		return u.WithOAuthCallback(ctx, req, &session, func(ctx context.Context, s *domain.OAuthState, oui *domain.OAuthUserInfo) (*db.User, error) {
			return u.repo.OAuthRegister(ctx, s.Platform, s.InviteCode, oui)
		})

	case consts.OAuthKindLogin:
		setting, err := u.repo.GetSetting(ctx)
		if err != nil {
			return "", err
		}
		return u.WithOAuthCallback(ctx, req, &session, func(ctx context.Context, s *domain.OAuthState, oui *domain.OAuthUserInfo) (*db.User, error) {
			if setting.EnableAutoLogin {
				return u.repo.SignUpOrIn(ctx, s.Platform, oui)
			}
			return u.repo.OAuthLogin(ctx, s.Platform, oui)
		})
	default:
		return "", errcode.ErrOAuthStateInvalid
	}
}

func (u *UserUsecase) FetchUserInfo(ctx context.Context, req *domain.OAuthCallbackReq, session *domain.OAuthState) (*domain.OAuthUserInfo, error) {
	setting, err := u.repo.GetSetting(ctx)
	if err != nil {
		return nil, err
	}

	cfg, err := u.getOAuthConfig(setting, session.Platform)
	if err != nil {
		return nil, err
	}

	oauth, err := oauth.NewOAuther(*cfg)
	if err != nil {
		return nil, err
	}
	userInfo, err := oauth.GetUserInfo(req.Code)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

type OAuthUserRepoHandle func(context.Context, *domain.OAuthState, *domain.OAuthUserInfo) (*db.User, error)

func (u *UserUsecase) WithOAuthCallback(ctx context.Context, req *domain.OAuthCallbackReq, session *domain.OAuthState, handle OAuthUserRepoHandle) (string, error) {
	info, err := u.FetchUserInfo(ctx, req, session)
	if err != nil {
		return "", err
	}

	user, err := handle(ctx, session, info)
	if err != nil {
		return "", err
	}

	apiKey, err := u.repo.GetOrCreateApiKey(ctx, user.ID.String())
	if err != nil {
		return "", err
	}

	redirect := session.RedirectURL

	if session.SessionID != "" {
		r, vsess, err := u.getVSCodeURL(ctx, session.SessionID, apiKey.Key, user.Username)
		if err != nil {
			return "", err
		}
		redirect = fmt.Sprintf("%s?redirect_url=%s", redirect, url.QueryEscape(r))
		if err := u.repo.SaveUserLoginHistory(ctx, user.ID.String(), req.IP, vsess); err != nil {
			u.logger.With("error", err).Error("save user login history")
		}
	}

	u.logger.Debug("oauth callback", "redirect", redirect)
	return redirect, nil
}
