package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
)

var (
	apiSecret   string
	outcomesVar string
	client      *stream.Client
)

func init() {
	rand.Seed(time.Now().Unix())
	flag.StringVar(&apiSecret, "secret", "", "Stream Chat API Secret")

	var outcomeNames []string
	for name, _ := range availableOutcomes {
		outcomeNames = append(outcomeNames, name)
	}

	defaultOutcomes := strings.Join(outcomeNames, ",")
	outcomeUsage := fmt.Sprintf("List of outcomes (commas) to use, allowed values: %s", defaultOutcomes)
	flag.StringVar(&outcomesVar, "outcomes", defaultOutcomes, outcomeUsage)
}

func writeToken(w http.ResponseWriter, token string) {
	payload := map[string]string{
		"token": token,
	}

	json.NewEncoder(w).Encode(payload)
}

func returnsGarbageToken(w http.ResponseWriter, _ string) {
	writeToken(w, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ6IjoieCJ9.aXtogGVI9k20geqrwSgwo9eHHHN1CRI6UlA4WXjZJPw")
}

func returnsValidToken(w http.ResponseWriter, userID string) {
	token, err := client.CreateToken(userID, time.Now().Add(time.Minute))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	writeToken(w, token)
}

func returnsExpiredToken(w http.ResponseWriter, userID string) {
	token, err := client.CreateToken(userID, time.Now().Add(-time.Minute))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	writeToken(w, token)
}

func serverErrors(w http.ResponseWriter, _ string) {
	w.WriteHeader(500)
}

func serverTimesOut(w http.ResponseWriter, _ string) {
	time.Sleep(time.Hour)
}

type tokenGen func(http.ResponseWriter, string)

var outcomes []tokenGen

var availableOutcomes = map[string]tokenGen{
	"serverTimesOut":      serverTimesOut,
	"serverErrors":        serverErrors,
	"returnsExpiredToken": returnsExpiredToken,
	"returnsValidToken":   returnsValidToken,
	"returnsGarbageToken": returnsGarbageToken,
}

func token(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.URL.Query()["userID"]

	if !ok || len(userID[0]) < 1 {
		log.Println("URL Param 'userID' is missing")
		w.WriteHeader(500)
		return
	}

	fn := outcomes[rand.Intn(len(outcomes))]
	log.Println(fmt.Sprintf("running outcome %v", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()))
	fn(w, userID[0])
}

func main() {
	var err error

	flag.Parse()

	if apiSecret == "" {
		log.Fatal("api secret param not provided")
	}

	for _, outcome := range strings.Split(outcomesVar, ",") {
		fn, ok := availableOutcomes[outcome]
		if ok {
			outcomes = append(outcomes, fn)
		}
	}

	if len(outcomes) == 0 {
		log.Fatal("did not provide a list of valid outcomes see --help")
	}

	client, err = stream.NewClient("key-does-matter", apiSecret)
	if err != nil {
		print(err)
		return
	}
	http.HandleFunc("/token", token)
	log.Println("listening on :8090")
	log.Println("selected " + outcomesVar + " outcomes")
	http.ListenAndServe(":8090", nil)
}
