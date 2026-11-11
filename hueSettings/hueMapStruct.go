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
func (h *HueMap) insertIntoHueMap(valueType, value, hue string) {
    switch valueType {
    case "keyword":
        h.Keywords[value] = hue
    case "symbols":
        h.Symbols[value] = hue
    case "strings":
        h.Strings[value] = hue
    case "comment":
        h.Comment = value
        h.CommentHue = hue
    }
}
