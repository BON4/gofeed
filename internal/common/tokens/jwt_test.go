package tokens

import (
	"testing"
	"time"
)

type Account struct {
	Username string
	Role     string
}

var SECRET = "qwebnmrtyxcv123hjkasd567op091234"

var generator Generator[Account]

func TestMain(m *testing.M) {
	var err error
	generator, err = NewJWTGenerator[Account](SECRET)

	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestJWT(t *testing.T) {
	accsesDur := time.Minute
	refreshDur := time.Hour

	account := Account{
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
