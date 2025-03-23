package domain

type TokenPayload struct {
	Email  string   `json:"email"`
	UserID int      `json:"user_id"`
	Role   UserRole `json:"role"`
}

type RefreshTokenPayload struct {
	UserID uint64 `json:"user_id"`
}
