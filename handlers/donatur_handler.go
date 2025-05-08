// handlers/donatur_handler.go
package handlers
import (
	"html/template"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// TemplateData adalah struktur untuk meneruskan data ke template HTML
type TemplateData struct {
	TotalDonatur    int64
	TotalDonasi     string
	KampanyeAktif   int
	TingkatPenyaluran string
	Donatur         []DonateurData
}

// DonateurData adalah struktur untuk data donatur yang akan ditampilkan
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

// GetLatestDonaturHandler menangani permintaan untuk menampilkan halaman donatur terbaru
func GetLatestDonaturHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mengambil data donasi dari database
		var donatur []struct {
			ID            uint
			Nama          string
			Email         string
			Jumlah        float64
			Kategori      string
			TanggalDonasi time.Time
		}

		// Mengambil 6 donasi terbaru dengan status pembayaran selesai
		// dan mengurutkannya berdasarkan tanggal donasi terbaru
		result := db.Table("donasi").
			Select("id, nama, email, jumlah, kategori, tanggal_donasi").
			Where("status_pembayaran = ?", "selesai").
			Order("tanggal_donasi desc").
			Limit(6).
			Find(&donatur)

		if result.Error != nil {
			http.Error(w, "Gagal mengambil data donatur: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Mengambil statistik untuk dashboard
		var totalDonatur int64
		db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Count(&totalDonatur)

		var totalDonasi float64
		db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Select("COALESCE(SUM(jumlah), 0)").Scan(&totalDonasi)

		// Data yang akan dikirim ke template
		data := TemplateData{
			TotalDonatur:     totalDonatur,
			TotalDonasi:      formatRupiah(totalDonasi),
			KampanyeAktif:    12,  // Nilai default atau dapat diambil dari database jika ada
			TingkatPenyaluran: "97%", // Nilai default atau dapat diambil dari database jika ada
			Donatur:          []DonateurData{},
		}

		// Proses data donatur untuk ditampilkan
		for _, d := range donatur {
			// Mencari donatur teratas (top donatur)
			isTop := false
			if len(data.Donatur) == 0 { // Anggap donatur pertama sebagai top donatur
				isTop = true
			}

			// Format waktu relatif
			waktuRelative := formatRelativeTime(d.TanggalDonasi)

			// Tambahkan ke data donatur
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
				Avatar:          "/api/placeholder/60/60", // Placeholder untuk avatar
			})
		}

		// Parse dan render template
		tmpl, err := template.ParseFiles("templates/donatur.html")
		if err != nil {
			http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Render template dengan data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Gagal merender template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Format angka ke format Rupiah
func formatRupiah(amount float64) string {
	if amount >= 1000000000 {
		return "Rp" + formatFloat(amount/1000000000) + " Milyar"
	} else if amount >= 1000000 {
		return "Rp" + formatFloat(amount/1000000) + " Juta"
	} else if amount >= 1000 {
		return "Rp" + formatFloat(amount/1000) + " Ribu"
	}
	return "Rp" + formatFloat(amount)
}

// Format float dengan presisi dua desimal
func formatFloat(num float64) string {
	return string(append([]byte{}, []byte(sprintf("%.2f", num))...))
}

// Format waktu relatif (contoh: "5 menit yang lalu")
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "baru saja"
	case diff < time.Hour:
		return sprintf("%d menit yang lalu", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return sprintf("%d jam yang lalu", int(diff.Hours()))
	default:
		return sprintf("%d hari yang lalu", int(diff.Hours()/24))
	}
}

// Struktur Donasi untuk mengambil data dari database
type Donasi struct {
	ID               uint      `gorm:"primaryKey"`
	Nama             string
	Email            string
	Telepon          string
	Jumlah           float64
	Kategori         string
	Keterangan       string
	Anonim           bool
	MetodePembayaran string
	StatusPembayaran string `gorm:"default:'menunggu'"`
	TanggalDonasi    time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}