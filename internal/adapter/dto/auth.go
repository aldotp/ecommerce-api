package dto

func NewAuthResponse(token string) AuthResponse {
	return AuthResponse{
		AccessToken: token,
	}
}

// authResponse represents an authentication response body
type AuthResponse struct {
	AccessToken string `json:"access_token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
