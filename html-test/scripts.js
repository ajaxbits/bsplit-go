/* scripts.js */
// Add any necessary JavaScript here, for example, to handle closing modals
document.addEventListener('click', function(event) {
    if (event.target.classList.contains('close')) {
        document.getElementById('modalContainer').innerHTML = '';
        document.getElementById('splitTypeModal').innerHTML = '';
    }
});
