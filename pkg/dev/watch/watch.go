package watch

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sync2 "knativetest/pkg/dev/sync"
	"knativetest/pkg/dev/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileWatcher struct {
	client     *util.ClientSet
	watcher    *fsnotify.Watcher
	podSyncer  sync2.Syncer
	doneChan   chan bool
	namespace  string
	deployment string
	pod        string
	container  string
	localDir   string
	remoteDir  string
}

type DevOptions interface {
	GetNamespace() string
	GetDeployment() string
	GetPod() string
	GetContainer() string
	GetLocalDir() string
	GetRemoteDir() string
}

func (w *FileWatcher) GetKubeContext() string {
	return ""
}

func (w *FileWatcher) GetKubeConfig() string {
	return w.client.Kubeconfig
}

func (w *FileWatcher) GetKubeNamespace() string {
	return w.namespace
}

func (w *FileWatcher) start() {
	defer w.watcher.Close()
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				from := event.Name
				to := filepath.Join(w.remoteDir, getRelativePath(from, w.localDir))
				log.Printf("event: %v, relative path: %s\n", event, to)
				if event.Op&fsnotify.Write == fsnotify.Write {
					go w.syncWithExec(from, to)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					go w.removeWithExec(from, to)
					w.unWatch(from)
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					go w.syncWithExec(from, to)
					w.watch(from)
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					go w.removeWithExec(from, to)
					w.unWatch(from)
				} else if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					go w.syncWithExec(from, to)
				}
			case err := <-w.watcher.Errors:
				log.Printf("error: %v\n", err)
			}
		}
	}()
	<-w.doneChan
	log.Printf("done chan receive stop singe")
}

func getRelativePath(f, base string) string {
	base = filepath.Join(base, string(os.PathSeparator))
	newPath := strings.ReplaceAll(f, base, "")
	folder := filepath.Base(base)
	return folder + newPath
}

func Watch(cli *util.ClientSet, option DevOptions) {
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	w := &FileWatcher{
		client:     cli,
		watcher:    newWatcher,
		doneChan:   make(chan bool, 1),
		namespace:  option.GetNamespace(),
		deployment: option.GetDeployment(),
		pod:        option.GetPod(),
		container:  option.GetContainer(),
		localDir:   option.GetLocalDir(),
		remoteDir:  option.GetRemoteDir(),
	}
	w.podSyncer = sync2.NewPodSyncer(w, w.namespace)

	w.watch(option.GetLocalDir())
	w.WriteFolder(option.GetLocalDir(), option.GetRemoteDir())
	w.start()
}

func (w *FileWatcher) unWatch(path string) {
	if e := w.watcher.Remove(path); e != nil {
		fmt.Println(e)
	}
}

func (w *FileWatcher) watch(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return
	}

	err = w.watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	// if directory have subdirectory, then needs to watch all subdirectory
	if fileInfo.IsDir() {
		list, _ := ioutil.ReadDir(path)
		for _, info := range list {
			if info.IsDir() {
				w.watch(filepath.Join(path, info.Name()))
			}
		}
	}
}

func (w *FileWatcher) syncWithExec(local, remote string) {
	fmt.Printf("try to sync local: %s, to remote: %s\n", local, remote)
	pod, container := w.getPodAndContainer()
	maps := map[string][]string{local: {remote}}
	copyFileFn := w.podSyncer.CopyFileCmd(context.Background(), pod, container, maps)
	if b, err := sync2.RunCmdOut(copyFileFn); err != nil {
		fmt.Printf("error sync: %v, log: %s\n", err.Error(), string(b))
	}
}

func (w *FileWatcher) removeWithExec(local, remote string) {
	fmt.Printf("try to remove local: %s, to remote: %s\n", local, remote)
	pod, container := w.getPodAndContainer()
	maps := map[string][]string{local: {remote}}
	copyFileFn := w.podSyncer.DeleteFileCmd(context.Background(), pod, container, maps)
	if b, err := sync2.RunCmdOut(copyFileFn); err != nil {
		fmt.Printf("error sync: %v, log: %s\n", err.Error(), string(b))
	}
}

func (w *FileWatcher) WriteFolder(local, remote string) {
	lastFolder := filepath.Base(local)
	remote = filepath.Join(remote, lastFolder)
	log.Printf("copy full, local: %s, remote: %s\n", local, remote)
	pod, container := w.getPodAndContainer()
	deleteFileCmd := w.podSyncer.DeleteFileCmd(context.Background(), pod, container, map[string][]string{local: {remote}})
	if b, err := sync2.RunCmdOut(deleteFileCmd); err != nil {
		fmt.Printf("error empty folder: %v, log: %s\n", err.Error(), string(b))
	}
	copyFileFn := w.podSyncer.CopyFolderCmd(context.Background(), pod, container, local, remote)
	if b, err := sync2.RunCmdOut(copyFileFn); err != nil {
		fmt.Printf("error copy folder: %v, log: %s\n", err.Error(), string(b))
	}
}

func (w *FileWatcher) getPodAndContainer() (corev1.Pod, corev1.Container) {
	if w.pod != "" && w.container != "" {
		return corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: w.pod, Namespace: w.namespace}}, corev1.Container{Name: w.container}
	}

	pods, err := w.client.ClientSet.CoreV1().Pods(w.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
	if err != nil {
		log.Fatalf("get kubedev pod error: %v", err)
	}
	if len(pods.Items) != 1 {
		log.Println("this should not happened")
	}
	return pods.Items[0], pods.Items[0].Spec.Containers[0]
}
