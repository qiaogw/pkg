// Copyright 2018 cloudy itcloudy@qq.com.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package tools

import (
	"crypto/sha256"
	"fmt"
)

// SHA256 为str生成SHA256哈希值
func SHA256(str string) (result string) {
	h := sha256.New()
	h.Write([]byte(str))
	result = fmt.Sprintf("%x", h.Sum(nil))
	return
}
