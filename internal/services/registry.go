package services

type Category string

const (
	CategoryDatabase  Category = "database"
	CategorySearch    Category = "search"
	CategoryCache     Category = "cache"
	CategoryQueue     Category = "queue"
	CategoryStorage   Category = "storage"
	CategoryBroadcast Category = "broadcasting"
	CategoryMail      Category = "mail"
)

type Service struct {
	Image        string
	DefaultTag   string
	Category     Category
	Port         int
	CanBeShared  bool
}

// All known services. Port is the container-internal (and default host-mapped) port.
var All = []Service{
	// Databases
	{Image: "mysql",       DefaultTag: "latest", Category: CategoryDatabase,  Port: 3306,  CanBeShared: true},
	{Image: "mariadb",     DefaultTag: "latest", Category: CategoryDatabase,  Port: 3306,  CanBeShared: true},
	{Image: "percona",     DefaultTag: "latest", Category: CategoryDatabase,  Port: 3306,  CanBeShared: true},
	{Image: "postgres",    DefaultTag: "latest", Category: CategoryDatabase,  Port: 5432,  CanBeShared: true},
	{Image: "mongo",       DefaultTag: "latest", Category: CategoryDatabase,  Port: 27017, CanBeShared: true},
	{Image: "cassandra",   DefaultTag: "latest", Category: CategoryDatabase,  Port: 9042,  CanBeShared: true},
	// Search
	{Image: "elasticsearch", DefaultTag: "latest", Category: CategorySearch,  Port: 9200,  CanBeShared: true},
	{Image: "opensearch",    DefaultTag: "latest", Category: CategorySearch,  Port: 9200,  CanBeShared: true},
	{Image: "meilisearch",   DefaultTag: "latest", Category: CategorySearch,  Port: 7700,  CanBeShared: true},
	{Image: "typesense",     DefaultTag: "latest", Category: CategorySearch,  Port: 8108,  CanBeShared: true},
	{Image: "solr",          DefaultTag: "latest", Category: CategorySearch,  Port: 8983,  CanBeShared: true},
	// Cache
	{Image: "redis",      DefaultTag: "latest", Category: CategoryCache,     Port: 6379,  CanBeShared: true},
	{Image: "memcached",  DefaultTag: "latest", Category: CategoryCache,     Port: 11211, CanBeShared: false},
	// Queue
	{Image: "rabbitmq",   DefaultTag: "latest", Category: CategoryQueue,     Port: 5672,  CanBeShared: true},
	{Image: "kafka",      DefaultTag: "latest", Category: CategoryQueue,     Port: 9092,  CanBeShared: true},
	// Storage
	{Image: "minio",      DefaultTag: "latest", Category: CategoryStorage,   Port: 9000,  CanBeShared: true},
	// Broadcasting
	{Image: "soketi",     DefaultTag: "latest", Category: CategoryBroadcast, Port: 6001,  CanBeShared: true},
	{Image: "reverb",     DefaultTag: "latest", Category: CategoryBroadcast, Port: 8080,  CanBeShared: true},
	// Mail
	{Image: "mailpit",    DefaultTag: "latest", Category: CategoryMail,      Port: 1025,  CanBeShared: true},
	{Image: "mailhog",    DefaultTag: "latest", Category: CategoryMail,      Port: 1025,  CanBeShared: true},
}

// Find returns the Service for the given image name, or false if unsupported.
func Find(image string) (Service, bool) {
	for _, s := range All {
		if s.Image == image {
			return s, true
		}
	}
	return Service{}, false
}

// IsSupported reports whether the image is a known service.
func IsSupported(image string) bool {
	_, ok := Find(image)
	return ok
}
