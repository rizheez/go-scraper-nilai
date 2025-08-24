package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

const (
	// Folder paths
	JSONFolder  = "nilai_json"
	ExcelFolder = "nilai_excel"

	// File names
	CookieFile    = "cookie.txt"
	JurusanFile   = "jurusan.json"
	MediaEndpoint = "/media.php"
	IndexEndpoint = "/index.php"
	LoginEndpoint = "/ceklogin.php?h="

	// HTTP methods
	GET  = "GET"
	POST = "POST"

	// Content types
	ContentTypeForm = "application/x-www-form-urlencoded"
	ContentTypeJSON = "application/json"

	// HTTP headers
	HeaderUserAgent      = "User-Agent"
	HeaderContentType    = "Content-Type"
	HeaderCookie         = "Cookie"
	HeaderXRequestedWith = "X-Requested-With"
	HeaderReferer        = "Referer"
	HeaderOrigin         = "Origin"
	HeaderAccept         = "Accept"

	// Header values
	XMLHttpRequest = "XMLHttpRequest"
	AcceptJSON     = "application/json, text/javascript, */*; q=0.01"
	CharsetUTF8    = "; charset=UTF-8"

	// Form fields
	FormUsername       = "username"
	FormPassword       = "password"
	FormValidation     = "validation"
	FormHideValidation = "hide_validation"
	FormHideIP         = "hide_ipnya"
	FormParam          = "param"
	FormCetak          = "cetak"
	FormPS             = "ps"
	FormPK             = "pk"
	FormSMTHN          = "smthn"

	// Cookie names
	CookiePHPSESSID = "PHPSESSID"

	// Default values
	DefaultIP  = "182.8.179.9"
	RegValue   = "REG"
	CetakValue = "1"

	// Excel headers
	ExcelNIM       = "NIM"
	ExcelNama      = "Nama Peserta"
	ExcelAngka     = "Angka"
	ExcelHuruf     = "Huruf"
	ExcelKehadiran = "Kehadiran"
	ExcelProjek    = "Projek"
	ExcelQuiz      = "Quiz"
	ExcelTugas     = "Tugas"
	ExcelUTS       = "UTS"
	ExcelUAS       = "UAS"

	// Progress bar
	ProgressBarLength = 40

	// Log levels
	LogInfo  = "[INFO]"
	LogError = "[ERROR]"
	LogWarn  = "[WARN]"
	LogDebug = "[DEBUG]"
)

var (
	baseURL string
	cookie  string
)

// log prints a message with a log level prefix
func log(level, message string) {
	fmt.Printf("%s %s\n", level, message)
}

// logf prints a formatted message with a log level prefix
func logf(level, format string, args ...interface{}) {
	fmt.Printf("%s %s\n", level, fmt.Sprintf(format, args...))
}

// Scraper holds the state and configuration for scraping
type Scraper struct {
	client  *http.Client
	baseURL string
	cookie  string
	config  *Config
}

// NewScraper creates a new Scraper instance
func NewScraper(config *Config) *Scraper {
	return &Scraper{
		client:  &http.Client{},
		baseURL: config.BaseURL,
		cookie:  "",
		config:  config,
	}
}

// DoRequest performs an HTTP request and returns the response body
func (s *Scraper) DoRequest(method, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, s.baseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}

	// Set headers
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

	// Perform request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request: %w", err)
	}
	defer res.Body.Close()

	// Check status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request gagal dengan status: %d", res.StatusCode)
	}

	// Read response body
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response: %w", err)
	}

	return respBody, nil
}

// User agent list (tidak berubah)
var userAgents = []string{
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) Gecko/20100101 Firefox/130.0",
	"Mozilla/5.0 (X11; Arch Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.6065.0 Safari/537.36",
	"Mozilla/5.0 (X11; Arch Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.7 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
}

type Jurusan struct {
	JrsID   string `json:"jrsid"`
	KodeJrs string `json:"kodejrs"`
	NamaJrs string `json:"namajrs"`
}

type MataKuliah struct {
	JID        string `json:"jid"`
	Namamk     string `json:"namamk"`
	SKS        string `json:"sks"`
	Semester   string `json:"semester"`
	Namadosen  string `json:"namadosen"`
	NamaKelas  string `json:"nama_kelas"`
	Hari       string `json:"namahari"`
	JamKuliah  string `json:"jamkuliah"`
	JmlPeserta string `json:"jmlpeserta"`
	Bisainput  int    `json:"bisainput"`
	Cetak      string `json:"cetak"`
	Infomk     string `json:"infomk"`
	Kelas      string `json:"kelas"`
}

