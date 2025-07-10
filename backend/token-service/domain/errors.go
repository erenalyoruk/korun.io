package domain

import "korun.io/shared/models"

var (
	ErrInvalidAccessToken   = &models.Error{Code: "T001", Message: "invalid access token"}
	ErrInvalidRefreshToken  = &models.Error{Code: "T002", Message: "invalid refresh token"}
	ErrRefreshTokenNotFound = &models.Error{Code: "T003", Message: "refresh token not found"}
	ErrTokenExpired         = &models.Error{Code: "T004", Message: "token has expired"}
	ErrTokenAlreadyRevoked  = &models.Error{Code: "T005", Message: "token already revoked"}
)
