package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"
	"time"
)

const (
	LoginCaptchaCookieName = "pantas_login_captcha"
	captchaLifetime        = 5 * time.Minute
	captchaAlphabet        = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

var captchaAdditionalData = []byte("pantas-login-captcha-v1")

// NewLoginCaptcha returns an encrypted, short-lived proof cookie and a visual
// challenge. The answer is never exposed in JSON or browser-readable storage.
func (s *Service) NewLoginCaptcha() (string, []byte, error) {
	return s.newLoginCaptcha(time.Now())
}

func (s *Service) newLoginCaptcha(now time.Time) (string, []byte, error) {
	codeBytes := make([]byte, 5)
	randomness := make([]byte, 64)
	if _, err := rand.Read(randomness); err != nil {
		return "", nil, err
	}
	for index := range codeBytes {
		codeBytes[index] = captchaAlphabet[int(randomness[index])%len(captchaAlphabet)]
	}
	code := string(codeBytes)
	payload := fmt.Sprintf("%d|%s", now.Add(captchaLifetime).Unix(), code)
	token, err := s.sealCaptcha([]byte(payload))
	if err != nil {
		return "", nil, err
	}
	imageBytes, err := renderCaptchaPNG(code, randomness[8:])
	if err != nil {
		return "", nil, err
	}
	return token, imageBytes, nil
}

func (s *Service) VerifyLoginCaptcha(token, answer string) bool {
	return s.verifyLoginCaptcha(token, answer, time.Now())
}

func (s *Service) verifyLoginCaptcha(token, answer string, now time.Time) bool {
	answer = strings.ToUpper(strings.TrimSpace(answer))
	if len(answer) != 5 {
		return false
	}
	payload, err := s.openCaptcha(token)
	if err != nil {
		return false
	}
	parts := strings.SplitN(string(payload), "|", 2)
	if len(parts) != 2 {
		return false
	}
	expiresAt, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || now.Unix() > expiresAt || expiresAt > now.Add(10*time.Minute).Unix() {
		return false
	}
	if len(parts[1]) != len(answer) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(parts[1]), []byte(answer)) == 1
}

