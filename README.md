# golang-fuzz-test

This repository is demo how to use golang to do fuzz test

## http handler sample

a http handler with calculate the highest value and return it

```golang
package fuzzing

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValuesRequest struct {
	Values []int `json:"values"`
}

func CalculateHightHandler(w http.ResponseWriter, r *http.Request) {
	var vr ValuesRequest

	if err := json.NewDecoder(r.Body).Decode(&vr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var high int

	for _, value := range vr.Values {
		if value > high {
			high = value
		}
	}

	// if high == 50 {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Something went wrong"))
	// 	return
	// }

	fmt.Fprintf(w, "%d", high)
}
```

## setup fuzz test

a fuzz test which will generate test data from fuzz seed

```golang
package fuzzing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FuzzCalculateHighesHandler(f *testing.F) {

	srv := httptest.NewServer(http.HandlerFunc(CalculateHightHandler))
	defer srv.Close()

	testCases := []ValuesRequest{
		{[]int{1, 2, 3, 4, 5, 6, 6, 7, 7, 8, 9, 10}},
		{[]int{-1, -2, -3, -4, -5, -6, -6, -7, -7, -8, -9, -10}},
		{[]int{-50, -2, -3, 4, -5, 6, 6, 7, 7, 8, 9, 10}},
		{[]int{10, 20, 30, 40, 50, 60, 60, 70, 70, 80, 90, 100}},
	}

	for _, testCase := range testCases {
		data, _ := json.Marshal(testCase)
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if !json.Valid(data) {
			t.Skip("Invalid json")
		}

		vr := ValuesRequest{}
		err := json.Unmarshal(data, &vr)
		if err != nil {
			t.Skip("Invalid json data")
		}
		resp, err := http.DefaultClient.Post(srv.URL, "application/json", bytes.NewBuffer(data))

		if err != nil {
			t.Errorf("Error reaching http API: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		var response int
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Errorf("Error: %v", err)
		}
	})
}

```