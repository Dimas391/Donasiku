// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Donasi struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Nama            string    `json:"nama"`
	Email           string    `json:"email"`
	Telepon         string    `json:"telepon"`
	Jumlah          float64   `json:"jumlah"`
	Kategori        string    `json:"kategori"`
	Keterangan      string    `json:"keterangan"`
	Anonim          bool      `json:"anonim"`
	MetodePembayaran string   `json:"metode_pembayaran"`
	StatusPembayaran string   `json:"status_pembayaran" gorm:"default:'menunggu'"`
	TanggalDonasi   time.Time `json:"tanggal_donasi"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Komunitas struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InfoPembayaran struct {
	Metode          string   `json:"metode"`
	NomorRekening   []string `json:"nomor_rekening,omitempty"`
	NamaPemilik     string   `json:"nama_pemilik,omitempty"`
	QRCodeURL       string   `json:"qr_code_url,omitempty"`
	PetunjukPembayaran string `json:"petunjuk_pembayaran"`
}

type PenerimaBantuan struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Nama        string    `json:"nama"`
	Kategori    string    `json:"kategori"`
	JumlahDana  float64   `json:"jumlah_dana"`
	TanggalTerima time.Time `json:"tanggal_terima"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *gorm.DB

func main() {
	initDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/donasi", getDonasi).Methods("GET")
	r.HandleFunc("/api/donasi", createDonasi).Methods("POST")
	r.HandleFunc("/api/donasi/{id}", getDonasiByID).Methods("GET")
	r.HandleFunc("/api/donasi/{id}", updateDonasi).Methods("PUT")
	r.HandleFunc("/api/donasi/{id}", deleteDonasi).Methods("DELETE")

	r.HandleFunc("/api/komunitas", getKomunitas).Methods("GET")
	r.HandleFunc("/api/komunitas", createKomunitas).Methods("POST")
	r.HandleFunc("/api/komunitas/{id}", getKomunitasByID).Methods("GET")

	r.HandleFunc("/api/donasi/kategori/{kategori}", getDonasiByKategori).Methods("GET")

	r.HandleFunc("/api/donasi/total", getTotalDonasi).Methods("GET")
	
	r.HandleFunc("/api/pembayaran/{metode}", getInfoPembayaran).Methods("GET")
	
	// Endpoint untuk update status pembayaran
	r.HandleFunc("/api/donasi/{id}/status", updateStatusPembayaran).Methods("PUT")

	r.HandleFunc("/api/donasi/kategori-total", getDonasiPerKategori).Methods("GET")
	r.HandleFunc("/api/dashboard/stats", getDashboardStats).Methods("GET")

	// Mengizinkan CORS untuk pengembangan
	r.Use(corsMiddleware)

	// Static file server untuk frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	

	// Mulai server
	port := "8080"
	fmt.Println("Server berjalan di port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	db.AutoMigrate(&Donasi{}, &Komunitas{}, &PenerimaBantuan{})
}

// Middleware untuk CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})

	// Add these functions to your main.go file

// Handler untuk mendapatkan statistik dashboard


// Lalu tambahkan fungsi ini di main() untuk mendaftarkan endpoint baru

}

