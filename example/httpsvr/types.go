package main

type UserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type User struct {
	UserReq
	ID string `json:"id"`
}

type UserRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type UserLoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	UserRes
	AccessToken string `json:"accessToken"`
}
