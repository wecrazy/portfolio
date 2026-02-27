/* ===================================================
   ADMIN JS -- Theme, i18n, Sidebar, Datatables
   Toast & Confirm are handled by toast.js (SweetAlert2)
   =================================================== */

document.addEventListener('DOMContentLoaded', function () {

    // ---------- HTMX Toast Events ----------
    // Server-side can fire HX-Trigger headers like:
    //   HX-Trigger: {"showToast":{"message":"Saved!","type":"success"}}
    document.body.addEventListener('showToast', function (e) {
        var detail = e.detail || {};
        var message = detail.message || 'Action completed';
        var type = detail.type || 'success';
        window.showToast(message, type);
    });

    // ---------- HTMX Network / Server Error Toasts ----------
    // Fired when the request never gets a response (server down, network loss, restart).
    document.body.addEventListener('htmx:sendError', function (e) {
        var dict = adminI18nCache[adminCurrentLang] || {};
        var msg = dict['admin.toast.connection_lost'] || 'Connection lost — the server may be restarting. Please try again.';
        window.showToast(msg, 'error');
    });

    // Fired when the server returns a non-2xx status (500, 503, 429, etc.).
    document.body.addEventListener('htmx:responseError', function (e) {
        var status = (e.detail && e.detail.xhr) ? e.detail.xhr.status : 0;
        var dict = adminI18nCache[adminCurrentLang] || {};
        var msg;
        if (status === 503) {
            msg = dict['admin.toast.server_overloaded'] || 'Server is temporarily overloaded. Please try again shortly.';
        } else if (status === 429) {
            msg = dict['admin.toast.rate_limited'] || 'Too many requests — please slow down.';
        } else if (status >= 500) {
            msg = dict['admin.toast.server_error'] || 'Server error (' + status + '). Please try again.';
        } else if (status === 401 || status === 403) {
            msg = dict['admin.toast.unauthorized'] || 'Session expired. Please log in again.';
            setTimeout(function () { window.location.href = '/admin/login'; }, 2000);
        } else {
            msg = dict['admin.toast.request_failed'] || 'Request failed (' + status + '). Please try again.';
        }
        window.showToast(msg, 'error');
    });

    // ---------- Admin i18n Engine ----------
    var adminI18nCache = {};
    var adminDefaultLang = document.documentElement.getAttribute('data-default-lang') || 'en';
    var adminCurrentLang = localStorage.getItem('admin-lang') || adminDefaultLang;

    function adminApplyTranslations(dict) {
        document.querySelectorAll('[data-i18n]').forEach(function (el) {
            var key = el.getAttribute('data-i18n');
            if (dict[key]) {
                el.textContent = dict[key];
            }
        });
        document.querySelectorAll('[data-i18n-placeholder]').forEach(function (el) {
            var key = el.getAttribute('data-i18n-placeholder');
            if (dict[key]) {
                el.setAttribute('placeholder', dict[key]);
            }
        });
        var langLabel = document.getElementById('adminLangLabel');
        if (langLabel) {
            langLabel.textContent = dict['lang.' + adminCurrentLang] || adminCurrentLang.toUpperCase();
        }
        document.querySelectorAll('.admin-lang-option').forEach(function (opt) {
            opt.classList.toggle('active', opt.getAttribute('data-lang') === adminCurrentLang);
        });
        document.documentElement.setAttribute('lang', adminCurrentLang);
    }

    function adminSetLanguage(lang) {
        adminCurrentLang = lang;
        localStorage.setItem('admin-lang', lang);
        document.cookie = 'lang=' + lang + '; path=/; SameSite=Lax; max-age=31536000';
        if (adminI18nCache[lang]) {
            adminApplyTranslations(adminI18nCache[lang]);
            return;
        }
        var ver = document.documentElement.getAttribute('data-app-version') || '';
        var langUrl = '/lang/' + lang + (ver ? '?v=' + ver : '');
        fetch(langUrl)
            .then(function (r) { return r.json(); })
            .then(function (dict) {
                adminI18nCache[lang] = dict;
                // Guard: user may have switched language while this fetch was in-flight.
                if (lang !== adminCurrentLang) return;
                adminApplyTranslations(dict);
            })
            .catch(function () { /* silent fallback */ });
    }

    // Initial load
    adminSetLanguage(adminCurrentLang);

    // Language switcher click handler
    document.querySelectorAll('.admin-lang-option').forEach(function (opt) {
        opt.addEventListener('click', function (e) {
            e.preventDefault();
            adminSetLanguage(this.getAttribute('data-lang'));
        });
    });

    // ---------- Toast from URL Query Param (after redirect) ----------
    (function () {
        var params = new URLSearchParams(window.location.search);
        var toastKey = params.get('toast');
        if (toastKey) {
            var messages = {
                'login_success': 'Welcome back!',
                'logout_success': 'You have been logged out.'
            };
            var msg = messages[toastKey] || toastKey;
            if (adminI18nCache[adminCurrentLang]) {
                var i18nKey = 'admin.toast.' + toastKey;
                if (adminI18nCache[adminCurrentLang][i18nKey]) {
                    msg = adminI18nCache[adminCurrentLang][i18nKey];
                }
            }
            var toastType = toastKey === 'logout_success' ? 'info' : 'success';
            setTimeout(function () { window.showToast(msg, toastType); }, 150);
            var url = new URL(window.location);
            url.searchParams.delete('toast');
            window.history.replaceState({}, '', url.pathname + url.search);
        }
    })();

    // ---------- Admin Theme Toggle ----------
    var adminThemeToggle = document.getElementById('adminThemeToggle');
    if (adminThemeToggle) {
        updateAdminThemeIcon(document.documentElement.getAttribute('data-bs-theme') || 'dark');

        adminThemeToggle.addEventListener('click', function () {
            var current = document.documentElement.getAttribute('data-bs-theme');
            var next = current === 'dark' ? 'light' : 'dark';
            document.documentElement.setAttribute('data-bs-theme', next);
            localStorage.setItem('admin-theme', next);
            updateAdminThemeIcon(next);
        });
    }

    function updateAdminThemeIcon(theme) {
        var toggle = document.getElementById('adminThemeToggle');
        if (!toggle) return;
        var icon = toggle.querySelector('i');
        if (!icon) return;
        icon.className = theme === 'light' ? 'bxf bx-moon' : 'bxf bx-sun';
    }

    // ---------- Sidebar Active Link ----------
    var currentPath = window.location.pathname;
    var sidebarLinks = document.querySelectorAll('.admin-sidebar .sidebar-nav .nav-link');
    sidebarLinks.forEach(function (link) {
        var href = link.getAttribute('href');
        if (href === currentPath || (href !== '/admin' && currentPath.startsWith(href))) {
            link.classList.add('active');
        } else if (href === '/admin' && currentPath === '/admin') {
            link.classList.add('active');
        }
    });

    // ---------- Sidebar Toggle (desktop collapse / mobile slide) ----------
    var sidebarToggle = document.getElementById('sidebar-toggle');
    var sidebar = document.getElementById('admin-sidebar');
    var adminMain = document.querySelector('.admin-main');
    var overlay = document.getElementById('sidebar-overlay');

    function isDesktop() { return window.innerWidth >= 992; }

    // Restore desktop collapsed state on load.
    if (isDesktop() && localStorage.getItem('admin-sidebar-collapsed') === 'true') {
        sidebar.classList.add('collapsed');
        if (adminMain) adminMain.classList.add('sidebar-collapsed');
    }

    if (sidebarToggle && sidebar) {
        sidebarToggle.addEventListener('click', function () {
            if (isDesktop()) {
                var collapsed = sidebar.classList.toggle('collapsed');
                if (adminMain) adminMain.classList.toggle('sidebar-collapsed', collapsed);
                localStorage.setItem('admin-sidebar-collapsed', collapsed);
            } else {
                sidebar.classList.toggle('show');
                if (overlay) overlay.classList.toggle('show');
            }
        });

        if (overlay) {
            overlay.addEventListener('click', function () {
                sidebar.classList.remove('show');
                overlay.classList.remove('show');
            });
        }
    }

    // On resize: sync state — no collapsed class on mobile, no show class on desktop.
    window.addEventListener('resize', function () {
        if (isDesktop()) {
            sidebar.classList.remove('show');
            if (overlay) overlay.classList.remove('show');
            var savedCollapsed = localStorage.getItem('admin-sidebar-collapsed') === 'true';
            sidebar.classList.toggle('collapsed', savedCollapsed);
            if (adminMain) adminMain.classList.toggle('sidebar-collapsed', savedCollapsed);
        } else {
            sidebar.classList.remove('collapsed');
            if (adminMain) adminMain.classList.remove('sidebar-collapsed');
        }
    });

    // ---------- HTMX Confirm (Swal2 dialog) ----------
    document.body.addEventListener('htmx:confirm', function (e) {
        var trigger = e.detail.elt;
        var confirmMsg = trigger.getAttribute('data-confirm');
        if (confirmMsg) {
            e.preventDefault();
            var dict = adminI18nCache[adminCurrentLang] || {};
            var title = dict['admin.confirm.title'] || 'Are you sure?';
            var btnText = dict['admin.confirm.delete_btn'] || 'Yes, delete';
            window.confirmDialog(
                title,
                confirmMsg,
                function () { e.detail.issueRequest(); },
                {
                    confirmButtonText: btnText,
                    confirmButtonColor: '#ef4444',
                }
            );
        }
    });

    // ---------- Re-apply i18n after HTMX swaps ----------
    document.body.addEventListener('htmx:afterSettle', function () {
        if (adminI18nCache[adminCurrentLang]) {
            adminApplyTranslations(adminI18nCache[adminCurrentLang]);
        } else {
            adminSetLanguage(adminCurrentLang);
        }
    });

    // ---------- Datatable Controls ----------
    var debounceTimers = {};

    function getListParams(entity) {
        var search = document.getElementById(entity + '-search');
        var perPage = document.getElementById(entity + '-per-page');
        var target = document.getElementById(entity + '-list');
        var sortBy = target ? (target.getAttribute('data-sort-by') || '') : '';
        var sortDir = target ? (target.getAttribute('data-sort-dir') || 'ASC') : 'ASC';
        return {
            search: search ? search.value : '',
            per_page: perPage ? perPage.value : '5',
            page: '1',
            sort_by: sortBy,
            sort_dir: sortDir
        };
    }

    window.reloadList = function (entity) {
        var params = getListParams(entity);
        var target = document.getElementById(entity + '-list');
        if (!target) return;
        var baseUrl = target.getAttribute('data-list-url') || target.getAttribute('hx-get');
        if (!baseUrl) return;
        baseUrl = baseUrl.split('?')[0];
        var qs = '?page=' + params.page +
                 '&per_page=' + params.per_page +
                 '&search=' + encodeURIComponent(params.search);
        if (params.sort_by) {
            qs += '&sort_by=' + params.sort_by + '&sort_dir=' + params.sort_dir;
        }
        htmx.ajax('GET', baseUrl + qs, { target: '#' + entity + '-list', swap: 'innerHTML' });
    };

    window.debounceReloadList = function (entity) {
        if (debounceTimers[entity]) clearTimeout(debounceTimers[entity]);
        debounceTimers[entity] = setTimeout(function () {
            window.reloadList(entity);
        }, 300);
    };

    window.sortList = function (entity, column) {
        var target = document.getElementById(entity + '-list');
        if (!target) return;
        var currentSort = target.getAttribute('data-sort-by') || '';
        var currentDir = target.getAttribute('data-sort-dir') || 'ASC';
        var newDir = (column === currentSort && currentDir === 'ASC') ? 'DESC' : 'ASC';
        target.setAttribute('data-sort-by', column);
        target.setAttribute('data-sort-dir', newDir);
        window.reloadList(entity);
    };

    window.loadPage = function (entity, page) {
        var params = getListParams(entity);
        var target = document.getElementById(entity + '-list');
        if (!target) return;
        var baseUrl = target.getAttribute('data-list-url') || target.getAttribute('hx-get');
        if (!baseUrl) return;
        baseUrl = baseUrl.split('?')[0];
        var qs = '?page=' + page +
                 '&per_page=' + params.per_page +
                 '&search=' + encodeURIComponent(params.search);
        if (params.sort_by) {
            qs += '&sort_by=' + params.sort_by + '&sort_dir=' + params.sort_dir;
        }
        htmx.ajax('GET', baseUrl + qs, { target: '#' + entity + '-list', swap: 'innerHTML' });
    };

});
