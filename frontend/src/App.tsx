import { useEffect, useState } from "react";
import "./App.css";
import { api, models } from "../wailsjs/go/models";
import { EventsEmit, EventsOn } from "../wailsjs/runtime";
import CommandPanel from "./components/CommandPanel";
import RealtimeGraph from "./components/RealtimeGraph";
import Terminal from "./components/Terminal";
import Portal from "./components/Portal";

function App() {
  const [appState, setAppState] = useState(Number);
  // command panel
  const [cmdPanelShow, setCmdPanelShow] = useState(false);
  const [code, setCode] = useState("");
  const connectDone = () => {
    setAppState(1);
  };

  return (
    <div id="App" className="container bg-gray-900 h-dvh">
      <div
        id="content"
        className={`h-full flex flex-row ${cmdPanelShow ? "blur-sm" : ""}`}
      >
        <div className="w-1/3">
          <Terminal />
        </div>
        {appState === 0 ? (
          <div className="w-2/3">
            <Portal setState={setAppState} />
          </div>
        ) : (
          <div>
            <p>TODO</p>
          </div>
        )}
      </div>
      <div
        id="command-panel"
        className={`${
          cmdPanelShow ? "" : "hidden"
        } w-screen fixed top-0 left-0`}
      >
        <CommandPanel setIsShow={setCmdPanelShow} setCode={setCode} />
      </div>
    </div>
  );
}

export default App;
