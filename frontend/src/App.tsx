import { useEffect, useState } from "react";
import "./App.css";
import CommandPanel from "./components/CommandPanel";
import Terminal from "./components/Terminal";
import Portal from "./components/Portal";
import Viewer from "./components/Viewer";
import StatusBar from "./components/StatusBar";
import KeyMessage from "./components/KeyMessage";

function App() {
  const [appState, setAppState] = useState(Number);
  // command panel
  const [cmdPanelShow, setCmdPanelShow] = useState(false);
  const [code, setCode] = useState("");
  const [showTerminal, setShowTerminal] = useState(false);
  const connectDone = () => {
    setAppState(1);
  };
  const terminalHandler = (e: KeyboardEvent) => {
    if (e.key === " ") {
      setShowTerminal(!showTerminal);
      e.preventDefault();
    }
  };
  useEffect(() => {
    document.addEventListener("keydown", terminalHandler);
    return () => {
      document.removeEventListener("keydown", terminalHandler);
    };
  });

  return (
    <div id="App" className="container bg-gray-900 h-dvh">
      <div className="flex flex-col h-full">
        <StatusBar Components={appState === 0 ? [] : [KeyMessage]} />
        <div
          id="content"
          className={`h-full flex flex-row ${cmdPanelShow ? "blur-sm" : ""}`}
        >
          <div className="flex w-full">
            {appState === 0 ? (
              <Portal connectDoneCallback={connectDone} />
            ) : (
              <Viewer Code={code} />
            )}
          </div>
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
      <div
        className={`fixed top-0 left-0 w-screen h-full opacity-75 ${
          showTerminal ? "" : "hidden"
        }`}
      >
        <Terminal />
      </div>
    </div>
  );
}

export default App;
