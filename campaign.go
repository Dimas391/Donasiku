package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"time"
	"encoding/json"
)

// Campaign struct to match the database schema
type Campaign struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	ImagePath     string    `json:"image_path"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	IsUrgent      bool      `json:"is_urgent"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Content       string    `json:"content"`
	TotalDonations float64 `json:"total_donations"` // Total donations for this campaign

}

// CampaignDisplayData is used for template rendering
type CampaignDisplayData struct {
    ID                uint    
    Title             string    
    Description       string    
    ShortDescription  string    
    Category          string   
    TargetAmount      float64  
    TargetFormatted   string    
    CurrentAmount     float64  
    CurrentFormatted  string   
    Progress          int      
    DaysRemaining     int      
    StartDate         time.Time
    EndDate           time.Time 
    IsUrgent          bool     
    IsActive          bool      
    ImagePath       string 
    DonationCount     int  
	HoursRemaining  int
	TimeRemaining   string
    CountdownTimer  string
}

// Helper function untuk parsing time
func parseTime(timeStr string) time.Time {
    layouts := []string{
        "2006-01-02 15:04:05-07",
        "2006-01-02T15:04:05-07:00",
        "2006-01-02 15:04:05.000000-07",
        "2006-01-02T15:04:05.000000-07:00",
    }
    
    for _, layout := range layouts {
        if t, err := time.Parse(layout, timeStr); err == nil {
            return t
        }
    }
    
    // Fallback ke waktu sekarang jika parsing gagal
    return time.Now()
}

func seedCampaigns() error {
    // Cek apakah data sudah ada
    var count int64
    db.Model(&Campaign{}).Count(&count)
    
    // Jika sudah ada data, skip seeding
    if count > 0 {
        fmt.Println("Campaign data already exists, skipping seed...")
        return nil
    }

    campaigns := []Campaign{
        {
            ID:             1,
            Title:          "Bantuan Banjir Hagu",
            Description:    "Ribuan keluarga terdampak banjir di Kalimantan Selatan membutuhkan bantuan mendesak untuk kebutuhan pokok, obat-obatan, dan tempat penampungan sementara",
            Category:       "Bencana Alam",
            ImagePath:      "../image/banjir.png",
            TargetAmount:   250000000,
            CurrentAmount:  308023467,
            IsUrgent:       true,
            EndDate:        parseTime("2025-06-30 23:59:59+00"),
            CreatedAt:      parseTime("2025-05-06 11:22:23.889886+00"),
            UpdatedAt:      parseTime("2025-05-15 02:44:45.932721+00"),
            Content:        "Banjir bandang telah melanda kawasan Hagu di Kalimantan Selatan sejak tiga hari lalu, mengakibatkan ribuan warga kehilangan tempat tinggal. Saat ini mereka tinggal di pengungsian dengan kondisi yang sangat memprihatinkan. Dana yang terkumpul akan digunakan untuk: 1) Penyediaan makanan siap saji, 2) Pakaian layak pakai, 3) Obat-obatan, 4) Perlengkapan bayi dan anak, 5) Selimut dan kasur lipat",
            TotalDonations: 0,
        },
        {
            ID:             2,
            Title:          "Pendidikan Anak Desa",
            Description:    "Bantu anak-anak desa terpencil mendapatkan akses pendidikan layak",
            Category:       "Pendidikan",
            ImagePath:      "../image/pendidikan.jpg",
            TargetAmount:   150000000,
            CurrentAmount:  75500000,
            IsUrgent:       false,
            EndDate:        parseTime("2025-06-05 11:22:23.897198+00"),
            CreatedAt:      parseTime("2025-05-06 11:22:23.898105+00"),
            UpdatedAt:      parseTime("2025-05-08 15:14:25.309527+00"),
            Content:        "Di daerah terpencil Kalimantan Timur, masih banyak anak yang harus berjalan berkilo-kilometer untuk mencapai sekolah. Banyak yang akhirnya putus sekolah karena kendala jarak dan biaya. Kampanye ini bertujuan untuk membangun perpustakaan dan pusat belajar di desa-desa terpencil.",
            TotalDonations: 0,
        },
        {
            ID:             3,
            Title:          "Panti Asuhan Cahaya Kasih",
            Description:    "Bantu renovasi panti asuhan yang menampung 50 anak yatim piatu",
            Category:       "Sosial",
            ImagePath:      "../image/panti.png",
            TargetAmount:   100000000,
            CurrentAmount:  45000000,
            IsUrgent:       false,
            EndDate:        parseTime("2025-05-26T11:22:23.897199+00:00"),
            CreatedAt:      parseTime("2025-05-06T11:22:23.901437+00:00"),
            UpdatedAt:      parseTime("2025-05-06T11:22:23.901437+00:00"),
            Content:        "Panti Asuhan Cahaya Kasih sudah berdiri sejak 15 tahun lalu dan telah membantu ratusan anak mendapatkan pendidikan dan kehidupan yang layak. Namun, bangunan panti yang sudah tua membutuhkan renovasi mendesak terutama di bagian atap yang bocor dan sistem sanitasi yang rusak.",
            TotalDonations: 0,
        },
        {
            ID:             4,
            Title:          "Donasi Kurban 2025",
            Description:    "Bantu saudara kita yang kurang mampu merayakan Idul Adha dengan menyediakan hewan kurban untuk wilayah terpencil.",
            Category:       "Keagamaan",
            ImagePath:      "../image/kurban.jpeg",
            TargetAmount:   100000000,
            CurrentAmount:  25000000,
            IsUrgent:       false,
            EndDate:        parseTime("2025-06-15T12:00:00+00:00"),
            CreatedAt:      parseTime("2025-05-08T10:35:04.058487+00:00"),
            UpdatedAt:      parseTime("2025-05-08T10:35:04.058487+00:00"),
            Content:        "Dana yang terkumpul akan digunakan untuk pembelian hewan kurban seperti kambing dan sapi, yang kemudian akan disalurkan ke daerah terpencil yang membutuhkan. Setiap kurban akan didistribusikan secara adil dan transparan kepada masyarakat yang membutuhkan.",
            TotalDonations: 0,
        },
        {
            ID:             5,
            Title:          "Bantu Komunitas Tunanetra",
            Description:    "Mari bantu komunitas tunanetra mendapatkan alat bantu baca Braille, tongkat, dan pelatihan keterampilan.",
            Category:       "Disabilitas",
            ImagePath:      "../image/tunanetra.png",
            TargetAmount:   75000000,
            CurrentAmount:  12000000,
            IsUrgent:       false,
            EndDate:        parseTime("2025-07-01T12:00:00+00:00"),
            CreatedAt:      parseTime("2025-05-08T10:36:30.296597+00:00"),
            UpdatedAt:      parseTime("2025-05-08T10:36:30.296597+00:00"),
            Content:        "Donasi akan digunakan untuk membeli alat bantu seperti buku Braille, tongkat pintar, serta memberikan pelatihan keterampilan kerja bagi para penyandang tunanetra agar mereka bisa lebih mandiri dan produktif di masyarakat.",
            TotalDonations: 0,
        },
    }

    // Insert semua campaign sekaligus
    if err := db.Create(&campaigns).Error; err != nil {
        return fmt.Errorf("failed to seed campaigns: %v", err)
    }

    fmt.Printf("Successfully seeded %d campaigns\n", len(campaigns))
    return nil
}