type Nilai struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Hadir    string `json:"hadir"`
	Projek   string `json:"projek"`
	Quiz     string `json:"quiz"`
	Tugas    string `json:"tugas"`
	UTS      string `json:"uts"`
	UAS      string `json:"uas"`
	NilAngka string `json:"nil_angka"`
	NilHuruf string `json:"nil_huruf"`
	KRSID    string `json:"krsid"`
	KHSID    string `json:"khsid"`
}

func main() {
	// --- Input username & password ---
	config, err := LoadConfig()
	if err != nil {
		logf(LogError, "Gagal load konfigurasi: %v", err)
		return
	}

	// Create scraper instance
	scraper := NewScraper(config)

	// Handle authentication
	if err := handleAuthentication(scraper); err != nil {
		logf(LogError, "Gagal autentikasi: %v", err)
		return
	}

	// --- Ambil semester ---
	semester, err := scraper.SelectSemester()
	if err != nil {
		logf(LogError, "Gagal memilih semester: %v", err)
		return
	}
	logf(LogInfo, "Semester dipilih: %s", semester)

	// --- Load jurusan ---
	jurusanList, err := loadJurusan()
	if err != nil {
		logf(LogError, "Gagal load jurusan: %v", err)
		return
	}

	// --- Proses scraping per jurusan ---
	if err := processJurusan(scraper, jurusanList, semester); err != nil {
		logf(LogError, "Gagal proses jurusan: %v", err)
		return
	}

	log(LogInfo, "\nSemua data berhasil disimpan di folder nilai & excel")
}

// handleAuthentication handles the authentication process
func handleAuthentication(scraper *Scraper) error {
	// --- LOAD COOKIE ---
	if err := loadCookie(); err != nil {
		logf(LogError, "Gagal load cookie: %v", err)
		// Continue anyway, as we might still be able to login
	}
	scraper.cookie = cookie

	// --- LOGIN OTOMATIS ---
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
	logf(LogInfo, "Login Sebagai: %s", scraper.config.Username)
	return nil
}

// loadJurusan loads the jurusan data from file
func loadJurusan() ([]Jurusan, error) {
	cfgFile, err := os.ReadFile(JurusanFile)
	if err != nil {
		logf(LogError, "Gagal baca %s: %v", JurusanFile, err)
		return nil, fmt.Errorf("gagal baca %s: %w", JurusanFile, err)
	}
	var jurusanList []Jurusan
	if err := json.Unmarshal(cfgFile, &jurusanList); err != nil {
		logf(LogError, "Gagal parsing JSON: %v", err)
		return nil, fmt.Errorf("gagal parsing JSON: %w", err)
	}
	return jurusanList, nil
}

// processJurusan processes the scraping for each jurusan
func processJurusan(scraper *Scraper, jurusanList []Jurusan, semester string) error {
	for _, jur := range jurusanList {
		if jur.KodeJrs == "" {
			continue
		}
		logf(LogInfo, "\nMulai scraping jurusan: %s", jur.NamaJrs)

		if err := scraper.SetProdi(jur.KodeJrs, RegValue, semester); err != nil {
			logf(LogError, "Gagal set prodi: %v", err)
			continue
		}

		respMK, err := scraper.GetRekapMK()
		if err != nil {
			logf(LogError, "Gagal ambil mata kuliah: %v", err)
			continue
		}
		logf(LogInfo, "Total MK: %d", len(respMK.Rows))

		folderJSON := filepath.Join(JSONFolder, jur.NamaJrs, semester)
		folderExcel := filepath.Join(ExcelFolder, jur.NamaJrs, semester)
		os.MkdirAll(folderJSON, os.ModePerm)
		os.MkdirAll(folderExcel, os.ModePerm)

		var mkList []MataKuliah
		totals := len(respMK.Rows)
		for _, mk := range respMK.Rows {
			if mk.Cetak == "1" {
				mkList = append(mkList, mk)
			}
		}

		totalMK := len(mkList)
		doneMK := 0
		lastPercent := 0
		for _, mk := range mkList {
			nilaiMK, err := scraper.GetListNilai(mk.Infomk)
			if err != nil {
				logf(LogError, "Gagal ambil nilai untuk MK: %s - %v", mk.Namamk, err)
				continue
			}
			if mk.Cetak != "1" {
				continue
			}
			replacer := strings.NewReplacer("/", "-", ":", "")
			namaGabungan := fmt.Sprintf("%s %s %s", mk.Namamk, mk.Kelas, mk.Namadosen)
			namaFile := replacer.Replace(namaGabungan)

			if err := writeJSON(filepath.Join(folderJSON, namaFile+".json"), nilaiMK); err != nil {
				logf(LogError, "Gagal tulis JSON: %v", err)
			}
			if err := writeExcel(filepath.Join(folderExcel, namaFile+".xlsx"), nilaiMK); err != nil {
				logf(LogError, "Gagal tulis Excel: %v", err)
			}

			doneMK++
			updateProgress(jur.NamaJrs, doneMK, totalMK, &lastPercent)
		}
		fmt.Println()
		logf(LogInfo, "Jurusan %s: berhasil simpan %d MK dari %d MK, skip %d MK karena status cetak = 0", jur.NamaJrs, doneMK, totals, totals-doneMK)
	}
	return nil
}

