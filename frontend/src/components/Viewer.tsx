type ViewerProps = {
  Code: string;
};
function Viewer(props: ViewerProps) {
  return (
    <div className="flex">
      <h1>Viewer</h1>
      <p>{props.Code}</p>
    </div>
  );
}

export default Viewer;
