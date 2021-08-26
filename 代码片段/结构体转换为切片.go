func StrctToSlice(f Prodata) []string {
	//把结构体转换为切片
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		ss[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return ss
}
