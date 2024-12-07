package postgresql

// The database schema to be used.
const scheme = "public."

// The names of tables in the database.
const (
	TableUsers           = scheme + "users"
	TableRefreshSessions = scheme + "refresh_sessions"
	TableRole            = scheme + "roles"
)
