import React, { useState, useEffect } from 'react';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  const [containers, setContainers] = useState([]);
  const [ip, setIp] = useState('');

  const apiUrl = 'http://localhost:8081';

  const fetchContainerData = async () => {
    try {
      const response = await axios.get(`${apiUrl}/containers/get`);
      setContainers(response.data || []);
    } catch (error) {
      console.error('Ошибка при получении данных:', error);
    }
  };

  const handleAddContainer = async (event) => {
    event.preventDefault();
    if (!ip) {
      alert('Пожалуйста, заполните IP адрес');
      return;
    }

    try {
      await axios.post(`${apiUrl}/containers/add`, {
        ip: ip,
      });

      await fetchContainerData();

      setIp('');
    } catch (error) {
      console.error('Ошибка при добавлении контейнера:', error);
    }
  };

  const updateContainerStatus = async (container) => {
    try {
      await axios.post(`${apiUrl}/containers/update`, {
        ip: container.ip,
        status: container.status,
      });
    } catch (error) {
      console.error('Ошибка при обновлении статуса контейнера:', error);
    }
  };

  const handleManualUpdate = async () => {
    try {
      const response = await axios.get(`${apiUrl}/containers/get`);
      const containersData = response.data;

      for (let container of containersData) {
        const status = await axios.post(`${apiUrl}/containers/update`, {
          ip: container.ip,
          status: container.status,
        });
      }

      await fetchContainerData();
    } catch (error) {
      console.error('Ошибка при ручном обновлении:', error);
    }
  };

  useEffect(() => {
    fetchContainerData();

    const interval = setInterval(fetchContainerData, 5000);

    return () => clearInterval(interval);
  }, []);

  return (
      <div className="container mt-5">
        <h1>Мониторинг контейнеров</h1>

        <form onSubmit={handleAddContainer} className="mb-4">
          <div className="form-group">
            <label htmlFor="ip">IP адрес контейнера</label>
            <input
                type="text"
                id="ip"
                className="form-control"
                value={ip}
                onChange={(e) => setIp(e.target.value)}
                placeholder="Введите IP адрес"
            />
          </div>
          <button type="submit" className="btn btn-primary mt-3">Добавить контейнер</button>
        </form>

        <button onClick={handleManualUpdate} className="btn btn-secondary mb-4">
          Обновить
        </button>

        <table className="table table-striped">
          <thead>
          <tr>
            <th scope="col">IP адрес</th>
            <th scope="col">Статус</th>
          </tr>
          </thead>
          <tbody>
          {Array.isArray(containers) && containers.length > 0 ? (
              containers.map((container, index) => (
                  <tr key={index}>
                    <td>{container.ip}</td>
                    <td>{container.status}</td>
                    <td>{container.lastChecked}</td>
                  </tr>
              ))
          ) : (
              <tr>
                <td colSpan="3">Нет данных для отображения</td>
              </tr>
          )}
          </tbody>
        </table>
      </div>
  );
}

export default App;
