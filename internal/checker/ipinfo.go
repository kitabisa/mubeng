package checker

type IPInfo struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Readme   string `json:"readme"`
	Region   string `json:"region"`
	Timezone string `json:"timezone"`
}