// getUrgentCampaigns fetches all urgent campaigns from database
func getUrgentCampaigns() ([]Campaign, error) {
	var campaigns []Campaign
	result := db.Where("is_urgent = ?", true).Find(&campaigns)
	return campaigns, result.Error
}

// getUrgentCampaign fetches a specific urgent campaign by ID
func getUrgentCampaign(id uint) (Campaign, error) {
	var campaign Campaign
	result := db.Where("id = ? AND is_urgent = ?", id, true).First(&campaign)
	return campaign, result.Error
}

func getNonUrgentCampaigns() ([]Campaign, error) {
    var campaigns []Campaign
    result := db.Where("is_urgent = ?", false).Find(&campaigns)
    if result.Error != nil {
        return nil, result.Error
    }
    return campaigns, nil
}


// prepareCampaignDisplayData formats campaign data for display
func prepareCampaignDisplayData(campaign Campaign) CampaignDisplayData {
	// Calculate progress percentage
	progress := int(math.Floor((campaign.CurrentAmount / campaign.TargetAmount) * 100))
	
	// Calculate remaining time
	now := time.Now()
	duration := campaign.EndDate.Sub(now)
	daysRemaining := int(duration.Hours() / 24)
	hoursRemaining := int(math.Floor(duration.Hours())) % 24
	
	// Format countdown timer (days:hours:minutes:seconds)
	minutes := int(math.Floor(duration.Minutes())) % 60
	seconds := int(math.Floor(duration.Seconds())) % 60
	countdownTimer := fmt.Sprintf("%02d:%02d:%02d:%02d", daysRemaining, hoursRemaining, minutes, seconds)
	
	// Format time remaining text
	timeRemaining := fmt.Sprintf("%d hari %d jam", daysRemaining, hoursRemaining)
	
	return CampaignDisplayData{
		ID:               campaign.ID,
		Title:            campaign.Title,
		Description:      campaign.Description,
		Category:         campaign.Category,
		ImagePath:        campaign.ImagePath,
		TargetAmount:     campaign.TargetAmount,
		TargetFormatted:  formatRupiah(campaign.TargetAmount),
		CurrentAmount:    campaign.CurrentAmount,
		CurrentFormatted: formatRupiah(campaign.CurrentAmount),
		IsUrgent:         campaign.IsUrgent,
		Progress:         progress,
		DaysRemaining:    daysRemaining,
		HoursRemaining:   hoursRemaining,
		TimeRemaining:    timeRemaining,
		EndDate:          campaign.EndDate,
		CountdownTimer:   countdownTimer,
	}
}

