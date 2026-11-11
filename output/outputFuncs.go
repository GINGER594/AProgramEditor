package output

import (
    "fmt"
    "strings"

    texteditor "aprogrameditor/textEditor"
    huesettings "aprogrameditor/hueSettings"
)

const resetHue string = "\x1b[0m"
const bgResetHue string = "\x1b[49m"
const fgResetHue string = "\x1b[39m"
const curHue string = "\x1b[47m"

//takes in line no. and total lines, returns padded line no. UI element
func generateLineNumUI(lineNum, lineCount int) string {
    strLineNum := fmt.Sprintf("%d", lineNum+1)
    strLineCount := fmt.Sprintf("%d", lineCount)
    whiteSpace := strings.Repeat(" ", len(strLineCount))
    if lineNum >= lineCount {
        return whiteSpace + "| "
    }
    whiteSpace = strings.Repeat(" ", len(strLineCount)-len(strLineNum))
    return whiteSpace + strLineNum + "| "
}

//inserts cursor into highlighted line
func insertCursor(line string, curX int) string {
    xAdj := 0
    realX := -1
    for x := 0; x < len(line); x++ {
        if line[x] == byte(27) {
            escCodeLen := len(fgResetHue)
            xAdj += escCodeLen
            x += escCodeLen - 1
        } else {
            if realX >= curX {
                break
            }
            realX += 1
        }
    }
    curX += xAdj
    if curX >= len(line)-1 {
        return line + curHue + " " + bgResetHue
    }
    return line[:curX+1] + curHue + string(line[curX+1]) + bgResetHue + line[curX+2:]
}

//removes all valid foreground hues
func removeFgHues(str string) string {
    for _, hue := range huesettings.ValidHues {
        str = strings.ReplaceAll(str, "\x1b["+hue+"m", "")
    }
    return str
}

//applies comment hues to line (including in-line comments, edge-cases e.g. comment-symbols in strings)
func applyCommentHue(line, commentSymbol, commentHue string, stringHues map[string]string) string {
    if commentSymbol == "" {
        return line
    }
    //checking if comment symbol is in a string, only applying hue map if not
    currentStrSymbol := ""
    for x := len(commentSymbol); x <= len(line); x++ {
        chr := string(line[x-len(commentSymbol)])
        if _, ok := stringHues[chr]; currentStrSymbol == "" && ok {
            currentStrSymbol = chr
        } else if chr == currentStrSymbol {
            currentStrSymbol = ""
        } else if line[x-len(commentSymbol):x] == commentSymbol && currentStrSymbol == "" {
            prefix := line[:x-len(commentSymbol)]
            suffix := line[x-len(commentSymbol):]
            return prefix + commentHue + removeFgHues(suffix) + fgResetHue
        }
    }
    return line
}

//applies string hues and removes any existing mid-string hues
func applyStringHues(line string, stringHues map[string]string) string {
    newLine := ""
    substrings := []string{}
    substring := ""
    for x := 0; x < len(line); x++ {
        chr := string(line[x])
        if len(substring) == 0 {
            if _, ok := stringHues[chr]; ok {
                newLine += string(byte(0))
                substring += chr
            } else {
                newLine += chr
            }
        } else {
            substring += chr
            if chr == "\\" { //skipping over esc chars in strings
                if x < len(line)-1 {
                    substring += string(line[x+1])
                }
                x += 1
            } else if strSymbol := string(substring[0]); chr == strSymbol {
                hue, _ := stringHues[strSymbol]
                substrings = append(substrings, hue+removeFgHues(substring)+fgResetHue)
                substring = ""
            }
        }
    }
    if len(substring) > 0 {
        hue, _ := stringHues[string(substring[0])]
        substrings = append(substrings, hue+removeFgHues(substring)+fgResetHue)
        substring = ""
    }
    for _, substring := range substrings {
        newLine = strings.Replace(newLine, string(byte(0)), substring, 1)
    }
    return newLine
}

