package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func processMHS(scraper *Scraper, jur Jurusan, semester string) error {
	// Set prodi sesuai jurusan dan semester
	if err := scraper.SetProdi(jur.KodeJrs, RegValue, semester); err != nil {
		return fmt.Errorf("gagal set prodi untuk jurusan %s: %w", jur.NamaJrs, err)
	}

	// Ambil data rekap mahasiswa
	resp, err := scraper.GetRekapMHS()
	if err != nil {
		return fmt.Errorf("gagal ambil rekap mahasiswa jurusan %s: %w", jur.NamaJrs, err)
	}

	// Konversi response ke slice Mahasiswa
	mhsList := append([]Mahasiswa(nil), resp.Rows...)

	// Kalau tidak ada mahasiswa, langsung return
	if len(mhsList) == 0 {
		logf(LogWarn, "Jurusan %s tidak ada data mahasiswa untuk semester %s", jur.NamaJrs, semester)
		return nil
	}

	// Filter mahasiswa berdasarkan tahun masuk (optional)
	filteredMhsList, tahun, err := filterMahasiswaByYear(mhsList)
	if err != nil {
		return fmt.Errorf("gagal filter mahasiswa: %w", err)
	}

	// Kalau setelah filter tidak ada mahasiswa
	if len(filteredMhsList) == 0 {
		logf(LogWarn, "Jurusan %s tidak ada data mahasiswa setelah filter untuk semester %s", jur.NamaJrs, semester)
		return nil
	}

	// Siapkan folder penyimpanan JSON & Excel
	folderJSON := filepath.Join(JSONFolder, jur.NamaJrs, "Mahasiswa")
	folderExcel := filepath.Join(ExcelFolder, jur.NamaJrs, "Mahasiswa")
	if err := os.MkdirAll(folderJSON, os.ModePerm); err != nil {
		return fmt.Errorf("gagal buat folder JSON: %w", err)
	}
	if err := os.MkdirAll(folderExcel, os.ModePerm); err != nil {
		return fmt.Errorf("gagal buat folder Excel: %w", err)
	}

	// Logging header
	printHeader(fmt.Sprintf("Scraping Mahasiswa Jurusan %s - Semester %s", jur.NamaJrs, semester), nil)
	logf(LogInfo, "Jurusan %s: berhasil ambil %d mahasiswa (dari %d total)", jur.NamaJrs, len(filteredMhsList), len(mhsList))

	// Write JSON file
	namaFile := sanitizeFilename("Mahasiswa")
	jsonPath := filepath.Join(folderJSON, namaFile+" "+tahun+".json")
	if err := writeJSONMHS(jsonPath, filteredMhsList); err != nil {
		return fmt.Errorf("gagal tulis JSON untuk jurusan %s: %w", jur.NamaJrs, err)
	}
	// Write Excel file
	excelPath := filepath.Join(folderExcel, namaFile+" "+tahun+".xlsx")
	if err := writeExcelMHS(excelPath, filteredMhsList); err != nil {
		return fmt.Errorf("gagal tulis Excel untuk jurusan %s: %w", jur.NamaJrs, err)
	}

	logf(LogInfo, "Berhasil simpan data mahasiswa ke: %s", jsonPath)
	logf(LogInfo, "Berhasil simpan data mahasiswa ke: %s", excelPath)

	return nil
}

func scrapeMHS(MHS, folderJSON, folderExcel string) {
	namaFile := sanitizeFilename(fmt.Sprintf("Mahasiswa"))

	if err := writeJSONMHS(filepath.Join(folderJSON, namaFile+".json"), MHS); err != nil {
		logf(LogError, "Gagal tulis JSON: %v", err)
	}
	// if err := writeExcelMHS(filepath.Join(folderExcel, namaFile+".xlsx"), MHS); err != nil {
	// 	logf(LogError, "Gagal tulis Excel: %v", err)
	// }
}

func writeJSONMHS(path string, data interface{}) error {
	file, _ := os.Create(path)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeExcelMHS(path string, data []Mahasiswa) error {
	xlsx := excelize.NewFile()
	sheet := "Sheet1"

	// Define headers for student data
	headers := []string{
		"NIM", "Nama", "Jurusan",
		"NIK", "Program Studi", "Angkatan",
		"Tempat Lahir", "Tanggal Lahir", "Gender", "Agama",
		"Email", "HP", "Alamat", "Status", "IPK", "SKS Total",
	}

	// Set headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		xlsx.SetCellValue(sheet, cell, h)
	}

	// Set data rows
	for i, mhs := range data {
		row := i + 2
		vals := []interface{}{
			mhs.NIM,
			mhs.Nama,
			mhs.NamaJrs,
			mhs.NoKTP,
			mhs.NamaPK,
			mhs.TanggalMasuk[:4], // Extract year from date
			mhs.TempatLahir,
			mhs.TanggalLahir,
			mhs.Gender,
			mhs.NamaAgama,
			mhs.Email,
			mhs.HP1,
			mhs.AlamatSurat1,
			mhs.StatusMhsKet,
			mhs.IPK,
			mhs.SKSTotal,
		}

		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			xlsx.SetCellValue(sheet, cell, v)
		}
	}

	return xlsx.SaveAs(path)
}

// filterMahasiswaByYear filters mahasiswa based on enrollment year
func filterMahasiswaByYear(mhsList []Mahasiswa) ([]Mahasiswa, string, error) {
	// Tampilkan opsi filter
	fmt.Println()
	fmt.Println("=================================")
	log(LogInfo, "Filter Mahasiswa berdasarkan Tahun Masuk:")
	fmt.Println("=================================")
	logf(LogInfo, "[1] Semua Tahun (Tanpa Filter)")
	logf(LogInfo, "[2] Filter berdasarkan Tahun Tertentu")
	fmt.Println("=================================")

	var pilihan int
	fmt.Printf("[INFO] Pilih opsi filter (1-2): ")
	_, err := fmt.Scan(&pilihan)
	if err != nil {
		return nil, "Semua Tahun", fmt.Errorf("gagal membaca input filter: %w", err)
	}

	switch pilihan {
	case 1:
		// Tidak ada filter, return semua data
		logf(LogInfo, "Filter: Mengambil semua data mahasiswa")
		return mhsList, "Semua Tahun", nil

	case 2:
		// Filter berdasarkan tahun
		var tahunFilter string
		fmt.Printf("[INFO] Masukkan tahun masuk (contoh: 2025, 2024, 2023): ")
		_, err := fmt.Scan(&tahunFilter)
		if err != nil {
			return nil, string(tahunFilter), fmt.Errorf("gagal membaca input tahun: %w", err)
		}

		// Validasi tahun
		if _, err := strconv.Atoi(tahunFilter); err != nil {
			return nil, string(tahunFilter), fmt.Errorf("tahun tidak valid: %s", tahunFilter)
		}

		logf(LogInfo, "Filter: Mengambil data mahasiswa tahun %s", tahunFilter)

		// Filter mahasiswa berdasarkan tahun masuk
		var filteredList []Mahasiswa
		for _, mhs := range mhsList {
			// Parse tanggal masuk (format: "2025-09-01")
			if strings.HasPrefix(mhs.TanggalMasuk, tahunFilter) {
				filteredList = append(filteredList, mhs)
			}
		}

		logf(LogInfo, "Ditemukan %d mahasiswa dari tahun %s (dari %d total)", len(filteredList), tahunFilter, len(mhsList))
		return filteredList, string(tahunFilter), nil

	default:
		return nil, "Semua Tahun", fmt.Errorf("pilihan filter tidak valid: %d", pilihan)
	}
}
