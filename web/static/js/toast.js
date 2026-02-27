/**
 * toast.js — global SweetAlert2-powered toast & confirm utilities.
 *
 * Exposes:
 *   window.showToast(msg, type)
 *     type: 'success' | 'info' | 'warning' | 'danger'
 *
 *   window.confirmDialog(title, text, callback)
 *     callback() is called only when the user confirms.
 *
 * Requires SweetAlert2 to be loaded first.
 */
(function () {
  'use strict';

  // ── Shared Swal mixin for toasts ─────────────────────────────────────────────
  var Toast = Swal.mixin({
    toast: true,
    position: 'bottom-end',
    showConfirmButton: false,
    showCloseButton: true,
    timer: 4500,
    timerProgressBar: true,
    width: '360px',
    padding: '0.85em 1.2em',
    didOpen: function (popup) {
      popup.addEventListener('mouseenter', Swal.stopTimer);
      popup.addEventListener('mouseleave', Swal.resumeTimer);
    },
  });

  // icon mapping: Swal2 supports 'success','error','warning','info','question'
  var iconMap = {
    success: 'success',
    info:    'info',
    warning: 'warning',
    danger:  'error',
    error:   'error',
  };

  /**
   * Show a toast notification.
   * @param {string} msg  - Message text
   * @param {string} type - 'success' | 'info' | 'warning' | 'danger'
   */
  window.showToast = function (msg, type) {
    Toast.fire({
      icon: iconMap[type] || 'info',
      title: msg,
    });
  };

  /**
   * Show a Swal2 confirmation dialog (replaces window.confirm).
   * @param {string}   title    - Bold heading
   * @param {string}   text     - Subtext / message
   * @param {Function} onConfirm - Called when user clicks "Confirm"
   * @param {Object}   [opts]   - Optional overrides merged into Swal.fire config
   */
  window.confirmDialog = function (title, text, onConfirm, opts) {
    Swal.fire(Object.assign({
      title: title || 'Are you sure?',
      text: text || '',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Yes, proceed',
      cancelButtonText: 'Cancel',
      confirmButtonColor: '#ef4444',
      cancelButtonColor: '#6c6cff',
      reverseButtons: true,
      focusCancel: true,
    }, opts || {})).then(function (result) {
      if (result.isConfirmed && typeof onConfirm === 'function') {
        onConfirm();
      }
    });
  };
})();
