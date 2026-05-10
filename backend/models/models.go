package models

import (
	"time"
)

// ==========================================
// 1. MASTER DATA & AUTENTIKASI
// ==========================================

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"` // tanda "-" agar password tidak ikut terkirim ke frontend
	NamaLengkap  string    `gorm:"type:varchar(150);not null" json:"nama_lengkap"`
	Role         string    `gorm:"type:varchar(50);not null" json:"role"` // ADMIN / BIDAN / KADER / USER
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`

	// Relasi ke tabel pasien untuk Role USER (Orang tua / Ibu hamil)
	Balitas   []Balita   `gorm:"foreignKey:UserID" json:"balitas,omitempty"`
	IbuHamils []IbuHamil `gorm:"foreignKey:UserID" json:"ibu_hamils,omitempty"`
}

type JadwalKegiatan struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NamaKegiatan string    `gorm:"type:varchar(200);not null" json:"nama_kegiatan"`
	Tanggal      time.Time `gorm:"type:date;not null" json:"tanggal"`
	JamMulai     string    `gorm:"type:varchar(10);not null" json:"jam_mulai"`
	JamSelesai   string    `gorm:"type:varchar(10);not null" json:"jam_selesai"`
	Lokasi       string    `gorm:"type:varchar(200);not null" json:"lokasi"`
	Status       string    `gorm:"type:varchar(50);default:'AKAN DATANG'" json:"status"` // AKAN DATANG / SELESAI / BATAL
	CreatedByID  string    `gorm:"type:uuid;not null" json:"created_by_id"`
	CreatedAt    time.Time `json:"created_at"`

	CreatedBy User `gorm:"foreignKey:CreatedByID" json:"created_by"`
}

// ==========================================
// 2. ENTITAS PESERTA POSYANDU
// ==========================================

type Balita struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        string    `gorm:"type:uuid;not null" json:"user_id"` // Akun milik orang tua (Role: USER)
	NIK           string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nik"`
	NamaBalita    string    `gorm:"type:varchar(150);not null" json:"nama_balita"`
	TanggalLahir  time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
	JenisKelamin  string    `gorm:"type:varchar(10);not null" json:"jenis_kelamin"` // L / P
	NamaOrangTua  string    `gorm:"type:varchar(150);not null" json:"nama_orang_tua"`
	Alamat        string    `gorm:"type:text" json:"alamat"`
	CreatedAt     time.Time `json:"created_at"`

	Pemeriksaans []PemeriksaanBalita `gorm:"foreignKey:BalitaID" json:"pemeriksaans,omitempty"`
	Imunisasis   []ImunisasiBalita   `gorm:"foreignKey:BalitaID" json:"imunisasis,omitempty"`
}

type IbuHamil struct {
	ID              string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          string    `gorm:"type:uuid;not null" json:"user_id"` // Akun login milik ibu hamil
	NIK             string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nik"`
	NamaIbu         string    `gorm:"type:varchar(150);not null" json:"nama_ibu"`
	TanggalLahir    time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
	HPL             time.Time `gorm:"type:date;not null" json:"hpl"` // Hari Perkiraan Lahir
	Alamat          string    `gorm:"type:text" json:"alamat"`
	StatusKehamilan string    `gorm:"type:varchar(50);default:'AKTIF'" json:"status_kehamilan"` // AKTIF / SUDAH LAHIR
	CreatedAt       time.Time `json:"created_at"`

	Pemeriksaans []PemeriksaanIbuHamil `gorm:"foreignKey:IbuHamilID" json:"pemeriksaans,omitempty"`
}

type Lansia struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NIK            string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nik"`
	NamaLengkap    string    `gorm:"type:varchar(150);not null" json:"nama_lengkap"`
	TanggalLahir   time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
	JenisKelamin   string    `gorm:"type:varchar(10);not null" json:"jenis_kelamin"` // L / P
	Alamat         string    `gorm:"type:text" json:"alamat"`
	NamaPendamping string    `gorm:"type:varchar(150)" json:"nama_pendamping"`
	KontakDarurat  string    `gorm:"type:varchar(50)" json:"kontak_darurat"`
	Status         string    `gorm:"type:varchar(50);default:'AKTIF'" json:"status"` // AKTIF / PINDAH / MENINGGAL
	CreatedAt      time.Time `json:"created_at"`

	Pemeriksaans []PemeriksaanLansia `gorm:"foreignKey:LansiaID" json:"pemeriksaans,omitempty"`
}

