package test

import (
	"encoding/json"
	"fmt"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

func parseFile(filename string, out interface{}) {
	b, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("read file error: %v\n", e)
	}
	var bb, _ = yaml.ToJSON(b)
	err := json.Unmarshal(bb, out)
	fmt.Printf("info: %v, error: %v", out, err)
}

func Test(t *testing.T) {
	task := &v1beta1.ClusterTask{}
	parseFile("yaml/task.yaml", task)
	fmt.Printf("task info: %v", task)
}
