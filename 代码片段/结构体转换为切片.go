func StrctToSlice(f Prodata) []string {
	//把结构体转换为切片
	//切片可以方便的写入
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		ss[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return ss
}
