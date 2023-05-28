package fileHandler

import (
  "fmt"
	"os"
	"path/filepath"
	"bufio"
)

const (
	todosDir      = ".todos"
	tasksFileName = "tasks.txt"
)

var initialTasks = []string {
	"[0] Welcome, press [?] to view the keybinds",
	"[1] You can use this program to keep track of what you want done",
	"[1] Don't forget to write any changes you make with [w]",
}

func CreateDirIfNotExists() {
	homeDir, _ := os.UserHomeDir()

	todosDirPath := filepath.Join(homeDir, todosDir) // Create ~/.todos/ (first time setup)
	_, err := os.Stat(todosDirPath)
	if os.IsNotExist(err) {
		_ = os.Mkdir(todosDirPath, 0700)
	}

	tasksFilePath := filepath.Join(todosDirPath, tasksFileName) 
	_, err = os.Stat(tasksFilePath)
	if os.IsNotExist(err) {
		file, _ := os.Create(tasksFilePath) // Create tasks.txt
		defer file.Close()

		for _, task := range initialTasks { 
			_, _ = file.WriteString(task + "\n") // Populate with dummy data
		}
	}
}

func SaveTasks(tasks []string) {
  homeDir, _ := os.UserHomeDir() 
  filePath := filepath.Join(homeDir, todosDir, tasksFileName)
  file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

  for _, task := range tasks {
    _, _ = file.WriteString(task + "\n")
  }
}

func LoadTasks() ([]string, error) {
  homeDir, err := os.UserHomeDir() 
  if err != nil {
    fmt.Println("Error getting home dir:", err)
    return nil, err
  }

  filePath := filepath.Join(homeDir, todosDir, tasksFileName) // Build filepath
  file, err := os.Open(filePath)
  if err != nil {
    fmt.Println("Error opening file:", err)
    return nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file) // Read lines from file, store in lines variable
  var lines []string
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return lines, nil //returs []string containing lines from Todo-list file
}

