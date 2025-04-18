package enum

type Status string

const (
	STATUS_ACTIVE   Status = "active"
	STATUS_INACTIVE Status = "inactive"
)

func GetStatus(t Status) Status {
	switch t {
	case STATUS_ACTIVE:
		return STATUS_ACTIVE
	default:
		return STATUS_INACTIVE
	}
}

