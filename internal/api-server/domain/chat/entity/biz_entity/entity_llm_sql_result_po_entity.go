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
