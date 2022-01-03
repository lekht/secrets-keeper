package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func handleTestRequest(w *httptest.ResponseRecorder, r *http.Request) {
	router := getRouter()
	router.ServeHTTP(w, r)
}

func TestIndexPage(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("index page is not 200", w.Code)
	}
}

func TestSaveMessage(t *testing.T) {
	testMessage := "foo"
	postData := strings.NewReader(fmt.Sprintf("message=%s", testMessage))
	request, _ := http.NewRequest("POST", "/", postData)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != http.StatusOK {
		t.Error("save is not 200", w.Code)
	}

	key := keyBuilder.Get()
	savedMessage, _ := keeper.Get(key)
	if savedMessage != testMessage {
		t.Error("message was not saved")
	}

	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)
	if !strings.Contains(string(data), key) {
		t.Error("result page without key")
	}
}

func TestReadMessage(t *testing.T) {
	testMessage := "AHAHAhhahaHHAHAH"
	key := keyBuilder.Get()
	keeper.Set(key, testMessage)
	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 200 {
		t.Error("read is not ok", w.Code)
	}
	result := w.Result()
	defer result.Body.Close()
	data, _ := ioutil.ReadAll(result.Body)
	if !strings.Contains(string(data), testMessage) {
		t.Error("result page without key")
	}

	_, err := keeper.Get(key)
	if err == nil {
		t.Error("keeper value must be empty")
	}
}

func TestReadMessageNotFound(t *testing.T) {
	key := keyBuilder.Get()
	request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", key), nil)
	w := httptest.NewRecorder()
	handleTestRequest(w, request)
	if w.Code != 404 {
		t.Error("empty message must be 404", w.Code)
	}
}
