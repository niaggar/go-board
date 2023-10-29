package main

import (
	"fmt"
	"go-board/logic"
	"go-board/logic/config"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	configsRoute, exportRoute := getBaseAppRoute()
	configsFiles := getNewExpConfigs(configsRoute)
	configsSelected := selectConfigsFiles(&configsFiles)

	var wg sync.WaitGroup

	startTime := time.Now()
	fmt.Printf("\nRunning %d experiments at %s\n", len(configsSelected), startTime.Format("2006-01-02-15-04-05"))
	for val := range configsSelected {
		gbConfigRoute := configsRoute + "/" + configsFiles[val]

		wg.Add(1)
		go executeGaltonBoard(gbConfigRoute, exportRoute, &wg)
	}

	wg.Wait()

	finalTime := time.Now()
	elapsed := finalTime.Sub(startTime)
	fmt.Printf("Total time: %v\n", elapsed)
}

func executeGaltonBoard(gbConfigRoute, exportRoute string, wg *sync.WaitGroup) {
	defer wg.Done()

	currentTime := time.Now()
	timeTxt := currentTime.Format("2006-01-02-15-04-05")
	gbConfig := config.GetConfiguration(gbConfigRoute)
	baseRoute := fmt.Sprintf("%s/exp-%s-%s", exportRoute, gbConfig.Experiment.Name, timeTxt)

	fmt.Printf("Execute of --%s-- at %s\n", gbConfig.Experiment.Name, timeTxt)
	if _, err := os.Stat(baseRoute); os.IsNotExist(err) {
		err := os.MkdirAll(baseRoute, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	var wgInternal sync.WaitGroup
	for i := 0; i < gbConfig.Experiment.Executions; i++ {
		wgInternal.Add(1)

		go func(pos int, baseRouteWg string) {
			defer wgInternal.Done()

			exportHistoRoute := baseRouteWg + fmt.Sprintf("/histo-%d.dat", pos)
			exportPathRoute := baseRouteWg + fmt.Sprintf("/path-%d.dat", pos)

			gbConfig.Experiment.ExportHistogram.Route = exportHistoRoute
			gbConfig.Experiment.ExportPaths.Route = exportPathRoute

			gb := logic.NewGaltonBoard(&gbConfig)
			gb.Run()
			gb.Finish()
		}(i, baseRoute)
	}

	wgInternal.Wait()

	finalTime := time.Now()
	elapsed := finalTime.Sub(currentTime)
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Exp --%s-- saved at: %s\n", gbConfig.Experiment.Name, baseRoute)
}

func parseSelection(selection string) []int {
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

func getNewExpConfigs(baseRoute string) map[int]string {
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

func selectConfigsFiles(files *map[int]string) []int {
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

	selected = parseSelection(selection)
	return selected
}

func getBaseAppRoute() (string, string) {
	globalConf := config.GetGlobalConfiguration()

	return globalConf.Base + globalConf.Configs, globalConf.Base + globalConf.Exports
}
