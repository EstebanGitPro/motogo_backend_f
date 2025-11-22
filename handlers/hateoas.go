package handlers

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
}

func BuildAccountLinks(baseURL string, accountID string) []Link {
	return []Link{
		{
			Href:   baseURL + "/motogo/api/v1/accounts/" + accountID,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   baseURL + "/motogo/api/v1/accounts/" + accountID,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   baseURL + "/motogo/api/v1/accounts/" + accountID,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   baseURL + "/motogo/api/v1/accounts",
			Rel:    "collection",
			Method: "GET",
		},
	}
}
