package app

import (
	"fmt"
	"go-board/logic"
	"go-board/logic/config"
	"go-board/utils"
	"go-board/utils/export"
	"log"
	"os"
	"sync"
	"time"
)

func RunConsole() {
	configsRoute, exportRoute := utils.GetBaseAppRoute()
	configsFiles := utils.GetNewExpConfigs(configsRoute)
	configsSelected := utils.SelectConfigsFiles(&configsFiles)

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
			exportHistoImg := baseRouteWg + fmt.Sprintf("/histo-%d.png", pos)

			gbConfig.Experiment.ExportHistogram.Route = exportHistoRoute
			gbConfig.Experiment.ExportPaths.Route = exportPathRoute

			gb := logic.NewGaltonBoard(&gbConfig)
			gb.RunAll()
			gb.Finish()

			export.PlotHistogram(exportHistoImg, exportHistoRoute)
		}(i, baseRoute)
	}

	wgInternal.Wait()

	finalTime := time.Now()
	elapsed := finalTime.Sub(currentTime)
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Exp --%s-- saved at: %s\n", gbConfig.Experiment.Name, baseRoute)
}
