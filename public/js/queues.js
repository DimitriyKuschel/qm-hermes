const TABLE_UPDATE_INTERVAL = 1000;

async function deleteQueue(name) {
    try {
        let response = await fetch(`/api/dashboard/queue/delete?queue_name=${name}`);
        return true;
    }
    catch (error) {
        console.error(error);
        return false;
    }
}

async function updateQueueTable() {
    const queueTable = document.getElementById('queuesTable');

    try {
        const detailedQueuesResponse = await fetch('/api/dashboard/queues/detailed');
        const jsonDetailedQueues = await detailedQueuesResponse.json();

        queueTable.innerHTML = '';

        jsonDetailedQueues.forEach((item, i) => {
            const tr = document.createElement('tr');
            tr.innerHTML = `<td>${i + 1}</td><td>${item.queue_name}</td><td>${item.queue_size}</td><td><a href="/queue?queue_name=${item.queue_name}">View</a></td><td><a href="/queue/delete?delete=${item.queue_name}" data-name="${item.queue_name}" class="delete">Delete</a> </td>`;
            queueTable.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
    }

    document.querySelectorAll('.delete').forEach((item) => {
        item.addEventListener('click', async (e) => {
            e.preventDefault();
            if (confirm(`Are you sure you want to delete ${item.dataset.name}?`)) {
                let res = await deleteQueue(item.dataset.name);
                if (res) {
                    window.location.reload();
                }
                else {
                    alert('Error deleting queue');
                }
            }
        });
    });
}

setInterval(updateQueueTable, TABLE_UPDATE_INTERVAL);