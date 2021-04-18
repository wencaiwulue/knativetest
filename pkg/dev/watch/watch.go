package watch

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/cp"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	sync2 "knativetest/pkg/dev/sync"
	"knativetest/pkg/dev/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var watcher *fsnotify.Watcher
var doneChan = make(chan bool, 1)
var localDir, remoteDir string

func Start() {
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("event: %s, rel path: %s\n", event, getRPath(event.Name, localDir))
				rpath := getRPath(event.Name, localDir)
				newpath := filepath.Join(remoteDir, rpath)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					syncWithExec(event.Name, newpath)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					removeWithExec(event.Name, newpath)
					fi, _ := os.Stat(event.Name)
					if fi.IsDir() {
						unWatch(fi.Name())
					}
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					syncWithExec(event.Name, newpath)
					fi, _ := os.Stat(event.Name)
					if fi.IsDir() {
						watch(fi.Name())
					}
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					syncWithExec(event.Name, newpath)
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					syncWithExec(event.Name, newpath)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	<-doneChan
	log.Printf("done chan receive stop singe")
}

func getRPath(f, base string) string {
	base = filepath.Join(base, string(os.PathSeparator))
	return strings.ReplaceAll(f, base, "")
}

func Watch(local, remote string) {
	localDir = local
	remoteDir = remote
	watch(localDir)
	syncWithExec(localDir, remote)
	Start()
}

func unWatch(path string) {
	if e := watcher.Remove(path); e != nil {
		fmt.Println(e)
	}
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

func syncWithExec(local, remote string) {
	fmt.Printf("try to sync local: %s, to remote: %s\n", local, remote)
	pods, err := util.Clients.ClientSet.CoreV1().Pods("test").
		List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
	if err != nil {
		log.Fatalf("get kubedev pod error: %v", err)
	}
	if len(pods.Items) != 1 {
		log.Println("this should not happened")
	}
	kubeconfig := util.Clients.Kubeconfig
	namespace := "test"
	podSyncer := sync2.NewPodSyncer(namespace, kubeconfig, "")
	maps := map[string][]string{local: {remote}}
	copyFileFn := podSyncer.CopyFileFn(context.Background(), pods.Items[0], pods.Items[0].Spec.Containers[0], maps)
	if b, err := sync2.RunCmdOut(copyFileFn); err != nil {
		fmt.Printf("error sync: %v, log: %s\n", err.Error(), string(b))
	}
}

func removeWithExec(local, remote string) {
	fmt.Printf("try to sync local: %s, to remote: %s\n", local, remote)
	pods, err := util.Clients.ClientSet.CoreV1().Pods("test").
		List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
	if err != nil {
		log.Fatalf("get kubedev pod error: %v", err)
	}
	if len(pods.Items) != 1 {
		log.Println("this should not happened")
	}
	kubeconfig := util.Clients.Kubeconfig
	namespace := "test"
	podSyncer := sync2.NewPodSyncer(namespace, kubeconfig, "")
	maps := map[string][]string{local: {remote}}
	copyFileFn := podSyncer.DeleteFileFn(context.Background(), pods.Items[0], pods.Items[0].Spec.Containers[0], maps)
	if b, err := sync2.RunCmdOut(copyFileFn); err != nil {
		fmt.Printf("error sync: %v, log: %s\n", err.Error(), string(b))
	}
}

func init() {
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	watcher = newWatcher
}

func copy(local, remote string) {
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
	err = opts.Complete(newFactory(), &cobra.Command{})
	if err != nil {
		log.Printf("complete error: %v\n", err)
	} else {
		log.Println("complete no error")
	}
	opts.Namespace = "test"
	opts.Container = "test"
	err = opts.Run([]string{"/Users/naison/Desktop", "test/test-54d97cbcd-792r4:/tmp"})
	if err != nil {
		log.Printf("error info: %v\n", err)
	} else {
		log.Println("sync no error")
	}
}

func newFactory() cmdutil.Factory {
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	path := "/Users/naison/codingtest"
	namespace := "test"
	kubeConfigFlags.Namespace = &namespace
	kubeConfigFlags.KubeConfig = &path
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	return cmdutil.NewFactory(matchVersionKubeConfigFlags)
}
