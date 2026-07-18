-- Generated from the supplied Rekapitulasi workbook.
-- Contains personal data (name and NIP): keep the repository private.
-- 1123 employees; 11 divisions; 46 sections.

begin;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
values ('KPU-TPK', 'KPU Bea dan Cukai Tipe A Tanjung Priok', '-', 'office', null, 0, true)
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-01', 'Bagian Umum', 'Bagian Umum - -', 'division', id, 10, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-02', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - -', 'division', id, 20, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-03', 'Bidang Keberatan', 'Bidang Keberatan - -', 'division', id, 30, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-04', 'Bidang Kepatuhan Internal', 'Bidang Kepatuhan Internal - -', 'division', id, 40, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-05', 'Bidang Pelayanan Fasilitas Pabean dan Cukai', 'Bidang Pelayanan Fasilitas Pabean dan Cukai - -', 'division', id, 50, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-06', 'Bidang Pelayanan Pabean dan Cukai I', 'Bidang Pelayanan Pabean dan Cukai I - -', 'division', id, 60, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-07', 'Bidang Pelayanan Pabean dan Cukai II', 'Bidang Pelayanan Pabean dan Cukai II - -', 'division', id, 70, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-08', 'Bidang Pelayanan Pabean dan Cukai III', 'Bidang Pelayanan Pabean dan Cukai III - -', 'division', id, 80, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-09', 'Bidang Pelayanan Pabean dan Cukai IV', 'Bidang Pelayanan Pabean dan Cukai IV - -', 'division', id, 90, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-10', 'Bidang Penindakan dan Penyidikan', 'Bidang Penindakan dan Penyidikan - -', 'division', id, 100, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'DIV-11', 'Bidang Perbendaharaan', 'Bidang Perbendaharaan - -', 'division', id, 110, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-01', 'Subbagian Dukungan Teknis', 'Bagian Umum - Subbagian Dukungan Teknis', 'section', id, 10, true
from public.units where code = 'DIV-01'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-02', 'Subbagian Keuangan', 'Bagian Umum - Subbagian Keuangan', 'section', id, 20, true
from public.units where code = 'DIV-01'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-03', 'Subbagian Sumber Daya Manusia', 'Bagian Umum - Subbagian Sumber Daya Manusia', 'section', id, 30, true
from public.units where code = 'DIV-01'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-04', 'Subbagian Tata Usaha dan Rumah Tangga', 'Bagian Umum - Subbagian Tata Usaha dan Rumah Tangga', 'section', id, 40, true
from public.units where code = 'DIV-01'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-05', 'Seksi Bimbingan Kepatuhan I', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - Seksi Bimbingan Kepatuhan I', 'section', id, 50, true
from public.units where code = 'DIV-02'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-06', 'Seksi Bimbingan Kepatuhan II', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - Seksi Bimbingan Kepatuhan II', 'section', id, 60, true
from public.units where code = 'DIV-02'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-07', 'Seksi Bimbingan Kepatuhan III', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - Seksi Bimbingan Kepatuhan III', 'section', id, 70, true
from public.units where code = 'DIV-02'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-08', 'Seksi Bimbingan Kepatuhan IV', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - Seksi Bimbingan Kepatuhan IV', 'section', id, 80, true
from public.units where code = 'DIV-02'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-09', 'Seksi Layanan Informasi', 'Bidang Bimbingan Kepatuhan dan Layanan Informasi - Seksi Layanan Informasi', 'section', id, 90, true
from public.units where code = 'DIV-02'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-10', 'Seksi Bantuan Hukum', 'Bidang Keberatan - Seksi Bantuan Hukum', 'section', id, 100, true
from public.units where code = 'DIV-03'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-11', 'Seksi Keberatan dan Banding I', 'Bidang Keberatan - Seksi Keberatan dan Banding I', 'section', id, 110, true
from public.units where code = 'DIV-03'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-12', 'Seksi Keberatan dan Banding II', 'Bidang Keberatan - Seksi Keberatan dan Banding II', 'section', id, 120, true
from public.units where code = 'DIV-03'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-13', 'Seksi Keberatan dan Banding III', 'Bidang Keberatan - Seksi Keberatan dan Banding III', 'section', id, 130, true
from public.units where code = 'DIV-03'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-14', 'Seksi Kepatuhan Pelaksanaan Tugas Administrasi', 'Bidang Kepatuhan Internal - Seksi Kepatuhan Pelaksanaan Tugas Administrasi', 'section', id, 140, true
from public.units where code = 'DIV-04'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-15', 'Seksi Kepatuhan Pelaksanaan Tugas Pelayanan', 'Bidang Kepatuhan Internal - Seksi Kepatuhan Pelaksanaan Tugas Pelayanan', 'section', id, 150, true
from public.units where code = 'DIV-04'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-16', 'Seksi Kepatuhan Pelaksanaan Tugas Pengawasan', 'Bidang Kepatuhan Internal - Seksi Kepatuhan Pelaksanaan Tugas Pengawasan', 'section', id, 160, true
from public.units where code = 'DIV-04'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-17', 'Seksi Perijinan dan Fasilitas Pabean dan Cukai I', 'Bidang Pelayanan Fasilitas Pabean dan Cukai - Seksi Perijinan dan Fasilitas Pabean dan Cukai I', 'section', id, 170, true
from public.units where code = 'DIV-05'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-18', 'Seksi Perijinan dan Fasilitas Pabean dan Cukai II', 'Bidang Pelayanan Fasilitas Pabean dan Cukai - Seksi Perijinan dan Fasilitas Pabean dan Cukai II', 'section', id, 180, true
from public.units where code = 'DIV-05'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-19', 'Seksi Perijinan dan Fasilitas Pabean dan Cukai III', 'Bidang Pelayanan Fasilitas Pabean dan Cukai - Seksi Perijinan dan Fasilitas Pabean dan Cukai III', 'section', id, 190, true
from public.units where code = 'DIV-05'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-20', 'Seksi Administrasi Manifes', 'Bidang Pelayanan Pabean dan Cukai I - Seksi Administrasi Manifes', 'section', id, 200, true
from public.units where code = 'DIV-06'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-21', 'Seksi Pabean dan Cukai I', 'Bidang Pelayanan Pabean dan Cukai I - Seksi Pabean dan Cukai I', 'section', id, 210, true
from public.units where code = 'DIV-06'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-22', 'Seksi Pabean dan Cukai II', 'Bidang Pelayanan Pabean dan Cukai I - Seksi Pabean dan Cukai II', 'section', id, 220, true
from public.units where code = 'DIV-06'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-23', 'Seksi Tempat Penimbunan', 'Bidang Pelayanan Pabean dan Cukai I - Seksi Tempat Penimbunan', 'section', id, 230, true
from public.units where code = 'DIV-06'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-24', 'Seksi Administrasi Manifes', 'Bidang Pelayanan Pabean dan Cukai II - Seksi Administrasi Manifes', 'section', id, 240, true
from public.units where code = 'DIV-07'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-25', 'Seksi Pabean dan Cukai I', 'Bidang Pelayanan Pabean dan Cukai II - Seksi Pabean dan Cukai I', 'section', id, 250, true
from public.units where code = 'DIV-07'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-26', 'Seksi Pabean dan Cukai II', 'Bidang Pelayanan Pabean dan Cukai II - Seksi Pabean dan Cukai II', 'section', id, 260, true
from public.units where code = 'DIV-07'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-27', 'Seksi Tempat Penimbunan', 'Bidang Pelayanan Pabean dan Cukai II - Seksi Tempat Penimbunan', 'section', id, 270, true
from public.units where code = 'DIV-07'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-28', 'Seksi Administrasi Manifes', 'Bidang Pelayanan Pabean dan Cukai III - Seksi Administrasi Manifes', 'section', id, 280, true
from public.units where code = 'DIV-08'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-29', 'Seksi Pabean dan Cukai I', 'Bidang Pelayanan Pabean dan Cukai III - Seksi Pabean dan Cukai I', 'section', id, 290, true
from public.units where code = 'DIV-08'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-30', 'Seksi Pabean dan Cukai II', 'Bidang Pelayanan Pabean dan Cukai III - Seksi Pabean dan Cukai II', 'section', id, 300, true
from public.units where code = 'DIV-08'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-31', 'Seksi Tempat Penimbunan', 'Bidang Pelayanan Pabean dan Cukai III - Seksi Tempat Penimbunan', 'section', id, 310, true
from public.units where code = 'DIV-08'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-32', 'Seksi Administrasi Manifes', 'Bidang Pelayanan Pabean dan Cukai IV - Seksi Administrasi Manifes', 'section', id, 320, true
from public.units where code = 'DIV-09'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-33', 'Seksi Pabean dan Cukai I', 'Bidang Pelayanan Pabean dan Cukai IV - Seksi Pabean dan Cukai I', 'section', id, 330, true
from public.units where code = 'DIV-09'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-34', 'Seksi Pabean dan Cukai II', 'Bidang Pelayanan Pabean dan Cukai IV - Seksi Pabean dan Cukai II', 'section', id, 340, true
from public.units where code = 'DIV-09'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-35', 'Seksi Tempat Penimbunan', 'Bidang Pelayanan Pabean dan Cukai IV - Seksi Tempat Penimbunan', 'section', id, 350, true
from public.units where code = 'DIV-09'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-36', 'Seksi Intelijen I', 'Bidang Penindakan dan Penyidikan - Seksi Intelijen I', 'section', id, 360, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-37', 'Seksi Intelijen II', 'Bidang Penindakan dan Penyidikan - Seksi Intelijen II', 'section', id, 370, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-38', 'Seksi Penindakan I', 'Bidang Penindakan dan Penyidikan - Seksi Penindakan I', 'section', id, 380, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-39', 'Seksi Penindakan II', 'Bidang Penindakan dan Penyidikan - Seksi Penindakan II', 'section', id, 390, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-40', 'Seksi Penindakan III', 'Bidang Penindakan dan Penyidikan - Seksi Penindakan III', 'section', id, 400, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-41', 'Seksi Penyidikan I', 'Bidang Penindakan dan Penyidikan - Seksi Penyidikan I', 'section', id, 410, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-42', 'Seksi Penyidikan II', 'Bidang Penindakan dan Penyidikan - Seksi Penyidikan II', 'section', id, 420, true
from public.units where code = 'DIV-10'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-43', 'Seksi Penagihan I', 'Bidang Perbendaharaan - Seksi Penagihan I', 'section', id, 430, true
from public.units where code = 'DIV-11'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-44', 'Seksi Penagihan II', 'Bidang Perbendaharaan - Seksi Penagihan II', 'section', id, 440, true
from public.units where code = 'DIV-11'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-45', 'Seksi Penerimaan dan Pengembalian I', 'Bidang Perbendaharaan - Seksi Penerimaan dan Pengembalian I', 'section', id, 450, true
from public.units where code = 'DIV-11'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'SEC-46', 'Seksi Penerimaan dan Pengembalian II', 'Bidang Perbendaharaan - Seksi Penerimaan dan Pengembalian II', 'section', id, 460, true
from public.units where code = 'DIV-11'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

insert into public.units (code, name, source_name, unit_type, parent_id, sort_order, is_active)
select 'FUNGSIONAL', 'Fungsional', 'Fungsional', 'functional', id, 990, true
from public.units where code = 'KPU-TPK'
on conflict (code) do update set
  name = excluded.name, source_name = excluded.source_name, unit_type = excluded.unit_type,
  parent_id = excluded.parent_id, sort_order = excluded.sort_order, is_active = true;

