import { useEffect, useState } from "react";
import "./App.css";
import { Connect, FetchStatus } from "../wailsjs/go/api/App";
import { api } from "../wailsjs/go/models";
import { EventsEmit, EventsOn } from "../wailsjs/runtime";
import CommandPanel from "./components/CommandPanel";

function App() {
  const [msg, setMsg] = useState("better tdx");
  const updateMsg = (msg: string) => {
    setMsg(msg);
  };
  const [isConnected, setIsConnected] = useState(false);
  const [serverMsg, setServerMsg] = useState("disconnected");

  const [allStockMarket, setAllStockMarket] = useState<api.StockMeta[]>([]);

  const updateIndexInfo = (result: api.IndexInfo) => {
    setMsg(result.Msg);
    setAllStockMarket(result.StockList);
  };

  function connect() {
    Connect("").then((msg) => {
      setServerMsg(msg);
    });
  }

  useEffect(() => {
    EventsOn(api.MsgKey.processMsg, (msg: api.ProcessInfo) => {
      updateMsg(msg.Msg);
    });
    EventsOn(api.MsgKey.connectionStatus, (connectionStatus: number) => {
      setIsConnected(connectionStatus > 0);
    });
  }, []);

  function fetchStatus() {
    FetchStatus().then(updateIndexInfo);
  }

  // command panel
  const [cmdPanelShow, setCmdPanelShow] = useState(false);

  return (
    <div id="App" className="container bg-gray-900 h-screen">
      <div id="content" className={`${cmdPanelShow ? "blur-sm" : ""}`}>
        <div id="status-bar" className="flex mb-2 bg-gray-700 space-x-2">
          <div
            className={`px-2 w-1/5 ${
              isConnected ? "bg-green-900" : "bg-red-900"
            }`}
          >
            <span>{isConnected ? "connected" : "disconnected"}</span>
          </div>
          <div className="w-1/5">
            <span>{serverMsg}</span>
          </div>
          <div className="w-3/5 text-left">
            <span>{msg}</span>
          </div>
        </div>
        <div>
          <div className="space-x-4">
            <button className="action-btn" onClick={connect}>
              Connect
            </button>
            <button className="action-btn" onClick={fetchStatus}>
              FetchStatus
            </button>
          </div>
        </div>
      </div>
      <div
        id="command-panel"
        className={`${
          cmdPanelShow ? "" : "hidden"
        } w-screen fixed top-0 left-0`}
      >
        <CommandPanel setIsShow={setCmdPanelShow} />
      </div>
    </div>
  );
}

export default App;
