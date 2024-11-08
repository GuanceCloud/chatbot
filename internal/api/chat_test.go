package api

import (
	"testing"

	"github.com/GuanceCloud/chatbot/pkg/utils"
	"gotest.tools/v3/assert"
)

var tk = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWNudF9hZjc4YjY4MjM4NjAxMWViOTIzMWVlMGQ5MDYwZTRiYyIsImV4cCI6MTczMTY2NTI4N30.8udNl1vi3QLduLpM58sGO1H3IsV06Gz1901Eb0ZELyE"

func TestS(t *testing.T) {
	v, err := utils.VerifyToken(tk)
	assert.NilError(t, err)
	assert.Equal(t, "acnt_af78b682386011eb9231ee0d9060e4bc", v["user_id"])
}