func initDB() {
	var err error
	maxRetries := 5
	retryDelay := time.Second * 3
	
	// Retry connection logic
	for retries := 0; retries < maxRetries; retries++ {
		// Koneksi ke database PostgreSQL di docker
		dsn := "host=db user=postgres password=postgres dbname=donasi port=5432 sslmode=disable TimeZone=Asia/Jakarta"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %d: Failed to connect to database: %v", retries+1, err)
			if retries < maxRetries-1 {
				log.Printf("Retrying in %v...", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			log.Fatal("All connection attempts failed:", err)
		} else {
			break
		}
	}

	// Migrasi database
	db.AutoMigrate(&Donasi{}, &Komunitas{})
	fmt.Println("Database berhasil dimigrasi")
}

// Handler untuk info pembayaran
func getInfoPembayaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	metode := vars["metode"]
	
	var infoPembayaran InfoPembayaran
	
	switch metode {
	case "transfer":
		infoPembayaran = InfoPembayaran{
			Metode: "transfer",
			NomorRekening: []string{
				"BCA: 1234567890",
				"Mandiri: 0987654321",
				"BNI: 1122334455",
			},
			NamaPemilik: "Platform Donasi",
			PetunjukPembayaran: "Selesaikan pembayaran dalam 24 jam",
		}
	case "ewallet":
		infoPembayaran = InfoPembayaran{
			Metode: "ewallet",
			NomorRekening: []string{
				"GoPay: 0812345678",
				"OVO: 0812345678",
				"Dana: 0812345678",
			},
			PetunjukPembayaran: "Selesaikan pembayaran dalam 24 jam",
		}
	case "qris":
		infoPembayaran = InfoPembayaran{
			Metode: "qris",
			QRCodeURL: "/img/qris-sample.png",
			PetunjukPembayaran: "Scan QR Code menggunakan aplikasi e-wallet favorit Anda",
		}
	default:
		http.Error(w, "Metode pembayaran tidak valid", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(infoPembayaran)
}

// Handler untuk donasi
func getDonasi(w http.ResponseWriter, r *http.Request) {
    var donasi []Donasi
    
    // Check if limit parameter is provided
    limitParam := r.URL.Query().Get("limit")
    
    query := db
    
    // Apply limit if provided
    if limitParam != "" {
        limit, err := strconv.Atoi(limitParam)
        if err == nil && limit > 0 {
            query = query.Limit(limit)
        }
    }
    
    // Order by most recent donations
    query = query.Order("tanggal_donasi desc")
    
    // Execute the query
    result := query.Find(&donasi)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(donasi)
}

func getDonasiByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var donasi Donasi
	result := db.First(&donasi, id)
	if result.Error != nil {
		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donasi)
}

func getDonasiByKategori(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kategori := vars["kategori"]

	var donasi []Donasi
	result := db.Where("kategori = ?", kategori).Find(&donasi)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donasi)
}

func getTotalDonasi(w http.ResponseWriter, r *http.Request) {
	// Menghitung total donasi
	type TotalDonasi struct {
		Total float64 `json:"total"`
		Count int64   `json:"count"`
	}

	var total TotalDonasi
	db.Model(&Donasi{}).Select("COALESCE(SUM(jumlah), 0) as total, COUNT(*) as count").Scan(&total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(total)
}

func createDonasi(w http.ResponseWriter, r *http.Request) {
	var donasi Donasi
	err := json.NewDecoder(r.Body).Decode(&donasi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi data donasi
	if donasi.Nama == "" || donasi.Email == "" || donasi.Jumlah <= 0 || donasi.Kategori == "" {
		http.Error(w, "Data donasi tidak lengkap", http.StatusBadRequest)
		return
	}

	// Set default value
	donasi.TanggalDonasi = time.Now()
	donasi.StatusPembayaran = "menunggu" // Status awal menunggu pembayaran
	
	// Simpan ke database
	result := db.Create(&donasi)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(donasi)
}

func updateStatusPembayaran(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var donasi Donasi
	result := db.First(&donasi, id)
	if result.Error != nil {
		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
		return
	}

	// Decode request body
	var statusUpdate struct {
		Status string `json:"status"`
	}
	
	err = json.NewDecoder(r.Body).Decode(&statusUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Update status pembayaran
	donasi.StatusPembayaran = statusUpdate.Status
	db.Save(&donasi)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donasi)
}

func updateDonasi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var donasi Donasi
	result := db.First(&donasi, id)
	if result.Error != nil {
		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
		return
	}

	var updatedDonasi Donasi
	err = json.NewDecoder(r.Body).Decode(&updatedDonasi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update field-field donasi sesuai form
	donasi.Nama = updatedDonasi.Nama
	donasi.Email = updatedDonasi.Email
	donasi.Telepon = updatedDonasi.Telepon
	donasi.Jumlah = updatedDonasi.Jumlah
	donasi.Kategori = updatedDonasi.Kategori
	donasi.Keterangan = updatedDonasi.Keterangan
	donasi.Anonim = updatedDonasi.Anonim
	donasi.MetodePembayaran = updatedDonasi.MetodePembayaran
	
	// Jangan izinkan perubahan status pembayaran melalui update umum
	// StatusPembayaran harus diubah melalui endpoint khusus

	db.Save(&donasi)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donasi)
}

func deleteDonasi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var donasi Donasi
	result := db.First(&donasi, id)
	if result.Error != nil {
		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
		return
	}

	db.Delete(&donasi)
	w.WriteHeader(http.StatusNoContent)
}

