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
