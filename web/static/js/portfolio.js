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
        if (langLabel) {
            // Always show the short 2-letter code in the button to avoid overflow.
            langLabel.textContent = currentLang.toUpperCase();
        }
        // Mark active option in dropdown.
        document.querySelectorAll('.lang-option').forEach(function (opt) {
            opt.classList.toggle('active', opt.getAttribute('data-lang') === currentLang);
        });
        // Update html lang attribute.
        document.documentElement.setAttribute('lang', currentLang);
    }

    // Translate elements that carry data-translate (DB-driven content).
    // Originals are stashed in data-translate-orig so we can restore on lang revert.
    function translateDynamicContent(lang) {
        var elements = Array.from(document.querySelectorAll('[data-translate]'));
        if (!elements.length) return;

        // Stash originals on first encounter.
        elements.forEach(function (el) {
            if (!el.hasAttribute('data-translate-orig')) {
                el.setAttribute('data-translate-orig', el.textContent.trim());
            }
        });

        // Revert to original when switching back to the content language (EN).
        // Use contentLang here, NOT defaultLang — the UI default preference may
        // differ from the language the DB content is actually stored in.
        if (lang === contentLang) {
            elements.forEach(function (el) {
                el.textContent = el.getAttribute('data-translate-orig');
            });
            return;
        }

        if (!translateCache[lang]) translateCache[lang] = {};

        // Split into already-cached and needs-fetch.
        var toFetch = [];
        var toFetchEls = [];
        elements.forEach(function (el) {
            var orig = el.getAttribute('data-translate-orig');
            if (!orig) return; // skip empty content
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
            // Always translate FROM contentLang (the DB storage language, 'en'),
            // not from defaultLang (the UI default), which may be different.
            body: JSON.stringify({ texts: toFetch, from: contentLang, to: lang })
        })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                if (!Array.isArray(data.translations)) return;
                // Guard: user may have switched language while this fetch was in-flight.
                if (lang !== currentLang) return;
                data.translations.forEach(function (translated, idx) {
                    var orig = toFetch[idx];
                    // Only cache & apply when MyMemory actually translated something.
                    // If it echoed back the original (API error / rate-limit), skip
                    // caching so the next language switch will retry the call.
                    if (translated && translated.trim() !== orig) {
                        translateCache[lang][orig] = translated;
                        persistTranslateCache();
                        toFetchEls[idx].textContent = translated;
                    }
                });
            })
            .catch(function () {
                // Silently fall back; original text stays.
            });
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
        swapThemeIcons(document.documentElement.getAttribute('data-bs-theme') || 'dark');
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
    if (themeToggle) {
        var savedTheme = localStorage.getItem('portfolio-theme');
        if (savedTheme) {
            document.documentElement.setAttribute('data-bs-theme', savedTheme);
            updateThemeIcon(savedTheme);
        } else {
            updateThemeIcon('dark');
        }

        themeToggle.addEventListener('click', function () {
            var current = document.documentElement.getAttribute('data-bs-theme');
            var next = current === 'dark' ? 'light' : 'dark';
            document.documentElement.setAttribute('data-bs-theme', next);
            localStorage.setItem('portfolio-theme', next);
            updateThemeIcon(next);
            swapThemeIcons(next);
        });
    }

    function updateThemeIcon(theme) {
        var toggle = document.getElementById('themeToggle');
        if (!toggle) return;
        var icon = toggle.querySelector('i');
        if (!icon) return;
        if (theme === 'light') {
            icon.className = 'bxf bx-moon';
        } else {
            icon.className = 'bxf bx-sun';
        }
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

});
