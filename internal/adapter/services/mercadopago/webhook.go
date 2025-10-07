package mercadopago

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type CheckWebhookSignatureInput struct {
	XSignature string
	XRequestId string
	DataID     string
	Secret     string
}

func CheckWebhookSignature(input CheckWebhookSignatureInput) bool {
	xSignature := input.XSignature
	xRequestId := input.XRequestId
	dataID := input.DataID
	parts := strings.Split(xSignature, ",")
	var ts, hash string
	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			if key == "ts" {
				ts = value
			} else if key == "v1" {
				hash = value
			}
		}
	}
	manifest := fmt.Sprintf("id:%v;request-id:%v;ts:%v;", dataID, xRequestId, ts)
	hmacHash := hmac.New(sha256.New, []byte(input.Secret))
	hmacHash.Write([]byte(manifest))
	sha := hex.EncodeToString(hmacHash.Sum(nil))

	fmt.Println(input.Secret)
	fmt.Println(os.Getenv("MP_WEBHOOK_SECRET"))
	fmt.Println(hash)
	fmt.Println(sha)
	return strings.EqualFold(hash, sha)
}
