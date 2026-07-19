-- PANTAS V6 - integritas struktur organisasi yang dapat dikelola admin.
-- Aman dijalankan pada database V5 yang sudah berisi unit dan pegawai.

begin;

create index if not exists units_parent_sort_idx
  on public.units (parent_id, sort_order, name);

create unique index if not exists units_one_office
  on public.units (unit_type)
  where unit_type = 'office';

create or replace function public.pantas_validate_unit_hierarchy()
returns trigger
language plpgsql
set search_path = public, pg_temp
as $$
declare
  v_parent_type text;
  v_parent_active boolean;
  v_active_members bigint;
  v_active_children bigint;
begin
  new.code := upper(btrim(new.code));
  new.name := btrim(new.name);
  new.source_name := nullif(btrim(new.source_name), '');

  if new.code !~ '^[A-Z0-9][A-Z0-9._-]{1,31}$' then
    raise exception using
      errcode = '23514',
      message = 'Kode unit harus 2-32 karakter: huruf kapital, angka, titik, garis bawah, atau tanda hubung.';
  end if;
  if char_length(new.name) < 2 or char_length(new.name) > 200 then
    raise exception using errcode = '23514', message = 'Nama unit harus terdiri dari 2-200 karakter.';
  end if;
  if char_length(coalesce(new.source_name, '')) > 300 then
    raise exception using errcode = '23514', message = 'Nama penempatan Excel maksimal 300 karakter.';
  end if;
  if new.parent_id = new.id then
    raise exception using errcode = '23514', message = 'Unit tidak dapat menjadi induk bagi dirinya sendiri.';
  end if;
  if tg_op = 'UPDATE' and new.unit_type <> old.unit_type then
    raise exception using errcode = '23514', message = 'Jenis unit tidak dapat diubah setelah unit dibuat.';
  end if;

  if new.unit_type = 'office' then
    if new.parent_id is not null then
      raise exception using errcode = '23514', message = 'Unit kantor tidak boleh memiliki induk.';
    end if;
  else
    if new.parent_id is null then
      raise exception using errcode = '23514', message = 'Unit selain kantor wajib memiliki unit induk.';
    end if;

    select unit_type, is_active
      into v_parent_type, v_parent_active
    from public.units
    where id = new.parent_id;

    if not found then
      raise exception using errcode = '23503', message = 'Unit induk tidak ditemukan.';
    end if;
    if new.unit_type = 'division' and v_parent_type <> 'office' then
      raise exception using errcode = '23514', message = 'Bidang/bagian harus berada di bawah unit kantor.';
    end if;
    if new.unit_type = 'section' and v_parent_type <> 'division' then
      raise exception using errcode = '23514', message = 'Seksi/subbagian harus berada di bawah bidang/bagian.';
    end if;
    if new.unit_type = 'functional' and v_parent_type <> 'office' then
      raise exception using errcode = '23514', message = 'Unit Fungsional harus berada di bawah unit kantor.';
    end if;
    if new.is_active and not v_parent_active then
      raise exception using errcode = '23514', message = 'Unit aktif tidak boleh berada di bawah induk nonaktif.';
    end if;
  end if;

  if tg_op = 'UPDATE' and old.is_active and not new.is_active then
    select count(*) into v_active_members
    from public.users
    where unit_id = new.id and is_active and deleted_at is null;

    select count(*) into v_active_children
    from public.units
    where parent_id = new.id and is_active;

    if v_active_members > 0 then
      raise exception using
        errcode = '23514',
        message = 'Unit masih memiliki pegawai aktif dan belum dapat dinonaktifkan.';
    end if;
    if v_active_children > 0 then
      raise exception using
        errcode = '23514',
        message = 'Unit masih memiliki unit anak aktif dan belum dapat dinonaktifkan.';
    end if;
  end if;

  return new;
end;
$$;

do $$
begin
  if exists (
    select 1
    from public.units u
    left join public.units p on p.id = u.parent_id
    where (u.unit_type = 'office' and u.parent_id is not null)
       or (u.unit_type = 'division' and (p.id is null or p.unit_type <> 'office'))
       or (u.unit_type = 'section' and (p.id is null or p.unit_type <> 'division'))
       or (u.unit_type = 'functional' and (p.id is null or p.unit_type <> 'office'))
       or (u.is_active and p.id is not null and not p.is_active)
       or (not u.is_active and exists (select 1 from public.units c where c.parent_id = u.id and c.is_active))
       or (not u.is_active and exists (select 1 from public.users x where x.unit_id = u.id and x.is_active and x.deleted_at is null))
       or u.code !~ '^[A-Z0-9][A-Z0-9._-]{1,31}$'
       or char_length(btrim(u.name)) not between 2 and 200
       or char_length(coalesce(btrim(u.source_name), '')) > 300
  ) then
    raise exception 'Migration V6 dibatalkan: struktur unit lama tidak valid. Perbaiki parent, tipe, kode, nama, atau status unit terlebih dahulu.';
  end if;
end;
$$;

drop trigger if exists units_validate_hierarchy on public.units;
create trigger units_validate_hierarchy
before insert or update on public.units
for each row execute function public.pantas_validate_unit_hierarchy();

create or replace function public.pantas_protect_system_unit_delete()
returns trigger
language plpgsql
set search_path = public, pg_temp
as $$
begin
  if old.unit_type in ('office', 'functional') then
    raise exception using
      errcode = '23514',
      message = 'Unit kantor dan Fungsional merupakan unit sistem dan tidak dapat dihapus.';
  end if;
  return old;
end;
$$;

drop trigger if exists units_protect_system_delete on public.units;
create trigger units_protect_system_delete
before delete on public.units
for each row execute function public.pantas_protect_system_unit_delete();

revoke all on function public.pantas_validate_unit_hierarchy() from public, anon, authenticated;
revoke all on function public.pantas_protect_system_unit_delete() from public, anon, authenticated;

commit;
