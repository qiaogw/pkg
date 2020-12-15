package tools

//RemoveRepeatedElement 数组中删除重复对象
func RemoveRepeatedElement(arr []interface{}) (newArr []interface{}) {
	newArr = make([]interface{}, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

//ListExists判断元素是否在数组内
func ListExists(arr []string, e string) (ok bool) {
	var set map[string]struct{}
	set = make(map[string]struct{})
	// 上面2部可替换为set := make(map[string]struct{})
	// 将list内容传递进map,只根据key判断，所以不需要关心value的值，用struct{}{}表示
	for _, value := range arr {
		set[value] = struct{}{}
	}
	// 检查元素是否在map
	_, ok = set[e]
	return
}

//ListExists  删除元素是否在数组内
func ListRemove(arr []string, e string) (newArr []string) {
	newArr = make([]string, 0)
	for _, v := range arr {
		if v != e {
			newArr = append(newArr, v)
		}
	}
	return newArr
}
