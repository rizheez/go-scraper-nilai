package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func (s *Scraper) SelectSemester() (string, error) {
	body, err := s.DoRequest(POST, "/_modul/aksi_umum.php?act=pilih_smtthnakd", nil)
	if err != nil {
		return "", err
	}
	var semesters []Semester
	if err := json.Unmarshal(body, &semesters); err != nil {
		return "", err
	}
	if len(semesters) == 0 {
		return "", fmt.Errorf("tidak ada semester")
	}

	printHeader("Daftar Semester", nil)
	for i, sm := range semesters {
		logf(LogInfo, "[%d] %s", i+1, sm.Keterangan)
	}

	var sel int
	fmt.Print("[INFO] Pilih semester (nomor): ")
	fmt.Scan(&sel)
	return semesters[sel-1].Smtthnakd, nil
}

func SelectJurusan() (Jurusan, error) {
	data, err := os.ReadFile(JurusanFile)
	if err != nil {
		return Jurusan{}, err
	}
	var jurusanList []Jurusan
	if err := json.Unmarshal(data, &jurusanList); err != nil {
		return Jurusan{}, err
	}
	if len(jurusanList) == 0 {
		return Jurusan{}, fmt.Errorf("jurusan kosong")
	}

	printHeader("Daftar Jurusan", nil)
	for i, j := range jurusanList {
		logf(LogInfo, "[%d] %s", i+1, j.NamaJrs)
	}

	var sel int
	fmt.Print("[INFO] Pilih jurusan (nomor): ")
	fmt.Scan(&sel)
	return jurusanList[sel-1], nil
}
