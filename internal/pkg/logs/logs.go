package logs

import (
	"log"
	"os"
)

var (
	lInfo    = log.New(os.Stdout, "[INFO] ", log.Lmsgprefix|log.LstdFlags|log.LUTC)
	lWarning = log.New(os.Stderr, "[WARNING] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lError   = log.New(os.Stderr, "[ERROR] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lFatal   = log.New(os.Stderr, "[FATAL] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lDebug   = log.New(os.Stderr, "[DEBUG] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
)

var (
	Info    = lInfo.Println
	Warning = lWarning.Println
	Error   = lError.Println
	Fatal   = lFatal.Fatalln
	Debug   = lDebug.Println
)
