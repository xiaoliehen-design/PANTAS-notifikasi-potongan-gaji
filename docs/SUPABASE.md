# Konfigurasi Supabase

## 1. Buat project

Pilih region yang dekat dengan service Render (contoh: Singapore), buat password database yang kuat, lalu simpan password pada password manager.

## 2. Jalankan SQL

Di **SQL Editor**, jalankan file berikut secara berurutan dan tunggu sampai masing-masing selesai:

1. `supabase/migrations/001_pantas_schema.sql`
2. `supabase/migrations/002_separate_admin_accounts.sql`
3. `supabase/seed/002_employees_from_reference.sql`

Migration pertama membuat skema utama. Migration kedua membuat identitas akun dan tabel administrator yang terpisah dari pegawai, serta memigrasikan referensi audit lama. Seed berikutnya membuat struktur organisasi dan 1.123 akun pegawai.

Seed mengandung data pribadi. Jangan menaruhnya di repository publik.

## 3. Verifikasi hasil

Jalankan di SQL Editor:

```sql
select unit_type, count(*) from public.units group by unit_type order by unit_type;
select position_role, count(*) from public.users group by position_role order by position_role;
select name, nip, position_role from public.users where position_role = 'office_head';
select username, name, is_active from public.admin_accounts;
select id, name, public from storage.buckets where id = 'pantas-appeals';
```

Hasil utama yang diharapkan sebelum aplikasi pertama kali dijalankan: 59 unit, 1.123 pengguna, satu `office_head` bernama Adhang Noegroho Adhi, dan bucket `public=false`. Baris `admin_accounts` dibuat oleh aplikasi saat startup setelah environment bootstrap admin diisi.

## 4. Connection string

Di dashboard Supabase pilih **Connect → Session pooler**. Gunakan port **5432**:

```text
postgresql://postgres.PROJECT_REF:PASSWORD@aws-REGION.pooler.supabase.com:5432/postgres?sslmode=require
```

Simpan sebagai `DATABASE_URL` di Render. PANTAS membatasi pool aplikasi menjadi maksimum 10 koneksi dan memakai simple protocol agar kompatibel dengan pooler.

Jangan menggunakan URL browser/API (`https://...supabase.co`) sebagai `DATABASE_URL`. Jangan menaruh password database di repository.

## 5. URL dan service key

Ambil dari **Project Settings / API**:

- Project URL → `SUPABASE_URL`
- Secret/service-role key → `SUPABASE_SERVICE_ROLE_KEY`

Key tersebut hanya boleh berada pada environment server Render. Jangan memakai anon key dan jangan menanam key ke `app.js`.

Backend menggunakan service-role key hanya untuk objek Storage. Bucket tetap privat; pengguna mengunduh dokumen melalui endpoint Go setelah pemeriksaan kepemilikan/lingkup.

## 6. RLS

Migration mengaktifkan RLS dan mencabut akses tabel/view dari role `anon` dan `authenticated`. Arsitektur PANTAS tidak melakukan query database langsung dari browser. Otorisasi aplikasi dilakukan oleh backend Go dengan sesi PANTAS.

Gunakan connection string role database yang disalin dari dashboard. Jangan membuat policy `select all` untuk mempermudah debugging.

## 7. Backup dan pemeliharaan

- Aktifkan backup/PITR sesuai paket dan kebijakan kantor.
- Uji restore sebelum go-live.
- Rotasi password database dan service-role key bila terpapar.
- Pantau ukuran `attendance_records`, `audit_logs`, dan `notification_jobs`.
- Bersihkan session, OTP, dan login attempt lama dengan job terjadwal bila volume mulai besar.
- Jalankan Security Advisor Supabase setelah migration.

## Mengulang seed

Seed aman dijalankan ulang untuk menyelaraskan nama, unit, jabatan, dan status aktif pegawai yang berasal dari workbook. Seed tidak mereset password yang sudah diganti dan tidak mencabut status admin yang sudah diberikan. Karena seed adalah snapshot awal, setelah go-live mutasi rutin sebaiknya dilakukan dari menu admin, bukan dengan menjalankan ulang seed lama.
