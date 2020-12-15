package tools

import (
	"reflect"
	"testing"
)

func TestSubstr(t *testing.T) {
	type args struct {
		str    string
		start  int
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"iamlinet", args{"iamlinet", 3, 2}, "li"},
		{"heisyeelle", args{"heisyeelle", 4, 3}, "yee"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substr(tt.args.str, tt.args.start, tt.args.length); got != tt.want {
				t.Errorf("Substr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstr2(t *testing.T) {
	type args struct {
		str   string
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"iamlinet", args{"iamlinet", 0, 2}, "ia"},
		{"heisyeelle", args{"heisyeelleddddddd", 4, 7}, "yee"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substr2(tt.args.str, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("Substr2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBetweenStr(t *testing.T) {
	type args struct {
		str    string
		substr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"iamlinet", args{"iamlinet", "ml"}, "mlinet,"},
		{"heisyee", args{"heisyee", "y"}, "yee"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBetweenStr(tt.args.str, tt.args.substr); got != tt.want {
				t.Errorf("GetBetweenStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt64(t *testing.T) {
	type args struct {
		min int64
		max int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{"1", args{1, 3}, 2},
		{"1000", args{1001, 1003}, 1002},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt64(tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("RandInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrim(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"i am tom", args{"i am tom"}, "iamtom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Strim(tt.args.str); got != tt.want {
				t.Errorf("Strim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnicode(t *testing.T) {
	type args struct {
		rs string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"iamg", args{"iamg"}, "iamg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unicode(tt.args.rs); got != tt.want {
				t.Errorf("Unicode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicatesAndEmpty(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		// TODO: Add test cases.
		{"12334", args{[]string{"1", "2", "3", "3", "4"}}, []string{"1", "2", "3", "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := RemoveDuplicatesAndEmpty(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("RemoveDuplicatesAndEmpty() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestSliToStr(t *testing.T) {
	type args struct {
		sl     string
		params []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		// TODO: Add test cases.
		{"12334", args{"name", []interface{}{"1", "2", "3", "3", "4"}}, "12334"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := SliToStr(tt.args.sl, tt.args.params...); gotStr != tt.wantStr {
				t.Errorf("SliToStr() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestGetTree(t *testing.T) {
	//pt := make(map[string]interface{})
	ps := []map[string]interface{}{
		{
			"Id":   1,
			"name": "1k",
			"Pid":  0,
		},
		{
			"Id":   2,
			"name": "2k",
			"Pid":  0,
		},
		{
			"Id":   3,
			"name": "3k",
			"Pid":  1,
		},
		{
			"Id":   4,
			"name": "4k",
			"Pid":  1,
		},
		{
			"Id":   5,
			"name": "5k",
			"Pid":  2,
		},
		{
			"Id":   6,
			"name": "6k",
			"Pid":  2,
		},
		{
			"Id":   7,
			"name": "7k",
			"Pid":  6,
		},
	}
	wt := map[string]interface{}{
		"children": []map[string]interface{}{
			{
				"Id":   1,
				"Pid":  0,
				"name": "1k",
				"children": []map[string]interface{}{
					{
						"Id":   3,
						"Pid":  1,
						"name": "3k",
					},
					{
						"Id":   4,
						"Pid":  1,
						"name": "4k",
					},
				},
			},
			{
				"Id":   2,
				"Pid":  0,
				"name": "2k",
				"children": []map[string]interface{}{
					{
						"Id":   5,
						"Pid":  2,
						"name": "5k",
					},
					{
						"Id":   6,
						"Pid":  2,
						"name": "6k",
						"children": []map[string]interface{}{
							{
								"Id":   7,
								"Pid":  6,
								"name": "7k",
							},
						},
					},
				},
			},
		},
	}

	type args struct {
		access []map[string]interface{}
		pid    string
		//pTree  *map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
		{"testint", args{ps, "0"}, wt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTree(tt.args.access, tt.args.pid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTree() = %v, want %v", got, tt.want)
			}

		})
	}
	//fmt.Println(pt)
}
