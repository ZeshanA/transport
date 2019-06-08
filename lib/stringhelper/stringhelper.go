package stringhelper

func SliceToInterface(slice *[]string) []interface{} {
	interfaceSlice := make([]interface{}, len(*slice))
	for i, v := range *slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

// SliceToSet converts a slice of strings into a "Set"
// of strings (implemented using a map). This is primarily
// to speed up checks like "does string x exist in the list".
func SliceToSet(slice []string) map[string]bool {
	set := map[string]bool{}
	for _, key := range slice {
		set[key] = true
	}
	return set
}

func AllAfter(after string, slice []string) []string {
	var result []string
	i := IndexOf(after, slice)
	return result[i+1:]
}

func IndexOf(str string, slice []string) int {
	for i, s := range slice {
		if s == str {
			return i
		}
	}
	return -1
}
