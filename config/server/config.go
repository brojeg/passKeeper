package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer
	ExternalDependency
	ServerAuth
	ServerLog
	Certificates
}
type HTTPServer struct {
	ServerPort string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
}
type ExternalDependency struct {
	Database string `env:"DATABASE_URI"`
}
type ServerAuth struct {
	JWTPassword    string `env:"JWT_PASSWORD"`
	ExpirationTime int    `env:"EXPIRATION_TIME" envDefault:"15"`
}
type ServerLog struct {
	Log string `env:"SERVER_LOG"`
}
type Certificates struct {
	TLSCertFile string `env:"TLSCERTFILE"`
	TLSKeyFile  string `env:"TLSKEYFILE"`
}

func NewServerConfig() *Config {
	sc := Config{}
	godotenv.Load(".env")
	err := env.Parse(&sc.ExternalDependency)
	env.Parse(&sc.HTTPServer)
	env.Parse(&sc.ServerAuth)

	_, envAdddressExists := os.LookupEnv("RUN_ADDRESS")
	_, envDBExists := os.LookupEnv("DATABASE_URI")
	_, envJWTPAsswordExists := os.LookupEnv("JWT_PASSWORD")
	_, envExpirationTimeExists := os.LookupEnv("EXPIRATION_TIME")
	_, envTLSCertFileExists := os.LookupEnv("TLSCERTFILE")
	_, envTLSKeyFileExists := os.LookupEnv("TLSKEYFILE")

	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
	flag.Func("a", "Server address (default localhost:8080)", func(flagValue string) error {
		if envAdddressExists {
			return nil
		}
		sc.ServerPort = flagValue
		return nil
	})
	flag.Func("d", "Postgres connection string (No default value)", func(flagValue string) error {
		if envDBExists {
			return nil
		}
		sc.Database = flagValue

		return nil
	})
	flag.Func("p", "Check for JWT", func(flagValue string) error {
		if envJWTPAsswordExists {
			return nil
		}
		sc.JWTPassword = flagValue
		return nil
	})
	flag.Func("t", "TTL for JWT token (default 15m", func(flagValue string) error {
		if envExpirationTimeExists {
			return nil
		}
		intVar, err := strconv.Atoi(flagValue)
		if err != nil {
			return err
		}
		sc.ExpirationTime = intVar
		return nil
	})
	flag.Func("tc", "TLS Cert file", func(flagValue string) error {
		if envTLSCertFileExists {
			return nil
		}
		sc.TLSCertFile = flagValue
		return nil
	})
	flag.Func("tk", "TLS Key file", func(flagValue string) error {
		if envTLSKeyFileExists {
			return nil
		}
		sc.TLSKeyFile = flagValue
		return nil
	})
	flag.Parse()

	return &sc
}
