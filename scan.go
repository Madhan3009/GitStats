package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
	"syscall"
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
			return true
		}
	}
	return false
}

func joinSlices(existingRepo []string, newRepo []string) []string {
	for _, i := range newRepo {
		if !sliceContains(existingRepo, i) {
			existingRepo = append(existingRepo, i)
		}
	}
	return existingRepo
}

func parseFileintoLines(filepath string) []string {
	f, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
		panic(err)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

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

func addNewSliceElementstoFile(filepath string, newRepo []string) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		panic(err)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	scanner := bufio.NewScanner(f)
	var existingRepo []string
	for scanner.Scan() {
		existingRepo = append(existingRepo, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	repos := joinSlices(existingRepo, newRepo)

	if err := f.Truncate(0); err != nil {
		panic(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		panic(err)
	}

	content := strings.Join(repos, "\n")
	if _, err := f.WriteString(content); err != nil {
		panic(err)
	}
}

func scan(folder string) {
	fmt.Print("Folders found : \n\n")
	repo := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementstoFile(filePath, repo)
	fmt.Printf("\nDone! Added %d new repos to %s.\n", len(repo), filePath)
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
