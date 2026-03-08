package depDtoReq

type CreateDepartment struct {
	Name     string `json:"name" validate:"required,min=1,max=200"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type UpdateDepartment struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	ParentID *uint   `json:"parent_id,omitempty"`
}
