import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { Connect, FetchStatus } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

function App() {
  const [msg, setMsg] = useState("...");
  const updateMsg = (msg: string) => {
    setMsg(msg);
  }
  const [isConnected, setIsConnected] = useState(false);
  const updateIsConnected = (IsConnected: boolean) => {
    setIsConnected(IsConnected);
  }
  const [stockCount, setStockCount] = useState(-1);
  const updateStockCount = (StockCount: number) => {
    setStockCount(StockCount);
  }

  const updateIndexInfo = (result: main.IndexInfo) => {
    setMsg(result.Msg);
    setIsConnected(result.IsConnected);
    setStockCount(result.StockCount);
  };

  function connect() {
    Connect("").then(updateMsg);
  }

  function fetchStatus() {
    FetchStatus().then(updateIndexInfo);
  }

  return (
    <div id="App">
      <div id="status-bar">
        <h1 className="text-3xl font-bold underline">better tdx</h1>
        <button className="btn" onClick={connect}>
          Connect
        </button>
        <button className="btn" onClick={fetchStatus}>FetchStatus</button>
        <p id="status-bar-connection">connection: {isConnected ? "connected" : "disconnected"}</p>
        <p id="status-bar-stock-cnt">sh: {stockCount}</p>
      </div>
    </div>
  );
}

export default App;
