package cursorUtils 

import (
  "regexp"
)

// this is insane
// How can there be no built in TextInput widget built into termui?
// Theres literally a pull request for it that's been sat there for 6+ years lol

func RemoveCursor(inputtedText string) (string) {
  pattern := regexp.MustCompile(`\[[^\]]*\]\(bg:white\)`)
  match := pattern.FindStringSubmatchIndex(inputtedText)
  charToPreserve := inputtedText[match[0] + 1]
  return inputtedText[:match[0]] + string(charToPreserve) + inputtedText[match[1]:]
}

func MoveCursor(inputtedText string, indexOfNextChar int) (string) {
  newString := RemoveCursor(inputtedText)
  return newString[:indexOfNextChar] + "[" + string(newString[indexOfNextChar]) + "](bg:white)" + newString[indexOfNextChar + 1:]
}

