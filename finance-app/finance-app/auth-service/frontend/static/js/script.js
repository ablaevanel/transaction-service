document.getElementById("register-form")?.addEventListener("submit", async function(e) {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const response = await fetch("http://localhost:8080/api/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
    });

    const result = await response.json();

    const messageElement = document.getElementById("message");

    if (response.ok) {
        messageElement.textContent = "Регистрация прошла успешно!";
        messageElement.style.color = "green"; 

        setTimeout(function() {
            window.location.href = "http://localhost:8080/api/auth/login"; 
        }, 2000);
    } else {
        messageElement.textContent = "Ошибка регистрации";
        messageElement.style.color = "red";
    }
});

const transactionsServiceUrl = "http://localhost:8081";

document.getElementById("login-form")?.addEventListener("submit", async function(e) {
    e.preventDefault();
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("http://localhost:8080/api/auth/login", { 
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        const result = await response.json();
        const messageElement = document.getElementById("message");

        if (response.ok) {
            console.log("Response OK:", result); 
            messageElement.textContent = "Вход прошел успешно!";
            messageElement.style.color = "green";
            
            setTimeout(() => {
                console.log("Redirecting to http://localhost:8081/api/transactions");
                window.location.href = `${transactionsServiceUrl}/api/transactions`;
            }, 2000);
        } else {
            messageElement.textContent = result.error || "Ошибка входа";
            messageElement.style.color = "red";
            console.log("Login error:", result);
        }
    } catch (error) {
        console.error("Ошибка при входе:", error);
        document.getElementById("message").textContent = "Ошибка при входе";
        document.getElementById("message").style.color = "red";
    }
});

function deleteCookie(name) {
    document.cookie = name + '=; Max-Age=0; path=/; domain=localhost;';
}

async function logout() {
    try {
        deleteCookie('jwtToken');

        const messageElement = document.getElementById('message');
        messageElement.textContent = 'Вы успешно вышли из системы.';
        messageElement.style.color = 'green';

        setTimeout(() => {
            window.location.href = 'http://localhost:8080/api/auth/login';
        }, 2000);
    } catch (error) {
        console.error('Ошибка при выходе из системы:', error);
        const messageElement = document.getElementById('message');
        messageElement.textContent = 'Ошибка при выходе из системы.';
        messageElement.style.color = 'red';
    }
}
document.getElementById('logout-button').addEventListener('click', logout);
