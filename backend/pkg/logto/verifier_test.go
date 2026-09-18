package logto

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestVerifierValidatesLogtoTokenAndOrganization(t *testing.T) {
	key, jwks := testRSAKey(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jwks)
	}))
	t.Cleanup(server.Close)

	verifier := newTestVerifier(t, server.URL)
	token := signedTestToken(t, key, jwt.MapClaims{
		"iss":             "https://logto.example/oidc",
		"aud":             "http://localhost:8000",
		"sub":             "user-1",
		"organization_id": "org-1",
		"exp":             time.Now().Add(time.Hour).Unix(),
		"iat":             time.Now().Unix(),
	})

	claims, err := verifier.Parse(t.Context(), "Bearer "+token)
	require.NoError(t, err)
	require.Equal(t, "user-1", claims.Subject)
	require.Equal(t, "org-1", claims.OrganizationID)
	require.Equal(t, "org-1", mustOrganizationID(t, verifier, claims))
}

func TestVerifierRejectsInvalidIssuerAudienceAndExpiration(t *testing.T) {
	key, jwks := testRSAKey(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jwks)
	}))
	t.Cleanup(server.Close)
	verifier := newTestVerifier(t, server.URL)

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "issuer",
			claims: jwt.MapClaims{
				"iss": "https://wrong.example/oidc", "aud": "http://localhost:8000", "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "audience",
			claims: jwt.MapClaims{
				"iss": "https://logto.example/oidc", "aud": "http://wrong.example", "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
			},
		},
		{
			name: "expired",
			claims: jwt.MapClaims{
				"iss": "https://logto.example/oidc", "aud": "http://localhost:8000", "sub": "user-1", "exp": time.Now().Add(-time.Minute).Unix(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := verifier.Parse(t.Context(), "Bearer "+signedTestToken(t, key, tt.claims))
			require.Error(t, err)
		})
	}
}

func TestVerifierRequiresOrganizationForTenantAccess(t *testing.T) {
	key, jwks := testRSAKey(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jwks)
	}))
	t.Cleanup(server.Close)
	verifier := newTestVerifier(t, server.URL)
	token := signedTestToken(t, key, jwt.MapClaims{
		"iss": "https://logto.example/oidc", "aud": "http://localhost:8000", "sub": "user-1", "exp": time.Now().Add(time.Hour).Unix(),
	})

	claims, err := verifier.Parse(t.Context(), "Bearer "+token)
	require.NoError(t, err)
	_, err = verifier.OrganizationID(claims)
	require.ErrorIs(t, err, ErrOrganizationRequired)
}

func TestNewVerifierRequiresSecurityConfiguration(t *testing.T) {
	_, err := NewVerifier(viper.New())
	require.EqualError(t, err, "security.logto.issuer is required")
}

func newTestVerifier(t *testing.T, jwksURL string) *Verifier {
	t.Helper()
	conf := viper.New()
	conf.Set("security.logto.issuer", "https://logto.example/oidc")
	conf.Set("security.logto.audience", "http://localhost:8000")
	conf.Set("security.logto.jwks_url", jwksURL)
	conf.Set("security.logto.organization_claim", "organization_id")
	verifier, err := NewVerifier(conf)
	require.NoError(t, err)
	return verifier
}

func testRSAKey(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	encode := func(value []byte) string { return base64.RawURLEncoding.EncodeToString(value) }
	n := key.PublicKey.N.Bytes()
	e := big.NewInt(int64(key.PublicKey.E)).Bytes()
	body, err := json.Marshal(jwksResponse{Keys: []jwk{{Kid: "test-key", Kty: "RSA", N: encode(n), E: encode(e)}}})
	require.NoError(t, err)
	return key, body
}

func signedTestToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key"
	signed, err := token.SignedString(key)
	require.NoError(t, err)
	return signed
}

func mustOrganizationID(t *testing.T, verifier *Verifier, claims *Claims) string {
	t.Helper()
	organizationID, err := verifier.OrganizationID(claims)
	require.NoError(t, err)
	return organizationID
}
