package main

import (
	"fmt"
	"time"
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
	FakValue   = "1" // Default faculty value

	// Excel headers
	ExcelNIM       = "Nim"
	ExcelNama      = "Nama Mahasiswa"
	ExcelKodeMK    = "Kode Mata Kuliah"
	ExcelNamaMK    = "Nama Mata Kuliah"
	ExcelSemester  = "Semester"
	ExcelKelas     = "Nama Kelas"
	ExcelAngka     = "Angka"
	ExcelHuruf     = "Huruf"
	ExcelKehadiran = "Aktivitas Partisipatif"
	ExcelProjek    = "Hasil Proyek"
	ExcelQuiz      = "Kognitif/ Pengetahuan Quiz"
	ExcelTugas     = "Kognitif/ Pengetahuan Tugas"
	ExcelUTS       = "Kognitif/ Pengetahuan Ujian Tengah Semester"
	ExcelUAS       = "Kognitif/ Pengetahuan Ujian Akhir Semester"
	ExcelKodePK    = "Kode Prodi Mahasiswa"
	ExcelNamaPK    = "Nama Prodi Mahasiswa"
	ExcelKodeProdi = "Kode Prodi Kelas"
	ExcelNamaProdi = "Nama Prodi Kelas"

	// Progress bar
	ProgressBarLength = 40

	// Log levels
	LogInfo    = "[INFO]"
	LogError   = "[ERROR]"
	LogWarn    = "[WARN]"
	LogDebug   = "[DEBUG]"
	LogWelcome = "[WELCOME]"
)

var (
	// baseURL string
	cookie string
)

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

func main() {
	// --- Input username & password ---
	config, err := LoadConfig()

	clearScreen()
	fmt.Println("=================================")
	logf(LogWelcome, "Scraper Nilai Akademik")
	fmt.Println("=================================")
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
	logf(LogInfo, "Login Sebagai: %s", scraper.config.Username)
	fmt.Println()

	// --- Ambil semester ---
	semester, err := scraper.SelectSemester()
	if err != nil {
		logf(LogError, "Gagal memilih semester: %v", err)
		return
	}
	logf(LogInfo, "Semester dipilih: %s", semester)
	fmt.Println()
	// --- Load jurusan ---
	jurusan, err := loadJurusan()
	if err != nil {
		logf(LogError, "Gagal load jurusan: %v", err)
		return
	}
	fmt.Println()

	// --- Pilih jenis scraping ---
	fmt.Println("=================================")
	log(LogInfo, "Pilih jenis scraping:")
	fmt.Println("=================================")
	logf(LogInfo, "[1] Process Jurusan (Scrape Nilai Mata Kuliah)")
	logf(LogInfo, "[2] Process Mahasiswa (Scrape Data Mahasiswa)")
	logf(LogInfo, "[3] Process Keduanya")
	fmt.Println("=================================")

	var pilihan int
	fmt.Printf("[INFO] Pilih opsi (1-3): ")
	_, err = fmt.Scan(&pilihan)
	if err != nil {
		logf(LogError, "Gagal membaca input: %v", err)
		return
	}

	if pilihan < 1 || pilihan > 3 {
		logf(LogError, "Pilihan invalid: %d. Pilih antara 1-3", pilihan)
		return
	}

	fmt.Println()
	start := time.Now()
	// --- Proses scraping sesuai pilihan ---
	switch pilihan {
	case 1:
		logf(LogInfo, "Memulai scraping Nilai Mata Kuliah...")
		if err := processJurusan(scraper, jurusan, semester); err != nil {
			logf(LogError, "Gagal proses jurusan: %v", err)
			return
		}
	case 2:
		logf(LogInfo, "Memulai scraping Data Mahasiswa...")
		if err := processMHS(scraper, jurusan, semester); err != nil {
			logf(LogError, "Gagal proses Mahasiswa: %v", err)
			return
		}
	case 3:
		logf(LogInfo, "Memulai scraping Keduanya...")
		if err := processJurusan(scraper, jurusan, semester); err != nil {
			logf(LogError, "Gagal proses jurusan: %v", err)
			return
		}
		if err := processMHS(scraper, jurusan, semester); err != nil {
			logf(LogError, "Gagal proses Mahasiswa: %v", err)
			return
		}
	}
	elapsed := time.Since(start)
	fmt.Println()
	fmt.Println("=====================================================================")
	log(LogInfo, "Semua data berhasil disimpan di folder")
	logf(LogInfo, "Waktu yang dibutuhkan: %s", formatDuration(elapsed))
	fmt.Println("=====================================================================")
}

// func setHeaders(req *http.Request) {
// 	ua := userAgents[rand.Intn(len(userAgents))]
// 	req.Header.Set(HeaderUserAgent, ua)
// 	req.Header.Set(HeaderXRequestedWith, XMLHttpRequest)
// 	req.Header.Set(HeaderReferer, baseURL+MediaEndpoint)
// 	req.Header.Set(HeaderOrigin, baseURL)
// 	req.Header.Set(HeaderAccept, AcceptJSON)
// 	req.Header.Set(HeaderContentType, ContentTypeForm+CharsetUTF8)
// 	req.Header.Set(HeaderCookie, cookie)
// }
// no changes
