package services

import "sort"

// Healthcheck defines container health check options.
type Healthcheck struct {
	Test     []string
	Interval string
	Timeout  string
	Retries  int
}

// Service defines a predefined service entry in the catalog.
type Service struct {
	Name        string
	Image       string
	Ports       []string
	Env         map[string]string
	Volumes     []string
	Healthcheck *Healthcheck
}

var catalog = map[string]Service{
	"postgres": {
		Name:    "postgres",
		Image:   "postgres:16-alpine",
		Ports:   []string{"5432:5432"},
		Env:     map[string]string{"POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "app"},
		Volumes: []string{"postgres_data:/var/lib/postgresql/data"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
	},
	"mysql": {
		Name:    "mysql",
		Image:   "mysql:8.4",
		Ports:   []string{"3306:3306"},
		Env:     map[string]string{"MYSQL_ROOT_PASSWORD": "root", "MYSQL_DATABASE": "app"},
		Volumes: []string{"mysql_data:/var/lib/mysql"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD", "mysqladmin", "ping", "-h", "localhost"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
	},
	"redis": {
		Name:    "redis",
		Image:   "redis:7.4-alpine",
		Ports:   []string{"6379:6379"},
		Volumes: []string{"redis_data:/data"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: "10s",
			Timeout:  "3s",
			Retries:  5,
		},
	},
	"mongodb": {
		Name:    "mongodb",
		Image:   "mongo:7.0",
		Ports:   []string{"27017:27017"},
		Volumes: []string{"mongodb_data:/data/db"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD", "mongosh", "--eval", "db.adminCommand('ping')"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
	},
	"rabbitmq": {
		Name:    "rabbitmq",
		Image:   "rabbitmq:3.13-management",
		Ports:   []string{"5672:5672", "15672:15672"},
		Volumes: []string{"rabbitmq_data:/var/lib/rabbitmq"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD", "rabbitmq-diagnostics", "ping"},
			Interval: "15s",
			Timeout:  "10s",
			Retries:  5,
		},
	},
	"minio": {
		Name:    "minio",
		Image:   "minio/minio:RELEASE.2024-06-13T22-53-53Z",
		Ports:   []string{"9000:9000", "9001:9001"},
		Env:     map[string]string{"MINIO_ROOT_USER": "minioadmin", "MINIO_ROOT_PASSWORD": "minioadmin"},
		Volumes: []string{"minio_data:/data"},
		Healthcheck: &Healthcheck{
			Test:     []string{"CMD", "mc", "ready", "local"},
			Interval: "15s",
			Timeout:  "10s",
			Retries:  5,
		},
	},
	"mailpit": {
		Name:  "mailpit",
		Image: "axllent/mailpit:v1.20",
		Ports: []string{"1025:1025", "8025:8025"},
	},
	"meilisearch": {
		Name:    "meilisearch",
		Image:   "getmeili/meilisearch:v1.9",
		Ports:   []string{"7700:7700"},
		Env:     map[string]string{"MEILI_MASTER_KEY": "masterKey"},
		Volumes: []string{"meili_data:/meili_data"},
	},
}

// Get returns the predefined service definition for the given service name, if present.
func Get(name string) (Service, bool) {
	s, ok := catalog[name]
	return s, ok
}

// Names returns a sorted list of all available service catalog names.
func Names() []string {
	out := make([]string, 0, len(catalog))
	for k := range catalog {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
