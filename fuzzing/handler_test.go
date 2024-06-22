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
