function getJWTToken() {
    const cookies = document.cookie.split('; ');
    for (let i = 0; i < cookies.length; i++) {
        const [key, value] = cookies[i].split('=');
        if (key === 'jwtToken') {
            return value;
        }
    }
    return null; 
}

async function addTransaction(event) {
    event.preventDefault(); 

    const token = getJWTToken();
    if (!token) {
        alert('Пожалуйста, войдите в систему');
        return;
    }

    const amount = parseFloat(document.getElementById('amount').value);
    const description = document.getElementById('description').value;
    const category = document.getElementById('category').value;

    if (!amount || !description || !category) {
        alert('Пожалуйста, заполните все поля');
        return;
    }

    const transactionData = {
        amount,
        description,
        category
    };

    try {
        const res = await fetch('http://localhost:8081/api/transactions/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(transactionData)
        });

        const result = await res.json();
        const messageElement = document.getElementById('message');
        
        if (res.ok) {
            messageElement.textContent = 'Транзакция успешно добавлена!';
            messageElement.style.color = 'green';
            setTimeout(() => window.location.href = 'http://localhost:8081/api/transactions', 2000);
        } else {
            messageElement.textContent = result.error || 'Ошибка при добавлении транзакции';
            messageElement.style.color = 'red';
        }
    } catch (error) {
        console.error('Ошибка при добавлении транзакции:', error);
        const messageElement = document.getElementById('message');
        messageElement.textContent = 'Ошибка при добавлении транзакции';
        messageElement.style.color = 'red';
    }
}

document.getElementById('add-transaction-form').addEventListener('submit', addTransaction);