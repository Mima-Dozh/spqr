package spqrparser

var reservedWords = map[string]int{
	"pools":               POOLS,
	"servers":             SERVERS,
	"clients":             CLIENTS,
	"client":              CLIENT,
	"databases":           DATABASES,
	"show":                SHOW,
	"stats":               STATS,
	"kill":                KILL,
	"column":              COLUMN,
	"columns":             COLUMNS,
	"shard":               SHARD,
	"rule":                RULE,
	"sharding":            SHARDING,
	"create":              CREATE,
	"add":                 ADD,
	"key":                 KEY,
	"range":               RANGE,
	"shards":              SHARDS,
	"key_ranges":          KEY_RANGES,
	"sharding_rules":      SHARDING_RULES,
	"lock":                LOCK,
	"unlock":              UNLOCK,
	"drop":                DROP,
	"all":                 ALL,
	"shutdown":            SHUTDOWN,
	"split":               SPLIT,
	"version":             VERSION,
	"from":                FROM,
	"by":                  BY,
	"to":                  TO,
	"with":                WITH,
	"unite":               UNITE,
	"listen":              LISTEN,
	"register":            REGISTER,
	"unregister":          UNREGISTER,
	"router":              ROUTER,
	"move":                MOVE,
	"routers":             ROUTERS,
	"address":             ADDRESS,
	"host":                HOST,
	"route":               ROUTE,
	"dataspace":           DATASPACE,
	"table":               TABLE,
	"hash":                HASH,
	"function":            FUNCTION,
	"backend_connections": BACKEND_CONNECTIONS,
	"where":               WHERE,
}
