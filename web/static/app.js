"use strict";

const state = {
  user: null,
  route: "dashboard",
  dashboard: null,
  units: null,
  sidebarCollapsed: readSidebarPreference(),
};

const $ = (selector, root = document) => root.querySelector(selector);
const $$ = (selector, root = document) => [...root.querySelectorAll(selector)];
const page = $("#page-content");

document.addEventListener("DOMContentLoaded", boot);

async function boot() {
  bindGlobalEvents();
  try {
    const result = await api("/api/auth/me", { quiet: true });
    state.user = result.user;
    showApp();
  } catch {
    showAuth();
  }
}

function bindGlobalEvents() {
  $("#login-form").addEventListener("submit", login);
  $("#forgot-form").addEventListener("submit", forgotPassword);
  $("#reset-form").addEventListener("submit", resetPassword);
  $("#show-forgot").addEventListener("click", () => setAuthView("forgot"));
  $$('[data-auth-back]').forEach(button => button.addEventListener("click", () => setAuthView("login")));
  $$(".password-toggle").forEach(button => button.addEventListener("click", () => {
    const input = button.closest(".password-wrap").querySelector("input");
    input.type = input.type === "password" ? "text" : "password";
  }));
  $("#logout-button").addEventListener("click", logout);
  $("#menu-button").addEventListener("click", () => $("#sidebar").classList.toggle("open"));
  $("#sidebar-collapse-button").addEventListener("click", toggleSidebar);
  $("#notification-button").addEventListener("click", openNotifications);
  $("#drawer-close").addEventListener("click", closeNotifications);
  $("#drawer-backdrop").addEventListener("click", closeNotifications);
  $("#profile-button").addEventListener("click", () => navigate("profile"));
  $("#password-form").addEventListener("submit", changePassword);
  window.addEventListener("hashchange", renderRoute);
  window.addEventListener("resize", () => {
    applySidebarState();
    if (state.dashboard) drawHistoryChart(state.dashboard.history || []);
  });
  document.addEventListener("click", async event => {
    const documentButton = event.target.closest("[data-documents]");
    if (documentButton) await revealDocuments(documentButton);
  });
}

async function login(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const button = $("button[type=submit]", form);
  setBusy(button, true, "Memeriksa…");
  try {
    const data = await api("/api/auth/login", { method: "POST", body: formJSON(form), auth: false });
    state.user = data.user;
    form.reset();
    showApp();
  } catch (error) {
    toast("Login gagal", error.message, "error");
  } finally {
    setBusy(button, false);
  }
}

async function forgotPassword(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const values = formJSON(form);
  const button = $("button[type=submit]", form);
  setBusy(button, true, "Mengirim…");
  try {
    const data = await api("/api/auth/forgot-password", { method: "POST", body: values, auth: false });
    const reset = $("#reset-form");
    reset.elements.nip.value = values.nip;
    reset.elements.channel.value = values.channel;
    setAuthView("reset");
    toast("Permintaan diterima", data.message, "success");
  } catch (error) {
    toast("Belum dapat mengirim kode", error.message, "error");
  } finally { setBusy(button, false); }
}

async function resetPassword(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const button = $("button[type=submit]", form);
  setBusy(button, true, "Menyimpan…");
  try {
    const data = await api("/api/auth/reset-password", { method: "POST", body: formJSON(form), auth: false });
    form.reset(); setAuthView("login");
    toast("Password diperbarui", data.message, "success");
  } catch (error) { toast("Gagal mengatur ulang", error.message, "error"); }
  finally { setBusy(button, false); }
}

async function changePassword(event) {
  event.preventDefault();
  const form = event.currentTarget;
  const button = $("button[type=submit]", form);
  setBusy(button, true, "Menyimpan…");
  try {
    const data = await api("/api/auth/change-password", { method: "POST", body: formJSON(form) });
    $("#password-dialog").close();
    state.user = null;
    showAuth();
    toast("Berhasil", data.message, "success");
  } catch (error) {
    toast("Password belum berubah", error.message, "error");
    if (error.code === "current_password_invalid") {
      form.elements.current_password.value = "";
      form.elements.current_password.setAttribute("aria-invalid", "true");
      form.elements.current_password.focus();
    }
  }
  finally { setBusy(button, false); }
}

async function logout() {
  try { await api("/api/auth/logout", { method: "POST" }); } catch {}
  state.user = null;
  showAuth();
}

function setAuthView(name) {
  ["login", "forgot", "reset"].forEach(view => $(`#${view}-view`).hidden = view !== name);
}

function showAuth() {
  $("#app-shell").hidden = true;
  $("#auth-shell").hidden = false;
  $("#sidebar").classList.remove("open");
  if ($("#password-dialog").open) $("#password-dialog").close();
  closeNotifications();
  history.replaceState(null, "", location.pathname);
  setAuthView("login");
}

function showApp() {
  $("#auth-shell").hidden = true;
  $("#app-shell").hidden = false;
  const initials = getInitials(state.user.name);
  $("#sidebar-avatar").textContent = initials;
  $("#top-avatar").textContent = initials;
  $("#sidebar-name").textContent = state.user.name;
  $("#sidebar-unit").textContent = state.user.unit_name;
  $("#top-name").textContent = state.user.name.split(" ")[0];
  applySidebarState();
  buildNavigation();
  if (state.user.must_change_password) {
    $("#password-dialog").showModal();
    return;
  }
  if (!location.hash) location.hash = state.user.is_admin ? "#monitoring" : "#dashboard";
  renderRoute();
}

function readSidebarPreference() {
  try { return localStorage.getItem("pantas_sidebar_collapsed") === "true"; }
  catch { return false; }
}

function applySidebarState() {
  const compact = state.sidebarCollapsed && window.matchMedia("(min-width: 861px)").matches;
  $("#app-shell").classList.toggle("sidebar-collapsed", compact);
  const button = $("#sidebar-collapse-button");
  button.setAttribute("aria-label", compact ? "Perbesar menu" : "Minimize menu");
  button.title = compact ? "Perbesar menu" : "Minimize menu";
  $("span", button).textContent = compact ? "›" : "‹";
}

function toggleSidebar() {
  state.sidebarCollapsed = !state.sidebarCollapsed;
  try { localStorage.setItem("pantas_sidebar_collapsed", String(state.sidebarCollapsed)); } catch {}
  applySidebarState();
}

const baseNavigation = [
  { section: "Pribadi" },
  { id: "dashboard", label: "Dashboard", icon: "⌂" },
  { id: "history", label: "Riwayat Potongan", icon: "⌁" },
  { id: "deductions", label: "Detail Potongan", icon: "≣" },
  { id: "appeals", label: "Banding Saya", icon: "◇" },
  { id: "profile", label: "Profil & Keamanan", icon: "○" },
];

function buildNavigation() {
  const items = state.user.is_admin ? [
    { section: "Administrator" },
    { id: "monitoring", label: "Monitoring Kantor", icon: "◎" },
    { id: "warnings", label: "Peringatan", icon: "△" },
    { id: "admin-reviews", label: "Keputusan Banding", icon: "◆" },
    { id: "admin-users", label: "Pengguna & Unit", icon: "◉" },
    { id: "admin-imports", label: "Import Data", icon: "⇧" },
    { id: "admin-parameters", label: "Parameter", icon: "⚙" },
    { id: "profile", label: "Profil & Keamanan", icon: "○" },
  ] : [...baseNavigation];
  if (!state.user.is_admin && isSupervisor()) items.push(
    { section: "Kepemimpinan" },
    { id: "monitoring", label: "Monitoring Unit", icon: "◎" },
    { id: "warnings", label: "Peringatan", icon: "△" },
    { id: "reviews", label: "Verifikasi Banding", icon: "✓" },
  );
  $("#main-nav").innerHTML = items.map(item => item.section
    ? `<div class="nav-section">${h(item.section)}</div>`
    : `<a class="nav-link" href="#${item.id}" data-route="${item.id}" title="${h(item.label)}" aria-label="${h(item.label)}"><span class="nav-icon">${item.icon}</span><span class="nav-label">${h(item.label)}</span></a>`
  ).join("");
  $$(".nav-link").forEach(link => link.addEventListener("click", () => $("#sidebar").classList.remove("open")));
}

