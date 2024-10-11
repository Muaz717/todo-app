package testmodels

import resp "github.com/Muaz717/todo-app/internal/lib/api/response"

type RegRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	resp.Response
	Token string `json:"token"`
}
