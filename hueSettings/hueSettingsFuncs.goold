package huesettings

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "slices"
    "strings"
)

const settingsPath string = "/APESettings.cfg"
const settingsFlagHelpMsg string = "\r\nrun with -s flag to edit APESettings"

var validTypes []string = []string{"keyword", "symbols", "strings", "comment"}
var validSymbols string = "!#$%&()*+,-./:;<=>?@[\\]^{|}~" //list of ascii symbols (excluding whitespace, underscores and string symbols)
var validStrings []string = []string{"\"", "'", "`"}

var ValidHues []string = []string{"30", "31", "32", "33", "34", "35", "36", "37", "38", "39"}

func CreateSettingsFileIfNotExists() error {
    settingsPath, err := GetSettingsPath()
    if err != nil {
        return err
    }

    _, err = os.Open(settingsPath)
    if err != nil {
        _, err = os.Create(settingsPath)
        if err != nil {
            return err
        }
    }
    return nil
}

//gets path to settings file in home dir
func GetSettingsPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return homeDir + settingsPath, nil
}

//parses line in settings file, returns each of its parts, returns error if line is malformed
func parseSettingsEntry(line string) (string, string, string, string, error) {
    tokens := strings.Split(line, " ")
    if len(tokens) != 4 {
        return "", "", "", "", fmt.Errorf("malformed entry: '%s'%s", line, settingsFlagHelpMsg)
    }

    ext, valueType, identifier, hue := tokens[0], tokens[1], tokens[2], tokens[3]
    if !slices.Contains(validTypes, valueType) {
        return "", "", "", "", fmt.Errorf("invalid type entry: '%s'%s", line, settingsFlagHelpMsg)
    }
    if valueType == "keyword" && strings.ContainsAny(identifier, validSymbols+"\"'`") {
        return "", "", "", "", fmt.Errorf("invalid keyword entry (contains non-underscore symbols): '%s'%s", line, settingsFlagHelpMsg)
    }
    if valueType == "symbols" && (!strings.ContainsAny(identifier, validSymbols) || len(identifier) != 1) {
        return "", "", "", "", fmt.Errorf("invalid symbols entry: '%s'%s", line, settingsFlagHelpMsg)
    }
    if valueType == "strings" && !slices.Contains(validStrings, identifier) {
        return "", "", "", "", fmt.Errorf("invalid strings entry (must be one of: \", ', `): '%s'%s", line, settingsFlagHelpMsg)
    }
    if !slices.Contains(ValidHues, hue) {
        return "", "", "", "", fmt.Errorf("invalid hue entry: '%s'%s", line, settingsFlagHelpMsg)
    }
    hue = "\x1b[" + hue + "m"
    return ext, valueType, identifier, hue, nil
}

//reads and parses syntax hue settings, returns settings as a HueMap for a chosen file ext
func GetSyntaxHueSettings(chosenExt string) (HueMap, error) {
    settingsPath, err := GetSettingsPath()
    if err != nil {
        return HueMap{}, err
    }
    iofile, err := os.Open(settingsPath)
    if err != nil {
        return HueMap{}, err
    }
    hueMap := HueMap{
        Keywords: map[string]string{},
        Symbols: map[string]string{},
        Strings: map[string]string{},
    }

    //reading and parsing file
    reader := bufio.NewReader(iofile)
    for {
        line, readErr := reader.ReadString(byte('\n'))
        line = strings.ReplaceAll(line, "\n", "")
        if line != "" {
            ext, valueType, identifier, hue, parseErr := parseSettingsEntry(line)
            if parseErr != nil {
                return HueMap{}, parseErr
            }
            if ext == chosenExt { //only adding entries with the chosen file ext
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
