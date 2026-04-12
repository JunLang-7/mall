package dto

type HelloReq struct {
	Name string `json:"name"`
}

type HelloResp struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}
