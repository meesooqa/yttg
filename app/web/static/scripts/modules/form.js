document.addEventListener('DOMContentLoaded', function() {
    const form = document.forms[0];
    if (form) {
        const firstTextInput = form.querySelector(
            'input[type="url"]'
        );
        if (firstTextInput) {
            firstTextInput.select();
        }
    }
});
