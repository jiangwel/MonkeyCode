package consts

const (
	UserActiveKeyFmt = "user:active:%s"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusLocked   UserStatus = "locked"
)

const (
	SessionName = "monkeycode_session"
)

type UserPlatform string

const (
	UserPlatformEmail    UserPlatform = "email"
	UserPlatformDingTalk UserPlatform = "dingtalk"
)

type OAuthKind string

const (
	OAuthKindSignUpOrIn OAuthKind = "signup_or_in"
)

type InviteCodeStatus string

const (
	InviteCodeStatusPending InviteCodeStatus = "pending"
	InviteCodeStatusUsed    InviteCodeStatus = "used"
)
