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

async function loadTransaction() {
    const transactionID = window.location.pathname.split('/').pop(); 
    console.log('Transaction ID:', transactionID);

    const token = getJWTToken();

    if (!token) {
        alert('Пожалуйста, войдите в систему');
        return;
    }

    try {
        const response = await fetch(`http://localhost:8081/api/transactions/${transactionID}/data`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const result = await response.json();
            const transaction = result.transaction;

            if (transaction) {
                document.getElementById('transaction-id').value = transaction._id;
                document.getElementById('amount').value = transaction.amount;
                document.getElementById('description').value = transaction.description;
                document.getElementById('category').value = transaction.category;
            } else {
                alert('Не удалось найти данные транзакции');
            }
        } else {
            alert('Не удалось загрузить данные транзакции');
        }
    } catch (error) {
        console.error('Ошибка при загрузке транзакции:', error);
        alert('Ошибка при загрузке данных');
    }
}

async function editTransaction(event) {
    const transactionID = window.location.pathname.split('/').pop(); 
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
        category,
    };

    try {
        const res = await fetch(`http://localhost:8081/api/transactions/${transactionID}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(transactionData),
        });

        const result = await res.json();
        const messageElement = document.getElementById('message');

        if (res.ok) {
            messageElement.textContent = 'Транзакция успешно отредактирована!';
            messageElement.style.color = 'green';
            setTimeout(() => window.location.href = 'http://localhost:8081/api/transactions', 2000); 
        } else {
            messageElement.textContent = result.error || 'Ошибка при редактировании транзакции';
            messageElement.style.color = 'red';
        }
    } catch (error) {
        console.error('Ошибка при редактировании транзакции:', error);
        const messageElement = document.getElementById('message');
        messageElement.textContent = 'Ошибка при редактировании транзакции';
        messageElement.style.color = 'red';
    }
}


window.addEventListener('load', loadTransaction);

document.getElementById('edit-transaction-form').addEventListener('submit', editTransaction);