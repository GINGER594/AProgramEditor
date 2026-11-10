package huesettings

//struct for all syntax hue data for one file ext
type HueMap struct {
    Keywords map[string]string
    Symbols map[string]string
    Strings map[string]string
    Comment string
    CommentHue string
}

//inserts parsed settings entry into hueMap
func (h *HueMap) insertIntoHueMap(valueType, identifier, hue string) {
    switch valueType {
    case "keyword":
        h.Keywords[identifier] = hue
    case "symbols":
        h.Symbols[identifier] = hue
    case "strings":
        h.Strings[identifier] = hue
    case "comment":
        h.Comment = identifier
        h.CommentHue = hue
    }
}
