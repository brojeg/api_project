package models

import (
	"net/http"
	"testing"
)

func TestRespond(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		data Response
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Respond(tt.args.w, tt.args.data)
		})
	}
}
