package domain

import (
	"time"

	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type AccountCredentials struct {
	Username string
	Role     string
}

type Token struct {
	AccessToken    string
	AccessPayload  *tokens.Payload[AccountCredentials]
	RefreshToken   string
	RefreshPayload *tokens.Payload[AccountCredentials]
}

type TokenFactoryConfig struct {
	RefreshTokenDuration time.Duration
	AccessTokenDuration  time.Duration
	TokenSecret          string
}

func (f TokenFactoryConfig) Validate() error {
	var err error

	if len(f.TokenSecret) < 32 {
		err = multierr.Append(
			err,
			errors.Errorf("TokenSecret length should be greater then 32, but is: %d", len(f.TokenSecret)),
		)
	}

	if f.AccessTokenDuration < time.Minute {
		err = multierr.Append(
			err,
			errors.Errorf("AccessTokenDuration should be greater then %s, but is: %s", time.Minute, f.AccessTokenDuration),
		)
	}

	if f.RefreshTokenDuration < time.Minute*10 {
		err = multierr.Append(
			err,
			errors.Errorf("RefreshTokenDuration should be greater then %s, but is: %s", time.Hour, f.RefreshTokenDuration),
		)
	}

	return err
}

type TokenFactory struct {
	tokenGen tokens.Generator[AccountCredentials]
	fc       TokenFactoryConfig
}

func NewAuthTokenFactory(fc TokenFactoryConfig) (*TokenFactory, error) {
	if err := fc.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config passed to factory")
	}

	gen, err := tokens.NewJWTGenerator[AccountCredentials](fc.TokenSecret)
	if err != nil {
		return nil, err
	}

	return &TokenFactory{fc: fc, tokenGen: gen}, nil
}

func (f *TokenFactory) NewTokenPair(instance AccountCredentials) (*Token, error) {
	acessToken, acessPayload, err := f.tokenGen.CreateToken(instance, f.fc.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := f.tokenGen.CreateToken(instance, f.fc.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:    acessToken,
		AccessPayload:  acessPayload,
		RefreshToken:   refreshToken,
		RefreshPayload: refreshPayload,
	}, nil
}

func (f *TokenFactory) NewAccesToken(refreshToken string) (*Token, error) {
	refreshPayload, err := f.tokenGen.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	acessToken, acessPayload, err := f.tokenGen.CreateToken(refreshPayload.GetInstance(), f.fc.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:    acessToken,
		AccessPayload:  acessPayload,
		RefreshToken:   refreshToken,
		RefreshPayload: refreshPayload,
	}, nil
}

func (f *TokenFactory) VerifyToken(token string) (*tokens.Payload[AccountCredentials], error) {
	return f.tokenGen.VerifyToken(token)
}
