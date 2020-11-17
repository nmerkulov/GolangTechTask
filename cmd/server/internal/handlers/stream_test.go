package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_streamHandler_CreateStream(t *testing.T) {
	//it is possible to generate mocks by gomock or mockery, but IMHO it is just fine to use inmem storage
	store := NewInMemStore()
	routes := Routes(store)
	ts := httptest.NewServer(routes)
	defer ts.Close()
	validJson := new(bytes.Buffer)
	_ = json.NewEncoder(validJson).Encode(CreateStream{
		Name:    "very name",
		BuffIDs: []uint64{1, 2, 3},
	})
	resp, body := testRequest(t, ts, http.MethodPost, "/stream", validJson)
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("status expected: 200\n got: %d", resp.StatusCode))
	}
	var stream Stream
	if err := json.Unmarshal(body, &stream); err != nil {
		t.Fatal(fmt.Errorf("testRequest#Unmarshal: %w", err))
	}

	if stream.Name != "very name" {
		t.Errorf("invalid name,\nexpected: very_name\n     got: %v", stream.Name)
	}
	s, err := store.GetStream(stream.ID)
	if err != nil {
		t.Errorf("unexpected error in storage: %v", err)
	}
	if !reflect.DeepEqual(s, stream) {
		t.Errorf("storage value and api value differs")
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp, nil
	}
	return resp, respBody
}
