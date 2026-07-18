# Upgrade administrator non-pegawai

Panduan ini digunakan bila migration `001` dan seed pegawai sudah pernah dijalankan serta aplikasi lama sudah berada di Render.

## Urutan wajib

1. Buat backup database Supabase.
2. Buka **Supabase → SQL Editor** dan jalankan seluruh isi:

   ```text
   supabase/migrations/002_separate_admin_accounts.sql
   supabase/migrations/003_fix_pgcrypto_schema.sql
   ```

   Jika migration `002` sudah pernah berhasil dan Render hanya gagal dengan error `gen_salt does not exist`, cukup jalankan migration `003`.

3. Pastikan query selesai tanpa error. Migration mempertahankan seluruh pegawai, presensi, banding, import, dan audit yang sudah ada. Flag admin lama pada pegawai dinonaktifkan.
4. Di **Render → Environment**, tambahkan:

   ```text
   BOOTSTRAP_ADMIN_USERNAME=admin.pantas
   BOOTSTRAP_ADMIN_PASSWORD=<password awal kuat>
   BOOTSTRAP_ADMIN_NAME=Administrator PANTAS
   ```

   Username harus 3–64 karakter, diawali huruf, dan hanya menggunakan huruf kecil, angka, titik, garis bawah, atau tanda hubung. Password harus 12–128 karakter, tidak memuat username, dan memakai minimal tiga jenis karakter.

5. Hapus environment lama `BOOTSTRAP_ADMIN_NIP` bila masih ada. Aplikasi baru mengabaikannya, tetapi penghapusan mencegah kebingungan operasional.
6. Unggah source terbaru ke GitHub dan deploy commit tersebut di Render.
7. Login memakai username dan password awal admin, lalu segera ganti password saat diminta.

## Verifikasi

Jalankan melalui SQL Editor:

```sql
select count(*) as jumlah_pegawai from public.users;

select username, name, is_active, must_change_password
from public.admin_accounts;

select count(*) as pegawai_dengan_flag_admin_lama
from public.users
where is_admin;
```

Hasil yang diharapkan setelah startup aplikasi: 1.123 pegawai, satu baris administrator, dan nol flag admin lama pada pegawai.
