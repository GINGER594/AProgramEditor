package texteditor

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "slices"
    "strings"

    "golang.org/x/term"

    cursor "aprogrameditor/cursor"
)

const tab string = "    "

const lowUILineA string = "^[: COMMENT LINE  | ^C: STORE LINES  | ^Z: STORE+DELETE LINES  | ^R: REPLACE ALL     | ^S: SAVE"
const lowUILineB string = "^]: COMMENT LINES | ^X: DELETE LINES | ^V: INSERT STORED LINES | ^L: TOGGLE LINE NO. | ^Q: QUIT"

const CommentUI string = "COMMENT LINES IN RANGE (inclusive) [low/high]: "
const StoreUI string = "STORE LINES IN RANGE (inclusive) [low/high]: "
const DeleteUI string = "DELETE LINES IN RANGE (inclusive) [low/high]: "
const StoreAndDeleteUI string = "STORE+DELETE LINES IN RANGE (inclusive) [low/high]: "
const ReplaceUI string = "REPLACE ALL ['a'/'b']: "

//struct for text editing
type TextEditor struct {
    path string
    comment string
    upperUI string
    lowerUI string
    text []string
    storedLines []string
    termWidth int
    termHeight int
    saved bool
    lineNumToggle bool
    cur cursor.Cursor
}

//takes in path, reads file into text field
func (t *TextEditor) readFile() error {
    iofile, err := os.Open(t.path)
    if err != nil {
        return err
    }
    defer iofile.Close()
    t.saved = true

    reader := bufio.NewReader(iofile)
    for {
        line, err := reader.ReadString(byte('\n'))
        //removing unwanted strings (carriage returns, new lines etc)
        line = strings.ReplaceAll(line, string(byte(0)), "")
        line = strings.ReplaceAll(line, "\r", "")
        line = strings.ReplaceAll(line, "\n", "")
        line = strings.ReplaceAll(line, "\t", tab)
        t.text = append(t.text, line)
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil
        }
    }
    return nil
}

//refreshes internal terminal dimensions/UI fields (t.termHeight depends on UIHeight etc)
func (t *TextEditor) UpdateTerminalAndUIFields() error {
    termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
    if err != nil {
        return err
    }
    UISeparator := strings.Repeat("_", termWidth)
    t.upperUI = fmt.Sprintf("\x1b[0mAPE : %s\r\n%s\r\n", t.path, UISeparator)
    t.lowerUI = fmt.Sprintf("\x1b[0m%s\r\n%s\r\n%s\r\n%s\r\n", UISeparator, lowUILineA, lowUILineB, UISeparator)
    UIHeight := (strings.Count(t.upperUI+t.lowerUI, "\n") + 1)
    t.termWidth, t.termHeight = termWidth, termHeight-(UIHeight+1) //adding 1 to UIHeight so ENTER does not add extra line during ^R input
    return nil
}

//sets up text editor obj (sets path, comment symbol, reads file)
func (t *TextEditor) SetDefaultValues(path, comment string) error {
    t.path = path
    t.comment = comment
    err := t.readFile()
    if err != nil {
        return err
    }
    err = t.UpdateTerminalAndUIFields()
    if err != nil {
        return err
    }
    t.lineNumToggle = true
    t.cur = cursor.Cursor{}
    t.cur.SetDefaultValues()
    return nil
}

func (t *TextEditor) Saved() bool {
    return t.saved
}

func (t *TextEditor) LineNumToggle() bool {
    return t.lineNumToggle
}

//returns fields necessary for output (including UI elements)
func (t *TextEditor) GetOutputFields() (string, string, []string, int, int, int, int, int) {
    return t.upperUI, t.lowerUI, t.text, t.cur.X(), t.cur.Y(), t.cur.TermY(), t.termWidth, t.termHeight
}

//text editing methods:

func (t *TextEditor) AdjustCursorY(n int) {
    t.cur.ScrollVertical(t.text, t.termHeight, n)
}

func (t *TextEditor) AdjustCursorX(n int) {
    t.cur.ScrollHorizontal(t.text, t.termHeight, n)
}

func (t *TextEditor) HomeCursor() {
    t.cur.ScrollToLineStart()
}

