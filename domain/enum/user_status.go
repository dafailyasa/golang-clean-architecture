package enum

type Status string

const (
	UserStatusActive    Status = "active"
	UserStatusInactive  Status = "inactive"
	UserStatusSuspended Status = "suspended"
)

func (s Status) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	}
	return false
}
