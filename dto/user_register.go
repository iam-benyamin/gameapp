package dto // Package dto means data transfer object

type RegisterRequest struct {
	Name        string `json:"name"`         // struct tag are like meta information and compiler will ignore them
	PhoneNumber string `json:"phone_number"` // but some package like json marshal will look at them
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}
