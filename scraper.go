package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

type Scraper struct {
	client  *fasthttp.Client
	baseURL string
	cookie  string
	config  *Config
}

func NewScraper(config *Config) *Scraper {
	return &Scraper{
		client:  &fasthttp.Client{},
		baseURL: config.BaseURL,
		config:  config,
	}
}

func (s *Scraper) DoRequest(method, endpoint string, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(s.baseURL + endpoint)
	req.Header.SetMethod(method)

	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)
	req.Header.Set(HeaderReferer, s.baseURL+MediaEndpoint)
	req.Header.Set(HeaderOrigin, s.baseURL)
	req.Header.Set(HeaderAccept, AcceptJSON)

	if len(body) > 0 {
		req.SetBody(body)
		req.Header.Set(HeaderContentType, ContentTypeForm+CharsetUTF8)
	}
	if s.cookie != "" {
		req.Header.Set(HeaderCookie, s.cookie)
	}

	if err := s.client.Do(req, res); err != nil {
		return nil, fmt.Errorf("gagal kirim request: %w", err)
	}

	if res.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("request gagal status: %d", res.StatusCode())
	}

	return res.Body(), nil
}

func (s *Scraper) IsSessionValid() bool {
	body, err := s.DoRequest(fasthttp.MethodGet, MediaEndpoint, nil)
	if err != nil {
		logf(LogError, "Gagal cek session: %v", err)
		return false
	}
	return !(strings.Contains(string(body), "login") || strings.Contains(string(body), "Username"))
}

func (s *Scraper) SetProdi(kodeProdi, kodePK, smthn string) error {
	form := buildForm(map[string]string{
		FormPS:    kodeProdi,
		FormPK:    kodePK,
		FormSMTHN: smthn,
	})
	_, err := s.DoRequest(fasthttp.MethodPost, "/_modul/mod_prodi_smthn/aksi_prodi_smthn.php", form)
	return err
}

func (s *Scraper) GetRekapMK() (*RekapMKResponse, error) {
	data := []byte("page=1&rows=300&sort=hari&order=asc")
	body, err := s.DoRequest(fasthttp.MethodPost, "/_modul/mod_nilmk/aksi_nilmk.php?act=rekapNILMK", data)
	if err != nil {
		return nil, err
	}
	var resp RekapMKResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scraper) GetListNilai(infomk string) ([]Nilai, error) {
	form := buildForm(map[string]string{
		FormParam: infomk,
		FormCetak: CetakValue,
	})
	body, err := s.DoRequest(fasthttp.MethodPost, "/_modul/mod_nilmk/aksi_nilmk.php?act=listNILMK", form)
	if err != nil {
		return nil, err
	}
	var hasil []Nilai
	if err := json.Unmarshal(body, &hasil); err != nil {
		return nil, err
	}
	return hasil, nil
}

// --- util form encoder ---
func buildForm(values map[string]string) []byte {
	form := ""
	for k, v := range values {
		if form != "" {
			form += "&"
		}
		form += fmt.Sprintf("%s=%s", k, v)
	}
	return []byte(form)
}

// loadJurusan loads the jurusan data from file
func loadJurusan() (Jurusan, error) {
	data, err := os.ReadFile(JurusanFile)
	if err != nil {
		return Jurusan{}, fmt.Errorf("gagal baca %s: %w", JurusanFile, err)
	}

	var jurusanList []Jurusan
	if err := json.Unmarshal(data, &jurusanList); err != nil {
		return Jurusan{}, fmt.Errorf("gagal parsing JSON jurusan: %w", err)
	}

	if len(jurusanList) == 0 {
		return Jurusan{}, fmt.Errorf("tidak ada jurusan yang tersedia")
	}

	fmt.Println()
	fmt.Println("=================================")
	log(LogInfo, "Daftar Jurusan:")
	fmt.Println("=================================")
	for i, j := range jurusanList {
		logf(LogInfo, "[%d] %s", i+1, j.NamaJrs)
	}

	var selection int
	fmt.Printf("[INFO] Pilih jurusan (nomor): ")
	_, err = fmt.Scan(&selection)
	if err != nil {
		return Jurusan{}, fmt.Errorf("gagal membaca input: %w", err)
	}
	if selection < 1 || selection > len(jurusanList) {
		return Jurusan{}, fmt.Errorf("pilihan invalid: %d", selection)
	}

	return jurusanList[selection-1], nil
}
