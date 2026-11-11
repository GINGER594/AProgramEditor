package main

import (
    "bufio"
    "fmt"
    "os"
    "slices"
    "strconv"
    "strings"

    "golang.org/x/term"

    huesettings "aprogrameditor/hueSettings"
    output "aprogrameditor/output"
    texteditor "aprogrameditor/textEditor"
)

const helpFlag string = "-h"
const settingsFlag string = "-s"
const debugFlag string = "-d"
const helpMsg string =
                "usage: ape [FILE] [OPTIONS]" +
                "\n\noptions:" +
                "\n-h: show this help text" +
                "\n\n-s: open settings:" +
                "\n    usage: [EXT] [TYPE] [VALUE] [ANSI HUE CODE]" +
                "\n    types:" +
                "\n        keyword    any string of letters or underscores" +
                "\n        symbols    any single ascii symbol (excluding whitespace, numbers, underscores and string symbols)" +
                "\n        strings    any one of: \", ', `" +
                "\n        comment    any string" +
                "\n    ANSI hue codes:" +
                "\n        black      30" +
                "\n        red        31" +
                "\n        green      32" +
                "\n        yellow     33" +
                "\n        blue       34" +
                "\n        magenta    35" +
                "\n        cyan       36" +
                "\n        white      37" +
                "\n        default    39" +
                "\n\n-d: debug mode (opens file without settings)" +
                "\n\nAPE v1" +
                "\n"

var keyByteMap map[string][]byte = map[string][]byte{
    "up": {27, 91, 65, 0},
    "down": {27, 91, 66, 0},
    "right": {27, 91, 67, 0},
    "left": {27, 91, 68, 0},
    "home": {27, 91, 72, 0},
    "end": {27, 91, 70, 0},
    "pg up": {27, 91, 53, 126},
    "pg dn": {27, 91, 54, 126},
    "backspace": {127, 0, 0, 0},
    "delete": {27, 91, 51, 126},
    "enter": {13, 0, 0, 0},
    "tab": {9, 0, 0, 0},
    "shift+tab": {27, 91, 90, 0},
    "^[": {27, 0, 0, 0},
    "^]": {29, 0, 0, 0},
    "^C": {3, 0, 0, 0},
    "^X": {24, 0, 0, 0},
    "^Z": {26, 0, 0, 0},
    "^V": {22, 0, 0, 0},
    "^R": {18, 0, 0, 0},
    "^L": {12, 0, 0, 0},
    "^S": {19, 0, 0, 0},
    "^Q": {17, 0, 0, 0},
}

//if file has not been saved the user will be given the option to save it
func saveAndExit(editor *texteditor.TextEditor) error {
    fmt.Printf("\x1bc")
    if editor.Saved() {
        return nil
    }

    fmt.Printf("file may contain unsaved changes - save [y/N]?: ")
    inp := ""
    fmt.Scanln(&inp) //using fmt scan to accept only the first arg in input
    fmt.Printf("\x1bc")
    if inp != "y" && inp != "Y" {
        return nil
    }
    err := editor.SaveFile()
    if err != nil {
        return err
    }
    return nil
}

//(to be used when term state = raw) switches terminal state, gets user input
func getNonRawInp(oldState *term.State, inpMsg string) (string, error) {
    err := term.Restore(int(os.Stdin.Fd()), oldState)
    if err != nil {
        return "", err
    }
    defer term.MakeRaw(int(os.Stdin.Fd()))
    fmt.Printf("\x1b[0m%s", inpMsg)
    inpScanner := bufio.NewReader(os.Stdin)
    inp, _, err := inpScanner.ReadLine()
    return string(inp), nil
}

//parses non-raw user input into a, b for TextEditor.ReplaceAllOccurrences, returns true if input is valid
func getReplaceAllInp(oldState *term.State, inpMsg string) (string, string, bool, error) {
    inp, err := getNonRawInp(oldState, inpMsg)
    if err != nil {
        return "", "", false, err
    }

    //parsing input
    splitInp := strings.Split(string(inp), "'/'")
    if len(splitInp) != 2 {
        return "", "", false, nil
    }
    a, b := splitInp[0], splitInp[1]
    trimA, trimB := strings.TrimPrefix(a, "'"), strings.TrimSuffix(b, "'")
    if a == trimA || b == trimB || len(trimA) == 0 { //checking for surrounding ' and if a is empty
        return "", "", false, nil
    }
    return trimA, trimB, true, nil
}

