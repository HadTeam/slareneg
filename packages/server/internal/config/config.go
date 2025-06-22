package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Auth     AuthConfig     `json:"auth"`
	Cache    CacheConfig    `json:"cache"`
	Game     GameConfig     `json:"game"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	StaticDir    string        `json:"staticDir"`
}

type AuthConfig struct {
	JWTSecret      string        `json:"jwtSecret"`
	TokenExpiry    time.Duration `json:"tokenExpiry"`
	BCryptCost     int           `json:"bcryptCost"`
	RateLimitRPS   int           `json:"rateLimitRPS"`
	RateLimitBurst int           `json:"rateLimitBurst"`
}

type CacheConfig struct {
	CleanupInterval time.Duration `json:"cleanupInterval"`
	DefaultTTL      time.Duration `json:"defaultTTL"`
	MaxMemoryMB     int           `json:"maxMemoryMB"`
}

type GameConfig struct {
	MaxRooms            int           `json:"maxRooms"`
	MaxPlayersPerRoom   int           `json:"maxPlayersPerRoom"`
	GameTimeout         time.Duration `json:"gameTimeout"`
	ReconnectTimeout    time.Duration `json:"reconnectTimeout"`
	HeartbeatInterval   time.Duration `json:"heartbeatInterval"`
	MatchmakingInterval time.Duration `json:"matchmakingInterval"`
}

type DatabaseConfig struct {
	Type         string `json:"type"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Database     string `json:"database"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
	SSLMode      string `json:"sslMode"`
}

type LoggingConfig struct {
	Level       string `json:"level"`
	Format      string `json:"format"`
	OutputFile  string `json:"outputFile"`
	MaxSize     int    `json:"maxSize"`
	MaxBackups  int    `json:"maxBackups"`
	MaxAge      int    `json:"maxAge"`
	Compress    bool   `json:"compress"`
	Development bool   `json:"development"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			StaticDir:    "./static",
		},
		Auth: AuthConfig{
			JWTSecret:      "your-secret-key-change-this-in-production",
			TokenExpiry:    24 * time.Hour,
			BCryptCost:     12,
			RateLimitRPS:   10,
			RateLimitBurst: 20,
		},
		Cache: CacheConfig{
			CleanupInterval: 10 * time.Minute,
			DefaultTTL:      1 * time.Hour,
			MaxMemoryMB:     100,
		},
		Game: GameConfig{
			MaxRooms:            100,
			MaxPlayersPerRoom:   8,
			GameTimeout:         30 * time.Minute,
			ReconnectTimeout:    30 * time.Second,
			HeartbeatInterval:   30 * time.Second,
			MatchmakingInterval: 5 * time.Second,
		},
		Database: DatabaseConfig{
			Type:         "sqlite",
			Host:         "localhost",
			Port:         5432,
			Database:     "slareneg.db",
			Username:     "",
			Password:     "",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			SSLMode:      "disable",
		},
		Logging: LoggingConfig{
			Level:       "info",
			Format:      "json",
			OutputFile:  "",
			MaxSize:     100,
			MaxBackups:  3,
			MaxAge:      28,
			Compress:    true,
			Development: false,
		},
	}
}

func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	config.loadFromEnv()

	return config, nil
}

func (c *Config) loadFromEnv() {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Server.Port = p
		}
	}
	if staticDir := os.Getenv("STATIC_DIR"); staticDir != "" {
		c.Server.StaticDir = staticDir
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		c.Auth.JWTSecret = jwtSecret
	}
	if tokenExpiry := os.Getenv("TOKEN_EXPIRY"); tokenExpiry != "" {
		if d, err := time.ParseDuration(tokenExpiry); err == nil {
			c.Auth.TokenExpiry = d
		}
	}

	if cleanupInterval := os.Getenv("CACHE_CLEANUP_INTERVAL"); cleanupInterval != "" {
		if d, err := time.ParseDuration(cleanupInterval); err == nil {
			c.Cache.CleanupInterval = d
		}
	}
	if defaultTTL := os.Getenv("CACHE_DEFAULT_TTL"); defaultTTL != "" {
		if d, err := time.ParseDuration(defaultTTL); err == nil {
			c.Cache.DefaultTTL = d
		}
	}

	if maxRooms := os.Getenv("GAME_MAX_ROOMS"); maxRooms != "" {
		if mr, err := strconv.Atoi(maxRooms); err == nil {
			c.Game.MaxRooms = mr
		}
	}
	if maxPlayersPerRoom := os.Getenv("GAME_MAX_PLAYERS_PER_ROOM"); maxPlayersPerRoom != "" {
		if mp, err := strconv.Atoi(maxPlayersPerRoom); err == nil {
			c.Game.MaxPlayersPerRoom = mp
		}
	}
	if gameTimeout := os.Getenv("GAME_TIMEOUT"); gameTimeout != "" {
		if d, err := time.ParseDuration(gameTimeout); err == nil {
			c.Game.GameTimeout = d
		}
	}
	if reconnectTimeout := os.Getenv("GAME_RECONNECT_TIMEOUT"); reconnectTimeout != "" {
		if d, err := time.ParseDuration(reconnectTimeout); err == nil {
			c.Game.ReconnectTimeout = d
		}
	}

	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		c.Database.Type = dbType
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		c.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if p, err := strconv.Atoi(dbPort); err == nil {
			c.Database.Port = p
		}
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		c.Database.Database = dbName
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		c.Database.Username = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		c.Database.Password = dbPass
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.Logging.Level = logLevel
	}
	if logFormat := os.Getenv("LOG_FORMAT"); logFormat != "" {
		c.Logging.Format = logFormat
	}
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		c.Logging.OutputFile = logFile
	}
	if dev := os.Getenv("LOG_DEVELOPMENT"); dev != "" {
		c.Logging.Development = dev == "true"
	}
}

func (c *Config) SaveConfig(configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	if c.Auth.JWTSecret == "your-secret-key-change-this-in-production" {
		return fmt.Errorf("please change the default JWT secret")
	}

	if c.Game.MaxRooms <= 0 {
		return fmt.Errorf("max rooms must be positive")
	}

	if c.Game.MaxPlayersPerRoom <= 0 {
		return fmt.Errorf("max players per room must be positive")
	}

	return nil
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetDatabaseDSN() string {
	switch c.Database.Type {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host,
			c.Database.Port,
			c.Database.Username,
			c.Database.Password,
			c.Database.Database,
			c.Database.SSLMode,
		)
	case "sqlite":
		return c.Database.Database
	default:
		return ""
	}
}
