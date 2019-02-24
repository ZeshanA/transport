package stringhelper

func SliceToInterface(slice *[]string) []interface{} {
	interfaceSlice := make([]interface{}, len(*slice))
	for i, v := range *slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
