package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func handleAuthentication(scraper *Scraper) error {
	if err := loadCookie(); err != nil {
		logf(LogError, "Gagal load cookie: %v", err)
	}
	scraper.cookie = cookie

	if scraper.cookie != "" && scraper.IsSessionValid() {
		log(LogInfo, "Cookie masih valid, skip login")
	} else {
		log(LogInfo, "Cookie Tidak Ditemukan / Cookie Tidak Valid")
		log(LogInfo, "Login ulang...")
		if !scraper.Login(scraper.config.Username, scraper.config.Password) {
			logf(LogError, "login gagal")
			return fmt.Errorf("login gagal")
		}
		if err := saveCookie(); err != nil {
			logf(LogWarn, "Gagal simpan cookie: %v", err)
		}
	}
	return nil
}

func (s *Scraper) Login(username, password string) bool {
	res, err := s.client.Get(s.baseURL + IndexEndpoint)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	session := ""
	for _, c := range res.Cookies() {
		if c.Name == CookiePHPSESSID {
			session = c.Value
		}
	}
	if session == "" {
		return false
	}

	hideValidation := generateValidation()
	hideIP := getRandomIP()

	data := url.Values{}
	data.Set(FormUsername, username)
	data.Set(FormPassword, password)
	data.Set(FormValidation, hideValidation)
	data.Set(FormHideValidation, hideValidation)
	data.Set(FormHideIP, hideIP)

	req, _ := http.NewRequest(POST, s.baseURL+LoginEndpoint, strings.NewReader(data.Encode()))
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderContentType, ContentTypeForm)
	req.Header.Set(HeaderCookie, CookiePHPSESSID+"="+session)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)

	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == CookiePHPSESSID {
			s.cookie = CookiePHPSESSID + "=" + c.Value
			cookie = s.cookie
		}
	}

	body, _ := io.ReadAll(resp.Body)
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

// --- AMBIL IP PUBLIK RANDOM ---
func getRandomIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log(LogWarn, "Gagal ambil IP publik, pakai default")
		return DefaultIP
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		log(LogWarn, "Gagal baca response IP publik, pakai default")
		return DefaultIP
	}

	return strings.TrimSpace(string(ip))
}

func loadCookie() error {
	data, err := os.ReadFile(CookieFile)
	if err != nil {
		// It's okay if the file doesn't exist, just return nil
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("gagal baca cookie.txt: %w", err)
	}
	cookie = string(data)
	log(LogInfo, "Cookie ditemukan!")
	return nil
}

// --- Simpan cookie ke file ---
func saveCookie() error {
	if cookie != "" {
		if err := os.WriteFile(CookieFile, []byte(cookie), 0644); err != nil {
			return fmt.Errorf("gagal simpan cookie: %w", err)
		}
	}
	return nil
}
// no changes
