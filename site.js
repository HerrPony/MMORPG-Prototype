// Version of the native client (from the update manifest next to this page), mobile menu and the screenshot lightbox.
(function () {
  var burger = document.querySelector('.burger'), nav = document.querySelector('nav.main');
  if (burger) burger.addEventListener('click', function () { var o = nav.classList.toggle('open'); burger.setAttribute('aria-expanded', o); });

  var plat = /Android|iPhone|iPad|iPod/i.test(navigator.userAgent) ? 'mobile' : /Windows/i.test(navigator.userAgent) ? 'windows' : /Linux|X11/i.test(navigator.userAgent) ? 'linux' : 'other';
  document.documentElement.setAttribute('data-platform', plat);
  document.querySelectorAll('[data-for]').forEach(function (el) { if (el.getAttribute('data-for') === plat) el.classList.add('recommended'); });

  fetch('updates/manifest-windows.json', { cache: 'no-cache' }).then(function (r) { return r.ok ? r.json() : Promise.reject(); }).then(function (m) {
    var d = m.released ? new Date(m.released).toLocaleDateString('de-DE', { day: 'numeric', month: 'long', year: 'numeric' }) : '';
    document.querySelectorAll('[data-version]').forEach(function (el) { el.textContent = 'Version ' + m.version + (d ? ' · ' + d : ''); });
    document.querySelectorAll('[data-version-num]').forEach(function (el) { el.textContent = m.version; });
  }).catch(function () {});

  var items = Array.prototype.slice.call(document.querySelectorAll('.gallery figure'));
  if (!items.length) return;
  var box = document.createElement('div'); box.className = 'lightbox'; box.setAttribute('role', 'dialog'); box.setAttribute('aria-modal', 'true');
  box.innerHTML = '<button class="close" aria-label="Schließen">×</button><button class="nav prev" aria-label="Zurück">‹</button><img alt=""><p></p><button class="nav next" aria-label="Weiter">›</button>';
  document.body.appendChild(box);
  var img = box.querySelector('img'), cap = box.querySelector('p'), cur = 0;
  function show(i) { cur = (i + items.length) % items.length; var f = items[cur]; img.src = f.querySelector('a, img').getAttribute('data-full') || f.querySelector('img').src; img.alt = f.querySelector('img').alt; cap.textContent = (f.querySelector('figcaption') || {}).textContent || ''; box.classList.add('open'); }
  items.forEach(function (f, i) { f.addEventListener('click', function () { show(i); }); });
  box.addEventListener('click', function (e) {
    if (e.target.classList.contains('prev')) show(cur - 1); else if (e.target.classList.contains('next')) show(cur + 1);
    else if (e.target === box || e.target.classList.contains('close')) box.classList.remove('open');
  });
  document.addEventListener('keydown', function (e) {
    if (!box.classList.contains('open')) return;
    if (e.key === 'Escape') box.classList.remove('open'); else if (e.key === 'ArrowLeft') show(cur - 1); else if (e.key === 'ArrowRight') show(cur + 1);
  });
})();
