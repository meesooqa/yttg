async function fetchStatuses() {
    try {
        let response = await fetch('/status');
        if (response.ok) {
console.log("OK");
            let data = await response.json();
            updateTable(data);
        } else {
console.log("not OK");
        }
    } catch (err) {
        console.error('Error while getting statuses:', err);
    }
}

function updateTable(data) {
    let tbody = document.getElementById('jobsBody');
    tbody.innerHTML = '';
    for (let id in data) {
        let row = document.createElement('tr');
        let cellID = document.createElement('td');
        let cellStatus = document.createElement('td');

        cellID.textContent = id;
        cellStatus.textContent = data[id];

        row.appendChild(cellID);
        row.appendChild(cellStatus);
        tbody.appendChild(row);
    }
}

setInterval(fetchStatuses, 1000);
window.onload = fetchStatuses;
