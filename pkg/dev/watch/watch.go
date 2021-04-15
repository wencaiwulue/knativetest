package watch

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/cp"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
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
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
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
	pods, err := util.Clients.ClientSet.CoreV1().Pods("test").
		List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
	if err != nil {
		log.Fatalf("get kubedev pod error: %v", err)
	}
	if len(pods.Items) != 1 {
		log.Println("this should not happened")
	}
	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	opts := cp.NewCopyOptions(ioStreams)
	//util.Clients.RestConfig.GroupVersion = &schema.GroupVersion{
	//    Group:   "apps",
	//    Version: "v1",
	//}

	/*ok, logs := util.WaitForCommandDone("kubectl cp /Users/naison/Desktop test/test-54d97cbcd-792r4:/tmp -n test --kubeconfig=/Users/naison/codingtest")
	  if !ok {
	      log.Println("logs: " + logs)
	  } else {
	      log.Println("sync ok")
	  }*/

	//util.Clients.RestConfig.NegotiatedSerializer = scheme.Codecs
	//opts.ClientConfig = util.Clients.RestConfig
	//opts.Clientset = util.Clients.ClientSet
	//opts.Namespace = "test"
	//opts.Container = "test"
	//opts.NoPreserve = true
	//r := "test" + "/" + pods.Items[0].Name + ":" + "/tmp/test.yaml"
	//log.Printf("try to sychronize file, local: %v, remote: %v\n", local, r)
	err = opts.Complete(newFactory(), &cobra.Command{})
	if err != nil {
		log.Printf("complete error: %v\n", err)
	} else {
		log.Println("complete no error")
	}
	opts.Namespace = "test"
	opts.Container = "test"
	err = opts.Run([]string{" /Users/naison/Desktop", "test/test-54d97cbcd-792r4:/tmp"})
	if err != nil {
		log.Printf("error info: %v\n", err)
	} else {
		log.Println("sync no error")
	}
}

func init() {
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	watcher = newWatcher
}

func newFactory() cmdutil.Factory {
	path := "/Users/naison/codingtest"
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	test := "test"
	kubeConfigFlags.Namespace = &test

	kubeConfigFlags.KubeConfig = &path
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	return f
}
