package main

import (
	"encoding/json"
)

type request1 struct {
	Method *string
	Number json.RawMessage `json:"number"`
}

type response1 struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}
