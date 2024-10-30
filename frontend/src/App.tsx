import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { Connect, FetchStatus } from "../wailsjs/go/main/App";
import { main } from "../wailsjs/go/models";

function App() {
  const [msg, setMsg] = useState("...");
  const updateMsg = (msg: string) => {
    setMsg(msg);
  };
  const [isConnected, setIsConnected] = useState(false);

  const [allStockMarket, setAllStockMarket] = useState<main.StockMarketList[]>(
    []
  );

  const updateIndexInfo = (result: main.IndexInfo) => {
    setMsg(result.Msg);
    setIsConnected(result.IsConnected);
    setAllStockMarket(result.AllStock);
  };

  function connect() {
    Connect("").then(updateMsg);
  }

  function fetchStatus() {
    FetchStatus().then(updateIndexInfo);
  }

  return (
    <div id="App" className="container">
      <div id="status-bar" className="container">
        <p className="gap-x-4">
          <span className="text-3xl font-bold underline">[better tdx]</span>
          <span>connection: {isConnected ? "connected" : "disconnected"}</span>
          <span></span>
          <span>msg: {msg}</span>
        </p>
        <button className="btn" onClick={connect}>
          Connect
        </button>
        <button className="btn" onClick={fetchStatus}>
          FetchStatus
        </button>
      </div>
      <div id="container" className="container">
        <div id="stock-list" className="w-1/3">
          {allStockMarket.map((item) => (
            <div>
              <p className="text-3xl font-bold underline">{item.MarketStr}</p>
              <div id={`stock-list-${item.Market}`}>
                {item.StockList.map((stock) => (
                  <p>
                    {stock.Desc} {stock.Code}
                  </p>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
