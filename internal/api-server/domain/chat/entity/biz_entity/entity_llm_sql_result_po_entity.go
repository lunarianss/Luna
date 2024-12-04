package biz_entity

type StatisticDailyConversations struct {
	Data []*StatisticDailyConversationsItem `json:"data"`
}

type StatisticDailyConversationsItem struct {
	Date  string `json:"date"`
	Count int64  `json:"conversation_count" gorm:"column:message_count"`
}
