package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("Returns customized message", func(t *testing.T) {
		expected := "Hi there, I love you!"

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/you", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		actual, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if expected != string(actual) {
			t.Fail()
		}
	})

	t.Run("Returns loaded page content", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/view/TestPage", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200. Actual %d\n", resp.StatusCode)
		}
	})

	t.Run("Redirects on wiki page does not exist", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/view/NewPage", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200. Actual %d\n", resp.StatusCode)
		}

		if resp.Request.URL.Path != "/edit/NewPage" {
			t.Errorf("Expected redirect to /edit/NewPage. Actual %s", resp.Request.URL.Path)
		}
	})
}

func TestHandler(t *testing.T) {
	t.Run("Redirects on wiki page does not exist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/view/NewPage", nil)
		w := httptest.NewRecorder()
		viewHandler(w, req)

		resp := w.Result()

		if resp.StatusCode != 302 {
			t.Errorf("Expected status code 302. Actual %d.", resp.StatusCode)
		}
		if resp.Header.Get("Location") != "/edit/NewPage" {
			t.Errorf("Expected /edit/NewPage. Actual %s.", resp.Header.Get("Location"))
		}
	})
}
