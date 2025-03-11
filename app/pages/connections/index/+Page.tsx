import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import { Link } from "../../../components/Link.jsx";
import "./Page.less"
import { apiUrl } from "../../../api/api.js";
import Ico from "../../../components/Ico.jsx";
import { onReloadActions } from "../actions.telefunc.js";
import ControlsReloadButton from "../../../components/ControlsReloadButton.jsx";



export default function Page() {

  const data = useData<Data>()
  return (
    <>
      <h1 className="page-title">DNS Servers</h1>
      <div className="grid">
        {
          data.Data.map(each => {
            return <div key={each.Name} className="connection-box">
              <h3>{each.Name}</h3>
              <p>{each.RemoteAddress}</p>
              <p>Unbound</p>
              <div className="bottom-bar">
                <ControlsReloadButton nodeId={each.RemoteAddress} updateIfItChanges={data} />
                <Link aria-label="Domain Blocks" data-balloon-pos="down" className="button" href={`/connections/${each.RemoteAddress}/blocked`}>
                  <Ico>block</Ico>
                </Link>
                <Link aria-label="Domain Redirects" data-balloon-pos="down" className="button" href={`/connections/${each.RemoteAddress}/redirects`}>
                  <Ico>airline_stops</Ico>
                </Link>
              </div>
            </div>
          })
        }
      </div>
    </>
  );
}
