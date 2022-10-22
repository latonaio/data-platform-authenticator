package response

type JWTResponseFormat struct {
	Jwt string `json:"jwt"`
}

type UserResponseFormat struct {
	BusinessPartner string `json:"business_partner"`
	LoginID         string `json:"login_id"`
}
