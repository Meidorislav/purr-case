package cases_dto

type DropEntry struct {
	SKU    string `json:"sku"`
	Weight int    `json:"weight"`
}

type CaseCustomAttributes struct {
	Type      string      `json:"type"`
	DropTable []DropEntry `json:"drop_table"`
}

type OpenCaseResponse struct {
	CaseSKU string `json:"case_sku"`
	WonSKU  string `json:"won_sku"`
}
