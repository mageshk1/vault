package main

import (
	"context"
	"script/database"
	"strconv"

	"github.com/rs/zerolog/log"
)

type RegistryDetails struct {
	WorkspaceID int64  `json:"workspace_id"`
	Name        string `json:"name"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Apikey      string `json:"apikey"`
}

func init() {
	database.GetInstance()
	database.GetVaultCon()

	var registry []RegistryDetails

	if err := database.GetInstance().Model(RegistryDetails{}).Select("workspace_id,name,user_name,password,api_key").Where(
		"status = ? and scan_status = ?", "Active", "Completed").Scan(&registry).Error; err != nil {
		log.Error().Msg("unable to fetch registry details: " + err.Error())
		return
	}

	for _, response := range registry {

		tenantId := strconv.FormatInt(response.WorkspaceID, 10)

		fullPath := tenantId + "/" + response.Name

		secretData := map[string]interface{}{
			"username": response.UserName,
			"password": response.Password,
			"apikey":   response.Apikey,
		}

		basePath := "/vuln-service"

		err := database.Client.KVv1(basePath).Put(context.Background(), fullPath, secretData)

		if err != nil {
			log.Error().Msg("error in stroing data to secret: " + err.Error())
			return
		}

		log.Info().Msg("data stored successfully")

	}

	// str := database.GetInstance().Model(&RegistryDetails{}).Migrator().DropColumn(&RegistryDetails{}, "user_name", "password", "api_key").Error()

	err := database.GetInstance().Exec("ALTER TABLE registry_details DROP COLUMN user_name,DROP COLUMN password, DROP COLUMN api_key").Error

	if err != nil {
		log.Error().Msg("Failed to drop columns:" + err.Error())
		return
	}

	log.Info().Msg("Columns dropped successfully")

}

func main() {

}
