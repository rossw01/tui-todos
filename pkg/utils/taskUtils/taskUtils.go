package taskUtils 

import (
  "regexp"
  "strings"
  "strconv"
	"github.com/gizak/termui/v3/widgets"
  "tui-todos/pkg/utils/cursorUtils"
  "tui-todos/pkg/utils/navUtils"
)

type TodoType int64

const (
  DefaultNewTaskText = "[ ](bg:white)"
  Task TodoType = 0
  Subtask TodoType = 1
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

func UpdateStringIndex(tasks *widgets.List, index int) () {
  for i := index; i < len(tasks.Rows); i++ {
    tasks.Rows[i] = "[" + strconv.Itoa(i) + "]" + tasks.Rows[i][strings.Index(tasks.Rows[i], "]") + 1:]
  }
}

func RemoveTask(tasks *widgets.List, ) {
  if len(tasks.Rows) > 0 && tasks.SelectedRow >= 0 && tasks.SelectedRow < len(tasks.Rows) {
    index := tasks.SelectedRow
    tasks.Rows = append(tasks.Rows[:tasks.SelectedRow], tasks.Rows[tasks.SelectedRow+1:]...)
    if !(len(tasks.Rows) > tasks.SelectedRow) {
      tasks.ScrollUp();
    }
    UpdateStringIndex(tasks, index)
  }
}

func InsertTask(tasks *widgets.List, offset int, newTaskType TodoType, newTask *widgets.Paragraph, termWidth int, termHeight int) {
  // Offset is for placing new task above/below the currently selected item
  if len(newTask.Text) > 0 {
    index := tasks.SelectedRow
    if len(tasks.Rows) == 0 {
      tasks.Rows = []string{"[0] " + cursorUtils.RemoveCursor(newTask.Text)}
      newTask.Text = DefaultNewTaskText
      navUtils.CloseAddTask(termWidth, termHeight, tasks)
    } else {
      var task string
      if newTaskType == Task {
        task = "[" + strconv.Itoa(index + 1 + offset) + "] " + cursorUtils.RemoveCursor(newTask.Text)
      } else if newTaskType == Subtask {
        task = "[" + strconv.Itoa(index + 1 + offset) + "] └─ " + cursorUtils.RemoveCursor(newTask.Text)
      }
      tasks.Rows = append(tasks.Rows[:index + 1 + offset], append([]string{task}, tasks.Rows[index + 1 + offset:]...)...)
    }
    newTask.Text = DefaultNewTaskText
    UpdateStringIndex(tasks, index + 1 + offset)
  }
}

func GetMatchedColour(selectedRow string) ([][]string) {
  regex := regexp.MustCompile(`\(fg.*?\)`)
  return regex.FindAllStringSubmatch(selectedRow, -1) // Find final regex string match
}

func CycleColour(selectedRow string) (string) {
  matches := GetMatchedColour(selectedRow)
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

