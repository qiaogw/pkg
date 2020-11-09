package hugo

import (
	"regexp"
	"testing"
)

func TestGetVaild(t *testing.T) {

	str := `{ 
	"title": "Bright",
		"date": "2020-08-05T18:45:40.379323+08:00",
		"author": "John Doe",
		"image": "images/blog/blog-post-1.jpg",
		"bg_image": "images/featue-bg.jpg",
		"categories": [
		"行业新闻"
],
"tags": [
"市场",
"政策"
],
"description": "this is meta description",
"draft": false,
"type": "post"
}


zxtestcv，Vaild[Email;Mi
`
	//str = "zxtestcv，Vaild[Email;Min(120)]()]  {  aaa ,{eee, {[ddd,fff,]} } }"
	var reg = regexp.MustCompile(`(?s:\{(.*?)})`)
	params := reg.FindStringSubmatch(str)
	for _, param := range params {
		t.Log(param)
	}
}
