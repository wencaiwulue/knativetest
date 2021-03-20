package testing

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
		fmt.Printf("read file error: %v\n", e.Error())
	}
	fmt.Println(string(b))
	var bb, _ = yaml.ToJSON(b)
	if err := json.Unmarshal(bb, out); err != nil {
		fmt.Printf("info: %v, error: %v", out, err)
	}
}

func TestParseFile(t *testing.T) {
	task := &v1beta1.ClusterTask{}
	parseFile("../../manifest/task.manifest", task)
	fmt.Printf("task info: %v\n", task)
}
