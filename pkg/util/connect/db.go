package connect

const (
	// Driver names
	SQLite   = "sqlite3"
	Postgres = "postgres"

	// Table names
	TableAccessSession        = "access_session"
	TableAccessPassword       = "user_password"
	TableNamespace            = "namespaces"
	TableNamespaceMemberships = "namespace_memberships"
	TableNamespaceDomains     = "namespace_domains"
	TableNamespaceConfig      = "namespace_config"
	TableManagerEdge          = "namespace_manager_edge"
)
