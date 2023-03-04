# test-crud-user-orders
Proyek ini dibuat untuk menjawab Soal Coding Test.

Proyek ini dibuat menggunakan salah satu contoh dari [repository `awesome-compose`](https://github.com/docker/awesome-compose/) dengan stack `nginx-golang-mariadb-redis` yang bisa anda lihat [di sini](https://github.com/docker/awesome-compose/tree/master/nginx-golang-mysql).

Teknologi yang digunakan pada Proyek ini adalah :
- Nginx : Web Server (proxy)
- MariaDB : SQL Database Server (MySQL)
- GoLang : Bahasa Pemrograman Backend (IDE)
- Echo : Framework yang menggunakan GoLang
- Redis : Cache in-memory
- Docker : Container
- Git : Sistem Pengontrol Versi (Kode)
- Postman : Dokumentasi API (File: `CRUD - User Orders.postman_collection.json`)

## Bagaimana Cara Menjalankan Proyek Ini?
Ada beberapa langkah mudah untuk menjalankan Proyek ini namun pastikan bahwa Komputer anda telah terinstall Docker dan Git.

### 1. Clone Proyek
```
$ mkdir PROJECT && cd PROJECT
$ git clone https://github.com/iFadli/test-crud-user-orders.git
... # tunggu hingga selesai # ...
Cloning into 'test-curd-user-orders'...
remote: Enumerating objects: xxx, done.
remote: Counting objects: 100% (xxx/xxx), done.
remote: Compressing objects: 100% (xxx/xxx), done.
remote: Total xxx (delta xx), reused xxx (delta xx), pack-reused x
Receiving objects: 100% (xxx/xxx), xx.xx KiB | xxx.xx KiB/s, done.
Resolving deltas: 100% (xx/xx), done.
$ cd test-crud-user-orders
```

### 2. * Atur File .env Jika Ingin Merubah DB Config
```
$ nano .env

SERVICE_PORT=8000

DB_HOST=db-hub.docker
DB_NAME=orders
DB_USER=root
DB_PASSWORD=BaGuViX91oo
DB_PORT=3306

REDIS_HOST=redis-hub.docker
REDIS_PORT=6379
REDIS_PASSWORD=

LOG_FILE=ServiceLog dateformat.log
```

### 3. Jalankan Proyek ini dengan docker-compose
```
$ docker-compose up -d
... # tunggu hingga selesai # ...
[+] Running 5/5
 ⠿ Network test-crud-user-orders_default      Created
 ⠿ Volume "test-crud-user-orders_db-data"     Created
 ⠿ Container test-crud-user-orders-redis-1    Started
 ⠿ Container test-crud-user-orders-db-1       Started
 ⠿ Container test-crud-user-orders-backend-1  Started
 ⠿ Container test-crud-user-orders-proxy-1    Started
```

### 4. Akses Proyek sesuai Kebutuhan
Pada pengaturan Docker Proyek ini, secara default akan meng-expose Port 2 Service yang digunakan; Yakni, Database (MariaDB : 3306) dan Web Server (Nginx : 8080).

Jika ingin menghubungkan Database dengan Tool Database Manager seperti DBeaver, anda dapat menyesuaikan konfirugasi dengan File `.env`.
#### !! DBeaver
Berikut langkah-langkah konfigurasi DBeaver :

>1. Buka `DBeaver`.
>2. Klik ikon `Connect to a Database` di Pojok-Kiri-Atas.
>3. Pada Kategori `Popular`, pilih `MariaDB` lalu klik Next.
>4. Isikan Konfigurasi sesuai dengan File `.env` pada Proyek ini.
>5. Jika berhasil, akan ada Ikon centang hijau pada daftar Database di sebelah kiri.
#### !! Postman
Selain pengaturan Database Manager dari luar Docker yang dapat mengakses ke service Database, di Proyek ini juga disematkan Postman Collection untuk memudahkan calon pengguna dalam pengoprasian API ini dengan membaca Dokumentasi API yang ada pada Postman Collection (pada File `CRUD - User Orders.postman_collection.json`.

Berikut cara membukanya :
>1. Pastikan di Desktop kamu telah terinstall [`Postman (Download)`](https://www.postman.com/downloads/).
>2. Download File Collection dari Github [`Postman Collection`](https://github.com/iFadli/test-crud-user-orders/blob/main/backend/CRUD%20-%20User%20Orders.postman_collection.json) atau Buka dari Repository ini yang telah kamu Clone sebelumnya.
>3. Buka Aplikasi Postman, Klik tombol `Import` pada Kiri Atas, Klik `Choose File` dan Pilih File Postman Collection `.json`.
>4. Jika telah berhasil meng-import Collection, akan muncul pada Sisi Kiri aplikasi Postman.
>5. Silahkan Fokuskan pada Direktori `CRUD - User Orders`.
>6. Sebelum mencoba Send Request, jangan lupa untuk mengaktifkan Service-nya terlebih dahulu.

## Daftar API

Sesuai dengan kebutuhan dari Soal Proyek ini, berikut daftar API-nya :
```
GET    /users/
GET    /users/:id
GET    /users/:id/order-histories
POST   /users/
PUT    /users/:id
DELETE /users/:id

GET    /order-items/
GET    /order-items/:id
POST   /order-items/
PUT    /order-items/:id
DELETE /order-items/:id

GET    /order-histories/
GET    /order-histories/:id
POST   /order-histories/
PUT    /order-histories/:id
DELETE /order-histories/:id
```

---
### Daftar Port Aktif
```
8080   Service Nginx
3306   Service MariaDB
```

Lakukan perintah berikut di Terminal untuk Menonaktifkan Proyek :
```
$ docker-compose down --volumes
```

---
### Extra's Info
1. Why clean architecture is good for your application ? if yes, please explain !
   - Seperti informasi pada Konsep dari Clean Architecture yang membagi Logika-Logika yang dibuat kedalam Layer Terpisah; Entity, Use Case, dan Interface.
   - Aplikasi ini dibuat dengan Memisahkan beberapa Logika kedalam Folder-Folder agar mempermudah pengembangan lebih lanjut. Sisi yang dipermudahnya pun bisa dari sisi Penambahan Logika-nya dan juga dari sisi mempermudah Pemahaman Developer Lain yang baru terjun kedalam Projek ini karena Setiap Folder memiliki Fungsinya sendiri.
2. How to scale up your application and when it needs to be ?
   - Scale Up disini, bercabang makna nya. Jika yang dimaksud adalah Scale Up jumlah Instance yang dibuat untuk menjalankan Service ini, maka diperlukan Scale Up ketika Usage Instance (seperti CPU, RAM, DISK, dan NETWORK) telah mencapai 85%. Nilai 85% tidak selalu menjadi nilai Mutlak, nilai ini tergantung Kecepatan DevOps dalam menyiapkan Instance Baru.