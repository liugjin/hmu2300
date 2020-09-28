package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/gwaylib/errors"
	"github.com/howeyc/fsnotify"
)

func WatchFile(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(errors.As(err))
		return
	}
	defer watcher.Close()
	if err := watcher.Watch(file); err != nil {
		log.Fatal(errors.As(err))
		return
	}

	end := make(chan os.Signal, 2)
	go func() {
		for {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.Println("error:", errors.As(err))
				end <- os.Interrupt
				return
			}
			fmt.Printf("now file data:\n%s\n", string(data))
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", errors.As(err))
			}
		}
	}()

	signal.Notify(end, os.Interrupt, os.Kill)
	<-end
}

func main() {
	file := flag.String("file", "", "listen file")

	flag.Parse()
	if len(*file) == 0 {
		flag.Usage()
		return
	}

	WatchFile(*file)
}
