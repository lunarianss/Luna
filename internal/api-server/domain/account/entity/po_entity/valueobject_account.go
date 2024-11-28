package po_entity

type AccountStatus string

const (
	PENDING       AccountStatus = "pending"
	UNINITIALIZED AccountStatus = "uninitialized"
	ACTIVE        AccountStatus = "active"
	BANNED        AccountStatus = "banned"
	CLOSED        AccountStatus = "closed"
)

type TenantStatus string

const (
	TENANT_NORMAL TenantStatus = "normal"
	ARCHIVE       TenantStatus = "archive"
)

type TenantAccountJoinRole string

const (
	OWNER            TenantAccountJoinRole = "owner"
	ADMIN            TenantAccountJoinRole = "admin"
	NORMAL           TenantAccountJoinRole = "normal"
	DATASET_OPERATOR TenantAccountJoinRole = "dataset_operator"
)
