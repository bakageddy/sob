package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type AppState struct {
	Start     int64 `json:"start"`
	reset     bool
	StatePath string
}

var app AppState

func init() {
	flag.BoolVar(&app.reset, "restart", false, "Restart streak")
	flag.Int64Var(&app.Start, "start_date", -1, "Set Start Date of Streak in UTC Time Format")
	flag.StringVar(&app.StatePath, "path", "$HOME/.local/state/sob.state", "Set Default Path to save history")
}

func save_to_file(app AppState) error {
	app.StatePath = os.ExpandEnv(app.StatePath)
	f, err := os.OpenFile(app.StatePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	err = enc.Encode(app)
	return err
}

func load_from_file() (*AppState, error) {
	app.StatePath = os.ExpandEnv(app.StatePath)
	f, err := os.Open(app.StatePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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
	if app.reset == true {
		if app.Start == -1 {
			app.Start = time.Now().Unix()
		}
		err := save_to_file(app)
		if err != nil {
			fmt.Println("Error trying to save file: ", err)
			os.Exit(1)
		}
	} else {
		app, err := load_from_file()
		if err != nil {
			fmt.Println("Error trying to load file: ", err)
			fmt.Println("Consider restarting your counter")
			os.Exit(1)
		}
		t := time.Unix(app.Start, 0)
		duration := time.Since(t)
		hours := int64(duration.Hours())
		days := hours / 24
		hours = hours % 24

		fmt.Printf("Days: %d Hours: %d\n", days, hours)
	}
}
