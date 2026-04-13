package main

import "encoding/json"

type StandardError struct {
	Domain  string `json:"domain"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *StandardError) Error() string {
	return e.Domain + ":" + e.Code + ":" + e.Message
}

func DeserializeStandardError(stderr []byte) (*StandardError, error) {
	var stdErr StandardError
	if err := json.Unmarshal(stderr, &stdErr); err != nil {
		return nil, err
	}
	return &stdErr, nil
}
