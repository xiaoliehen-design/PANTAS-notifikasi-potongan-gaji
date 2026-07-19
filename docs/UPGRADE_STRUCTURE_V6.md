# Upgrade pengelolaan struktur organisasi V6

Panduan ini digunakan untuk database PANTAS V5 yang sudah berjalan dan telah memiliki data pegawai maupun potongan.

## Urutan upgrade

1. Buat backup database Supabase dan pastikan dapat dipulihkan.
2. Hentikan sementara service PANTAS di Render agar tidak ada perubahan pengguna atau import saat migration dijalankan.
3. Buka **Supabase → SQL Editor**, lalu jalankan seluruh isi:

   ```text
   supabase/migrations/004_manage_organization_units.sql
   ```

4. Migration memeriksa struktur lama terlebih dahulu. Jika ada parent, tipe, kode, nama, atau status yang tidak valid, transaksi dibatalkan tanpa perubahan. Perbaiki data yang disebutkan sebelum mencoba lagi.
5. Deploy source V6 di Render.
6. Login sebagai administrator dan buka **Struktur Organisasi**.

Migration ini tidak menghapus atau mengubah pegawai, presensi, potongan, banding, maupun akun administrator.

## Aturan operasional

- Bidang/bagian selalu berada langsung di bawah unit kantor.
- Seksi/subbagian wajib memilih tepat satu bidang/bagian induk.
- Mengubah bidang induk seksi juga mengubah cakupan monitoring dan verifikasi kepala bidang.
- Nama penempatan Excel harus sama dengan nilai penempatan pada workbook. Mode otomatis menghasilkan `Nama Bidang - -` untuk bidang dan `Nama Bidang - Nama Seksi` untuk seksi.
- Unit dengan pegawai aktif atau seksi aktif tidak dapat dinonaktifkan.
- Unit dengan pegawai, unit anak, atau riwayat mutasi tidak dapat dihapus. Nonaktifkan unit bila histori harus dipertahankan.
- Unit kantor dan Fungsional merupakan unit sistem dan tidak dapat diubah atau dihapus dari menu ini.

## Verifikasi

Jalankan melalui SQL Editor:

```sql
select u.code, u.name, u.unit_type, p.name as parent_name, u.is_active
from public.units u
left join public.units p on p.id = u.parent_id
order by coalesce(p.sort_order, u.sort_order), u.unit_type, u.sort_order;

select count(*) as seksi_tanpa_bidang
from public.units s
left join public.units p on p.id = s.parent_id and p.unit_type = 'division'
where s.unit_type = 'section' and p.id is null;
```

Hasil `seksi_tanpa_bidang` harus `0`.
