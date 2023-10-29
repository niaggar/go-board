package config

import (
	"encoding/json"
	"fmt"
	"go-board/models"
	"os"
)

func GetConfiguration(route string) models.ConfigurationFile {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config models.ConfigurationFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func GetGlobalConfiguration() models.ConfigurationGlobal {
	file, err := os.Open("./.config.json")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config models.ConfigurationGlobal
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}
