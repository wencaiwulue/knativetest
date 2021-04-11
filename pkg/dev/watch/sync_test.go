package watch

import (
	"github.com/fsnotify/fsnotify"
	"k8s.io/utils/exec"
	"log"
	"testing"
)

func TestSync(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				nginx()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/path/to/file1")
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add("/path/to/file2") //也可以监听文件夹
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
func nginx() {
	cmd := exec.New().Command("/usr/local/bin/lunchy", "restart", "sync")
	cmd.Run()
}