func (s *Service) captchaAEAD() (cipher.AEAD, error) {
	key := sha256.Sum256([]byte("pantas-login-captcha|" + s.cfg.AppSecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func (s *Service) sealCaptcha(payload []byte) (string, error) {
	aead, err := s.captchaAEAD()
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	sealed := aead.Seal(nil, nonce, payload, captchaAdditionalData)
	combined := append(nonce, sealed...)
	return base64.RawURLEncoding.EncodeToString(combined), nil
}

func (s *Service) openCaptcha(token string) ([]byte, error) {
	aead, err := s.captchaAEAD()
	if err != nil {
		return nil, err
	}
	combined, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(combined) <= aead.NonceSize() {
		return nil, fmt.Errorf("invalid captcha token")
	}
	nonce, sealed := combined[:aead.NonceSize()], combined[aead.NonceSize():]
	return aead.Open(nil, nonce, sealed, captchaAdditionalData)
}

var captchaGlyphs = map[byte][7]uint8{
	'2': {0b11110, 0b00001, 0b00001, 0b01110, 0b10000, 0b10000, 0b11111},
	'3': {0b11110, 0b00001, 0b00001, 0b01110, 0b00001, 0b00001, 0b11110},
	'4': {0b10010, 0b10010, 0b10010, 0b11111, 0b00010, 0b00010, 0b00010},
	'5': {0b11111, 0b10000, 0b10000, 0b11110, 0b00001, 0b00001, 0b11110},
	'6': {0b01111, 0b10000, 0b10000, 0b11110, 0b10001, 0b10001, 0b01110},
	'7': {0b11111, 0b00001, 0b00010, 0b00100, 0b01000, 0b01000, 0b01000},
	'8': {0b01110, 0b10001, 0b10001, 0b01110, 0b10001, 0b10001, 0b01110},
	'9': {0b01110, 0b10001, 0b10001, 0b01111, 0b00001, 0b00001, 0b11110},
	'A': {0b01110, 0b10001, 0b10001, 0b11111, 0b10001, 0b10001, 0b10001},
	'B': {0b11110, 0b10001, 0b10001, 0b11110, 0b10001, 0b10001, 0b11110},
	'C': {0b01111, 0b10000, 0b10000, 0b10000, 0b10000, 0b10000, 0b01111},
	'D': {0b11110, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b11110},
	'E': {0b11111, 0b10000, 0b10000, 0b11110, 0b10000, 0b10000, 0b11111},
	'F': {0b11111, 0b10000, 0b10000, 0b11110, 0b10000, 0b10000, 0b10000},
	'G': {0b01111, 0b10000, 0b10000, 0b10111, 0b10001, 0b10001, 0b01110},
	'H': {0b10001, 0b10001, 0b10001, 0b11111, 0b10001, 0b10001, 0b10001},
	'J': {0b00111, 0b00010, 0b00010, 0b00010, 0b10010, 0b10010, 0b01100},
	'K': {0b10001, 0b10010, 0b10100, 0b11000, 0b10100, 0b10010, 0b10001},
	'L': {0b10000, 0b10000, 0b10000, 0b10000, 0b10000, 0b10000, 0b11111},
	'M': {0b10001, 0b11011, 0b10101, 0b10101, 0b10001, 0b10001, 0b10001},
	'N': {0b10001, 0b11001, 0b11001, 0b10101, 0b10011, 0b10011, 0b10001},
	'P': {0b11110, 0b10001, 0b10001, 0b11110, 0b10000, 0b10000, 0b10000},
	'Q': {0b01110, 0b10001, 0b10001, 0b10001, 0b10101, 0b10010, 0b01101},
	'R': {0b11110, 0b10001, 0b10001, 0b11110, 0b10100, 0b10010, 0b10001},
	'S': {0b01111, 0b10000, 0b10000, 0b01110, 0b00001, 0b00001, 0b11110},
	'T': {0b11111, 0b00100, 0b00100, 0b00100, 0b00100, 0b00100, 0b00100},
	'U': {0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b01110},
	'V': {0b10001, 0b10001, 0b10001, 0b10001, 0b10001, 0b01010, 0b00100},
	'W': {0b10001, 0b10001, 0b10001, 0b10101, 0b10101, 0b10101, 0b01010},
	'X': {0b10001, 0b10001, 0b01010, 0b00100, 0b01010, 0b10001, 0b10001},
	'Y': {0b10001, 0b10001, 0b01010, 0b00100, 0b00100, 0b00100, 0b00100},
	'Z': {0b11111, 0b00001, 0b00010, 0b00100, 0b01000, 0b10000, 0b11111},
}

func renderCaptchaPNG(code string, randomness []byte) ([]byte, error) {
	const width, height, scale = 220, 64, 5
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	background := color.RGBA{R: 244, G: 248, B: 251, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvas.SetRGBA(x, y, background)
		}
	}
	lineColor := color.RGBA{R: 183, G: 212, B: 231, A: 255}
	for index := 0; index < 6; index++ {
		offset := index * 4
		drawCaptchaLine(canvas, int(randomness[offset])%width, int(randomness[offset+1])%height, int(randomness[offset+2])%width, int(randomness[offset+3])%height, lineColor)
	}
	palette := []color.RGBA{{R: 11, G: 59, B: 102, A: 255}, {R: 18, G: 111, B: 168, A: 255}, {R: 36, G: 59, B: 83, A: 255}}
	for index := 0; index < len(code); index++ {
		glyph := captchaGlyphs[code[index]]
		originX := 15 + index*40 + (int(randomness[(30+index)%len(randomness)])%5 - 2)
		originY := 14 + (int(randomness[(38+index)%len(randomness)])%7 - 3)
		skew := int(randomness[(46+index)%len(randomness)])%3 - 1
		for row, bits := range glyph {
			for column := 0; column < 5; column++ {
				if bits&(1<<uint(4-column)) == 0 {
					continue
				}
				x := originX + column*scale + (row-3)*skew
				y := originY + row*scale
				for dy := 0; dy < scale-1; dy++ {
					for dx := 0; dx < scale-1; dx++ {
						if image.Pt(x+dx, y+dy).In(canvas.Bounds()) {
							canvas.SetRGBA(x+dx, y+dy, palette[index%len(palette)])
						}
					}
				}
			}
		}
	}
	drawCaptchaLine(canvas, 8, 34, 211, 29+int(randomness[55%len(randomness)])%10, color.RGBA{R: 229, G: 174, B: 61, A: 255})
	var encoded bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	if err := encoder.Encode(&encoded, canvas); err != nil {
		return nil, err
	}
	return encoded.Bytes(), nil
}

func drawCaptchaLine(canvas *image.RGBA, x0, y0, x1, y1 int, lineColor color.RGBA) {
	dx, sx := absInt(x1-x0), -1
	if x0 < x1 {
		sx = 1
	}
	dy, sy := -absInt(y1-y0), -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy
	for {
		if image.Pt(x0, y0).In(canvas.Bounds()) {
			canvas.SetRGBA(x0, y0, lineColor)
		}
		if x0 == x1 && y0 == y1 {
			return
		}
		twice := 2 * err
		if twice >= dy {
			err += dy
			x0 += sx
		}
		if twice <= dx {
			err += dx
			y0 += sy
		}
	}
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
