/* ===================================================
   PORTFOLIO JS -- Scroll, AOS, Skills, Theme, i18n
   =================================================== */

document.addEventListener('DOMContentLoaded', function () {

    // ---------- i18n Engine ----------
    var i18nCache = {};
    // translateCache[lang][origText] = translatedText
    // Load from localStorage so successful translations persist across page reloads
    // without re-hitting the free-tier API on every visit.
    var TRANSLATE_CACHE_KEY = 'portfolio-translate-cache-v2';
    var translateCache = (function () {
        try {
            return JSON.parse(localStorage.getItem(TRANSLATE_CACHE_KEY)) || {};
        } catch (e) {
            return {};
        }
    }());
    function persistTranslateCache() {
        try { localStorage.setItem(TRANSLATE_CACHE_KEY, JSON.stringify(translateCache)); } catch (e) { /* quota */ }
    }
    var defaultLang = document.documentElement.getAttribute('data-default-lang') || 'en';
    // contentLang = the language the DB content is stored in (always 'en').
    // This is intentionally separate from defaultLang (UI preference) so that
    // even if the site default is changed to another language, translation still
    // knows the correct source language to send to the API.
    var contentLang = document.documentElement.getAttribute('data-content-lang') || 'en';
    var currentLang = localStorage.getItem('portfolio-lang') || defaultLang;

    // Expose a global translation lookup so external scripts (events.js, toast.js)
    // can get the current language string without needing access to this closure.
    // window.t(key, fallback) returns the translated string or fallback.
    window.t = function (key, fallback) {
        var d = i18nCache[currentLang];
        return (d && d[key]) || fallback || key;
    };

    function applyTranslations(dict) {
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
        // Update the lang dropdown label with full localized language name.
        var langLabel = document.getElementById('langLabel');
        var langFlag = document.getElementById('langFlag');
        if (langLabel) {
            langLabel.textContent = currentLang.toUpperCase();
        }
        // Mark active option in dropdown.
        document.querySelectorAll('.lang-option').forEach(function (opt) {
            opt.classList.toggle('active', opt.getAttribute('data-lang') === currentLang);
        });
        if (langFlag) {
            var activeOpt = document.querySelector('.lang-option[data-lang="' + currentLang + '"]');
            langFlag.textContent = activeOpt ? (activeOpt.getAttribute('data-flag') || '🌐') : '🌐';
        }
        // Update html lang attribute.
        document.documentElement.setAttribute('lang', currentLang);
    }

    function translateElements(elements, lang, originalAttr) {
        if (!elements.length) return;

        elements.forEach(function (el) {
            if (!el.hasAttribute(originalAttr)) {
                el.setAttribute(originalAttr, el.textContent.trim());
            }
        });

        if (lang === contentLang) {
            elements.forEach(function (el) {
                el.textContent = el.getAttribute(originalAttr);
            });
            return;
        }

        if (!translateCache[lang]) translateCache[lang] = {};

        var toFetch = [];
        var toFetchEls = [];
        elements.forEach(function (el) {
            var orig = el.getAttribute(originalAttr);
            if (!orig) return;
            if (translateCache[lang][orig] !== undefined) {
                el.textContent = translateCache[lang][orig];
            } else {
                toFetch.push(orig);
                toFetchEls.push(el);
            }
        });

        if (!toFetch.length) return;

        fetch('/api/translate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ texts: toFetch, from: contentLang, to: lang })
        })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                if (!Array.isArray(data.translations)) return;
                if (lang !== currentLang) return;

                data.translations.forEach(function (translated, idx) {
                    var orig = toFetch[idx];
                    var targetEl = toFetchEls[idx];
                    if (!targetEl || !targetEl.isConnected) return;
                    if (translated && translated.trim() !== orig) {
                        translateCache[lang][orig] = translated;
                        persistTranslateCache();
                        targetEl.textContent = translated;
                    }
                });
            })
            .catch(function () {
                // Silently fall back; original text stays.
            });
    }

    // helper used by HTMX-injected fragments (error messages, comment batches)
    // to translate just that piece of DOM without touching the rest of the page.
    function applyInlineTranslations(el) {
        if (!el || !i18nCache[currentLang]) return;
        var dict = i18nCache[currentLang];
        el.querySelectorAll('[data-i18n]').forEach(function (node) {
            var key = node.getAttribute('data-i18n');
            if (dict[key]) node.textContent = dict[key];
        });
        el.querySelectorAll('[data-i18n-placeholder]').forEach(function (node) {
            var key = node.getAttribute('data-i18n-placeholder');
            if (dict[key]) node.setAttribute('placeholder', dict[key]);
        });
    }
    // expose for inline snippets in templates
    window.applyInlineTranslations = applyInlineTranslations;

    function initTooltips(scope) {
        if (typeof bootstrap === 'undefined' || !bootstrap.Tooltip) return;
        var root = scope || document;
        root.querySelectorAll('[data-bs-toggle="tooltip"]').forEach(function (el) {
            if (!bootstrap.Tooltip.getInstance(el)) {
                new bootstrap.Tooltip(el);
            }
        });
    }

    function initInteractivePanels(scope) {
        var root = scope || document;
        root.querySelectorAll('.interactive-panel').forEach(function (panel) {
            if (panel.hasAttribute('data-interactive-init')) return;
            panel.setAttribute('data-interactive-init', 'true');

            panel.addEventListener('mousemove', function (e) {
                var rect = panel.getBoundingClientRect();
                var pointerX = ((e.clientX - rect.left) / rect.width) * 100;
                var pointerY = ((e.clientY - rect.top) / rect.height) * 100;
                panel.style.setProperty('--pointer-x', pointerX + '%');
                panel.style.setProperty('--pointer-y', pointerY + '%');
            });

            panel.addEventListener('mouseleave', function () {
                panel.style.removeProperty('--pointer-x');
                panel.style.removeProperty('--pointer-y');
            });
        });
    }

    function animateCount(el) {
        if (el.getAttribute('data-counted') === 'true') return;

        var target = parseInt(el.getAttribute('data-count') || '0', 10);
        if (!Number.isFinite(target)) return;

        el.setAttribute('data-counted', 'true');
        var duration = 900;
        var start = null;

        function frame(ts) {
            if (!start) start = ts;
            var progress = Math.min((ts - start) / duration, 1);
            var eased = 1 - Math.pow(1 - progress, 3);
            el.textContent = String(Math.round(target * eased));
            if (progress < 1) {
                requestAnimationFrame(frame);
            }
        }

        requestAnimationFrame(frame);
    }

    function initCountUps() {
        var counters = document.querySelectorAll('.count-up');
        if (!counters.length) return;

        if (!('IntersectionObserver' in window)) {
            counters.forEach(function (counter) {
                animateCount(counter);
            });
            return;
        }

        var observer = new IntersectionObserver(function (entries, obs) {
            entries.forEach(function (entry) {
                if (entry.isIntersecting) {
                    animateCount(entry.target);
                    obs.unobserve(entry.target);
                }
            });
        }, { threshold: 0.5 });

        counters.forEach(function (counter) {
            if (counter.getAttribute('data-counted') !== 'true') {
                observer.observe(counter);
            }
        });
    }

    function setExperienceToggleState(button, expanded) {
        var label = button.querySelector('span');
        var icon = button.querySelector('i');
        var key = expanded ? 'common.see_less' : 'common.see_more';
        var fallback = expanded ? 'See less' : 'See more';

        if (label) {
            label.setAttribute('data-i18n', key);
            label.textContent = window.t(key, fallback);
        }
        if (icon) {
            icon.className = expanded ? 'bx bx-chevron-up' : 'bx bx-chevron-down';
        }
        button.classList.toggle('is-open', expanded);
    }

    function initExperienceToggles(scope) {
        var root = scope || document;
        root.querySelectorAll('.experience-toggle').forEach(function (button) {
            if (button.hasAttribute('data-toggle-init')) return;
            button.setAttribute('data-toggle-init', 'true');

            setExperienceToggleState(button, button.getAttribute('aria-expanded') === 'true');

            var targetSelector = button.getAttribute('data-bs-target');
            if (!targetSelector) return;

            var target = document.querySelector(targetSelector);
            if (!target) return;

            target.addEventListener('show.bs.collapse', function () {
                setExperienceToggleState(button, true);
            });
            target.addEventListener('hide.bs.collapse', function () {
                setExperienceToggleState(button, false);
            });
        });
    }

    // Translate elements that carry data-translate (DB-driven content).
    // Originals are stashed in data-translate-orig so we can restore on lang revert.
    function translateDynamicContent(lang) {
        translateElements(Array.from(document.querySelectorAll('[data-translate]')), lang, 'data-translate-orig');
    }

    function setLanguage(lang) {
        currentLang = lang;
        localStorage.setItem('portfolio-lang', lang);
        document.cookie = 'lang=' + lang + '; path=/; SameSite=Lax; max-age=31536000';
        translateDynamicContent(lang);
        if (i18nCache[lang]) {
            applyTranslations(i18nCache[lang]);
            return;
        }
        var ver = document.documentElement.getAttribute('data-app-version') || '';
        var langUrl = '/lang/' + lang + (ver ? '?v=' + ver : '');
        fetch(langUrl)
            .then(function (r) { return r.json(); })
            .then(function (dict) {
                i18nCache[lang] = dict;
                // Guard: user may have switched language while this fetch was in-flight.
                if (lang !== currentLang) return;
                applyTranslations(dict);
            })
            .catch(function () {
                // Silently fall back; default text stays.
            });
    }

    // Initial load: apply saved language.
    setLanguage(currentLang);

    // Language switcher click handler.
    document.querySelectorAll('.lang-option').forEach(function (opt) {
        opt.addEventListener('click', function (e) {
            e.preventDefault();
            setLanguage(this.getAttribute('data-lang'));
        });
    });

    // Re-apply translations after HTMX injects new content (comments, project cards, etc.).
    document.body.addEventListener('htmx:afterSettle', function () {
        if (i18nCache[currentLang]) {
            applyTranslations(i18nCache[currentLang]);
        }
        translateDynamicContent(currentLang);
        initSeeMoreButtons();
        initExperienceToggles(document);
        initTooltips(document);
        initInteractivePanels(document);
        if (typeof AOS !== 'undefined') {
            AOS.refreshHard();
        }
        var activeTheme = document.documentElement.getAttribute('data-bs-theme') || 'dark';
        swapThemeIcons(activeTheme);
        syncThemeToggleState(activeTheme);
    });

    // ---------- Scroll Progress Bar ----------
    var scrollProgress = document.getElementById('scrollProgress');
    if (scrollProgress) {
        function updateScrollProgress() {
            var scrollTop = window.scrollY;
            var docHeight = document.documentElement.scrollHeight - window.innerHeight;
            if (docHeight > 0) {
                var percent = (scrollTop / docHeight) * 100;
                scrollProgress.style.width = percent + '%';
            }
        }
        window.addEventListener('scroll', function () {
            requestAnimationFrame(updateScrollProgress);
        });
        updateScrollProgress();
    }

    // ---------- Navbar Scroll Effect ----------
    var navbar = document.getElementById('mainNav');
    if (navbar) {
        function handleNavbarScroll() {
            if (window.scrollY > 50) {
                navbar.classList.add('scrolled');
            } else {
                navbar.classList.remove('scrolled');
            }
        }
        window.addEventListener('scroll', function () {
            requestAnimationFrame(handleNavbarScroll);
        });
        handleNavbarScroll();
    }

    // ---------- Mobile detection helper ----------
    // reused by resume link interception and certificate click handler
    function isMobileBrowser() {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    }

    // ---------- Mobile Navbar Auto-Close ----------
    // Bootstrap collapses its own dropdown menus on outside click, but does NOT
    // collapse the main navbar toggle (#navContent) automatically. Fix it here.
    var navContent = document.getElementById('navContent');
    function collapseNavbar() {
        if (!navContent) return;
        if (navContent.classList.contains('show')) {
            // Use Bootstrap's Collapse API so the animation plays correctly.
            var bsCollapse = bootstrap.Collapse.getInstance(navContent);
            if (bsCollapse) {
                bsCollapse.hide();
            } else {
                new bootstrap.Collapse(navContent, { toggle: false }).hide();
            }
        }
    }

    // Close when clicking anywhere outside the navbar.
    document.addEventListener('click', function (e) {
        if (navbar && !navbar.contains(e.target)) {
            collapseNavbar();
        }
    });

    // Close when a nav anchor link is clicked (scrolls to section, menu should close).
    if (navContent) {
        navContent.querySelectorAll('a.nav-link').forEach(function (link) {
            link.addEventListener('click', function () {
                collapseNavbar();
            });
        });
    }

    // Close when the theme toggle is clicked (it lives inside the navbar).
    var themeToggleBtn = document.getElementById('themeToggle');
    if (themeToggleBtn) {
        themeToggleBtn.addEventListener('click', function () {
            collapseNavbar();
        });
    }

    // Close when a language option is selected (also inside the navbar).
    document.querySelectorAll('.lang-option').forEach(function (opt) {
        opt.addEventListener('click', function () {
            collapseNavbar();
        });
    });

    // ---------- 3D Tilt Effect ----------
    var tiltElements = document.querySelectorAll('.hero-portrait-card');
    tiltElements.forEach(function(el) {
        el.addEventListener('mousemove', function(e) {
            var rect = el.getBoundingClientRect();
            var x = e.clientX - rect.left - rect.width/2;
            var y = e.clientY - rect.top - rect.height/2;
            el.style.transform = 'perspective(1000px) rotateX(' + (-y / 28) + 'deg) rotateY(' + (x / 24) + 'deg) translateY(-4px)';
            el.style.transition = 'none';
        });
        el.addEventListener('mouseleave', function() {
            el.style.transform = '';
            el.style.transition = 'transform 0.5s ease';
        });
    });

    // ---------- AOS Init ----------
    if (typeof AOS !== 'undefined') {
        AOS.init({
            duration: 800,
            once: true,
            offset: 80,
            easing: 'ease-out-cubic'
        });
    }

    // ---------- Skill Bars Fill via IntersectionObserver ----------
    var skillsSection = document.getElementById('skills');
    if (skillsSection) {
        var skillsFilled = false;
        var skillObserver = new IntersectionObserver(function (entries) {
            entries.forEach(function (entry) {
                if (entry.isIntersecting && !skillsFilled) {
                    skillsFilled = true;
                    var bars = skillsSection.querySelectorAll('.skill-bar-fill');
                    bars.forEach(function (bar) {
                        var progress = parseInt(bar.getAttribute('data-progress') || 0, 10);
                        bar.style.width = progress + '%';
                        // Colour tier: 1-25 red, 26-50 orange, 51-75 blue, 76-89 purple-cyan, 90-100 primary gradient
                        var level;
                        if (progress <= 25) level = 'beginner';
                        else if (progress <= 50) level = 'elementary';
                        else if (progress <= 75) level = 'intermediate';
                        else if (progress <= 89) level = 'advanced';
                        else level = 'expert';
                        bar.setAttribute('data-level', level);
                    });
                }
            });
        }, { threshold: 0.2 });
        skillObserver.observe(skillsSection);
    }

    // ---------- Active Nav Link via IntersectionObserver ----------
    var sections = document.querySelectorAll('section[id]');
    var navLinks = document.querySelectorAll('#mainNav .nav-link[href^="#"]');

    if (sections.length > 0 && navLinks.length > 0) {
        var sectionObserver = new IntersectionObserver(function (entries) {
            entries.forEach(function (entry) {
                if (entry.isIntersecting) {
                    var id = entry.target.getAttribute('id');
                    navLinks.forEach(function (link) {
                        link.classList.remove('active');
                        if (link.getAttribute('href') === '#' + id) {
                            link.classList.add('active');
                        }
                    });
                }
            });
        }, {
            rootMargin: '-20% 0px -60% 0px',
            threshold: 0
        });

        sections.forEach(function (section) {
            sectionObserver.observe(section);
        });
    }

    // ---------- Theme Toggle ----------
    var themeToggle = document.getElementById('themeToggle');
    var themeToggleAnimationTimer = null;
    if (themeToggle) {
        var savedTheme = localStorage.getItem('portfolio-theme');
        if (savedTheme) {
            document.documentElement.setAttribute('data-bs-theme', savedTheme);
            syncThemeToggleState(savedTheme);
        } else {
            syncThemeToggleState('dark');
        }

        themeToggle.addEventListener('click', function () {
            var current = document.documentElement.getAttribute('data-bs-theme');
            var next = current === 'dark' ? 'light' : 'dark';
            animateThemeToggle(themeToggle, current, next);
            document.documentElement.setAttribute('data-bs-theme', next);
            localStorage.setItem('portfolio-theme', next);
            syncThemeToggleState(next);
            swapThemeIcons(next);
        });
    }

    function animateThemeToggle(toggle, currentTheme, nextTheme) {
        var thumb = toggle.querySelector('.theme-toggle-thumb');
        var track = toggle.querySelector('.theme-toggle-track');
        if (!thumb || !track) return;

        var maxX = Math.max(track.clientWidth - thumb.clientWidth - 10, 0);
        var startX = currentTheme === 'light' ? maxX : 0;
        var endX = nextTheme === 'light' ? maxX : 0;
        var bounceX = nextTheme === 'light' ? endX + 6 : Math.max(endX - 6, 0);
        var atmosphereShift = nextTheme === 'light' ? '-6px' : '6px';

        toggle.style.setProperty('--theme-toggle-start-x', startX + 'px');
        toggle.style.setProperty('--theme-toggle-bounce-x', bounceX + 'px');
        toggle.style.setProperty('--theme-toggle-end-x', endX + 'px');
        toggle.style.setProperty('--theme-toggle-atmosphere-shift', atmosphereShift);

        toggle.classList.remove('is-animating');
        void toggle.offsetWidth;
        toggle.classList.add('is-animating');

        if (themeToggleAnimationTimer) {
            window.clearTimeout(themeToggleAnimationTimer);
        }
        themeToggleAnimationTimer = window.setTimeout(function () {
            toggle.classList.remove('is-animating');
        }, 650);
    }

    function syncThemeToggleState(theme) {
        var toggle = document.getElementById('themeToggle');
        if (!toggle) return;
        toggle.setAttribute('data-theme', theme);
        toggle.setAttribute('aria-pressed', theme === 'light' ? 'true' : 'false');
        toggle.setAttribute('title', theme === 'light' ? 'Switch to dark theme' : 'Switch to light theme');
    }

    // ---------- Theme-aware Icon Swapping ----------
    function swapThemeIcons(theme) {
        document.querySelectorAll('.theme-icon').forEach(function (img) {
            var lightSrc = img.getAttribute('data-src-light');
            var darkSrc = img.getAttribute('data-src-dark');
            if (theme === 'light' && lightSrc) {
                img.src = lightSrc;
            } else if (theme === 'dark' && darkSrc) {
                img.src = darkSrc;
            }
        });
    }

    // Apply on initial load.
    var initialTheme = document.documentElement.getAttribute('data-bs-theme') || 'dark';
    swapThemeIcons(initialTheme);
    syncThemeToggleState(initialTheme);
    initTooltips(document);
    initInteractivePanels(document);
    initCountUps();
    initExperienceToggles(document);

    // ---------- See More / See Less Toggle ----------
    function initSeeMoreButtons() {
        document.querySelectorAll('.truncate-wrapper').forEach(function (wrapper) {
            var textEl = wrapper.querySelector('.text-truncate-clamp');
            var btn = wrapper.querySelector('.see-more-btn');
            if (!textEl || !btn) return;
            // Skip if already initialized.
            if (btn.hasAttribute('data-initialized')) return;
            btn.setAttribute('data-initialized', 'true');

            // Check if text actually overflows.
            if (textEl.scrollHeight > textEl.clientHeight + 1) {
                btn.style.display = 'inline';
            } else {
                btn.style.display = 'none';
                return;
            }

            btn.addEventListener('click', function () {
                var isExpanded = textEl.classList.contains('expanded');
                textEl.classList.toggle('expanded');
                if (isExpanded) {
                    btn.setAttribute('data-i18n', 'common.see_more');
                    btn.textContent = 'See more';
                } else {
                    btn.setAttribute('data-i18n', 'common.see_less');
                    btn.textContent = 'See less';
                }
                // Re-apply i18n for the button text.
                if (i18nCache[currentLang]) {
                    var key = btn.getAttribute('data-i18n');
                    if (i18nCache[currentLang][key]) {
                        btn.textContent = i18nCache[currentLang][key];
                    }
                }
            });
        });
    }

    // Delay init slightly to allow AOS animations to settle.
    setTimeout(initSeeMoreButtons, 500);

    // ---------- Comment DOM Windowing ----------
    // After each batch-load swap, cap visible comment items at MAX_COMMENT_ITEMS.
    // Excess items are trimmed from the TOP (newest) so the user's current
    // "older" reading position is preserved. A "See newest" link reloads page 1.
    (function () {
        var MAX_COMMENT_ITEMS = 30;

        document.addEventListener('htmx:afterSwap', function (evt) {
            var list = document.getElementById('comment-list');
            if (!list) return;

            // Identify the swap kind:
            //   target === #comment-list  → initial load or new-comment prepend → skip trim
            //   target === #comment-load-more (old node) → batch load → maybe trim
            var targetId = (evt.detail && evt.detail.target) ? evt.detail.target.id : (evt.target ? evt.target.id : '');
            if (targetId === 'comment-list' || targetId === '') return;

            var items = list.querySelectorAll('.comment-item');
            if (items.length <= MAX_COMMENT_ITEMS) return;

            // Trim oldest rendered items from the top of the list
            var excess = items.length - MAX_COMMENT_ITEMS;
            for (var i = 0; i < excess; i++) {
                if (items[i]) items[i].remove();
            }

            // Reveal the windowed indicator with a "See newest" reset link
            var windowLabel = document.getElementById('comment-window-label');
            if (windowLabel) {
                windowLabel.style.display = '';
                // Apply current i18n if already loaded
                if (typeof applyInlineTranslations === 'function') {
                    // translate the window label itself
                    applyInlineTranslations(windowLabel);
                }
            }
        });

        // Delegated click for the "See newest" reset button (lives inside HTMX content)
        document.addEventListener('click', function (evt) {
            var btn = evt.target.closest('#comment-reset-btn');
            if (!btn) return;
            evt.preventDefault();
            htmx.ajax('GET', '/comments', { target: '#comment-list', swap: 'innerHTML' });
            var windowLabel = document.getElementById('comment-window-label');
            if (windowLabel) windowLabel.style.display = 'none';
        });
    }());

    // ---------- Smooth Scroll for Anchor Links ----------
    document.querySelectorAll('a[href^="#"]').forEach(function (anchor) {
        anchor.addEventListener('click', function (e) {
            var targetId = this.getAttribute('href');
            if (targetId === '#') return;
            var target = document.querySelector(targetId);
            if (target) {
                e.preventDefault();
                var navHeight = navbar ? navbar.offsetHeight : 0;
                var targetPosition = target.getBoundingClientRect().top + window.scrollY - navHeight;
                window.scrollTo({
                    top: targetPosition,
                    behavior: 'smooth'
                });

                // Close mobile nav if open
                var navCollapse = document.getElementById('navContent');
                if (navCollapse && navCollapse.classList.contains('show')) {
                    var bsCollapse = bootstrap.Collapse.getInstance(navCollapse);
                    if (bsCollapse) {
                        bsCollapse.hide();
                    }
                }
            }
        });
    });

    // ensure pdf.js worker is configured
    if (window.pdfjsLib) {
        pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.9.179/pdf.worker.min.js';
    }

    // also intercept resume button clicks and redirect mobile users
    document.querySelectorAll('a[href="/resume"]').forEach(function(link){
        link.addEventListener('click', function(e){
            if (!isMobileBrowser()) return;
            var raw = link.getAttribute('data-resume-url');
            if (raw) {
                e.preventDefault();
                // open the direct PDF link instead of the viewer on mobile
                window.open(raw, '_blank');
            }
        });
    });
    // ---------- Certificate preview click handler ----------
    document.addEventListener('click', function(e) {
        const el = e.target.closest('.cert-view');
        if (!el) return;
        e.preventDefault();
        e.stopPropagation();
        let url = el.getAttribute('data-url');
        const type = el.getAttribute('data-type');
        const title = el.getAttribute('data-title') || '';
        // console.log('cert click', {url,type,title});
        const modalElem = document.getElementById('certModal');
        if (!modalElem) return;
        // ensure body can't scroll while modal open (bootstrap should do this
        // automatically via .modal-open but we reinforce it for our fullscreen
        // case).
        document.body.style.overflow = 'hidden';
        const modal = new bootstrap.Modal(modalElem);
        const body = modalElem.querySelector('.modal-body');
        modalElem.querySelector('.modal-title').textContent = title;
            // helper to detect mobile browsers via user agent
        function isMobileBrowser() {
            return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
        }

        if (type === 'pdf') {
            // on phones we can't reliably embed a PDF; open it directly instead
            if (isMobileBrowser()) {
                window.open(url, '_blank');
                return;
            }

            // the server already proxies, so just embed whatever we get
            // height calc subtracts approx header height (56px) so we don't get
            // a scrollbar on the entire modal; blob/pdf workers handle the rest.
            body.innerHTML = '<embed src="' + url + '" type="application/pdf" ' +
                             'width="100%" style="height:calc(100vh - 56px);"></embed>' +
                             '<p class="mt-2 text-end"><a href="' + url + '" target="_blank" class="btn btn-sm btn-outline-light">Open in new tab</a></p>';
            // intercept events directly on embed in case its own scrollability is
            // limited (unzoomed PDF); this prevents wheel/touch from leaking out.
            setTimeout(function() {
                const emb = body.querySelector('embed');
                if (emb) {
                    const block = function(ev) {
                        ev.preventDefault();
                        ev.stopPropagation();
                        return false;
                    };
                    emb.addEventListener('wheel', block, {passive:false});
                    emb.addEventListener('touchmove', block, {passive:false});
                    // remove when modal closes
                    modalElem.addEventListener('hidden.bs.modal', function() {
                        emb.removeEventListener('wheel', block);
                        emb.removeEventListener('touchmove', block);
                    }, {once:true});
                }
            }, 0);
        } else {
            body.innerHTML = '<img src="' + url + '" class="img-fluid w-100 h-auto">';
        }
        modal.show();
        modalElem.addEventListener('hidden.bs.modal', function() {
            document.body.style.overflow = '';
            // remove event listeners we added below
            modalElem.removeEventListener('wheel', stopScroll, {passive: false});
            modalElem.removeEventListener('touchmove', stopScroll, {passive: false});
        }, {once: true});

        // prevent any scroll event from reaching the document while the
        // modal is open by capturing at the document level.  This avoids
        // leaks from <embed> or other deep elements that don't bubble to
        // modalElem itself.
        function stopScroll(ev) {
            ev.preventDefault();
            ev.stopPropagation();
            return false;
        }
        document.addEventListener('wheel', stopScroll, {passive: false, capture: true});
        document.addEventListener('touchmove', stopScroll, {passive: false, capture: true});
        modalElem.addEventListener('hidden.bs.modal', function() {
            document.body.style.overflow = '';
            document.removeEventListener('wheel', stopScroll, {capture: true});
            document.removeEventListener('touchmove', stopScroll, {capture: true});
        }, {once: true});
    });

    // ---------- PDF thumbnail rendering using pdf.js ----------
    function renderPdfThumb(canvas, id) {
        if (!window.pdfjsLib) return;
        const url = '/cert/preview?id=' + id;
        const loadingTask = pdfjsLib.getDocument(url);
        loadingTask.promise.then(function(pdf) {
            pdf.getPage(1).then(function(page) {
                const viewport = page.getViewport({ scale: 1 });
                const ctx = canvas.getContext('2d');
                // use parent width to account for responsive sizing
                const parent = canvas.parentElement;
                const width = (parent ? parent.clientWidth : canvas.clientWidth) || viewport.width;
                const scale = width / viewport.width;
                const scaled = page.getViewport({ scale });
                canvas.width = scaled.width;
                canvas.height = scaled.height;
                page.render({ canvasContext: ctx, viewport: scaled });
            });
        }).catch(function(err){
            console.warn('pdf thumb failed', err);
        });
    }
    function initPdfThumbs() {
        document.querySelectorAll('canvas[data-pdf-id]').forEach(function(c){
            const id = c.getAttribute('data-pdf-id');
            renderPdfThumb(c, id);
        });
    }
    // initial render
    initPdfThumbs();
    // re-render after HTMX swaps partials
    document.addEventListener('htmx:afterSwap', function(evt) {
        if (evt.detail && evt.detail.target && evt.detail.target.id === 'certificate-list-area') {
            initPdfThumbs();
        }
    });

});
