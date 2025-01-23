package biz_entity

type IndexType string

const (
	PARAGRAPH_INDEX    IndexType = "text_model"
	QA_INDEX           IndexType = "qa_model"
	PARENT_CHILD_INDEX IndexType = "parent_child_index"
	SUMMARY_INDEX      IndexType = "summary_index"
)
