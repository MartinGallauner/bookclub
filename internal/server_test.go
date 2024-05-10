package internal

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_jwtMiddleware(t *testing.T) {

	type args struct {
		next http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO:: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jwtMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jwtMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
