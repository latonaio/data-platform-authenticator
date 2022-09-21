package configs

import (
	"os"
	"strconv"
)

type Configs struct {
	Server     Server
	Database   Database
	Jwt        Jwt
	PrivateKey string
}

type Server struct {
	Port             string
	ShutdownWaitTIme int
}

type Database struct {
	HostName     string
	Port         string
	UserName     string
	UserPassword string
	MaxOpenCon   int
	MaxIdleCon   int
	MaxLifeTime  int
	Name         string
	TableName    string
}

type Jwt struct {
	Exp int64
}

func New() (*Configs, error) {
	cfg := &Configs{
		Server: Server{
			Port:             os.Getenv("SERVER_PORT"),
			ShutdownWaitTIme: parseInt(os.Getenv("SHUT_DOWN_WAIT_TIME")),
		},
		Database: Database{
			HostName:     os.Getenv("HOST_NAME"),
			Port:         os.Getenv("DB_PORT"),
			UserName:     os.Getenv("USER_NAME"),
			UserPassword: os.Getenv("USER_PASSWORD"),
			MaxOpenCon:   parseInt(os.Getenv("MAX_OPEN_CON")),
			MaxIdleCon:   parseInt(os.Getenv("MAX_IDLE_CON")),
			MaxLifeTime:  parseInt(os.Getenv("MAX_LIFE_TIME")),
			Name:         os.Getenv("NAME"),
			TableName:    os.Getenv("TABLE_NAME"),
		},
		Jwt: Jwt{
			Exp: parseInt64(os.Getenv("EXP")),
		},
		PrivateKey: os.Getenv("PRIVATE_KEY"),
	}

	return cfg, nil
}

func parseInt(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}

func parseInt64(value string) int64 {
	v, _ := strconv.ParseUint(value, 10, 64)
	return int64(v)
}
