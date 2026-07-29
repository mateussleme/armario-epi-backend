package viaonda

type TagEntry struct {
	Id  string `json:"tag_id"`
	Epc string `json:"tag_epc"`
}

type GetEntriesResponse struct {
	RecordCount int        `json:"record_count"`
	Records     []TagEntry `json:"dados"`
}

type GpioResponse struct {
	Status string `json:"status gpio"`
}
