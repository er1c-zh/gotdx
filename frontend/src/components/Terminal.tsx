import { useEffect, useRef, useState } from "react";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import { api, models } from "../../wailsjs/go/models";

type LogLine = {
  timestamp: string;
  msg: string;
};
function Terminal() {
  const [logList, setLogList] = useState<LogLine[]>([]);
  const maxLogLine = 100;
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const cancel = EventsOn(
      api.MsgKey.processMsg,
      (msg: models.ProcessInfo) => {
        if (msg.Type === 0) {
          return;
        }
        setLogList((list) => {
          while (list.length > maxLogLine) {
            list.shift();
          }
          return [
            ...list,
            { timestamp: new Date().toLocaleString(), msg: msg.Msg },
          ];
        });
        ref.current?.scrollTo(0, 9999);
      }
    );
    return () => cancel();
  }, []);

  return (
    <div className="flex flex-col w-full h-full">
      <div
        ref={ref}
        className="terminal bg-gray-800 
    w-full h-full max-h-30 text-left overflow-auto"
      >
        {logList.map((log, index) => (
          <div key={index} className="flex flex-row">
            <div className="w-1/3">{log.timestamp}</div>
            <div className="w-2/3">{log.msg}</div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Terminal;
