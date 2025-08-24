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

var (
	baseURL string
	cookie  string
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

	LoadConfig()

	username := USERNAME
	password := PASSWORD
	baseURL = BaseURL
	// --- LOAD COOKIE ---
	loadCookie()

	// --- LOGIN OTOMATIS ---
	if cookie != "" && isSessionValid() {
		fmt.Println("[INFO] Cookie masih valid, skip login")
	} else {
		fmt.Println("[INFO] Login ulang...")
		if !login(username, password) {
			fmt.Println("[ERROR] Login gagal, cek username/password")
			return
		}
		saveCookie()
	}

	fmt.Println("[INFO] Login Sebagai:", username)

	// --- Ambil semester ---
	semester := pilihSemester()
	fmt.Println("[INFO] Semester dipilih:", semester)

	// --- Load jurusan ---
	cfgFile, err := os.ReadFile("jurusan.json")
	if err != nil {
		fmt.Println("[ERROR] Gagal baca jurusan.json:", err)
		return
	}
	var jurusanList []Jurusan
	if err := json.Unmarshal(cfgFile, &jurusanList); err != nil {
		fmt.Println("[ERROR] Gagal parsing JSON:", err)
		return
	}

	// --- Proses scraping per jurusan ---
	for _, jur := range jurusanList {
		if jur.KodeJrs == "" {
			continue
		}
		fmt.Printf("\n[INFO] Mulai scraping jurusan: %s\n", jur.NamaJrs)

		if err := setProdi(jur.KodeJrs, "REG", semester); err != nil {
			fmt.Println("[ERROR] Gagal set prodi:", err)
			continue
		}

		respMK := getRekapMK()
		if respMK == nil {
			fmt.Println("[ERROR] Gagal ambil mata kuliah")
			continue
		}
		fmt.Printf("Total MK: %d\n", len(respMK.Rows))

		folderJSON := filepath.Join("nilai_json", jur.NamaJrs, semester)
		folderExcel := filepath.Join("nilai_excel", jur.NamaJrs, semester)
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
			nilaiMK := getListNilai(mk.Infomk)
			if mk.Cetak != "1" {
				continue
			}
			replacer := strings.NewReplacer("/", "-", ":", "")
			namaGabungan := fmt.Sprintf("%s %s %s", mk.Namamk, mk.Kelas, mk.Namadosen)
			namaFile := replacer.Replace(namaGabungan)

			writeJSON(filepath.Join(folderJSON, namaFile+".json"), nilaiMK)
			writeExcel(filepath.Join(folderExcel, namaFile+".xlsx"), nilaiMK)

			doneMK++
			percent := doneMK * 100 / totalMK
			if percent-lastPercent >= 5 || percent == 100 {
				barLen := 40
				pos := percent * barLen / 100
				bar := strings.Repeat("=", pos) + strings.Repeat(" ", barLen-pos)
				fmt.Printf("\r[PROGRESS] Mata kuliah %s: [%s] %d%% (%d/%d)", jur.NamaJrs, bar, percent, doneMK, totalMK)
				lastPercent = percent
			}
		}
		fmt.Println()
		fmt.Printf("[INFO] Jurusan %s: berhasil simpan %d MK dari %d MK, skip %d MK karena status cetak = 0\n", jur.NamaJrs, doneMK, totals, totals-doneMK)
	}

	fmt.Println("\n[INFO] Semua data berhasil disimpan di folder nilai & excel")
}
func isSessionValid() bool {
	req, _ := http.NewRequest("GET", baseURL+"/media.php", nil)
	setHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Gagal cek session:", err)
		return false
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if strings.Contains(string(body), "login") || strings.Contains(string(body), "Username") {
		return false
	}
	return true
}

// --- Load cookie dari file ---
func loadCookie() {
	data, err := os.ReadFile("cookie.txt")
	if err == nil {
		cookie = string(data)
		fmt.Println("[INFO] Cookie ditemukan!")
	}
}

// --- Simpan cookie ke file ---
func saveCookie() {
	if cookie != "" {
		_ = os.WriteFile("cookie.txt", []byte(cookie), 0644)
	}
}

