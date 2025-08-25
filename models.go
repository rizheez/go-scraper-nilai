package main

type Jurusan struct {
	JrsID   string `json:"jrsid"`
	KodeJrs string `json:"kodejrs"`
	NamaJrs string `json:"namajrs"`
}

type Semester struct {
	Keterangan string `json:"keterangan"`
	Smtthnakd  string `json:"smtthnakd"`
}

type MataKuliah struct {
	JID       string `json:"jid"`
	Namamk    string `json:"namamk"`
	Kelas     string `json:"kelas"`
	Namadosen string `json:"namadosen"`
	Cetak     string `json:"cetak"`
	Infomk    string `json:"infomk"`
}

type Nilai struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	NilAngka string `json:"nil_angka"`
	NilHuruf string `json:"nil_huruf"`
	Hadir    string `json:"hadir"`
	Projek   string `json:"projek"`
	Quiz     string `json:"quiz"`
	Tugas    string `json:"tugas"`
	UTS      string `json:"uts"`
	UAS      string `json:"uas"`
}

type RekapMKResponse struct {
	Total int          `json:"total"`
	Rows  []MataKuliah `json:"rows"`
}
