-- PANTAS — perbaikan bootstrap admin pada Supabase.
-- Aman dijalankan ulang. Migration ini tidak mengubah atau menghapus data pegawai.

begin;

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
      extensions.crypt(p_initial_password, extensions.gen_salt('bf', 12)), true, true
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

commit;
