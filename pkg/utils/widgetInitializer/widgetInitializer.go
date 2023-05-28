package widgetInitializer

import (
	"tui-todos/pkg/utils/navUtils"
	"tui-todos/pkg/utils/taskUtils"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
) 

const (
  helpWindowText =  `
    [i/o] Create new task below currently selected task
    [a] Create new subtask
    [O] Create new task above currently selected task
    [gg] Jump to first task in Todolist
    [G] Jump to last task in Todolist
    [x/d] Delete currently selected task
    [w] Write changes
    [Tab] Cycle text colour of currently selected task
    [q/C-c] Quit program
    [Escape] Cancel new task creation
    [Escape] Close help window
  `
)

func SetupTasksWidget(tasks *widgets.List, todoItems []string, termWidth int, termHeight int) {
  tasks.Title = "Tasks: "
  tasks.Rows = todoItems
	tasks.TextStyle = ui.NewStyle(ui.ColorYellow)
	tasks.WrapText = true; 
  tasks.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	tasks.SetRect(0, 0, termWidth, termHeight)
	tasks.WrapText = true 
}

func SetupNewTaskWidget(newTask *widgets.Paragraph, termWidth int, termHeight int) {
  newTask.Title = "Create a new task"
  newTask.Text = taskUtils.DefaultNewTaskText 
  newTask.SetRect(0, termHeight - navUtils.NewTaskHeight, termWidth, termHeight)
}

func SetupHelpWidget(help *widgets.Paragraph, termWidth int, termHeight int) {
  help.Title = "Keybinds: "
  help.Text = helpWindowText
}
