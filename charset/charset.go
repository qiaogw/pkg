package charset

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/axgle/mahonia"
	//"github.com/admpub/chardet"
	//sc "github.com/admpub/mahonia"
	runewidth "github.com/mattn/go-runewidth"
)

//字符集转换，从fromEnc 到 toEnc
func Convert(fromEnc string, toEnc string, b []byte) ([]byte, error) {
	if !Validate(fromEnc) {
		return nil, errors.New(`Unsuppored charset: ` + fromEnc)
	}
	if !Validate(toEnc) {
		return nil, errors.New(`Unsuppored charset: ` + toEnc)
	}
	dec := mahonia.NewDecoder(fromEnc)
	s := dec.ConvertString(string(b))
	enc := mahonia.NewEncoder(toEnc)
	s = enc.ConvertString(s)
	b = []byte(s)
	return b, nil
}

func Validate(enc string) bool {
	return mahonia.GetCharset(enc) != nil
}

func Truncate(str string, width int) string {
	w := 0
	b := []byte(str)
	var buf bytes.Buffer
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rw := runewidth.RuneWidth(r)
		if w+rw > width {
			break
		}
		buf.WriteRune(r)
		w += rw
		b = b[size:]
	}
	return buf.String()
}

func With(str string) int {
	return runewidth.StringWidth(str)
}

func RuneWith(str string) int {
	return utf8.RuneCountInString(str)
}
