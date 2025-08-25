package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func handleAuthentication(scraper *Scraper) error {
	if err := loadCookie(); err != nil {
		logf(LogError, "Gagal load cookie: %v", err)
	}
	scraper.cookie = cookie

	if scraper.cookie != "" && scraper.IsSessionValid() {
		log(LogInfo, "Cookie masih valid, skip login")
	} else {
		log(LogInfo, "Login ulang...")
		if !scraper.Login(scraper.config.Username, scraper.config.Password) {
			return fmt.Errorf("login gagal")
		}
		if err := saveCookie(); err != nil {
			logf(LogWarn, "Gagal simpan cookie: %v", err)
		}
	}
	return nil
}

func (s *Scraper) Login(username, password string) bool {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(s.baseURL + IndexEndpoint)
	req.Header.SetMethod(fasthttp.MethodGet)

	if err := s.client.Do(req, res); err != nil {
		return false
	}

	session := string(res.Header.PeekCookie(CookiePHPSESSID))
	if session == "" {
		return false
	}

	hideValidation := generateValidation()
	hideIP := getRandomIP()

	data := buildForm(map[string]string{
		FormUsername:       username,
		FormPassword:       password,
		FormValidation:     hideValidation,
		FormHideValidation: hideValidation,
		FormHideIP:         hideIP,
	})

	req.Reset()
	req.SetRequestURI(s.baseURL + LoginEndpoint)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(data)

	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderContentType, ContentTypeForm)
	req.Header.Set(HeaderCookie, CookiePHPSESSID+"="+session)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)

	if err := s.client.Do(req, res); err != nil {
		return false
	}

	cookies := res.Header.PeekCookie(CookiePHPSESSID)
	if len(cookies) > 0 {
		trim := strings.SplitN(string(cookies), ";", 2)
		s.cookie = trim[0]
		cookie = s.cookie
	}

	body := res.Body()
	return strings.Contains(string(body), `"success":true`)
}

func generateValidation() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, 3)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return string(code)
}

// --- Ambil IP publik random ---
func getRandomIP() string {
	status, body, err := fasthttp.Get(nil, "https://api.ipify.org")
	if err != nil || status != fasthttp.StatusOK {
		log(LogWarn, "Gagal ambil IP publik, pakai default")
		return DefaultIP
	}
	return strings.TrimSpace(string(body))
}

func loadCookie() error {
	data, err := os.ReadFile(CookieFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("gagal baca cookie.txt: %w", err)
	}
	cookie = string(data)
	log(LogInfo, "Cookie ditemukan!")
	return nil
}

func saveCookie() error {
	if cookie != "" {
		if err := os.WriteFile(CookieFile, []byte(cookie), 0644); err != nil {
			return fmt.Errorf("gagal simpan cookie: %w", err)
		}
	}
	return nil
}
