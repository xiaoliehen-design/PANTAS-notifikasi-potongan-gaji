# PANTAS

PANTAS — **Pemantauan Absensi dan Tunjangan Secara Akumulatif** — adalah aplikasi internal KPU Bea dan Cukai Tipe A Tanjung Priok untuk menyajikan potongan tukin secara privat, ringkas, dan berjenjang.

Backend ditulis dengan Go, data utama disimpan di PostgreSQL Supabase, dokumen banding disimpan pada bucket Supabase Storage privat, dan aplikasi siap dibangun sebagai satu container di Render.

> **Penting:** seed pegawai memuat nama dan NIP. Gunakan repository GitHub **private**, batasi akses project Supabase/Render, dan jangan mengirim service-role key melalui chat atau email.

## Fitur

- Pegawai login memakai NIP; password awal sama dengan NIP dan wajib diganti pada login pertama.
- Administrator merupakan akun sistem terpisah, login memakai username dan password awal dari environment Render, serta tidak tercatat sebagai pegawai.
- Profil untuk mengubah password serta memverifikasi email dan nomor HP pemulihan.
- Dashboard otomatis mengikuti periode terbaru yang dipublikasikan admin.
- Potongan berjalan, detail per tanggal, grafik riwayat default 12 bulan, dan filter rentang bulan/tahun.
- Ringkasan hari kerja (`P=1`, `M=1`, `PM=2`), lembur (`L1=1`, `L2=1`), cuti, dan off.
- Banding per hari potongan dengan kategori alasan dan dokumen PDF/JPG/PNG opsional.
- Verifikasi atasan per hari, lalu keputusan final admin per hari beserta komentar opsional.
- Monitoring berjenjang untuk kepala seksi/subbagian, kepala bidang/bagian, Kepala Kantor, dan admin.
- Total potongan pada cakupan kewenangan dan agregat kantor untuk admin.
- Peringatan individu dan unit dengan parameter yang hanya dapat diubah admin.
- Admin dapat menambah, memindahkan, menonaktifkan/menghapus, dan mereset password pengguna.
- Admin dapat menambah dan mengubah aturan potongan langsung dari Parameter Sistem.
- Import workbook bulanan dengan validasi format, NIP, tanggal, duplikasi, unit, hash file, staging, dan publikasi atomik.
- Notifikasi dalam aplikasi dan email setelah periode dipublikasikan; angka notifikasi dihapus ketika panel notifikasi dibuka.
- Pengiriman OTP email/nomor HP diverifikasi langsung terhadap respons provider dan kegagalan tidak memutus sesi pengguna.
- Email mendukung Gmail SMTP tanpa domain maupun Resend dengan domain terverifikasi; nomor HP mendukung Twilio SMS atau webhook internal.
- Audit log untuk perubahan penting.

## Data awal dari workbook

File `supabase/seed/002_employees_from_reference.sql` dibentuk dari workbook rekapitulasi yang diberikan dan berisi:

- 1.123 pegawai;
- 59 unit: 1 kantor, 11 bidang/bagian, 46 seksi/subbagian, dan 1 kelompok fungsional;
- Adhang Noegroho Adhi sebagai Kepala Kantor;
- kepala bidang/bagian sebagai pegawai pertama pada kelompok `Bidang/Bagian - -`;
- kepala seksi/subbagian sebagai pegawai pertama pada setiap kelompok seksi/subbagian;
- 296 pegawai pada kelompok Fungsional.

Contoh yang diminta juga terpetakan: Heru Prayitno sebagai Kepala Bagian Umum dan Misnawi sebagai Kepala Subbagian Dukungan Teknis.

## Mulai cepat: Supabase → GitHub → Render

1. Buat project Supabase baru.
2. Di **SQL Editor**, jalankan berurutan:
   - `supabase/migrations/001_pantas_schema.sql`
   - `supabase/migrations/002_separate_admin_accounts.sql`
   - `supabase/migrations/003_fix_pgcrypto_schema.sql`
   - `supabase/seed/002_employees_from_reference.sql`
