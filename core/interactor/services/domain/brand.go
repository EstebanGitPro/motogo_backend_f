package domain

// Brand represents a motorcycle brand from the catalog
type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Message codes for brands module
const (
	MsgBrandsRetrieved = "MOD_B_BRD_EXI_00001"
)
