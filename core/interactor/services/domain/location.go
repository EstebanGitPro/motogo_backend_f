package domain

// Department represents a Colombian department
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// City represents a city within a department
type City struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DepartmentID string `json:"department_id"`
}

// Message codes for location module
const (
	MsgDepartmentsRetrieved = "MOD_L_DEP_EXI_00001"
	MsgCitiesRetrieved      = "MOD_L_CIT_EXI_00001"
)
