package response

type JWTResponseFormat struct {
	Jwt string `json:"jwt"`
}

type UserVerifyResponseFormat struct {
	EmailAddress string `json:"email_address"`
}

type UserDetailResponseFormat struct {
	EmailAddress          string `json:"email_address"`
	BusinessPartner       int    `json:"business_partner"`
	BusinessPartnerName   string `json:"business_partner_name"`
	BusinessUserFirstName string `json:"business_user_first_name"`
	BusinessUserLastName  string `json:"business_user_last_name"`
	BusinessUserFullName  string `json:"business_user_full_name"`
	Language              string `json:"language"`
}
