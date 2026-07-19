package auth

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bcpriok/pantas/internal/config"
)

func TestLoginCaptchaIsEncryptedCaseInsensitiveAndExpires(t *testing.T) {
	service := &Service{cfg: config.Config{AppSecret: "01234567890123456789012345678901"}}
	now := time.Date(2026, time.July, 19, 12, 0, 0, 0, time.UTC)
	token, imageBytes, err := service.newLoginCaptcha(now)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(imageBytes, []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}) {
		t.Fatal("captcha image is not PNG")
	}
	payload, err := service.openCaptcha(token)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(payload), "|", 2)
	if len(parts) != 2 || len(parts[1]) != 5 {
		t.Fatalf("invalid encrypted captcha payload: %q", payload)
	}
	answer := parts[1]
	if !service.verifyLoginCaptcha(token, strings.ToLower(answer), now.Add(time.Minute)) {
		t.Fatal("valid captcha answer was rejected")
	}
	if service.verifyLoginCaptcha(token, answer, now.Add(6*time.Minute)) {
		t.Fatal("expired captcha answer was accepted")
	}
	if service.verifyLoginCaptcha(token, "WRONG", now.Add(time.Minute)) {
		t.Fatal("incorrect captcha answer was accepted")
	}
}
