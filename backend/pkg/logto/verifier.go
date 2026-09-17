package logto

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
)

var ErrOrganizationRequired = errors.New("logto organization claim is required")

type Config struct {
	Issuer            string
	Audience          string
	JWKSURL           string
	OrganizationClaim string
}

type Claims struct {
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}

type Verifier struct {
	config Config
	client *http.Client

	mu   sync.RWMutex
	keys map[string]any
	sf   singleflight.Group
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

func NewVerifier(conf *viper.Viper) (*Verifier, error) {
	c := Config{
		Issuer:            strings.TrimRight(conf.GetString("security.logto.issuer"), "/"),
		Audience:          strings.TrimSpace(conf.GetString("security.logto.audience")),
		JWKSURL:           strings.TrimSpace(conf.GetString("security.logto.jwks_url")),
		OrganizationClaim: strings.TrimSpace(conf.GetString("security.logto.organization_claim")),
	}
	if c.Issuer == "" {
		return nil, errors.New("security.logto.issuer is required")
	}
	if c.Audience == "" {
		return nil, errors.New("security.logto.audience is required")
	}
	if c.JWKSURL == "" {
		return nil, errors.New("security.logto.jwks_url is required")
	}
	if c.OrganizationClaim == "" {
		return nil, errors.New("security.logto.organization_claim is required")
	}
	return &Verifier{config: c, client: http.DefaultClient, keys: make(map[string]any)}, nil
}

func (v *Verifier) Parse(ctx context.Context, authorization string) (*Claims, error) {
	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if tokenString == "" {
		return nil, errors.New("authorization token is required")
	}

	rawClaims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, rawClaims, func(token *jwt.Token) (any, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA, *jwt.SigningMethodECDSA:
		default:
			return nil, fmt.Errorf("unsupported logto signing method %q", token.Method.Alg())
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("logto token has no key id")
		}
		return v.publicKey(ctx, kid)
	}, jwt.WithIssuer(v.config.Issuer), jwt.WithAudience(v.config.Audience), jwt.WithValidMethods([]string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}))
	if err != nil || !token.Valid {
		if err == nil {
			err = errors.New("logto token is invalid")
		}
		return nil, err
	}
	subject, _ := rawClaims["sub"].(string)
	if subject == "" {
		return nil, errors.New("logto token subject is required")
	}
	organizationID, _ := rawClaims[v.config.OrganizationClaim].(string)
	claims := &Claims{
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
		},
	}
	return claims, nil
}

func (v *Verifier) OrganizationID(claims *Claims) (string, error) {
	if claims == nil || claims.OrganizationID == "" {
		return "", ErrOrganizationRequired
	}
	return claims.OrganizationID, nil
}

func (v *Verifier) publicKey(ctx context.Context, kid string) (any, error) {
	v.mu.RLock()
	key := v.keys[kid]
	v.mu.RUnlock()
	if key != nil {
		return key, nil
	}

	res, err, _ := v.sf.Do(v.config.JWKSURL, func() (any, error) {
		v.mu.RLock()
		cached := v.keys[kid]
		v.mu.RUnlock()
		if cached != nil {
			return cached, nil
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.config.JWKSURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := v.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch logto jwks: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch logto jwks: unexpected status %s", resp.Status)
		}

		var body jwksResponse
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode logto jwks: %w", err)
		}
		keys := make(map[string]any, len(body.Keys))
		for _, item := range body.Keys {
			if item.Kid == "" {
				continue
			}
			var parsed any
			var err error
			switch item.Kty {
			case "RSA":
				parsed, err = parseRSAKey(item)
			case "EC":
				parsed, err = parseECKey(item)
			default:
				continue
			}
			if err != nil {
				return nil, err
			}
			keys[item.Kid] = parsed
		}
		v.mu.Lock()
		v.keys = keys
		key = v.keys[kid]
		v.mu.Unlock()
		if key == nil {
			return nil, fmt.Errorf("logto jwks key %q not found", kid)
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func parseECKey(item jwk) (*ecdsa.PublicKey, error) {
	curves := map[string]elliptic.Curve{
		"P-256": elliptic.P256(),
		"P-384": elliptic.P384(),
		"P-521": elliptic.P521(),
	}
	curve, ok := curves[item.Crv]
	if !ok {
		return nil, fmt.Errorf("unsupported Logto EC curve %q", item.Crv)
	}
	decode := func(value string) ([]byte, error) {
		return base64.RawURLEncoding.DecodeString(value)
	}
	x, err := decode(item.X)
	if err != nil {
		return nil, fmt.Errorf("decode Logto EC x coordinate: %w", err)
	}
	y, err := decode(item.Y)
	if err != nil {
		return nil, fmt.Errorf("decode Logto EC y coordinate: %w", err)
	}
	if len(x) == 0 || len(y) == 0 {
		return nil, errors.New("Logto EC key is incomplete")
	}
	return &ecdsa.PublicKey{Curve: curve, X: new(big.Int).SetBytes(x), Y: new(big.Int).SetBytes(y)}, nil
}

func parseRSAKey(item jwk) (*rsa.PublicKey, error) {
	decode := func(value string) ([]byte, error) {
		return base64.RawURLEncoding.DecodeString(value)
	}
	n, err := decode(item.N)
	if err != nil {
		return nil, fmt.Errorf("decode logto jwks modulus: %w", err)
	}
	e, err := decode(item.E)
	if err != nil {
		return nil, fmt.Errorf("decode logto jwks exponent: %w", err)
	}
	if len(n) == 0 || len(e) == 0 {
		return nil, errors.New("logto jwks RSA key is incomplete")
	}
	exponent := 0
	for _, value := range e {
		exponent = exponent<<8 | int(value)
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: exponent}, nil
}
