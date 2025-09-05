package logger

import (
	"context"
	"encoding/json"
	"log"
	"os"
)

var (
	debugMode = true
	logger    = log.New(os.Stdout, "[APP] ", log.LstdFlags|log.Lshortfile)
)

func Debug(message string, param any) {
	if !debugMode {
		return
	}
	logger.Printf("%s: %+v\n", message, param)
}

func DebugJSON(ctx context.Context, message string, param any) {
	if !debugMode {
		return
	}

	bytes, err := json.MarshalIndent(param, "", "  ")
	if err != nil {
		logger.Printf("DebugJSON error: %v", err)
		return
	}

	logger.Printf("%s:\n%s\n", message, string(bytes))
}
