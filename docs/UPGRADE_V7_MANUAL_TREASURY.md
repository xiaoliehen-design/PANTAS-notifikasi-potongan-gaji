# Upgrade V7: cuti, koreksi manual, dan perbendaharaan

Panduan ini digunakan untuk instalasi PANTAS V6 yang sudah berjalan. Migration V7 mempertahankan pegawai, periode, potongan, banding, dan akun yang sudah ada.

## 1. Backup

Buat backup Supabase sebelum deployment. Pastikan ada jalur restore yang sudah diuji, terutama sebelum memakai fitur hapus data satu periode.

## 2. Jalankan migration

Di Supabase SQL Editor, jalankan seluruh isi:

```text
supabase/migrations/005_manual_deductions_treasury.sql
```

Migration menambahkan peran admin perbendaharaan, metadata sumber/koreksi pada data presensi, tabel riwayat koreksi, aturan cuti/izin, dan membuat ulang tiga view efektif. Seluruh perubahan berada dalam satu transaksi.

Verifikasi:

```sql
select code, label, rate
from public.deduction_rules
where source_field in ('leave', 'status')
order by sort_order;

select column_name
from information_schema.columns
where table_schema = 'public'
  and table_name = 'attendance_records'
  and column_name in ('record_source', 'last_manual_reason_code', 'updated_by', 'updated_at');
```

## 3. Tambahkan environment Render

Isi ketiga nilai berikut. Username dan password harus berbeda dari akun admin sistem.

```env
BOOTSTRAP_TREASURY_USERNAME=perbendaharaan.pantas
BOOTSTRAP_TREASURY_PASSWORD=password-awal-kuat-minimal-12-karakter
BOOTSTRAP_TREASURY_NAME=Admin Perbendaharaan
```

Deploy source V7 atau restart service. Aplikasi membuat akun perbendaharaan saat startup, tanpa menimpa password akun yang sudah pernah dibuat.

## 4. UAT minimum

1. Login admin sistem dan pastikan menu **Kelola Potongan** tampil.
2. Tambahkan satu potongan manual untuk pegawai uji, lalu edit persentasenya dengan ketiga pilihan alasan. Pilihan **Lainnya** wajib memunculkan dan mewajibkan kotak penjelasan.
3. Pastikan data yang sudah memiliki banding tidak dapat dikoreksi manual.
4. Login admin perbendaharaan, pastikan hanya rekap periode berjalan dan profil yang dapat diakses, lalu buka hasil ekspor Excel.
5. Pada database uji, uji hapus data periode. Penghapusan memerlukan alasan dan teks konfirmasi persis, menghapus semua batch/potongan/banding periode tersebut, serta menyisakan audit tindakan tingkat periode.

## 5. Catatan operasional penghapusan

Gunakan **Hapus Data Bulan Ini** hanya untuk kasus workbook yang salah unggah dan setelah backup. Setelah dihapus, periode tidak lagi dipublikasikan dan harus diimpor serta dipublikasikan ulang dari workbook yang benar. Dokumen banding terkait ikut dihapus dari Storage sejauh layanan Storage dapat dijangkau; kegagalan pembersihan file dilaporkan sebagai peringatan agar dapat ditindaklanjuti.
