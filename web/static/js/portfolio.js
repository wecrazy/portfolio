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
        // Update the lang dropdown label.
        var langLabel = document.getElementById('langLabel');
        if (langLabel) {
            langLabel.textContent = currentLang.toUpperCase();
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
        if (i18nCache[lang]) {
            applyTranslations(i18nCache[lang]);
            return;
        }
        fetch('/static/lang/' + lang + '.json')
            .then(function (r) { return r.json(); })
            .then(function (dict) {
                i18nCache[lang] = dict;
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
        });
    }

    function updateThemeIcon(theme) {
        var toggle = document.getElementById('themeToggle');
        if (!toggle) return;
        var icon = toggle.querySelector('i');
        if (!icon) return;
        if (theme === 'light') {
            icon.className = 'bi bi-moon-fill';
        } else {
            icon.className = 'bi bi-sun-fill';
        }
    }

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
