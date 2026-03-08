package empDtoReq

type CreateEmployee struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	HiredAt  string `json:"hired_at,omitempty"`
}