// Handler untuk komunitas
func getKomunitas(w http.ResponseWriter, r *http.Request) {
	var komunitas []Komunitas
	result := db.Find(&komunitas)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(komunitas)
}

func getKomunitasByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var komunitas Komunitas
	result := db.First(&komunitas, id)
	if result.Error != nil {
		http.Error(w, "Komunitas tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(komunitas)
}

// Handler untuk mendapatkan total donasi per kategori
func getDonasiPerKategori(w http.ResponseWriter, r *http.Request) {
    type KategoriTotal struct {
        Kategori string  `json:"kategori"`
        Total    float64 `json:"total"`
        Count    int64   `json:"count"`
    }
    
    var result []KategoriTotal
    
    // Query untuk menghitung total per kategori
    rows, err := db.Model(&Donasi{}).
        Select("kategori, COALESCE(SUM(jumlah), 0) as total, COUNT(*) as count").
        Group("kategori").
        Rows()
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    // Iterasi hasil query
    for rows.Next() {
        var kt KategoriTotal
        if err := db.ScanRows(rows, &kt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        result = append(result, kt)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func createKomunitas(w http.ResponseWriter, r *http.Request) {
	var komunitas Komunitas
	err := json.NewDecoder(r.Body).Decode(&komunitas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&komunitas)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(komunitas)
}

func getDashboardStats(w http.ResponseWriter, r *http.Request) {
	type DashboardStats struct {
		TotalDonasi        float64 `json:"total_donasi"`
		JumlahDonatur      int64   `json:"jumlah_donatur"`
		PenerimaBantuan    int64   `json:"penerima_bantuan"`
		DistribusiKategori []struct {
			Kategori string  `json:"kategori"`
			Total    float64 `json:"total"`
			Persen   float64 `json:"persen"`
		} `json:"distribusi_kategori"`
	}

	var stats DashboardStats

	// Hitung total donasi
	db.Model(&Donasi{}).
		Where("status_pembayaran = ?", "menunggu").
		Select("COALESCE(SUM(jumlah), 0)").
		Scan(&stats.TotalDonasi)

	// Hitung jumlah donatur yang unik (berdasarkan email)
	db.Model(&Donasi{}).
		Where("status_pembayaran = ?", "menunggu").
		Distinct("email").
		Count(&stats.JumlahDonatur)

	// Penerima bantuan (misalnya dari model/tabel lain yang belum dibuat)
	// Contoh: Jika menggunakan data komunitas sebagai penerima
	db.Model(&Komunitas{}).Count(&stats.PenerimaBantuan)

	// Distribusi per kategori
	var kategoriTotals []struct {
		Kategori string  `json:"kategori"`
		Total    float64 `json:"total"`
	}

	db.Model(&Donasi{}).
		Where("status_pembayaran = ?", "menunggu").
		Select("kategori, COALESCE(SUM(jumlah), 0) as total").
		Group("kategori").
		Scan(&kategoriTotals)

	// Hitung persentase untuk setiap kategori
	if stats.TotalDonasi > 0 {
		for _, kt := range kategoriTotals {
			stats.DistribusiKategori = append(stats.DistribusiKategori, struct {
				Kategori string  `json:"kategori"`
				Total    float64 `json:"total"`
				Persen   float64 `json:"persen"`
			}{
				Kategori: kt.Kategori,
				Total:    kt.Total,
				Persen:   (kt.Total / stats.TotalDonasi) * 100,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