// --- LOGIN OTOMATIS (ambil PHPSESSID) ---
func login(username, password string) bool {
	req, _ := http.NewRequest("GET", baseURL+"/index.php", nil)
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Gagal request index.php:", err)
		return false
	}
	defer res.Body.Close()

	// Ambil PHPSESSID awal
	initialSession := ""
	for _, c := range res.Cookies() {
		if c.Name == "PHPSESSID" {
			initialSession = c.Value
		}
	}
	if initialSession == "" {
		fmt.Println("[ERROR] PHPSESSID awal tidak ditemukan")
		return false
	}

	// Generate hide_validation & IP random
	hideValidation := generateValidation()
	hideIP := getRandomIP()

	// POST login
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("validation", hideValidation)
	data.Set("hide_validation", hideValidation)
	data.Set("hide_ipnya", hideIP)

	req2, _ := http.NewRequest("POST", baseURL+"/ceklogin.php?h=", strings.NewReader(data.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Set("Cookie", "PHPSESSID="+initialSession)
	req2.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req2.Header.Set("X-Requested-With", "XMLHttpRequest")
	req2.Header.Set("Referer", baseURL+"/index.php")

	res2, err := client.Do(req2)
	if err != nil {
		fmt.Println("[ERROR] Gagal POST login:", err)
		return false
	}
	defer res2.Body.Close()

	// Ambil cookie baru dari response login
	for _, c := range res2.Cookies() {
		if c.Name == "PHPSESSID" {
			cookie = "PHPSESSID=" + c.Value
			fmt.Println("[INFO] Login berhasil!")
			break
		}
	}

	saveCookie()
	respBody, _ := io.ReadAll(res2.Body)
	// fmt.Println("[DEBUG] login response:", string(respBody))
	// fmt.Println("[DEBUG] cookie:", cookie)

	return strings.Contains(string(respBody), `"success":true`)
}

// --- GENERATE HIDE_VALIDATION ---
func generateValidation() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 3)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

// --- AMBIL IP PUBLIK RANDOM ---
func getRandomIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println("[WARN] Gagal ambil IP publik, pakai default")
		return "182.8.179.9"
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(ip))
}

// --- Fungsi baru untuk ambil semester ---
func pilihSemester() string {
	req, _ := http.NewRequest("POST", baseURL+"/_modul/aksi_umum.php?act=pilih_smtthnakd", nil)
	setHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Gagal ambil semester:", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var semList []struct {
		Keterangan string `json:"keterangan"`
		Smtthnakd  string `json:"smtthnakd"`
	}
	if err := json.Unmarshal(body, &semList); err != nil {
		fmt.Println("[ERROR] Gagal parsing JSON semester:", err)
		os.Exit(1)
	}

	fmt.Println("Daftar Semester:")
	for i, s := range semList {
		fmt.Printf("[%d] %s\n", i+1, s.Keterangan)
	}

	var pilih int
	fmt.Print("Pilih semester (nomor): ")
	fmt.Scan(&pilih)
	if pilih < 1 || pilih > len(semList) {
		fmt.Println("[ERROR] Pilihan invalid")
		os.Exit(1)
	}

	return semList[pilih-1].Smtthnakd
}

// --- Set prodi ---
func setProdi(kodeProdi, kodePK, smthn string) error {
	form := url.Values{}
	form.Set("ps", kodeProdi)
	form.Set("pk", kodePK)
	form.Set("smthn", smthn)

	req, _ := http.NewRequest("POST", baseURL+"/_modul/mod_prodi_smthn/aksi_prodi_smthn.php", strings.NewReader(form.Encode()))
	setHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	io.ReadAll(res.Body)
	return nil
}

// --- Fungsi lain tetap sama ---
func getRekapMK() *struct {
	Total int          `json:"total"`
	Rows  []MataKuliah `json:"rows"`
} {
	data := "page=1&rows=300&sort=hari&order=asc"
	req, _ := http.NewRequest("POST", baseURL+"/_modul/mod_nilmk/aksi_nilmk.php?act=rekapNILMK", strings.NewReader(data))
	setHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var resp struct {
		Total int          `json:"total"`
		Rows  []MataKuliah `json:"rows"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("[ERROR] Gagal parsing JSON:", err)
		return nil
	}
	return &resp
}

func getListNilai(infomk string) []Nilai {
	form := url.Values{}
	form.Set("param", infomk)
	form.Set("cetak", "1")

	req, _ := http.NewRequest("POST", baseURL+"/_modul/mod_nilmk/aksi_nilmk.php?act=listNILMK", strings.NewReader(form.Encode()))
	setHeaders(req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var hasil []Nilai
	if err := json.Unmarshal(body, &hasil); err != nil {
		fmt.Println("[ERROR] Gagal parsing JSON listNILMK:", err)
		return nil
	}
	return hasil
}

func writeJSON(path string, data interface{}) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("[ERROR] Gagal buat file:", path, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Println("[ERROR] Gagal tulis JSON:", err)
	}
}

func writeExcel(path string, data []Nilai) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{"NIM", "Nama Peserta", "Angka", "Huruf", "Kehadiran", "Projek", "Quiz", "Tugas", "UTS", "UAS"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, n := range data {
		row := i + 2
		values := []interface{}{n.NIM, n.Nama, n.NilAngka, n.NilHuruf, n.Hadir, n.Projek, n.Quiz, n.Tugas, n.UTS, n.UAS}
		for j, v := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	if err := f.SaveAs(path); err != nil {
		fmt.Println("[ERROR] Gagal simpan Excel:", path, err)
	}
}

func setHeaders(req *http.Request) {
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", ua)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", baseURL+"/media.php")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", cookie)
}
