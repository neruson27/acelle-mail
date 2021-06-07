package middlewares

import "context"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ctxUserID     = contextKey("session-user-id")
	ctxCompanyID  = contextKey("session-company-id")
	ctxPrivileges = contextKey("session-privileges")
)

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxUserID).(string)
	return userID, ok
}

func GetCompanyID(ctx context.Context) (string, bool) {
	companyID, ok := ctx.Value(ctxCompanyID).(string)
	return companyID, ok
}

func GetPrivileges(ctx context.Context) ([]string, bool) {
	privileges, ok := ctx.Value(ctxPrivileges).([]string)
	return privileges, ok
}

func SetUserID(ctx context.Context, claims *SessionClaims) context.Context {
	if claims.UserID == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxUserID, claims.UserID)
}

func SetCompanyID(ctx context.Context, claims *SessionClaims) context.Context {
	if claims.CompanyID == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxCompanyID, claims.CompanyID)
}

func SetPrivileges(ctx context.Context, claims *SessionClaims) context.Context {
	if len(claims.Privileges) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ctxPrivileges, claims.Privileges)
}

func ContextAppendValues(ctx context.Context, claims *SessionClaims, appendFunctions ...func(ctx2 context.Context, claims *SessionClaims) context.Context) context.Context {
	for _, fn := range appendFunctions {
		ctx = fn(ctx, claims)
	}
	return ctx
}
