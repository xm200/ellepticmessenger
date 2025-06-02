document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorElement = document.getElementById('errorMessage').value;

    try {
        const response = await fetch('/login', {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({username, password})
        });
        if (response.ok) {
            window.location.href = '/home';
        } else {
            const data = await response.json();
            errorElement.textContent = data.error || 'Login failed';
        }
    } catch (err) {
        errorElement.textContent = "Network error, try again later";
    }
})