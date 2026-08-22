package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddress string

	UploadDirectory string
	OutputDirectory string

	WorkerCount int

	MaxUploadSize int64

	MaxImageWidth int

	FileRetention   time.Duration
	CleanupInterval time.Duration
	CommandTimeout  time.Duration

	JobStore string

	DatabaseURL            string
	DatabaseMaxConns       int
	DatabaseConnectTimeout time.Duration

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryFolder    string
}

func Load() Config {
	return Config{
		// Render (and most cloud platforms) inject PORT, not SERVER_ADDRESS.
		// Check PORT first so the server binds on the correct port in production.
		ServerAddress: buildServerAddress(),

		UploadDirectory: getEnv(
			"UPLOAD_DIRECTORY",
			"uploads/originals",
		),

		OutputDirectory: getEnv(
			"OUTPUT_DIRECTORY",
			"uploads/processed",
		),

		WorkerCount: getEnvInt(
			"WORKER_COUNT",
			3,
		),

		MaxUploadSize: int64(
			getEnvInt(
				"MAX_UPLOAD_SIZE_MB",
				32,
			),
		) * 1024 * 1024,

		MaxImageWidth: getEnvInt(
			"MAX_IMAGE_WIDTH",
			5000,
		),

		FileRetention: time.Duration(
			getEnvInt(
				"FILE_RETENTION_HOURS",
				24,
			),
		) * time.Hour,

		CleanupInterval: time.Duration(
			getEnvInt(
				"CLEANUP_INTERVAL_MINUTES",
				60,
			),
		) * time.Minute,

		CommandTimeout: time.Duration(
			getEnvInt(
				"COMMAND_TIMEOUT_SECONDS",
				60,
			),
		) * time.Second,

		JobStore: getEnv(
			"JOB_STORE",
			"postgres",
		),

		DatabaseURL: buildDatabaseURL(),

		DatabaseMaxConns: getEnvInt(
			"DATABASE_MAX_CONNS",
			10,
		),

		DatabaseConnectTimeout: time.Duration(
			getEnvInt(
				"DATABASE_CONNECT_TIMEOUT_SECONDS",
				10,
			),
		) * time.Second,

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryFolder:    getEnv("CLOUDINARY_FOLDER", "featherweight"),
	}
}

// buildDatabaseURL returns DATABASE_URL directly if it is set.
// Otherwise, it builds a PostgreSQL connection URL from the DB_* variables.
func buildDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "featherweight")
	password := getEnv("DB_PASSWORD", "featherweight")
	name := getEnv("DB_NAME", "featherweight")
	sslMode := getEnv("DB_SSLMODE", "disable")

	return "postgres://" +
		user + ":" +
		password + "@" +
		host + ":" +
		port + "/" +
		name +
		"?sslmode=" + sslMode
}

// buildServerAddress resolves the address the HTTP server should listen on.
// Priority: PORT (injected by Render) → SERVER_ADDRESS → :8080 (local dev).
func buildServerAddress() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return getEnv("SERVER_ADDRESS", ":8080")
}

func getEnv(
	key string,
	defaultValue string,
) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(
	key string,
	defaultValue int,
) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsedValue
}
