package main

import (
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/klog/v2"
)

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.Errorf("Error while proxying request: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func main() {
	bpUrl := "https://api-backplane.apps.siwu-bp.gjia.s1.devshift.org/"
	url, err := url.Parse(bpUrl)
	if err != nil {
		panic(err)
	}

	responder := &responder{}
	proxy := proxy.NewUpgradeAwareHandler(url, nil, false, false, responder)
	proxy.UseRequestLocation = true
	proxy.UseLocationHost = true
	proxy.AppendLocationPath = false

	http.Handle("/", proxy)
	klog.Fatal(http.ListenAndServe(":8090", nil))
}
