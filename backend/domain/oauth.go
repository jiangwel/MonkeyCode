package domain

import "github.com/chaitin/MonkeyCode/backend/consts"

type OAuther interface {
	GetAuthorizeURL() (state string, url string)
	GetUserInfo(code string) (*OAuthUserInfo, error)
}

type OAuthConfig struct {
	Debug        bool
	Platform     consts.UserPlatform
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type OAuthUserInfo struct {
	ID        string `json:"id"`
	UnionID   string `json:"union_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type OAuthSignUpOrInReq struct {
	Platform    consts.UserPlatform `json:"platform" query:"platform" validate:"required"` // 第三方平台 dingtalk
	SessionID   string              `json:"session_id" query:"session_id"`                 // 会话ID
	RedirectURL string              `json:"redirect_url" query:"redirect_url"`             // 登录成功后跳转的 URL
}

type OAuthCallbackReq struct {
	State string `json:"state" query:"state" validate:"required"`
	Code  string `json:"code" query:"code" validate:"required"`
}

type OAuthURLResp struct {
	URL string `json:"url"`
}

type OAuthState struct {
	Kind        consts.OAuthKind    `json:"kind" query:"kind" validate:"required"`         // 注册或登录
	SessionID   string              `json:"session_id"`                                    // 会话ID
	Platform    consts.UserPlatform `json:"platform" query:"platform" validate:"required"` // 第三方平台 dingtalk
	RedirectURL string              `json:"redirect_url" query:"redirect_url"`             // 登录成功后跳转的 URL
}
