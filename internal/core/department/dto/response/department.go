package depDtoResp

import (
	empModel "org-structure-api/internal/core/employee/model"
	"time"
)

type Department struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`

	Employees []*empModel.Employee `json:"employees,omitempty"`
	Children  []*Department        `json:"children,omitempty"`
}
