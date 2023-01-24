package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"golang.org/x/exp/slices"
)

//go:embed deviceids.txt
var deviceFile embed.FS

func deviceCheck() error {
	file, _ := deviceFile.ReadFile("deviceids.txt")
	deviceIds := strings.Split(string(file), "\n")
	id, err := machineid.ProtectedID("**TaskMaster**")
	if err != nil {
		return err
	}
	if !slices.Contains(deviceIds, id) {
		return ErrUnauth
	}
	return nil
}

func createSiteDir(site string) error {
	if err := os.Mkdir(site, 0755); err != nil {
		return err
	}
	if err := createTaskTemplate(site, false); err != nil {
		return err
	}
	return nil
}

func getFolderList() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	dirList := []string{}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			dirList = append(dirList, fileInfo.Name())
		}
	}
	return dirList, nil
}

func initialCheck() error {
	updated := false

	dirList, err := getFolderList()
	if err != nil {
		return err
	}

	for _, site := range SITES {
		if !slices.Contains(dirList, site) {
			if err = createSiteDir(site); err != nil {
				return err
			}
			updated = true
		}
	}

	if updated {
		return ErrAppUpdated
	}

	return nil
}

func createTaskTemplate(site string, new bool) error {
	files, err := os.ReadDir(fmt.Sprintf("./%s", site))
	if err != nil {
		return err
	}

	fileList := []string{}
	for _, fileInfo := range files {
		fileList = append(fileList, fileInfo.Name())
	}

	if slices.Contains(fileList, "tasks.csv") && new {
		return ErrTasksFileExists
	}

	if !slices.Contains(fileList, "tasks.csv") {
		err = os.WriteFile(fmt.Sprintf("./%s/tasks.csv", site), Configuration[site].tasksHeaders, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadCsv(site string, file string) ([][]string, error) {
	f, err := os.Open(fmt.Sprintf("./%s/%s.csv", site, file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	tasks, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
