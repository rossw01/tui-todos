package main

import (
	"fmt"
	"log"
	"tui-todos/pkg/utils/cursorUtils"
	"tui-todos/pkg/utils/fileHandler"
	"tui-todos/pkg/utils/navUtils"
	"tui-todos/pkg/utils/taskUtils"
	"tui-todos/pkg/utils/widgetInitializer"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

  termWidth, termHeight := ui.TerminalDimensions()

  fileHandler.CreateDirIfNotExists()
  todoItems, err := fileHandler.LoadTasks()
  if err != nil {
    fmt.Println("Error loading tasks:", err)
    return
  }
	tasks := widgets.NewList()
  widgetInitializer.SetupTasksWidget(tasks, todoItems, termWidth, termHeight)

  newTask := widgets.NewParagraph()
  widgetInitializer.SetupNewTaskWidget(newTask, termWidth, termHeight)

  help := widgets.NewParagraph()
  widgetInitializer.SetupHelpWidget(help, termWidth, termHeight)

  ui.Render(tasks) // Initial UI render

	previousKey := "" // This is used for 'gg' binding for jumping to the first item
  newTaskOpen := false
  newTaskIndex := 0 // 0 -> Below, -1 -> Above
  cursorIndex := 0
	uiEvents := ui.PollEvents()

  var newTaskType taskUtils.TodoType
  // Input handling
	for {
		e := <-uiEvents
    if e.ID == "<C-c>" {
      return
    } 
    // Task-view Inputs
    if !newTaskOpen {
      switch e.ID {
      case "q":
        return
      case "j", "<Down>":
        if len(tasks.Rows) > 0 {
          tasks.ScrollDown()
        }
      case "k", "<Up>":
        if len(tasks.Rows) > 0 { tasks.ScrollUp() }
      case "<Tab>":
        if len(tasks.Rows) > 0 {
          tasks.Rows[tasks.SelectedRow] = taskUtils.CycleColour(tasks.Rows[tasks.SelectedRow])
        }
      case "<C-d>":
        if len(tasks.Rows) > 0 { tasks.ScrollHalfPageDown() }
      case "<C-u>":
        if len(tasks.Rows) > 0 { tasks.ScrollHalfPageUp() }
      case "<C-f>":
        if len(tasks.Rows) > 0 { tasks.ScrollPageDown() }
      case "<C-b>":
        if len(tasks.Rows) > 0 { tasks.ScrollPageUp() }
      case "w":
        fileHandler.SaveTasks(tasks.Rows)
      case "g":
        if previousKey == "g" { if len(tasks.Rows) > 0 { tasks.ScrollTop() } }
      case "<Home>":
        if len(tasks.Rows) > 0 { tasks.ScrollTop() }
      case "G", "<End>":
        if len(tasks.Rows) > 0 { tasks.ScrollBottom() }
      case "i", "o":
        newTaskIndex = 0
        newTaskOpen = true
        newTaskType = taskUtils.Task
        navUtils.OpenAddTask(termWidth, termHeight, tasks, newTask)
      case "a":
        newTaskIndex = 0
        newTaskOpen = true
        newTaskType = taskUtils.Subtask
        navUtils.OpenAddTask(termWidth, termHeight, tasks, newTask)
      case "O":
        newTaskIndex = -1
        newTaskOpen = true
        newTaskType = taskUtils.Task
        navUtils.OpenAddTask(termWidth, termHeight, tasks, newTask)
      case "d", "x":
        taskUtils.RemoveTask(tasks)
      case "?":
        navUtils.OpenHelp(termWidth, termHeight, help)
      case "<Escape>":
        navUtils.CloseHelp(help)
      }
      previousKey = e.ID
    } else if newTaskOpen {
      // New-task window inputs
      switch e.ID {
        case "<Escape>":
          if newTaskOpen {
            newTask.Text = taskUtils.DefaultNewTaskText
            newTaskOpen = false
            navUtils.CloseAddTask(termWidth, termHeight, tasks)
            cursorIndex = 0
          }
        case "<Backspace>":
          if len(newTask.Text) > len(taskUtils.DefaultNewTaskText) && cursorIndex > 0 {
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
          if len(newTask.Text) > len(taskUtils.DefaultNewTaskText) {
            taskUtils.InsertTask(tasks,  newTaskIndex, newTaskType, newTask, termWidth, termHeight);
            newTaskOpen = false
            navUtils.CloseAddTask(termWidth, termHeight, tasks)
            cursorIndex = 0
          }
        case "<Left>":
          if len(newTask.Text) > len(taskUtils.DefaultNewTaskText) && cursorIndex > 0 {
            newTask.Text = cursorUtils.MoveCursor(newTask.Text, cursorIndex - 1);
            if cursorIndex > 0 { cursorIndex-- }
          }
        case "<Right>":
          if len(newTask.Text) > len(taskUtils.DefaultNewTaskText) && cursorIndex < (len(newTask.Text) - len(taskUtils.DefaultNewTaskText)) {
            newTask.Text = cursorUtils.MoveCursor(newTask.Text, cursorIndex + 1);
            if cursorIndex < (len(newTask.Text) - len(taskUtils.DefaultNewTaskText)) {
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
