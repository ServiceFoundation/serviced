package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/control-center/serviced/config"
	"github.com/control-center/serviced/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/zenoss/glog"
	"net/http"
	"strings"
	"sync"
)

type jwtAuth0Claims struct {
	Issuer    string   `json:"iss,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	Audience  []string `json:"aud,omitempty"`
	Groups    []string `json:"https://zenoss.com/groups,omitempty"`
	Subject   string   `json:"sub,omitempty"`
}

func (t *jwtAuth0Claims) Valid() error {
	if t.Expired() {
		return ErrAuth0TokenExpired
	}
	opts := config.GetOptions()
	expectedIssuer := fmt.Sprintf("https://%s/", opts.Auth0Domain)
	if t.Issuer != expectedIssuer {
		return ErrAuth0TokenBadIssuer
	}
	if !utils.StringInSlice(opts.Auth0Audience, t.Audience) {
		return ErrAuth0TokenBadAudience
	}
	return nil
}

type Auth0Token interface {
	HasAdminAccess() bool
	User() string
	Expiration() int64
}

type jwtAuth0RestToken struct {
	*jwtAuth0Claims
	authIdentity Identity
	restToken    string
}

func (t *jwtAuth0Claims) Expired() bool {
	now := jwt.TimeFunc().UTC().Unix()
	return now >= t.ExpiresAt
}

func (t *jwtAuth0Claims) Expiration() int64 {
	return t.ExpiresAt
}

func (t *jwtAuth0Claims) User() string {
	// Auth0 returns username in the subj field in the form <source>|<username>.
	// Strip off the source and only return the username.
	// Per review comment, there may be two '|' characters in some cases. We are
	// using the last field, which should be the username.
	fields := strings.Split(t.Subject, "|")
	return fields[len(fields)-1]
}

func (t *jwtAuth0Claims) HasAdminAccess() bool {
	opts := config.GetOptions()
	auth0Group := opts.Auth0Group
	if !utils.StringInSlice(auth0Group, t.Groups) {
		glog.Warning("Auth0 Admin access denied - '" + auth0Group + "' not found in Groups.")
		return false
	}
	return true
}

type JSONWebkeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Jwks struct {
	Keys []JSONWebkeys `json:"keys"`
	m    sync.Mutex
}

func (j *Jwks) refreshIfEmpty() {
	j.m.Lock()
	defer j.m.Unlock()
	if len(j.Keys) == 0 {
		glog.V(0).Info("Fetching jwks key from auth0")
		opts := config.GetOptions()
		auth0Domain := opts.Auth0Domain
		resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain))

		if err != nil {
			glog.Warning("error getting well-known jwks: ", err)
			return
		}
		defer resp.Body.Close()

		var newjwks = Jwks{}
		err = json.NewDecoder(resp.Body).Decode(&newjwks)

		if err != nil {
			glog.Warning("error decoding JWKS keys from JSON: ", err)
			return
		}
		j.Keys = newjwks.Keys
	}
	return
}

var auth0Jwks = &Jwks{}

func getPemCert(token *jwt.Token) ([]byte, error) {
	cert := ""
	if auth0Jwks == nil {
		var jwks = Jwks{}
		auth0Jwks = &jwks
	}
	auth0Jwks.refreshIfEmpty()

	x5c := auth0Jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == auth0Jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		glog.Warning("Unable to find appropriate key.")
		err := errors.New("unable to find appropriate key")
		return []byte(cert), err
	}

	return []byte(cert), nil
}

func getRSAPublicKey(token *jwt.Token) (*rsa.PublicKey, error) {
	certBytes, err := getPemCert(token)
	if err != nil {
		glog.Warning("error getting Pem Cert from auth0: ", err)
		return nil, err
	}
	block, _ := pem.Decode(certBytes)
	var cert *x509.Certificate
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		glog.Warning("error parsing certificate: ", err)
		return nil, err
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return rsaPublicKey, nil
}

/*
	See https://auth0.com/docs/api-auth/tutorials/verify-access-token for information on
	validating auth0 tokens. Per https://jwt.io/, the jwt-go library validates exp,
	but not iss or sub.
*/
func ParseAuth0Token(token string) (Auth0Token, error) {
	claims := &jwtAuth0Claims{}
	identity := &jwtIdentity{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm matches the key
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			glog.Warning("error getting RSA key from PEM: ", ErrInvalidSigningMethod)
			return nil, ErrInvalidSigningMethod
		}

		// extract public key from token
		key, err := getRSAPublicKey(token)
		if err != nil {
			glog.Warning("error getting RSA key from PEM: ", err)
			return nil, fmt.Errorf("error getting RSA key from PEM: %v\n", err)
		}
		return key, nil
	})
	if err != nil {
		if verr, ok := err.(*jwt.ValidationError); ok {
			glog.Warning("Validation error from jwt.ParseWIthClaims(): ", verr)
			if verr.Inner != nil && (verr.Inner == ErrIdentityTokenExpired || verr.Inner == ErrIdentityTokenBadSig) {
				return nil, verr.Inner
			}
			if verr.Errors&jwt.ValidationErrorExpired != 0 || verr.Inner != nil && verr.Inner == ErrRestTokenExpired {
				return nil, ErrRestTokenExpired
			}
			if verr.Errors&(jwt.ValidationErrorSignatureInvalid|jwt.ValidationErrorUnverifiable) != 0 {
				return nil, ErrRestTokenBadSig
			}
			if verr.Errors&(jwt.ValidationErrorMalformed) != 0 {
				return nil, ErrBadRestToken
			}
			if verr.Inner != nil {
				return nil, verr.Inner
			}
			if verr != nil {
				return nil, verr
			}
		}
		return nil, err
	}
	if claims, ok := parsed.Claims.(*jwtAuth0Claims); ok && parsed.Valid {
		restToken := &jwtAuth0RestToken{}
		restToken.jwtAuth0Claims = claims
		restToken.authIdentity = identity
		restToken.restToken = token
		return restToken, nil
	}
	glog.Warning("ParseAuth0Token: ", ErrIdentityTokenInvalid)
	return nil, ErrIdentityTokenInvalid
}