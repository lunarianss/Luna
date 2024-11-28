package po_entity

type ProviderType string

type ProviderQuotaType string

const (
	CUSTOM ProviderType = "custom"
	SYSTEM ProviderType = "system"
)

const (
	PAID ProviderQuotaType = "paid"

	FREE ProviderQuotaType = "free"

	TRIAL ProviderQuotaType = "trial"
)
