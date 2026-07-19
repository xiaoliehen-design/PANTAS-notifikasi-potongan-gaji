# Pengaturan Email dan Notifikasi Nomor HP

Panduan ini menggunakan URL aplikasi:

```text
https://pantas-notifikasi-potongan-gaji.onrender.com
```

URL tersebut digunakan sebagai `APP_URL`. Subdomain `onrender.com` bukan domain email dan tidak digunakan pada `EMAIL_FROM`.

## A. Gmail SMTP tanpa domain

### 1. Siapkan akun pengirim

Gunakan akun Gmail khusus aplikasi, bukan akun pribadi. Contoh:

```text
pantas.notifikasi@gmail.com
```

Aktifkan **2-Step Verification**, kemudian buka **App Passwords** dan buat App Password bernama `PANTAS Render`. Simpan 16 karakter yang ditampilkan. Jika menu App Password tidak ada, akun mungkin dibatasi administrator atau memakai Advanced Protection; gunakan akun lain yang diizinkan atau Google OAuth/relay resmi.

### 2. Isi Environment Render

Pada service PANTAS pilih **Environment**, lalu isi:

```env
APP_URL=https://pantas-notifikasi-potongan-gaji.onrender.com
EMAIL_PROVIDER=smtp
EMAIL_FROM=PANTAS <pantas.notifikasi@gmail.com>
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=pantas.notifikasi@gmail.com
SMTP_PASSWORD=xxxxxxxxxxxxxxxx
SMTP_TLS_MODE=starttls
```

`EMAIL_FROM` dan `SMTP_USERNAME` sebaiknya menggunakan alamat yang sama. Jangan menambahkan tanda kutip pada kolom value Render. `SMTP_PASSWORD` adalah App Password tanpa spasi, bukan password login Google.

Klik **Save Changes**, kemudian deploy ulang service.

### 3. Uji email

1. Login PANTAS menggunakan akun pegawai.
2. Buka **Profil & Keamanan → Email → Ubah**.
3. Masukkan email tujuan.
4. Masukkan password PANTAS yang sama dengan password login saat ini.
5. Klik **Kirim kode** dan periksa inbox/spam.
6. Masukkan OTP enam digit.

Jika muncul `Password PANTAS salah`, backend belum mencoba mengirim email. Jika muncul `Belum dapat mengirim kode`, periksa Render Logs dan `notification_jobs.last_error`.

## B. SMS nomor HP melalui Twilio

Twilio tidak memerlukan domain web. Kanal ini dipakai untuk OTP verifikasi nomor HP sekaligus notifikasi saat admin mempublikasikan periode. Namun, SMS bersifat berbayar setelah trial dan ketersediaan pengirim dapat bergantung pada negara/operator.

### 1. Siapkan Twilio

1. Buat akun Twilio.
2. Catat **Account SID**.
3. Untuk trial, catat **Auth Token** dan verifikasi nomor HP tujuan pada Twilio Console.
4. Siapkan Twilio Phone Number atau Messaging Service yang dapat mengirim SMS ke Indonesia.
5. Untuk produksi, buat API Key dan API Secret khusus PANTAS; jangan memakai Auth Token utama jika tidak diperlukan.

### 2. Konfigurasi yang dianjurkan

Dengan Messaging Service:

```env
PHONE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_API_KEY=SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_API_SECRET=secret-yang-ditampilkan-saat-key-dibuat
TWILIO_MESSAGING_SERVICE_SID=MGxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_FROM_NUMBER=
TWILIO_API_BASE_URL=https://api.twilio.com/2010-04-01
```

Untuk trial tanpa API key:

```env
PHONE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=auth-token-twilio
TWILIO_API_KEY=
TWILIO_API_SECRET=
TWILIO_MESSAGING_SERVICE_SID=
TWILIO_FROM_NUMBER=+1xxxxxxxxxx
TWILIO_API_BASE_URL=https://api.twilio.com/2010-04-01
```

Nomor `TWILIO_FROM_NUMBER` harus merupakan nomor milik akun Twilio, bukan nomor pribadi Indonesia. Bila memakai `TWILIO_MESSAGING_SERVICE_SID`, `TWILIO_FROM_NUMBER` boleh dikosongkan.

### 3. Uji nomor HP

1. Login PANTAS.
2. Buka **Profil & Keamanan → Nomor HP → Ubah**.
3. Masukkan `0812...` atau `+62812...`.
4. Masukkan password PANTAS.
5. Klik **Kirim kode**.
6. Masukkan kode SMS yang diterima.

Pada akun Twilio trial, nomor tujuan harus diverifikasi lebih dahulu. Untuk seluruh pegawai, upgrade akun dan pastikan sender/messaging service dapat mengirim ke Indonesia.

## C. Notifikasi saat admin mempublikasikan periode

Publikasi yang berhasil dilakukan dalam satu transaksi database. Setelah batch aktif:

1. setiap pegawai aktif memperoleh notifikasi dalam aplikasi;
2. pegawai dengan `email_verified_at` memperoleh job email;
3. pegawai dengan `phone_verified_at` memperoleh job SMS;
4. pegawai yang mempunyai kedua kontak memperoleh email dan SMS;
5. email atau nomor yang belum terverifikasi tidak digunakan.

Email dan SMS hanya menyatakan bahwa data periode tersedia. Nilai atau rincian potongan tidak dicantumkan dan hanya dapat dilihat setelah login ke PANTAS. Upload yang baru berstatus draft tidak mengirim notifikasi; pengiriman baru dijadwalkan setelah admin menekan **Publikasikan**.

## D. Webhook internal sebagai alternatif

Jika kantor memiliki gateway SMS/WhatsApp sendiri, gunakan:

```env
PHONE_PROVIDER=webhook
PHONE_OTP_WEBHOOK_URL=https://gateway-internal.example.go.id/otp
PHONE_OTP_WEBHOOK_TOKEN=token-rahasia
```

PANTAS mengirim JSON:

```json
{
  "to": "+628123456789",
  "message": "pesan OTP dalam teks biasa",
  "template": "contact_otp"
}
```

Gateway harus mengembalikan HTTP 2xx hanya setelah permintaan pengiriman diterima.

Webhook yang dipakai untuk notifikasi periode harus menerima `template` bernilai `period_published`, selain `contact_otp`.

## E. Diagnosis database

Jalankan di Supabase SQL Editor tanpa memilih kolom `payload`, karena payload berisi OTP:

```sql
select created_at, channel, destination, template_code,
       status, attempts, sent_at, last_error
from public.notification_jobs
order by created_at desc
limit 20;
```

- `sent`: provider menerima permintaan.
- `pending`: akan dicoba kembali.
- `failed`: gagal setelah batas percobaan; lihat `last_error`.
- `cancelled`: OTP lama digantikan oleh permintaan yang lebih baru.
- `processing`: sedang dikirim atau menunggu pemulihan lock.
