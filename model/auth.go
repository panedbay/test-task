package model

type AuthRequest struct {
	Password string `json:"password"`

	Username string `json:"username"`
}

type AuthResponse struct {
	Token string `json:"token,omitempty"`
}
