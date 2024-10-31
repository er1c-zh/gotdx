import * as d3 from "d3";

type RealtimeGraphProps = {
  code: string;
  data: number[];
  width: number;
  height: number;
  marginTop: number;
  marginRight: number;
  marginBottom: number;
  marginLeft: number;
};
function RealtimeGraph(props: RealtimeGraphProps) {
  const x = d3.scaleLinear(
    [0, props.data.length - 1],
    [props.marginLeft, props.width - props.marginRight]
  );
  const y = d3.scaleLinear(
    d3.extent(props.data).map((d) => d ?? 0),
    [props.height - props.marginBottom, props.marginTop]
  );
  const line = d3.line((d, i) => x(i), y);
  return (
    <div>
      <svg width={props.width} height={props.height}>
        <path
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
          d={line(props.data) ?? undefined}
        />
        <g fill="white" stroke="currentColor" strokeWidth="1.5">
          {props.data.map((d, i) => (
            <circle key={i} cx={x(i)} cy={y(d)} r="2.5" />
          ))}
        </g>
      </svg>
    </div>
  );
}

export default RealtimeGraph;
