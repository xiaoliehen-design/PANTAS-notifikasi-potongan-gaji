-- PANTAS — Pemantauan Absensi dan Tunjangan Secara Akumulatif
-- Jalankan satu kali melalui Supabase SQL Editor pada project baru.

begin;

create schema if not exists extensions;
create extension if not exists pgcrypto with schema extensions;

create or replace function public.pantas_set_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

create table if not exists public.units (
  id uuid primary key default gen_random_uuid(),
  code text not null unique,
  name text not null,
  source_name text unique,
  unit_type text not null check (unit_type in ('office', 'division', 'section', 'functional')),
  parent_id uuid references public.units(id) on delete restrict,
  sort_order integer not null default 0,
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists public.users (
  id uuid primary key default gen_random_uuid(),
  nip text not null unique check (nip ~ '^[0-9]{18}$'),
  name text not null,
  unit_id uuid not null references public.units(id) on delete restrict,
  position_role text not null default 'staff'
    check (position_role in ('staff', 'section_head', 'division_head', 'office_head', 'functional')),
  is_admin boolean not null default false,
  email text,
  email_verified_at timestamptz,
  phone_e164 text check (phone_e164 is null or phone_e164 ~ '^\+[1-9][0-9]{7,14}$'),
  phone_verified_at timestamptz,
  password_hash text,
  must_change_password boolean not null default true,
  is_active boolean not null default true,
  last_login_at timestamptz,
  deleted_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint users_email_normalized check (email is null or email = lower(btrim(email)))
);

create unique index if not exists users_email_unique
  on public.users (lower(email)) where email is not null and deleted_at is null;
create unique index if not exists users_phone_unique
  on public.users (phone_e164) where phone_e164 is not null and deleted_at is null;
create index if not exists users_unit_active_idx on public.users (unit_id, is_active);
create index if not exists users_role_active_idx on public.users (position_role, is_active);
create unique index if not exists users_one_section_head_per_unit
  on public.users (unit_id) where position_role = 'section_head' and is_active and deleted_at is null;
create unique index if not exists users_one_division_head_per_unit
  on public.users (unit_id) where position_role = 'division_head' and is_active and deleted_at is null;
create unique index if not exists users_one_office_head
  on public.users (position_role) where position_role = 'office_head' and is_active and deleted_at is null;

create table if not exists public.user_assignment_history (
  id bigserial primary key,
  user_id uuid not null references public.users(id) on delete restrict,
  previous_unit_id uuid references public.units(id) on delete restrict,
  new_unit_id uuid not null references public.units(id) on delete restrict,
  previous_role text,
  new_role text not null,
  changed_by uuid references public.users(id) on delete set null,
  reason text,
  changed_at timestamptz not null default now()
);

create table if not exists public.sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.users(id) on delete cascade,
  token_hash bytea not null unique,
  csrf_hash bytea not null,
  ip_address inet,
  user_agent text,
  created_at timestamptz not null default now(),
  last_seen_at timestamptz not null default now(),
  expires_at timestamptz not null,
  revoked_at timestamptz
);

create index if not exists sessions_user_active_idx on public.sessions (user_id, expires_at)
  where revoked_at is null;

create table if not exists public.login_attempts (
  id bigserial primary key,
  nip text,
  ip_address inet,
  was_successful boolean not null,
  occurred_at timestamptz not null default now()
);

create index if not exists login_attempts_nip_time_idx on public.login_attempts (nip, occurred_at desc);
create index if not exists login_attempts_ip_time_idx on public.login_attempts (ip_address, occurred_at desc);

