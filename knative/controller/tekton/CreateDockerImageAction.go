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

	_ = r.ParseForm()
	zipFile, _, err := r.FormFile("file")
	if err != nil {
		log.Printf(err.Error())
	}
	defer func() {
		_ = zipFile.Close()
	}()

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Printf(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"" + "/node-hello"},
		Remove:     true,
	}
	res, err := dockerClient.ImageBuild(r.Context(), zipFile, opts)
	if err != nil {
		log.Printf(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	scanner := bufio.NewScanner(res.Body)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	_, _ = w.Write([]byte(lastLine))
}
