import "../App.css";
import { useEffect, useRef, useState } from "react";

interface CommandPanelProps {
  setIsShow: React.Dispatch<React.SetStateAction<boolean>>;
}

function CommandPanel(props: CommandPanelProps) {
  const [cmd, setCmd] = useState("");
  let isShow = false;
  const inputRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    inputRef.current?.focus();
    document.addEventListener("keyup", (e) => {
      if (isShow && e.key === "Escape") {
        isShow = false;
        props.setIsShow(false);
        setCmd("");
      } else if (!isShow && /^[0-9a-zA-Z]+$/.test(e.key)) {
        isShow = true;
        props.setIsShow(true);
      }
    });
  });
  return (
    <div id="command-panel-root" className="container w-full h-screen flex">
      <div className="w-1/6"></div>
      <div className="w-2/3 bg-gray-700 mt-36 h-fit p-4 rounded border-gray-600 border-4">
        <input
          ref={inputRef}
          type="text"
          className="w-full text-4xl bg-gray-800 p-8"
          value={cmd}
          onFocus={(e) => {
            e.target.setSelectionRange(100, 100);
          }}
          onChange={(e) => {
            if (e.target.value.length > 0) {
              setCmd(e.target.value);
            }
          }}
        />
      </div>
      <div className="w-1/6"></div>
    </div>
  );
}

export default CommandPanel;
