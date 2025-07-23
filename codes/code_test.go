package codes_test

import (
	"errors"
	"gatesvr/codes"

	"testing"
)

func TestConvert(t *testing.T) {
	code := codes.Convert(errors.New("rpc error: code = Unknown desc = code error: code = 10 desc = account exists"))

	t.Log(code)
}
