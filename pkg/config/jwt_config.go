package config

type JWTConfig struct {
	SecretKey string
	TokenLifetime int
}