package database

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once
var dba *gorm.DB
var Client *api.Client

// GetInstance - Returns a DB instance
func GetInstance() *gorm.DB {
	once.Do(func() {

		user := ""
		password := ""

		host := ""
		port := ""
		dbname := ""
		schema := ""

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s", host, user, password, dbname, port, schema)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		dba = db
		if err != nil {
			log.Panic().Msgf("Error connecting to the database at %s:%s/%s", host, port, dbname)
		}

		log.Info().Msgf("Successfully established connection to %s:%s/%s", host, port, dbname)
	})
	return dba
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetVaultCon() *api.Client {
	token := "" // add vault token
	vaultAddress := "http://localhost:8200/"

	client, err := api.NewClient(&api.Config{Address: vaultAddress, HttpClient: httpClient})

	if err != nil {
		log.Error().Msg("error connecting to vault: " + err.Error())
		return client
	}

	client.SetToken(token)

	Client = client

	return Client

}
