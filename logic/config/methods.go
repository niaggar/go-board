package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetNewConfiguration(route string) BaseConfig {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config BaseConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func GetCurrentConfiguration(route string) StateConfig {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config StateConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func GetGlobalConfiguration() GlobalConfig {
	file, err := os.Open("./.config.json")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config GlobalConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func SaveCurrentConfiguration(route string, config StateConfig) {
	file, err := os.Create(route)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		panic(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error al codificar el archivo JSON:", err)
		panic(err)
	}
}
