package integration

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServerContentLengthResponses(t *testing.T) {
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/about", "Identity")
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/about", "identity")
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/about", "gzip")
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/about", "Gzip")
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/mse6/get", "gzip")
	MethodHasContentLengthAndBody(t, "GET", "http://localhost:8080/mse6/get", "identity")

	MethodHasZeroContentLengthAndNoBody(t, "OPTIONS", "http://localhost:8080/mse6/options", "identity")
	MethodHasZeroContentLengthAndNoBody(t, "OPTIONS", "http://localhost:8080/mse6/options", "gzip")
	//golang removes content-length from http 204 response.
	MethodHasNoContentLengthHeaderAndNoBody(t, "OPTIONS", "http://localhost:8080/mse6/options?code=204", "identity")
	MethodHasNoContentLengthHeaderAndNoBody(t, "OPTIONS", "http://localhost:8080/mse6/options?code=204", "gzip")
}

func MethodHasContentLengthAndBody(t *testing.T, method string, url string, acceptEncoding string) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept-Encoding", acceptEncoding)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error connecting to server for method %s, url %s, cause: %s", method, url, err)
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	got := int(resp.ContentLength)
	notwant := -1
	if got == notwant {
		t.Errorf("illegal response for method %s, url %s received Content-Length got %d", method, url, got)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if got != len(bodyBytes) {
		t.Errorf("illegal response for method %s, url %s content-length %d does not match body size %d", method, url, got, len(bodyBytes))
	}
}

func MethodHasZeroContentLengthAndNoBody(t *testing.T, method, url string, acceptEncoding string) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept-Encoding", acceptEncoding)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error connecting to server for method %s, url %s, cause: %s", method, url, err)
	}

	if resp != nil && resp.Body != nil {
		bod, _ := ioutil.ReadAll(resp.Body)
		if len(bod) > 0 {
			t.Errorf("illegal response contains body for method %s, got: %v", method, bod)
		}
		defer resp.Body.Close()
	}

	//not that we don't trust you golang, but test the actual http header sent
	got2 := resp.Header.Get("Content-Length")
	want2 := "0"
	if got2 != want2 {
		t.Errorf("illegal response for method %s Content-Length, want %s got %s", method, want2, got2)
	}
}

func MethodHasNoContentLengthHeaderAndNoBody(t *testing.T, method, url string, acceptEncoding string) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept-Encoding", acceptEncoding)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error connecting to server for method %s, url %s, cause: %s", method, url, err)
	}

	if resp != nil && resp.Body != nil {
		bod, _ := ioutil.ReadAll(resp.Body)
		if len(bod) > 0 {
			t.Errorf("illegal response contains body for method %s, got: %v", method, bod)
		}
		defer resp.Body.Close()
	}

	got2 := resp.Header.Get("Content-Length")
	want2 := ""
	if got2 != want2 {
		t.Errorf("illegal response for method %s Content-Length, want %s got %s", method, want2, got2)
	}
}