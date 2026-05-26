// Smooth scroll for anchor links
document.addEventListener('click', function(e) {
    var link = e.target.closest('a[href^="#"]');
    if (!link) return;
    var target = document.querySelector(link.getAttribute('href'));
    if (target) {
        e.preventDefault();
        target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        history.pushState(null, '', link.getAttribute('href'));
    }
});
