# Konfigurasi Render

## Deployment Blueprint

1. Buat repository GitHub **private** dan pastikan `render.yaml` berada di root.
2. Di Render pilih **New → Blueprint**.
3. Hubungkan repository dan branch `main`.
4. Render membaca `render.yaml`, membangun `Dockerfile`, dan membuat web service `pantas-priok` di region Singapore.
5. Isi environment yang diminta, lalu jalankan deploy.

Blueprint menggunakan plan `starter`, health check `/healthz`, dan auto-deploy setelah seluruh pemeriksaan CI GitHub lulus. Plan dapat disesuaikan dengan kebijakan dan beban kantor.

## Environment wajib

| Key | Nilai |
|---|---|
| `APP_ENV` | `production` |
| `APP_URL` | URL HTTPS final, tanpa `/` di belakang |
| `APP_SECRET` | Dibuat otomatis oleh Blueprint; minimal 32 karakter acak |
| `DATABASE_URL` | Supabase Session pooler port 5432 + `sslmode=require` |
| `SUPABASE_URL` | `https://PROJECT_REF.supabase.co` |
| `SUPABASE_SERVICE_ROLE_KEY` | Secret/service-role key, server-side saja |
| `SUPABASE_STORAGE_BUCKET` | `pantas-appeals` |
| `BOOTSTRAP_ADMIN_USERNAME` | Username admin, misalnya `admin.pantas` |
| `BOOTSTRAP_ADMIN_PASSWORD` | Password awal kuat, 12–128 karakter dan minimal tiga jenis karakter |
| `BOOTSTRAP_ADMIN_NAME` | Nama tampilan administrator |
| `EMAIL_PROVIDER` | `smtp` untuk Gmail tanpa domain atau `resend` untuk domain terverifikasi |
| `EMAIL_FROM` | Nama dan alamat pengirim; untuk Gmail harus memakai akun Gmail yang diautentikasi |

Administrator disimpan pada `admin_accounts`, bukan `users`, sehingga tidak memerlukan NIP dan tidak ikut dalam monitoring pegawai. Username dinormalisasi menjadi huruf kecil. Password environment hanya digunakan saat akun pertama kali dibuat dan tidak menimpa password yang sudah diganti melalui aplikasi.

## Environment opsional/default

| Key | Default | Keterangan |
|---|---:|---|
| `SMTP_HOST` | kosong | Untuk Gmail isi `smtp.gmail.com` |
| `SMTP_PORT` | `587` | Port STARTTLS Gmail |
| `SMTP_USERNAME` | kosong | Alamat Gmail khusus PANTAS |
| `SMTP_PASSWORD` | kosong | App Password Google 16 karakter |
| `SMTP_TLS_MODE` | `starttls` | Gunakan `starttls` untuk port 587; `implicit` untuk port 465 |
| `RESEND_API_URL` | `https://api.resend.com/emails` | Endpoint alternatif Resend |
| `PHONE_PROVIDER` | `auto` | Pilih `twilio` atau `webhook` untuk OTP nomor HP dan notifikasi publikasi melalui SMS |
| `TWILIO_ACCOUNT_SID` | kosong | Account SID Twilio |
| `TWILIO_API_KEY` / `TWILIO_API_SECRET` | kosong | Kredensial yang dianjurkan untuk produksi |
| `TWILIO_AUTH_TOKEN` | kosong | Alternatif untuk uji coba; jangan dipublikasikan |
| `TWILIO_MESSAGING_SERVICE_SID` | kosong | Messaging Service SID; dapat diganti `TWILIO_FROM_NUMBER` |
| `PHONE_OTP_WEBHOOK_URL` | kosong | Alternatif endpoint SMS/WhatsApp internal |
| `PHONE_OTP_WEBHOOK_TOKEN` | kosong | Bearer token webhook |
| `TRUST_PROXY` | `true` | Memakai IP pertama `X-Forwarded-For` dari Render |
| `COOKIE_SECURE` | `true` | Wajib `true` pada HTTPS produksi |
| `SESSION_TTL` | `12h` | Masa maksimum sesi |
| `SESSION_IDLE_TTL` | `30m` | Batas tidak aktif sebelum wajib login kembali |
| `MAX_EXCEL_BYTES` | `20971520` | 20 MB |
| `MAX_DOCUMENT_BYTES` | `5242880` | 5 MB |
| `WORKER_INTERVAL` | `5s` | Interval worker notifikasi |

## Urutan deployment yang aman

1. Jalankan migration `001`, `002`, `003`, lalu seed Supabase.
2. Buat Blueprint Render dan isi semua secret.
3. Deploy pertama.
4. Salin URL `https://...onrender.com`, set sebagai `APP_URL`, lalu redeploy.
5. Buka `/healthz`; respons harus `{"status":"ok","database":"ok"}`.
6. Login memakai username/password bootstrap admin dan segera ganti password awal.
7. Uji import pada satu periode, banding satu hari, verifikasi atasan, dan keputusan admin.

Jika memakai custom domain, ubah `APP_URL` ke custom domain HTTPS yang benar. Origin check PANTAS sengaja menolak request mutasi dari origin lain.

## Reset darurat password admin

