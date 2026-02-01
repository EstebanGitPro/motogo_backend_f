package domain

type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DepartmentID string `json:"department_id"`
}
