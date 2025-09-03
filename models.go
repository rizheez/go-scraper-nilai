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
	KodeJrs   string `json:"kodejrs"`
	KodeMK    string `json:"kodemk"`
	KodePK    string `json:"kodepk"`
	Smtthnakd string `json:"smtthnakd"`
	NamaJrs   string `json:"namajrs"`
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

type BobotMK struct {
	MataKuliah MataKuliah `json:"mata_kuliah"`
	Bobot      Bobot      `json:"bobot"`
}

type Bobot struct {
	Hadir  string `json:"hdr"`
	Projek string `json:"projek"`
	Quiz   string `json:"quiz"`
	Tugas  string `json:"tgs"`
	UTS    string `json:"uts"`
	UAS    string `json:"uas"`
}

type RekapMKResponse struct {
	Total int          `json:"total"`
	Rows  []MataKuliah `json:"rows"`
}

type RekapMHSResponse struct {
	Total int         `json:"total"`
	Rows  []Mahasiswa `json:"rows"`
}

type Mahasiswa struct {
	IDMhs                   string  `json:"idmhs"`
	FakID                   string  `json:"fakid"`
	JrsID                   string  `json:"jrsid"`
	PkID                    string  `json:"pkid"`
	MnID                    string  `json:"mnid"`
	KodeFak                 string  `json:"kodefak"`
	KodeJrs                 string  `json:"kodejrs"`
	KodePk                  string  `json:"kodepk"`
	KodeMn                  string  `json:"kodemn"`
	KodePa                  string  `json:"kodepa"`
	Kurikulum               string  `json:"kurikulum"`
	NoSel                   string  `json:"nosel"`
	NIM                     string  `json:"nim"`
	NoTranskrip             string  `json:"no_transkrip"`
	NoPin                   string  `json:"no_pin"`
	Nama                    string  `json:"nama"`
	TempatLahir             string  `json:"tem_lahir"`
	TanggalLahir            string  `json:"tgl_lahir"`
	Gender                  string  `json:"gender"`
	Agama                   string  `json:"agama"`
	Marital                 string  `json:"marital"`
	NoKTP                   string  `json:"no_ktp"`
	AlamatSurat1            string  `json:"alm1_surat"`
	AlamatSurat2            string  `json:"alm2_surat"`
	RT_RW_Surat             string  `json:"rtrw_surat"`
	KotaSurat               string  `json:"kot_surat"`
	KodePosSurat            string  `json:"kdp_surat"`
	Telepon                 string  `json:"telepon"`
	HP1                     string  `json:"hp1"`
	HP2                     string  `json:"hp2"`
	Email                   string  `json:"email"`
	Tinggal                 string  `json:"tinggal"`
	NamaAyah                string  `json:"nama_ayah"`
	NamaIbu                 string  `json:"nama_ibu"`
	KerjaAyah               string  `json:"kerja_ayah"`
	KerjaIbu                string  `json:"kerja_ibu"`
	DidikAyah               string  `json:"didik_ayah"`
	DidikIbu                string  `json:"didik_ibu"`
	NikAyah                 string  `json:"nik_ayah"`
	NikIbu                  string  `json:"nik_ibu"`
	TanggalLahirAyah        string  `json:"tanggal_lahir_ayah"`
	TanggalLahirIbu         string  `json:"tanggal_lahir_ibu"`
	IdDidikAyah             string  `json:"id_didik_ayah"`
	IdDidikIbu              string  `json:"id_didik_ibu"`
	IdPenghasilanAyah       string  `json:"id_penghasilan_ayah"`
	IdPenghasilanIbu        string  `json:"id_penghasilan_ibu"`
	IdKerjaAyah             string  `json:"id_kerja_ayah"`
	IdKerjaIbu              string  `json:"id_kerja_ibu"`
	IdNPWPMhs               string  `json:"id_npwp_mhs"`
	AlamatOrtu              string  `json:"alamat_ortu"`
	KotaOrtu                string  `json:"kota_ortu"`
	KodePosOrtu             string  `json:"kodepos_ortu"`
	TelpOrtu                string  `json:"telp_ortu"`
	HPOrtu                  string  `json:"hp_ortu"`
	NamaSekolah             string  `json:"nama_sklh"`
	AlamatSekolah           string  `json:"alam_sklh"`
	JenisSekolah            string  `json:"jj_sklh"`
	Perusahaan              string  `json:"perusahaan"`
	AlamatPerusahaan        string  `json:"alm_perush"`
	KotaPerusahaan          string  `json:"kot_perush"`
	KodePosPerusahaan       string  `json:"kdp_perush"`
	TelpPerusahaan          string  `json:"tlp_perush"`
	FaxPerusahaan           string  `json:"fax_perush"`
	KDPTIMSMHS              string  `json:"kdptimsmhs"`
	KDJENMSMHS              string  `json:"kdjenmsmhs"`
	KDPSTMSMHS              string  `json:"kdpstmsmhs"`
	NIMHSMsmhs              string  `json:"nimhsmsmhs"`
	NMMHSMSMHS              string  `json:"nmmhsmsmhs"`
	ShiftMSMHS              string  `json:"shiftmsmhs"`
	TempatLahirMSMHS        string  `json:"tplhrmsmhs"`
	TglLahirMSMHS           string  `json:"tglhrmsmhs"`
	KDJekMSMHS              string  `json:"kdjekmsmhs"`
	TahunMSMHS              string  `json:"tahunmsmhs"`
	PeriodeSMTHN            string  `json:"smawlmsmhs"`
	BTSTUMSMHS              string  `json:"btstumsmhs"`
	ASMMMSMHS               string  `json:"assmamsmhs"`
	TGMSKMSMHS              string  `json:"tgmskmsmhs"`
	TGLLSMSMHS              string  `json:"tgllsmsmhs"`
	STMHSMSMHS              string  `json:"stmhsmsmhs"`
	STPIDMSMHS              string  `json:"stpidmsmhs"`
	SKSDIMSMHS              string  `json:"sksdimsmhs"`
	ASNIMMSMHS              string  `json:"asnimmsmhs"`
	ASPTIMSMHS              string  `json:"asptimsmhs"`
	ASJENMSMHS              string  `json:"asjenmsmhs"`
	ASPSTMSMHS              string  `json:"aspstmsmhs"`
	SMTHNLulus              string  `json:"smthnlulus"`
	NoSKLulus               string  `json:"nosklulus"`
	IDPerguruanTinggiAsal   string  `json:"id_perguruan_tinggi_asal"`
	NamaPerguruanTinggiAsal string  `json:"nama_perguruan_tinggi_asal"`
	IDProdiAsal             string  `json:"id_prodi_asal"`
	NamaProgramStudiAsal    string  `json:"nama_program_studi_asal"`
	Foto                    *string `json:"foto"`
	AsnMPTI                 string  `json:"asnmpti"`
	AsnmPST                 *string `json:"asnmpst"`
	JumMK                   string  `json:"jummk"`
	JumSKS                  string  `json:"jumsks"`
	JumUTU                  string  `json:"jumutu"`
	IPK                     string  `json:"ipk"`
	Logika                  string  `json:"logika"`
	CreateDate              string  `json:"createdate"`
	ModDate                 string  `json:"moddate"`
	RT                      string  `json:"rt"`
	RW                      string  `json:"rw"`
	Jalan                   string  `json:"jalan"`
	Dusun                   string  `json:"dusun"`
	KodePos                 string  `json:"kode_pos"`
	Kelurahan               string  `json:"kelurahan"`
	IDWilayah               string  `json:"id_wilayah"`
	NamaWilayah             string  `json:"nama_wilayah"`
	BiayaMasuk              string  `json:"biaya_masuk"`
	KdDaftar                string  `json:"kd_daftar"`
	KodeAgama               string  `json:"kode_agama"`
	NamaAgama               string  `json:"nama_agama"`
	IDRegPD                 string  `json:"id_reg_pd"`
	IDMahasiswa             string  `json:"id_mahasiswa"`
	IDJalurMasuk            string  `json:"id_jalur_masuk"`
	IDJnsTinggal            string  `json:"id_jns_tinggal"`
	IDJnsKeluar             string  `json:"id_jns_keluar"`
	IDJenjDidik             string  `json:"id_jenj_didik"`
	IDJnsDaftar             string  `json:"id_jns_daftar"`
	IDAlatTransport         string  `json:"id_alat_transport"`
	IDPenghasilan           string  `json:"id_penghasilan"`
	IDPembiayaan            string  `json:"id_pembiayaan"`
	IDKPS                   string  `json:"id_kps"`
	TotalTagihan            string  `json:"total_tagihan"`
	NamaFak                 string  `json:"namafak"`
	NamaJrs                 string  `json:"namajrs"`
	Jenjang                 string  `json:"jenjang"`
	NamaJJG                 string  `json:"nama_jjg"`
	KDJen                   string  `json:"kdjen"`
	KDPST                   string  `json:"kdpst"`
	BatasStudi              string  `json:"batastudi"`
	NamaPK                  string  `json:"namapk"`
	NamaMN                  *string `json:"namamn"`
	Group                   string  `json:"group"`
	TanggalMasuk            string  `json:"tgl_masuk"`
	StatusMhsKet            string  `json:"status_mhs_ket"`
	TempatTanggalLahir      string  `json:"temtgl_lahir"`
	SKSTotal                int     `json:"sks_total"`
}
