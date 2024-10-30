import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { FetchStatus } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

function App() {
  // connection status
  const [connectionStatus, setConnectionStatus] = useState("Disconnected");
  const updateConnectionStatus = (result: main.Status) => {
    setConnectionStatus(result.Msg);
  };

  function tryConnect() {
    FetchStatus().then(updateConnectionStatus);
    console.log("tryConnect");
  }

  return (
    <div id="App">
      <div id="status-bar">
        <button className="btn" onClick={tryConnect}>
          Connect
        </button>
        <p id="status-bar-connection">connection: {connectionStatus}</p>
      </div>
    </div>
  );
}

export default App;
