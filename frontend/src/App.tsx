import { useEffect, useState } from "react";
import "./App.css";
import { api, models } from "../wailsjs/go/models";
import { EventsEmit, EventsOn } from "../wailsjs/runtime";
import CommandPanel from "./components/CommandPanel";
import Terminal from "./components/Terminal";
import Portal from "./components/Portal";
import Viewer from "./components/Viewer";
import StatusBar from "./components/StatusBar";

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
      <StatusBar Components={[]} />
      <div
        id="content"
        className={`h-full flex flex-row ${cmdPanelShow ? "blur-sm" : ""}`}
      >
        <div className="w-1/3">
          <Terminal />
        </div>
        <div className="w-2/3">
          {appState === 0 ? (
            <Portal connectDoneCallback={connectDone} />
          ) : (
            <Viewer Code={code} />
          )}
        </div>
      </div>
      <div
        id="command-panel"
        className={`${
          cmdPanelShow ? "" : "hidden"
        } w-screen fixed top-0 left-0`}
      >
        <CommandPanel
          IsShow={cmdPanelShow}
          setIsShow={setCmdPanelShow}
          setCode={setCode}
        />
      </div>
    </div>
  );
}

export default App;
