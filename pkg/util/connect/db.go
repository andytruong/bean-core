package connect

const (
	// Driver names
	SQLite   = "sqlite3"
	Postgres = "postgres"

	// Table names
	TableAccessSession        = "access_session"
	TableNamespace            = "namespaces"
	TableNamespaceMemberships = "namespace_memberships"
	TableNamespaceDomains     = "namespace_domains"
	TableNamespaceConfig      = "namespace_config"
)
