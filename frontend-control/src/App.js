import './App.css';
import axios from 'axios';
import React, { useEffect, useState } from 'react';

const textareaStyle = {
  maxWidth: '70%', // 70% of the screen width
  width: '100%', // Ensures it takes up available space within the maxWidth
  height: '200px', // Set a height if needed
  whiteSpace: 'pre-line', // Preserves newlines in the logs
  border: '1px solid #ccc', // Optional: Style the textarea border
  padding: '10px', // Optional: Add some padding inside the textarea
  fontSize: '16px', // Optional: Set the font size
};


function App() {

  const [status, setStatus] = useState('stopped');
  const [logs, setLogs] = useState("");
  const [statusColor, setStatusColor] = useState('dark_red');

  const fetchData = async () => {
    try {
      const response = await axios.get('http://192.168.86.230:31201/status');
      setStatus(response.data);

      if (response.data === "started") {
        setStatusColor('dark_green')
      } else {
        setStatusColor('dark_red')
      }

    } catch (error) {
      console.error('Error fetching data:', error)
      setStatus("error fetching data");
    }
  };

  const fetchLogs = async () => {
    try {
      const response = await axios.get('http://192.168.86.230:31201/logs');
      if (response.data !== "\n") {
        setLogs(response.data);
      }
    } catch (error) {
      console.error('Error fetching data:', error)
      setLogs("error fetching data");
    }
  };


  useEffect(() => {
    fetchData();
    fetchLogs();

    var intervalTime = 2000;
    if (status === "stopped") {
      intervalTime = 4000;
    };

    const intervalId = setInterval(() => {
      fetchData();
      fetchLogs();
    }, intervalTime);

    return () => clearInterval(intervalId);
  }, [status]);


  const startRequest = async () => {
    try {
      
      const response = await axios.post('http://192.168.86.230:31201/start')

      const result = await response.text();
      console.log(result.data);
      fetchData();
    } catch (error) {
      console.log(error);
    }
  };

  const stopRequest = async () => {
    try {
      
      const response = await axios.post('http://192.168.86.230:31201/stop')

      const result = await response.text();
      console.log(result.data);
      fetchData();
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1 className='blue'>minecraft control panel</h1>
        <h2>
          The server is: <b className={statusColor}>{status}</b>!
        </h2>
        <div className="minecraft-button-group">
          <button className="minecraft-button" onClick={startRequest}>start</button>
          <button className="minecraft-button" onClick={stopRequest}>stop</button>
        </div>

        <textarea className='dark_red' style={{whiteSpace: 'pre-line', overflowY: 'scroll' }} readOnly value={logs} style={textareaStyle}>
        </textarea>

      </header>
    </div>
  );
}

export default App;
