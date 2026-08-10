package texteditor

import (
    "bufio"
    "fmt"
    "io"
    cursor "ninthcircleoftext/cursor"
    "os"
    "slices"
    "strings"

    "golang.org/x/term"
)

const lowUILine1 string = "^?: COMMENT LINE | ^C: STORE LINES         | ^R: REPLACE ALL     | ^S: SAVE"
const lowUILine2 string = "^X: DELETE LINES | ^V: INSERT STORED LINES | ^T: TOGGLE LINE NO. | ^Q: QUIT"

const DeleteUI string = "DELETE LINES IN RANGE (inclusive) [low/high]: "
const CopyUI string = "STORE LINES IN RANGE (inclusive) [low/high]: "
const ReplaceUI string = "REPLACE ALL ['a'/'b']: "

//struct for handling text editing (insert, delete etc)
type TextEditor struct { //fields are organised by type
    path string
    tab string
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

//takes in a path & reads the file into the text field of the TextEditor struct
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
        //removing any unwanted strings (carriage returns, new lines etc)
        line = strings.ReplaceAll(line, string(byte(0)), "")
        line = strings.ReplaceAll(line, "\r", "")
        line = strings.ReplaceAll(line, "\n", "")
        line = strings.ReplaceAll(line, "\t", t.tab)
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

//refreshes the internal terminal dimension & UI fields
func (t *TextEditor) UpdateTerminalAndUIFields() error {
    termWidth, termheight, err := term.GetSize(int(os.Stdout.Fd()))
    if err != nil {
        return err
    }
    UISeparator := strings.Repeat("_", termWidth)
    t.upperUI = fmt.Sprintf("\x1b[0m%s\r\n%s\r\n", t.path, UISeparator)
    t.lowerUI = fmt.Sprintf("\x1b[0m%s\r\n%s\r\n%s\r\n%s\r\n", UISeparator, lowUILine1, lowUILine2, UISeparator)
    UIHeight := (strings.Count(t.upperUI+t.lowerUI, "\n") + 1)
    t.termWidth, t.termHeight = termWidth, termheight-(UIHeight+1) //adding 1 to UIHeight so that during ^R input, ENTER does not add an extra line
    return nil
}

//sets up the text editor (sets path, comment symbol, reads file)
func (t *TextEditor) SetDefaultValues(path, comment string) error {
    t.path = path
    t.tab = "    "
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

//returns the fields necessary for output (including UI elements)
func (t *TextEditor) GetOutputFields() (string, string, []string, int, int, int, int) {
    return t.upperUI, t.lowerUI, t.text, t.cur.X(), t.cur.Y(), t.cur.TermY(), t.termHeight
}

//ACTUAL TEXT-EDITING METHODS START HERE:
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

//iteration ensures correct view scrolling over lines
func (t *TextEditor) PageDownCursor() {
    for i := 0; i < t.termHeight; i++ {
        t.cur.ScrollVertical(t.text, t.termHeight, 1)
    }
}

//handles backspace when at the start of a line (removes line)
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

//handles delete when at the end of a line (removes line)
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
func (t *TextEditor) InsertLine() {
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
    t.InsertString(t.tab)
}

//CTRL-KEY FUNCTIONS:
//comments any uncommented line, uncomments any commented line
func (t *TextEditor) commentLine() {
    line := t.text[t.cur.Y()]
    if trimmedLine := strings.TrimSpace(line); strings.HasPrefix(trimmedLine, t.comment) { //uncommenting line if it starts with the comment symbol
        if t.cur.X() < len(t.comment) {
            t.cur.ScrollToLineStart()
        } else {
            t.cur.ScrollHorizontal(t.text, t.termHeight, -len(t.comment))
        }
        t.text[t.cur.Y()] = strings.Replace(line, t.comment, "", 1)
    } else { //commenting line
        t.text[t.cur.Y()] = t.comment + line
        t.cur.ScrollHorizontal(t.text, t.termHeight, len(t.comment))
    }
    t.saved = false
}

func (t *TextEditor) CommentLine() {
    if t.comment != "" { //only commenting line if comment is not empty (to stop saved from being incorrectly set to true)
        t.commentLine()
    }
}

//deletes lines between low & high & moves the cursor to (0, low)
func (t *TextEditor) DeleteLines(low, high int) {
    if low <= high {
        low = max(0, low-1)
        high = min(len(t.text)-1, high-1)
        t.cur.ScrollVertical(t.text, t.termHeight, low-t.cur.Y()-1)
        t.cur.ScrollToLineStart()
        t.text = append(t.text[:low], t.text[high+1:]...)
        if len(t.text) <= 0 {
            t.text = []string{""}
        }
        t.saved = false
    }
}

//stores lines between low & high in the storedLines field
func (t *TextEditor) StoreLines(low, high int) {
    if low <= high {
        low = max(0, low-1)
        high = min(len(t.text)-1, high-1)
        t.storedLines = slices.Clone(t.text[low : high+1])
    }
}

//inserts the storedLines field into the text field at the next line after cursor Y & moves the cursor to the end of the newly-inserted lines
func (t *TextEditor) InsertStoredLines() {
    if len(t.storedLines) > 0 {
        t.text = append(t.text[:t.cur.Y()+1], append(t.storedLines, t.text[t.cur.Y()+1:]...)...)

        //scrolling to the end of the inserted lines
        t.cur.ScrollVertical(t.text, t.termHeight, len(t.storedLines))
        t.cur.ScrollToLineEnd(t.text)
        t.saved = false
    }
}

//replaces all occurrences of a with b in the text field
func (t *TextEditor) ReplaceAllOccurrences(a, b string) {
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

//method for toggling on/off line numbers
func (t *TextEditor) ToggleLineNumField() {
    t.lineNumToggle = !t.lineNumToggle
}

//writes the file from the text field to the file at the path field
func (t *TextEditor) SaveFile() error {
    iofile, err := os.Create(t.path)
    if err != nil {
        return err
    }
    defer iofile.Close()

    file := strings.Join(t.text, "\n")
    writer := bufio.NewWriter(iofile)
    _, err = writer.Write([]byte(file))
    writer.Flush() //flushing the writer regardless of write-error status
    if err != nil {
        return err
    }
    t.saved = true
    return nil
}
