package main

import (
  "fmt"
  "bufio"
	"log"
  "strconv"
  "strings"
  "os"
  "path/filepath"
  "regexp"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	todosDir      = ".todos"
	tasksFileName = "tasks.txt"
  defaultNewTaskText = "[ ](bg:white)"
)

const (
  helpWindowText =  `
    [a/i/o] Create new task below currently selected task
    [O] Create new task above currently selected task
    [x/d] Delete currently selected task
    [w] Write changes
    [Tab] Cycle text colour of currently selected task
    [q/C-c] Quit program
    [Escape] Cancel new task creation
    [Escape] Close help window
  `
)
var fontColours = []string {
  "(fg:red)",
  "(fg:green)",
  "(fg:yellow)",
  "(fg:blue)",
  "(fg:magenta)",
  "(fg:cyan)",
  "(fg:white)",
  "(fg:default)",
}

var initialTasks = []string {
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

func getMatchedColour(selectedRow string) ([][]string) {
  regex := regexp.MustCompile(`\(fg.*?\)`)
  return regex.FindAllStringSubmatch(selectedRow, -1) // Find final regex string match
}

func cycleColour(selectedRow string) (string) {
  matches := getMatchedColour(selectedRow)

  if len(matches) > 0 {
    matchedColour := matches[len(matches) - 1][0] // Get last regex match
		selectedRowWithoutColour := strings.Replace(selectedRow, matchedColour, "", -1)

    fontColoursIndex := 0

    for i, s := range fontColours {
      if strings.Contains(s, matchedColour) {
        fontColoursIndex = i
        break
      }
    }

    if fontColoursIndex >= len(fontColours) - 1 {
      fontColoursIndex = 0
    } else {
      fontColoursIndex++
    }
    return selectedRowWithoutColour + fontColours[fontColoursIndex]
  } else {
    // Since there's no colour appended to the end, just append 1st fontColour and return
    leftBracketIndex := strings.Index(selectedRow, "]")
    modifiedStr := selectedRow[:leftBracketIndex + 2] +
      "[" + selectedRow[leftBracketIndex + 2:]
    return modifiedStr + "]" + fontColours[0]
  }
}

func removeCursor(inputtedText string) (string) {
  pattern := regexp.MustCompile(`\[[^\]]*\]\(bg:white\)`)
  match := pattern.FindStringSubmatchIndex(inputtedText)
  charToPreserve := inputtedText[match[0] + 1]
  return inputtedText[:match[0]] + string(charToPreserve) + inputtedText[match[1]:]
}

func moveCursor(inputtedText string, indexOfNextChar int) (string) {
  // this is insane
  // How can there be no built in TextInput widget built into termui?
  // Theres literally a pull request for it that's been sat there for 6+ years lol
  // why do i have to make this schizo madman function
  newString := removeCursor(inputtedText)
  return newString[:indexOfNextChar] + "[" + string(newString[indexOfNextChar]) + "](bg:white)" + newString[indexOfNextChar + 1:]
}

func openHelp(termWidth int, termHeight int, helpWindow *widgets.Paragraph) {
  helpWindow.SetRect(termWidth / 8, termHeight / 8, (termWidth / 8) * 7, (termHeight / 8) * 7 )
  ui.Render(helpWindow)
}

func closeHelp(helpWindow *widgets.Paragraph) {
  helpWindow.SetRect(0,0,0,0)
  ui.Render(helpWindow)
}
  
  // var openAddTask = func() {
  //   tasks.SetRect(0, 0, termWidth, termHeight - newTaskHeight)
  //   ui.Render(tasks, newTask) // Test if this can be removed
  //   newTaskOpen = true
  // }
  //
  // var closeAddTask = func() {
  //   tasks.SetRect(0, 0, termWidth, termHeight)
  //   ui.Render(tasks) // Test if this can be removed
  //   newTaskOpen = false
  // }

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

  termWidth, termHeight := ui.TerminalDimensions()

  newTaskOpen := false // Used to check which inputs should be handled
	previousKey := "" // This is used for 'gg' binding for jumping to the first item
  newTaskIndex := 0 // 0 -> Below, -1 -> Above
  cursorIndex := 0
	uiEvents := ui.PollEvents()

  createDirIfNotExists()

  todoItems, err := loadTasks()
  if err != nil {
    fmt.Println("Error loading tasks:", err)
    return
  }
	tasks := widgets.NewList()
  tasks.Title = "Tasks: "
  tasks.Rows = todoItems
	tasks.TextStyle = ui.NewStyle(ui.ColorYellow)
	tasks.WrapText = true; 
  tasks.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	tasks.SetRect(0, 0, termWidth, termHeight)
	tasks.WrapText = true 

  newTask := widgets.NewParagraph()
  newTask.Title = "Create a new task"
  newTask.Text = defaultNewTaskText 
  newTaskHeight := 4
  newTask.SetRect(0, termHeight - newTaskHeight, termWidth, termHeight)

  help := widgets.NewParagraph()
  help.Title = "Keybinds: "
  help.Text = helpWindowText

  var openAddTask = func() {
    tasks.SetRect(0, 0, termWidth, termHeight - newTaskHeight)
    ui.Render(tasks, newTask, help)
    newTaskOpen = true
  }

  var closeAddTask = func() {
    tasks.SetRect(0, 0, termWidth, termHeight)
    ui.Render(tasks, help) // Test if this can be removed
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

  var insertTask = func(offset int) { 
  // Offset is for placing new task above/below the currently selected item
    if len(newTask.Text) > 0 {
      index := tasks.SelectedRow
      if len(tasks.Rows) == 0 {
        tasks.Rows = []string{"[0] " + newTask.Text}
        newTask.Text = defaultNewTaskText
        closeAddTask()
      } else {
        task := "[" + strconv.Itoa(index + 1 + offset) + "] " + removeCursor(newTask.Text)
        tasks.Rows = append(tasks.Rows[:index + 1 + offset], append([]string{task}, tasks.Rows[index + 1 + offset:]...)...)
      }
      newTask.Text = defaultNewTaskText
      updateStringIndex(index + 1 + offset) // Not sure if doing the +1 is an optimisation or a slowdown lol
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
    } 
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
      case "<Tab>":
        if len(tasks.Rows) > 0 {
          tasks.Rows[tasks.SelectedRow] = cycleColour(tasks.Rows[tasks.SelectedRow])
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
      case "a", "i", "o":
        newTaskIndex = 0;
        openAddTask();
      case "O":
        newTaskIndex = -1;
        openAddTask();
      case "d", "x":
        removeTask();
      case "?":
        openHelp(termWidth, termHeight, help)
      case "<Escape>":
        closeHelp(help)
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
            newTask.Text = defaultNewTaskText
            closeAddTask();
            cursorIndex = 0
          }
        case "<Backspace>":
          if len(newTask.Text) > len(defaultNewTaskText) && cursorIndex > 0 {
            newTask.Text = newTask.Text[:cursorIndex - 1] + newTask.Text[cursorIndex:]
            if cursorIndex > 0 { cursorIndex-- }
          }
        case "<Space>":
          newTask.Text = newTask.Text[:cursorIndex] + " " + newTask.Text[cursorIndex:]
          cursorIndex++
        case "<S-Space>": // This one isn't working in some terminals
          newTask.Text = newTask.Text[:cursorIndex] + " " + newTask.Text[cursorIndex:]
          cursorIndex++
        case "<Enter>":
          if len(newTask.Text) > len(defaultNewTaskText) {
            insertTask(newTaskIndex);
            closeAddTask()
            cursorIndex = 0
          }
        case "<Left>":
          if len(newTask.Text) > len(defaultNewTaskText) && cursorIndex > 0 {
            newTask.Text = moveCursor(newTask.Text, cursorIndex - 1);
            if cursorIndex > 0 { cursorIndex-- }
          }
        case "<Right>":
          if len(newTask.Text) > len(defaultNewTaskText) && cursorIndex < (len(newTask.Text) - len(defaultNewTaskText)) {
            newTask.Text = moveCursor(newTask.Text, cursorIndex + 1);
            if cursorIndex < (len(newTask.Text) - len(defaultNewTaskText)) {
              cursorIndex++
            }
          } 
        default:
          if len(string(e.ID)) == 1 { // Gets rid of things like <End> being inserted and breaking everything
            newTask.Text = newTask.Text[:cursorIndex] + string(e.ID) + newTask.Text[cursorIndex:]
            cursorIndex++
          }
      }
      ui.Render(newTask, help)
    }
    ui.Render(tasks, help)
  }
}
