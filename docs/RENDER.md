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
| `RESEND_API_KEY` | API key Resend |
| `EMAIL_FROM` | Pengirim dari domain yang sudah diverifikasi |

Administrator disimpan pada `admin_accounts`, bukan `users`, sehingga tidak memerlukan NIP dan tidak ikut dalam monitoring pegawai. Username dinormalisasi menjadi huruf kecil. Password environment hanya digunakan saat akun pertama kali dibuat dan tidak menimpa password yang sudah diganti melalui aplikasi.

## Environment opsional/default

| Key | Default | Keterangan |
|---|---:|---|
| `PHONE_OTP_WEBHOOK_URL` | kosong | Endpoint SMS/WhatsApp internal |
| `PHONE_OTP_WEBHOOK_TOKEN` | kosong | Bearer token webhook |
| `TRUST_PROXY` | `true` | Memakai IP pertama `X-Forwarded-For` dari Render |
| `COOKIE_SECURE` | `true` | Wajib `true` pada HTTPS produksi |
| `SESSION_TTL` | `12h` | Masa maksimum sesi |
| `SESSION_IDLE_TTL` | `2h` | Batas tidak aktif |
| `MAX_EXCEL_BYTES` | `20971520` | 20 MB |
| `MAX_DOCUMENT_BYTES` | `5242880` | 5 MB |
| `WORKER_INTERVAL` | `5s` | Interval worker notifikasi |

## Urutan deployment yang aman

1. Jalankan migration `001`, migration `002`, lalu seed Supabase.
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
set password_hash = crypt('Password-Sementara-2026!', gen_salt('bf', 12)),
    must_change_password = true
where username = 'admin.pantas';

update public.sessions
set revoked_at = now()
where user_id = (
  select account_id from public.admin_accounts where username = 'admin.pantas'
);
```

Ganti username dan password sementara pada contoh tersebut. Jangan menyimpan query yang sudah berisi password nyata.

## Email

PANTAS mengantrekan email di database lalu worker pada proses Go mengirimkannya melalui Resend. Email hanya dikirim kepada akun yang emailnya telah diverifikasi.

Sebelum go-live:

- verifikasi domain pengirim di Resend;
- gunakan alamat seperti `PANTAS <notifikasi@domain.go.id>`;
- uji OTP, publikasi periode, pengajuan banding, dan hasil banding;
- pantau baris `failed` pada `public.notification_jobs`.

## Webhook nomor HP

Jika diaktifkan, PANTAS mengirim `POST` JSON:

```json
{
  "to": "+628123456789",
  "message": "pesan OTP yang sudah menjadi teks biasa",
  "template": "contact_otp"
}
```

Header `Authorization: Bearer <PHONE_OTP_WEBHOOK_TOKEN>` ditambahkan bila token diisi. Provider harus mengembalikan status 2xx. Tanpa webhook, fitur pemulihan lewat nomor HP tidak dapat mengirim OTP; pemulihan email tetap berfungsi.

## Troubleshooting

- **Service gagal startup:** lihat log `configuration invalid`; periksa `APP_URL`, `APP_SECRET`, Supabase URL/key, dan `DATABASE_URL`.
- **Database unavailable:** pastikan password benar, connection string memakai pooler port 5432, dan `sslmode=require`.
- **Health check 503:** migration belum selesai atau database tidak dapat dijangkau.
- **Dokumen gagal upload:** cek bucket, `SUPABASE_URL`, dan service-role key.
- **Email berstatus failed:** cek `RESEND_API_KEY`, domain pengirim, dan isi `last_error` pada `notification_jobs`.
- **403 invalid_origin:** nilai `APP_URL` tidak sama dengan origin URL yang dibuka pengguna.
