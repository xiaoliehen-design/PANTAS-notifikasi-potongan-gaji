# Format workbook import

## Ketentuan wajib

| Elemen | Nilai yang diterima |
|---|---|
| Ekstensi | `.xlsx` |
| Ukuran default | Maksimum 20 MB, dapat diubah melalui `MAX_EXCEL_BYTES` |
| Jumlah sheet | Tepat satu |
| Nama sheet | `DETAIL WFH WFO` |
| Periode | Cell `B2`, contoh `16 Juni 2026 s.d. 15 Juli 2026` |
| Tanggal cetak | Cell `B3` (disimpan sebagai metadata) |
| Header | `A4:O4`, harus sama persis |
| Data | Mulai baris 5 |

Urutan header:

| Kolom | Header | Contoh isi |
|---|---|---|
| A | Tanggal | tanggal Excel |
| B | Nama | Nama pegawai |
| C | NIP | 18 digit |
| D | Bidang | Bidang/Bagian |
| E | Locus Penempatan | struktur lengkap dari sumber |
| F | Jam Masuk | waktu Excel atau kosong |
| G | Jam Pulang | waktu Excel atau kosong |
| H | TL | `TL1`, `TL2`, `TL3`, `LA`, atau kosong |
| I | PSW | `PSW1`–`PSW4`, `LA`, atau kosong |
| J | Shift | `P`, `M`, `PM`, `L1`, `L2`, `OFF`, atau nilai sumber lain |
| K | Status | antara lain `I`, `TK`, atau kosong |
| L | Cuti | jenis cuti atau kosong |
| M | Penugasan | teks sumber atau kosong |
| N | Konfirmasi | teks sumber atau kosong |
| O | Keterangan | teks sumber atau kosong |

## Aturan potongan awal

Default migration berisi tarif yang ditemukan pada workbook rekapitulasi:

| Sumber | Kode | Tarif |
|---|---|---:|
| TL | TL1 | 1% |
| TL | TL2 | 1,25% |
| TL | TL3 / LA | 2,5% |
| PSW | PSW1 | 0,5% |
| PSW | PSW2 | 1% |
| PSW | PSW3 | 1,25% |
| PSW | PSW4 / LA | 2,5% |
| Cuti | Cuti Alasan Penting Dipotong | 5% |
| Cuti | Cuti Besar Dipotong | 2,5% |
| Cuti | Cuti Sakit Dipotong | 2,5% |
| Status | I / TK | 5% |

Jika satu hari memuat beberapa kode yang memiliki tarif, tarif hari tersebut dijumlahkan dan setiap komponennya disimpan. Admin dapat mengubah label, tarif, atau status aktif aturan sebelum import berikutnya; publikasi lama tidak dihitung ulang.

## Validasi sebelum publikasi

PANTAS menolak publikasi bila:

- sheet atau header berubah;
- periode tidak valid atau tanggal data berada di luar periode;
- NIP bukan 18 digit atau belum ada pada master pengguna;
- pasangan NIP+tanggal muncul lebih dari sekali;
- nama/tanggal/jam penting tidak dapat dibaca;
- file bukan XLSX.

Perbedaan nama atau unit dengan master pengguna ditampilkan sebagai peringatan. Baris kosong diabaikan dan dihitung pada preview.

## Pemulihan workbook contoh

Workbook `Upload dokumen.xlsx` yang diberikan berisi sheet utama lengkap tetapi central directory ZIP terpotong. Parser terlebih dahulu mencoba pembacaan XLSX normal. Bila gagal, parser memulihkan local entries yang lengkap, hanya menerima metode ZIP yang aman, kemudian menjalankan validasi yang sama. Preview akan menampilkan `recovered_partial_container`; admin tetap harus memeriksa jumlah baris dan pegawai sebelum publikasi.