func (t *TextEditor) EndCursor() {
    t.cur.ScrollToLineEnd(t.text)
}

//iteration ensures correct view scrolling over lines
func (t *TextEditor) PageUpCursor() {
    for i := 0; i < t.termHeight; i++ {
        t.cur.ScrollVertical(t.text, t.termHeight, -1)
    }
}

func (t *TextEditor) PageDownCursor() {
    for i := 0; i < t.termHeight; i++ {
        t.cur.ScrollVertical(t.text, t.termHeight, 1)
    }
}

//handles backspace when line start (removes line)
func (t *TextEditor) backspaceLine() {
    if t.cur.Y() > 0 {
        line := t.text[t.cur.Y()]
        t.text = append(t.text[:t.cur.Y()], t.text[t.cur.Y()+1:]...)
        t.cur.ScrollVertical(t.text, t.termHeight, -1)
        t.cur.ScrollToLineEnd(t.text)
        t.text[t.cur.Y()] += line
    }
}

func (t *TextEditor) Backspace() {
    if t.cur.X() <= -1 {
        t.backspaceLine()
    } else {
        line := t.text[t.cur.Y()]
        line = line[:t.cur.X()] + line[t.cur.X()+1:]
        t.text[t.cur.Y()] = line
        t.cur.ScrollHorizontal(t.text, t.termHeight, -1)
    }
    t.saved = false
}

//handles delete when at line end (removes line)
func (t *TextEditor) deleteLine() {
    if t.cur.Y() < len(t.text)-1 {
        line := t.text[t.cur.Y()+1]
        t.text = append(t.text[:t.cur.Y()+1], t.text[t.cur.Y()+2:]...)
        t.text[t.cur.Y()] += line
    }
}

func (t *TextEditor) Delete() {
    if t.cur.X() >= len(t.text[t.cur.Y()])-1 {
        t.deleteLine()
    } else {
        line := t.text[t.cur.Y()]
        line = line[:t.cur.X()+1] + line[t.cur.X()+2:]
        t.text[t.cur.Y()] = line
    }
    t.saved = false
}

//inserts \n
func (t *TextEditor) InsertNewLn() {
    line := t.text[t.cur.Y()]
    t.text[t.cur.Y()] = line[:t.cur.X()+1]
    t.text = slices.Insert(t.text, t.cur.Y()+1, line[t.cur.X()+1:])
    t.cur.ScrollVertical(t.text, t.termHeight, 1)
    t.cur.ScrollToLineStart()
    t.saved = false
}

func (t *TextEditor) InsertString(s string) {
    line := t.text[t.cur.Y()]
    if line == "" || t.cur.X() >= len(line) {
        line += s
    } else {
        line = line[:t.cur.X()+1] + s + line[t.cur.X()+1:]
    }
    t.text[t.cur.Y()] = line
    t.cur.ScrollHorizontal(t.text, t.termHeight, len(s))
    t.saved = false
}

func (t *TextEditor) InsertTab() {
    t.InsertString(tab)
}

//removes tab from line start (if prefix exists)
func (t *TextEditor) Untab() {
    line := t.text[t.cur.Y()]
    line, found := strings.CutPrefix(line, tab)
    if found {
        if t.cur.X() < len(tab) {
            t.cur.ScrollToLineStart()
        } else {
            t.cur.ScrollHorizontal(t.text, t.termHeight, -len(tab))
        }
        t.text[t.cur.Y()] = line
        t.saved = false
    }
}

//ctrl-key methods:

//comments line y (and scrolls accordingly if it is the cursor line)
func (t *TextEditor) commentLine(y int) {
    line := t.text[y]
    t.text[y] = t.comment + line
    if y == t.cur.Y() {
        t.cur.ScrollHorizontal(t.text, t.termHeight, len(t.comment))
    }
    t.saved = false
}

//decomments line y (and scrolls accordingly if it is the cursor line)
func (t *TextEditor) decommentLine(y int) {
    line := t.text[y][len(t.comment):]
    if y == t.cur.Y() {
        if t.cur.X() < len(t.comment) {
            t.cur.ScrollToLineStart()
        } else {
            t.cur.ScrollHorizontal(t.text, t.termHeight, -len(t.comment))
        }
    }
    t.text[y] = line
    t.saved = false
}

