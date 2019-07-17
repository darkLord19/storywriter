package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

//Input struct parses json input recieved to /story
type Input struct {
	Words    []string `json:"words"`
	Template string   `json:"template"`
}

//PopulateTemplate will substitute placeholders(&1,&2,&3...) with values from words list
func (in *Input) PopulateTemplate() (string, error) {

	var response string

	re := regexp.MustCompile("&[0-9]+")
	indices := re.FindAllIndex([]byte(in.Template), -1)

	l := len(indices)
	n := len(in.Words)

	prevEnd := 0
	for i := 0; i < l; i++ {
		x := indices[i][0]
		y := indices[i][1]
		wordIDx, _ := strconv.Atoi(in.Template[x+1 : y])
		//Return error if wordindex in placeholder is out of range of words list
		if wordIDx > n {
			return "", fmt.Errorf("range error: %s", in.Template[x:y])
		}
		response += in.Template[prevEnd:x] + in.Words[wordIDx-1]
		prevEnd = y
	}
	response += in.Template[prevEnd:]

	return response, nil
}

func main() {
	http.HandleFunc("/story", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var in Input
			json.NewDecoder(r.Body).Decode(&in)
			response, err := in.PopulateTemplate()
			if err != nil {
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			json.NewEncoder(w).Encode(response)
		default:
			fmt.Fprintf(w, "Sorry, only POST method is supported.")
		}
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
