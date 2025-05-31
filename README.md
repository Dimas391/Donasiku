# 📦 Aplikasi Donasiku
Platform Donasi Komunitas adalah aplikasi berbasis web yang menyediakan sistem lengkap untuk mengelola donasi dan bantuan yang transparan. Dibangun dengan Go, PostgreSQL, dan dikemas dalam Docker untuk memudahkan deployment.

---

## 🚀 Fitur Utama

- ✅ Tampilan total donasi, jumlah donatur, dan penerima manfaat
- ✅ Pencatatan donasi dengan berbagai metode pembayaran
- ✅ Klasifikasi donasi berdasarkan kategori (Pendidikan, Kesehatan, Bencana Alam, Sosial)
- ✅ Grafik distribusi donasi berdasarkan kategori
- ✅ Dukungan berbagai metode (Transfer Bank, eWallet, QRIS)
- ✅ Pencatatan detil alokasi dana dan penerima manfaat

---

## 🚀 Cara Menjalankan

### 🔁 Menggunakan Docker Compose

1. **Clone repository**
   
   Clone repository ke komputer lokal kalian
   ```bah
   git clone https://github.com/Dimas391/Donasiku.git
   ``` 
3. **Masuk ke dalam direktori proyek**
   
   Navigasikan ke folder proyek yang baru saja di-clone
   ```bash
   cd Donasiku
   ```
5. **Build image Docker jika belum dibuild sebelumnya**
   
   Jalankan perintah berikut untuk membangun image Docker dari Dockerfile:
   ```bash
   docker-compose build
   ```
7. **Jalankan aplikasi menggubakan docker**
   
   Jalankan perintah berikut untuk menjalankan aplikasi donasiku
   ```bash
   docker-compose up -d
   ```
9. **Akes browser**
    
   Buka browser untuk menjalankan aplikasi
   ```bash
   http://localhost:8080
   ```
   
## 🛠️ Teknologi yang Digunakan

### 🔧 Backend

- **Go (Golang)** – Bahasa pemrograman backend
- **Gorilla Mux** – HTTP router dan dispatcher untuk REST API
- **GORM** – Object Relational Mapper untuk PostgreSQL
- **PostgreSQL** – Database relasional 

### 🎨 Frontend

- **HTML** – tampilan antarmuka website
- **CSS** - Styling dan layout tampilan website
- **Bootstrap 5** – Framework CSS untuk tampilan responsif
- **JavaScript** – Menambahkan interaktivitas pada sisi client
- **CSS** – Digunakan untuk grafik dan visualisasi data yang ringan

### ⚙️ DevOps

- **Docker** – Untuk containerisasi aplikasi agar lebih portabel
- **Docker Compose** – Mengelola dan menjalankan beberapa container sekaligus
---

## ✍️ Kontribusi

Kami sangat menghargai kontribusi dari anda. Jika anda ingin berkontribusi, ikuti langkah-langkah berikut:

1. Fork repositori ini
2. Buat branch untuk fitur atau perbaikan yang kamu ingin tambahkan (contoh: `git checkout -b fitur/menambahkan-grafik`)
3. Lakukan perubahan yang diperlukan
4. Commit perubahan dengan pesan yang jelas (`git commit -m "Menambahkan grafik distribusi donasi"`)
5. Push ke branchmu (`git push origin fitur/menambahkan-grafik`)
6. Buat pull request ke repositori utama

---

## 📄 Lisensi

Aplikasi ini dilisensikan di bawah **MIT License**. Untuk informasi lebih lanjut, lihat [LICENSE](LICENSE).

---

Jika ada yang perlu ditanyakan atau butuh bantuan lebih lanjut, jangan ragu untuk membuka *issue* di repositori atau menghubungi kami!

---



   


