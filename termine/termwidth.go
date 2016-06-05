// +build !linux

package main

func TermWidth() (width int, ok bool) {
	return 0, false
}
