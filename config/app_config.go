package config

type Config struct {
	AppName    string `envconfig:"APP_NAME" default:"HertzApp"`
	Port       int    `envconfig:"PORT" default:"8080"`
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     int    `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	JWTSecret  string `envconfig:"JWT_SECRET" required:"true"`
	ElasticURL string `envconfig:"ELASTIC_URL" required:"true" default:"http://localhost:9200"`
}
