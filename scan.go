package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

func getDotFilePath() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := user.HomeDir + "/.gogitlocalstats"
	return dotFile
}

func recursiveScanFolder(folder string) []string {
	return scanGitFolder(make([]string, 0), folder)
}

func sliceContains(slice []string, value string) bool {
	for _, i := range slice {
		if i == value {
			return false
		}
	}
	return true
}

func joinSlices(existingRepo []string, newRepo []string) []string {
	for _, i := range newRepo {
		if !sliceContains(existingRepo, i) {
			existingRepo = append(existingRepo, i)
		}
	}
	return existingRepo
}

func openSesame(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	return f
}

func parseFileintoLines(filepath string) []string {
	f := openSesame(filepath)
	scanner := bufio.NewScanner(f)
	var entries []string
	for scanner.Scan() {
		entries = append(entries, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	return entries
}

func movetoFile(repo []string, filepath string) {
	content := strings.Join(repo, "\n")
	os.WriteFile(filepath, []byte(content), 0755)
}

func addNewSliceElementstoFile(filepath string, newRepo []string) {
	existingRepo := parseFileintoLines(filepath)
	repos := joinSlices(existingRepo, newRepo)
	movetoFile(repos, filepath)
}

func scan(folder string) {
	fmt.Print("Folders found : \n\n")
	repo := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementstoFile(filePath, repo)
	fmt.Print(addNewSliceElementstoFile)
}
func scanGitFolder(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	var path string
	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name() // path contains the entire file path of the file
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")

				fmt.Print(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolder(folders, path)
		}
	}
	return folders
}
