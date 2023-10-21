package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetNewConfiguration(route string) NewConfig {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config NewConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func GetCurrentConfiguration(route string) CurrentConfig {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		panic(err)
	}

	defer file.Close()

	var config CurrentConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error al decodificar el archivo JSON:", err)
		panic(err)
	}

	return config
}

func SaveCurrentConfiguration(route string, config CurrentConfig) {
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
