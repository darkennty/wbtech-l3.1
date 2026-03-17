async function createNotification(e) {
    e.preventDefault();

    const dateVal = document.getElementById('send_date').value;
    const timeVal = document.getElementById('send_time').value;

    if (!dateVal || !timeVal) {
        alert('Please select both date and time');
        return;
    }

    const iso = new Date(`${dateVal}T${timeVal}`).toISOString();

    const payload = {
        channel: document.getElementById('channel').value,
        recipient: document.getElementById('recipient') ? document.getElementById('recipient').value : '',
        message: document.getElementById('message').value,
        send_at: iso,
    };
    const res = await fetch('/notify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
    });
    if (!res.ok) {
        const text = await res.text();
        alert('Error creating notification: ' + (text || res.status));
        return;
    }
    const data = await res.json();
    alert('Created notification with id: ' + data.result.id);
    await loadNotifications();
}

async function loadNotifications() {
    const res = await fetch('/notify');
    if (!res.ok) {
        return;
    }

    const data = await res.json();
    const items = data.result.notifications;

    const tbody = document.getElementById('notifications-body');
    tbody.innerHTML = '';

    items.forEach((n) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${n.id}</td>
            <td>${n.channel}</td>
            <td>${n.recipient}</td>
            <td>${n.payload}</td>
            <td>${n.scheduled_at}</td>
            <td>${n.status}</td>
            <td>${n.retry_count}</td>
            <td>
            ${n.status === 'scheduled' || n.status === 'pending'
              ? '<button data-id="' + n.id + '">Cancel</button>'
              : '<button data-id="' + n.id + '">Delete</button>'}
            </td>
        `;

        const btn = tr.querySelector('button');
        if (btn) {
            btn.addEventListener('click', async () => {
            await fetch('/notify/' + n.id, { method: 'PUT' });
            await loadNotifications();  
            });
        }

        tbody.appendChild(tr);
    });
}

document.getElementById('create-form').addEventListener('submit', createNotification);
document.getElementById('refresh').addEventListener('click', loadNotifications);

setInterval(loadNotifications, 5000);
loadNotifications();