// ==========================================
// 3. ENTITAS TRANSAKSIONAL (REKAM MEDIS)
// ==========================================

type PemeriksaanBalita struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BalitaID       string    `gorm:"type:uuid;not null" json:"balita_id"`
	TanggalPeriksa time.Time `gorm:"type:date;not null" json:"tanggal_periksa"`
	BeratBadan     float64   `gorm:"type:decimal(5,2);not null" json:"berat_badan"`
	TinggiBadan    float64   `gorm:"type:decimal(5,2);not null" json:"tinggi_badan"`
	LingkarKepala  float64   `gorm:"type:decimal(5,2)" json:"lingkar_kepala"`
	StatusGizi     string    `gorm:"type:varchar(50);not null" json:"status_gizi"` // NORMAL / GIZI BAIK / GIZI KURANG / STUNTING
	Catatan        string    `gorm:"type:text" json:"catatan"`
	DiperiksaOleh  string    `gorm:"type:uuid;not null" json:"diperiksa_oleh"` // ID Kader/Bidan
	CreatedAt      time.Time `json:"created_at"`

	Pemeriksa User `gorm:"foreignKey:DiperiksaOleh" json:"pemeriksa"`
}

type ImunisasiBalita struct {
	ID               string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BalitaID         string    `gorm:"type:uuid;not null" json:"balita_id"`
	JenisImunisasi   string    `gorm:"type:varchar(100);not null" json:"jenis_imunisasi"` // BCG, DPT, dll
	TanggalImunisasi time.Time `gorm:"type:date;not null" json:"tanggal_imunisasi"`
	Catatan          string    `gorm:"type:text" json:"catatan"`
	DicatatOleh      string    `gorm:"type:uuid;not null" json:"dicatat_oleh"`
	CreatedAt        time.Time `json:"created_at"`

	Pencatat User `gorm:"foreignKey:DicatatOleh" json:"pencatat"`
}

type PemeriksaanIbuHamil struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IbuHamilID     string    `gorm:"type:uuid;not null" json:"ibu_hamil_id"`
	TanggalPeriksa time.Time `gorm:"type:date;not null" json:"tanggal_periksa"`
	UsiaKehamilan  int       `gorm:"not null" json:"usia_kehamilan"` // format minggu
	TekananDarah   string    `gorm:"type:varchar(20);not null" json:"tekanan_darah"` // e.g. "120/80"
	BeratBadan     float64   `gorm:"type:decimal(5,2);not null" json:"berat_badan"`
	Catatan        string    `gorm:"type:text" json:"catatan"`
	DiperiksaOleh  string    `gorm:"type:uuid;not null" json:"diperiksa_oleh"`
	CreatedAt      time.Time `json:"created_at"`

	Pemeriksa User `gorm:"foreignKey:DiperiksaOleh" json:"pemeriksa"`
}

type PemeriksaanLansia struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LansiaID       string    `gorm:"type:uuid;not null" json:"lansia_id"`
	TanggalPeriksa time.Time `gorm:"type:date;not null" json:"tanggal_periksa"`
	TekananDarah   string    `gorm:"type:varchar(20);not null" json:"tekanan_darah"`
	GulaDarah      float64   `gorm:"type:decimal(6,2)" json:"gula_darah"` // mg/dL
	Kolesterol     float64   `gorm:"type:decimal(6,2)" json:"kolesterol"` // mg/dL
	Catatan        string    `gorm:"type:text" json:"catatan"`
	DiperiksaOleh  string    `gorm:"type:uuid;not null" json:"diperiksa_oleh"`
	CreatedAt      time.Time `json:"created_at"`

	Pemeriksa User `gorm:"foreignKey:DiperiksaOleh" json:"pemeriksa"`
}