package domain

const (
	AdminRole    int32 = 1
	SecurityRole int32 = 2

	StudentRole int32 = 4
	GuestRole   int32 = 5
)

type Role struct {
	ID   int32
	Name string
}
