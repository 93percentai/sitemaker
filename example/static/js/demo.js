// demo.js — loaded via [INCLUDES]: # (/js/demo.js)
document.addEventListener('DOMContentLoaded', function() {
  var el = document.getElementById('js-include-demo');
  if (!el) return;

  var colors = ['#2563eb', '#059669', '#d97706', '#dc2626', '#7c3aed', '#0891b2'];
  var idx = 0;

  el.innerHTML = '<p style="font-size:1.25rem;font-weight:600;margin:0;" id="demo-text">This text is styled by demo.js</p>' +
    '<p style="font-size:0.875rem;color:#6b7280;margin:0.25rem 0 0;">Loaded via the INCLUDES directive as an external JS file.</p>' +
    '<button id="demo-btn" style="margin-top:0.75rem;padding:0.5rem 1rem;background:#2563eb;color:white;border:none;border-radius:0.25rem;cursor:pointer;">Change Color</button>';

  document.getElementById('demo-btn').addEventListener('click', function() {
    idx = (idx + 1) % colors.length;
    document.getElementById('demo-text').style.color = colors[idx];
    this.style.background = colors[idx];
  });
});
