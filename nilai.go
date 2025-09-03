package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

const WorkerCount = 5

func processJurusan(scraper *Scraper, jur Jurusan, semester string) error {
	if err := scraper.SetProdi(jur.KodeJrs, RegValue, semester); err != nil {
		return err
	}

	resp, err := scraper.GetRekapMK()
	if err != nil {
		return err
	}

	// filter MK cetak=1
	var skip int
	var mkList []MataKuliah
	for _, mk := range resp.Rows {
		if mk.Cetak == "1" {
			mkList = append(mkList, mk)
		} else {
			skip++
		}
	}

	total := len(mkList)
	if total == 0 {
		logf(LogWarn, "Jurusan %s tidak ada MK dengan cetak=1", jur.NamaJrs)
		return nil
	}
	all := len(resp.Rows)
	folderJSON := filepath.Join(JSONFolder, jur.NamaJrs, semester)
	folderExcel := filepath.Join(ExcelFolder, jur.NamaJrs, semester)
	os.MkdirAll(folderJSON, os.ModePerm)
	os.MkdirAll(folderExcel, os.ModePerm)

	var wg sync.WaitGroup
	done := 0
	last := 0
	mu := sync.Mutex{}
	printHeader("Scraping Jurusan", nil)
	logf("[SCRAPING]", "Mulai scraping jurusan: %s", jur.NamaJrs)
	for _, mk := range mkList {
		wg.Add(1)
		go func(mk MataKuliah) {
			defer wg.Done()
			scrapeMK(scraper, mk, folderJSON, folderExcel)

			mu.Lock()
			done++
			updateProgress(jur.NamaJrs, done, total, &last)
			mu.Unlock()
		}(mk)
	}

	wg.Wait()
	fmt.Println()
	logf(LogInfo, "Jurusan %s: berhasil simpan %d MK dari %d MK, skip %d MK karena status cetak = 0", jur.NamaJrs, done, all, skip)
	return nil
}

func scrapeMK(scraper *Scraper, mk MataKuliah, folderJSON, folderExcel string) {
	nilai, err := scraper.GetListNilai(mk.Infomk)
	if err != nil {
		logf(LogError, "Gagal ambil nilai MK %s: %v", mk.Namamk, err)
		return
	}
	infomk := strings.Split(mk.Infomk, "#")
	fak := infomk[0]
	// Get bobot data
	bobotData, err := scraper.GetBobotMK(fak, mk.KodeJrs, mk.KodePK, mk.Kelas, mk.KodeMK)
	if err != nil {
		logf(LogWarn, "Gagal ambil bobot MK %s: %v", mk.Namamk, err)
		// Continue with empty bobot data
		bobotData = Bobot{}
	}

	namaFile := sanitizeFilename(fmt.Sprintf("%s R%s %s", mk.Namamk, mk.Kelas, mk.Namadosen))

	// Write nilai data
	if err := writeJSON(filepath.Join(folderJSON, namaFile+".json"), nilai); err != nil {
		logf(LogError, "Gagal tulis JSON nilai: %v", err)
	}
	if err := writeExcel(filepath.Join(folderExcel, namaFile+".xlsx"), nilai, mk); err != nil {
		logf(LogError, "Gagal tulis Excel nilai: %v", err)
	}

	// Write bobot data
	bobotMK := BobotMK{
		MataKuliah: mk,
		Bobot:      bobotData,
	}
	namaFileBobot := sanitizeFilename(fmt.Sprintf("%s R%s %s_bobot", mk.Namamk, mk.Kelas, mk.Namadosen))
	if err := writeJSON(filepath.Join(folderJSON, namaFileBobot+".json"), bobotMK); err != nil {
		logf(LogError, "Gagal tulis JSON bobot: %v", err)
	}
	if err := writeBobotExcel(filepath.Join(folderExcel, namaFileBobot+".xlsx"), bobotMK); err != nil {
		logf(LogError, "Gagal tulis Excel bobot: %v", err)
	}
}

func writeJSON(path string, data interface{}) error {
	file, _ := os.Create(path)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeExcel(path string, data []Nilai, mk MataKuliah) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{ExcelNIM, ExcelNama, ExcelKodeMK, ExcelNamaMK, ExcelSemester, ExcelKelas, ExcelAngka, ExcelHuruf, ExcelKehadiran, ExcelProjek, ExcelQuiz, ExcelTugas, ExcelUTS, ExcelUAS, ExcelKodePK, ExcelNamaPK, ExcelKodeProdi, ExcelNamaProdi}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, n := range data {
		row := i + 2
		vals := []interface{}{n.NIM, n.Nama, mk.KodeMK, mk.Namamk, mk.Smtthnakd, mk.Kelas, n.NilAngka, n.NilHuruf, n.Hadir, n.Projek, n.Quiz, n.Tugas, n.UTS, n.UAS, mk.KodeJrs, mk.NamaJrs, mk.KodeJrs, mk.NamaJrs}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}
	return f.SaveAs(path)
}

func writeBobotExcel(path string, data BobotMK) error {
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Set mata kuliah information headers
	headers := []string{"Mata Kuliah", "Kelas", "Dosen", "Kode MK", "Kode Prodi", "Kode PK"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Set mata kuliah information values
	mkVals := []interface{}{data.MataKuliah.Namamk, data.MataKuliah.Kelas, data.MataKuliah.Namadosen, data.MataKuliah.KodeMK, data.MataKuliah.KodeJrs, data.MataKuliah.KodePK}
	for i, v := range mkVals {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, v)
	}

	// Add bobot headers starting from row 4
	f.SetCellValue(sheet, "A4", "Bobot (%)")

	// Add bobot data starting from row 5
	bobotComponents := []string{"Hadir", "Projek", "Quiz", "Tugas", "UTS", "UAS"}
	bobotValues := []string{data.Bobot.Hadir, data.Bobot.Projek, data.Bobot.Quiz, data.Bobot.Tugas, data.Bobot.UTS, data.Bobot.UAS}
	for i, component := range bobotComponents {

		// Set component name in column A
		cell, _ := excelize.CoordinatesToCellName(i+1, 5)
		f.SetCellValue(sheet, cell, component)
		// Set bobot value in column B
	}
	for i, value := range bobotValues {

		cell, _ := excelize.CoordinatesToCellName(i+1, 6)
		f.SetCellValue(sheet, cell, value)
	}

	return f.SaveAs(path)
}

// no changes
