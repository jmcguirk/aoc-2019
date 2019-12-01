package main

import (
	"fmt"
	"log"
	"time"
)

type logWrapper struct {
	StartTime time.Time;
}

func (this *logWrapper) Init(){
	this.StartTime = time.Now();
}

func (this *logWrapper) Info(format string, args ...interface{}){

	prefix := fmt.Sprintf("[INFO] %s - ", this.SecondsSinceStatup());
	suffix := fmt.Sprintf(format, args...);
	fmt.Println(prefix + suffix);
}

func (this *logWrapper) Fatal(format string, args ...interface{}){

	prefix := fmt.Sprintf("[FATAL] %s - ", this.SecondsSinceStatup());
	suffix := fmt.Sprintf(format, args...);
	log.Fatal(prefix + suffix);
}

func (this *logWrapper) FatalError(err error){

	prefix := fmt.Sprintf("[FATAL] %s - ", this.SecondsSinceStatup());
	suffix := fmt.Sprintf(err.Error());
	log.Fatal(prefix + suffix);
}

func (this *logWrapper) SecondsSinceStatup() string{
	return FormatDuration(time.Now().Sub(this.StartTime));
}

var Log = &logWrapper{};
