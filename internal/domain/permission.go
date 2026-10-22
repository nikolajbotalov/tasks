package domain

type Permission string

const (
	PermissionTaskRead      Permission = "task:read"
	PermissionTaskCreate    Permission = "task:create"
	PermissionTaskUpdate    Permission = "task:update"
	PermissionTaskDelete    Permission = "task:delete"
	PermissionCommentRead   Permission = "comment:read"
	PermissionCommentCreate Permission = "comment:create"
	PermissionUserRead      Permission = "user:read"
	PermissionUserUpdate    Permission = "user:update"
)

type RolePermissions map[Role]map[Permission]struct{}

func (rp RolePermissions) Has(role Role, perm Permission) bool {
	perms, ok := rp[role]
	if !ok {
		return false
	}

	_, ok = perms[perm]
	return ok
}