3. Salin connection string **Session pooler port 5432** dan tambahkan `sslmode=require` bila belum ada.
4. Buat repository GitHub **private**, lalu unggah seluruh isi folder ini ke root repository.
5. Di Render pilih **New → Blueprint**, hubungkan repository, dan gunakan `render.yaml`.
6. Isi seluruh secret yang ditandai `sync: false`; petunjuk lengkap ada di [docs/RENDER.md](docs/RENDER.md).
7. Setelah URL Render tersedia, set `APP_URL` ke URL HTTPS tersebut lalu redeploy.
8. Login dengan `BOOTSTRAP_ADMIN_USERNAME` dan `BOOTSTRAP_ADMIN_PASSWORD`, lalu segera ganti password awal.
9. Buka **Import Data**, unggah workbook bulanan, periksa preview, lalu klik **Publikasikan**.

Petunjuk terperinci:

- [Konfigurasi Supabase](docs/SUPABASE.md)
- [Konfigurasi Render](docs/RENDER.md)
- [Pengaturan email dan OTP nomor HP](docs/EMAIL_PHONE_SETUP.md)
- [Upgrade administrator non-pegawai](docs/UPGRADE_ADMIN.md)
- [Format workbook import](docs/EXCEL_FORMAT.md)
- [Arsitektur dan matriks akses](docs/ARCHITECTURE.md)
- [Keamanan dan operasi](docs/SECURITY.md)

## Format Excel import

Importer sengaja ketat agar hasil tidak bergeser diam-diam:

- satu sheet bernama tepat `DETAIL WFH WFO`;
- periode pada `B2`, misalnya `16 Juni 2026 s.d. 15 Juli 2026`;
- header tepat pada `A4:O4`:

```text
Tanggal | Nama | NIP | Bidang | Locus Penempatan | Jam Masuk | Jam Pulang |
TL | PSW | Shift | Status | Cuti | Penugasan | Konfirmasi | Keterangan
```

File contoh yang diberikan memiliki bagian akhir container XLSX yang tidak lengkap, tetapi sheet utama utuh. Importer memiliki mode pemulihan khusus, lalu tetap memvalidasi seluruh header, tanggal, NIP, dan baris sebelum data dapat dipublikasikan. Detail ada di [docs/EXCEL_FORMAT.md](docs/EXCEL_FORMAT.md).

## Menjalankan lokal

Persyaratan: Go 1.26.5 dan project Supabase yang sudah dimigrasikan.

```bash
cp .env.example .env
# isi nilai .env, lalu ekspor dengan cara yang sesuai shell Anda
go run ./cmd/server
```

Aplikasi tersedia di `http://localhost:10000`. Untuk lokal, ubah `APP_ENV=development`, `APP_URL=http://localhost:10000`, dan `COOKIE_SECURE=false`.

## Pengujian

```bash
make check
```

Untuk menguji langsung workbook contoh yang diberikan:

```bash
PANTAS_SAMPLE_XLSX='/path/Upload dokumen.xlsx' go test ./internal/importer -run TestSuppliedWorkbookWhenConfigured -v
```

CI GitHub menjalankan pemeriksaan format, `go vet`, seluruh test, dan build pada setiap pull request. Build menonaktifkan penyisipan metadata VCS agar hasil konsisten di GitHub Actions dan Docker.

## Struktur repository

```text
cmd/server/                  entrypoint aplikasi
internal/auth/               sesi, password, OTP, CSRF, rate limit
internal/httpapi/            API, otorisasi, dashboard, banding, admin
internal/importer/           parser XLSX, perhitungan, staging, publikasi
internal/mailer/             antrean Gmail/Resend dan Twilio/webhook nomor HP
internal/storage/            akses bucket privat Supabase
supabase/migrations/         skema database dan default parameter
supabase/seed/               unit, nama, NIP, dan jabatan awal
web/static/                  antarmuka responsif tanpa build Node.js
docs/                        panduan deployment dan operasi
```

## Catatan produksi

Paket ini sudah melewati build, `go vet`, unit test parser, pengujian workbook contoh (33.690 baris aktif), dan eksekusi migration+seed pada runtime PostgreSQL kompatibel. Sebelum digunakan untuk data tukin riil, lakukan UAT dengan perwakilan pegawai/atasan/admin, verifikasi kembali penetapan pejabat hasil inferensi workbook, aktifkan backup, dan lakukan penilaian keamanan serta kepatuhan internal.
