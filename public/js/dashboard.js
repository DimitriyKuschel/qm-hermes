const CHART_UPDATE_INTERVAL = 1000;
const MAX_CHART_LABELS = 60;

// Function to fetch data from API
async function fetchData(url) {
  const response = await fetch(url);
  return response.json();
}

// Function to update counts and detailed queues
function updateCountsAndQueues(data) {
  const queueCount = document.getElementById('queueCount');
  const clientsCount = document.getElementById('clientsCount');
  const messageCount = document.getElementById('messageCount');
  const queuesTable = document.getElementById('queuesTable');

  queueCount.innerHTML = data.queueCount;
  clientsCount.innerHTML = data.clientsCount;
  messageCount.innerHTML = data.messageCount;

  queuesTable.innerHTML = '';

  data.detailedQueues.forEach((item, i) => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${i + 1}</td><td>${item.queue_name}</td><td>${item.queue_size}</td><td><a href="/queue?queue_name=${item.queue_name}">View</a></td><td><a href="/queue/delete?delete=${item.queue_name}" class="delete" data-name="${item.queue_name}">Delete</a> </td>`;
    queuesTable.appendChild(tr);
  });
}

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

// Function to create and initialize the chart
function createChart() {
  const ctx = document.getElementById('chart');
  return new Chart(ctx, {
    type: 'line',
    data: {
      labels: [],
      datasets: [
        {
          data: [],
          label: 'Messages',
          lineTension: 0,
          backgroundColor: 'transparent',
          borderColor: '#50CD89',
          borderWidth: 4,
          pointBackgroundColor: '#50CD89'
        },
        {
          data: [],
          label: 'Queues',
          lineTension: 0,
          backgroundColor: 'transparent',
          borderColor: '#3E97FF',
          borderWidth: 4,
          pointBackgroundColor: '#3E97FF'
        },
        {
          data: [],
          label: 'Subscribers',
          lineTension: 0,
          backgroundColor: 'transparent',
          borderColor: '#F1416C',
          borderWidth: 4,
          pointBackgroundColor: '#F1416C'
        }
      ]
    },
    options: {
      plugins: {
        legend: {
          display: true
        },
        tooltip: {
          boxPadding: 3
        }
      }
    }
  });
}

// Main code
const chart = createChart();

setInterval(async () => {
  try {
    const queueCount = fetchData('/api/dashboard/queues');
    const clientsCount = fetchData('/api/dashboard/clients');
    const messageCount = fetchData('/api/dashboard/messages');
    const detailedQueues = fetchData('/api/dashboard/queues/detailed');

    const data = {
      queueCount: (await queueCount).length,
      clientsCount: (await clientsCount).count,
      messageCount: (await messageCount).count,
      detailedQueues: await detailedQueues
    };

    updateCountsAndQueues(data);

    const labels = chart.data.labels;
    const datasets = chart.data.datasets;
    const dataValues = [data.messageCount, data.queueCount, data.clientsCount];

    if (labels.length > MAX_CHART_LABELS) {
      labels.shift();
      datasets.forEach(dataset => dataset.data.shift());
    }

    labels.push(new Date(Date.now()).toLocaleTimeString());
    dataValues.forEach((value, index) => datasets[index].data.push(value));

    chart.update();

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

  } catch (error) {
    console.error(error);
  }
}, CHART_UPDATE_INTERVAL);

