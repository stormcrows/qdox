package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	want := &QueryResponse{
		Query: "wild weekend",
		Results: []Result{
			{Path: "static/Grand Teton National Park.txt", Similarity: "92"},
			{Path: "static/Around the End - Ralph Henry Barbour.txt", Similarity: "40"},
		},
	}

	testResponse(t, "wild weekend", "5", "0.3", want)
}

func TestQueryWithDifferentN(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	want := &QueryResponse{
		Query: "wild weekend",
		Results: []Result{
			{Path: "static/Grand Teton National Park.txt", Similarity: "92"},
		},
	}

	testResponse(t, "wild weekend", "1", "0.3", want)
}

func TestQueryWithDifferentThreshold(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	want := &QueryResponse{
		Query: "wild weekend",
		Results: []Result{
			{Path: "static/Grand Teton National Park.txt", Similarity: "92"},
		},
	}

	testResponse(t, "wild weekend", "5", "0.5", want)
}

func TestParamsErrors(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testParamsError(t, "", "", "")
	// query
	testParamsError(t, "", "5", "0.3")
	// n
	testParamsError(t, "n is zero", "0", "0.3")
	testParamsError(t, "n less than zero", "-1", "0.3")
	// threshold
	testParamsError(t, "negative threshold", "5", "-0.1")
}

func TestParamsOK(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testParamsOK(t, "this is fine", "1", "0.0")
	testParamsOK(t, "default all", "", "")
	testParamsOK(t, "default threshold", "1", "")
	testParamsOK(t, "default n", "", "0.3")
	testParamsOK(t, "default n and zero threshold", "", "0.0")
}

func testParamsError(t *testing.T, q string, n string, threshold string) {
	param := make(url.Values)
	param.Set("q", q)
	param.Set("n", n)
	param.Set("threshold", threshold)

	req, err := http.NewRequest("GET", "/query?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(QueryHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "incorrect status code")
}

func testParamsOK(t *testing.T, q string, n string, threshold string) {
	param := make(url.Values)
	param.Set("q", q)
	param.Set("n", n)
	param.Set("threshold", threshold)

	req, err := http.NewRequest("GET", "/query?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(QueryHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "incorrect status code")
}

func testResponse(t *testing.T, q string, n string, threshold string, want *QueryResponse) {
	param := make(url.Values)
	param.Set("q", q)
	param.Set("n", n)
	param.Set("threshold", threshold)

	req, err := http.NewRequest("GET", "/query?"+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(QueryHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "incorrect status code")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "incorrect Content-Type")

	resp := rr.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	qresp := QueryResponse{Query: "", Results: make([]Result, 0)}

	err = json.Unmarshal(body, &qresp)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, *want, qresp, "query response different from expected")
}