// IsSessionValid checks if the current session is still valid
func (s *Scraper) IsSessionValid() bool {
	body, err := s.DoRequest(GET, MediaEndpoint, nil)
	if err != nil {
		logf(LogError, "Gagal cek session: %v", err)
		return false
	}

	if strings.Contains(string(body), "login") || strings.Contains(string(body), "Username") {
		return false
	}
	return true
}

// --- Load cookie dari file ---
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

// Login performs the login process to the SIAKAD system
func (s *Scraper) Login(username, password string) bool {
	// GET request to index.php
	res, err := s.client.Get(s.baseURL + IndexEndpoint)
	if err != nil {
		logf(LogError, "Gagal request index.php: %v", err)
		return false
	}
	defer res.Body.Close()

	// Ambil PHPSESSID awal
	initialSession := ""
	for _, c := range res.Cookies() {
		if c.Name == CookiePHPSESSID {
			initialSession = c.Value
		}
	}
	if initialSession == "" {
		log(LogError, "PHPSESSID awal tidak ditemukan")
		return false
	}

	// Generate hide_validation & IP random
	hideValidation := generateValidation()
	hideIP := getRandomIP()

	// POST login
	data := url.Values{}
	data.Set(FormUsername, username)
	data.Set(FormPassword, password)
	data.Set(FormValidation, hideValidation)
	data.Set(FormHideValidation, hideValidation)
	data.Set(FormHideIP, hideIP)

	req, err := http.NewRequest(POST, s.baseURL+LoginEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		logf(LogError, "Gagal membuat request login: %v", err)
		return false
	}

	// Set headers for login request
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderContentType, ContentTypeForm)
	req.Header.Set(HeaderCookie, CookiePHPSESSID+"="+initialSession)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)
	req.Header.Set(HeaderReferer, s.baseURL+IndexEndpoint)

	resp, err := s.client.Do(req)
	if err != nil {
		logf(LogError, "Gagal POST login: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Ambil cookie baru dari response login
	for _, c := range resp.Cookies() {
		if c.Name == CookiePHPSESSID {
			s.cookie = CookiePHPSESSID + "=" + c.Value
			cookie = s.cookie // Update global cookie for backward compatibility
			log(LogInfo, "Login berhasil!")
			break
		}
	}

	if err := saveCookie(); err != nil {
		logf(LogWarn, "Gagal simpan cookie: %v", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logf(LogWarn, "Gagal baca response body: %v", err)
		// Continue anyway, as login might still be successful
	}

	// logf(LogDebug, "login response: %s", string(respBody))
	// logf(LogDebug, "cookie: %s", cookie)

	return strings.Contains(string(respBody), `"success":true`)
}

// --- GENERATE HIDE_VALIDATION ---
func generateValidation() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// Use a more secure random number generator
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

// Semester represents a semester option
type Semester struct {
	Keterangan string `json:"keterangan"`
	Smtthnakd  string `json:"smtthnakd"`
}

// SelectSemester allows the user to select a semester from a list
func (s *Scraper) SelectSemester() (string, error) {
	body, err := s.DoRequest(POST, "/_modul/aksi_umum.php?act=pilih_smtthnakd", nil)
	if err != nil {
		return "", fmt.Errorf("gagal ambil semester: %w", err)
	}

	var semesters []Semester
	if err := json.Unmarshal(body, &semesters); err != nil {
		return "", fmt.Errorf("gagal parsing JSON semester: %w", err)
	}

	if len(semesters) == 0 {
		return "", fmt.Errorf("tidak ada semester yang tersedia")
	}

	log(LogInfo, "Daftar Semester:")
	for i, semester := range semesters {
		logf(LogInfo, "[%d] %s", i+1, semester.Keterangan)
	}

	var selection int
	log(LogInfo, "Pilih semester (nomor): ")
	_, err = fmt.Scan(&selection)
	if err != nil {
		return "", fmt.Errorf("gagal membaca input: %w", err)
	}
	if selection < 1 || selection > len(semesters) {
		return "", fmt.Errorf("pilihan invalid: %d", selection)
	}

	return semesters[selection-1].Smtthnakd, nil
}

// SetProdi sets the program study (prodi) for a given semester
func (s *Scraper) SetProdi(kodeProdi, kodePK, smthn string) error {
	form := url.Values{}
	form.Set(FormPS, kodeProdi)
	form.Set(FormPK, kodePK)
	form.Set(FormSMTHN, smthn)

	_, err := s.DoRequest(POST, "/_modul/mod_prodi_smthn/aksi_prodi_smthn.php", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("gagal set prodi: %w", err)
	}

	return nil
}

// RekapMKResponse represents the response from getRekapMK
type RekapMKResponse struct {
	Total int          `json:"total"`
	Rows  []MataKuliah `json:"rows"`
}

// GetRekapMK retrieves the recap of courses (mata kuliah)
func (s *Scraper) GetRekapMK() (*RekapMKResponse, error) {
	data := "page=1&rows=300&sort=hari&order=asc"
	body, err := s.DoRequest(POST, "/_modul/mod_nilmk/aksi_nilmk.php?act=rekapNILMK", strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gagal ambil rekap MK: %w", err)
	}

	var resp RekapMKResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON: %w", err)
	}
	return &resp, nil
}

// GetListNilai retrieves the list of grades (nilai) for a given course
func (s *Scraper) GetListNilai(infomk string) ([]Nilai, error) {
	form := url.Values{}
	form.Set(FormParam, infomk)
	form.Set(FormCetak, CetakValue)

	body, err := s.DoRequest(POST, "/_modul/mod_nilmk/aksi_nilmk.php?act=listNILMK", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("gagal ambil list nilai: %w", err)
	}

	var hasil []Nilai
	if err := json.Unmarshal(body, &hasil); err != nil {
		return nil, fmt.Errorf("gagal parsing JSON listNILMK: %w", err)
	}
	return hasil, nil
}

func writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		logf(LogError, "Gagal buat file %s: %v", path, err)
		return fmt.Errorf("gagal buat file %s: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		logf(LogError, "Gagal tulis JSON ke %s: %v", path, err)
		return fmt.Errorf("gagal tulis JSON ke %s: %w", path, err)
	}
	return nil
}

