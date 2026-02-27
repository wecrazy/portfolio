/* ===================================================
   PORTFOLIO JS -- Scroll, AOS, Skills, Theme, i18n
   =================================================== */

document.addEventListener('DOMContentLoaded', function () {

    // ---------- i18n Engine ----------
    var i18nCache = {};
    var defaultLang = document.documentElement.getAttribute('data-default-lang') || 'en';
    var currentLang = localStorage.getItem('portfolio-lang') || defaultLang;

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
            langLabel.textContent = dict['lang.' + currentLang] || currentLang.toUpperCase();
        }
        // Mark active option in dropdown.
        document.querySelectorAll('.lang-option').forEach(function (opt) {
            opt.classList.toggle('active', opt.getAttribute('data-lang') === currentLang);
        });
        // Update html lang attribute.
        document.documentElement.setAttribute('lang', currentLang);
    }

    function setLanguage(lang) {
        currentLang = lang;
        localStorage.setItem('portfolio-lang', lang);
        document.cookie = 'lang=' + lang + '; path=/; SameSite=Lax; max-age=31536000';
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
                        var progress = bar.getAttribute('data-progress') || 0;
                        bar.style.width = progress + '%';
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
