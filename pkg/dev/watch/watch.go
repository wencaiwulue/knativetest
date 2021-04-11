package watch

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/exec"
	"knativetest/pkg/dev/cp"
	"knativetest/pkg/dev/util"
	"log"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher
var doneChan = make(chan bool, 1)

func Start() {
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					sync(event.Name, "/tmp"+event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					sync(event.Name, "/tmp"+event.Name)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					sync(event.Name, "/tmp"+event.Name)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					sync(event.Name, "/tmp"+event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	<-doneChan
	log.Printf("done chan receive stop singe")
}

func Watch(path string) {
	watch(path)
	Start()
}

func watch(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return
	}

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	// if directory have subdirectory, then needs to watch all subdirectory
	if fileInfo.IsDir() {
		list, _ := ioutil.ReadDir(path)
		for _, info := range list {
			if info.IsDir() {
				watch(filepath.Join(path, info.Name()))
			}
		}
	}
}

func sync(local, remote string) {
	//cmd := exec.New().Command("/usr/local/bin/lunchy", "restart", "sync")
	//cmd.Run()
	log.Printf("try to sychronize file, local: %v, remote: %v\n", local, remote)
	pods, err := util.Clients.ClientSet.CoreV1().Pods("test").List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
	if err != nil {
		log.Fatalf("get kubedev pod error: %v", err)
	}
	if len(pods.Items) != 1 {
		log.Println("this should not happened")
	}
	ioStreams, _, _, _ := genericclioptions.NewTestIOStreams()
	cmd := cp.NewCmdCp(ioStreams)
	opts := cp.NewCopyOptions(ioStreams)
	src := cp.FileSpec{
		File: local,
	}
	dest := cp.FileSpec{
		Namespace: "test",
		PodName:   pods.Items[0].Name,
		File:      remote,
	}
	_ = opts.Complete(cmd)
	options := &exec.ExecOptions{}
	if err = opts.CopyToPod(src, dest, options); err != nil {
		// todo bugs here
		// /Users/naison/go/pkg/mod/k8s.io/client-go@v0.18.8/tools/remotecommand/remotecommand.go:108
		fmt.Printf("copy to pod error: %v\n", err.Error())
	}

}

func init() {
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	watcher = newWatcher
}
