package tekton

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"net/http"
)

type CreateDockerImageAction struct {
	http.Handler
}

func (a *CreateDockerImageAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Print("here1")

	_ = r.ParseForm()
	zipFile, _, err := r.FormFile("file")
	if err != nil {
		log.Printf(err.Error())
		fmt.Print(err.Error())
	}
	defer func() {
		_ = zipFile.Close()
	}()

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Printf(err.Error())
		fmt.Print(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile_admissionwebhook",
		Tags:       []string{"" + "/node-hello"},
		Remove:     true,
	}
	fmt.Print("here2")

	res, err := dockerClient.ImageBuild(r.Context(), zipFile, opts)
	if err != nil {
		log.Printf(err.Error())
		fmt.Print(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fmt.Print("here3")

	scanner := bufio.NewScanner(res.Body)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}
	fmt.Print("here4")

	_, _ = w.Write([]byte(lastLine))
}
