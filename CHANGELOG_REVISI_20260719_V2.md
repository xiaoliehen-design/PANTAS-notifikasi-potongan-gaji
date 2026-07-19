# Revisi PANTAS 19 Juli 2026 — Email, Nomor HP, dan Notifikasi

## Perubahan

1. Email OTP kini mendukung dua provider:
   - Gmail/SMTP tanpa domain sendiri;
   - Resend untuk domain yang telah diverifikasi.
2. Pesan kesalahan password pada perubahan kontak dibedakan dari kegagalan provider. Password yang dimaksud adalah password login PANTAS.
3. Input password perubahan kontak memiliki tombol tampil/sembunyikan dan petunjuk yang lebih jelas.
4. OTP nomor HP mendukung Twilio SMS langsung serta webhook SMS/WhatsApp internal.
5. Nomor Indonesia dinormalisasi ke format `+62...` sebelum dikirim.
6. Angka badge notifikasi langsung hilang ketika panel notifikasi dibuka; semua item ditandai dibaca di database.
7. Item notifikasi yang mempunyai tujuan dapat diklik untuk membuka halaman terkait.
8. Permintaan OTP baru membatalkan job OTP lama yang masih menunggu, sehingga kode lama tidak terkirim terlambat.
9. Dokumentasi konfigurasi ditambahkan pada `docs/EMAIL_PHONE_SETUP.md`.

## Environment minimum untuk Gmail

```env
APP_URL=https://pantas-notifikasi-potongan-gaji.onrender.com
EMAIL_PROVIDER=smtp
EMAIL_FROM=PANTAS <pantas.notifikasi@gmail.com>
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=pantas.notifikasi@gmail.com
SMTP_PASSWORD=app-password-google-16-karakter
SMTP_TLS_MODE=starttls
```

## Environment minimum untuk Twilio trial

```env
PHONE_PROVIDER=twilio
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=auth-token-twilio
TWILIO_FROM_NUMBER=+1xxxxxxxxxx
TWILIO_API_BASE_URL=https://api.twilio.com/2010-04-01
```

Nomor tujuan wajib diverifikasi di Twilio selama akun masih trial. Untuk produksi, gunakan API key/secret serta Messaging Service SID.

## Database

Tidak ada migration baru. Struktur database versi sebelumnya tetap digunakan.