//^[
func (t *TextEditor) CommentCursorLine() {
    if t.comment != "" {
        line := t.text[t.cur.Y()]
        if strings.HasPrefix(line, t.comment) {
            t.decommentLine(t.cur.Y())
        } else {
            t.commentLine(t.cur.Y())
        }
    }
}

//returns true if over half of the lines in a given range are commented
func (t *TextEditor) overHalfLinesCommented(low, high int) bool {
    commentedCount := 0
    for y := low; y < high; y++ {
        if strings.HasPrefix(t.text[y], t.comment) {
            commentedCount += 1
        }
    }
    if commentedCount > (high-low) / 2 {
        return true
    }
    return false
}

//^] comments/decomments lines between low and high (will decomment if over half of the lines are uncommented and comment otherwise)
func (t *TextEditor) CommentLines(low, high int) {
    if t.comment != "" && low <= high {
        low = max(0, low-1)
        high = min(len(t.text), high)
        overHalfCommented := t.overHalfLinesCommented(low, high)
        for y := low; y < high; y++ {
            if overHalfCommented && strings.HasPrefix(t.text[y], t.comment) {
                t.decommentLine(y)
            } else if !overHalfCommented {
                t.commentLine(y)
            }
        }
    }
}

//stores lines between low and high in storedLines field
func (t *TextEditor) storeLines(low, high int) {
    low = max(0, low-1)
    high = min(len(t.text), high)
    t.storedLines = slices.Clone(t.text[low:high])
}

//^C
func (t *TextEditor) StoreLines(low, high int) {
    if low <= high {
        t.storeLines(low, high)
    }
}

//deletes lines between low and high, moves cursor to (0, low)
func (t *TextEditor) deleteLines(low, high int) {
    low = max(0, low-1)
    high = min(len(t.text), high)
    t.cur.ScrollVertical(t.text, t.termHeight, low-t.cur.Y()-1)
    t.cur.ScrollToLineStart()
    t.text = append(t.text[:low], t.text[high:]...)
    if len(t.text) <= 0 {
        t.text = []string{""}
    }
    t.saved = false
}

//^X
func (t *TextEditor) DeleteLines(low, high int) {
    if low <= high {
        t.deleteLines(low, high)
    }
}

//^Z
func (t *TextEditor) StoreAndDeleteLines(low, high int) {
    if low <= high {
        t.storeLines(low, high)
        t.deleteLines(low, high)
    }
}

//^V inserts storedLines field into text field, moves cursor to (0, c.Y()+len(storedLines))
func (t *TextEditor) InsertStoredLines() {
    if len(t.storedLines) > 0 {
        t.text = append(t.text[:t.cur.Y()], append(t.storedLines, t.text[t.cur.Y():]...)...)

        //scrolling to (0, c.Y()+len(storedLines))
        t.cur.ScrollVertical(t.text, t.termHeight, len(t.storedLines))
        t.cur.ScrollToLineStart()
        t.saved = false
    }
}

//^R replaces all occurrences of a with b in text field
func (t *TextEditor) ReplaceAllOccurrences(a, b string) {
    if a != b { //only replacing if a != b (to prevent t.saved from being incorrectly set)
        for y, line := range t.text {
            line = strings.ReplaceAll(line, a, b)
            t.text[y] = line
            if y == t.cur.Y() {
                if t.cur.X() > len(line) {
                    t.cur.ScrollToLineEnd(t.text)
                }
            }
        }
        t.saved = false
    }
}

//^L
func (t *TextEditor) ToggleLineNumField() {
    t.lineNumToggle = !t.lineNumToggle
}

//^S writes file from text field to file at path field
func (t *TextEditor) SaveFile() error {
    iofile, err := os.Create(t.path)
    if err != nil {
        return err
    }
    defer iofile.Close()

    file := strings.Join(t.text, "\n")
    writer := bufio.NewWriter(iofile)
    _, err = writer.Write([]byte(file))
    writer.Flush() //flushing writer regardless of write-error status
    if err != nil {
        return err
    }
    t.saved = true
    return nil
}
