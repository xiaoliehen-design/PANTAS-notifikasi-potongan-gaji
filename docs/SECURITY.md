# Keamanan dan operasi

PANTAS menangani NIP, presensi, potongan tukin, alasan banding, dan dokumen pendukung. Seluruhnya harus diperlakukan sebagai data internal sensitif.

## Kontrol yang diterapkan

- Tidak ada pendaftaran publik; akun hanya dibuat admin.
- Password awal NIP hanya berlaku sampai pengguna dipaksa menggantinya.
- Password setelah perubahan di-hash dengan bcrypt melalui `pgcrypto` cost 12.
- Sesi memakai token acak 256-bit yang hanya disimpan sebagai hash di database.
- Cookie sesi `HttpOnly`, `Secure`, dan `SameSite=Lax` pada produksi.
- CSRF double-submit yang juga dibandingkan dengan hash sesi.
- Pemeriksaan same-origin untuk request mutasi.
- Rate limit login per NIP dan IP; OTP berlaku 10 menit dan dibatasi lima percobaan.
- Reset password dan perubahan kontak mencabut/memverifikasi kredensial sesuai alur.
- Pemeriksaan lingkup organisasi dilakukan di setiap endpoint monitoring, review, dan dokumen.
- Dokumen diperiksa dari isi aktual, bukan hanya ekstensi; tipe dibatasi PDF/JPG/PNG dan ukuran default 5 MB.
- Bucket Storage privat; objek tidak dibagikan sebagai URL publik.
- Header CSP, HSTS, anti-frame, no-sniff, referrer, dan permissions policy.
- Perubahan admin, import, banding, dan mutasi penting dicatat pada audit log.
- Penghapusan pengguna adalah soft delete agar integritas riwayat tetap terjaga.

## Risiko password awal NIP

Penggunaan NIP sebagai password awal mengikuti kebutuhan yang diberikan, tetapi NIP bukan secret. Mitigasi yang sudah diterapkan adalah wajib ganti password sebelum endpoint lain dapat digunakan, rate limit login, dan kemampuan admin mencabut sesi/reset akun.

Untuk go-live, berikan akses secara bertahap dan minta pegawai melakukan login pertama dalam jangka waktu terbatas. Bila kebijakan memungkinkan di masa depan, ganti password awal NIP dengan token aktivasi acak satu kali.

## Checklist go-live

- Repository GitHub private dan aksesnya ditinjau.
- Penetapan 58 pejabat hasil urutan workbook diverifikasi oleh SDM/admin.
- Hanya personel berwenang yang memiliki akses Supabase/Render/Resend.
- Service-role key dan database password hanya berada pada secret environment.
- `APP_URL` HTTPS benar; `COOKIE_SECURE=true`.
- Backup dan uji restore selesai.
- UAT role staf, seluruh tingkat atasan, Kepala Kantor, dan admin selesai.
- Uji bahwa satu pegawai tidak dapat membaca data/dokumen pegawai lain.
- Uji bahwa kepala seksi/bidang tidak dapat memperluas scope dengan mengganti URL.
- Uji email dan webhook nomor HP dengan data non-produksi.
- Kebijakan retensi dokumen banding, audit log, dan data presensi ditetapkan.
- Prosedur insiden dan rotasi key ditetapkan.

## Respons insiden

Jika database password atau service-role key diduga bocor:

1. cabut/rotasi secret di Supabase;
2. perbarui environment Render dan redeploy;
3. cabut seluruh sesi aktif dengan `update public.sessions set revoked_at=now() where revoked_at is null;`;
4. tinjau audit log, log Render, dan objek Storage;
5. ikuti prosedur pelaporan insiden internal.

Jangan memasukkan nilai secret ke issue GitHub, screenshot, log aplikasi, atau dokumen operasional.
