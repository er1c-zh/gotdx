import { v2 } from "../../wailsjs/go/models";
import CandleStickView from "./CandleStick";

type ViewerProps = {
  Code: string;
};
function Viewer(props: ViewerProps) {
  return (
    <div className="flex flex-row w-full h-full">
      <div className="flex flex-col w-36 h-full bg-gray-500">
        <p>{props.Code}</p>
      </div>
      <div className="flex flex-row w-full">
        <div className="flex flex-col h-full w-1/2">
          <div className="flex h-1/2 w-full min-w-full bg-red-400">
            <CandleStickView
              code={props.Code}
              period={v2.CandleStickPeriodType.CandleStickPeriodType1Day}
            />
          </div>
          <div className="flex h-1/2 w-full">
            <CandleStickView
              code={props.Code}
              period={v2.CandleStickPeriodType.CandleStickPeriodType1Day}
            />
          </div>
        </div>
        <div className="flex flex-col h-full w-1/2">quote and tick</div>
      </div>
    </div>
  );
}

export default Viewer;