create table if not exists public.recovery_otps (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.users(id) on delete cascade,
  purpose text not null check (purpose in ('password_reset', 'verify_email', 'verify_phone')),
  channel text not null check (channel in ('email', 'phone')),
  destination text not null,
  otp_hash bytea not null,
  requested_ip inet,
  attempts integer not null default 0 check (attempts >= 0),
  expires_at timestamptz not null,
  consumed_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists recovery_otps_lookup_idx
  on public.recovery_otps (user_id, purpose, channel, created_at desc);

create table if not exists public.pending_contact_changes (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.users(id) on delete cascade,
  channel text not null check (channel in ('email', 'phone')),
  destination text not null,
  otp_hash bytea not null,
  attempts integer not null default 0,
  expires_at timestamptz not null,
  consumed_at timestamptz,
  created_at timestamptz not null default now()
);

create table if not exists public.reporting_periods (
  id uuid primary key default gen_random_uuid(),
  label text not null,
  period_start date not null,
  period_end date not null,
  published_batch_id uuid,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (period_start, period_end),
  check (period_end >= period_start),
  check (period_end - period_start <= 62)
);

create table if not exists public.import_batches (
  id uuid primary key default gen_random_uuid(),
  period_id uuid not null references public.reporting_periods(id) on delete restrict,
  version integer not null,
  original_filename text not null,
  file_sha256 text not null check (file_sha256 ~ '^[a-f0-9]{64}$'),
  file_size_bytes bigint not null check (file_size_bytes > 0),
  sheet_name text not null,
  integrity_status text not null default 'valid'
    check (integrity_status in ('valid', 'recovered_partial_container')),
  status text not null default 'draft'
    check (status in ('draft', 'published', 'superseded', 'rejected')),
  row_count integer not null default 0,
  employee_count integer not null default 0,
  deduction_day_count integer not null default 0,
  total_deduction_rate numeric(12,6) not null default 0,
  warning_summary jsonb not null default '{}'::jsonb,
  created_by uuid not null references public.users(id) on delete restrict,
  published_by uuid references public.users(id) on delete restrict,
  created_at timestamptz not null default now(),
  published_at timestamptz,
  unique (period_id, version)
);

alter table public.reporting_periods
  drop constraint if exists reporting_periods_published_batch_id_fkey;
alter table public.reporting_periods
  add constraint reporting_periods_published_batch_id_fkey
  foreign key (published_batch_id) references public.import_batches(id) on delete set null;

create unique index if not exists import_batches_published_hash_unique
  on public.import_batches (file_sha256) where status = 'published';
create index if not exists import_batches_period_status_idx on public.import_batches (period_id, status, version desc);

create table if not exists public.deduction_rules (
  id uuid primary key default gen_random_uuid(),
  source_field text not null check (source_field in ('late', 'early_leave', 'leave', 'status', 'shift')),
  code text not null,
  label text not null,
  rate numeric(8,6) not null check (rate >= 0 and rate <= 1),
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (source_field, code)
);

create table if not exists public.attendance_records (
  id bigserial primary key,
  batch_id uuid not null references public.import_batches(id) on delete cascade,
  period_id uuid not null references public.reporting_periods(id) on delete restrict,
  user_id uuid not null references public.users(id) on delete restrict,
  source_row integer not null check (source_row >= 5),
  work_date date not null,
  check_in time,
  check_out time,
  late_code text,
  early_leave_code text,
  shift_code text,
  attendance_status text,
  leave_type text,
  assignment_type text,
  source_confirmation text,
  notes text,
  source_division text,
  source_placement text,
  deduction_rate numeric(8,6) not null default 0 check (deduction_rate >= 0 and deduction_rate <= 1),
  deduction_components jsonb not null default '[]'::jsonb,
  created_at timestamptz not null default now(),
  unique (batch_id, user_id, work_date)
);

create index if not exists attendance_batch_user_date_idx
  on public.attendance_records (batch_id, user_id, work_date);
create index if not exists attendance_period_user_idx
  on public.attendance_records (period_id, user_id);
create index if not exists attendance_user_deduction_idx
  on public.attendance_records (user_id, period_id, deduction_rate) where deduction_rate > 0;

create table if not exists public.appeal_reason_categories (
  id uuid primary key default gen_random_uuid(),
  code text not null unique,
  label text not null,
  description text,
  is_active boolean not null default true,
  sort_order integer not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists public.appeals (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.users(id) on delete restrict,
  period_id uuid not null references public.reporting_periods(id) on delete restrict,
  status text not null default 'submitted'
    check (status in ('submitted', 'supervisor_review', 'admin_review', 'finalized', 'cancelled')),
  submitted_at timestamptz not null default now(),
  finalized_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, period_id)
);

create table if not exists public.appeal_items (
  id uuid primary key default gen_random_uuid(),
  appeal_id uuid not null references public.appeals(id) on delete cascade,
  attendance_record_id bigint not null references public.attendance_records(id) on delete restrict,
  reason_category_id uuid not null references public.appeal_reason_categories(id) on delete restrict,
  explanation text not null check (char_length(btrim(explanation)) between 10 and 3000),
  original_deduction_rate numeric(8,6) not null check (original_deduction_rate > 0 and original_deduction_rate <= 1),
  supervisor_status text not null default 'pending'
    check (supervisor_status in ('pending', 'accepted', 'rejected')),
  supervisor_by uuid references public.users(id) on delete restrict,
  supervisor_comment text check (supervisor_comment is null or char_length(supervisor_comment) <= 2000),
  supervisor_reviewed_at timestamptz,
  admin_status text not null default 'pending'
    check (admin_status in ('pending', 'approved', 'rejected')),
  admin_by uuid references public.users(id) on delete restrict,
  admin_comment text check (admin_comment is null or char_length(admin_comment) <= 2000),
  adjusted_deduction_rate numeric(8,6) check (adjusted_deduction_rate is null or (adjusted_deduction_rate >= 0 and adjusted_deduction_rate <= 1)),
  admin_reviewed_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (attendance_record_id)
);

create index if not exists appeal_items_supervisor_queue_idx
  on public.appeal_items (supervisor_status, created_at) where supervisor_status = 'pending';
create index if not exists appeal_items_admin_queue_idx
  on public.appeal_items (admin_status, created_at) where supervisor_status <> 'pending' and admin_status = 'pending';

create table if not exists public.appeal_documents (
  id uuid primary key default gen_random_uuid(),
  appeal_item_id uuid not null references public.appeal_items(id) on delete cascade,
  storage_path text not null unique,
  original_filename text not null,
  mime_type text not null check (mime_type in ('application/pdf', 'image/jpeg', 'image/png')),
  size_bytes bigint not null check (size_bytes > 0 and size_bytes <= 5242880),
  sha256 text not null check (sha256 ~ '^[a-f0-9]{64}$'),
  uploaded_by uuid not null references public.users(id) on delete restrict,
  created_at timestamptz not null default now()
);

create table if not exists public.parameters (
  key text primary key,
  category text not null,
  label text not null,
  description text,
  value_json jsonb not null,
  value_type text not null check (value_type in ('integer', 'decimal', 'percent', 'json', 'boolean')),
  updated_by uuid references public.users(id) on delete set null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists public.notifications (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.users(id) on delete cascade,
  kind text not null,
  title text not null,
  body text not null,
  action_url text,
  read_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists notifications_user_unread_idx
  on public.notifications (user_id, created_at desc) where read_at is null;

create table if not exists public.notification_jobs (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references public.users(id) on delete cascade,
  channel text not null check (channel in ('email', 'phone')),
  destination text not null,
  template_code text not null,
  payload jsonb not null default '{}'::jsonb,
  status text not null default 'pending' check (status in ('pending', 'processing', 'sent', 'failed', 'cancelled')),
  attempts integer not null default 0,
  next_attempt_at timestamptz not null default now(),
  locked_at timestamptz,
  last_error text,
  sent_at timestamptz,
  created_at timestamptz not null default now()
);

create index if not exists notification_jobs_worker_idx
  on public.notification_jobs (status, next_attempt_at, created_at)
  where status in ('pending', 'processing');

create table if not exists public.audit_logs (
  id bigserial primary key,
  actor_id uuid references public.users(id) on delete set null,
  action text not null,
  entity_type text not null,
  entity_id text,
  metadata jsonb not null default '{}'::jsonb,
  ip_address inet,
  occurred_at timestamptz not null default now()
);

create index if not exists audit_logs_actor_time_idx on public.audit_logs (actor_id, occurred_at desc);
create index if not exists audit_logs_entity_idx on public.audit_logs (entity_type, entity_id, occurred_at desc);

drop trigger if exists units_updated_at on public.units;
create trigger units_updated_at before update on public.units
for each row execute function public.pantas_set_updated_at();
drop trigger if exists users_updated_at on public.users;
create trigger users_updated_at before update on public.users
for each row execute function public.pantas_set_updated_at();
drop trigger if exists periods_updated_at on public.reporting_periods;
create trigger periods_updated_at before update on public.reporting_periods
for each row execute function public.pantas_set_updated_at();
drop trigger if exists rules_updated_at on public.deduction_rules;
create trigger rules_updated_at before update on public.deduction_rules
for each row execute function public.pantas_set_updated_at();
drop trigger if exists reasons_updated_at on public.appeal_reason_categories;
create trigger reasons_updated_at before update on public.appeal_reason_categories
for each row execute function public.pantas_set_updated_at();
drop trigger if exists appeals_updated_at on public.appeals;
create trigger appeals_updated_at before update on public.appeals
for each row execute function public.pantas_set_updated_at();
drop trigger if exists appeal_items_updated_at on public.appeal_items;
create trigger appeal_items_updated_at before update on public.appeal_items
for each row execute function public.pantas_set_updated_at();
drop trigger if exists parameters_updated_at on public.parameters;
create trigger parameters_updated_at before update on public.parameters
for each row execute function public.pantas_set_updated_at();

insert into public.units (code, name, source_name, unit_type, sort_order)
values ('KPU-TPK', 'KPU Bea dan Cukai Tipe A Tanjung Priok', '-', 'office', 0)
on conflict (code) do update set name = excluded.name, source_name = excluded.source_name;

insert into public.deduction_rules (source_field, code, label, rate, sort_order)
values
  ('late', 'TL1', 'Terlambat tingkat 1', 0.010000, 10),
  ('late', 'TL2', 'Terlambat tingkat 2', 0.012500, 20),
  ('late', 'TL3', 'Terlambat tingkat 3', 0.025000, 30),
  ('late', 'LA', 'Tidak melakukan presensi masuk', 0.025000, 40),
  ('early_leave', 'PSW1', 'Pulang sebelum waktunya tingkat 1', 0.005000, 50),
  ('early_leave', 'PSW2', 'Pulang sebelum waktunya tingkat 2', 0.010000, 60),
  ('early_leave', 'PSW3', 'Pulang sebelum waktunya tingkat 3', 0.012500, 70),
  ('early_leave', 'PSW4', 'Pulang sebelum waktunya tingkat 4', 0.025000, 80),
  ('early_leave', 'LA', 'Tidak melakukan presensi pulang', 0.025000, 90),
  ('leave', 'Cuti Alasan Penting Dipotong', 'Cuti alasan penting dengan potongan', 0.050000, 100),
  ('leave', 'Cuti Besar Dipotong', 'Cuti besar dengan potongan', 0.025000, 110),
  ('leave', 'Cuti Sakit Dipotong', 'Cuti sakit dengan potongan', 0.025000, 120),
  ('status', 'I', 'Izin tidak masuk', 0.050000, 130),
  ('status', 'TK', 'Tanpa keterangan', 0.050000, 140)
on conflict (source_field, code) do update
set label = excluded.label, rate = excluded.rate, sort_order = excluded.sort_order;

insert into public.appeal_reason_categories (code, label, description, sort_order)
values
  ('personal_negligence', 'Kelalaian pribadi', 'Termasuk kendala atau kesalahan perangkat pribadi.', 10),
  ('assignment_letter', 'Surat tugas', 'Ketidaksesuaian terkait pelaksanaan surat tugas.', 20),
  ('attendance_system', 'Sistem presensi error', 'Gangguan pada sistem atau layanan presensi.', 30),
  ('force_majeure', 'Keadaan kahar', 'Keadaan di luar kendali yang dapat dibuktikan.', 40),
  ('leave_off', 'Cuti/off', 'Hari seharusnya tercatat sebagai cuti atau off.', 50),
  ('other', 'Lainnya', 'Alasan lain yang relevan dan dapat dijelaskan.', 60)
on conflict (code) do update
set label = excluded.label, description = excluded.description, sort_order = excluded.sort_order;

insert into public.parameters (key, category, label, description, value_json, value_type)
values
  ('individual_anomaly_lookback_months', 'warning', 'Periode anomali individu', 'Jumlah periode terdahulu untuk memastikan pegawai biasanya tidak memiliki potongan.', '6', 'integer'),
  ('individual_anomaly_prior_max_rate', 'warning', 'Batas riwayat anomali individu', 'Total potongan maksimum pada periode acuan agar potongan periode berjalan dinilai anomali.', '0', 'percent'),
  ('bad_habit_consecutive_periods', 'warning', 'Periode kebiasaan buruk', 'Jumlah periode berturut-turut dengan potongan.', '3', 'integer'),
  ('aggregate_spike_lookback_months', 'warning', 'Periode acuan lonjakan unit', 'Jumlah periode sebelumnya untuk menghitung rata-rata unit.', '6', 'integer'),
  ('aggregate_spike_multiplier', 'warning', 'Pengali lonjakan unit', 'Potongan periode berjalan dibanding rata-rata periode acuan.', '2.0', 'decimal'),
  ('aggregate_spike_min_delta', 'warning', 'Selisih minimum lonjakan unit', 'Selisih minimum agar lonjakan kecil tidak ditandai.', '0.005', 'percent'),
  ('aggregate_average_threshold', 'warning', 'Batas rata-rata potongan unit', 'Peringatan jika rata-rata potongan pegawai pada unit melebihi nilai ini.', '0.005', 'percent')
on conflict (key) do update
set category = excluded.category, label = excluded.label, description = excluded.description,
    value_json = excluded.value_json, value_type = excluded.value_type;

create or replace view public.published_attendance
with (security_invoker = true)
as
select ar.*
from public.attendance_records ar
join public.reporting_periods rp on rp.id = ar.period_id
where rp.published_batch_id = ar.batch_id;

create or replace view public.effective_attendance
with (security_invoker = true)
as
select
  pa.*,
  case
    when ai.admin_status = 'approved' then coalesce(ai.adjusted_deduction_rate, 0)
    else pa.deduction_rate
  end as effective_deduction_rate,
  ai.id as appeal_item_id,
  ai.supervisor_status,
  ai.admin_status
from public.published_attendance pa
left join public.appeal_items ai on ai.attendance_record_id = pa.id;

create or replace view public.monthly_user_summary
with (security_invoker = true)
as
select
  ea.period_id,
  ea.user_id,
  sum(ea.deduction_rate)::numeric(12,6) as original_deduction_rate,
  sum(ea.effective_deduction_rate)::numeric(12,6) as effective_deduction_rate,
  count(*) filter (where ea.deduction_rate > 0) as deduction_days,
  count(*) filter (where ea.shift_code = 'OFF') as off_days,
  count(*) filter (where nullif(btrim(ea.leave_type), '') is not null) as leave_days,
  sum(case ea.shift_code when 'P' then 1 when 'M' then 1 when 'PM' then 2 else 0 end) as work_days,
  count(*) filter (where ea.shift_code in ('L1', 'L2')) as overtime_days
from public.effective_attendance ea
group by ea.period_id, ea.user_id;

create or replace function public.pantas_bootstrap_admin(p_nip text, p_name text)
returns uuid
language plpgsql
security definer
set search_path = public
as $$
declare
  v_unit_id uuid;
  v_user_id uuid;
begin
  if p_nip !~ '^[0-9]{18}$' then
    raise exception 'NIP harus berisi tepat 18 digit';
  end if;

  select id into v_unit_id from public.units where code = 'KPU-TPK';
  insert into public.users (nip, name, unit_id, position_role, is_admin, password_hash, must_change_password)
  values (p_nip, btrim(p_name), v_unit_id, 'staff', true, null, true)
  on conflict (nip) do update
  set name = excluded.name, is_admin = true, is_active = true, deleted_at = null
  returning id into v_user_id;
  return v_user_id;
end;
$$;

revoke all on function public.pantas_bootstrap_admin(text, text) from public, anon, authenticated;

do $$
declare
  table_name text;
begin
  foreach table_name in array array[
    'units','users','user_assignment_history','sessions','login_attempts','recovery_otps',
    'pending_contact_changes','reporting_periods','import_batches','deduction_rules',
    'attendance_records','appeal_reason_categories','appeals','appeal_items','appeal_documents',
    'parameters','notifications','notification_jobs','audit_logs'
  ]
  loop
    execute format('alter table public.%I enable row level security', table_name);
    execute format('revoke all on table public.%I from anon, authenticated', table_name);
  end loop;
end;
$$;

revoke all on public.published_attendance from anon, authenticated;
revoke all on public.effective_attendance from anon, authenticated;
revoke all on public.monthly_user_summary from anon, authenticated;

-- Bucket privat; akses file hanya melalui backend PANTAS dengan service-role key.
insert into storage.buckets (id, name, public, file_size_limit, allowed_mime_types)
values (
  'pantas-appeals',
  'pantas-appeals',
  false,
  5242880,
  array['application/pdf', 'image/jpeg', 'image/png']
)
on conflict (id) do update
set public = false,
    file_size_limit = excluded.file_size_limit,
    allowed_mime_types = excluded.allowed_mime_types;

commit;
