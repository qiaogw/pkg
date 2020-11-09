package filemanager

import "os"

type SortByFileType []os.FileInfo

func (s SortByFileType) Len() int { return len(s) }
func (s SortByFileType) Less(i, j int) bool {
	if s[i].IsDir() {
		if !s[j].IsDir() {
			return true
		}
	} else if s[j].IsDir() {
		if !s[i].IsDir() {
			return false
		}
	}
	return s[i].Name() < s[j].Name()
}
func (s SortByFileType) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByModTime []os.FileInfo

func (s SortByModTime) Len() int { return len(s) }
func (s SortByModTime) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}
func (s SortByModTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByModTimeDesc []os.FileInfo

func (s SortByModTimeDesc) Len() int { return len(s) }
func (s SortByModTimeDesc) Less(i, j int) bool {
	return s[i].ModTime().UnixNano() > s[j].ModTime().UnixNano()
}
func (s SortByModTimeDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortByNameDesc []os.FileInfo

func (s SortByNameDesc) Len() int { return len(s) }
func (s SortByNameDesc) Less(i, j int) bool {
	return s[i].Name() > s[j].Name()
}
func (s SortByNameDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
