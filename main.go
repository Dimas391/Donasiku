// main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"math"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Donasi struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Nama             string    `json:"nama"`
	Email            string    `json:"email"`
	Telepon          string    `json:"telepon"`
	Jumlah           float64   `json:"jumlah"`
	Kategori         string    `json:"kategori"`
	Keterangan       string    `json:"keterangan"`
	Anonim           bool      `json:"anonim"`
	MetodePembayaran string    `json:"metode_pembayaran"`
	StatusPembayaran string    `json:"status_pembayaran" gorm:"default:'menunggu'"`
	TanggalDonasi    time.Time `json:"tanggal_donasi"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Komunitas struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nama      string    `json:"nama"`
	Deskripsi string    `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InfoPembayaran struct {
	Metode             string   `json:"metode"`
	NomorRekening      []string `json:"nomor_rekening,omitempty"`
	NamaPemilik        string   `json:"nama_pemilik,omitempty"`
	QRCodeURL          string   `json:"qr_code_url,omitempty"`
	PetunjukPembayaran string   `json:"petunjuk_pembayaran"`
}

type PenerimaBantuan struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Nama          string    `json:"nama"`
	Kategori      string    `json:"kategori"`
	JumlahDana    float64   `json:"jumlah_dana"`
	TanggalTerima time.Time `json:"tanggal_terima"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TemplateData struct {
	TotalDonatur      int64
	TotalDonasi       string
	KampanyeAktif     int
	TingkatPenyaluran string
	Donatur           []DonateurData
	UrgentCampaigns   []CampaignDisplayData
	NonUrgentCampaigns []CampaignDisplayData 
}

type DonateurData struct {
	ID              uint
	Nama            string
	Email           string
	Jumlah          float64
	JumlahFormatted string
	KampanyeNama    string
	WaktuDonasi     time.Time
	WaktuRelative   string
	IsTopDonatur    bool
	Avatar          string
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
	r.HandleFunc("/api/donasi/{id}/status", updateStatusPembayaran).Methods("PUT")
	r.HandleFunc("/api/donasi/kategori-total", getDonasiPerKategori).Methods("GET")
	r.HandleFunc("/api/dashboard/stats", getDashboardStats).Methods("GET")

	r.HandleFunc("/", getLatestDonaturPage).Methods("GET")
	r.HandleFunc("/donatur", getLatestDonaturPage).Methods("GET")
	r.HandleFunc("/semua-donatur", getAllDonaturPage).Methods("GET")

	r.HandleFunc("/api/campaigns/urgent", getUrgentCampaignsAPI).Methods("GET")
    r.HandleFunc("/api/campaigns/{id}", getCampaignByIDAPI).Methods("GET")
    
    // Page routes
    r.HandleFunc("/", getHomePage).Methods("GET")
    r.HandleFunc("/donatur", getLatestDonaturPage).Methods("GET")
    r.HandleFunc("/semua-donatur", getAllDonaturPage).Methods("GET")
    r.HandleFunc("/kampanye-mendesak", getUrgentCampaignsPage).Methods("GET")
    r.HandleFunc("/detail_donasi/{id}", getCampaignDetailPage).Methods("GET")
	r.HandleFunc("/api/quick-donate", quickDonate).Methods("POST")

    
	
    // Add a new route to display all campaigns with filters
    r.HandleFunc("/semua-kampanye", func(w http.ResponseWriter, r *http.Request) {
        // Redirect to filtered-donations with default parameters
        http.Redirect(w, r, "/filtered-donations?category=all&sort=newest", http.StatusSeeOther)
    }).Methods("GET")
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("image"))))

	r.Use(corsMiddleware)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := "8080"
	fmt.Println("Server berjalan di port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

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
}

func getCampaignByIDAPI(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "ID tidak valid", http.StatusBadRequest)
        return
    }
    
    campaign, err := getUrgentCampaign(uint(id))
    if err != nil {
        http.Error(w, "Kampanye tidak ditemukan: "+err.Error(), http.StatusNotFound)
        return
    }
	
    
    displayData := prepareCampaignDisplayData(campaign)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(displayData)
}

func getCampaignByID(id uint) (Campaign, error) {
    var campaign Campaign
    result := db.First(&campaign, id)
    return campaign, result.Error
}


func getCampaignDetailPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "ID tidak valid", http.StatusBadRequest)
        return
    }
    
    campaign, err := getCampaignByID(uint(id))
    if err != nil {
        http.Error(w, "Kampanye tidak ditemukan: "+err.Error(), http.StatusNotFound)
        return
    }
    
    displayData := prepareCampaignDisplayData(campaign)
    
    tmpl, err := template.ParseFiles("static/detail_donasi.html")
    if err != nil {
        http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    err = tmpl.Execute(w, displayData)
    if err != nil {
        http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}


func getLatestDonaturPage(w http.ResponseWriter, r *http.Request) {
	var donatur []struct {
		ID            uint
		Nama          string
		Email         string
		Jumlah        float64
		Kategori      string
		TanggalDonasi time.Time
	}

	queryResult := db.Table("donasis").
		Select("id, nama, email, jumlah, kategori, tanggal_donasi").
		Where("status_pembayaran = ?", "selesai").
		Order("tanggal_donasi desc").
		Limit(6).
		Find(&donatur)

	if queryResult.Error != nil {
		http.Error(w, "Gagal mengambil data donatur: "+queryResult.Error.Error(), http.StatusInternalServerError)
		return
	}

	var totalDonatur int64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Count(&totalDonatur)

	var totalDonasi float64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Select("COALESCE(SUM(jumlah), 0)").Scan(&totalDonasi)

	campaigns, err := getUrgentCampaigns()
    if err != nil {
        http.Error(w, "Gagal mengambil data kampanye mendesak: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    var urgentCampaigns []CampaignDisplayData
    for _, campaign := range campaigns {
        urgentCampaigns = append(urgentCampaigns, prepareCampaignDisplayData(campaign))
    }

	var displayCampaigns []CampaignDisplayData
    for _, campaign := range campaigns {
        displayCampaigns = append(displayCampaigns, prepareCampaignDisplayData(campaign))
    }

	var nonUrgentCampaigns []Campaign
	resultNonUrgent := db.Where("is_urgent = ?", false).Find(&nonUrgentCampaigns)
	if resultNonUrgent.Error != nil {
		http.Error(w, "Gagal mengambil data kampanye non-urgent: "+resultNonUrgent.Error.Error(), http.StatusInternalServerError)
		return
	}

	var nonUrgentCampaignsDisplay []CampaignDisplayData
	for _, campaign := range nonUrgentCampaigns {
		nonUrgentCampaignsDisplay = append(nonUrgentCampaignsDisplay, prepareCampaignDisplayData(campaign))
	}

	data := TemplateData{
		TotalDonatur:      totalDonatur,
		TotalDonasi:       formatRupiah(totalDonasi),
		KampanyeAktif:     12,
		TingkatPenyaluran: "97%",
		Donatur:           []DonateurData{},
		UrgentCampaigns:   urgentCampaigns,
		NonUrgentCampaigns: nonUrgentCampaignsDisplay,
	}

	for i, d := range donatur {
		isTop := false
		if i == 0 {
			isTop = true
		}

		waktuRelative := formatRelativeTime(d.TanggalDonasi)

		data.Donatur = append(data.Donatur, DonateurData{
			ID:              d.ID,
			Nama:            d.Nama,
			Email:           d.Email,
			Jumlah:          d.Jumlah,
			JumlahFormatted: formatRupiah(d.Jumlah),
			KampanyeNama:    d.Kategori,
			WaktuDonasi:     d.TanggalDonasi,
			WaktuRelative:   waktuRelative,
			IsTopDonatur:    isTop,
			Avatar:          "/static/img/avatar-placeholder.png",
		})
	}

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAllDonaturPage(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 20
	category := r.URL.Query().Get("category")
    sortBy := r.URL.Query().Get("sort")
    searchQuery := r.URL.Query().Get("q")

	if sortBy == "" {
        sortBy = "terbaru"
    }

	// Parse query parameters for pagination
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if s, err := strconv.Atoi(sizeParam); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	offset := (page - 1) * pageSize

	query := db.Table("donasis").
	Select("id, nama, email, jumlah, kategori, tanggal_donasi").
	Where("status_pembayaran = ?", "selesai")

		// Apply filters
		if category != "" && category != "all" {
			query = query.Where("kategori = ?", category)
		}

		if searchQuery != "" {
			query = query.Where("nama ILIKE ?", "%"+searchQuery+"%")
		}

		// Apply sorting
		switch sortBy {
		case "terlama":
			query = query.Order("tanggal_donasi asc")
		case "jumlah_tinggi":
			query = query.Order("jumlah desc")
		case "jumlah_rendah":
			query = query.Order("jumlah asc")
		default: // terbaru
			query = query.Order("tanggal_donasi desc")
		}


	// Get donatur data with pagination
	var donatur []struct {
		ID            uint
		Nama          string
		Email         string
		Jumlah        float64
		Kategori      string
		TanggalDonasi time.Time
	}

	queryResult := db.Table("donasis").
		Select("id, nama, email, jumlah, kategori, tanggal_donasi").
		Where("status_pembayaran = ?", "selesai").
		Order("tanggal_donasi desc").
		Offset(offset).
		Limit(pageSize).
		Find(&donatur)

	if queryResult.Error != nil {
		http.Error(w, "Gagal mengambil data donatur: "+queryResult.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Get total count for pagination
	var totalDonatur int64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Count(&totalDonatur)

	// Get total donation amount
	var totalDonasi float64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Select("COALESCE(SUM(jumlah), 0)").Scan(&totalDonasi)

	// Prepare template data
	data := TemplateData{
		TotalDonatur:      totalDonatur,
		TotalDonasi:       formatRupiah(totalDonasi),
		KampanyeAktif:     12,
		TingkatPenyaluran: "97%",
		Donatur:           []DonateurData{},
	}

	// // Transform donation data for the template
	// for _, d := range donatur {
	// 	waktuRelative := formatRelativeTime(d.TanggalDonasi)

	// 	data.Donatur = append(data.Donatur, DonateurData{
	// 		ID:              d.ID,
	// 		Nama:            d.Nama,
	// 		Email:           d.Email,
	// 		Jumlah:          d.Jumlah,
	// 		JumlahFormatted: formatRupiah(d.Jumlah),
	// 		KampanyeNama:    d.Kategori,
	// 		WaktuDonasi:     d.TanggalDonasi,
	// 		WaktuRelative:   waktuRelative,
	// 		Avatar:          "/static/img/avatar-placeholder.png",
	// 	})
	// }

	// Add pagination data to the template
	for i, d := range donatur {
        isTop := false
        if i == 0 && page == 1 && (sortBy == "terbaru" || sortBy == "jumlah_tinggi") {
            isTop = true
        }

        waktuRelative := formatRelativeTime(d.TanggalDonasi)

        data.Donatur = append(data.Donatur, DonateurData{
            ID:              d.ID,
            Nama:            d.Nama,
            Email:           d.Email,
            Jumlah:          d.Jumlah,
            JumlahFormatted: formatRupiah(d.Jumlah),
            KampanyeNama:    d.Kategori,
            WaktuDonasi:     d.TanggalDonasi,
            WaktuRelative:   waktuRelative,
            IsTopDonatur:    isTop,
            Avatar:          "/static/image/orang.jpg",
        })
    }

    // Calculate pagination links
    totalPages := int(math.Ceil(float64(totalDonatur) / float64(pageSize)))
    
    // Create a slice for pagination numbers
    pagination := make([]int, 0)
    
    // Show up to 5 page numbers centered around current page
    startPage := page - 2
    endPage := page + 2
    
    if startPage < 1 {
        startPage = 1
        endPage = int(math.Min(float64(5), float64(totalPages)))
    }
    
    if endPage > totalPages {
        endPage = totalPages
        startPage = int(math.Max(1, float64(totalPages-4)))
    }
    
    for i := startPage; i <= endPage; i++ {
        pagination = append(pagination, i)
    }

    tmpl, err := template.ParseFiles("static/semua-donatur.html")
    if err != nil {
        http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, map[string]interface{}{
        "Data":         data,
        "CurrentPage":  page,
        "TotalPages":   totalPages,
        "HasNextPage":  (page * pageSize) < int(totalDonatur),
        "HasPrevPage":  page > 1,
        "NextPage":     page + 1,
        "PrevPage":     page - 1,
        "Pagination":   pagination,
        "Category":     category,
        "Sort":         sortBy,
        "SearchQuery":  searchQuery,
    })
    
    if err != nil {
        http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
        return
    }


	err = tmpl.Execute(w, map[string]interface{}{
		"Data":         data,
		"CurrentPage":  page,
		"TotalPages":   int(math.Ceil(float64(totalDonatur) / float64(pageSize))),
		"HasNextPage":  (page * pageSize) < int(totalDonatur),
		"HasPrevPage":  page > 1,
		"NextPage":     page + 1,
		"PrevPage":     page - 1,
	})
	
	if err != nil {
		http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func formatRupiah(amount float64) string {
	if amount >= 1000000000 {
		return "Rp" + fmt.Sprintf("%.1f", amount/1000000000) + " Milyar"
	} else if amount >= 1000000 {
		return "Rp" + fmt.Sprintf("%.1f", amount/1000000) + " Juta"
	} else if amount >= 1000 {
		return "Rp" + fmt.Sprintf("%.1f", amount/1000) + " Ribu"
	}
	return "Rp" + fmt.Sprintf("%.0f", amount)
}

func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "baru saja"
	case diff < time.Hour:
		return fmt.Sprintf("%d menit yang lalu", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d jam yang lalu", int(diff.Hours()))
	default:
		return fmt.Sprintf("%d hari yang lalu", int(diff.Hours()/24))
	}
}

func initDB() {
	var err error
	maxRetries := 5
	retryDelay := time.Second * 3

	for retries := 0; retries < maxRetries; retries++ {
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

	db.AutoMigrate(&Donasi{}, &Komunitas{}, &PenerimaBantuan{}, &Campaign{})
	fmt.Println("Database berhasil dimigrasi")
}

// [Rest of your handler functions remain the same...]
// getInfoPembayaran, getDonasi, getDonasiByID, etc.

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
			NamaPemilik:        "Platform Donasi",
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
			Metode:             "qris",
			QRCodeURL:          "/img/qris-sample.png",
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
	donasi.StatusPembayaran = "selesai"

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
		Where("status_pembayaran = ?", "selesai").
		Select("COALESCE(SUM(jumlah), 0)").
		Scan(&stats.TotalDonasi)

	// Hitung jumlah donatur yang unik (berdasarkan email)
	db.Model(&Donasi{}).
		Where("status_pembayaran = ?", "selesai").
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
		Where("status_pembayaran = ?", "selesai").
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

// GetDonorsByCampaign fetches donors for a specific campaign
func GetDonorsByCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Get the latest donors for this campaign
	var donors []struct {
		ID            uint
		Nama          string
		Email         string
		Jumlah        float64
		Anonim        bool
		TanggalDonasi time.Time
	}

	// You might need to adjust this query based on your actual database schema
	// This assumes you have a campaign_id field in the donasi table
	queryResult := db.Table("donasis").
		Select("id, nama, email, jumlah, anonim, tanggal_donasi").
		Where("status_pembayaran = ? AND campaign_id = ?", "selesai", campaignID).
		Order("tanggal_donasi desc").
		Limit(10).
		Find(&donors)

	if queryResult.Error != nil {
		http.Error(w, "Gagal mengambil data donatur: "+queryResult.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Format the donor data for response
	var donorData []struct {
		ID              uint
		Nama            string
		Email           string
		Jumlah          float64
		JumlahFormatted string
		Anonim          bool
		WaktuDonasi     time.Time
		WaktuRelative   string
	}

	for _, d := range donors {
		donorData = append(donorData, struct {
			ID              uint
			Nama            string
			Email           string
			Jumlah          float64
			JumlahFormatted string
			Anonim          bool
			WaktuDonasi     time.Time
			WaktuRelative   string
		}{
			ID:              d.ID,
			Nama:            d.Nama,
			Email:           d.Email,
			Jumlah:          d.Jumlah,
			JumlahFormatted: formatRupiah(d.Jumlah),
			Anonim:          d.Anonim,
			WaktuDonasi:     d.TanggalDonasi,
			WaktuRelative:   formatRelativeTime(d.TanggalDonasi),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(donorData)
}

// Add the campaign ID field to the Donasi struct
func addCampaignIDFieldToDonasi() {
	// This is a placeholder function to represent the schema change
	// In reality, you would need to modify your database schema and struct definition
	
	// Example of modified Donasi struct:
	/*
	type Donasi struct {
		ID               uint      `json:"id" gorm:"primaryKey"`
		Nama             string    `json:"nama"`
		Email            string    `json:"email"`
		Telepon          string    `json:"telepon"`
		Jumlah           float64   `json:"jumlah"`
		Kategori         string    `json:"kategori"`
		Keterangan       string    `json:"keterangan"`
		Anonim           bool      `json:"anonim"`
		MetodePembayaran string    `json:"metode_pembayaran"`
		StatusPembayaran string    `json:"status_pembayaran" gorm:"default:'menunggu'"`
		TanggalDonasi    time.Time `json:"tanggal_donasi"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		CampaignID       uint      `json:"campaign_id"` // Added field
	}
	*/
}

