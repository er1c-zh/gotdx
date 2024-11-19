import { LogInfo } from "../../wailsjs/runtime/runtime";
import "../App.css";
import { useEffect, useRef, useState } from "react";

interface CommandPanelProps {
  setIsShow: React.Dispatch<React.SetStateAction<boolean>>;
  IsShow: boolean;
  setCode: React.Dispatch<React.SetStateAction<string>>;
}

function CommandPanel(props: CommandPanelProps) {
  const [cmd, setCmd] = useState("");
  const [show, setShow] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const inputHandler = (e: KeyboardEvent) => {
    LogInfo(e.key);
    LogInfo(`${props.IsShow}`);
    if (props.IsShow && e.key === "Escape") {
      props.setIsShow(false);
      setCmd("");
      e.preventDefault();
    } else if (!props.IsShow && /^[0-9a-zA-Z]+$/.test(e.key)) {
      props.setIsShow(true);
    } else if (props.IsShow && e.key === "Enter") {
      props.setIsShow(false);
      props.setCode(cmd);
      setCmd("");
      e.preventDefault();
    }
  };
  useEffect(() => {
    inputRef.current?.focus();
    document.addEventListener("keydown", inputHandler);
    return () => {
      document.removeEventListener("keydown", inputHandler);
    };
  });
  return (
    <div id="command-panel-root" className="container w-full h-screen flex">
      <div className="w-1/6"></div>
      <div className="w-2/3 bg-gray-700 mt-36 h-fit p-4 rounded border-gray-600 border-4">
        <input
          value={cmd}
          ref={inputRef}
          className="w-full text-4xl bg-gray-800 p-8 text-left overflow-x-auto"
          autoComplete="off"
          autoCorrect="off"
          autoCapitalize="off"
          onChange={(e) => {
            setCmd(e.target.value);
          }}
        ></input>
      </div>
      <div className="w-1/6"></div>
    </div>
  );
}

export default CommandPanel;
