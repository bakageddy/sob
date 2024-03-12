package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type AppState struct {
	Start int64 `json:"start"`
	reset bool
}

var app AppState

func init() {
	flag.BoolVar(&app.reset, "restart", false, "Restart streak")
	flag.Int64Var(&app.Start, "start_date", -1, "Set Start Date of Streak in UTC Time Format")
}

func save_to_file(app AppState) error {
	f, err := os.OpenFile("./state.json", os.O_WRONLY | os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(app)
	return err
}

func load_from_file() (*AppState, error) {
	f, err := os.Open("./state.json")
	if err != nil {
		return nil, err
	}
	var state AppState
	dec := json.NewDecoder(f)
	err = dec.Decode(&state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func main() {
	flag.Parse()
	if (app.reset != true && app.Start == -1) {
		app, err := load_from_file()
		if err == nil {
			t := time.Unix(int64(app.Start), 0)
			fmt.Println("Streak: ", time.Since(t))
		} else {
			fmt.Println("Error: ", err)
		}
	} else if (app.reset == true || app.Start != -1) {
		app.Start = time.Now().Unix()
		err := save_to_file(app)
		if err != nil {
			fmt.Println(err)
		}
	}
}
