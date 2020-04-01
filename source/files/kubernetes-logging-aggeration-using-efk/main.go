package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	logpath = "/data/log/log-app/"
)

func setLogFiles() *os.File {
	var logfile string

	if pathstat, err := os.Stat(logpath); os.IsNotExist(err) {
		logfile = "info.log"
	} else {
		if pathstat.IsDir() {
			logfile = logpath + "info.log"
		}
	}

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error(err.Error)
	}
	return f
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	logfile := setLogFiles()
	defer logfile.Close()
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
		NoColors:        true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.Println()
		logrus.Printf("received signal %s", sig)
		done <- true
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			logrus.SetOutput(logfile)
			logrus.Printf("this log to file %s", logfile.Name())
			msg := `this is multiline log
                                line 2
                                line 3`
			logrus.Println(msg)
			logrus.SetOutput(os.Stdout)
			logrus.Println("this log to stdout")
		}
	}()
	logrus.Println("awaiting signal")
	<-done

	logrus.Println("done.")
}