// Function to get campaigns with total donations
// func getAlokasiDana() ([]Campaign, error) {
//     var campaigns []Campaign

//     // Execute the SQL query
//     query := `
//     SELECT 
//         campaigns.id, 
//         campaigns.title, 
//         campaigns.category, 
//         campaigns.target_amount, 
//         COALESCE(SUM(donasis.jumlah), 0) AS total_donations
//     FROM 
//         campaigns
//     LEFT JOIN 
//         donasis ON donasis.campaign_id = campaigns.id
//     GROUP BY 
//         campaigns.id, campaigns.title, campaigns.category, campaigns.target_amount;
//     `

//     // Use the database connection to query the data
//     rows, err := db.Raw(query).Rows()  // Assuming db is your *gorm.DB instance
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()

//     // Loop through the rows and populate the campaigns slice
//     for rows.Next() {
//         var campaign Campaign
//         // Scan the result into the Campaign struct
//         if err := rows.Scan(&campaign.ID, &campaign.Title, &campaign.Category, &campaign.TargetAmount, &campaign.TotalDonations); err != nil {
//             return nil, err
//         }
//         campaigns = append(campaigns, campaign)
//     }

//     // Check for errors after iterating through the rows
//     if err := rows.Err(); err != nil {
//         return nil, err
//     }

//     return campaigns, nil
// }

// getUrgentCampaignsPage renders the page with urgent campaigns
func getUrgentCampaignsPage(w http.ResponseWriter, r *http.Request) {
	campaigns, err := getUrgentCampaigns()
	if err != nil {
		http.Error(w, "Gagal mengambil data kampanye mendesak: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	var displayCampaigns []CampaignDisplayData
	for _, campaign := range campaigns {
		displayCampaigns = append(displayCampaigns, prepareCampaignDisplayData(campaign))
	}
	
	data := struct {
		UrgentCampaigns []CampaignDisplayData
	}{
		UrgentCampaigns: displayCampaigns,
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

// getHomePage renders the home page with urgent campaigns section
func getHomePage(w http.ResponseWriter, r *http.Request) {
	campaigns, err := getUrgentCampaigns()
	if err != nil {
		http.Error(w, "Gagal mengambil data kampanye mendesak: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	var urgentCampaigns []CampaignDisplayData
	for _, campaign := range campaigns {
		urgentCampaigns = append(urgentCampaigns, prepareCampaignDisplayData(campaign))
	}
	
	// Get other data for homepage
	var totalDonatur int64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Count(&totalDonatur)
	
	var totalDonasi float64
	db.Model(&Donasi{}).Where("status_pembayaran = ?", "selesai").Select("COALESCE(SUM(jumlah), 0)").Scan(&totalDonasi)
	
	var campaignCount int64
	db.Model(&Campaign{}).Count(&campaignCount)
	
	data := struct {
		TotalDonatur      int64
		TotalDonasi       string
		KampanyeAktif     int64
		TingkatPenyaluran string
		UrgentCampaigns   []CampaignDisplayData
	}{
		TotalDonatur:      totalDonatur,
		TotalDonasi:       formatRupiah(totalDonasi),
		KampanyeAktif:     campaignCount,
		TingkatPenyaluran: "97%",
		UrgentCampaigns:   urgentCampaigns,
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

// REST API endpoint to get urgent campaigns
func getUrgentCampaignsAPI(w http.ResponseWriter, r *http.Request) {
	campaigns, err := getUrgentCampaigns()
	if err != nil {
		http.Error(w, "Gagal mengambil data kampanye mendesak: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	var displayCampaigns []CampaignDisplayData
	for _, campaign := range campaigns {
		displayCampaigns = append(displayCampaigns, prepareCampaignDisplayData(campaign))
	}
	
	// Set content type and encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(displayCampaigns)
}