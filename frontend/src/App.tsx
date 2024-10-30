import { useState } from "react";
import "./App.css";
import { Connect, FetchStatus } from "../wailsjs/go/api/App";
import { api } from "../wailsjs/go/models";

function App() {
  const [msg, setMsg] = useState("...");
  const updateMsg = (msg: string) => {
    setMsg(msg);
  };
  const [isConnected, setIsConnected] = useState(false);
  const [serverMsg, setServerMsg] = useState("...");

  const [allStockMarket, setAllStockMarket] = useState<api.StockMeta[]>([]);

  const updateIndexInfo = (result: api.IndexInfo) => {
    setMsg(result.Msg);
    setIsConnected(result.IsConnected);
    setAllStockMarket(result.StockList);
  };

  function connect() {
    Connect("").then((msg) => {
      setServerMsg(msg);
    });
  }

  function fetchStatus() {
    
    FetchStatus().then(updateIndexInfo);
  }

  return (
    <div id="App" className="container">
      <div id="status-bar" className="flex my-2">
        <div className={`w-1/4 ${isConnected ? "bg-green-900" : "bg-red-900"}`}>
          <span>{isConnected ? "connected" : "disconnected"}</span>
        </div>
        <div className="w-2/4">
          <span>{msg}</span>
        </div>
        <div className="w-1/4">
          <span>{serverMsg}</span>
        </div>
      </div>
      <div>
        <div className="space-x-4">
          <button className="btn" onClick={connect}>
            Connect
          </button>
          <button className="btn" onClick={fetchStatus}>
            FetchStatus
          </button>
        </div>
      </div>
      <div id="container" className="container">
        <div className="w-full">
          <p className="text-3xl font-bold underline">
            Stock Market List {allStockMarket.length}{" "}
          </p>
        </div>
        <div id="stock-list" className="w-full space-y-2">
          {allStockMarket.map((stock) => (
            <div key={stock.Code}>
              <p>
                {stock.Market} {stock.Desc} {stock.Code}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default App;