-- Demote seeded heads first so the partial unique indexes remain valid on re-run.
update public.users set position_role = 'staff'
where nip in (
  '197111201992011002', '196810141996031001', '197505141996031002', '198704052010122005', '199003162012101002', '199209062012101001', '199301102013101002', '199305162013101004',
  '199501172015021005', '199502092015021001', '199505112015121002', '199506102015021001', '199507162015021001', '199507232015021001', '199507302015021002', '199601032015021001',
  '199607072016121001', '199711212018012001', '199802202018121002', '199807202021011001', '199808292019121002', '199810102019121001', '199810312019121001', '199908312021011001',
  '200006112019121001', '200006162019122003', '200008072019121002', '197404091995031001', '198711122014021005', '199202282014111005', '199507162016122001', '199702272018011003',
  '199706282018011004', '199712012018121002', '199812302018122002', '199911262018122001', '199911302019122002', '200102052019122001', '197304221992122001', '198301242004121003',
  '198510052004121004', '198711172007011002', '198807082009121002', '199006052010011001', '199008252010011002', '199009202010121005', '199103192010011002', '199109032013101002',
  '199207182013102001', '199210092014112003', '199302042012101001', '199310122013101002', '199403022016121001', '199412032015021004', '199502112015021001', '199506052015021005',
  '199506262016121003', '199507212015021005', '199606212015121001', '199607252015021001', '199608102015122002', '199608222015122002', '199705282018011004', '199707252018011002',
  '199804202018011001', '199805172018122001', '199806152018012002', '199905102019122002', '199905232021011001', '199906222019122002', '199907172018121001', '200001082019122001',
  '200001152018121002', '200003182019122004', '200004202019122001', '200005022019122001', '200005022021011001', '200008292019121001', '197608221997031001', '197708101997031002',
  '197804092000011001', '198501262003122004', '199310232013101001', '199503202016122001', '199504202015021001', '199505082015021003', '199512052015021002', '199701132016121001',
  '199807152018121001', '199809142018122001', '199906082018122002', '199907102018122002', '199912222018121001', '200003312019122001', '200004182018121001', '200107112019121001',
  '197207091998031001', '197505031995022001', '198109212009011013', '198301232003122001', '198812042010121003', '198901222010011003', '199006052012102001', '199403042013101001',
  '199604252015021001', '199605182015021001', '199607252015021002', '199803292018121001', '199805122018011001', '199901242018011002', '199905242018121003', '200005092019122001',
  '197209021993011003', '198904042008121001', '199004262010121003', '199103262010121002', '199206132014112001', '199209262012101002', '199309172018011005', '199506142015122001',
  '199709032018011002', '199801022018011003', '199901092018122003', '200001072019122001', '200004252019122001', '200005032019121001', '197411051994021001', '198007272005011002',
  '198306012003122003', '199404272016121002', '199410292018012005', '199505142015021002', '199507082015021004', '199510302016121001', '199602242018011002', '199703012015122001',
  '199704012018121001', '199901052021011004', '199902152021011001', '199904042018122002', '199906092019121001', '199907182021011001', '199909102018122001', '197005161996032001',
  '199101242013102001', '199111042014111002', '199204012012101004', '199206122012101002', '199305032018011002', '199505152016122001', '199512102015021001', '199512212015021002',
  '199603032015122001', '199603052018012004', '199604072016121001', '199610112018012002', '199708132015122002', '199908162019121002', '199908182018122001', '199908312021011002',
  '199910052021011001', '200002232018122001', '196910221990121001', '199303172013101002', '199605292018011011', '199607102015121001', '199608202015121002', '199801272018011002',
  '199801292018011001', '199808182018012001', '199809102018011003', '199903102018122001', '199907012018121002', '199908252018121001', '200005182019122002', '200007212019121003',
  '196809061988121001', '197708191998031001', '198705262010122006', '199102082010121002', '199305032013101002', '199407262018011003', '199502012015021001', '199505202015021003',
  '199612182015122002', '197503301996031001', '198711262014022002', '199310142013101002', '199401252013101002', '199501012015021003', '199503122015021002', '199510072015021002',
  '199703152018012001', '197508151996031001', '199002032014021004', '199201132014111002', '199202232013101002', '199204282012101003', '199401272013101001', '199903062018122001',
  '200003202019122001', '197604201996031001', '198505102003122003', '199203142014111002', '199308212013101001', '199505142015021001', '199806082018012001', '199912202021011001',
  '197106161992011004', '197204261993021001', '199006062009121001', '199210202018011001', '199507082015122002', '199509032015021002', '199602202018121001', '199604232015121001',
  '199611262015122001', '199806172019121001', '199807292019121001', '197508071996021001', '199306082013101001', '199307292013101003', '199505142015022001', '199507192015121001',
  '199908252022011001', '197007281996032001', '199301122013101003', '199401052015021004', '199403262015021002', '199405042015021002', '199409092016121001', '199501092016121001',
  '199501182018012004', '199506112015021001', '199511282016121001', '199602042015021001', '199604232015021003', '199711272019121001', '199805052021011001', '197001101996031001',
  '197611231996021002', '198205022001122002', '198912052010011001', '199410072015021001', '199412042015021004', '199611232016122001', '199802142019121003', '199812062018121003',
  '200001022019122001', '197911252000121001', '199011102015022003', '199310142015021001', '199311092015021001', '199501072015021006', '199706282018121001', '199708192018121002',
  '199905042018122001', '199907012018121003', '200102232019121001', '200105182019122001', '197503131996021001', '199207272012101003', '199309032013101002', '199311242013101002',
  '199505112018012004', '199604152015021001', '199605282015122003', '199901092018121002', '199902212019122001', '199903252021011001', '196809241989121001', '197706171999031001',
  '199505142018011003', '199602132016121001', '199612212015122001', '199703222019121001', '199706222016121001', '199710072019121002', '199807292021011002', '199808222019121001',
  '199809062018011001', '199905252021011001', '199910192019122001', '200009172019122001', '198008162002121001', '198911232013101001', '199204052014111002', '199210152013101009',
  '199505172015021003', '199512132015121002', '199606192015121004', '199609222018011002', '199709032016122001', '199805132018121002', '199806152018012003', '199807242018011001',
  '199901122021011002', '199905032019121002', '199906232018122002', '199910022019122001', '199912072018122001', '200007142019122002', '197404211994021001', '199003302010011005',
  '199208152010121001', '199501172015021001', '199503182015021001', '199505292015021001', '199510202015121003', '199511282015121003', '199512052015021001', '199603092016122001',
  '199803232018121001', '199807272018011003', '199907082018121001', '199907112019122001', '199910262018121001', '199912122019122002', '200003102019121002', '200004172019122001',
  '200005082019122001', '200008162018122001', '200008172018121001', '200008262019122001', '197308191996021001', '198105262005011001', '199003232010121009', '199010052012101003',
  '199602172018011002', '199609132015121003', '199701252018011002', '199806112019121002', '199809072021011001', '200002202018122001', '196808191988121001', '197708252000011002',
  '199503022015021001', '199605232015122002', '199607152015121001', '199702222015121002', '199707042018011002', '199808282021011001', '199812082018011001', '199903012019121001',
  '199907132018122002', '199910252018121002', '200004082019122001', '197406141994031002', '199406282015021002', '199408042015021001', '199409052015021002', '199410152015021003',
  '199412012015021002', '199501172015021002', '199501242015021001', '199502042015021008', '199502072015021002', '199502142015121001', '199503132015021003', '199503142015021003',
  '199504042015021004', '199505102015021002', '199506042015021002', '199507012015021007', '199507142015021006', '199507212015021008', '199508102015021001', '199508112015021001',
  '199509112015121002', '199510082015021002', '199601212018011004', '199602132015021002', '199603242015121003', '199604062015121002', '199604192015121002', '199605142015121003',
  '199609042015121003', '199610062016121001', '199611022015121002', '199611142015121005', '199701132016121002', '199702102018011004', '199702142018011001', '199703082016121001',
  '199703122018121001', '199704022015121001', '199705082016121001', '199706032016121002', '199706122016121001', '199708312016121001', '199709052018011001', '199710312018011002',
  '199711102016121001', '199711262016121001', '199712022016121001', '199802032018011002', '199802232018011001', '199803082018011002', '199804282018011003', '199805042018011001',
  '199805212019121001', '199806112018121002', '199807062018121001', '199807122018121001', '199808252018011001', '199809192018121001', '199809242018121001', '199812132018011003',
  '199812272018011001', '199901062018011001', '199901242018121001', '199902122018121001', '199903032018121001', '199903032018121002', '199903052018121001', '199903122019121001',
  '199904292018122002', '199905072018011001', '199905242018121001', '199906162019121001', '199906192019121001', '199907152018121001', '199907222018121002', '199907302018121005',
  '199909012018121001', '199910282018121001', '199910292018121003', '199911152018121002', '199911172019121002', '199912192018121001', '199912282018121003', '200003222019121003',
  '200003272018121001', '200004142018121001', '200007252019121001', '200008092019121001', '200009012019121001', '200009302019121001', '200012012019122001', '197702151999031001',
  '199708072018011002', '199809212018011001', '199812012018122001', '199906052018011001', '199909052018122002', '200001152019122003', '200003062019121001', '200008092019122001',
  '200009072019122002', '197304121993022001', '199406102015021001', '199502232015021003', '199506272016121001', '199603192016121001', '199603222015121003', '199605182015121001',
  '199609072015121002', '199808042018011001', '199808082021011001', '199809032018012003', '199812212019122001', '199901102018121001', '199904132021011001', '199908272019121001',
  '200004122019122001', '196912251996031001', '197512081996031004', '199412142015021001', '199505092015021003', '199511102016121002', '199612092016121001', '199702012018121001',
  '199706042018012001', '199710162018012001', '199711082018011001', '199807142021011003', '200001082018122001', '200001142019121002', '200002012021011001', '197202111998031002',
  '198803212007011001', '199007182010011001', '199106222012101001', '199201182010121004', '199201252013101003', '199207162013101001', '199209282012101001', '199310032012101001',
  '199311292013101001', '199401042015021001', '199409172015021002', '199501092015021002', '199501132015021003', '199505132018011002', '199505252015021002', '199506142015021001',
  '199506252016121002', '199506292015021001', '199508072015021002', '199508132016121001', '199510152015021002', '199511012015121001', '199609182015121002', '199611172015121002',
  '199703272016121001', '199805072018121002', '199807012021011003', '199809112018011002', '199809142018121001', '199812182018011002', '199812252018121001', '199904072021011001',
  '199908032018121001', '199909052021011002', '199909182018121001', '199911172019121001', '200001032021011001', '197406091994021001', '198204052005011001', '199004222012102001',
  '199202162013101001', '199202192012101003', '199411142015021003', '199502132015021002', '199503092016121002', '199503092016121003', '199503192015021003', '199510042015021007',
  '199512032016121001', '199602252018011002', '199603092015021001', '199605062015121004', '199801282018012001', '199804082019121001', '199811102021011001', '199905292018121002',
  '199908262019122001', '200001142019122001', '200104212019122001', '197107191992011001', '199204132013101001', '199308102013101004', '199409012015021004', '199412102015021003',
  '199501222015021002', '199510152018011002', '199512012015021001', '199611042015121004', '199704142018121001', '199711202019121001', '199811012021011001', '199812062018012001',
  '200004122019121001', '197010011990121001', '197409301994021001', '198508012003122003', '199407042015021004', '199512132016121001', '199602102015122001', '199605292018011006',
  '199710252019121001', '199911082018121001', '200005232019122003', '200010242019122001', '197710252002121001', '199408072015021001', '199411112015021003', '199507042015021003',
  '199507092015021004', '199507192015021004', '199510132015021004', '199511122015021001', '199803102018011001', '199804292018011001', '199901192018121002', '199905162018121003',
  '199908252018122004', '199909182019122003', '200007022019122004', '200105022019121001', '197603241996031001', '199205272012101001', '199207042012101002', '199407082015022001',
  '199412232015021002', '199505212015021001', '199505282015021005', '199507112016122002', '199604272015121005', '199605102015121001', '199608172015021001', '199807162018122001',
  '199811272018011002', '199904282018122001', '199906022018122001', '199911172019122001', '200002052019122001', '200008162019122002', '197502031999031003', '199208292015021001',
  '199407212015021004', '199601062015021001', '199611072016121001', '199710242018121001', '199711272018011003', '199812112021011001', '199904282018121002', '199906062018121001',
  '199908012019122001', '200001082018122002', '198003012001121001', '198203202003121001', '199105052012101001', '199208222012101001', '199209152013101002', '199305222013101001',
  '199403112015021002', '199403312015021001', '199404022015021002', '199404272015021001', '199409302018011001', '199501302015021001', '199505152015021002', '199505312015021001',
  '199508042015021004', '199508092015021004', '199509022015021001', '199511102018011003', '199512032015021001', '199602032018011004', '199603032018011002', '199603112018011002',
  '199607262015121001', '199608132018121001', '199609042018121001', '199611172016121001', '199701082015121001', '199703132018121003', '199703242018011001', '199705272019121001',
  '199708032018121001', '199708062018121002', '199711082018011003', '199711262018011002', '199711262018121001', '199803232019121001', '199804222021011001', '199806052018121002',
  '199809072019121001', '199809192019121002', '199904092018121003', '199904172021011001', '199905102021011001', '200005122019121001', '197806142000121001', '198210152003121003',
  '198908062010011002', '198911242010011001', '199009292014111001', '199011082010011001', '199101202010011001', '199203302013101002', '199303052013101002', '199306032013101004',
  '199309032013101001', '199311232015021001', '199311302015021001', '199501042015021001', '199501082018011002', '199501112015021002', '199504192015121001', '199504292015021002',
  '199506252015021003', '199604032015121002', '199607042016121001', '199607172015121002', '199607282015121001', '199611042018121001', '199612262015121001', '199801282018011002',
  '199804142018012001', '199806122018011002', '199807222018011003', '199906272018121001', '199907292018121003', '199909172021011002', '198311122009011007', '198612212008121001',
  '199109172012101001', '199204292012101001', '199205142013101001', '199308182013101002', '199402152015021002', '199404032015021001', '199410232016121002', '199502202015021001',
  '199502242015021004', '199504082015021001', '199508072015021001', '199510042016121001', '199512152015121002', '199601012018011003', '199606032015021001', '199709172018011003',
  '199712212018121001', '199807172019121001', '199902052021011002', '199902102019122001', '199911022018121003', '197409041995031001', '198704162007101001', '199106032014111001',
  '199107192010121002', '199211202012101001', '199311102015021002', '199311282013101001', '199402222015021002', '199501312018011003', '199502182015021004', '199505182015021004',
  '199506142015021002', '199510302015121003', '199511052015121002', '199511252018011003', '199512302015121002', '199601192015021001', '199604012015021002', '199604292018011003',
  '199605172018011002', '199609042018011001', '199703242016121002', '199710112018121001', '199710302016121001', '199801022016121001', '199804042021011001', '199805072018011003',
  '199805132018012001', '199808262018011002', '199809172019121001', '199901092018121001', '199902232021011002', '199903152018011001', '199907062018121002', '197702152000011001',
  '198506072003122002', '198911172012101002', '199112162012101001', '199207112014111001', '199304102013101002', '199308312013101001', '199312052015021001', '199312212015021001',
  '199502272016121001', '199507222016121001', '199509132015121002', '199512212015121001', '199601202016121002', '199801152018011001', '199804022018121002', '199809162019121001',
  '199902252021011001', '199910162019121001', '199910292018121002', '200002072019121001', '200006292019122001', '197711112000011003', '199101282010011001', '199502232015021002',
  '199503052016121002', '199504112015021004', '199511252018011006', '199905052018122004', '199905192021011002', '199905292021011001', '197605261996021001', '199011122013101002',
  '199502082015021003', '199504012018012001', '199602142018011005', '199610162018011002', '199811092018121001', '199812272021011001', '196901151996031001', '197004031996032002',
  '199403202013101002', '199503182015021002', '199506012015021001', '199508292015021002', '199907252018122001', '199912162019121001', '200005242019122002', '197308061994022003',
  '199202152013101004', '199502212018011004', '199510122016121001', '199802122018011002', '199809032018122002', '199901272021011002', '200003162018122001', '197201211996032001',
  '199505092015121002', '199511042016122001', '199807062018011001', '199807082018122001', '200007242019121001', '197305221992121001', '198512302006041002', '198802262009121006',
  '199207142013102001', '199302262013101004', '199404222015021003', '199503112015021003', '199505302015021003', '199512152015021002', '199602122015122003', '199705092018121002',
  '199712072018012002', '199712092018012001', '199712102018011002', '199712282018122003', '199802112016122001', '199808132018121001', '199810092021011001', '199901212018121001',
  '199908312018122001', '199911102018122001', '200002102022011001', '197009301990121002', '197201161992121001', '197207281999031001', '197304171992121001', '197312301994021001',
  '197402051994021001', '197406131994021002', '197408251997031002', '197501121997031002', '197502141994022001', '197503031997031001', '197503081995031002', '197504102005011001',
  '197505071997031003', '197505302005011001', '197507091997031001', '197507092005011001', '197508161996021001', '197510301999031001', '197511251999031001', '197601111997031001',
  '197601161996021003', '197602271997031001', '197603301996021003', '197604281999031001', '197605091998031002', '197606292003121001', '197607281998031001', '197608051997031001',
  '197608201999031002', '197608281996021002', '197609211996031003', '197610011997031001', '197610111997031002', '197701101999031001', '197701121999031001', '197702031997031001',
  '197702051999031002', '197702052000011001', '197702091997031001', '197704011997031001', '197705141997031002', '197705191997031001', '197706101998031001', '197706151999031001',
  '197707281998031001', '197709081997031001', '197709091999031001', '197709172000011002', '197710101997031001', '197712101998031001', '197712191997031001', '197801041998031002',
  '197801091998031002', '197802222000012001', '197802262003122001', '197803182000011001', '197805062000011001', '197806072003121001', '197806191999031001', '197807122000011002',
  '197807162003121001', '197808042000121002', '197810101998031002', '197812121998031001', '197903042003121002', '197904061999031003', '197905042001121001', '197906242001121001',
  '197907042000121001', '197908122001121002', '197908181998031001', '197908252001121001', '197911252000121003', '197912182001121001', '198002142003121001', '198002192001121001',
  '198002272003121001', '198004282001121002', '198005122000121002', '198006122002122001', '198006202001121002', '198007222005011001', '198008272001121004', '198010092000121002',
  '198010102000121001', '198010102001121002', '198010302002122002', '198012112002121002', '198101012003121002', '198101102003122002', '198101182003121002', '198101312000121001',
  '198108112002121001', '198110182002121002', '198110232003121001', '198110272002121001', '198111252003121001', '198111292003121001', '198202072003121001', '198202102001121001',
  '198202262003121001', '198204022003122001', '198204102001121001', '198204162003122001', '198205262003121001', '198206262003121002', '198207042004121001', '198207112001121002',
  '198207212001121001', '198210212004121002', '198212082003121001', '198304222002121002', '198307282002121002', '198308062010121003', '198308152004121004', '198308312009011008',
  '198310182004121001', '198311262006021001', '198401132010121005', '198402282007011001', '198403262009011006', '198405202004121002', '198407182004121001', '198408122009011005',
  '198412162006021002', '198501112006021003', '198504172003121002', '198504182006021002', '198505212010121007', '198505232007011002', '198507272006021001', '198508162009011005',
  '198509062007011001', '198509212010121004', '198510252006021004', '198510302010121002', '198511092007011001', '198512132007011002', '198602092007101001', '198602202006021003',
  '198603232007101002', '198603282007101001', '198604102007101001', '198604112006021004', '198606082014021006', '198608052006021002', '198610082006021004', '198612042010121005',
  '198703122007101002', '198705112007101001', '198709182007101001', '198710082015021001', '198712122008121002', '198803242015021003', '198805262007101003', '198810252009121004',
  '198901172010011004', '198902162010011002', '198910052014021006', '198910272012101001', '198912232010011002', '199001192010011001', '199004242014022002', '199005142010011001',
  '199005252010011002', '199008242012101002', '199009232014022002', '199010032012102002', '199010272010011001', '199012152010121001', '199101272012101002', '199103312010121003',
  '199105182012101002', '199105272009121001', '199106012012101001', '199107232010121002', '199108242012101001', '199109172012101002', '199109182014021001', '199109212010121002',
  '199111062012101001', '199111102013102001', '199111142013101002', '199111272013101001', '199111282012101001', '199112112012101001', '199202072012101001', '199202072013101001',
  '199202242013101003', '199203202012101001', '199204212013101002', '199204222012101001', '199205022012101001', '199206072012101003', '199208082013101001', '199208132013101001',
  '199209062012101002', '199209292013101002', '199211022013101001', '199211102013101002', '199211152012101001', '199212152018011006', '199301102013101004', '199304122013101002',
  '199304262013101001', '199307142013101001', '199308072013101001', '199309022013101002', '199309212018012004', '199309272015021002', '199310042015021001', '199310282013101001',
  '199311062013101003', '199401052013101001', '199401252015021001', '199401262015021001', '199402102015021001', '199404242013101001', '199404252015021004', '199405232015021003',
  '199405242015021003', '199408252015021002', '199410022015021002', '199412162015021002', '199412312015021005', '199501052015021002', '199501232015021002', '199502042016121002',
  '199502172015021007', '199502252015021001', '199503212016121001', '199503232015021003', '199503252015021002', '199504072015021002', '199504142015021001', '199504242015021007',
  '199505102015021005', '199506032015021001', '199507022015021001', '199507182015121001', '199508072015021005', '199508112015021003', '199509052015021001', '199509182015021003',
  '199509222018011003', '199510102015121002', '199510212015021001', '199510212015121001', '199510282015021004', '199511042018011002', '199511122015021002', '199511132015021002',
  '199512132015021001', '199512202015021001', '199601092015121002', '199601172015021001', '199601172015021002', '199602182015121003', '199602182015121004', '199602232015021001',
  '199602282015021001', '199603162015121002', '199603262015021002', '199604082015121002', '199604242015021001', '199604262018011004', '199605312015121002', '199606132015121002',
  '199606202018121002', '199606232018011004', '199607292015021001', '199609082015121001', '199610192015121001', '199611152018121001', '199702102016121002', '199703142015121001',
  '199703142016121001', '199707042015121001', '199708122018121001', '199709172019121001', '199709182018011002', '199801062018011001', '199803142018011001', '199807192018011001',
  '199809022018011001', '199809182018011001', '199809252019121001', '199810282018011001', '199901012018121002', '199901152018011002', '199901252018011001', '199905072018121005',
  '200001222018121001', '200006082019121001', '200006202019121002'
);

