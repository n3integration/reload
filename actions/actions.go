package actions

import (
	"log"
	"os"
)

var (
	logger        = log.New(os.Stdout, "[reload] ", 0)
	immediate     = false
	colorGreen    = string([]byte{27, 91, 57, 55, 59, 51, 50, 59, 49, 109})
	colorRed      = string([]byte{27, 91, 57, 55, 59, 51, 49, 59, 49, 109})
	colorReset    = string([]byte{27, 91, 48, 109})
	notifications = false
)
