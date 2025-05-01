// script.js - For donation platform dashboard
document.addEventListener('DOMContentLoaded', function() {
    // Fetch total donation statistics
    fetchTotalDonation();
    
    // Fetch latest donations for the table
    fetchLatestDonations();
    
    // Set up distribution chart data
    fetchDistributionData();
});

// Function to fetch total donation data
function fetchTotalDonation() {
    fetch('/api/donasi/total')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Update the total donation amount
            document.querySelector('.stats-item:nth-child(1) .stats-number').textContent = 
                'Rp ' + formatNumber(data.total);
            
            // Update the total donor count
            document.querySelector('.stats-item:nth-child(2) .stats-number').textContent = 
                data.count;
            
            // Note: The "Penerima Manfaat" (Beneficiaries) count might need
            // a separate API endpoint if you're tracking that
        })
        .catch(error => {
            console.error('Error fetching total donation data:', error);
        });
}

// Function to fetch latest donations for the table
function fetchLatestDonations() {
    fetch('/api/donasi?limit=5')  // Assuming we can limit results, add this parameter to your API
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(donations => {
            const tableBody = document.querySelector('#donasiTable tbody');
            tableBody.innerHTML = ''; // Clear existing rows
            
            donations.forEach(donation => {
                // Skip anonymous donations if anonim is true
                const displayName = donation.anonim ? 'Anonim' : donation.nama;
                
                // Format date
                const donationDate = new Date(donation.tanggal_donasi);
                const formattedDate = donationDate.getDate() + ' ' + 
                                     getMonthNameIndonesian(donationDate.getMonth()) + ' ' +
                                     donationDate.getFullYear();
                
                // Create table row
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${displayName}</td>
                    <td>Rp ${formatNumber(donation.jumlah)}</td>
                    <td>${donation.kategori}</td>
                    <td>${formattedDate}</td>
                `;
                
                tableBody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching latest donations:', error);
        });
}

// Function to fetch distribution data by category (optimized version)
function fetchDistributionData() {
    fetch('/api/donasi/kategori-total')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(categoryData => {
            // Calculate total across all categories
            const totalAmount = categoryData.reduce((sum, item) => sum + item.total, 0);
            
            // Map of expected categories to their display names
            const categoryMap = {
                'Bencana': {index: 0, color: '#ff9a76'},
                'Pendidikan': {index: 1, color: '#78c6a3'},
                'Kesehatan': {index: 2, color: '#93b5e1'}
            };
            
            // Calculate percentages and update the UI
            categoryData.forEach(item => {
                if (categoryMap[item.kategori]) {
                    const percentage = ((item.total / totalAmount) * 100).toFixed(1);
                    const index = categoryMap[item.kategori].index;
                    
                    // Update the percentage in the legend
                    document.querySelector(`.category-item:nth-child(${index + 1}) .category-percent`).textContent = 
                        percentage + '%';
                }
            });
            
            // For a complete implementation, you would regenerate the SVG paths
            // based on the actual percentages here
            updatePieChart(categoryData, totalAmount, categoryMap);
        })
        .catch(error => {
            console.error('Error fetching distribution data:', error);
        });
}

// Function to update SVG pie chart
function updatePieChart(categoryData, totalAmount, categoryMap) {
    // This is a simplified implementation - a complete one would 
    // calculate actual SVG path data for the pie segments
    
    // First, prepare data with percentages and angles
    let startAngle = 0;
    const chartData = [];
    
    categoryData.forEach(item => {
        if (categoryMap[item.kategori]) {
            const percentage = (item.total / totalAmount);
            const endAngle = startAngle + (percentage * 2 * Math.PI);
            
            chartData.push({
                category: item.kategori,
                startAngle: startAngle,
                endAngle: endAngle,
                percentage: percentage,
                color: categoryMap[item.kategori].color
            });
            
            startAngle = endAngle;
        }
    });
    
    // Generate SVG paths using the calculated angles
    const radius = 150;
    const centerX = 200;
    const centerY = 200;
    
    // Update each path in the pie chart
    chartData.forEach((slice, index) => {
        // Calculate SVG arc path
        const startX = centerX + radius * Math.sin(slice.startAngle);
        const startY = centerY - radius * Math.cos(slice.startAngle);
        const endX = centerX + radius * Math.sin(slice.endAngle);
        const endY = centerY - radius * Math.cos(slice.endAngle);
        
        // Determine if the arc should be drawn as a large arc (more than 180 degrees)
        const largeArcFlag = slice.percentage > 0.5 ? 1 : 0;
        
        // Create SVG path string
        const pathData = `M${centerX},${centerY} L${startX},${startY} A${radius},${radius} 0 ${largeArcFlag},1 ${endX},${endY} Z`;
        
        // Find the appropriate SVG path element using appropriate selector
        // This assumes the paths are in the same order as our categories
        const paths = document.querySelectorAll('svg path');
        const pathElement = paths[index];
        
        if (pathElement) {
            pathElement.setAttribute('d', pathData);
        }
    });
}

// Helper function to format numbers with thousand separators
function formatNumber(number) {
    return number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
}

// Helper function to get Indonesian month names
function getMonthNameIndonesian(monthIndex) {
    const monthNames = [
        'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
        'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    return monthNames[monthIndex];
}

