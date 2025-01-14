package logger

import (
	"fmt"
)

func Log(msg string) {
	fmt.Println(msg)
}

func Error(msg any) {
	fmt.Println(fmt.Sprintf("ERROR: %v", msg))
}