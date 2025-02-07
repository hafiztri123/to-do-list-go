package logger

import (
	
	"os"
    "github.com/rs/zerolog/log" 
	"github.com/rs/zerolog"
)

func Init(){
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger= zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}