Admin tidak memakai fitur lupa password pegawai. Jika password admin terlupa dan tidak ada admin lain, jalankan melalui SQL Editor Supabase lalu login dan segera ganti kembali:

```sql
update public.admin_accounts
set password_hash = extensions.crypt(
      'Password-Sementara-2026!',
      extensions.gen_salt('bf', 12)
    ),
    must_change_password = true
where username = 'admin.pantas';

update public.sessions
set revoked_at = now()
where user_id = (
  select account_id from public.admin_accounts where username = 'admin.pantas'
);
```

Ganti username dan password sementara pada contoh tersebut. Jangan menyimpan query yang sudah berisi password nyata.

## Email Gmail tanpa domain

Untuk kondisi yang hanya mempunyai URL `onrender.com`, gunakan Gmail SMTP. URL Render tetap menjadi `APP_URL` dan tidak boleh dijadikan alamat pengirim email. Buat akun Gmail khusus PANTAS, aktifkan verifikasi dua langkah, lalu buat App Password.

Isi Environment Render:

```env
EMAIL_PROVIDER=smtp
EMAIL_FROM=PANTAS <pantas.notifikasi@gmail.com>
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=pantas.notifikasi@gmail.com
SMTP_PASSWORD=app-password-16-karakter
SMTP_TLS_MODE=starttls
```

Jangan mengisi `SMTP_PASSWORD` dengan password login Gmail. Jangan menaruh App Password di GitHub. Masukkan `EMAIL_FROM` tanpa tanda kutip pembuka/penutup dan gunakan key `SMTP_TLS_MODE` (bukan `SMTP_TLS`). Versi ini tetap menormalisasi format lama agar deployment yang sudah ada tidak langsung gagal. Jika `RESEND_API_KEY` lama masih tersimpan, nilai tersebut boleh tetap ada selama `EMAIL_PROVIDER` diisi `smtp`.

PANTAS meminta provider menerima OTP sebelum formulir kode ditampilkan. Kegagalan dicatat pada `notification_jobs` dan pengguna tetap berada di halaman yang sama. Pesan "Password PANTAS salah" berarti pengiriman belum dicoba karena password yang dimasukkan berbeda dari password login PANTAS.

## SMS nomor HP melalui Twilio

Set `PHONE_PROVIDER=twilio`, lalu isi Account SID, API key/secret, serta Messaging Service SID atau nomor pengirim Twilio. Kanal yang sama dipakai untuk OTP nomor HP dan notifikasi periode yang dipublikasikan. Nomor tujuan PANTAS dinormalisasi menjadi format `+62...`.

```env
PHONE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_API_KEY=SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_API_SECRET=secret-api-key
TWILIO_MESSAGING_SERVICE_SID=MGxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_FROM_NUMBER=
TWILIO_API_BASE_URL=https://api.twilio.com/2010-04-01
```

Untuk uji coba, `TWILIO_AUTH_TOKEN` dapat digunakan sebagai pengganti pasangan API key/secret. Akun trial hanya dapat mengirim ke nomor tujuan yang sudah diverifikasi di Twilio. Detail langkah ada di `docs/EMAIL_PHONE_SETUP.md`.

## Webhook nomor HP

Jika diaktifkan, PANTAS mengirim `POST` JSON:

```json
{
  "to": "+628123456789",
  "message": "pesan OTP atau notifikasi yang sudah menjadi teks biasa",
  "template": "contact_otp atau period_published"
}
```

Header `Authorization: Bearer <PHONE_OTP_WEBHOOK_TOKEN>` ditambahkan bila token diisi. Provider harus mengembalikan status 2xx. Tanpa webhook, tombol penambahan nomor HP akan menampilkan bahwa kanal belum dikonfigurasi; sistem tidak lagi memberi pesan seolah-olah OTP telah terkirim.

## Troubleshooting

- **Service gagal startup:** lihat log `configuration invalid`; periksa `APP_URL`, `APP_SECRET`, Supabase URL/key, dan `DATABASE_URL`.
- **Database unavailable:** pastikan password benar, connection string memakai pooler port 5432, dan `sslmode=require`.
- **Health check 503:** migration belum selesai atau database tidak dapat dijangkau.
- **Dokumen gagal upload:** cek bucket, `SUPABASE_URL`, dan service-role key.
- **Email berstatus failed:** untuk Gmail cek `EMAIL_PROVIDER`, pastikan `EMAIL_FROM` tidak dibungkus tanda kutip, alamat pengirim sama dengan `SMTP_USERNAME`, App Password berisi tepat 16 karakter tanpa spasi, key TLS bernama `SMTP_TLS_MODE`, lalu periksa `last_error`; untuk Resend cek API key/domain.
- **OTP email tidak masuk walau status sent:** periksa folder spam dan aktivitas akun/provider. Status `sent` berarti provider menerima permintaan, bukan jaminan inbox.
- **OTP/notifikasi nomor HP tidak terkirim:** untuk Twilio cek Messaging Logs dan pastikan akun trial sudah memverifikasi nomor tujuan; untuk webhook pastikan endpoint menerima `contact_otp` serta `period_published` dan mengembalikan HTTP 2xx.
- **403 invalid_origin:** nilai `APP_URL` tidak sama dengan origin URL yang dibuka pengguna.
