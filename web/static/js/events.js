/**
 * events.js — WebSocket client for real-time comment updates
 * and server shutdown / restart notifications.
 *
 * Messages received from /ws/comments (JSON):
 *   {"type":"comment",  "data":{"id":1,"user":"Alice","refresh":true}}
 *   {"type":"shutdown", "data":{"message":"Server is restarting…"}}
 *
 * Requires toast.js (window.showToast) to be loaded first.
 */
(function () {
  'use strict';

  if (!window.WebSocket) return; // browser too old

  var ws;
  var reconnectDelay = 3000;

  function connect() {
    var proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(proto + '//' + location.host + '/ws/comments');

    ws.onmessage = function (e) {
      try {
        var event = JSON.parse(e.data);

        if (event.type === 'comment') {
          var d = event.data || {};
          if (d.refresh && window.htmx) {
            var list = document.getElementById('comment-list');
            if (list) htmx.trigger(list, 'refresh');
          }
          if (d.user) {
            var commentMsg = window.t
              ? window.t('ws.new_comment', '{user} left a comment').replace('{user}', d.user)
              : d.user + ' left a comment';
            window.showToast(commentMsg, 'info');
          }

        } else if (event.type === 'shutdown') {
          var msg = window.t
            ? window.t('ws.shutdown', 'Server is restarting — please wait…')
            : (event.data && event.data.message) || 'Server is restarting — please wait…';
          window.showToast(msg, 'warning');
          ws.close();
          setTimeout(function () { connect(); }, reconnectDelay * 3);
        }
      } catch (_) {}
    };

    ws.onerror = function () {};
    ws.onclose = function () {
      setTimeout(function () { connect(); }, reconnectDelay);
    };
  }

  connect();
})();
