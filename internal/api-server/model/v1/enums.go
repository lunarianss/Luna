package model

type CreatedByRole string

const (
	AccountCreatedByRole CreatedByRole = "account"
	EndUserCreatedByRole CreatedByRole = "end_user"
)

type UserFrom string

const (
	AccountUserFrom UserFrom = "account"
	EndUserUserFrom UserFrom = "end-user"
)

type WorkflowRunTriggeredFrom string

const (
	DebuggingWorkflowRunTriggeredFrom WorkflowRunTriggeredFrom = "debugging"
	AppRunWorkflowRunTriggeredFrom    WorkflowRunTriggeredFrom = "app-run"
)
