import { useEffect, useState } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { api, models } from "../../wailsjs/go/models";

type StatusBarProps = {
  Components: React.ComponentType<any>[];
};
function StatusBar(props: StatusBarProps) {
  const [time, setTime] = useState("");
  const [serverInfo, setServerInfo] = useState<models.ServerStatus>(
    models.ServerStatus.createFrom({})
  );
  useEffect(() => {
    const ticker = setInterval(() => {
      setTime(new Date().toLocaleTimeString());
    }, 500);
    const cancel = EventsOn(
      api.MsgKey.serverStatus,
      (info: models.ServerStatus) => {
        setServerInfo(info);
      }
    );
    return () => {
      clearInterval(ticker);
      cancel();
    };
  }, []);
  return (
    <div className="flex flex-row w-full bg-gray-800">
      <div className="w-24 bg-yellow-900">{time}</div>
      <div
        className={`w-auto max-w-48 overflow-x-hidden truncate ... ${
          serverInfo.Connected ? "bg-green-700" : "bg-red-900"
        } text-left px-2`}
      >
        {serverInfo.ServerInfo
          ? serverInfo.ServerInfo
          : serverInfo.Connected
          ? "Connected"
          : "Disconnected"}
      </div>
      <div>
        {props.Components.map((C: React.ComponentType<any>) => {
          return <C />;
        })}
      </div>
    </div>
  );
}

export default StatusBar;
