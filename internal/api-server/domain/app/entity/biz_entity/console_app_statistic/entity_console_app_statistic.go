package biz_entity

type StatisticDailyConversations struct {
	Data []*StatisticDailyConversationsItem `json:"data"`
}

type StatisticDailyUser struct {
	Data []*StatisticDailyUsersItem `json:"data"`
}

type StatisticAverageInteraction struct {
	Data []*StatisticAverageInteractionItem `json:"data"`
}

type StatisticTokenCosts struct {
	Data []*StatisticTokenCostsItem `json:"data"`
}

type StatisticTokenCostsItem struct {
	Currency   string  `json:"currency"`
	Date       string  `json:"date"`
	TotalPrice float64 `json:"total_price" gorm:"column:total_price"`
	TotalCount int64   `json:"token_count" gorm:"column:token_count"`
}

type StatisticDailyConversationsItem struct {
	Date  string `json:"date"`
	Count int64  `json:"conversation_count" gorm:"column:message_count"`
}

type StatisticDailyUsersItem struct {
	Date  string `json:"date"`
	Count int64  `json:"terminal_count" gorm:"column:message_count"`
}

type StatisticAverageInteractionItem struct {
	Date  string  `json:"date"`
	Count float64 `json:"interactions" gorm:"column:interactions"`
}