//applies symbol hues to any symbol not part of a comment symbol or preceded by an esc sequence
func applySymbolHues(line, commentSymbol string, symbolHues map[string]string) string {
    for x := 0; x < len(line); x++ {
        b := line[x]
        if b == 27 { //checking for esc sequence
            x += 1
            continue
        }
        if len(commentSymbol) > 0 && x+len(commentSymbol) <= len(line) { //checking if symbol is part of a comment symbol
            if line[x:x+len(commentSymbol)] == commentSymbol {
                x += len(commentSymbol)-1
                continue
            }
        }
        if hue, ok := symbolHues[string(b)]; ok {
            symbol := hue + string(b) + fgResetHue
            line = line[:x] + symbol + line[x+1:]
            x += len(symbol)-1
        }
    }
    return line
}

//returns true if a is not an underscore and not a letter
func isNotUnderscoreOrLetter(a byte) bool {
    if (a != 95) && !(a >= 65 && a <= 90) && !(a >= 97 && a <= 122) {
        return true
    }
    return false
}

//returns longest keyword (and hue) (with no neighbouring letters or underscores) prefix of a given string
func getKeywordPrefix(str string, keywordHues map[string]string) (string, string) {
    foundKeyword := ""
    for keyword := range keywordHues {
        if strings.HasPrefix(str, keyword) {
            nextByte := byte((str+" ")[len(keyword)])
            if isNotUnderscoreOrLetter(nextByte) && len(keyword) > len(foundKeyword) {
                foundKeyword = keyword
            }
        }
    }
    hue, _ := keywordHues[foundKeyword]
    return foundKeyword, hue
}

func applyKeywordHues(line string, keywordHues map[string]string) string {
    var newLine strings.Builder
    line = " " + line //adding space to line so each previous byte can be processed
    for x := 1; x <= len(line); x++ { //iterating over each byte in line and checking for a keyword if b is a non-underscore, non-letter
        b := line[x-1]
        if x == 1 || isNotUnderscoreOrLetter(b) {
            if keyword, hue := getKeywordPrefix(line[x:], keywordHues); keyword != "" {
                newLine.WriteString(string(b) + hue + keyword + fgResetHue)
                x += len(keyword)
                continue
            }
        }
        newLine.WriteString(string(b))
    }
    return strings.TrimPrefix(newLine.String(), " ") //removing extra space
}

func applyHueMap(line string, hueMap *huesettings.HueMap) string {
    line = applyKeywordHues(line, hueMap.Keywords) //keywords
    line = applySymbolHues(line, hueMap.Comment, hueMap.Symbols) //symbols
    line = applyStringHues(line, hueMap.Strings) //strings (applied 2nd-last as any strings need to have all hues removed)
    line = applyCommentHue(line, hueMap.Comment, hueMap.CommentHue, hueMap.Strings) //comments (applied last as any comments need to have all hues removed)
    return line
}

//formats the stored file for output and returns it
func formatTextForTerminal(editor *texteditor.TextEditor, hueMap *huesettings.HueMap) string {
    upperUI, lowerUI, text, curX, curY, termY, termWidth, termHeight := editor.GetOutputFields()

    ftextBuilder := strings.Builder{}
    ftextBuilder.WriteString(resetHue + upperUI)

    //iterating over visible lines
    for y := termY; y < termY+termHeight; y++ {
        lineNumUI := ""
        if editor.LineNumToggle() {
            lineNumUI = generateLineNumUI(y, len(text))
        }
        line := ""
        if y <= len(text)-1 {
            line = text[y]
            //accounting for line wrapping
            termHeight -= (len(lineNumUI)+len(line)) / termWidth //removing 1 line from the bottom of the screen for each line wrap
            if y >= termY+termHeight && len(lineNumUI)+len(line) >= termWidth { //stopping line-wrapping for last visible line
                line = line[:termWidth-len(lineNumUI)-1]
            }
            //applying syntax hues
            line = applyHueMap(line, hueMap)
            if y == curY {
                line = insertCursor(line, curX)
            }
        }
        ftextBuilder.WriteString(resetHue + lineNumUI + line + "\r\n")
    }
    ftextBuilder.WriteString(resetHue + lowerUI)
    return ftextBuilder.String()
}

//clears terminal, outputs text
func OutEditorText(editor *texteditor.TextEditor, hueMap *huesettings.HueMap) {
    ftext := formatTextForTerminal(editor, hueMap)
    fmt.Printf("\x1bc%s", ftext)
}
