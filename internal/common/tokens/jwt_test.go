package tokens

import (
	"testing"
	"time"
)

var SECRET = "qwebnmrtyxcv123hjkasd567op091234"

var generator Generator[InstanceCredentials]

func TestMain(m *testing.M) {
	var err error
	generator, err = NewJWTGenerator[InstanceCredentials](SECRET)

	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestJWT(t *testing.T) {
	accsesDur := time.Minute
	refreshDur := time.Hour

	account := InstanceCredentials{
		Username: "admin",
		Role:     "admin",
	}

	acessToken, acessPayload, err := generator.CreateToken(account, accsesDur)
	if err != nil {
		t.Fatal(err)
		return
	}

	refreshToken, refreshPayload, err := generator.CreateToken(account, refreshDur)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("access tok: %+v\npayload: %+v", acessToken, acessPayload)
	t.Logf("refresh tok: %+v\npayload: %+v", refreshToken, refreshPayload)
}
