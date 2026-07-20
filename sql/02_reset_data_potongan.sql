-- PANTAS V6 - RESET DATA POTONGAN
--
-- CAKUPAN YANG DIHAPUS:
--   1. Periode pelaporan dan seluruh batch impor.
--   2. Seluruh catatan presensi/potongan.
--   3. Seluruh banding dan metadata dokumen banding.
--   4. Notifikasi, antrean email/SMS, dan audit yang terkait potongan/banding.
--
-- CAKUPAN YANG DIPERTAHANKAN:
--   1. Data pegawai, unit, jabatan, password, email, dan nomor HP.
--   2. Akun administrator dan sesi login yang sedang aktif.
--   3. Aturan potongan, kategori alasan banding, dan parameter sistem.
--
-- PENTING:
--   * Buat backup database sebelum menjalankan file ini.
--   * Hentikan sementara service PANTAS di Render agar tidak ada impor,
--     publish, banding, atau worker notifikasi yang berjalan bersamaan.
--   * SQL ini menghapus metadata dokumen banding. Berkas fisik pada bucket
--     Storage "pantas-appeals" tidak dihapus langsung oleh SQL. Daftar path
--     yang perlu dihapus manual ditampilkan pada hasil akhir.
--   * Untuk membuka pengaman, ubah nilai BELUM_DIKONFIRMASI di bawah menjadi
--     RESET_DATA_POTONGAN, lalu jalankan seluruh file di Supabase SQL Editor.

begin;

do $pantas_reset_guard$
declare
  v_confirmation constant text := 'BELUM_DIKONFIRMASI';
begin
  if v_confirmation <> 'RESET_DATA_POTONGAN' then
    raise exception using
      message = 'Reset dibatalkan: ubah BELUM_DIKONFIRMASI menjadi RESET_DATA_POTONGAN terlebih dahulu.';
  end if;

  if to_regclass('public.reporting_periods') is null
     or to_regclass('public.import_batches') is null
     or to_regclass('public.attendance_records') is null
     or to_regclass('public.appeals') is null then
    raise exception 'Reset dibatalkan: skema PANTAS V6 tidak ditemukan atau belum lengkap.';
  end if;
end;
$pantas_reset_guard$;

lock table
  public.reporting_periods,
  public.import_batches,
  public.attendance_records,
  public.appeals,
  public.appeal_items,
  public.appeal_documents,
  public.notifications,
  public.notification_jobs,
  public.audit_logs
in access exclusive mode;

drop table if exists pg_temp.pantas_storage_paths_to_delete;
create temporary table pantas_storage_paths_to_delete
on commit preserve rows
as
select storage_path
from public.appeal_documents
order by storage_path;

-- Hapus pesan yang dibuat oleh proses publish dan banding saja. OTP kontak,
-- OTP reset password, serta notifikasi lain tidak ikut terhapus.
delete from public.notification_jobs
where template_code in (
  'period_published',
  'appeal_submitted',
  'appeal_reviewed'
);

delete from public.notifications
where kind in (
  'period_published',
  'appeal_submitted',
  'appeal_admin_queue',
  'appeal_reviewed'
);

delete from public.audit_logs
where action like 'import.%'
   or action like 'appeal.%'
   or entity_type in ('import_batch', 'appeal', 'appeal_item');

-- Seluruh tabel yang saling terhubung dicantumkan dalam satu TRUNCATE agar
-- foreign key dua arah reporting_periods <-> import_batches tetap aman.
truncate table
  public.appeal_documents,
  public.appeal_items,
  public.appeals,
  public.attendance_records,
  public.import_batches,
  public.reporting_periods
restart identity;

commit;

-- Baris RESET_SELESAI menandakan transaksi berhasil. Jika terdapat baris
-- HAPUS_MANUAL_DI_STORAGE, hapus path tersebut dari bucket pantas-appeals.
select 'RESET_SELESAI'::text as status, null::text as storage_path
union all
select 'HAPUS_MANUAL_DI_STORAGE', storage_path
from pg_temp.pantas_storage_paths_to_delete
order by status, storage_path nulls first;
