document.getElementById("calcForm").addEventListener("submit", function(e) {
  e.preventDefault();
  var expr = document.getElementById("expression").value;
  document.getElementById("result").innerHTML = "";
  document.getElementById("status").innerHTML = "<div class='d-flex align-items-center'><strong>Отправка запроса...</strong><div class='spinner-border ms-2 text-primary' role='status'><span class='visually-hidden'>Загрузка...</span></div></div>";
  fetch("/api/v1/calculate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ expression: expr })
  })
  .then(res => res.json())
  .then(data => {
    document.getElementById("result").innerHTML = "<div class='alert alert-success'>Задача принята. ID: " + data.id + "</div>";
    pollExpressionStatus(data.id);
    updateRequestsList(); // Обновляем список запросов после добавления нового
  })
  .catch(err => {
    console.error("Ошибка отправки запроса", err);
    document.getElementById("result").innerHTML = "<div class='alert alert-danger'>Ошибка отправки запроса</div>";
    document.getElementById("status").innerHTML = "";
  });
});

function pollExpressionStatus(id) {
  var pollInterval = setInterval(function() {
    fetch("/api/v1/expressions/" + id)
      .then(res => res.json())
      .then(data => {
        if(data.status === "Выполнено.") {
          document.getElementById("status").innerHTML = "<div class='alert alert-success'>Результат: " + data.result + "</div>";
          clearInterval(pollInterval);
        } else {
          document.getElementById("status").innerHTML = "<div class='alert alert-info'>Статус: " + data.status + "</div>";
          if (data.status === "Ошибка.") {
            clearInterval(pollInterval);
          }
        }
      })
      .catch(err => {
        console.error("Ошибка получения статуса", err);
        document.getElementById("status").innerHTML = "<div class='alert alert-danger'>Ошибка получения статуса</div>";
        clearInterval(pollInterval);
      });
  }, 2000);
}

let currentRequests = new Set(); // Храним ID текущих запросов

function updateRequestsList() {
  console.log("Обновление списка запросов..."); 
  fetch("/api/v1/expressions")
    .then(res => {
      if (!res.ok) {
        throw new Error("Ошибка сети: " + res.statusText);
      }
      return res.json();
    })
    .then(data => {
      console.log("Получены данные:", data); 
      const requestsList = document.getElementById("requestsList");

      // Проходим по новым данным
      data.forEach(request => {
        if (!currentRequests.has(request.id)) {
          // Если запрос новый, добавляем его в начало списка
          const listItem = document.createElement("li");
          listItem.className = "overflow-auto list-group-item d-flex flex-column justify-content-between align-items-start";
          listItem.innerHTML = `
            <div class="overflow-auto"><strong>ID:</strong> ${request.id}</div>
            <div class="overflow-auto"><strong>Результат:</strong> ${request.result || "—"}</div>
            <div class="overflow-auto mt-2">
              <span class="badge bg-${request.status === 'Выполнено.' ? 'success' : request.status === 'Ошибка.' ? 'danger' : 'info'}">${request.status}</span>
            </div>
          `;
          requestsList.prepend(listItem); // Добавляем в начало списка
          currentRequests.add(request.id); // Запоминаем ID нового запроса
        }
      });

      // Обновляем статусы существующих запросов
      const existingItems = requestsList.querySelectorAll("li");
      existingItems.forEach(item => {
        const id = item.querySelector("div").textContent.split(": ")[1]; 
        const request = data.find(req => req.id == id); // Находим соответствующий запрос в данных
        if (request) {
          // Обновляем статус и результат
          const statusBadge = item.querySelector(".badge");
          statusBadge.className = `badge bg-${request.status === 'Выполнено.' ? 'success' : request.status === 'Ошибка.' ? 'danger' : 'info'}`;
          statusBadge.textContent = request.status;
          const resultDiv = item.querySelector("div:nth-child(2)");
          resultDiv.innerHTML = `<strong>Результат:</strong> ${request.result || "—"}`;
        }
      });

      // Проверяем, все ли запросы выполнены
      const allDone = data.every(request => request.status === 'Выполнено.' || request.status === 'Ошибка.');
      if (!allDone) {
        // Если не все выполнены, продолжаем обновление
        setTimeout(updateRequestsList, 2000); // Обновляем каждые 2 секунды
      }
    })
    .catch(err => {
      console.error("Ошибка получения списка запросов", err);
      const requestsList = document.getElementById("requestsList");
      requestsList.innerHTML = `<li class="list-group-item text-danger">Ошибка загрузки списка запросов: ${err.message}</li>`;
    });
}

updateRequestsList();

// Обновляем список запросов при загрузке страницы
document.addEventListener("DOMContentLoaded", updateRequestsList);