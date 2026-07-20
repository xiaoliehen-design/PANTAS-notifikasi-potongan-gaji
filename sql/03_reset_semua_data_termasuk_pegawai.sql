-- PANTAS V6 - RESET SEMUA DATA TERMASUK PEGAWAI - VERSI SIAP JALANKAN
--
-- PERINGATAN: file ini tidak memerlukan perubahan kode konfirmasi.
-- Menjalankan seluruh file akan langsung menghapus seluruh data pegawai dan
-- data operasional. Buat backup dan hentikan service Render terlebih dahulu.
--
-- CAKUPAN YANG DIHAPUS:
--   1. Seluruh data pegawai beserta password, email, nomor HP, OTP, dan sesi.
--   2. Seluruh data potongan, periode, batch impor, banding, dan dokumen.
--   3. Seluruh notifikasi, antrean pengiriman, login attempt, riwayat
--      penempatan, dan audit log.
--
-- CAKUPAN YANG DIPERTAHANKAN:
--   1. Skema database, function, trigger, view, index, dan RLS.
--   2. Unit organisasi, aturan potongan, kategori alasan, dan parameter.
--   3. Akun administrator terpisah beserta password administrator.
--      Seluruh sesi admin tetap dihapus sehingga admin harus login kembali.
--
-- PENTING:
--   * Buat backup database sebelum menjalankan file ini.
--   * Hentikan sementara service PANTAS di Render.
--   * SQL ini menghapus metadata dokumen banding. Berkas fisik pada bucket
--     Storage "pantas-appeals" tidak dihapus langsung oleh SQL. Daftar path
--     yang perlu dihapus manual ditampilkan pada hasil akhir.
--   * Versi ini sudah dikonfirmasi dan dapat langsung dijalankan seluruhnya
--     melalui Supabase SQL Editor.

begin;

do $pantas_reset_guard$
declare
  v_confirmation constant text := 'RESET_DATA_PEGAWAI_DAN_POTONGAN';
begin
  if v_confirmation <> 'RESET_DATA_PEGAWAI_DAN_POTONGAN' then
    raise exception using
      message = 'Reset dibatalkan: kode konfirmasi internal file tidak sesuai.';
  end if;

  if to_regclass('public.users') is null
     or to_regclass('public.accounts') is null
     or to_regclass('public.admin_accounts') is null
     or to_regclass('public.attendance_records') is null then
    raise exception 'Reset dibatalkan: skema PANTAS V6 tidak ditemukan atau belum lengkap.';
  end if;
end;
$pantas_reset_guard$;

lock table
  public.accounts,
  public.admin_accounts,
  public.users,
  public.user_assignment_history,
  public.sessions,
  public.login_attempts,
  public.recovery_otps,
  public.pending_contact_changes,
  public.reporting_periods,
  public.import_batches,
  public.attendance_records,
  public.appeals,
  public.appeal_items,
  public.appeal_documents,
  public.parameters,
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

-- Parameter sistem dipertahankan. Putuskan referensi updated_by yang mungkin
-- menunjuk akun pegawai sebelum akun pegawai tersebut dihapus.
update public.parameters p
set updated_by = null
where p.updated_by in (
  select a.id
  from public.accounts a
  where a.account_type = 'user'
);

-- Semua tabel berisi data operasional dan identitas pegawai dicantumkan dalam
-- satu TRUNCATE supaya tidak memakai CASCADE yang berisiko menyentuh tabel lain.
truncate table
  public.appeal_documents,
  public.appeal_items,
  public.appeals,
  public.attendance_records,
  public.import_batches,
  public.reporting_periods,
  public.user_assignment_history,
  public.recovery_otps,
  public.pending_contact_changes,
  public.sessions,
  public.login_attempts,
  public.notifications,
  public.notification_jobs,
  public.audit_logs,
  public.users
restart identity;

-- users.id dan accounts.id sengaja dipisahkan sejak skema V5. Setelah seluruh
-- profil pegawai dihapus, bersihkan identitas akun bertipe user. Akun bertipe
-- admin dan baris admin_accounts tetap dipertahankan.
delete from public.accounts
where account_type = 'user'
  and not exists (
    select 1
    from public.admin_accounts aa
    where aa.account_id = accounts.id
  );

update public.admin_accounts
set last_login_at = null;

commit;

-- Baris RESET_SELESAI menandakan transaksi berhasil. Jika terdapat baris
-- HAPUS_MANUAL_DI_STORAGE, hapus path tersebut dari bucket pantas-appeals.
select 'RESET_SELESAI'::text as status, null::text as storage_path
union all
select 'HAPUS_MANUAL_DI_STORAGE', storage_path
from pg_temp.pantas_storage_paths_to_delete
order by status, storage_path nulls first;
