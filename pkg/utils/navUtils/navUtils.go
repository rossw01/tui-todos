package navUtils

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const NewTaskHeight = 4

func OpenAddTask(termWidth int, termHeight int, tasks *widgets.List, newTask *widgets.Paragraph) { 
  tasks.SetRect(0, 0, termWidth, termHeight - NewTaskHeight)
  ui.Render(tasks, newTask);
  // newTaskOpen = true
}

func CloseAddTask(termWidth int, termHeight int, tasks *widgets.List) {
  tasks.SetRect(0, 0, termWidth, termHeight)
  ui.Render(tasks)
}

func OpenHelp(termWidth int, termHeight int, helpWindow *widgets.Paragraph) {
  helpWindow.SetRect(termWidth / 8, termHeight / 8, (termWidth / 8) * 7, (termHeight / 8) * 7 )
  ui.Render(helpWindow)
}

func CloseHelp(helpWindow *widgets.Paragraph) {
  helpWindow.SetRect(0,0,0,0)
  ui.Render(helpWindow)
}

