package common

type Data struct {
	Serial int `json:"serial"`
	Symbol string `json:"symbol"`
	Company string `json:"company"`
	Sector string `json:"sector"`
	Industry string `json:"industry"`
	Country string `json:"country"`
	MCap float64 `json:"mcap"`
	PE float64 `json:"pe"`
	Price float64 `json:"price"`
	Change float64 `json:"change"`
	Volume float64 `json:"volume"`
	LoadID int `json:"loadId"`
}

type ResultData struct {
	Symbol string `json:"symbol"`
	Company string `json:"company"`
	Industry string `json:"industry"`
	ChangeSMA float64 `json:"changeSma"`
	ChangeVMA float64 `json:"changeVma"`
	PriceSMA float64 `json:"priceSma"`
	PriceVMA float64 `json:"priceVma"`
	Count int `json:"count"`
	URL string `json:"url"`
	LoadIDs string `json:"loadIds"`
}

type ResultDataList []*ResultData

func (rd ResultDataList) Len() int {
	return len(rd)
}

// Swap is part of sort.Interface.
func (rd ResultDataList) Swap(i, j int) {
	rd[i], rd[j] = rd[j], rd[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (rd ResultDataList) Less(i, j int) bool {
	return rd[i].ChangeSMA > rd[j].ChangeSMA
	//return rd[i].Count > rd[j].Count
}
