package qrcode

import (
	"encoding/base64"
	"fmt"

	goqrcode "github.com/skip2/go-qrcode"
)

func CreateQRCode(data string) (string, error) {
	png, err := goqrcode.Encode(data, goqrcode.Medium, 256)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(png)
	return fmt.Sprintf("data:image/png;base64,%s", encoded), nil
}
