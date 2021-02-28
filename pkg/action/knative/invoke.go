package knative

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"test/pkg/action"
	"time"
)

type InvokeAction struct {
	action.Action
	Namespace string
	Name      string
	Tag       string
}

// TODO try to use one of four event trigger mode
func (c *InvokeAction) Process(ctx context.Context) interface{} {
	url := "istio-ingressgateway.istio-system.svc.cluster.local"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(""))
	if err != nil {
		return nil
	}
	m := ctx.Value("http").(map[string]interface{})
	reqq := m["http.request"].(*http.Request)
	body := m["http.request.body"].(string)
	fmt.Printf("body: %v\n", body)
	buf := new(bytes.Buffer)
	_ = reqq.Header.Write(buf)
	_ = req.Header.WriteSubset(buf, map[string]bool{
		"Host":              true,
		"Content-Length":    true,
		"Transfer-Encoding": true,
	})
	req.Header.Set("RequestId", uuid.New().String())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Knative-Serving-Tag", c.Tag)
	req.Header.Set("Host", fmt.Sprintf("%s.%s.xip.io", c.Name, c.Namespace))

	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(startTime)
	if err != nil {
		return nil
	}
	fmt.Printf("time cost: %v\n", latency.String())

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return respBody
}
