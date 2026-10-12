package domain

type Role string

const (
	RoleViewer Role = "viewer"
	RoleUser   Role = "user"
	RoleAdmin  Role = "admin"
)

func (r Role) IsValid() bool {
	return r == RoleViewer || r == RoleUser || r == RoleAdmin
}
