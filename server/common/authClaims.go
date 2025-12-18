package common

import "github.com/golang-jwt/jwt/v5"

type StackAuthClaims struct {
	ProjectId      string `json:"project_id"`
	BranchId       string `json:"branch_id"`
	RefreshTokenId string `json:"refresh_token_id"`
	Role           string `json:"role"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	SelectedTeamId string `json:"selected_team_id"`
	IsAnonymous    bool   `json:"is_anonymous"`
	jwt.RegisteredClaims
}
