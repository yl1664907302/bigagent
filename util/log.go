package util

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[bigagent]", log.Lshortfile|log.Ldate|log.Ltime)
