package logs

import (
	"log"
	"os"
)

var (
	lInfo    = log.New(os.Stdout, "[info] ", log.Lmsgprefix|log.LstdFlags|log.LUTC)
	lWarning = log.New(os.Stderr, "[warning] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lError   = log.New(os.Stderr, "[error] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lFatal   = log.New(os.Stderr, "[fatal] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	lDebug   = log.New(os.Stderr, "[debug] ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
)

var (
	Info    = lInfo.Println
	Warning = lWarning.Println
	Error   = lError.Println
	Fatal   = lFatal.Fatalln
	Debug   = lDebug.Println
)
