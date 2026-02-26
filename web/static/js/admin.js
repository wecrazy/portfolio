/* ===================================================
   ADMIN JS -- HTMX Toasts, Sidebar, Confirm Modal
   =================================================== */

document.addEventListener('DOMContentLoaded', function () {

    // ---------- HTMX Toast Events ----------
    // Listen for HTMX afterSettle to trigger toast notifications.
    // Backend can send HX-Trigger header with showToast event.
    document.body.addEventListener('showToast', function (e) {
        var detail = e.detail || {};
        var message = detail.message || 'Action completed';
        var type = detail.type || 'success';
        showToast(message, type);
    });

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

        // Close button handler
        var closeBtn = toast.querySelector('.btn-close');
        if (closeBtn) {
            closeBtn.addEventListener('click', function () {
                removeToast(toast);
            });
        }

        // Auto-dismiss after 4 seconds
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

    // ---------- Sidebar Active Link ----------
    var currentPath = window.location.pathname;
    var sidebarLinks = document.querySelectorAll('.admin-sidebar .nav-link');
    sidebarLinks.forEach(function (link) {
        var href = link.getAttribute('href');
        if (href === currentPath || (href !== '/admin' && currentPath.startsWith(href))) {
            link.classList.add('active');
        } else if (href === '/admin' && currentPath === '/admin') {
            link.classList.add('active');
        }
    });

    // ---------- Sidebar Toggle (Mobile) ----------
    var sidebarToggle = document.querySelector('.sidebar-toggle');
    var sidebar = document.querySelector('.admin-sidebar');
    var overlay = document.querySelector('.sidebar-overlay');

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
    // Usage: add data-confirm="Are you sure?" to any HTMX element.
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

});
