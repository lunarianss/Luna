package dto

type StatisticQuery struct {
	Start string `json:"start" form:"start"`
	End   string `json:"end" form:"end"`
}
