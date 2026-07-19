# Pengaturan Email dan Notifikasi Nomor HP

Panduan ini menggunakan URL aplikasi:

```text
https://pantas-notifikasi-potongan-gaji.onrender.com
```

URL tersebut digunakan sebagai `APP_URL`. Subdomain `onrender.com` bukan domain email dan tidak digunakan pada `EMAIL_FROM`.

## A. Render Free: Brevo API melalui HTTPS (disarankan)

[Render Free memblokir trafik keluar pada port 25, 465, dan 587](https://render.com/docs/free). Karena Gmail SMTP memakai port 465/587, konfigurasi host, TLS, username, dan App Password yang benar tetap akan berakhir dengan timeout. Gunakan provider transactional email berbasis HTTPS pada port 443.

### 1. Siapkan sender Brevo

1. Buat akun Brevo, ikuti [panduan transactional email API](https://developers.brevo.com/docs/send-a-transactional-email), dan buka **Settings → Senders, Domains & Dedicated IPs**.
2. Tambahkan alamat email khusus PANTAS sebagai sender.
3. Masukkan kode verifikasi yang dikirim Brevo ke alamat tersebut.
4. Buka **SMTP & API → API Keys**, buat API key khusus PANTAS, lalu simpan nilainya langsung di secret Environment Render.

Alamat sender dapat berupa Gmail. Untuk deliverability jangka panjang, domain kantor yang diautentikasi tetap lebih baik.

### 2. Isi Environment Render

```env
EMAIL_PROVIDER=brevo
EMAIL_FROM=PANTAS <alamat-sender-yang-sudah-diverifikasi@example.com>
BREVO_API_KEY=api-key-brevo-baru
BREVO_API_URL=https://api.brevo.com/v3/smtp/email
```

`SMTP_*` boleh tetap ada selama `EMAIL_PROVIDER` diisi `brevo`; nilai SMTP tersebut tidak akan digunakan. Klik **Save Changes**, lalu **Manual Deploy → Deploy latest commit**.

### 3. Uji email

1. Login PANTAS menggunakan akun pegawai.
2. Buka **Profil & Keamanan → Email → Ubah**.
3. Masukkan email tujuan dan password PANTAS saat ini.
4. Klik **Kirim kode**, lalu periksa inbox, spam, serta **Brevo → Transactional → Logs**.
5. Masukkan OTP enam digit.

Jika API key ditolak atau sender belum diverifikasi, PANTAS menampilkan diagnosis yang aman dan spesifik. Respons lengkap provider tetap disimpan pada `notification_jobs.last_error` untuk administrator.

## B. Gmail SMTP pada Render berbayar atau lokal

Gunakan akun Gmail khusus aplikasi, aktifkan **2-Step Verification**, lalu buat App Password. App Password adalah 16 karakter yang berbeda dari password login Gmail.

Isi environment berikut hanya jika layanan hosting mengizinkan koneksi keluar SMTP:

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

Konfigurasi ini tidak akan bekerja pada Render Free karena port 587 diblokir oleh platform. Jika muncul kegagalan autentikasi pada layanan yang mengizinkan SMTP, buat App Password baru dan pastikan alamat `EMAIL_FROM` sama dengan akun Gmail yang diautentikasi.

## C. SMS nomor HP melalui Twilio

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

## D. Notifikasi saat admin mempublikasikan periode

Publikasi yang berhasil dilakukan dalam satu transaksi database. Setelah batch aktif:

1. setiap pegawai aktif memperoleh notifikasi dalam aplikasi;
2. pegawai dengan `email_verified_at` memperoleh job email;
3. pegawai dengan `phone_verified_at` memperoleh job SMS;
4. pegawai yang mempunyai kedua kontak memperoleh email dan SMS;
5. email atau nomor yang belum terverifikasi tidak digunakan.

Email dan SMS hanya menyatakan bahwa data periode tersedia. Nilai atau rincian potongan tidak dicantumkan dan hanya dapat dilihat setelah login ke PANTAS. Upload yang baru berstatus draft tidak mengirim notifikasi; pengiriman baru dijadwalkan setelah admin menekan **Publikasikan**.

## E. Webhook internal sebagai alternatif

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

## F. Diagnosis database

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
