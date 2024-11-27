package po_entity

type TenantJoinResult struct {
	ID           string                 `json:"tenant_id" gorm:"column:tenant_id"`
	Name         string                 `json:"tenant_name" gorm:"column:tenant_name"`
	Plan         string                 `json:"tenant_plan" gorm:"column:tenant_plan"`
	Status       string                 `json:"tenant_status" gorm:"column:tenant_status"`
	CreatedAt    int64                  `json:"tenant_created_at" gorm:"column:tenant_created_at"`
	UpdatedAt    int64                  `json:"tenant_updated_at" gorm:"column:tenant_updated_at"`
	CustomConfig map[string]interface{} `json:"tenant_custom_config" gorm:"column:tenant_custom_config;serializer:json"`
	Role         string                 `json:"tenant_join_role" gorm:"column:tenant_join_role"`
}
