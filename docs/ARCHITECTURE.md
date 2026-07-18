# Arsitektur PANTAS

## Komponen

```mermaid
flowchart TD
    B[Browser pegawai] -->|HTTPS, cookie sesi, CSRF| G[Go web service di Render]
    G -->|Postgres/TLS| D[(Supabase PostgreSQL)]
    G -->|service-role server-side| S[(Supabase Storage privat)]
    G -->|antrean notifikasi| D
    G -->|HTTPS| E[Resend / webhook nomor HP]
```

Go menyajikan API dan aset web dari satu binary. Browser tidak menerima connection string database atau service-role key. Seluruh pemeriksaan kepemilikan dan lingkup organisasi dilakukan kembali di backend; menyembunyikan menu pada antarmuka bukan mekanisme otorisasi.

## Matriks akses

| Peran | Data pribadi | Monitoring | Banding yang diverifikasi |
|---|---|---|---|
| Staf | Diri sendiri | Tidak ada | Tidak ada |
| Kepala Seksi/Subbagian | Diri sendiri | Pegawai pada seksi/subbagian secara individual | Staf pada unit yang sama |
| Kepala Bidang/Bagian | Diri sendiri | Tiap seksi/subbagian sebagai agregat | Kepala seksi/subbagian di bawah bidang/bagian |
| Kepala Kantor | Diri sendiri | Bidang/bagian sebagai agregat; Fungsional agregat dengan drill-down | Kepala bidang/bagian dan Fungsional |
| Administrator sistem | Tidak memiliki data pegawai pribadi | Seluruh kantor dan drill-down seluruh unit | Keputusan final semua banding |

Administrator adalah akun sistem pada `admin_accounts`, bukan pegawai pada `users`. Akun pegawai tidak dapat diberi hak admin dari menu pengelolaan pengguna. Semua endpoint admin tetap memerlukan principal bertipe admin.

## Lingkup total monitoring

- Kepala seksi/subbagian: dirinya dan seluruh staf pada unit yang sama.
- Kepala bidang/bagian: pegawai yang berada langsung pada bidang/bagian serta seluruh seksi/subbagian anak.
- Kepala Kantor: seluruh pegawai aktif kantor; administrator sistem dapat melihat cakupan kantor yang sama tanpa ikut dihitung sebagai pegawai.
- Kartu agregat bidang mencakup kepala bidang, semua kepala seksi, dan semua anggota seksi.

## Alur import

```mermaid
flowchart TD
    U[Admin unggah XLSX] --> V[Validasi container, sheet, header, periode, NIP, tanggal]
    V --> R[Hitung komponen potongan dari aturan aktif]
    R --> P[Preview dan staging draft]
    P -->|Admin publikasi| A[Transaksi atomik: aktifkan batch dan periode]
    A --> N[Dashboard berpindah dan notifikasi dijadwalkan]
```

Import baru untuk periode yang sama menghasilkan versi baru. Versi publikasi lama hanya dapat diganti bila belum ada banding pada periode tersebut. File yang sama tidak dapat dipublikasikan dua kali.

## Alur banding

```mermaid
stateDiagram-v2
    [*] --> VerifikasiAtasan: Pegawai mengajukan per tanggal
    VerifikasiAtasan --> KeputusanAdmin: Semua tanggal telah diverifikasi
    KeputusanAdmin --> Selesai: Semua tanggal diputuskan
    Selesai --> [*]
```

Setiap tanggal memiliki status, komentar, dokumen, dan potongan hasil sendiri. Penolakan admin mempertahankan potongan awal. Persetujuan admin dapat mengubah potongan menjadi nol atau nilai yang lebih kecil dari potongan awal.

## Model data utama

- `accounts`, `admin_accounts`: identitas autentikasi dan administrator non-pegawai.
- `units`, `users`, `user_assignment_history`: struktur organisasi, pegawai, dan mutasi.
- `reporting_periods`, `import_batches`, `attendance_records`: data bulanan berversi.
- `appeals`, `appeal_items`, `appeal_documents`: proses banding per hari.
- `deduction_rules`, `appeal_reason_categories`, `parameters`: konfigurasi admin.
- `sessions`, `login_attempts`, `recovery_otps`, `pending_contact_changes`: keamanan akun.
- `notifications`, `notification_jobs`: notifikasi aplikasi dan provider eksternal.
- `audit_logs`: jejak perubahan penting.

View `published_attendance`, `effective_attendance`, dan `monthly_user_summary` memastikan dashboard hanya membaca batch yang sedang dipublikasikan dan otomatis memakai hasil banding yang disetujui.
