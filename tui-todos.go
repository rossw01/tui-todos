package main

import (
  "fmt"
  "bufio"
	"log"
  "strconv"
  "strings"
  "os"
  "path/filepath"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	todosDir      = ".todos"
	tasksFileName = "tasks.txt"
)

var initialTasks = []string{
	"[0] test0",
	"[1] [test1](fg:blue)",
	"[2] [test2](fg:red)",
	"[3] [test3](fg:white) output",
	"[4] go to shop",
	"[5] buy loads of eggs",
	"[6] make massive omlette",
	"[7] consume",
	"[8] do dishes",
	"[9] grow big muscles from the protein",
	"[10] ???",
	"[11] profit",
}

func createDirIfNotExists() {
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

func saveTasks(tasks []string) {
  homeDir, _ := os.UserHomeDir() 
  filePath := filepath.Join(homeDir, todosDir, tasksFileName)
  file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

  for _, task := range tasks {
    _, _ = file.WriteString(task + "\n")
  }
}

func loadTasks() ([]string, error) {
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

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

  termWidth, termHeight := ui.TerminalDimensions()

	tasks := widgets.NewList()
  tasks.Title = "Tasks:"
  newTaskOpen := false // Used to check which inputs should be handled
	previousKey := "" // This is used for 'gg' binding for jumping to the first item
	uiEvents := ui.PollEvents()

  createDirIfNotExists()

  todoItems, err := loadTasks()
  if err != nil {
    fmt.Println("Error loading tasks:", err)
    return
  }
  tasks.Rows = todoItems

  // Task Window Styling
	tasks.TextStyle = ui.NewStyle(ui.ColorYellow)
	tasks.WrapText = true; 
  tasks.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	tasks.SetRect(0, 0, termWidth, termHeight)

  newTask := widgets.NewParagraph()
  // New Task Window Styling
  newTask.Title = "Create a new task"
  newTask.Text = ""
	tasks.WrapText = true 
  newTaskHeight := 4
  newTask.SetRect(0, termHeight - newTaskHeight, termWidth, termHeight)

  var openAddTask = func() {
    tasks.SetRect(0, 0, termWidth, termHeight - newTaskHeight)
    ui.Render(tasks, newTask)
    newTaskOpen = true
  }

  var closeAddTask = func() {
    tasks.SetRect(0, 0, termWidth, termHeight)
    ui.Render(tasks)
    newTaskOpen = false
  }

  var updateStringIndex = func(index int) {
		for i := index; i < len(tasks.Rows); i++ {
			tasks.Rows[i] = "[" + strconv.Itoa(i) + "]" + tasks.Rows[i][strings.Index(tasks.Rows[i], "]")+1:]
		}
  }

  var removeTask = func() {
    if len(tasks.Rows) > 0 && tasks.SelectedRow >= 0 && tasks.SelectedRow < len(tasks.Rows) {
      index := tasks.SelectedRow
      tasks.Rows = append(tasks.Rows[:tasks.SelectedRow], tasks.Rows[tasks.SelectedRow+1:]...)
      if !(len(tasks.Rows) > tasks.SelectedRow) {
        tasks.ScrollUp();
      }
      updateStringIndex(index)
    }
  }

  var insertTask = func() {
    if len(newTask.Text) > 0 {
      index := tasks.SelectedRow
      if len(tasks.Rows) == 0 {
        tasks.Rows = []string{"[0] " + newTask.Text}
        newTask.Text = ""
        closeAddTask()
      } else {
        task := "[" + strconv.Itoa(index+1) + "] " + newTask.Text
        tasks.Rows = append(tasks.Rows[:index+1], append([]string{task}, tasks.Rows[index+1:]...)...)
      }
      newTask.Text = ""
      updateStringIndex(index + 1) // Not sure if doing the +1 is an optimisation or a slowdown lol
    }
  }

  ui.Render(tasks) // Initial UI render

  // Input handling
	for {
		e := <-uiEvents

    // Global Inputs
    switch e.ID {
      case "<C-c>": // You should always be able to close the program at any time
        return      // with <C-c>
    } // First renderer

    if !newTaskOpen {

      // Task-view Inputs
      switch e.ID {
      case "q", "<C-c>": // I find it annoying if the same key to exit text input mode
        return           // is the same as the key used to close the program.
      case "j", "<Down>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollDown()
        }
      case "k", "<Up>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollUp()
        }
      case "<C-d>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollHalfPageDown()
        }
      case "<C-u>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollHalfPageUp()
        }
      case "<C-f>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollPageDown()
        }
      case "<C-b>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollPageUp()
        }
      case "w":
        saveTasks(tasks.Rows)
      case "g":
        if previousKey == "g" {
          if len(tasks.Rows) > 0 {
            tasks.ScrollTop()
          }
        }
      case "<Home>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollTop()
        }
      case "G", "<End>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollBottom()
        }
      case "a", "i":
        openAddTask();
      case "d", "x":
        removeTask();
      }
      // used to help check for "gg"
      if previousKey == "g" {
        previousKey = ""
      } else {
        previousKey = e.ID
      }
    } else {
      // Add-task window inputs
      switch e.ID {

        case "<Escape>":
          if newTaskOpen {
            newTask.Text = ""
            closeAddTask();
          }
        case "<Backspace>":
          if len(newTask.Text) > 0 {
            newTask.Text = newTask.Text[:len(newTask.Text)-1]
            ui.Render(newTask)
          }
        case "<Space>":
          newTask.Text += " "
          ui.Render(newTask)
        case "<S-Space>": // This one isn't working, shift space still closes newTask
          newTask.Text += " "
          ui.Render(newTask)
        case "<Enter>":
          if len(newTask.Text) > 0 {
            insertTask();
            closeAddTask()
          }
        default:
          newTask.Text += e.ID
          ui.Render(newTask)
      }
    }
    ui.Render(tasks)
  }
}
