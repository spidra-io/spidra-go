package spidra

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyWebhook validates that a webhook request was signed by Spidra using HMAC-SHA256.
// signatureHeader is expected to be in the format "sha256=<hex>" as sent in X-Spidra-Signature.
func VerifyWebhook(payload []byte, signatureHeader, secret string) bool {
	if len(payload) == 0 || signatureHeader == "" || secret == "" {
		return false
	}

	sig := strings.TrimSpace(signatureHeader)
	if strings.HasPrefix(sig, "sha256=") {
		sig = strings.TrimPrefix(sig, "sha256=")
	}

	sig = strings.ToLower(sig)
	if len(sig) != 64 {
		return false
	}

	expectedSigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	computedSigBytes := mac.Sum(nil)

	return hmac.Equal(expectedSigBytes, computedSigBytes)
}

// VerifyWebhookString validates a webhook request when the payload is a string.
func VerifyWebhookString(payload, signatureHeader, secret string) bool {
	return VerifyWebhook([]byte(payload), signatureHeader, secret)
}

