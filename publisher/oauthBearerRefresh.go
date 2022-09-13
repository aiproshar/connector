package publisher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"regexp"
	"time"
)

var (
	// Regex for sasl.oauthbearer.config, which constrains it to be
	// 1 or more name=value pairs with optional ignored whitespace
	oauthbearerConfigRegex = regexp.MustCompile("^(\\s*(\\w+)\\s*=\\s*(\\w+))+\\s*$")
	// Regex used to extract name=value pairs from sasl.oauthbearer.config
	oauthbearerNameEqualsValueRegex = regexp.MustCompile("(\\w+)\\s*=\\s*(\\w+)")
)

const (
	principalClaimNameKey = "principalClaimName"
	principalKey          = "principal"
	joseHeaderEncoded     = "eyJhbGciOiJub25lIn0" // {"alg":"none"}
)

// handleOAuthBearerTokenRefreshEvent generates an unsecured JWT based on the configuration defined
// in sasl.oauthbearer.config and sets the token on the client for use in any future authentication attempt.
// It must be invoked whenever kafka.OAuthBearerTokenRefresh appears on the client's event channel,
// which will occur whenever the client requires a token (i.e. when it first starts and when the
// previously-received token is 80% of the way to its expiration time).
func handleOAuthBearerTokenRefreshEvent(client kafka.Handle, e kafka.OAuthBearerTokenRefresh) {
	oauthBearerToken, retrieveErr := retrieveUnsecuredToken(e)
	if retrieveErr != nil {
		fmt.Fprintf(os.Stderr, "%% Token retrieval error: %v\n", retrieveErr)
		client.SetOAuthBearerTokenFailure(retrieveErr.Error())
	} else {
		setTokenError := client.SetOAuthBearerToken(oauthBearerToken)
		if setTokenError != nil {
			fmt.Fprintf(os.Stderr, "%% Error setting token and extensions: %v\n", setTokenError)
			client.SetOAuthBearerTokenFailure(setTokenError.Error())
		}
	}
}

func retrieveUnsecuredToken(e kafka.OAuthBearerTokenRefresh) (kafka.OAuthBearerToken, error) {
	config := e.Config
	if !oauthbearerConfigRegex.MatchString(config) {
		return kafka.OAuthBearerToken{}, fmt.Errorf("ignoring event %T due to malformed config: %s", e, config)
	}
	// set up initial map with default values
	oauthbearerConfigMap := map[string]string{
		principalClaimNameKey: "sub",
	}
	// parse the provided config and store name=value pairs in the map
	for _, kv := range oauthbearerNameEqualsValueRegex.FindAllStringSubmatch(config, -1) {
		oauthbearerConfigMap[kv[1]] = kv[2]
	}
	principalClaimName := oauthbearerConfigMap[principalClaimNameKey]
	principal := oauthbearerConfigMap[principalKey]
	// regexp is such that principalClaimName cannot end up blank,
	// so check for a blank principal (which will happen if it isn't specified)
	if principal == "" {
		return kafka.OAuthBearerToken{}, fmt.Errorf("ignoring event %T: no %s: %s", e, principalKey, config)
	}
	// do not proceed if there are any unknown name=value pairs
	if len(oauthbearerConfigMap) > 2 {
		return kafka.OAuthBearerToken{}, fmt.Errorf("ignoring event %T: unrecognized key(s): %s", e, config)
	}

	now := time.Now()
	nowSecondsSinceEpoch := now.Unix()

	// The token lifetime needs to be long enough to allow connection and a broker metadata query.
	// We then exit immediately after that, so no additional token refreshes will occur.
	// Therefore set the lifetime to be an hour (though anything on the order of a minute or more
	// would be fine).
	expiration := now.Add(time.Second * time.Duration(3600))
	expirationSecondsSinceEpoch := expiration.Unix()

	oauthbearerMapForJSON := map[string]interface{}{
		principalClaimName: principal,
		"iat":              nowSecondsSinceEpoch,
		"exp":              expirationSecondsSinceEpoch,
	}
	claimsJSON, _ := json.Marshal(oauthbearerMapForJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	jwsCompactSerialization := joseHeaderEncoded + "." + encodedClaims + "."
	extensions := map[string]string{}
	oauthBearerToken := kafka.OAuthBearerToken{
		TokenValue: jwsCompactSerialization,
		Expiration: expiration,
		Principal:  principal,
		Extensions: extensions,
	}
	return oauthBearerToken, nil
}
