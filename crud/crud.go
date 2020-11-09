package crud

import (
	"reflect"
)

type TagBody struct {
	Title    string `json:"title"`
	Key      string `json:"key"`
	Sortable bool   `json:"sortable"`
	Label    string `json:"label"`
	Prop     string `json:"prop"`
}

//GetTag orm模型中获取tag
func GetTag(v interface{}, tag string) (ta []TagBody) {
	t := reflect.TypeOf(v).Elem()
	//beego.Debug(t)
	var tb TagBody
	for i := 0; i < t.NumField(); i++ {
		tg := t.Field(i).Tag.Get(tag) //将tag输出出来
		if tg != "" {
			tb.Title = tg
			tb.Key = t.Field(i).Name
			key := t.Field(i).Tag.Get("json")
			if key != "" {
				tb.Key = key
			}
			tb.Sortable = true
			tb.Label = tg
			tb.Prop = tb.Key
			if tb.Key == "Type" {
				tb.Key = "Tname"
				tb.Prop = "Tname"
			}
			if tb.Key == "Pid" {
				tb.Key = "Pname"
				tb.Prop = "Pname"
			}
			ta = append(ta, tb)
		}
	}
	return ta
}

//GetTag kv获取tag
func GetTagSelf(data []map[string]interface{}) (ta []TagBody) {
	var tb TagBody
	for k, _ := range data[0] {
		tb.Title = k
		tb.Label = k
		tb.Prop = k
		tb.Key = k
		ta = append(ta, tb)
	}
	return ta
}
