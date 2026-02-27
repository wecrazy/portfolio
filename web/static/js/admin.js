/* ===================================================
   ADMIN JS -- Theme, i18n, Toasts, Sidebar, Datatables
   =================================================== */

document.addEventListener('DOMContentLoaded', function () {

    // ---------- Toast System ----------
    function showToast(message, type) {
        var container = document.querySelector('.toast-container');
        if (!container) {
            container = document.createElement('div');
            container.className = 'toast-container';
            document.body.appendChild(container);
        }

        var toast = document.createElement('div');
        toast.className = 'toast toast-' + type + ' show';
        toast.setAttribute('role', 'alert');
        toast.innerHTML =
            '<div class="toast-body d-flex align-items-center justify-content-between py-3 px-4">' +
                '<span>' + message + '</span>' +
                '<button type="button" class="btn-close btn-close-white ms-3" aria-label="Close"></button>' +
            '</div>';

        container.appendChild(toast);

        var closeBtn = toast.querySelector('.btn-close');
        if (closeBtn) {
            closeBtn.addEventListener('click', function () {
                removeToast(toast);
            });
        }

        setTimeout(function () {
            removeToast(toast);
        }, 4000);
    }

    function removeToast(toast) {
        if (!toast || !toast.parentNode) return;
        toast.style.opacity = '0';
        toast.style.transform = 'translateX(100%)';
        toast.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
        setTimeout(function () {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 300);
    }

    // ---------- HTMX Toast Events ----------
    document.body.addEventListener('showToast', function (e) {
        var detail = e.detail || {};
        var message = detail.message || 'Action completed';
        var type = detail.type || 'success';
        showToast(message, type);
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
            langLabel.textContent = adminCurrentLang.toUpperCase();
        }
        document.querySelectorAll('.admin-lang-option').forEach(function (opt) {
            opt.classList.toggle('active', opt.getAttribute('data-lang') === adminCurrentLang);
        });
        document.documentElement.setAttribute('lang', adminCurrentLang);
    }

    function adminSetLanguage(lang) {
        adminCurrentLang = lang;
        localStorage.setItem('admin-lang', lang);
        if (adminI18nCache[lang]) {
            adminApplyTranslations(adminI18nCache[lang]);
            return;
        }
        var ver = document.documentElement.getAttribute('data-app-version') || '';
        var langUrl = '/static/lang/' + lang + '.json' + (ver ? '?v=' + ver : '');
        fetch(langUrl)
            .then(function (r) { return r.json(); })
            .then(function (dict) {
                adminI18nCache[lang] = dict;
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
            // Use i18n if available
            if (adminI18nCache[adminCurrentLang]) {
                var i18nKey = 'admin.toast.' + toastKey;
                if (adminI18nCache[adminCurrentLang][i18nKey]) {
                    msg = adminI18nCache[adminCurrentLang][i18nKey];
                }
            }
            setTimeout(function () { showToast(msg, 'success'); }, 100);
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
        icon.className = theme === 'light' ? 'bi bi-moon-fill' : 'bi bi-sun-fill';
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

    // ---------- Sidebar Toggle (Mobile) ----------
    var sidebarToggle = document.getElementById('sidebar-toggle');
    var sidebar = document.getElementById('admin-sidebar');
    var overlay = document.getElementById('sidebar-overlay');

    if (sidebarToggle && sidebar) {
        sidebarToggle.addEventListener('click', function () {
            sidebar.classList.toggle('show');
            if (overlay) {
                overlay.classList.toggle('show');
            }
        });

        if (overlay) {
            overlay.addEventListener('click', function () {
                sidebar.classList.remove('show');
                overlay.classList.remove('show');
            });
        }
    }

    // ---------- Confirm Modal Integration ----------
    document.body.addEventListener('htmx:confirm', function (e) {
        var trigger = e.detail.elt;
        var confirmMsg = trigger.getAttribute('data-confirm');
        if (confirmMsg) {
            e.preventDefault();
            if (window.confirm(confirmMsg)) {
                e.detail.issueRequest();
            }
        }
    });

    // ---------- Re-apply i18n after HTMX swaps ----------
    document.body.addEventListener('htmx:afterSettle', function () {
        if (adminI18nCache[adminCurrentLang]) {
            // Cache warm — apply instantly
            adminApplyTranslations(adminI18nCache[adminCurrentLang]);
        } else {
            // Lang JSON still in-flight; re-trigger so translations apply once fetch resolves
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
