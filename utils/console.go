package utils

import (
	"fmt"
	"go-board/logic/config"
	"os"
	"strconv"
	"strings"
)

func GetNewExpConfigs(baseRoute string) map[int]string {
	directory, err := os.ReadDir(baseRoute)
	if err != nil {
		fmt.Print(err.Error())
		return make(map[int]string)
	}

	files := make(map[int]string)

	for idx, entry := range directory {
		if !entry.IsDir() {
			name := entry.Name()
			files[idx] = name
		}
	}

	return files
}

func GetBaseAppRoute() (string, string) {
	globalConf := config.GetGlobalConfiguration()

	return globalConf.Base + globalConf.Configs, globalConf.Base + globalConf.Exports
}

func SelectConfigsFiles(files *map[int]string) []int {
	selected := make([]int, 0)

	fmt.Printf("Select the config files to execute:\n")
	fmt.Printf("Idx \t Config name\n")
	for idx, name := range *files {
		fmt.Printf("[%d] \t %s\n", idx, name)
	}

	fmt.Println("Write your selection like: 1,2,5,3.")
	var selection string
	_, err := fmt.Scan(&selection)
	if err != nil {
		panic(err)
	}

	selected = ParseSelection(selection)
	return selected
}

func ParseSelection(selection string) []int {
	items := strings.Split(selection, ",")
	index := make([]int, 0)

	for _, item := range items {
		conv, err := strconv.Atoi(item)
		if err == nil {
			index = append(index, conv)
		}
	}

	return index
}
