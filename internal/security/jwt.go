package security

import (
	"chopper/internal/domain"
	"fmt"
	"time"

	jwtPackage "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Jwt struct {
	secret         []byte
	expirationTime time.Duration
	issuer         string
	audience       string
}

func NewJwt(secret []byte, expirationTime time.Duration, issuer, audience string) *Jwt {
	return &Jwt{
		secret:         secret,
		expirationTime: expirationTime,
		issuer:         issuer,
		audience:       audience,
	}
}

type UserClaims struct {
	Id       uuid.UUID
	Username string
	Email    string
	Role     domain.Role
	jwtPackage.RegisteredClaims
}

func (j *Jwt) GenerateToken(id uuid.UUID, username, email string, role domain.Role) (string, error) {
	claims := UserClaims{
		Id:       id,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwtPackage.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwtPackage.NewNumericDate(time.Now()),
			ExpiresAt: jwtPackage.NewNumericDate(time.Now().Add(j.expirationTime)),
			Audience: jwtPackage.ClaimStrings{
				j.audience,
			},
		},
	}
	token := jwtPackage.NewWithClaims(jwtPackage.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *Jwt) ValidateToken(signedToken string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwtPackage.ParseWithClaims(signedToken, claims, func(signedToken *jwtPackage.Token) (any, error) {
		if _, ok := signedToken.Method.(*jwtPackage.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("wrong signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}
	if claims.Issuer != j.issuer {
		return claims, fmt.Errorf("wrong issuer")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return claims, fmt.Errorf("token is expired")
	}
	if token.Method.Alg() != jwtPackage.SigningMethodHS256.Alg() {
		return claims, fmt.Errorf("wrong signing method")
	}
	hasAudience := func(data jwtPackage.ClaimStrings, needAudience string) bool {
		for _, s := range data {
			if s == needAudience {
				return true
			}
		}
		return false
	}
	if hasAudience := hasAudience(claims.Audience, j.audience); !hasAudience {
		return claims, fmt.Errorf("wrong audience")
	}
	return claims, nil
}
