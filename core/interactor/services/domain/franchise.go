package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

type Franchise struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Branches    []Branch `json:"branches,omitempty"`
}

func (f *Franchise) SetID() {
	f.ID = uuid.Generate()
}

func (f *Franchise) ToLogger() []string {
	return []string{
		"id:" + f.ID,
		"name:" + f.Name,
	}
}
