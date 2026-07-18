begin;

-- Identitas autentikasi bersama. ID akun pegawai dipertahankan sama dengan
-- users.id agar seluruh audit dan referensi historis tetap valid.
create table if not exists public.accounts (
  id uuid primary key default gen_random_uuid(),
  account_type text not null check (account_type in ('user', 'admin')),
  created_at timestamptz not null default now()
);

insert into public.accounts (id, account_type)
select id, 'user' from public.users
on conflict (id) do nothing;

create or replace function public.pantas_ensure_user_account()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
begin
  insert into public.accounts (id, account_type)
  values (new.id, 'user')
  on conflict (id) do nothing;
  return new;
end;
$$;

drop trigger if exists users_ensure_account on public.users;
create trigger users_ensure_account
before insert on public.users
for each row execute function public.pantas_ensure_user_account();

create table if not exists public.admin_accounts (
  account_id uuid primary key references public.accounts(id) on delete cascade,
  username text not null unique
    check (username = lower(username) and username ~ '^[a-z][a-z0-9._-]{2,63}$'),
  name text not null check (char_length(btrim(name)) between 2 and 200),
  password_hash text not null,
  must_change_password boolean not null default true,
  is_active boolean not null default true,
  last_login_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

drop trigger if exists admin_accounts_updated_at on public.admin_accounts;
create trigger admin_accounts_updated_at
before update on public.admin_accounts
for each row execute function public.pantas_set_updated_at();

-- Kolom-kolom berikut menyimpan pelaku tindakan. Referensinya dipindahkan dari
-- profil pegawai ke identitas akun sehingga dapat menunjuk admin non-pegawai.
alter table public.sessions drop constraint if exists sessions_user_id_fkey;
alter table public.sessions
  add constraint sessions_user_id_fkey foreign key (user_id)
  references public.accounts(id) on delete cascade;

alter table public.user_assignment_history drop constraint if exists user_assignment_history_changed_by_fkey;
alter table public.user_assignment_history
  add constraint user_assignment_history_changed_by_fkey foreign key (changed_by)
  references public.accounts(id) on delete set null;

alter table public.import_batches drop constraint if exists import_batches_created_by_fkey;
alter table public.import_batches
  add constraint import_batches_created_by_fkey foreign key (created_by)
  references public.accounts(id) on delete restrict;
alter table public.import_batches drop constraint if exists import_batches_published_by_fkey;
alter table public.import_batches
  add constraint import_batches_published_by_fkey foreign key (published_by)
  references public.accounts(id) on delete restrict;

alter table public.appeal_items drop constraint if exists appeal_items_admin_by_fkey;
alter table public.appeal_items
  add constraint appeal_items_admin_by_fkey foreign key (admin_by)
  references public.accounts(id) on delete restrict;

alter table public.parameters drop constraint if exists parameters_updated_by_fkey;
alter table public.parameters
  add constraint parameters_updated_by_fkey foreign key (updated_by)
  references public.accounts(id) on delete set null;

alter table public.notifications drop constraint if exists notifications_user_id_fkey;
alter table public.notifications
  add constraint notifications_user_id_fkey foreign key (user_id)
  references public.accounts(id) on delete cascade;

alter table public.notification_jobs drop constraint if exists notification_jobs_user_id_fkey;
alter table public.notification_jobs
  add constraint notification_jobs_user_id_fkey foreign key (user_id)
  references public.accounts(id) on delete cascade;

alter table public.audit_logs drop constraint if exists audit_logs_actor_id_fkey;
alter table public.audit_logs
  add constraint audit_logs_actor_id_fkey foreign key (actor_id)
  references public.accounts(id) on delete set null;

-- Flag lama dipertahankan hanya untuk kompatibilitas data, tetapi tidak lagi
-- memberikan hak administrator dan seluruh nilainya dinonaktifkan.
update public.users set is_admin = false where is_admin;

drop function if exists public.pantas_bootstrap_admin(text, text);
create or replace function public.pantas_bootstrap_admin(
  p_username text,
  p_name text,
  p_initial_password text
)
returns uuid
language plpgsql
security definer
set search_path = public
as $$
declare
  v_username text := lower(btrim(p_username));
  v_account_id uuid;
begin
  if v_username !~ '^[a-z][a-z0-9._-]{2,63}$' then
    raise exception 'Username admin harus 3-64 karakter dan diawali huruf';
  end if;
  if char_length(p_initial_password) < 12 or char_length(p_initial_password) > 128 then
    raise exception 'Password awal admin harus 12-128 karakter';
  end if;
  if char_length(btrim(p_name)) < 2 then
    raise exception 'Nama admin tidak valid';
  end if;

  select account_id into v_account_id
  from public.admin_accounts
  where username = v_username;

  if v_account_id is null then
    v_account_id := gen_random_uuid();
    insert into public.accounts (id, account_type) values (v_account_id, 'admin');
    insert into public.admin_accounts (
      account_id, username, name, password_hash, must_change_password, is_active
    ) values (
      v_account_id, v_username, btrim(p_name),
      crypt(p_initial_password, gen_salt('bf', 12)), true, true
    );
  else
    update public.admin_accounts
    set name = btrim(p_name), is_active = true
    where account_id = v_account_id;
  end if;

  return v_account_id;
end;
$$;

revoke all on function public.pantas_bootstrap_admin(text, text, text)
from public, anon, authenticated;
revoke all on function public.pantas_ensure_user_account()
from public, anon, authenticated;

alter table public.accounts enable row level security;
alter table public.admin_accounts enable row level security;
revoke all on table public.accounts from anon, authenticated;
revoke all on table public.admin_accounts from anon, authenticated;

commit;