// Real-time update for campaign stats
func updateCampaignStats(campaignID uint, amount float64) error {
	// Get the current campaign
	var campaign Campaign
	if err := db.First(&campaign, campaignID).Error; err != nil {
		return err
	}

	// Update the current amount
	campaign.CurrentAmount += amount
	
	// Save the changes
	if err := db.Save(&campaign).Error; err != nil {
		return err
	}

	return nil
}

// Handle donation creation with campaign updates
func createDonasiWithCampaignUpdate(w http.ResponseWriter, r *http.Request) {
	var donasi Donasi
	err := json.NewDecoder(r.Body).Decode(&donasi)
	if err != nil {
		http.Error(w, "Format data tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set donation date to current time
	donasi.TanggalDonasi = time.Now()

	// Create the donation
	if err := db.Create(&donasi).Error; err != nil {
		http.Error(w, "Gagal membuat donasi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If the payment status is "selesai", update campaign stats
	if donasi.StatusPembayaran == "selesai" && donasi.ID > 0 {
		if err := updateCampaignStats(donasi.ID, donasi.Jumlah); err != nil {
			// Log the error but continue
			fmt.Printf("Error updating campaign stats: %v\n", err)
		}
	}

	// Return the created donation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(donasi)
}

// Handler to update donation status and update campaign stats if needed
// func updateDonasiStatusWithCampaignUpdate(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.ParseUint(vars["id"], 10, 64)
// 	if err != nil {
// 		http.Error(w, "ID tidak valid", http.StatusBadRequest)
// 		return
// 	}

// 	// Get the existing donation
// 	var existingDonasi Donasi
// 	if err := db.First(&existingDonasi, id).Error; err != nil {
// 		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
// 		return
// 	}

// 	// Parse the status update
// 	var statusUpdate struct {
// 		StatusPembayaran string `json:"status_pembayaran"`
// 	}
	
// 	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
// 		http.Error(w, "Format data tidak valid: "+err.Error(), http.StatusBadRequest)
// 		return
//  	}

// 	// Save old status for comparison
// 	// oldStatus := existingDonasi.StatusPembayaran

// 	// Update the status
// 	existingDonasi.StatusPembayaran = statusUpdate.StatusPembayaran
// 	existingDonasi.UpdatedAt = time.Now()

// 	// Save the changes
// 	if err := db.Save(&existingDonasi).Error; err != nil {
// 		http.Error(w, "Gagal memperbarui status: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

	// Handler to update donation status and update campaign stats if needed
func updateDonasiStatusWithCampaignUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	// Get the existing donation
	var existingDonasi Donasi
	if err := db.First(&existingDonasi, id).Error; err != nil {
		http.Error(w, "Donasi tidak ditemukan", http.StatusNotFound)
		return
	}

	// Parse the status update
	var statusUpdate struct {
		StatusPembayaran string `json:"status_pembayaran"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, "Format data tidak valid: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Save old status for comparison
	oldStatus := existingDonasi.StatusPembayaran

	// Update the status
	existingDonasi.StatusPembayaran = statusUpdate.StatusPembayaran
	existingDonasi.UpdatedAt = time.Now()

	// Save the changes
	if err := db.Save(&existingDonasi).Error; err != nil {
		http.Error(w, "Gagal memperbarui status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If status changes from non-selesai to selesai, update campaign stats
	if oldStatus != "selesai" && existingDonasi.StatusPembayaran == "selesai" && existingDonasi.ID > 0 {
		if err := updateCampaignStats(existingDonasi.ID, existingDonasi.Jumlah); err != nil {
			// Log the error but continue
			fmt.Printf("Error updating campaign stats: %v\n", err)
		}
	}

	// If status changes from selesai to non-selesai, reverse the campaign stats update
	if oldStatus == "selesai" && existingDonasi.StatusPembayaran != "selesai" && existingDonasi.ID > 0 {
		if err := updateCampaignStats(existingDonasi.ID, -existingDonasi.Jumlah); err != nil {
			// Log the error but continue
			fmt.Printf("Error updating campaign stats: %v\n", err)
		}
	}

	// Return the updated donation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingDonasi)
}

// Handler for quick donations - add this to main.go
func quickDonate(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    var quickDonation struct {
        Jumlah    float64 `json:"jumlah"`
        KampanyeID uint    `json:"kampanye_id"`
    }

    err := json.NewDecoder(r.Body).Decode(&quickDonation)
    if err != nil {
        http.Error(w, "Invalid request data: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate data
    if quickDonation.Jumlah <= 0 {
        http.Error(w, "Jumlah donasi harus lebih dari 0", http.StatusBadRequest)
        return
    }

    // Get campaign name for the selected campaign ID
    var campaignName string
    if quickDonation.KampanyeID > 0 {
        var campaign Campaign
        result := db.First(&campaign, quickDonation.KampanyeID)
        if result.Error != nil {
            http.Error(w, "Kampanye tidak ditemukan", http.StatusBadRequest)
            return
        }
        campaignName = campaign.Title
    } else {
        campaignName = "Donasi Umum"
    }

    // Create a new donation with available information
    // Note: For quick donations, we set minimal information
    // The user can complete the rest later if needed
    donasi := Donasi{
        Nama:             "Donatur Kilat", // Default name for quick donors
        Email:            "anonymous@example.com", // Default email
        Jumlah:           quickDonation.Jumlah,
        Kategori:         campaignName,
        Keterangan:       "Donasi Kilat",
        Anonim:           true, // Quick donations are anonymous by default
        MetodePembayaran: "Transfer Bank", // Default payment method
        StatusPembayaran: "selesai", // Mark as completed immediately for simplicity
        TanggalDonasi:    time.Now(),
    }

    // Save to database
    result := db.Create(&donasi)
    if result.Error != nil {
        http.Error(w, "Failed to create donation: "+result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Update campaign stats if a campaign was selected
    if quickDonation.KampanyeID > 0 {
        if err := updateCampaignStats(quickDonation.KampanyeID, quickDonation.Jumlah); err != nil {
            // Log the error but continue
            fmt.Printf("Error updating campaign stats: %v\n", err)
        }
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Donasi berhasil",
        "donasi_id": donasi.ID,
        "jumlah": donasi.Jumlah,
    })
}

