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