function navigate(route, params = "") { location.hash = `#${route}${params}`; }

async function renderRoute() {
  if (!state.user) return;
  const raw = location.hash.replace(/^#/, "") || "dashboard";
  const [route] = raw.split("?");
  if (state.user.is_admin && ["dashboard", "history", "deductions", "appeals", "reviews"].includes(route)) return navigate("monitoring");
  state.route = route;
  $$(".nav-link").forEach(link => link.classList.toggle("active", link.dataset.route === route));
  const labels = {
    dashboard: ["Ringkasan pribadi", "Dashboard"], history: ["Analisis pribadi", "Riwayat Potongan"],
    deductions: ["Transparansi", "Detail Potongan"], appeals: ["Koreksi data", "Banding Saya"],
    profile: ["Akun", "Profil & Keamanan"], monitoring: ["Kepemimpinan", "Monitoring Unit"],
    warnings: ["Deteksi dini", "Peringatan"], reviews: ["Alur banding", "Verifikasi Atasan"],
    "admin-reviews": ["Alur banding", "Keputusan Administrator"], "admin-users": ["Administrasi", "Pengguna & Unit"],
    "admin-imports": ["Administrasi", "Import Data Bulanan"], "admin-parameters": ["Administrasi", "Parameter Sistem"],
  };
  const label = labels[route] || ["PANTAS", "Halaman"];
  $("#page-kicker").textContent = label[0]; $("#page-title").textContent = label[1];
  page.innerHTML = `<div class="loading"><div class="spinner" aria-label="Memuat"></div></div>`;
  page.focus();
  try {
    const handlers = {
      dashboard: renderDashboard, history: renderHistory, deductions: renderDeductions,
      appeals: renderAppeals, profile: renderProfile, monitoring: renderMonitoring,
      warnings: renderWarnings, reviews: () => renderReviews(false),
      "admin-reviews": () => renderReviews(true), "admin-users": renderAdminUsers,
      "admin-imports": renderAdminImports, "admin-parameters": renderAdminParameters,
    };
    if (!handlers[route]) return navigate("dashboard");
    await handlers[route]();
  } catch (error) {
    page.innerHTML = errorState(error.message);
  }
}

async function renderDashboard() {
  const [dashboard, historyData] = await Promise.all([api("/api/dashboard"), api("/api/history")]);
  dashboard.history = historyData.points;
  state.dashboard = dashboard;
  updateNotificationBadge(dashboard.unread_notifications);
  if (!dashboard.current_period) {
    page.innerHTML = pageIntro("Belum ada periode berjalan", "Dashboard akan aktif setelah administrator mempublikasikan data bulanan.") + emptyState("□", "Data belum dipublikasikan", "Administrator dapat mengunggah workbook melalui menu Import Data.");
    return;
  }
  const summary = dashboard.summary;
  const appealCallout = dashboard.can_appeal
    ? `<div class="callout"><span>◇</span><div><strong>Ada potongan pada periode ini</strong><br>${dashboard.pending_appeal_items ? `${dashboard.pending_appeal_items} hari masih diproses.` : `Anda dapat mengajukan penjelasan per tanggal.`}<div class="space-top"><button class="button button-small button-secondary" data-go="appeals">Buka menu banding</button></div></div></div>`
    : `<div class="callout callout-info"><span>✓</span><div><strong>Tidak ada potongan</strong><br>Tidak diperlukan pengajuan banding pada periode ini.</div></div>`;
  page.innerHTML = `
    ${pageIntro(`Halo, ${h(firstName(state.user.name))}`, "Berikut ringkasan presensi dan potongan Anda.", periodPill(dashboard.current_period.label))}
    <section class="grid kpi-grid">
      ${kpi("Potongan efektif", percent(summary.effective_deduction), "Setelah keputusan banding", "%", summary.effective_deduction > 0 ? "kpi-danger" : "kpi-green")}
      ${kpi("Hari kerja", numberID(summary.work_days), "P=1 · M=1 · PM=2", "▣", "")}
      ${kpi("Hari lembur", numberID(summary.overtime_days), "L1 dan L2", "↗", "kpi-gold")}
      ${kpi("Cuti / Off", `${summary.leave_days} / ${summary.off_days}`, "Hari tercatat", "○", "kpi-green")}
    </section>
    <section class="grid content-grid">
      <div class="panel"><div class="panel-header"><div><h3>Tren 12 periode terakhir</h3><p>Potongan awal dibanding potongan efektif</p></div><button class="button button-small button-ghost" data-go="history">Atur rentang</button></div><div class="panel-body"><div class="chart-wrap"><canvas id="history-chart" role="img" aria-label="Grafik riwayat potongan"></canvas></div></div></div>
      <div class="grid">
        <div class="panel"><div class="panel-header"><div><h3>Ringkasan periode</h3><p>${h(dashboard.current_period.start)} — ${h(dashboard.current_period.end)}</p></div></div><div class="panel-body summary-list">
          ${summaryRow("Potongan awal", percent(summary.original_deduction))}${summaryRow("Hari terkena potongan", summary.deduction_days)}${summaryRow("Cuti", summary.leave_days)}${summaryRow("Off", summary.off_days)}
        </div></div>${appealCallout}
      </div>
    </section>
    <section class="panel section-top"><div class="panel-header"><div><h3>Detail potongan periode berjalan</h3><p>Hanya tanggal dengan potongan yang ditampilkan</p></div><button class="button button-small button-ghost" data-go="deductions">Lihat lengkap</button></div>${deductionTable(dashboard.deductions, 5)}</section>`;
  bindGoButtons();
  requestAnimationFrame(() => drawHistoryChart(historyData.points));
}

async function renderHistory() {
  const query = new URLSearchParams(location.hash.split("?")[1] || "");
  const now = new Date();
  const to = query.get("to") || `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,"0")}`;
  const start = new Date(`${to}-01T00:00:00`); start.setMonth(start.getMonth()-11);
  const from = query.get("from") || `${start.getFullYear()}-${String(start.getMonth()+1).padStart(2,"0")}`;
  const data = await api(`/api/history?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`);
  page.innerHTML = `${pageIntro("Riwayat potongan", "Gunakan rentang bulan untuk melihat pola dalam periode tertentu.")}
    <section class="panel"><form id="history-filter" class="filter-bar"><label class="field"><span>Dari bulan</span><input type="month" name="from" value="${h(from)}" required></label><label class="field"><span>Sampai bulan</span><input type="month" name="to" value="${h(to)}" required></label><button class="button button-primary" type="submit">Terapkan</button></form><div class="panel-body"><div class="chart-wrap"><canvas id="history-chart"></canvas></div></div>${historyTable(data.points)}</section>`;
  $("#history-filter").addEventListener("submit", event => { event.preventDefault(); const values=formJSON(event.currentTarget); navigate("history", `?from=${values.from}&to=${values.to}`); });
  requestAnimationFrame(() => drawHistoryChart(data.points));
}

async function renderDeductions() {
  const data = await api("/api/deductions");
  page.innerHTML = `${pageIntro("Detail potongan", "Setiap kode, jam presensi, dan status koreksi ditampilkan per tanggal.", `<button class="button button-secondary" data-go="appeals">Ajukan banding</button>`)}<section class="panel">${deductionTable(data.items)}</section>`;
  bindGoButtons();
}

async function renderMonitoring() {
  const query = new URLSearchParams(location.hash.split("?")[1] || "");
  const unit = query.get("unit_id");
  const data = await api(`/api/monitoring${unit ? `?unit_id=${encodeURIComponent(unit)}` : ""}`);
  if (!data.current_period) { page.innerHTML = emptyState("□","Belum ada periode","Monitoring aktif setelah data dipublikasikan."); return; }
  const back = unit ? `<button class="button button-ghost" data-monitor-back>← Kembali ke agregat</button>` : "";
  const content = data.mode === "people" ? peopleMonitoringTable(data.items) : aggregateMonitoringCards(data.items);
  page.innerHTML = `${pageIntro(data.mode === "people" ? "Detail pegawai" : "Ringkasan unit", "Cakupan mengikuti jabatan dan struktur organisasi Anda.", `${back}${periodPill(data.current_period.label)}`)}${monitoringTotals(data.totals)}<div class="section-top">${content}</div>`;
  $("[data-monitor-back]")?.addEventListener("click", () => navigate("monitoring"));
  $$('[data-monitor-unit]').forEach(button => button.addEventListener("click", () => navigate("monitoring", `?unit_id=${button.dataset.monitorUnit}`)));
}

async function renderWarnings() {
  const data = await api("/api/warnings");
  const individual = data.individual.length ? data.individual.map(warningCard).join("") : emptyState("✓", "Tidak ada peringatan individu", "Tidak ada pola yang memenuhi parameter saat ini.");
  const aggregate = data.aggregate.length ? data.aggregate.map(warningCard).join("") : emptyState("✓", "Tidak ada peringatan unit", "Rata-rata unit masih di bawah batas parameter.");
  page.innerHTML = `${pageIntro("Peringatan presensi", "Deteksi berbasis parameter membantu atasan memberi perhatian lebih dini.", data.period ? periodPill(data.period.label) : "")}
    <div class="grid content-grid"><section class="panel"><div class="panel-header"><div><h3>Anggota</h3><p>Anomali individu dan potongan berturut-turut</p></div><span class="status status-warning">${data.individual.length} peringatan</span></div><div class="panel-body card-list">${individual}</div></section><section class="panel"><div class="panel-header"><div><h3>Keseluruhan unit</h3><p>Lonjakan dan batas rata-rata</p></div></div><div class="panel-body card-list">${aggregate}</div></section></div>`;
}

async function renderAppeals() {
  const [historyData, options] = await Promise.all([api("/api/appeals"), api("/api/appeals/options")]);
  const existing = historyData.items.length ? historyData.items.map(appealHistoryCard).join("") : `<div class="callout callout-info">Belum ada riwayat banding.</div>`;
  const form = options.days.length ? appealForm(options) : emptyState("◇", "Tidak ada hari yang dapat dibanding", "Hari tanpa potongan atau yang sudah pernah diajukan tidak ditampilkan.");
  page.innerHTML = `${pageIntro("Banding potongan", "Isi penjelasan terpisah untuk setiap tanggal. Dokumen bersifat opsional.")}
    <div class="grid content-grid"><section><div class="section-heading"><h3>Pengajuan baru</h3><p>${options.days.length} hari tersedia</p></div>${form}</section><section><div class="section-heading"><h3>Riwayat</h3><p>Status verifikasi dan keputusan</p></div><div class="card-list">${existing}</div></section></div>`;
  $("#appeal-form")?.addEventListener("submit", event => submitAppeal(event, options));
}

async function submitAppeal(event, options) {
  event.preventDefault();
  const form = event.currentTarget;
  const button = $("button[type=submit]", form);
  const items = options.days.map(day => ({
    attendance_id: day.attendance_id,
    reason_id: form.elements[`reason_${day.attendance_id}`].value,
    explanation: form.elements[`explanation_${day.attendance_id}`].value,
  }));
  setBusy(button, true, "Mengajukan…");
  try {
    const result = await api("/api/appeals", { method: "POST", body: { period_id: options.period_id, items } });
    const createdByAttendance = new Map(result.items.map(item => [String(item.attendance_id), item.id]));
    const failed = [];
    for (const day of options.days) {
      const file = form.elements[`document_${day.attendance_id}`].files[0];
      if (!file) continue;
      try {
        await api(`/api/appeals/items/${createdByAttendance.get(String(day.attendance_id))}/document`, { method: "POST", rawBody: file, headers: { "Content-Type": file.type || "application/octet-stream", "X-Filename": encodeURIComponent(file.name) } });
      } catch { failed.push(formatDate(day.date)); }
    }
    toast("Banding diajukan", failed.length ? `Pengajuan tersimpan, tetapi dokumen tanggal ${failed.join(", ")} gagal diunggah.` : "Seluruh penjelasan dan dokumen berhasil disimpan.", failed.length ? "error" : "success");
    await renderAppeals();
  } catch (error) { toast("Banding belum terkirim", error.message, "error"); }
  finally { setBusy(button, false); }
}

async function renderReviews(admin) {
  const data = await api(admin ? "/api/reviews/admin" : "/api/reviews/supervisor");
  const cards = data.items.length ? data.items.map(item => reviewCard(item, admin)).join("") : emptyState("✓", "Antrean kosong", admin ? "Belum ada banding yang menunggu keputusan administrator." : "Belum ada banding yang menunggu verifikasi Anda.");
  page.innerHTML = `${pageIntro(admin ? "Keputusan final" : "Verifikasi atasan", admin ? "Keputusan dapat menghapus atau menyesuaikan potongan efektif per hari." : "Periksa alasan dan dokumen setiap tanggal secara terpisah.")}<section class="review-grid">${cards}</section>`;
  $$('[data-review-submit]').forEach(button => button.addEventListener("click", () => submitReview(button, admin)));
}

async function submitReview(button, admin) {
  const card = button.closest(".review-card");
  const id = button.dataset.reviewSubmit;
  const decision = button.dataset.decision;
  const comment = $("textarea", card).value;
  const body = { decision, comment };
  if (admin && decision === "approved") {
    const adjusted = $("input[name=adjusted_rate]", card)?.value;
    if (adjusted !== "") body.adjusted_rate = Number(adjusted) / 100;
  }
  const ok = await confirmAction(decision === "rejected" ? "Tolak banding?" : "Simpan keputusan?", "Keputusan diterapkan untuk tanggal ini dan tercatat pada audit trail.", decision === "rejected");
  if (!ok) return;
  setBusy(button, true, "Menyimpan…");
  try {
    await api(`${admin ? "/api/reviews/admin" : "/api/reviews/supervisor"}/${id}`, { method: "POST", body });
    toast("Keputusan tersimpan", "Antrean telah diperbarui.", "success");
    await renderReviews(admin);
  } catch (error) { toast("Belum tersimpan", error.message, "error"); }
  finally { setBusy(button, false); }
}

async function renderAdminUsers(pageNumber = 1) {
  if (!state.units) state.units = (await api("/api/admin/units")).items;
  const data = await api(`/api/admin/users?page=${pageNumber}&limit=50`);
  page.innerHTML = `${pageIntro("Kelola pengguna", "Tambah, pindahkan, nonaktifkan, atau reset password pegawai.", `<button class="button button-primary" data-new-user>+ Tambah pengguna</button>`)}
    <section class="panel"><form id="user-search" class="filter-bar"><label class="field grow"><span>Pencarian</span><input name="q" placeholder="Nama, NIP, atau unit"></label><button class="button button-secondary">Cari</button></form>${adminUserTable(data.items)}<div class="pagination"><button class="button button-small button-ghost" data-user-page="${Math.max(1,data.page-1)}" ${data.page<=1?"disabled":""}>Sebelumnya</button><span class="period-pill">${data.page} · ${numberID(data.total)} pengguna</span><button class="button button-small button-ghost" data-user-page="${data.page+1}" ${data.page*data.limit>=data.total?"disabled":""}>Berikutnya</button></div></section>`;
  $("[data-new-user]").addEventListener("click", () => openUserDialog(null));
  $("#user-search").addEventListener("submit", async event => { event.preventDefault(); const q=event.currentTarget.elements.q.value; const result=await api(`/api/admin/users?q=${encodeURIComponent(q)}&limit=100`); $(".table-wrap").outerHTML=adminUserTable(result.items); bindUserActions(); });
  $$('[data-user-page]').forEach(button => button.addEventListener("click", () => renderAdminUsers(Number(button.dataset.userPage))));
  bindUserActions();
}

function bindUserActions() {
  $$('[data-edit-user]').forEach(button => button.addEventListener("click", () => openUserDialog(JSON.parse(button.dataset.user))));
  $$('[data-reset-user]').forEach(button => button.addEventListener("click", async () => {
    if (!await confirmAction("Reset password?", "Password akan kembali menjadi NIP dan seluruh sesi pengguna dicabut.")) return;
    await api(`/api/admin/users/${button.dataset.resetUser}/reset-password`, { method:"POST" }); toast("Password direset","Pengguna wajib menggantinya saat login.","success");
  }));
  $$('[data-delete-user]').forEach(button => button.addEventListener("click", async () => {
    if (!await confirmAction("Hapus pengguna?", "Akun dinonaktifkan, kontak dihapus, tetapi riwayat tetap dipertahankan.", true)) return;
    await api(`/api/admin/users/${button.dataset.deleteUser}`, { method:"DELETE" }); toast("Pengguna dihapus","Riwayat tetap tersimpan.","success"); await renderAdminUsers();
  }));
}

function openUserDialog(user) {
  const dialog = dynamicDialog(`${user ? "Ubah" : "Tambah"} pengguna`, `
    <form class="stack-lg" id="user-dialog-form">
      ${!user ? `<label class="field"><span>NIP</span><input name="nip" maxlength="18" inputmode="numeric" required></label>` : ""}
      <label class="field"><span>Nama</span><input name="name" value="${h(user?.name||"")}" required></label>
      <label class="field"><span>Unit</span><select name="unit_id" required>${unitOptions(state.units,user?.unit_id)}</select></label>
      <label class="field"><span>Jabatan</span><select name="role">${roleOptions(user?.role||"staff")}</select></label>
      ${user ? `<label class="check-row"><input type="checkbox" name="is_active" ${user.is_active?"checked":""}> Akun aktif</label><label class="field"><span>Dasar perubahan</span><input name="reason" placeholder="Opsional"></label>` : ""}
      <button class="button button-primary" type="submit">Simpan</button>
    </form>`);
  $("form",dialog).addEventListener("submit",async event=>{event.preventDefault();const values=formJSON(event.currentTarget);if(user)values.is_active=event.currentTarget.elements.is_active.checked;try{await api(user?`/api/admin/users/${user.id}`:"/api/admin/users",{method:user?"PATCH":"POST",body:values});dialog.close();dialog.remove();toast("Berhasil","Data pengguna disimpan.","success");await renderAdminUsers();}catch(error){toast("Belum tersimpan",error.message,"error");}});
}

async function renderAdminImports() {
  const data = await api("/api/admin/imports");
  page.innerHTML = `${pageIntro("Import data bulanan", "PANTAS memvalidasi format, NIP, unit, periode, dan duplikasi sebelum publikasi.")}
    <div class="grid content-grid"><section class="panel"><div class="panel-header"><div><h3>Unggah workbook</h3><p>Format harus sama dengan Upload dokumen.xlsx</p></div></div><div class="panel-body"><label class="import-drop"><input id="import-file" type="file" accept=".xlsx"><div><div class="empty-icon">⇧</div><h3>Pilih atau seret file Excel</h3><p>Maksimum 20 MB · sheet DETAIL WFH WFO</p></div></label><div id="import-preview"></div></div></section><section class="panel"><div class="panel-header"><div><h3>Riwayat import</h3><p>Versi draft dan publikasi</p></div></div><div class="panel-body card-list">${data.items.length?data.items.map(importCard).join(""):emptyState("□","Belum ada import","")}</div></section></div>`;
  $("#import-file").addEventListener("change", event => previewImport(event.target.files[0]));
  bindImportActions();
}

async function previewImport(file) {
  if (!file) return;
  const target=$("#import-preview"); target.innerHTML=`<div class="loading compact"><div class="spinner"></div></div>`;
  const form=new FormData();form.append("file",file);
  try{const data=await api("/api/admin/imports/preview",{method:"POST",rawBody:form});target.innerHTML=importPreview(data.preview);bindImportActions();}
  catch(error){const preview=error.data?.preview;target.innerHTML=`<div class="callout callout-danger section-top"><span>!</span><div><strong>Import belum siap</strong><br>${h(error.message)}</div></div>${preview?importPreview(preview):""}`;}
}

function bindImportActions(){
  $$('[data-publish-import]').forEach(button=>button.addEventListener("click",async()=>{if(!await confirmAction("Publikasikan periode?","Dashboard seluruh pengguna akan berubah dan email dijadwalkan."))return;setBusy(button,true,"Menerbitkan…");try{const data=await api(`/api/admin/imports/${button.dataset.publishImport}/publish`,{method:"POST"});toast("Data dipublikasikan",data.message,"success");await renderAdminImports();}catch(error){toast("Gagal publikasi",error.message,"error");}finally{setBusy(button,false);}}));
  $$('[data-reject-import]').forEach(button=>button.addEventListener("click",async()=>{if(!await confirmAction("Batalkan draft?","Data staging pada draft tidak akan dipublikasikan.",true))return;await api(`/api/admin/imports/${button.dataset.rejectImport}`,{method:"DELETE"});await renderAdminImports();}));
}

async function renderAdminParameters() {
  const [parameters,rules,reasons]=await Promise.all([api("/api/admin/parameters"),api("/api/admin/rules"),api("/api/admin/reasons")]);
  page.innerHTML=`${pageIntro("Parameter sistem","Perubahan hanya dapat dilakukan administrator dan dicatat pada audit trail.")}
    <section class="panel"><div class="panel-header"><div><h3>Deteksi peringatan & dashboard</h3><p>Persentase disimpan sebagai nilai 0 sampai 1</p></div></div><div class="panel-body settings-grid">${parameters.items.map(parameterCard).join("")}</div></section>
    <section class="panel section-top"><div class="panel-header"><div><h3>Aturan potongan</h3><p>Tarif digunakan saat import berikutnya</p></div><button class="button button-small button-primary" data-add-rule>+ Tambah aturan</button></div><div class="table-wrap"><table class="data-table"><thead><tr><th>Sumber</th><th>Kode</th><th>Label</th><th>Tarif</th><th>Aktif</th><th></th></tr></thead><tbody>${rules.items.map(ruleRow).join("")}</tbody></table></div></section>
    <section class="panel section-top"><div class="panel-header"><div><h3>Kategori alasan banding</h3><p>Dapat ditambah atau dinonaktifkan</p></div><button class="button button-small button-primary" data-add-reason>+ Tambah</button></div><div class="table-wrap"><table class="data-table"><thead><tr><th>Kode</th><th>Label</th><th>Deskripsi</th><th>Aktif</th><th></th></tr></thead><tbody>${reasons.items.map(reasonRow).join("")}</tbody></table></div></section>`;
  $$('[data-save-parameter]').forEach(button=>button.addEventListener("click",()=>saveParameter(button)));
  $$('[data-save-rule]').forEach(button=>button.addEventListener("click",()=>saveRule(button)));
  $("[data-add-rule]").addEventListener("click",openRuleDialog);
  $$('[data-edit-reason]').forEach(button=>button.addEventListener("click",()=>openReasonDialog(JSON.parse(button.dataset.reason))));
  $("[data-add-reason]").addEventListener("click",()=>openReasonDialog(null));
}

async function saveParameter(button){const card=button.closest(".setting-card");const raw=$("input",card).value;let value;try{value=JSON.parse(raw);}catch{return toast("Nilai tidak valid","Gunakan JSON valid, misalnya 6, 0.005, atau {\"P\":1}.","error");}try{await api(`/api/admin/parameters/${encodeURIComponent(button.dataset.saveParameter)}`,{method:"PATCH",body:{value}});toast("Parameter disimpan","Berlaku pada perhitungan berikutnya.","success");}catch(error){toast("Belum tersimpan",error.message,"error");}}
async function saveRule(button){const row=button.closest("tr");const rate=Number($("[name=rate]",row).value)/100;const label=$("[name=label]",row).value;const is_active=$("[name=active]",row).checked;try{await api(`/api/admin/rules/${button.dataset.saveRule}`,{method:"PATCH",body:{label,rate,is_active}});toast("Aturan disimpan","Import lama tidak diubah.","success");}catch(error){toast("Belum tersimpan",error.message,"error");}}

function openRuleDialog(){
  const dialog=dynamicDialog("Tambah aturan potongan",`<form class="stack-lg" id="rule-dialog-form">
    <label class="field"><span>Sumber data</span><select name="source" required><option value="late">Terlambat</option><option value="early_leave">Pulang sebelum waktunya</option><option value="leave">Cuti</option><option value="status">Status presensi</option><option value="shift">Shift</option></select></label>
    <label class="field"><span>Kode</span><input name="code" maxlength="100" placeholder="Contoh: TL4" required><small>Harus sama persis dengan kode pada file import.</small></label>
    <label class="field"><span>Label</span><input name="label" maxlength="200" placeholder="Nama aturan yang mudah dipahami" required></label>
    <label class="field"><span>Tarif potongan (%)</span><input name="rate" type="number" min="0" max="100" step="0.01" placeholder="0,50" required></label>
    <label class="field"><span>Urutan (opsional)</span><input name="sort_order" type="number" min="-100000" max="100000" placeholder="Otomatis di urutan terakhir"><small>Kosongkan agar sistem menentukan urutan.</small></label>
    <label class="check-row"><input name="is_active" type="checkbox" checked> Langsung aktif untuk import berikutnya</label>
    <button class="button button-primary" type="submit">Tambah aturan</button>
  </form>`);
  $("form",dialog).addEventListener("submit",async event=>{
    event.preventDefault();
    const form=event.currentTarget;
    const button=$("button[type=submit]",form);
    const values=formJSON(form);
    const body={source:values.source,code:values.code,label:values.label,rate:Number(values.rate)/100,is_active:form.elements.is_active.checked};
    if(values.sort_order!=="")body.sort_order=Number(values.sort_order);
    setBusy(button,true,"Menambahkan…");
    try{
      await api("/api/admin/rules",{method:"POST",body});
      dialog.close();dialog.remove();
      toast("Aturan ditambahkan","Aturan baru digunakan pada import berikutnya.","success");
      await renderAdminParameters();
    }catch(error){toast("Aturan belum ditambahkan",error.message,"error");}
    finally{if(form.isConnected)setBusy(button,false);}
  });
}

function openReasonDialog(reason){const dialog=dynamicDialog(`${reason?"Ubah":"Tambah"} kategori alasan`,`<form class="stack-lg">${!reason?`<label class="field"><span>Kode</span><input name="code" pattern="[a-z0-9_]+" required></label>`:""}<label class="field"><span>Label</span><input name="label" value="${h(reason?.label||"")}" required></label><label class="field"><span>Deskripsi</span><textarea name="description">${h(reason?.description||"")}</textarea></label><label class="field"><span>Urutan</span><input name="sort_order" type="number" value="${reason?.sort_order||0}"></label>${reason?`<label class="check-row"><input name="is_active" type="checkbox" ${reason.is_active?"checked":""}> Aktif</label>`:""}<button class="button button-primary">Simpan</button></form>`);$("form",dialog).addEventListener("submit",async event=>{event.preventDefault();const values=formJSON(event.currentTarget);values.sort_order=Number(values.sort_order);if(reason)values.is_active=event.currentTarget.elements.is_active.checked;try{await api(reason?`/api/admin/reasons/${reason.id}`:"/api/admin/reasons",{method:reason?"PATCH":"POST",body:values});dialog.close();dialog.remove();await renderAdminParameters();}catch(error){toast("Belum tersimpan",error.message,"error");}});}

async function renderProfile(){
  if(state.user.is_admin){
    page.innerHTML=`${pageIntro("Profil & keamanan","Administrator adalah akun sistem yang terpisah dari data pegawai.")}
      <div class="grid profile-grid"><section class="panel profile-hero"><div class="avatar">${h(getInitials(state.user.name))}</div><h2>${h(state.user.name)}</h2><p>@${h(state.user.username)}</p><span class="status status-active">Administrator Sistem</span></section><section class="panel"><div class="panel-header"><div><h3>Keamanan akun</h3><p>Username admin tidak menggunakan NIP</p></div></div><div class="panel-body"><div class="contact-card"><div><strong>Password</strong><small>Gunakan password unik dan ganti secara berkala.</small></div><button class="button button-small button-secondary" data-change-password>Ganti</button></div></div></section></div>`;
    $("[data-change-password]").addEventListener("click",()=>$("#password-dialog").showModal());
    return;
  }
  page.innerHTML=`${pageIntro("Profil & keamanan","Kontak terverifikasi digunakan untuk pemulihan password.")}
    <div class="grid profile-grid"><section class="panel profile-hero"><div class="avatar">${h(getInitials(state.user.name))}</div><h2>${h(state.user.name)}</h2><p>${h(state.user.nip)}</p><span class="status status-active">${h(roleLabel(state.user.position_role))}</span><p>${h(state.user.unit_name)}</p></section><section class="panel"><div class="panel-header"><div><h3>Kontak pemulihan</h3><p>Perubahan memerlukan password dan kode OTP</p></div></div><div class="panel-body"><div class="contact-card"><div><strong>Email</strong><small>${h(state.user.email||"Belum ditambahkan")} · ${state.user.email_verified?"Terverifikasi":"Belum terverifikasi"}</small></div><button class="button button-small button-secondary" data-contact="email">Ubah</button></div><div class="contact-card"><div><strong>Nomor HP</strong><small>${h(state.user.phone||"Belum ditambahkan")} · ${state.user.phone_verified?"Terverifikasi":"Belum terverifikasi"}</small></div><button class="button button-small button-secondary" data-contact="phone">Ubah</button></div><div class="contact-card"><div><strong>Password</strong><small>Ganti berkala dan jangan gunakan NIP.</small></div><button class="button button-small button-secondary" data-change-password>Ganti</button></div></div></section></div>`;
  $$('[data-contact]').forEach(button=>button.addEventListener("click",()=>openContactDialog(button.dataset.contact)));
  $("[data-change-password]").addEventListener("click",()=>$("#password-dialog").showModal());
}

function openContactDialog(channel) {
  const label = channel === "email" ? "Email" : "Nomor HP";
  const dialog = dynamicDialog(`Verifikasi ${label}`, `<form class="stack-lg" data-contact-start><label class="field"><span>${label} baru</span><input name="destination" type="${channel === "email" ? "email" : "tel"}" required></label><label class="field"><span>Password saat ini</span><input name="current_password" type="password" required></label><button class="button button-primary">Kirim kode</button></form>`);
  const startForm = $("[data-contact-start]", dialog);
  startForm.addEventListener("submit", async event => {
    event.preventDefault();
    const form = event.currentTarget;
    const button = $("button[type=submit]", form);
    const values = formJSON(form);
    setBusy(button, true, "Mengirim…");
    try {
      const delivery = await api("/api/profile/contact/start", { method: "POST", body: { channel, ...values } });
      form.outerHTML = `<form class="stack-lg" data-contact-verify><label class="field"><span>Kode 6 digit</span><input class="otp-input" name="otp" inputmode="numeric" pattern="[0-9]{6}" maxlength="6" autocomplete="one-time-code" required></label><button class="button button-primary">Verifikasi</button></form>`;
      toast("Kode terkirim", delivery.message, "success");
      const verifyForm = $("[data-contact-verify]", dialog);
      verifyForm.addEventListener("submit", async verifyEvent => {
        verifyEvent.preventDefault();
        const otpForm = verifyEvent.currentTarget;
        const verifyButton = $("button[type=submit]", otpForm);
        const otp = otpForm.elements.otp.value;
        setBusy(verifyButton, true, "Memverifikasi…");
        try {
          await api("/api/profile/contact/verify", { method: "POST", body: { channel, otp } });
          dialog.close();
          dialog.remove();
          state.user = (await api("/api/auth/me")).user;
          toast("Kontak terverifikasi", "Kontak dapat digunakan untuk pemulihan.", "success");
          await renderProfile();
        } catch (error) {
          toast("Kode salah", error.message, "error");
        } finally {
          setBusy(verifyButton, false);
        }
      });
    } catch (error) {
      toast("Belum dapat mengirim", error.message, "error");
      if (error.code === "current_password_invalid" && form.isConnected) {
        form.elements.current_password.value = "";
        form.elements.current_password.setAttribute("aria-invalid", "true");
        form.elements.current_password.focus();
      }
    } finally {
      if (form.isConnected) setBusy(button, false);
    }
  });
}

async function openNotifications(){const drawer=$("#notification-drawer");drawer.hidden=false;$("#drawer-backdrop").hidden=false;const data=await api("/api/notifications");$("#notification-list").innerHTML=data.items.length?data.items.map(notificationItem).join(""):emptyState("♢","Belum ada notifikasi","");$$('[data-notification]',drawer).forEach(item=>item.addEventListener("click",async()=>{if(!item.classList.contains("unread"))return;await api(`/api/notifications/${item.dataset.notification}/read`,{method:"POST"});item.classList.remove("unread");}));}
function closeNotifications(){$("#notification-drawer").hidden=true;$("#drawer-backdrop").hidden=true;}

async function revealDocuments(button){const target=button.parentElement.querySelector("[data-document-list]");if(!target)return;try{const data=await api(`/api/appeals/items/${button.dataset.documents}/documents`);target.innerHTML=data.items.length?data.items.map(item=>`<a class="button button-small button-ghost" href="/api/documents/${item.id}" target="_blank" rel="noopener">${h(item.filename)} · ${fileSize(item.size)}</a>`).join(""):"<small>Tidak ada dokumen.</small>";}catch(error){toast("Dokumen belum dapat dibuka",error.message,"error");}}

function drawHistoryChart(points){
  const canvas=$("#history-chart");if(!canvas)return;const rect=canvas.getBoundingClientRect();const dpr=window.devicePixelRatio||1;canvas.width=Math.max(300,rect.width*dpr);canvas.height=265*dpr;const ctx=canvas.getContext("2d");ctx.scale(dpr,dpr);const width=rect.width,height=265,pad={l:44,r:18,t:18,b:42};ctx.clearRect(0,0,width,height);if(!points.length){ctx.fillStyle="#708399";ctx.font="12px system-ui";ctx.textAlign="center";ctx.fillText("Belum ada riwayat",width/2,height/2);return;}const maxRate=Math.max(.01,...points.map(p=>Math.max(p.original,p.effective)));const niceMax=Math.ceil(maxRate*100/2)*2/100;ctx.strokeStyle="#e5ebf0";ctx.lineWidth=1;ctx.fillStyle="#708399";ctx.font="10px system-ui";ctx.textAlign="right";for(let i=0;i<=4;i++){const y=pad.t+(height-pad.t-pad.b)*i/4;ctx.beginPath();ctx.moveTo(pad.l,y);ctx.lineTo(width-pad.r,y);ctx.stroke();ctx.fillText(percent(niceMax*(1-i/4)),pad.l-7,y+3);}const x=i=>points.length===1?(pad.l+width-pad.r)/2:pad.l+(width-pad.l-pad.r)*i/(points.length-1);const y=value=>pad.t+(height-pad.t-pad.b)*(1-value/niceMax);const series=(key,color)=>{ctx.strokeStyle=color;ctx.lineWidth=2.5;ctx.lineJoin="round";ctx.beginPath();points.forEach((p,i)=>i?ctx.lineTo(x(i),y(p[key])):ctx.moveTo(x(i),y(p[key])));ctx.stroke();points.forEach((p,i)=>{ctx.fillStyle="#fff";ctx.strokeStyle=color;ctx.lineWidth=2;ctx.beginPath();ctx.arc(x(i),y(p[key]),3.5,0,Math.PI*2);ctx.fill();ctx.stroke();});};series("original","#c97822");series("effective","#1479bf");ctx.fillStyle="#708399";ctx.textAlign="center";points.forEach((p,i)=>{if(points.length<=12||i%Math.ceil(points.length/12)===0)ctx.fillText(shortMonth(p.label),x(i),height-17);});ctx.textAlign="left";ctx.fillStyle="#c97822";ctx.fillRect(pad.l,2,10,3);ctx.fillStyle="#40556b";ctx.fillText("Awal",pad.l+15,7);ctx.fillStyle="#1479bf";ctx.fillRect(pad.l+60,2,10,3);ctx.fillStyle="#40556b";ctx.fillText("Efektif",pad.l+75,7);
}

function pageIntro(title,subtitle,action=""){return `<header class="page-intro"><div><span class="eyebrow">PANTAS</span><h2>${h(title)}</h2><p>${h(subtitle)}</p></div>${action?`<div class="page-actions">${action}</div>`:""}</header>`;}
function periodPill(label){return `<span class="period-pill">◷ ${h(label)}</span>`;}
function kpi(label,value,note,icon,kind){return `<article class="kpi-card ${kind}"><div class="kpi-label"><span>${h(label)}</span><span class="kpi-icon">${icon}</span></div><div class="kpi-value">${h(String(value))}</div><div class="kpi-note">${h(note)}</div></article>`;}
function summaryRow(label,value){return `<div class="summary-row"><span>${h(String(label))}</span><strong>${h(String(value))}</strong></div>`;}
function emptyState(icon,title,text){return `<div class="empty-state"><div><div class="empty-icon">${icon}</div><h3>${h(title)}</h3><p>${h(text)}</p></div></div>`;}
function errorState(message){return pageIntro("Halaman belum dapat dimuat","Coba beberapa saat lagi.")+`<div class="callout callout-danger"><span>!</span><div><strong>Terjadi kendala</strong><br>${h(message)}</div></div>`;}
function deductionTable(items,limit){const shown=limit?items.slice(0,limit):items;if(!shown.length)return emptyState("✓","Tidak ada potongan","Tidak ada tanggal dengan potongan pada periode ini.");return `<div class="table-wrap"><table class="data-table"><thead><tr><th>Tanggal</th><th>Jam masuk</th><th>Jam pulang</th><th>Kode</th><th>Potongan awal</th><th>Efektif</th><th>Status banding</th></tr></thead><tbody>${shown.map(item=>`<tr><td><strong>${formatDate(item.date)}</strong></td><td>${h(item.check_in||"—")}</td><td>${h(item.check_out||"—")}</td><td><span class="code-list">${(item.components||[]).map(c=>`<span class="code-chip">${h(c.code)}</span>`).join("")}</span></td><td class="rate">${percent(item.original)}</td><td class="rate ${item.effective===0?"rate-zero":""}">${percent(item.effective)}</td><td>${item.admin_status?status(item.admin_status):"—"}</td></tr>`).join("")}</tbody></table></div>`;}
function historyTable(points){if(!points.length)return emptyState("⌁","Belum ada data","Periode yang dipilih belum memiliki publikasi.");return `<div class="table-wrap"><table class="data-table"><thead><tr><th>Periode</th><th>Rentang</th><th class="number">Hari potong</th><th class="number">Awal</th><th class="number">Efektif</th></tr></thead><tbody>${points.slice().reverse().map(p=>`<tr><td><strong>${h(p.label)}</strong></td><td>${formatDate(p.start)} – ${formatDate(p.end)}</td><td class="number">${p.deduction_days}</td><td class="number rate">${percent(p.original)}</td><td class="number">${percent(p.effective)}</td></tr>`).join("")}</tbody></table></div>`;}
function peopleMonitoringTable(items){return `<section class="panel">${items.length?`<div class="table-wrap"><table class="data-table"><thead><tr><th>Pegawai</th><th>NIP</th><th>Jabatan</th><th class="number">Hari kerja</th><th class="number">Hari potong</th><th class="number">Potongan</th></tr></thead><tbody>${items.map(item=>`<tr><td><strong>${h(item.name)}</strong><br><small>${h(item.unit_name)}</small></td><td>${h(item.nip)}</td><td>${h(roleLabel(item.role))}</td><td class="number">${item.work_days}</td><td class="number">${item.deduction_days}</td><td class="number rate">${percent(item.effective)}</td></tr>`).join("")}</tbody></table></div>`:emptyState("○","Unit kosong","")}</section>`;}
function monitoringTotals(totals){if(!totals)return"";return `<section class="grid kpi-grid">${kpi("Total potongan efektif",percent(totals.effective),"Akumulasi cakupan kewenangan","%",totals.effective>0?"kpi-danger":"kpi-green")}${kpi("Rata-rata pegawai",percent(totals.average),"Termasuk pegawai tanpa potongan","◎","")}${kpi("Pegawai dipantau",numberID(totals.members),"Akun aktif dalam cakupan","◉","kpi-green")}${kpi("Hari terkena potongan",numberID(totals.deduction_days),"Akumulasi seluruh pegawai","△","kpi-gold")}</section>`;}
function aggregateMonitoringCards(items){return `<section class="grid kpi-grid">${items.map(item=>`<article class="data-card"><div class="data-card-header"><div><span class="eyebrow">${h(item.unit_type)}</span><h3>${h(item.unit_name)}</h3></div><span class="status ${item.average>.005?"status-warning":"status-approved"}">${percent(item.average)} rata-rata</span></div><p>${numberID(item.members)} pegawai · ${numberID(item.deduction_days)} hari potongan</p><div class="summary-list">${summaryRow("Total efektif",percent(item.effective))}${summaryRow("Total awal",percent(item.original))}</div>${item.detail_allowed?`<button class="button button-small button-secondary space-top" data-monitor-unit="${item.unit_id}">Lihat detail</button>`:""}</article>`).join("")}</section>`;}
function warningCard(item){return `<article class="data-card warning-card"><div class="warning-symbol ${item.severity==="danger"?"danger":""}">!</div><div><h4>${h(item.name||item.unit)}</h4><p>${h(item.message)}</p><div class="card-meta">${item.nip?`<span>${h(item.nip)}</span>`:""}${item.unit&&item.name?`<span>· ${h(item.unit)}</span>`:""}</div></div><strong class="rate">${percent(item.current_rate??item.current_average??0)}</strong></article>`;}
function appealForm(options){return `<form id="appeal-form" class="card-list">${options.days.map(day=>`<article class="appeal-day"><div class="appeal-day-summary"><span class="eyebrow">Tanggal potongan</span><div class="appeal-date">${formatDateLong(day.date)}</div><div class="code-list">${day.components.map(c=>`<span class="code-chip">${h(c.code)} · ${percent(c.rate)}</span>`).join("")}</div><p class="rate">Total ${percent(day.rate)}</p></div><div class="appeal-form-fields"><label class="field"><span>Kategori alasan</span><select name="reason_${day.attendance_id}" required><option value="">Pilih alasan</option>${options.reasons.map(r=>`<option value="${r.id}">${h(r.label)}</option>`).join("")}</select></label><label class="field"><span>Penjelasan</span><textarea name="explanation_${day.attendance_id}" minlength="10" maxlength="3000" placeholder="Jelaskan kronologi dan alasan koreksi untuk tanggal ini…" required></textarea></label><label class="file-drop"><span>＋</span><span><strong>Dokumen pendukung (opsional)</strong><small>PDF, JPG, atau PNG · maks. 5 MB</small></span><input name="document_${day.attendance_id}" type="file" accept="application/pdf,image/jpeg,image/png"></label></div></article>`).join("")}<button class="button button-primary" type="submit">Ajukan ${options.days.length} hari banding</button></form>`;}
function appealHistoryCard(appeal){return `<article class="data-card"><div class="data-card-header"><div><span class="eyebrow">${h(appeal.period_label)}</span><h3>${appeal.items.length} hari diajukan</h3></div>${status(appeal.status)}</div><div class="card-list space-top">${appeal.items.map(item=>`<div class="summary-row"><div><strong>${formatDate(item.date)}</strong><small class="block">${h(item.reason)} · ${statusText(item.admin_status)}</small><div data-document-list></div></div><div class="text-right"><span class="rate">${percent(item.original)} → ${percent(item.adjusted)}</span>${item.document_count?`<button class="link-button block" data-documents="${item.id}">${item.document_count} dokumen</button>`:""}</div></div>`).join("")}</div></article>`;}
function reviewCard(item,admin){return `<article class="data-card review-card"><div class="data-card-header"><div><span class="eyebrow">${h(item.period)} · ${formatDate(item.date)}</span><h3>${h(item.name)}</h3><p>${h(item.nip)} · ${h(item.unit)}</p></div><span class="rate">${percent(item.rate)}</span></div><div><span class="code-chip">${h(item.reason)}</span><p>${h(item.explanation)}</p></div>${item.document_count?`<div><button class="button button-small button-ghost" data-documents="${item.id}">Buka ${item.document_count} dokumen</button><div data-document-list class="page-actions space-top"></div></div>`:""}${admin?`<div class="callout callout-info"><span>✓</span><div><strong>Verifikasi atasan: ${h(statusText(item.supervisor_status))}</strong><br>${h(item.supervisor_comment||"Tanpa komentar")}</div></div>`:""}<div class="review-actions"><label class="field"><span>Komentar (opsional)</span><textarea maxlength="2000" placeholder="Catatan keputusan…"></textarea></label>${admin?`<label class="field"><span>Potongan hasil jika disetujui (%)</span><input name="adjusted_rate" type="number" min="0" max="${item.rate*100}" step="0.01" value="0"></label>`:""}<div class="review-buttons"><button class="button button-danger" data-review-submit="${item.id}" data-decision="${admin?"rejected":"rejected"}">Tolak</button><button class="button button-primary" data-review-submit="${item.id}" data-decision="${admin?"approved":"accepted"}">${admin?"Setujui":"Terima verifikasi"}</button></div></div></article>`;}
function adminUserTable(items){return `<div class="table-wrap"><table class="data-table"><thead><tr><th>Nama / NIP</th><th>Unit</th><th>Jabatan</th><th>Status</th><th>Kontak</th><th></th></tr></thead><tbody>${items.map(user=>`<tr><td><strong>${h(user.name)}</strong><br><small>${h(user.nip)}</small></td><td>${h(user.unit_name)}</td><td>${h(roleLabel(user.role))}</td><td>${status(user.is_active?"active":"inactive")}${user.must_change_password?`<br><small>password awal</small>`:""}</td><td>${user.email_verified?"Email ✓":"—"}${user.phone_verified?" · HP ✓":""}</td><td><div class="page-actions"><button class="button button-small button-ghost" data-edit-user data-user='${attributeJSON(user)}'>Ubah</button><button class="button button-small button-ghost" data-reset-user="${user.id}">Reset</button><button class="button button-small button-ghost" data-delete-user="${user.id}">Hapus</button></div></td></tr>`).join("")}</tbody></table></div>`;}
function importCard(item){return `<article class="data-card"><div class="data-card-header"><div><span class="eyebrow">Versi ${item.version}</span><h4>${h(item.label)}</h4></div>${status(item.status)}</div><p>${h(item.filename)}</p><div class="card-meta"><span>${numberID(item.rows)} baris</span><span>· ${numberID(item.employees)} pegawai</span><span>· ${h(item.integrity)}</span></div>${item.status==="draft"?`<div class="page-actions space-top"><button class="button button-small button-primary" data-publish-import="${item.id}">Publikasikan</button><button class="button button-small button-ghost" data-reject-import="${item.id}">Batalkan</button></div>`:""}</article>`;}
function importPreview(p){const ready=p.ready_to_publish;return `<div class="callout ${ready?"callout-info":"callout-danger"} section-top"><span>${ready?"✓":"!"}</span><div><strong>${ready?"Validasi berhasil":"Perlu perbaikan"}</strong><br>${h(p.period_label)} · ${h(p.integrity_status)}</div></div><div class="preview-grid">${previewStat("Baris aktif",numberID(p.rows))}${previewStat("Pegawai",numberID(p.employees))}${previewStat("Hari potongan",numberID(p.deduction_days))}${previewStat("Total akumulasi",percent(p.total_deduction))}${previewStat("Baris kosong diabaikan",numberID(p.blank_rows_ignored))}${previewStat("Beda unit",numberID(p.unit_mismatches?.length||0))}</div>${p.warnings?.length?`<div class="callout"><span>i</span><div>${p.warnings.map(h).join("<br>")}</div></div>`:""}${p.missing_users?.length?`<div class="callout callout-danger"><span>!</span><div><strong>${p.missing_users.length} NIP belum terdaftar</strong><br>${p.missing_users.slice(0,8).map(x=>`${h(x.name)} (${h(x.nip)})`).join("<br>")}</div></div>`:""}${ready?`<button class="button button-primary button-full section-top" data-publish-import="${p.batch_id}">Publikasikan ${h(p.period_label)}</button>`:""}`;}
function previewStat(label,value){return `<div class="preview-stat"><small>${h(label)}</small><strong>${h(String(value))}</strong></div>`;}
function parameterCard(item){return `<article class="data-card setting-card"><div><h4>${h(item.label)}</h4><p>${h(item.description)}</p><small>${h(item.key)} · ${h(item.value_type)}</small></div><input class="table-input" value="${h(JSON.stringify(item.value))}" aria-label="${h(item.label)}"><button class="button button-small button-secondary" data-save-parameter="${h(item.key)}">Simpan</button></article>`;}
function ruleRow(item){return `<tr><td>${h(item.source)}</td><td><span class="code-chip">${h(item.code)}</span></td><td><input class="table-input" name="label" value="${h(item.label)}"></td><td><input class="table-input" name="rate" type="number" min="0" max="100" step="0.01" value="${item.rate*100}">%</td><td><input name="active" type="checkbox" ${item.is_active?"checked":""}></td><td><button class="button button-small button-secondary" data-save-rule="${item.id}">Simpan</button></td></tr>`;}
function reasonRow(item){return `<tr><td>${h(item.code)}</td><td>${h(item.label)}</td><td>${h(item.description)}</td><td>${status(item.is_active?"active":"inactive")}</td><td><button class="button button-small button-ghost" data-edit-reason data-reason='${attributeJSON(item)}'>Ubah</button></td></tr>`;}
function notificationItem(item){return `<article class="notification-item ${item.read?"":"unread"}" data-notification="${item.id}"><h4>${h(item.title)}</h4><p>${h(item.body)}</p><time>${formatDateTime(item.created_at)}</time></article>`;}
function unitOptions(units,selected){return units.filter(u=>u.is_active).map(u=>`<option value="${u.id}" ${u.id===selected?"selected":""}>${h(u.name)} · ${h(u.type)}</option>`).join("");}
function roleOptions(selected){return [["staff","Staf"],["section_head","Kepala Seksi/Subbagian"],["division_head","Kepala Bidang/Bagian"],["office_head","Kepala Kantor"],["functional","Fungsional"]].map(([value,label])=>`<option value="${value}" ${value===selected?"selected":""}>${label}</option>`).join("");}
function status(value){return `<span class="status status-${h(value)}">${h(statusText(value))}</span>`;}
function statusText(value){return ({pending:"Menunggu",accepted:"Diterima",approved:"Disetujui",rejected:"Ditolak",submitted:"Diajukan",supervisor_review:"Verifikasi atasan",admin_review:"Keputusan admin",finalized:"Selesai",published:"Dipublikasikan",draft:"Draft",superseded:"Diganti",active:"Aktif",inactive:"Nonaktif"})[value]||value||"—";}
function roleLabel(value){return ({staff:"Staf",section_head:"Kepala Seksi/Subbagian",division_head:"Kepala Bidang/Bagian",office_head:"Kepala Kantor",functional:"Fungsional",admin:"Administrator Sistem"})[value]||value;}
function bindGoButtons(){$$('[data-go]').forEach(button=>button.addEventListener("click",()=>navigate(button.dataset.go)));}

async function api(path, options={}) {
  const headers={Accept:"application/json",...(options.headers||{})};
  let body=options.rawBody;
  if(options.body!==undefined){headers["Content-Type"]="application/json";body=JSON.stringify(options.body);}
  if(options.auth!==false && !["GET","HEAD"].includes(options.method||"GET")){const csrf=getCookie("pantas_csrf");if(csrf)headers["X-CSRF-Token"]=csrf;}
  const response=await fetch(path,{method:options.method||"GET",headers,body,credentials:"same-origin"});
  if(response.status===204)return null;
  const type=response.headers.get("content-type")||"";const data=type.includes("application/json")?await response.json():null;
  if(!response.ok){const code=data?.error?.code;if(response.status===401&&code==="unauthenticated"&&options.auth!==false){state.user=null;showAuth();}if(response.status===428&&state.user){$("#password-dialog").showModal();}const error=new Error(data?.error?.message||`Permintaan gagal (${response.status})`);error.code=code;error.data=data;throw error;}
  return data;
}

function dynamicDialog(title,content){const dialog=document.createElement("dialog");dialog.className="modal";dialog.innerHTML=`<div class="modal-card"><div class="data-card-header"><h2>${h(title)}</h2><button class="icon-button" type="button" data-close>×</button></div>${content}</div>`;document.body.append(dialog);$("[data-close]",dialog).addEventListener("click",()=>{dialog.close();dialog.remove();});dialog.addEventListener("cancel",()=>setTimeout(()=>dialog.remove(),0));dialog.showModal();return dialog;}
function confirmAction(title,message,danger=false){return new Promise(resolve=>{const dialog=$("#confirm-dialog");$("#confirm-title").textContent=title;$("#confirm-message").textContent=message;$("#confirm-accept").className=`button ${danger?"button-danger":"button-primary"}`;dialog.addEventListener("close",()=>resolve(dialog.returnValue==="confirm"),{once:true});dialog.showModal();});}
function toast(title,message,type="success"){const element=document.createElement("div");element.className=`toast toast-${type}`;element.innerHTML=`<span>${type==="success"?"✓":"!"}</span><div><strong>${h(title)}</strong><p>${h(message)}</p></div>`;$("#toast-region").append(element);setTimeout(()=>element.remove(),5500);}
function setBusy(button,busy,label){if(!button)return;if(busy){button.dataset.original=button.innerHTML;button.disabled=true;button.textContent=label||"Memproses…";}else{button.disabled=false;if(button.dataset.original)button.innerHTML=button.dataset.original;}}
function formJSON(form){return Object.fromEntries(new FormData(form).entries());}
function getCookie(name){return document.cookie.split("; ").find(row=>row.startsWith(`${name}=`))?.split("=").slice(1).join("=")||"";}
function h(value){return String(value??"").replace(/[&<>'"]/g,char=>({"&":"&amp;","<":"&lt;",">":"&gt;","'":"&#39;",'"':"&quot;"})[char]);}
function attributeJSON(value){return h(JSON.stringify(value));}
function percent(value){return new Intl.NumberFormat("id-ID",{style:"percent",minimumFractionDigits:value&&Math.abs(value)<.01?2:1,maximumFractionDigits:2}).format(Number(value)||0);}
function numberID(value){return new Intl.NumberFormat("id-ID").format(Number(value)||0);}
function formatDate(value){if(!value)return"—";return new Intl.DateTimeFormat("id-ID",{day:"2-digit",month:"short",year:"numeric"}).format(new Date(`${value}T00:00:00`));}
function formatDateLong(value){return new Intl.DateTimeFormat("id-ID",{weekday:"long",day:"numeric",month:"long",year:"numeric"}).format(new Date(`${value}T00:00:00`));}
function formatDateTime(value){return new Intl.DateTimeFormat("id-ID",{day:"numeric",month:"short",hour:"2-digit",minute:"2-digit"}).format(new Date(value));}
function shortMonth(label){return label.split(" ")[0].slice(0,3);}
function fileSize(value){return value>1048576?`${(value/1048576).toFixed(1)} MB`:`${Math.ceil(value/1024)} KB`;}
function getInitials(name){return String(name).split(/\s+/).filter(Boolean).slice(0,2).map(part=>part[0]).join("");}
function firstName(name){return String(name).split(/\s+/)[0];}
function updateNotificationBadge(count){const badge=$("#notification-badge");badge.hidden=!count;badge.textContent=count>99?"99+":count;}
function isSupervisor(){return state.user?.is_admin||["section_head","division_head","office_head"].includes(state.user?.position_role);}