//parses non-raw user input into numerical a, b for TextEditor methods that take in numerical range, returns true if input is valid
func getRangeInp(oldState *term.State, inpMsg string) (int, int, bool, error) {
    inp, err := getNonRawInp(oldState, inpMsg)
    if err != nil {
        return 0, 0, false, err
    }

    //parsing input
    splitInp := strings.Split(inp, "/")
    if len(splitInp) != 2 {
        return 0, 0, false, nil
    }
    a, err := strconv.Atoi(splitInp[0])
    if err != nil {
        return 0, 0, false, nil
    }
    b, err := strconv.Atoi(splitInp[1])
    if err != nil {
        return 0, 0, false, nil
    }
    return a, b, true, nil
}

//handles raw inputs, calling edit/display functions/methods
func editMode(editor *texteditor.TextEditor, hueMap *huesettings.HueMap) error {
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        return err
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState)

    inpScanner := bufio.NewReader(os.Stdin)
    for editLoop := true; editLoop; {
        err := editor.UpdateTerminalAndUIFields()
        if err != nil {
            return err
        }
        output.OutEditorText(editor, hueMap)

        inp := make([]byte, 4)
        inpScanner.Read(inp)
        switch {
        case slices.Equal(inp, keyByteMap["up"]):
            editor.AdjustCursorY(-1)
        case slices.Equal(inp, keyByteMap["down"]):
            editor.AdjustCursorY(1)
        case slices.Equal(inp, keyByteMap["right"]):
            editor.AdjustCursorX(1)
        case slices.Equal(inp, keyByteMap["left"]):
            editor.AdjustCursorX(-1)
        case slices.Equal(inp, keyByteMap["home"]):
            editor.HomeCursor()
        case slices.Equal(inp, keyByteMap["end"]):
            editor.EndCursor()
        case slices.Equal(inp, keyByteMap["pg up"]):
            editor.PageUpCursor()
        case slices.Equal(inp, keyByteMap["pg dn"]):
            editor.PageDownCursor()
        case slices.Equal(inp, keyByteMap["backspace"]):
            editor.Backspace()
        case slices.Equal(inp, keyByteMap["delete"]):
            editor.Delete()
        case slices.Equal(inp, keyByteMap["enter"]):
            editor.InsertNewLn()
        case slices.Equal(inp, keyByteMap["tab"]):
            editor.InsertTab()
        case slices.Equal(inp, keyByteMap["shift+tab"]):
            editor.Untab()
        //ctrl key methods
        case slices.Equal(inp, keyByteMap["^["]):
            editor.CommentCursorLine()
        case slices.Equal(inp, keyByteMap["^]"]):
            a, b, valid, err := getRangeInp(oldState, texteditor.CommentUI)
            if err != nil {
                return err
            } else if valid {
                editor.CommentLines(a, b)
            }
        case slices.Equal(inp, keyByteMap["^C"]):
            a, b, valid, err := getRangeInp(oldState, texteditor.StoreUI)
            if err != nil {
                return err
            } else if valid {
                editor.StoreLines(a, b)
            }
        case slices.Equal(inp, keyByteMap["^X"]):
            a, b, valid, err := getRangeInp(oldState, texteditor.DeleteUI)
            if err != nil {
                return err
            } else if valid {
                editor.DeleteLines(a, b)
            }
        case slices.Equal(inp, keyByteMap["^Z"]):
            a, b, valid, err := getRangeInp(oldState, texteditor.StoreAndDeleteUI)
            if err != nil {
                return err
            } else if valid {
                editor.StoreAndDeleteLines(a, b)
            }
        case slices.Equal(inp, keyByteMap["^V"]):
            editor.InsertStoredLines()
        case slices.Equal(inp, keyByteMap["^R"]):
            a, b, valid, err := getReplaceAllInp(oldState, texteditor.ReplaceUI)
            if err != nil {
                return err
            } else if valid {
                editor.ReplaceAllOccurrences(a, b)
            }
        case slices.Equal(inp, keyByteMap["^L"]):
            editor.ToggleLineNumField()
        case slices.Equal(inp, keyByteMap["^S"]):
            err := editor.SaveFile()
            if err != nil {
                return err
            }
        case slices.Equal(inp, keyByteMap["^Q"]):
            editLoop = false
        default:
            //checking if inp is a symbol or letter and not an unwanted ctrl/esc sequence
            if inp[0] >= 32 && inp[0] <= 126 && inp[1] == 0 && inp[2] == 0 && inp[3] == 0 {
                editor.InsertString(string(inp[0]))
            }
        }
    }
    return nil
}

