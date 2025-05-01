// static/js/script.js
document.addEventListener('DOMContentLoaded', function() {
    // Load donasi data
    loadDonasi();
    
    // Form submission
    document.getElementById('donasiForm').addEventListener('submit', function(e) {
        e.preventDefault();
        saveDonasi();
    });
});

// Load donasi data from API
function loadDonasi() {
    fetch('/api/donasi')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector('#donasiTable tbody');
            tableBody.innerHTML = '';
            
            data.forEach(donasi => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${donasi.id}</td>
                    <td>${donasi.nama}</td>
                    <td>Rp ${numberWithCommas(donasi.jumlah)}</td>
                    <td>${donasi.kategori}</td>
                    <td>${formatDate(donasi.tanggal_donasi)}</td>
                    <td>
                        <button class="btn btn-sm btn-info" onclick="editDonasi(${donasi.id})">Edit</button>
                        <button class="btn btn-sm btn-danger" onclick="deleteDonasi(${donasi.id})">Hapus</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error:', error));
}

// Save new donasi
function saveDonasi() {
    const donasi = {
        nama: document.getElementById('nama').value,
        email: document.getElementById('email').value,
        jumlah: parseFloat(document.getElementById('jumlah').value),
        kategori: document.getElementById('kategori').value,
        keterangan: document.getElementById('keterangan').value
    };
    
    fetch('/api/donasi', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(donasi),
    })
    .then(response => response.json())
    .then(data => {
        // Reset form
        document.getElementById('donasiForm').reset();
        // Reload donasi list
        loadDonasi();
        alert('Donasi berhasil disimpan!');
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Gagal menyimpan donasi.');
    });
}

// Edit donasi
function editDonasi(id) {
    fetch(`/api/donasi/${id}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('nama').value = data.nama;
            document.getElementById('email').value = data.email;
            document.getElementById('jumlah').value = data.jumlah;
            document.getElementById('kategori').value = data.kategori;
            document.getElementById('keterangan').value = data.keterangan;
            
            // Change form to update mode
            const form = document.getElementById('donasiForm');
            form.onsubmit = function(e) {
                e.preventDefault();
                updateDonasi(id);
            };
        })
        .catch(error => console.error('Error:', error));
}

// Update donasi
function updateDonasi(id) {
    const donasi = {
        nama: document.getElementById('nama').value,
        email: document.getElementById('email').value,
        jumlah: parseFloat(document.getElementById('jumlah').value),
        kategori: document.getElementById('kategori').value,
        keterangan: document.getElementById('keterangan').value
    };
    
    fetch(`/api/donasi/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(donasi),
    })
    .then(response => response.json())
    .then(data => {
        // Reset form
        document.getElementById('donasiForm').reset();
        // Reset form to save mode
        document.getElementById('donasiForm').onsubmit = function(e) {
            e.preventDefault();
            saveDonasi();
        };
        // Reload donasi list
        loadDonasi();
        alert('Donasi berhasil diupdate!');
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Gagal mengupdate donasi.');
    });
}

// Delete donasi
function deleteDonasi(id) {
    if (confirm('Apakah Anda yakin ingin menghapus donasi ini?')) {
        fetch(`/api/donasi/${id}`, {
            method: 'DELETE',
        })
        .then(response => {
            if (response.ok) {
                loadDonasi();
                alert('Donasi berhasil dihapus.');
            } else {
                throw new Error('Gagal menghapus donasi.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert(error.message);
        });
    }
}

// Helper functions
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('id-ID', { 
        day: 'numeric', 
        month: 'long', 
        year: 'numeric' 
    });
}

function numberWithCommas(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
}