func writeExcel(path string, data []Nilai) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{ExcelNIM, ExcelNama, ExcelAngka, ExcelHuruf, ExcelKehadiran, ExcelProjek, ExcelQuiz, ExcelTugas, ExcelUTS, ExcelUAS}
	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			logf(LogError, "Gagal membuat koordinat cell: %v", err)
			return fmt.Errorf("gagal membuat koordinat cell: %w", err)
		}
		f.SetCellValue(sheet, cell, h)
	}

	for i, n := range data {
		row := i + 2
		values := []interface{}{n.NIM, n.Nama, n.NilAngka, n.NilHuruf, n.Hadir, n.Projek, n.Quiz, n.Tugas, n.UTS, n.UAS}
		for j, v := range values {
			cell, err := excelize.CoordinatesToCellName(j+1, row)
			if err != nil {
				logf(LogError, "Gagal membuat koordinat cell: %v", err)
				return fmt.Errorf("gagal membuat koordinat cell: %w", err)
			}
			f.SetCellValue(sheet, cell, v)
		}
	}

	if err := f.SaveAs(path); err != nil {
		logf(LogError, "Gagal simpan Excel %s: %v", path, err)
		return fmt.Errorf("gagal simpan Excel %s: %w", path, err)
	}
	return nil
}

func setHeaders(req *http.Request) {
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set(HeaderUserAgent, ua)
	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)
	req.Header.Set(HeaderReferer, baseURL+MediaEndpoint)
	req.Header.Set(HeaderOrigin, baseURL)
	req.Header.Set(HeaderAccept, AcceptJSON)
	req.Header.Set(HeaderContentType, ContentTypeForm+CharsetUTF8)
	req.Header.Set(HeaderCookie, cookie)
}

// updateProgress updates the progress bar display
func updateProgress(jurusan string, done, total int, lastPercent *int) {
	percent := done * 100 / total
	if percent-*lastPercent >= 5 || percent == 100 {
		barLen := ProgressBarLength
		pos := percent * barLen / 100
		bar := strings.Repeat("=", pos) + strings.Repeat(" ", barLen-pos)
		fmt.Printf("\r[PROGRESS] Mata kuliah %s: [%s] %d%% (%d/%d)", jurusan, bar, percent, done, total)
		*lastPercent = percent
	}
}
