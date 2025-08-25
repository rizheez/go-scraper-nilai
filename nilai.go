package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	namaFile := sanitizeFilename(fmt.Sprintf("%s R%s %s", mk.Namamk, mk.Kelas, mk.Namadosen))

	if err := writeJSON(filepath.Join(folderJSON, namaFile+".json"), nilai); err != nil {
		logf(LogError, "Gagal tulis JSON: %v", err)
	}
	if err := writeExcel(filepath.Join(folderExcel, namaFile+".xlsx"), nilai); err != nil {
		logf(LogError, "Gagal tulis Excel: %v", err)
	}
}

func writeJSON(path string, data interface{}) error {
	file, _ := os.Create(path)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func writeExcel(path string, data []Nilai) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	headers := []string{ExcelNIM, ExcelNama, ExcelAngka, ExcelHuruf, ExcelKehadiran, ExcelProjek, ExcelQuiz, ExcelTugas, ExcelUTS, ExcelUAS}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, n := range data {
		row := i + 2
		vals := []interface{}{n.NIM, n.Nama, n.NilAngka, n.NilHuruf, n.Hadir, n.Projek, n.Quiz, n.Tugas, n.UTS, n.UAS}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}
	return f.SaveAs(path)
}
