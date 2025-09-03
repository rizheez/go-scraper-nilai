package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Scraper struct {
	client  *http.Client
	baseURL string
	cookie  string
	config  *Config
}

func NewScraper(config *Config) *Scraper {
	return &Scraper{
		client:  &http.Client{},
		baseURL: config.BaseURL,
		config:  config,
	}
}

func (s *Scraper) DoRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, s.baseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}

	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)
	req.Header.Set(HeaderReferer, s.baseURL+MediaEndpoint)
	req.Header.Set(HeaderOrigin, s.baseURL)
	req.Header.Set(HeaderAccept, AcceptJSON)
	if body != nil {
		req.Header.Set(HeaderContentType, ContentTypeForm+CharsetUTF8)
	}
	if s.cookie != "" {
		req.Header.Set(HeaderCookie, s.cookie)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal kirim request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request gagal status: %d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func (s *Scraper) GetBobotMK(fak, kodeProdi, kodePK, kls, kmk string) (Bobot, error) {
	data := "fak=" + fak + "&jrs=" + kodeProdi + "&prg=" + kodePK + "&kls=" + kls + "&kmk=" + kmk
	body, err := s.DoRequest(POST, "/_modul/mod_nilmk/aksi_nilmk.php?act=loadBOBOT", strings.NewReader(data))
	if err != nil {
		return Bobot{}, err
	}
	var bobotMK Bobot
	if err := json.Unmarshal(body, &bobotMK); err != nil {
		return Bobot{}, err
	}
	return bobotMK, nil

}

func (s *Scraper) GetRekapMHS() (*RekapMHSResponse, error) {
	data := "page=1&rows=500&"
	body, err := s.DoRequest(POST, "/_modul/mod_datamhs/aksi_datamhs.php?act=list", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	var resp RekapMHSResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Scraper) IsSessionValid() bool {
	body, err := s.DoRequest(GET, MediaEndpoint, nil)
	if err != nil {
		logf(LogError, "Gagal cek session: %v", err)
		return false
	}
	// logf(LogInfo, "Response Body : %s", string(body))
	return !(strings.Contains(string(body), "window.location = 'index.php'") || strings.Contains(string(body), "login") || strings.Contains(string(body), "Username"))
}

func (s *Scraper) SetProdi(kodeProdi, kodePK, smthn string) error {
	form := url.Values{}
	form.Set(FormPS, kodeProdi)
	form.Set(FormPK, kodePK)
	form.Set(FormSMTHN, smthn)
	_, err := s.DoRequest(POST, "/_modul/mod_prodi_smthn/aksi_prodi_smthn.php", strings.NewReader(form.Encode()))
	return err
}

func (s *Scraper) GetRekapMK() (*RekapMKResponse, error) {
	data := "page=1&rows=300&sort=hari&order=asc"
	body, err := s.DoRequest(POST, "/_modul/mod_nilmk/aksi_nilmk.php?act=rekapNILMK", strings.NewReader(data))
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
	form := url.Values{}
	form.Set(FormParam, infomk)
	form.Set(FormCetak, CetakValue)
	body, err := s.DoRequest(POST, "/_modul/mod_nilmk/aksi_nilmk.php?act=listNILMK", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	var hasil []Nilai
	if err := json.Unmarshal(body, &hasil); err != nil {
		return nil, err
	}
	return hasil, nil
}

// loadJurusan loads the jurusan data from file
func loadJurusan() (Jurusan, error) {
	// baca file jurusan.json (atau bisa juga dari API kalau ada)
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

	// tampilkan daftar
	fmt.Println()
	fmt.Println("=================================")
	log(LogInfo, "Daftar Jurusan:")
	fmt.Println("=================================")
	for i, j := range jurusanList {
		logf(LogInfo, "[%d] %s", i+1, j.NamaJrs)
	}

	// pilih input
	var selection int
	fmt.Printf("[INFO] Pilih jurusan (nomor): ")
	_, err = fmt.Scan(&selection)
	if err != nil {
		return Jurusan{}, fmt.Errorf("gagal membaca input: %w", err)
	}

	if selection < 1 || selection > len(jurusanList) {
		return Jurusan{}, fmt.Errorf("pilihan invalid: %d", selection)
	}

	// return jurusan terpilih
	return jurusanList[selection-1], nil
	// cfgFile, err := os.ReadFile(JurusanFile)
	// if err != nil {
	// 	logf(LogError, "Gagal baca %s: %v", JurusanFile, err)
	// 	return nil, fmt.Errorf("gagal baca %s: %w", JurusanFile, err)
	// }
	// var jurusanList []Jurusan
	// if err := json.Unmarshal(cfgFile, &jurusanList); err != nil {
	// 	logf(LogError, "Gagal parsing JSON: %v", err)
	// 	return nil, fmt.Errorf("gagal parsing JSON: %w", err)
	// }
	// return jurusanList, nil
}
// no changes
