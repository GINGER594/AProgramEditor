package huesettings

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "slices"
    "strings"
)

const settingsPath string = "/NCTSettings.cfg"

var validTypes []string = []string{"keyword", "symbols", "strings", "comment", "mlincom"}
//list of ascii symbols (excluding whitespace, underscores and string symbols
var validSymbols string = "!#$%&()*+,-./:;<=>?@[\\]^_{|}~"
var validStrings []string = []string{"\"", "'", "`"}
 
var ValidHues []string = []string{"30", "31", "32", "33", "34", "35", "36", "37", "38", "39"}

//creates an internal settings file if it does not already exist
func CreateInternalSettingsFile() error {
    internalSettingsPath, err := GetInternalSettingsPath()
    if err != nil {
        return err
    }

    _, err = os.Open(internalSettingsPath)
    if err != nil {
        _, err = os.Create(internalSettingsPath)
        if err != nil {
            return err
        }
    }
    return nil
}

//creates a path to the settings file in the home directory
func GetInternalSettingsPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return homeDir + settingsPath, nil
}

//parses an individual line in the settings file & returns each of its parts. returns an error if the line is malformed
func parseSettingsEntry(line string) (string, string, string, string, error) {
    tokens := strings.Split(line, "|")
    if len(tokens) != 4 {
        return "", "", "", "", fmt.Errorf("malformed entry: '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }

    ext, valueType, identifier, hue := tokens[0], tokens[1], tokens[2], tokens[3]
    if !slices.Contains(validTypes, valueType) {
        return "", "", "", "", fmt.Errorf("invalid type entry: '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    if valueType == "keyword" && strings.ContainsAny(identifier, validSymbols+" \"'`") {
        return "", "", "", "", fmt.Errorf("invalid keyword entry (contains non-underscore symbols): '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    if valueType == "symbols" && (!strings.ContainsAny(identifier, validSymbols) || len(identifier) != 1) {
        return "", "", "", "", fmt.Errorf("invalid symbols entry: '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    if valueType == "strings" && !slices.Contains(validStrings, identifier) {
        return "", "", "", "", fmt.Errorf("invalid strings entry (must be one of: \", ', `): '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    if valueType == "mlincom" && len(strings.Split(identifier, " ")) != 2 {
        return "", "", "", "", fmt.Errorf("invalid multi-line-comment entry (must contain start/end split by ' '): '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    if !slices.Contains(ValidHues, hue) {
        return "", "", "", "", fmt.Errorf("invalid hue entry: '%s'\r\nrun with the '-settings' flag to edit NCTSettings", line)
    }
    hue = "\x1b[" + hue + "m"
    return ext, valueType, identifier, hue, nil
}

//reads & parses syntax hue settings, returns settings as a HueMap, filtered by file extension
func GetSyntaxHueSettings(chosenExt string) (HueMap, error) {
    internalSettingsPath, err := GetInternalSettingsPath()
    if err != nil {
        return HueMap{}, err
    }
    iofile, err := os.Open(internalSettingsPath)
    if err != nil {
        return HueMap{}, err
    }
    hueMap := HueMap{
        Keywords: map[string]string{},
        Symbols: map[string]string{},
        Strings: map[string]string{},
    }

    //reading & parsing the file
    reader := bufio.NewReader(iofile)
    for {
        line, readErr := reader.ReadString(byte('\n'))
        line = strings.ReplaceAll(line, "\n", "")
        if line != "" {
            ext, valueType, identifier, hue, parseErr := parseSettingsEntry(line)
            if parseErr != nil {
                return HueMap{}, parseErr
            }
            if ext == chosenExt { //only adding entries for the desired file type
                hueMap.insertIntoHueMap(valueType, identifier, hue)
            }
        }
        if readErr == io.EOF {
            break
        }
        if readErr != nil {
            return HueMap{}, readErr
        }
    }
    return hueMap, nil
}
