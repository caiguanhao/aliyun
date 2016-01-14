package main

import (
	"fmt"
)

type ECSResponseError struct {
	Code      string `json:"Code"`
	HostID    string `json:"HostId"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
}

func (err *ECSResponseError) Error() string {
	return fmt.Sprintf("[FATAL] %s: %s", err.Code, err.Message)
}