//returns path, file ext, help flag status, settingsFlag status and debug flag status (opening file without reading settings)
func parseOsArgs() (string, string, bool, bool, bool, error) {
    //handling special flags (help and settings)
    if slices.Contains(os.Args, helpFlag) {
        return "", "", true, false, false, nil
    } else if slices.Contains(os.Args, settingsFlag) {
        path, err := huesettings.GetSettingsPath()
        if err != nil {
            return "", "", false, false, false, fmt.Errorf("error getting path to settings file: %s", err.Error())
        }
        return path, "", false, true, false, nil
    }

    if len(os.Args) != 2 && len(os.Args) != 3 {
        return "", "", false, false, false, fmt.Errorf("invalid no. of args recieved - expected: 1, got: %d", len(os.Args)-1)
    }
    //getting file extension
    path, ext := os.Args[1], ""
    unhiddenPath, _ := strings.CutPrefix(path, ".") //removing '.' prefix (hidden files)
    splitPath := strings.Split(unhiddenPath, ".")
    if len(splitPath) > 0 {
        ext = strings.Split(unhiddenPath, ".")[len(splitPath)-1]
    }
    //handling '-d' modifier flag
    debugMode := false
    if len(os.Args) == 3 {
        if os.Args[2] == debugFlag {
            debugMode = true
        } else {
            return "", "", false, false, false, fmt.Errorf("invalid modifier flag: '%s'", os.Args[2])
        }
    }
    return path, ext, false, false, debugMode, nil
}

func runTextEditor() error {
    //parsing args
    path, ext, help, settingsMode, debugMode, err := parseOsArgs()
    if err != nil {
        return fmt.Errorf("error during arg validation: %s", err.Error())
    }
    if help {
        fmt.Printf(helpMsg)
        return nil
    }
    hueMap := huesettings.HueMap{
        Keywords: map[string]string{},
        Symbols: map[string]string{},
        Strings: map[string]string{},
    }
    if !debugMode {
        //ensuring existence of settings file (if not debug mode)
        err = huesettings.CreateSettingsFileIfNotExists()
        if err != nil {
            return fmt.Errorf("error ensuring existence of settings file: %s", err.Error())
        }
        //processing syntax hue settings (if not settings or debug mode)
        if !settingsMode {
            hueMap, err = huesettings.GetSyntaxHueSettings(ext)
            if err != nil {
                return fmt.Errorf("error processing settings: %s", err.Error())
            }
        }
    }
    //setting up editor
    editor := texteditor.TextEditor{}
    err = editor.SetDefaultValues(path, hueMap.Comment)
    if err != nil {
        return fmt.Errorf("error during text-editor set-up: %s", err.Error())
    }
    //edit mode
    err = editMode(&editor, &hueMap)
    if err != nil {
        return fmt.Errorf("error during file-edit session: %s", err.Error())
    }
    //exiting program
    err = saveAndExit(&editor)
    if err != nil {
        return fmt.Errorf("error during save/exit: %s", err.Error())
    }
    return nil
}

func main() {
    err := runTextEditor()
    if err != nil {
        fmt.Printf("%s\nrun with -h flag for more info\n", err.Error())
    }
}