insert into public.users (nip, name, unit_id, position_role, password_hash, must_change_password, is_active)
values
  ('197111201992011002', 'Adhang Noegroho Adhi', (select id from public.units where code = 'KPU-TPK'), 'office_head', null, true, true),
  ('196810141996031001', 'Heru Prayitno', (select id from public.units where code = 'DIV-01'), 'division_head', null, true, true),
  ('197505141996031002', 'Misnawi', (select id from public.units where code = 'SEC-01'), 'section_head', null, true, true),
  ('198704052010122005', 'Lismawarti Nurfitri', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199003162012101002', 'Rohmad', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199209062012101001', 'Fattah Yubib Mubarok', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199301102013101002', 'Muhammad Syaifuddin Ikhsan', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199305162013101004', 'I Made Agus Merthayasa Utama', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199501172015021005', 'M. Naufal Hanif', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199502092015021001', 'Karunia Aditama', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199505112015121002', 'Fungki Putra Syamsuddin', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199506102015021001', 'Wisnu Albar Dwiwibowo', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199507162015021001', 'Dwiyan Kurnianto Saputro', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199507232015021001', 'Rizki Akbar Wijaya', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199507302015021002', 'Dhimas Wicaksono', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199601032015021001', 'Sandiyo Sunarko', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199607072016121001', 'Reidhon Lanri Luciano Subagio', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199711212018012001', 'Dyah Ayu Murni Shaleha', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199802202018121002', 'Arta Atmawijaya Kartoloh', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199807202021011001', 'Muhammad Iqbal', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199808292019121002', 'Ahmad Rifqi Nur Rosyid', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199810102019121001', 'Agung Budi Saputro', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199810312019121001', 'Brilliant Arthur Sebastian', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('199908312021011001', 'Raihan Naufal Aziz', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('200006112019121001', 'Bagas Nur Rohman', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('200006162019122003', 'Leny Randi Barkah', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('200008072019121002', 'Refifareli', (select id from public.units where code = 'SEC-01'), 'staff', null, true, true),
  ('197404091995031001', 'Ariansyah', (select id from public.units where code = 'SEC-02'), 'section_head', null, true, true),
  ('198711122014021005', 'Iis Iswandy', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199202282014111005', 'Yohannes Heryanto', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199507162016122001', 'Eunike Sianturi', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199702272018011003', 'Fahry Amrizal Pawitra', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199706282018011004', 'Dandy Yunanza Anggara', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199712012018121002', 'Keen Deswandy Saragih', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199812302018122002', 'Cindy Mellinda', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199911262018122001', 'Zetha Flandira Martan', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('199911302019122002', 'Dina Novitasari', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('200102052019122001', 'Vera Fazlina Saragih', (select id from public.units where code = 'SEC-02'), 'staff', null, true, true),
  ('197304221992122001', 'Rini Setiyowati', (select id from public.units where code = 'SEC-03'), 'section_head', null, true, true),
  ('198301242004121003', 'Dimas Pratama', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('198510052004121004', 'Muhammad Royhan', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('198711172007011002', 'Mu''Ammar Ilyas', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('198807082009121002', 'Ganang Sutawijaya', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199006052010011001', 'Irfan Nur Ilman', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199008252010011002', 'Alfin Yudistira', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199009202010121005', 'Richan Cahya Pribadi', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199103192010011002', 'Fitra Aidin Nuansa Muassa', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199109032013101002', 'Yan Rizqi Kurniawan', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199207182013102001', 'Utami', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199210092014112003', 'Enni Sayekti', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199302042012101001', 'Muhammad Sahri Aziz', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199310122013101002', 'Lukman Nulhakim', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199403022016121001', 'Burhan Pirzada', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199412032015021004', 'Muhammad Nurfadilah', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199502112015021001', 'Dimas Wahyu Susanto', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199506052015021005', 'Sandy Putra Godlas Siahaan', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199506262016121003', 'Ronald Guntara', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199507212015021005', 'Arya Laksamana Nugroho', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199606212015121001', 'Muhamad Fadhil Kurniawan', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199607252015021001', 'Bhisma Haryo Samodro', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199608102015122002', 'Ummah Hamidah', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199608222015122002', 'Maulida Khomariyah', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199705282018011004', 'Hilmy Abyansyah Nugraha', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199707252018011002', 'Irjayat Reza Mardika', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199804202018011001', 'Yudha Prasetya', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199805172018122001', 'Hariyani Kurnia Dewi', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199806152018012002', 'Zunia Nafi''Atullina', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199905102019122002', 'Aldarida Anggita Kusumastuti', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199905232021011001', 'Chalvin Sitepu', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199906222019122002', 'Heny Kurniawati', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('199907172018121001', 'Ginting Praba Waluyo', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200001082019122001', 'Indira Dian Fadhilah', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200001152018121002', 'Muhammad Aris Azhimi', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200003182019122004', 'Imala Islam Madania', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200004202019122001', 'Ainun Fatihah Salsabila', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200005022019122001', 'Ading Meidika Rilanti', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200005022021011001', 'Irfan Fadilah', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('200008292019121001', 'Hawley Naufal Muhammad', (select id from public.units where code = 'SEC-03'), 'staff', null, true, true),
  ('197608221997031001', 'Agus Praminto', (select id from public.units where code = 'SEC-04'), 'section_head', null, true, true),
  ('197708101997031002', 'Romanus Rano Agus T.H.', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('197804092000011001', 'Parasian Silitonga', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('198501262003122004', 'Nuraini Fitriyanti', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199310232013101001', 'Widi Arsandi', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199503202016122001', 'Bibiana Tri Widiastuti', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199504202015021001', 'Albani Ahmad', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199505082015021003', 'Muhammad Syahrial', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199512052015021002', 'Rizal Abdillah', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199701132016121001', 'R. Lukman Ludiansyah', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199807152018121001', 'Mochammad Hisyam', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199809142018122001', 'Nanda Dwi Hidayati', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199906082018122002', 'Rema Jesika Br Sitepu', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199907102018122002', 'Galuh Ayshandra Karina Putri', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('199912222018121001', 'Fauzi Ichad Wiekaldie', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('200003312019122001', 'Nanda Marlina', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('200004182018121001', 'Muhammad Togawa Rasyid Ridla', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('200107112019121001', 'Barakah Syahadat', (select id from public.units where code = 'SEC-04'), 'staff', null, true, true),
  ('197207091998031001', 'Niko Budhi Darma', (select id from public.units where code = 'DIV-02'), 'division_head', null, true, true),
  ('197505031995022001', 'Heny Rusindarti', (select id from public.units where code = 'SEC-05'), 'section_head', null, true, true),
  ('198109212009011013', 'Dwi Hantono S.W', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('198301232003122001', 'Wiwin Rahayu', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('198812042010121003', 'Ogi Boi Sarjanto Sitohang', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('198901222010011003', 'Misbahul Amin', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199006052012102001', 'Hesty Yuniasih Pratama', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199403042013101001', 'Iqbal Aji Harjuna', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199604252015021001', 'Helga Candra Negara', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199605182015021001', 'Nurul Huda Zainal Mutaqim', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199607252015021002', 'Harry Mauladi', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199803292018121001', 'Andrian Izza Prayudhi', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199805122018011001', 'Hana Adi Nirwaana', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199901242018011002', 'Musthafa ''Azzam Tsabit', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('199905242018121003', 'Muhammad Yusuf Dwi Rizaldy', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('200005092019122001', 'Andarin Andar Prasasti', (select id from public.units where code = 'SEC-05'), 'staff', null, true, true),
  ('197209021993011003', 'Syamsu Priatmojo', (select id from public.units where code = 'SEC-06'), 'section_head', null, true, true),
  ('198904042008121001', 'Yayan Yuliandi Yunahar', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199004262010121003', 'Haposan Indra Wesly Pasaribu', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199103262010121002', 'Tegdi Subanda Manullang', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199206132014112001', 'Reisa Devi Maharani', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199209262012101002', 'A. Akbar Kurniawan Fr', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199309172018011005', 'Ahmad Yosep Setiaji', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199506142015122001', 'Laili Nugrahaeni', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199709032018011002', 'Ikhsan Akbar Triadi', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199801022018011003', 'Alvin Rianto', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('199901092018122003', 'Winda Aris Setyowati', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('200001072019122001', 'Amalia Damayanti', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('200004252019122001', 'Niche Srirugyatuz Zahra', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('200005032019121001', 'Alexander Ananda Risky Kurniawan', (select id from public.units where code = 'SEC-06'), 'staff', null, true, true),
  ('197411051994021001', 'Yusep Sasmita', (select id from public.units where code = 'SEC-07'), 'section_head', null, true, true),
  ('198007272005011002', 'Kabrie Aliakim', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('198306012003122003', 'Karti', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199404272016121002', 'Hendrian Yustio Nugroho', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199410292018012005', 'Rezkita Fajriati Gani', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199505142015021002', 'Muhammad Ridwan Yusuf', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199507082015021004', 'Faisal Nashuha Adlin', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199510302016121001', 'Rizki Abdillah H. Siregar', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199602242018011002', 'Ibnu Rizal Rabbani', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199703012015122001', 'Tivara Merliana Putri', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199704012018121001', 'Firman Agung Aji Setiawan', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199901052021011004', 'Adnantiya Asfahany', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199902152021011001', 'Rifqi Darmawan', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199904042018122002', 'Melvy Permatasari', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199906092019121001', 'Patar Mangasi Sihaloho', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199907182021011001', 'Uzlifat Dinu Salata', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('199909102018122001', 'Alifiandri Zaneta Putri', (select id from public.units where code = 'SEC-07'), 'staff', null, true, true),
  ('197005161996032001', 'Ambar Susilowati', (select id from public.units where code = 'SEC-08'), 'section_head', null, true, true),
  ('199101242013102001', 'Aris Chusnul Rachmawati', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199111042014111002', 'Didit Rahadi Setiadi', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199204012012101004', 'Rizki Rachmatullah Catur Putra', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199206122012101002', 'Erwin Siahaan', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199305032018011002', 'Bahrul Hanif', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199505152016122001', 'Andinia May Cahya', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199512102015021001', 'Mohamad Adnan Ghifari', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199512212015021002', 'Dhani Muhammad Syahputra Barus', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199603032015122001', 'Merlina Irvana Fitri', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199603052018012004', 'Chintyas Agrandis Anantaka', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199604072016121001', 'Dulmalik', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199610112018012002', 'Navyana Hena Oktavia', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199708132015122002', 'Roselyne Taruli Simanjuntak', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199908162019121002', 'Habil Alrasyid', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199908182018122001', 'Prischa Agustina', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199908312021011002', 'Rifan Fahriansyah', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('199910052021011001', 'Amin Prawiro Madhani', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('200002232018122001', 'Hannysya', (select id from public.units where code = 'SEC-08'), 'staff', null, true, true),
  ('196910221990121001', 'Rinto Setiawan', (select id from public.units where code = 'SEC-09'), 'section_head', null, true, true),
  ('199303172013101002', 'Soma Ainur Rahma', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199605292018011011', 'Yoga Haryo Pratikto', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199607102015121001', 'Pipen Dewantoro Wibowo', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199608202015121002', 'Patrick Nackok', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199801272018011002', 'Hasya Luthfi Akbarudin', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199801292018011001', 'Firmansyah Bagus Wibisono', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199808182018012001', 'Yoan Agustin', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199809102018011003', 'Raden Mas Fajrul Falah', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199903102018122001', 'Intan Bella Pratiwi', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199907012018121002', 'Muhammad Ardiansyah', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('199908252018121001', 'Adib Rahmat', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('200005182019122002', 'Amanda Ayu Nur Fajrina', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('200007212019121003', 'Redho Putra Erlandi', (select id from public.units where code = 'SEC-09'), 'staff', null, true, true),
  ('196809061988121001', 'Edy Susetyo', (select id from public.units where code = 'DIV-03'), 'division_head', null, true, true),
  ('197708191998031001', 'Carl Augustinus Hothinca Soutihon Tampubolon', (select id from public.units where code = 'SEC-10'), 'section_head', null, true, true),
  ('198705262010122006', 'Dhika Widyartika', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199102082010121002', 'Febri Heriansyah', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199305032013101002', 'Aris Pranata', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199407262018011003', 'Zulfadli Zulfikar Zen', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199502012015021001', 'Resnansyah Muhammad', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199505202015021003', 'Bayu Triyogo', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('199612182015122002', 'Tasya Maurizka Putri Pramusjanto', (select id from public.units where code = 'SEC-10'), 'staff', null, true, true),
  ('197503301996031001', 'Zubair Alimustaka', (select id from public.units where code = 'SEC-11'), 'section_head', null, true, true),
  ('198711262014022002', 'Ummi Isnaini', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199310142013101002', 'Adie Santoso', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199401252013101002', 'Sultoni', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199501012015021003', 'Movi Herliansyah', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199503122015021002', 'Mohammad Luthfi Pryantomo', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199510072015021002', 'Agni Nala Pamungkas', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('199703152018012001', 'Ruth Sorta Mutiara Nainggolan', (select id from public.units where code = 'SEC-11'), 'staff', null, true, true),
  ('197508151996031001', 'Jliteng Wibowo', (select id from public.units where code = 'SEC-12'), 'section_head', null, true, true),
  ('199002032014021004', 'Rendy Hardy Syaputra', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('199201132014111002', 'Yos Ricki Yanuar', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('199202232013101002', 'Andre', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('199204282012101003', 'Andhika Anggie Hutomo', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('199401272013101001', 'Fanny Avianuari', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('199903062018122001', 'Azella Dina Rosani', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('200003202019122001', 'Thessalonica Ega Devara', (select id from public.units where code = 'SEC-12'), 'staff', null, true, true),
  ('197604201996031001', 'Imam Supriadi', (select id from public.units where code = 'SEC-13'), 'section_head', null, true, true),
  ('198505102003122003', 'Triyani', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('199203142014111002', 'Achmad Barik Romadlon Nawawi', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('199308212013101001', 'Riyan Andriyanto', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('199505142015021001', 'Muhammad Bilal Ichsanul Putra', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('199806082018012001', 'Paramita Yuniasta', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('199912202021011001', 'Fikri Ramadhan', (select id from public.units where code = 'SEC-13'), 'staff', null, true, true),
  ('197106161992011004', 'Muhamad Irwan', (select id from public.units where code = 'DIV-04'), 'division_head', null, true, true),
  ('197204261993021001', 'Rahmanto', (select id from public.units where code = 'SEC-14'), 'section_head', null, true, true),
  ('199006062009121001', 'Dandyo Ferry Sandria Tanjung', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199210202018011001', 'Farouk Badri Al Baehaki', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199507082015122002', 'Isabella Adiawati', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199509032015021002', 'Syaifuddin Tansa Wicaksana', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199602202018121001', 'Beny Rich A. Tampubolon', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199604232015121001', 'Ridho Saputra', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199611262015122001', 'Indah Kurnia Taqwati', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199806172019121001', 'Imam Murtaqi', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('199807292019121001', 'Muhammad Bayu Yuliansyah', (select id from public.units where code = 'SEC-14'), 'staff', null, true, true),
  ('197508071996021001', 'Teguh Pribadi', (select id from public.units where code = 'SEC-15'), 'section_head', null, true, true),
  ('199306082013101001', 'Mohamad Ibnu Mulia Reshafahmi', (select id from public.units where code = 'SEC-15'), 'staff', null, true, true),
  ('199307292013101003', 'Edi Prabowo', (select id from public.units where code = 'SEC-15'), 'staff', null, true, true),
  ('199505142015022001', 'Meiliana Stephany Silalahi', (select id from public.units where code = 'SEC-15'), 'staff', null, true, true),
  ('199507192015121001', 'Maju P. Sitorus', (select id from public.units where code = 'SEC-15'), 'staff', null, true, true),
  ('199908252022011001', 'Mohammad Alfat Husnurozak', (select id from public.units where code = 'SEC-15'), 'staff', null, true, true),
  ('197007281996032001', 'Imelda Malik', (select id from public.units where code = 'SEC-16'), 'section_head', null, true, true),
  ('199301122013101003', 'Muhammad Faishal', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199401052015021004', 'Hilmansyah Damanik', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199403262015021002', 'Anggit Wicaksono Putro', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199405042015021002', 'Hafiz Mulyahadi', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199409092016121001', 'Ferial Hendi Nugroho', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199501092016121001', 'Heru Setyo Utomo', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199501182018012004', 'Winda Widyaningsih', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199506112015021001', 'Yostra Herdiawan', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199511282016121001', 'Wibowo Atmojo', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199602042015021001', 'Arif Suhendri', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199604232015021003', 'Irfan Fathoni Lukman', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199711272019121001', 'Hasan Ba''Isa', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('199805052021011001', 'Habib Mugni Al Sa''Ad', (select id from public.units where code = 'SEC-16'), 'staff', null, true, true),
  ('197001101996031001', 'Swoko Adi', (select id from public.units where code = 'DIV-05'), 'division_head', null, true, true),
  ('197611231996021002', 'Arif Rifani', (select id from public.units where code = 'SEC-17'), 'section_head', null, true, true),
  ('198205022001122002', 'Riska Maya Kurniasari', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('198912052010011001', 'Ardi Raharja Nur', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('199410072015021001', 'Handis Oktavian', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('199412042015021004', 'Younan Nur Zonanda', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('199611232016122001', 'Nurani Setiawati', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('199802142019121003', 'Febrian Dharma Putra Sakti', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('199812062018121003', 'Muhammad Ridwan', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('200001022019122001', 'Devi Mutia Zulfarizki', (select id from public.units where code = 'SEC-17'), 'staff', null, true, true),
  ('197911252000121001', 'Muchtian Purwoko', (select id from public.units where code = 'SEC-18'), 'section_head', null, true, true),
  ('199011102015022003', 'Dewi Susliana Fauzani', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199310142015021001', 'Ali Prakoso Wibowo', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199311092015021001', 'Sri Cahyo Prasojo', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199501072015021006', 'Wipra Prasetya Adil Syahputra', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199706282018121001', 'Riza Yusnizar Aliza', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199708192018121002', 'Muammar Rifki Anwar', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199905042018122001', 'Safira Amalia Fauziah', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('199907012018121003', 'Umar Abdillah Wicaksono', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('200102232019121001', 'Ilham Yusuf Pratama', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('200105182019122001', 'Meilena Cahya Wulandari', (select id from public.units where code = 'SEC-18'), 'staff', null, true, true),
  ('197503131996021001', 'Omben Subarlih', (select id from public.units where code = 'SEC-19'), 'section_head', null, true, true),
  ('199207272012101003', 'Bagus Dwi Priantoro', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199309032013101002', 'Helmy Agiel', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199311242013101002', 'Bayu Eka Setya Dharma', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199505112018012004', 'Chintya Megantari', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199604152015021001', 'Shidiq Maulana', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199605282015122003', 'Putri Ratriani', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199901092018121002', 'Ibnu Mahdilevi Soeparto', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199902212019122001', 'Ajeng Pramesti', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('199903252021011001', 'Enos Saen A. Ginting', (select id from public.units where code = 'SEC-19'), 'staff', null, true, true),
  ('196809241989121001', 'Pantjoro Agoeng', (select id from public.units where code = 'DIV-06'), 'division_head', null, true, true),
  ('197706171999031001', 'Agus Madi', (select id from public.units where code = 'SEC-20'), 'section_head', null, true, true),
  ('199505142018011003', 'Denis Arif Pambudi', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199602132016121001', 'Febri Rusiana Ramadan', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199612212015122001', 'Mega Puji Anggraeni', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199703222019121001', 'Kukuh Muhammad Farhan Putera', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199706222016121001', 'Ammin Mubarok', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199710072019121002', 'Muhamad Ichwan Nurdin Ashari', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199807292021011002', 'Ahdaf Ginala Amri', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199808222019121001', 'Rus Hertanto', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199809062018011001', 'Ragas Aziz Kurniawan', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199905252021011001', 'Rizky Adi Nugroho', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('199910192019122001', 'Riri Okta Purbaningsih', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('200009172019122001', 'Rizky Wati Situmorang', (select id from public.units where code = 'SEC-20'), 'staff', null, true, true),
  ('198008162002121001', 'Agus Setiawan', (select id from public.units where code = 'SEC-21'), 'section_head', null, true, true),
  ('198911232013101001', 'Anton Tabah', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199204052014111002', 'Agung Nurwanto', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199210152013101009', 'Prana Putra Dewa', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199505172015021003', 'Al Kafi Samhan Hindami', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199512132015121002', 'Hafizh Ar-Rasyid Kuncoro', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199606192015121004', 'Muhammad Arbi Juniarto', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199609222018011002', 'Jona Galatians Pandiangan', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199709032016122001', 'Fildatun Ni`Mah', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199805132018121002', 'Aditya Widi Saputra', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199806152018012003', 'Dona Krisanti', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199807242018011001', 'Muhamad Nurpasca Ageng Tama', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199901122021011002', 'Nur Rohman Fauzan', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199905032019121002', 'Irwan Tambunan', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199906232018122002', 'Nindya Pini Siwi', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199910022019122001', 'Greiny Widia Damanik', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('199912072018122001', 'Kamilia Bilqis Adillah', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('200007142019122002', 'Lisna Riyani', (select id from public.units where code = 'SEC-21'), 'staff', null, true, true),
  ('197404211994021001', 'Fuad Muftie', (select id from public.units where code = 'SEC-22'), 'section_head', null, true, true),
  ('199003302010011005', 'Heldi Ramadhan', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199208152010121001', 'M.Dio Ariansyah', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199501172015021001', 'Shubkhi Dzulfiqar', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199503182015021001', 'Fahmi Wahyu Trihasno', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199505292015021001', 'Fajar Priyo Utomo', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199510202015121003', 'Sholih Adam Gumelar', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199511282015121003', 'Olioz Novan Rio Zakaria', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199512052015021001', 'Fiqri Alfarisi', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199603092016122001', 'Ines Karina', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199803232018121001', 'Toni Arie Gilang Wicaksono', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199807272018011003', 'Irfan Taufiq Sudiro', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199907082018121001', 'Delvis Arif Maulana', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199907112019122001', 'Nabila Abdurahman', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199910262018121001', 'Fikri Fathin Daulay', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('199912122019122002', 'Hilma Taqi Muljami', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200003102019121002', 'Rif`At Reyhansyah', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200004172019122001', 'Dyah Widyasari', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200005082019122001', 'Meiriska Hariwinto', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200008162018122001', 'Assyfa Nabila Gamal', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200008172018121001', 'Haqzen Aulia Mahardhika', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('200008262019122001', 'Amariza Laraswati', (select id from public.units where code = 'SEC-22'), 'staff', null, true, true),
  ('197308191996021001', 'Agustinus Rahmad Subagyo', (select id from public.units where code = 'SEC-23'), 'section_head', null, true, true),
  ('198105262005011001', 'Ardhyaloka', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199003232010121009', 'Panji Witoko', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199010052012101003', 'Harris Maulana Rakasiwi', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199602172018011002', 'Royska Nurmuhammad', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199609132015121003', 'Rifaldo Bontong', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199701252018011002', 'Bebeto Ramadhan Saputro', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199806112019121002', 'Felix F Sidabalok', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('199809072021011001', 'Teguh Pratama', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('200002202018122001', 'Egta Ayu Fadhlillah Sugiarto', (select id from public.units where code = 'SEC-23'), 'staff', null, true, true),
  ('196808191988121001', 'Agustyan Umardani', (select id from public.units where code = 'DIV-07'), 'division_head', null, true, true),
  ('197708252000011002', 'Heri', (select id from public.units where code = 'SEC-24'), 'section_head', null, true, true),
  ('199503022015021001', 'Wisnu Suryo Negoro', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199605232015122002', 'Geby Marselvia', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199607152015121001', 'Yulian Asiddiqi Kahfi', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199702222015121002', 'Rayhan Muhammad Fadilah', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199707042018011002', 'Indra Bagas Pratama', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199808282021011001', 'Syafrizal Bagus Firmansyah', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199812082018011001', 'Muhammad Farhansyah Yavi Anwar', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199903012019121001', 'Nadhif An Naufal Azmi', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199907132018122002', 'Tri Evi Yulianty', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('199910252018121002', 'Abdur Rahman', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('200004082019122001', 'Dewi Puspitasari', (select id from public.units where code = 'SEC-24'), 'staff', null, true, true),
  ('197406141994031002', 'Efan Sandy Akbar', (select id from public.units where code = 'SEC-25'), 'section_head', null, true, true),
  ('199406282015021002', 'Fahmi Hayyat', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199408042015021001', 'Virqly Khaybar', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199409052015021002', 'Bennarivo', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199410152015021003', 'Trishna Yodi Pratama Sinambela', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199412012015021002', 'Nady Herdian Brahmantya', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199501172015021002', 'Falsafa Amal Islami', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199501242015021001', 'Yanualdo Ramagesti Shanenda', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199502042015021008', 'Mohammad Fauzan Ramadhan', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199502072015021002', 'Muhamad Wildan Ramdhani Ardliyanshah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199502142015121001', 'Zuhud Ruhullah Moussovi', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199503132015021003', 'Irfan Satriyo Darmanto', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199503142015021003', 'Egi Pratama', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199504042015021004', 'Ilham Annas Darajat Rahmatullah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199505102015021002', 'Bagja Baharudin', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199506042015021002', 'Richy Naviri Putra Panjaitan', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199507012015021007', 'Zain Chanor S.', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199507142015021006', 'Andrean Julio Firmansah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199507212015021008', 'Fathul Huda', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199508102015021001', 'Rido Fadila', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199508112015021001', 'Purnomo Gumelar', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199509112015121002', 'Alwi Shihab', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199510082015021002', 'M. Adriel Oktodio Pratama', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199601212018011004', 'Dimas Ismanuari', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199602132015021002', 'Rahmat Nugroho', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199603242015121003', 'Rusydan Fauzani Akbar', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199604062015121002', 'Romy Aricson Bijak', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199604192015121002', 'M. Nizar Adiatma Pratama', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199605142015121003', 'Yusuf Nur Faisal', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199609042015121003', 'Neza Panji Volorous', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199610062016121001', 'Muhammad Ardin Prawira', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199611022015121002', 'Denny Arief Firmansyah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199611142015121005', 'Yoan Agung Wibawa', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199701132016121002', 'Riza Ilham Arifin', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199702102018011004', 'Faizal Fahmi', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199702142018011001', 'Sacmika Sabiila Rozaq', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199703082016121001', 'Getar Samodra Alfaridzi', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199703122018121001', 'Denis Setya Putra', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199704022015121001', 'Frampy Gilbert Lesnussa', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199705082016121001', 'Rizky Mei Herlias', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199706032016121002', 'Muhamad Nail Iman', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199706122016121001', 'Yoga Prabandanu', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199708312016121001', 'Ilham Satiya Wijaya', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199709052018011001', 'Adeseptian Risberg Prakoso Silaen', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199710312018011002', 'Lucky Yohans Gultom', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199711102016121001', 'Aji Pandu Mahardika', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199711262016121001', 'Fadhil Muhammad Yusuf', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199712022016121001', 'Farid Al Qadri', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199802032018011002', 'Febriza Raditya Pradipta', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199802232018011001', 'Oscar Halomoan Panggabean', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199803082018011002', 'Muhammad Rifki Adriyanto', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199804282018011003', 'Viky Fadila Surya Iranto', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199805042018011001', 'Muhammad Bobby Firmansyah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199805212019121001', 'Bima Budi Wibowo', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199806112018121002', 'Dimas Probo Dewanto', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199807062018121001', 'Mohammad Aryo Dwi Pangestu', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199807122018121001', 'Carlos Trias Kinan Risyad', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199808252018011001', 'Muhammad Widyan Zulal', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199809192018121001', 'Sidik Ramadi', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199809242018121001', 'Candra Ardika', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199812132018011003', 'Firhan Bayu Adiyuana', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199812272018011001', 'Ramadhan Prima Kuntjoro', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199901062018011001', 'Muhammad Pangayom Nan Tabah Handoyo', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199901242018121001', 'Andisena Panggah Komandona', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199902122018121001', 'Farizaki Adhilleon', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199903032018121001', 'Frassetyo Marizda', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199903032018121002', 'Danar Duha Fadila', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199903052018121001', 'Ilyas Ezar Saputra', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199903122019121001', 'Nalindra Hanung Brillianingtyas', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199904292018122002', 'Indah Jaya Br. Saragih', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199905072018011001', 'Haidar Farras', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199905242018121001', 'David Perdana', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199906162019121001', 'Wirandika Aulia Hasan', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199906192019121001', 'Bayu Putra Feranda', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199907152018121001', 'Adiaksa Mardaup Simanjuntak', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199907222018121002', 'Farizh Alfarizhi Hidayat', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199907302018121005', 'Farid Bagus Megantara', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199909012018121001', 'Vikrama Krishna', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199910282018121001', 'Pandu Wachyu Pamungkas', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199910292018121003', 'Abdul Hafidz Yogatama', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199911152018121002', 'Dhimas Probo Kuncorojati', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199911172019121002', 'Alwi Tondi Rifaldo Harahap', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199912192018121001', 'Kukuh Hidayat', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('199912282018121003', 'Odas Gujarat', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200003222019121003', 'Muhammad Alwi Rofiqi', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200003272018121001', 'Firman Ibrahim', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200004142018121001', 'Rio Brata Eka Husada', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200007252019121001', 'Fadhil Yulianto', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200008092019121001', 'Anggid Aji Wicaksono', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200009012019121001', 'Akbar Saefullah', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200009302019121001', 'Muhammad Rizki Ardhana', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('200012012019122001', 'Ryzka Ramadhani', (select id from public.units where code = 'SEC-25'), 'staff', null, true, true),
  ('197702151999031001', 'Rosyidan Syah', (select id from public.units where code = 'SEC-26'), 'section_head', null, true, true),
  ('199708072018011002', 'Mohamad Bani Alifianto', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('199809212018011001', 'Ahmad Safi`i', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('199812012018122001', 'Nugroho Laraswati', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('199906052018011001', 'Finka Fajar Abdillah', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('199909052018122002', 'Fanny Fatikasari', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('200001152019122003', 'Isnaini Fadhilah', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('200003062019121001', 'Fanfan Rudi Afandi', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('200008092019122001', 'Marinda Min Amrina Rosada', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('200009072019122002', 'Septya Ningrum Rahayu', (select id from public.units where code = 'SEC-26'), 'staff', null, true, true),
  ('197304121993022001', 'Alvina Christine Zebua', (select id from public.units where code = 'SEC-27'), 'section_head', null, true, true),
  ('199406102015021001', 'Graha Putra Kusuma', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199502232015021003', 'Febri Ramdhan Sidiq', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199506272016121001', 'Ahmad Arya Kusuma', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199603192016121001', 'Firmansyah Dian Pradana', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199603222015121003', 'Viki Bayu Prasetya', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199605182015121001', 'Muhamad Faisal Fatih', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199609072015121002', 'Muhammad Ilman Shiddiq', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199808042018011001', 'Bagoes Triantoro Ajie', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199808082021011001', 'Burhanudin Abdillah', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199809032018012003', 'Imtihan Legati', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199812212019122001', 'Nia Susilowati', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199901102018121001', 'Fauzan Maulana', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199904132021011001', 'Fadlan Bani Husna', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('199908272019121001', 'Akbar Tanjung', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('200004122019122001', 'Agnes Dwi Pratiwi', (select id from public.units where code = 'SEC-27'), 'staff', null, true, true),
  ('196912251996031001', 'Toto Raharjo', (select id from public.units where code = 'DIV-08'), 'division_head', null, true, true),
  ('197512081996031004', 'Erwin Bangun Maruli Tua', (select id from public.units where code = 'SEC-28'), 'section_head', null, true, true),
  ('199412142015021001', 'Odie Mahendra', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199505092015021003', 'Harits Manazili', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199511102016121002', 'Pahlevi Adhi Nugraha', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199612092016121001', 'Gagat Ridwan Wicaksana', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199702012018121001', 'Ananda Alfin', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199706042018012001', 'Nurul Khafifah', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199710162018012001', 'Nanda Oktaviani Putri', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199711082018011001', 'Andro Java', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('199807142021011003', 'Muhammad Ardan Felani', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('200001082018122001', 'Bella Fitri Melinia', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('200001142019121002', 'Rizky Millenio Jagantina', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('200002012021011001', 'Fadhilah Achmad Riano', (select id from public.units where code = 'SEC-28'), 'staff', null, true, true),
  ('197202111998031002', 'Budi Satria', (select id from public.units where code = 'SEC-29'), 'section_head', null, true, true),
  ('198803212007011001', 'Andre Karel Lewis', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199007182010011001', 'Ridho Syahputra', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199106222012101001', 'Ignasius Dwi Ariputra Lengkong', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199201182010121004', 'Firman Setiawan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199201252013101003', 'Pandu Mahasyah', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199207162013101001', 'Suhastomo', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199209282012101001', 'Ludfan Kusuma Wardani', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199310032012101001', 'Ocktopardomuan Sidabutar', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199311292013101001', 'Ganes Yudha Dwiyan Prakasa', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199401042015021001', 'Faizal Armandha', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199409172015021002', 'Mahaji Suryoyudanto', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199501092015021002', 'Gerrid Enggar Pracoyo', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199501132015021003', 'Difi Triyandi', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199505132018011002', 'Arrizkana Hadi Yayan Nurrosyd', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199505252015021002', 'M. Muin Tanzil', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199506142015021001', 'Azkha Kurnia Indrajaya', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199506252016121002', 'Nur Utomo Taufiq', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199506292015021001', 'Muhammad Imam Nawawi', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199508072015021002', 'Muhammad Nursetiawan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199508132016121001', 'Muhammad Ilham Budiman', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199510152015021002', 'Muhammad Nizar Arifullah', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199511012015121001', 'Muh. Fakhrul Wahyu Pratama Murad', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199609182015121002', 'Muhammad Amran Wahiduddin Nurlatifa', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199611172015121002', 'Rais Asif Hamdan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199703272016121001', 'Muhammad Chandra Kumar', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199805072018121002', 'Ronaldo Sanjaya Hutagalung', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199807012021011003', 'Agung Pratama Putra', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199809112018011002', 'Naufal Farhan Nefawan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199809142018121001', 'Puput Pujiono', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199812182018011002', 'Galih Cahyo Kuncoro', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199812252018121001', 'Kristian Jimmy Hamonangan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199904072021011001', 'Rangga Pramusinto', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199908032018121001', 'Patrick Wisnu Aryoputro', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199909052021011002', 'Agung Satriyo Wibowo', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199909182018121001', 'Muhammad Naufal Ariq', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('199911172019121001', 'Dwiky Vendanata', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('200001032021011001', 'Fakhri Ridwan', (select id from public.units where code = 'SEC-29'), 'staff', null, true, true),
  ('197406091994021001', 'Yunianto Laesul Afdol', (select id from public.units where code = 'SEC-30'), 'section_head', null, true, true),
  ('198204052005011001', 'Lucky Imawan Supriyadi', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199004222012102001', 'Melita Yuniarti', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199202162013101001', 'Syahid Izzudin', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199202192012101003', 'Ian Kurnianto', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199411142015021003', 'Ahmad Mahfuzh', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199502132015021002', 'Taruna Bagus Harefa', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199503092016121002', 'Reza Tri Anugrah Putra', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199503092016121003', 'Ahmad Mustaqfirin', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199503192015021003', 'Yudha Kusumawardhana', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199510042015021007', 'A. Imam Asyraf Amiruddin', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199512032016121001', 'Ryan Maliki', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199602252018011002', 'Febrian Putra Wibawa', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199603092015021001', 'Syafiq Akrom', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199605062015121004', 'Rendra Rezki Purwandani', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199801282018012001', 'Belinda Fitri Susatyo', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199804082019121001', 'Muhammad Sulthoni Ulumuddin', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199811102021011001', 'Mochammad Kemal Thaariq', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199905292018121002', 'Satriya Abdi Nugraha', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('199908262019122001', 'Rizki Fahrianti', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('200001142019122001', 'Maulania Yustika Bahaji', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('200104212019122001', 'Ratna Widianti', (select id from public.units where code = 'SEC-30'), 'staff', null, true, true),
  ('197107191992011001', 'Wahyu Setyono Widyobroto', (select id from public.units where code = 'SEC-31'), 'section_head', null, true, true),
  ('199204132013101001', 'Agus Tri Pijanarko', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199308102013101004', 'M Rizaldi Prasetya', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199409012015021004', 'Septian Delta Andreas', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199412102015021003', 'Rigid Eka Pambudi', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199501222015021002', 'Ginanjar Rijal Rosyid', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199510152018011002', 'Anggit Sasmito Kresno', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199512012015021001', 'Anggit Dimastiko Hidayat', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199611042015121004', 'Irhamsyah', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199704142018121001', 'Wahyu Budi Santoso', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199711202019121001', 'Hervista Bagas Pratama', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199811012021011001', 'Yudha Pratama', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('199812062018012001', 'Desy Florencya', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('200004122019121001', 'Ahmad Ziyad Rizqullah', (select id from public.units where code = 'SEC-31'), 'staff', null, true, true),
  ('197010011990121001', 'Kurnia Saktiyono', (select id from public.units where code = 'DIV-09'), 'division_head', null, true, true),
  ('197409301994021001', 'Lanang Dwi Wirawan', (select id from public.units where code = 'SEC-32'), 'section_head', null, true, true),
  ('198508012003122003', 'Ati Utami Indrawati', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199407042015021004', 'Julio Hagi Kennedy', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199512132016121001', 'Dimas Rizal Marindis', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199602102015122001', 'Wella Febryantos', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199605292018011006', 'Maulana Irfan Afif', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199710252019121001', 'Fajar Ikramsyah', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('199911082018121001', 'Sigit Nurcahyono', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('200005232019122003', 'Ozora Nabilah Sari Izdihar', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('200010242019122001', 'Yumna Alifatun Arda', (select id from public.units where code = 'SEC-32'), 'staff', null, true, true),
  ('197710252002121001', 'Khairudin', (select id from public.units where code = 'SEC-33'), 'section_head', null, true, true),
  ('199408072015021001', 'Rangga Anugrah Rachmaputra', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199411112015021003', 'Ali Reza', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199507042015021003', 'Aditya Cahya Widiyono', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199507092015021004', 'Acep Nurul Hakim', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199507192015021004', 'Danang Pringgo Kisworo', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199510132015021004', 'Ilham Insani Isma Rozaq', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199511122015021001', 'Muhamad Agung Fathurrahman', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199803102018011001', 'Aldian Swanida Jayaha', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199804292018011001', 'Alan Prihariyanto', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199901192018121002', 'Riyan Andi Laga', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199905162018121003', 'Mohamad Wahyu Rizal Adnan', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199908252018122004', 'Elisa Linetta', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('199909182019122003', 'Revita Monica M. Simanjuntak', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('200007022019122004', 'Rika Boru Panggabean', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('200105022019121001', 'Ifan Rizki Pramudita', (select id from public.units where code = 'SEC-33'), 'staff', null, true, true),
  ('197603241996031001', 'Andrianto', (select id from public.units where code = 'SEC-34'), 'section_head', null, true, true),
  ('199205272012101001', 'M. Taufan Arif Putra', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199207042012101002', 'Faisal Muharman', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199407082015022001', 'Yulita Amalia Triyana', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199412232015021002', 'Mochammad Dozen Khaya', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199505212015021001', 'Herdian Rifqi Priharso', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199505282015021005', 'Ahmad Zaid', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199507112016122002', 'Rifqi Yulan Husnia', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199604272015121005', 'Kasbia Putra Pamungkas', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199605102015121001', 'Anas Fahrizal', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199608172015021001', 'Rangga Prasasti Herjuna', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199807162018122001', 'Ninda Maulidiasari', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199811272018011002', 'Hendra', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199904282018122001', 'Ersa Adisty Yuniar', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199906022018122001', 'Versy Tasia Geraldine', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('199911172019122001', 'Imrotus Sholeha Zaqiya', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('200002052019122001', 'Clarisa Putri Viandini', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('200008162019122002', 'Nia Alda Dewi', (select id from public.units where code = 'SEC-34'), 'staff', null, true, true),
  ('197502031999031003', 'Samid', (select id from public.units where code = 'SEC-35'), 'section_head', null, true, true),
  ('199208292015021001', 'Tubagus Muhammad Dimas Kurniawan', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199407212015021004', 'Candra Nurrochman', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199601062015021001', 'Muhammad Zaky Firdaus Andipratama', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199611072016121001', 'Fiqih Komarullah', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199710242018121001', 'Muhammad Fauzi', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199711272018011003', 'Ilham Muhamad Fajar', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199812112021011001', 'Hatta Khoiruka', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199904282018121002', 'Muhammad Abdul Aziz Subagyo', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199906062018121001', 'M. Agung Triwijaya', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('199908012019122001', 'Febi Alisia Fransiska Sagala', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('200001082018122002', 'Safitri Anis Setya Panarima', (select id from public.units where code = 'SEC-35'), 'staff', null, true, true),
  ('198003012001121001', 'Ichlas Maradona', (select id from public.units where code = 'DIV-10'), 'division_head', null, true, true),
  ('198203202003121001', 'Sigit Tri Hatmoko', (select id from public.units where code = 'SEC-36'), 'section_head', null, true, true),
  ('199105052012101001', 'Bismi Mohammad Alief Fathansyah', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199208222012101001', 'Muhammad Danu Fiza', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199209152013101002', 'Muhsin Z Muchtar', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199305222013101001', 'Jaka Utama', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199403112015021002', 'Abdurrahman Basyiruddin Rabbani', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199403312015021001', 'Faisal Saputra', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199404022015021002', 'Aditya Gumay', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199404272015021001', 'Nico Alvindra', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199409302018011001', 'Mochammad Divin Vidyaka Putra', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199501302015021001', 'Rico Fernando Sitompul', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199505152015021002', 'Anggi Kurnia Fajri', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199505312015021001', 'Rizqy Nur Wahid', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199508042015021004', 'Abidzar Almadira', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199508092015021004', 'Jouwardy Alfredo', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199509022015021001', 'Ilham Akbar', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199511102018011003', 'Muhammad Faesal Hudha Simanjuntak', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199512032015021001', 'Mesua Smithian S Depari', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199602032018011004', 'M. Rezky Ramadhan', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199603032018011002', 'Hermansyah Eko Widianto', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199603112018011002', 'Riyanda Taufiqurrohman', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199607262015121001', 'Diaz Bayu Samudra', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199608132018121001', 'Muhamad Anandi Kharisma', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199609042018121001', 'Restu Galih Sajati', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199611172016121001', 'Riung Garendra Argatama', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199701082015121001', 'Muhammad Fadhila Azzam', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199703132018121003', 'Rakhmat Habibie', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199703242018011001', 'Chandra Try Youdhana', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199705272019121001', 'Sulthan Van Diori', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199708032018121001', 'Agung Prade Simanjuntak', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199708062018121002', 'Langlang Hartanjaya', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199711082018011003', 'Nazal Amrullah', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199711262018011002', 'Yaga Dewantara', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199711262018121001', 'Bimo Dwi Noviantoro', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199803232019121001', 'Wahyu Dwika Mahendra', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199804222021011001', 'Ilyas Hanafi', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199806052018121002', 'Muhammad Ilham Akbar', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199809072019121001', 'Firman Zidane Shahroni', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199809192019121002', 'Arief Ahmad Abdul Azis', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199904092018121003', 'Denniz Decaprio Aprildo', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199904172021011001', 'Yusuf Iqbal Maulana', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('199905102021011001', 'Daffa Alfana Prananda', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('200005122019121001', 'Muhammad Faris Pratama', (select id from public.units where code = 'SEC-36'), 'staff', null, true, true),
  ('197806142000121001', 'Didik Mujiyono', (select id from public.units where code = 'SEC-37'), 'section_head', null, true, true),
  ('198210152003121003', 'Lukman Kamil', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('198908062010011002', 'Zulda Agusta', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('198911242010011001', 'Dio Yogasidi', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199009292014111001', 'Pramadya Purwa Andhika', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199011082010011001', 'Rendra Mulya', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199101202010011001', 'Eko Anestio Darmadi', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199203302013101002', 'Chandra Irawan Alfianto', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199303052013101002', 'Dita Riyan Dwi Saputra', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199306032013101004', 'Hendro Widhiyarto', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199309032013101001', 'Septian Markus Esa', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199311232015021001', 'Wisnu Wijanarko', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199311302015021001', 'Waskita Wahyu Murti', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199501042015021001', 'Hamzah Kalam Sumeru', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199501082018011002', 'Dimas Narendra Anwar', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199501112015021002', 'Aris Dharmawan', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199504192015121001', 'Dio Permadiyanto', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199504292015021002', 'Fajar Triantoro', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199506252015021003', 'Bagus Setiawan', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199604032015121002', 'Ganang Rusdian Apriaji', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199607042016121001', 'Brevy Fauzi', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199607172015121002', 'Mohammad Ifan Adinsyah Taufik', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199607282015121001', 'Machrobi Yulianto', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199611042018121001', 'Rizki Rachmad Hidayat', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199612262015121001', 'Naufal Hafif Hasibuan', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199801282018011002', 'Reza Alvin Ramadhan', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199804142018012001', 'Gita Pratiwi', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199806122018011002', 'Muhammad Fathi Fawwaz', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199807222018011003', 'Yusuf Ihsan Haqqi Habibbullah', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199906272018121001', 'Eghatama Sadriansyah', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199907292018121003', 'Ananda Alrif`an Ritonga', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('199909172021011002', 'Haris Suganda', (select id from public.units where code = 'SEC-37'), 'staff', null, true, true),
  ('198311122009011007', 'Andri Noverianto', (select id from public.units where code = 'SEC-38'), 'section_head', null, true, true),
  ('198612212008121001', 'Rocky Paruhum Siahaan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199109172012101001', 'Imawan Avicena', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199204292012101001', 'Paulus Wijaya Sitorus', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199205142013101001', 'Bayu Nurmay Prazyugi', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199308182013101002', 'Gilang Aditya Utama', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199402152015021002', 'Sendi Sela Putra', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199404032015021001', 'Fariz Fadhilah Hasya Aghnia Azta Hasibuan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199410232016121002', 'Anugrah Eka Prasetya', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199502202015021001', 'Robi Febrianto Irawan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199502242015021004', 'Febri Ramadhan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199504082015021001', 'Muhammad Hanif', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199508072015021001', 'Akmal Setiawan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199510042016121001', 'Pranata Fransiskus Tarigan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199512152015121002', 'Ruben Yonathan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199601012018011003', 'Ahmad Fajri', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199606032015021001', 'Helmi Kurniawan', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199709172018011003', 'Dimas Adicahyo Prakoso', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199712212018121001', 'Yogi Pratama Putra', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199807172019121001', 'Jeremy Tolopan Martua', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199902052021011002', 'Stefi Sutha Syawal Dilapanga', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199902102019122001', 'Fernanda Pramudianti', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('199911022018121003', 'Indra Karunia Putra', (select id from public.units where code = 'SEC-38'), 'staff', null, true, true),
  ('197409041995031001', 'Andy Gunawan', (select id from public.units where code = 'SEC-39'), 'section_head', null, true, true),
  ('198704162007101001', 'Dede Raya Apriyanto Saragih', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199106032014111001', 'Aditya Juniyansyah', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199107192010121002', 'Muhlisin', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199211202012101001', 'Wildan Rosyadi', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199311102015021002', 'M. Ronny Meliala', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199311282013101001', 'Ahdiat Rizkiansyah', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199402222015021002', 'Bisma Sandhi Yudhanto', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199501312018011003', 'Nur Cahyono', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199502182015021004', 'Aditya Bayu Wicaksono', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199505182015021004', 'Andre Fernando L Pontoh', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199506142015021002', 'Achmad Kamil', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199510302015121003', 'M. Tegar Abadi', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199511052015121002', 'Cakra Bhayu Mardani', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199511252018011003', 'Benny Wiranto', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199512302015121002', 'Silvester Danu Dirgantara', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199601192015021001', 'Samuel Edward Pangondian', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199604012015021002', 'Brian Galih Syifaul Qulub', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199604292018011003', 'Johannes Hasoloan Simanjuntak', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199605172018011002', 'Muhamad Tegar Damanta', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199609042018011001', 'Ade Wardiman Suparman', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199703242016121002', 'Ariqhi Sinatria Haryawan', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199710112018121001', 'Mahardhika Aryadhana Jattin', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199710302016121001', 'Joshua Putra Sedianto', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199801022016121001', 'Ahmad Rizky Ramadhan', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199804042021011001', 'Shofwan Palwisaputra', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199805072018011003', 'Kevin Saghiira Hermawan', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199805132018012001', 'Rizky Ayu Puspa Ningrum', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199808262018011002', 'Daniel Amaarta Hasiholan Silaen', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199809172019121001', 'Bosman Hasudungan Simanullang', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199901092018121001', 'Moch. Yoga Pradana', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199902232021011002', 'Ramadhika Rizqi Wirambara', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199903152018011001', 'Tegar Pambudi', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('199907062018121002', 'Brilian Surya Ramadhan', (select id from public.units where code = 'SEC-39'), 'staff', null, true, true),
  ('197702152000011001', 'Henki', (select id from public.units where code = 'SEC-40'), 'section_head', null, true, true),
  ('198506072003122002', 'Djulaeha', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('198911172012101002', 'Hendra Hafiz Kartasasmita', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199112162012101001', 'Chandra Kusuma Yuda', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199207112014111001', 'Galih Justhian Wirajoyo', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199304102013101002', 'Noorman Aditya', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199308312013101001', 'Syaiful Afandy', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199312052015021001', 'Mufqi Ardio Alifnanda', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199312212015021001', 'Aulia Anshari', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199502272016121001', 'Multazam Ramadani', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199507222016121001', 'Caesar Sandyarino', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199509132015121002', 'Mohamad Budi Santoso', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199512212015121001', 'Mochammad Nizar', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199601202016121002', 'Ryan Chaidir Kamarullah', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199801152018011001', 'Renaldy Anicetus', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199804022018121002', 'Mohamad Zulhadi Syauqi Alfarisi', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199809162019121001', 'Bagas Septio Surya Saputra', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199902252021011001', 'Tito Bayu Noor Pambudi', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199910162019121001', 'Dian Oktavian Sanjaya', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('199910292018121002', 'Haris Kurniawan', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('200002072019121001', 'Muhammad Daffa Muzhaffar', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('200006292019122001', 'Alfia Nurrul Izza', (select id from public.units where code = 'SEC-40'), 'staff', null, true, true),
  ('197711112000011003', 'Samino', (select id from public.units where code = 'SEC-41'), 'section_head', null, true, true),
  ('199101282010011001', 'Lucky Agung Firdaus Hidayat', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199502232015021002', 'Erri Fadjar Pradhana', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199503052016121002', 'Raihan Erfanda Zantio Nasution', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199504112015021004', 'Adhifa Mizan Ghifary', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199511252018011006', 'Dipdha Saptagita Pupadewa', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199905052018122004', 'Siti Fadhilah Humairah', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199905192021011002', 'Michael Reza Vrederick Sihaloho', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('199905292021011001', 'Jonathan Alex Rumahorbo', (select id from public.units where code = 'SEC-41'), 'staff', null, true, true),
  ('197605261996021001', 'Suhartoyo', (select id from public.units where code = 'SEC-42'), 'section_head', null, true, true),
  ('199011122013101002', 'Baiquni Aryodamar', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199502082015021003', 'Rachmad Hidayat', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199504012018012001', 'Nurdiana Riyadlatul Jannah', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199602142018011005', 'Febrian Luxfi Harada', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199610162018011002', 'Andrey Doohan', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199811092018121001', 'Dwi Mayasin Al Rasyid', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('199812272021011001', 'Ken Abdulah Aziz Romadhoni', (select id from public.units where code = 'SEC-42'), 'staff', null, true, true),
  ('196901151996031001', 'Andi Hermawan', (select id from public.units where code = 'DIV-11'), 'division_head', null, true, true),
  ('197004031996032002', 'Roslindawaty Br. Ginting', (select id from public.units where code = 'SEC-43'), 'section_head', null, true, true),
  ('199403202013101002', 'Muhammad Vito Jati Pascalri', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('199503182015021002', 'Aidil Vachri', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('199506012015021001', 'Moch Arif Fauzan', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('199508292015021002', 'Muhamad Rais Alvin', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('199907252018122001', 'Nabylia Adhara Laksita', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('199912162019121001', 'Aldira Wahid Ramadhan', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('200005242019122002', 'Meilya Rizky Utami', (select id from public.units where code = 'SEC-43'), 'staff', null, true, true),
  ('197308061994022003', 'Agustina Marpaung', (select id from public.units where code = 'SEC-44'), 'section_head', null, true, true),
  ('199202152013101004', 'Dinar Suherman', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('199502212018011004', 'Ferrian Faiz Ramadhan', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('199510122016121001', 'Christophorus Gian Baskara', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('199802122018011002', 'Dedy Purwanto', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('199809032018122002', 'Anita Lasnauli', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('199901272021011002', 'Naufal Fadhlurrochman', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('200003162018122001', 'Ida Nur Hajah', (select id from public.units where code = 'SEC-44'), 'staff', null, true, true),
  ('197201211996032001', 'Yenni Rachmawati', (select id from public.units where code = 'SEC-45'), 'section_head', null, true, true),
  ('199505092015121002', 'Agung Maulana Syahrudin', (select id from public.units where code = 'SEC-45'), 'staff', null, true, true),
  ('199511042016122001', 'Khansa Ufairah', (select id from public.units where code = 'SEC-45'), 'staff', null, true, true),
  ('199807062018011001', 'Reyhan Mauludi Arissaputra', (select id from public.units where code = 'SEC-45'), 'staff', null, true, true),
  ('199807082018122001', 'Sylraini Pramono Putri', (select id from public.units where code = 'SEC-45'), 'staff', null, true, true),
  ('200007242019121001', 'Fairuz Dinar Aqshal', (select id from public.units where code = 'SEC-45'), 'staff', null, true, true),
  ('197305221992121001', 'Erli Haryanto', (select id from public.units where code = 'SEC-46'), 'section_head', null, true, true),
  ('198512302006041002', 'Muhamad Ilham', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('198802262009121006', 'Aloysius Lambok Siahaan', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199207142013102001', 'Kartiyasa Arifta Putri', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199302262013101004', 'Bayu Tri Widodo', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199404222015021003', 'Abdul Rahman Zain', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199503112015021003', 'Fithra Abdulloh Salsabila', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199505302015021003', 'Kurniawan Cahyo Nugroho', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199512152015021002', 'Brillian Widiyanto', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199602122015122003', 'Najah Dhiya Kamilina', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199705092018121002', 'Muhammad Abil Zebramawi', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199712072018012002', 'Rahmatul Ulfa', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199712092018012001', 'Syifa Salsabila', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199712102018011002', 'Aziz Ghani Wahyu Santoso', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199712282018122003', 'Ristya Nurharisa', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199802112016122001', 'Exa Febrin Manurung', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199808132018121001', 'Aditya Indrawan Prasetyo', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199810092021011001', 'Wahyu Guntoro Adi Susanto', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199901212018121001', 'Ilham Widhi Pamungkas', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199908312018122001', 'Rifka Agustina', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('199911102018122001', 'Milka Jeges T.M. Br. Sihombing', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('200002102022011001', 'Yogi M. Simanjuntak', (select id from public.units where code = 'SEC-46'), 'staff', null, true, true),
  ('197009301990121002', 'Hari Prabowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197201161992121001', 'Jonathan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197207281999031001', 'Erwindra Rachmawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197304171992121001', 'Ikbal', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197312301994021001', 'Mufti Widadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197402051994021001', 'Sugiono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197406131994021002', 'Agus Sulistyo E.S.', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197408251997031002', 'Diding Toyibudin', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197501121997031002', 'Sujadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197502141994022001', 'Widia Hastuti', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197503031997031001', 'Shodikin', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197503081995031002', 'Panangian Mangaratua Marpaung', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197504102005011001', 'Eko Wigiyanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197505071997031003', 'Akhmad Yadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197505302005011001', 'Hasrul Mahyuzar Melayu', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197507091997031001', 'Budi Darmanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197507092005011001', 'Purnomo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197508161996021001', 'Arief Ferdian Syah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197510301999031001', 'Supriyono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197511251999031001', 'Noviardi Hidayat', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197601111997031001', 'Edie Purwanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197601161996021003', 'Arri Wisnu Tri Kumoro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197602271997031001', 'Andi Prabawa', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197603301996021003', 'Muhammad Sulaiman Dasril', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197604281999031001', 'Yanwar Ariyadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197605091998031002', 'Sudarso', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197606292003121001', 'Toto Purwanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197607281998031001', 'Muhamad Rijal', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197608051997031001', 'Agus Supriatna Pongoh', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197608201999031002', 'Muhammad Sholakhudin', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197608281996021002', 'Moh. Deni Ramdhan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197609211996031003', 'Evrizal Mandala Ronaldo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197610011997031001', 'Mohamad Baidowi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197610111997031002', 'Darma Setiawan Saragih', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197701101999031001', 'Selamet Wibisono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197701121999031001', 'Hendrawan Istanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197702031997031001', 'Arif Tri Handoko', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197702051999031002', 'Teguh Purwono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197702052000011001', 'Buana Tugas Sanjaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197702091997031001', 'Akhmad Kuncoro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197704011997031001', 'Martua Afrido Bona Feri Sianturi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197705141997031002', 'Toto Raharjo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197705191997031001', 'Sofyar Banuaraja Ritonga', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197706101998031001', 'Muktar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197706151999031001', 'Endro Eswo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197707281998031001', 'Kristanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197709081997031001', 'Disco Valuasy Harefa', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197709091999031001', 'Agung Wibawa', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197709172000011002', 'Peterus Daimura Silalahi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197710101997031001', 'Momon Rusmono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197712101998031001', 'Muksin Ridwani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197712191997031001', 'Supendi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197801041998031002', 'Ciptono Setia Budi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197801091998031002', 'Elvis Parlindungan Sianturi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197802222000012001', 'Dwi Wahyu Handayani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197802262003122001', 'Kariyani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197803182000011001', 'Eko Budyanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197805062000011001', 'Anton Wirawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197806072003121001', 'Friadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197806191999031001', 'Sugeng Cahyono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197807122000011002', 'Esti Dwi Yulianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197807162003121001', 'Agung Saputro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197808042000121002', 'Darwanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197810101998031002', 'Budi Budiana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197812121998031001', 'Didik Eko Wahyudi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197903042003121002', 'Dodi Cahyadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197904061999031003', 'Sunissan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197905042001121001', 'Budi Santoso', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197906242001121001', 'Teguh Ahmad Ikhsan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197907042000121001', 'Krisna Julianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197908122001121002', 'Heru Agus Widarto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197908181998031001', 'Rd. Agus Suganda', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197908252001121001', 'Ade Fitriansyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197911252000121003', 'Zulkarnain', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('197912182001121001', 'Surata', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198002142003121001', 'Ari Subagyo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198002192001121001', 'Muhamad Hadi Ismanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198002272003121001', 'Iwan Darmawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198004282001121002', 'Muhammad Hanifah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198005122000121002', 'Inwan Manaf', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198006122002122001', 'Virgilia Letare Trishanti', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198006202001121002', 'Yudi Purnama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198007222005011001', 'Sarif Tinambunan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198008272001121004', 'Ardiansyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198010092000121002', 'Luthfi Purnama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198010102000121001', 'Nizar Utama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198010102001121002', 'Sunarko', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198010302002122002', 'Murtini', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198012112002121002', 'Muhammad Shoufran', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198101012003121002', 'Haryono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198101102003122002', 'Eka Rinarti', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198101182003121002', 'Achmad Robani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198101312000121001', 'Abdul Gafur', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198108112002121001', 'Agus Rifa''i', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198110182002121002', 'Sutrisna', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198110232003121001', 'Faried Amrullah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198110272002121001', 'Deddy Mendai Zuhriansyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198111252003121001', 'Dedek Susanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198111292003121001', 'Antony Adhe Sanjaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198202072003121001', 'Rudi Firmansyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198202102001121001', 'Hartono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198202262003121001', 'Ghazali Wijaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198204022003122001', 'Mariah Qibtiyyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198204102001121001', 'Muparrih', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198204162003122001', 'Ema Ratna Sari', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198205262003121001', 'Andy Marwan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198206262003121002', 'Arianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198207042004121001', 'Eko Achmad Santoso', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198207112001121002', 'Farid Najhi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198207212001121001', 'Yayat Ruhiyat', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198210212004121002', 'Daris Purnomo Jati', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198212082003121001', 'Triyanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198304222002121002', 'Arif Wicaksono', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198307282002121002', 'Sudianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198308062010121003', 'Wisnu Widyotomo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198308152004121004', 'Slamet Widodo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198308312009011008', 'Lawrentus', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198310182004121001', 'Akhmad Kosasi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198311262006021001', 'Andi Wahyudi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198401132010121005', 'Rian Aqsa Hermawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198402282007011001', 'Heru Ferdian', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198403262009011006', 'Wahyu Hidayat', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198405202004121002', 'Wiratno', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198407182004121001', 'Faqih Yusuf', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198408122009011005', 'Tatak Suryaputra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198412162006021002', 'Reano', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198501112006021003', 'Akbar Nugraha', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198504172003121002', 'Mohamad Hendra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198504182006021002', 'Arerisza Yuanfala', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198505212010121007', 'Ramses', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198505232007011002', 'Antonius Ade Permana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198507272006021001', 'Munandar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198508162009011005', 'Agus Suprianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198509062007011001', 'Ajar Septian Aditama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198509212010121004', 'Yoga Anggoro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198510252006021004', 'Sehat Daulay', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198510302010121002', 'Budi Dwi Oktianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198511092007011001', 'Eko Bagus Syafarudin', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198512132007011002', 'Brian Pujianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198602092007101001', 'Hidayat Subkhani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198602202006021003', 'Harry Susanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198603232007101002', 'Marthin Pitua', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198603282007101001', 'Mukhammad Nurhudah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198604102007101001', 'Muhammad Fajar Shidiq', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198604112006021004', 'Frengky Sahat Binsar Sitompul', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198606082014021006', 'Jaka Iswan Fitriawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198608052006021002', 'Sigit Satriyo Wibowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198610082006021004', 'Ade Yudha Utama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198612042010121005', 'Tjipto Aji Sudarso', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198703122007101002', 'Bekti Sulistianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198705112007101001', 'Ari Purnomo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198709182007101001', 'Mhd. Fuad Salim Nasution', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198710082015021001', 'Deni Rio Fandra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198712122008121002', 'Nicholas Lumban Tobing', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198803242015021003', 'Rendhy Koescahyadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198805262007101003', 'Zein Husein Siregar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198810252009121004', 'Muhammad Galih Permadi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198901172010011004', 'Rizki Ferdian Syah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198902162010011002', 'Dwicky Sofyan Hardhini', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198910052014021006', 'Maulana Hariyudha', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198910272012101001', 'M. Akmal', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('198912232010011002', 'Putra Aji Nugroho', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199001192010011001', 'Moch. Zainudin Efendi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199004242014022002', 'Lara Sri Yeni', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199005142010011001', 'Moh. Fajri Pratama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199005252010011002', 'Manik Semesta', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199008242012101002', 'Agristanda Surya Kusuma', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199009232014022002', 'Triana Putrie Vinansari', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199010032012102002', 'Reza Konesti', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199010272010011001', 'Eftrian Dika', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199012152010121001', 'Zulfiqi Fauzi Wibowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199101272012101002', 'Ahmad Dhabith Jihaduddin Shabran', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199103312010121003', 'Martiga Dwi Ananto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199105182012101002', 'Pademak Siringo Ringo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199105272009121001', 'Istian Prabowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199106012012101001', 'Benni Daniel Sihite', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199107232010121002', 'Allamaski Mochammad', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199108242012101001', 'Enda Putra Nanda Sembiring', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199109172012101002', 'Irfan Shahbhana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199109182014021001', 'Cahyo Baskoro Indra Maulana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199109212010121002', 'Robby Maulana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199111062012101001', 'Puspita Adhi Nugraha', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199111102013102001', 'Rossy Ananda', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199111142013101002', 'Dio Wijayanto Nugroho', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199111272013101001', 'Wisnu Wahyu Wardhana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199111282012101001', 'Teuku Muhammad Ridha', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199112112012101001', 'R. Kautsar Firdausi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199202072012101001', 'Febri Kurniawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199202072013101001', 'Dahlan Pamuji', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199202242013101003', 'Muchamad Dwi Susilo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199203202012101001', 'Rifan Primardian Ramadhani', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199204212013101002', 'Anselmus Kartino Medja', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199204222012101001', 'Ardhiansyah Fuad Asrurrosyid', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199205022012101001', 'Gemilang Bagas Putra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199206072012101003', 'Olverio Marshal Arief', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199208082013101001', 'Christy Agustian Situmorang', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199208132013101001', 'Fahmi Fahrur Rozi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199209062012101002', 'Rian Ajiwijaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199209292013101002', 'Rian Wijaya Sirait', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199211022013101001', 'Nopia Setia Putra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199211102013101002', 'Aris Aditya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199211152012101001', 'Jaya Nopiantho Sinulingga', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199212152018011006', 'Dezky Muji Setyo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199301102013101004', 'Edi Saputro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199304122013101002', 'Ahmad Fauzi Basuki', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199304262013101001', 'Anhar Prasetya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199307142013101001', 'Wahyu Pahlawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199308072013101001', 'Syaefullah Nur Ahmad P', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199309022013101002', 'Hedi Maulana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199309212018012004', 'Ayu Demmy Karinta', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199309272015021002', 'Wage Brantiawan Vaksiandra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199310042015021001', 'Arpan Hidayat', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199310282013101001', 'Ulwan Zaki', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199311062013101003', 'Novian Adi Nugroho', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199401052013101001', 'Fransisko', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199401252015021001', 'Indra Mangapon Purba', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199401262015021001', 'Hafiz Ali', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199402102015021001', 'Chairul Rizka Harahap', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199404242013101001', 'Pri Adiyanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199404252015021004', 'Ngabdul Khakim', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199405232015021003', 'Nara Praba Wardaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199405242015021003', 'Irwan Ardiansyah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199408252015021002', 'Fazlur Rahman', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199410022015021002', 'Rizky Eka Pradana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199412162015021002', 'Febri Ramadhoni Saputra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199412312015021005', 'Yusri Muhammad', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199501052015021002', 'Eko David Prasetiyo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199501232015021002', 'Welly Tamarodo Gultom', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199502042016121002', 'Muhammad Ardy Febriant Widagdo Atmaja', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199502172015021007', 'Hilman Nuzzy Al Mustaqim', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199502252015021001', 'Hanif Amirul Hakim', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199503212016121001', 'Aldo Monang Simanjuntak', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199503232015021003', 'Herbet Alfrin Simanjuntak', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199503252015021002', 'Agung Tri Wibowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199504072015021002', 'Herlambang Suko Prayogi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199504142015021001', 'A. Basofi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199504242015021007', 'Hafiizh Ha Razzaag', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199505102015021005', 'Agung Muhamad Reza', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199506032015021001', 'Muhammad Alif Amidhan Ganda Wijaya', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199507022015021001', 'Suryadin Sanusi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199507182015121001', 'Muhammad Sahal Savana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199508072015021005', 'Rivaldi Yudistira Bratanegara', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199508112015021003', 'Panji Erindra Hanggoro Raras', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199509052015021001', 'Muhammad Alfath Wijayanto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199509182015021003', 'Ridho Moch. Zain', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199509222018011003', 'Shofwan Zuhdi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199510102015121002', 'Ahmad Fiqri', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199510212015021001', 'Dwi Artha Oky Setyawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199510212015121001', 'Muhammad Reza', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199510282015021004', 'Ibnu Sholeh Nurul Firdaus', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199511042018011002', 'Friski Alexander Siburian', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199511122015021002', 'Brian Prasetya Kurniawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199511132015021002', 'Raden Muhammad Raka Aditama Poernama', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199512132015021001', 'Natanael Ignasio Ronoko', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199512202015021001', 'Rully Achmad Upoyo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199601092015121002', 'Arifan Anwar Pandia', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199601172015021001', 'I Gede Adika Prameswara Putra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199601172015021002', 'Galih Bekti Prabowo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199602182015121003', 'Rian Surya Angga Permana', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199602182015121004', 'Ardhi Febri Ferdyan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199602232015021001', 'Muhammad Hafizh Ezyoni', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199602282015021001', 'Putut Nur Alfianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199603162015121002', 'Yesack Adiyones Ngili', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199603262015021002', 'Muhammad Reza Aka Putra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199604082015121002', 'Muhammad Dzulham Fadhil', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199604242015021001', 'Faisal Ahmad Ilhami', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199604262018011004', 'Faisal Bustanul Arifin', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199605312015121002', 'Dody  Krisna Rianto', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199606132015121002', 'Abu Hanifa Al Ubaydah', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199606202018121002', 'Rifan Dwi Darmawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199606232018011004', 'Juan Jeremia Hasudungan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199607292015021001', 'Hashry Rizaddin Ahmad', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199609082015121001', 'Forester Sianipar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199610192015121001', 'Muhammad Aldy Mahron Nasution', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199611152018121001', 'Mochammad Falah Akbar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199702102016121002', 'Muhammad Fahrul Al Qadri', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199703142015121001', 'Henok Ginda Morgan Simarmata', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199703142016121001', 'Muhammad Fahmi Pratama Putra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199707042015121001', 'Ayyubi Atmara', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199708122018121001', 'Douglas Bungaran', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199709172019121001', 'Deandro Mikail Sebriano', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199709182018011002', 'Bagas Fathoni Mirhan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199801062018011001', 'Muhamad Umar Sena', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199803142018011001', 'Weby Ulul Albab', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199807192018011001', 'Achmad Dzulfikar Maulidi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199809022018011001', 'Yehezkiel Christian Prasetyo', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199809182018011001', 'Sihar Alberd', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199809252019121001', 'Muhammad Farhan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199810282018011001', 'Irvan Setyo Pratama Sunarko Putro', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199901012018121002', 'Adhitya Gabe Butar Butar', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199901152018011002', 'Muhammad Fakhri Ramadhandy', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199901252018011001', 'Radheva Hafizh Kurniawan', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('199905072018121005', 'Yusri Mahendra', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('200001222018121001', 'Billfrit Gregerius Situmorang', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('200006082019121001', 'Muhammad Farhan Habibi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true),
  ('200006202019121002', 'Muhammad Arsyvivaldi', (select id from public.units where code = 'FUNGSIONAL'), 'functional', null, true, true)
on conflict (nip) do update set
  name = excluded.name,
  unit_id = excluded.unit_id,
  position_role = excluded.position_role,
  is_active = true,
  deleted_at = null;

commit;
