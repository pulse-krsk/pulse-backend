package postgresql

// The database schema to be used.
const scheme = "public."

// The names of tables in the database.
const (
	TableUsers              = scheme + "users"
	TableRefreshSessions    = scheme + "refresh_sessions"
	TableRoles              = scheme + "roles"
	TableEventTypes         = scheme + "event_types"
	TableUsersFavoriteTypes = scheme + "users_favorite_types"
)
