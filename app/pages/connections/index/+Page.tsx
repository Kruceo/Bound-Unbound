import { useData } from "vike-react/useData";
import type { Data } from "./+data.js";
import { Link } from "../../../components/Link.jsx";
import "./Page.less"
import Ico from "../../../components/Ico.jsx";
import ControlsReloadButton from "../../../components/ControlsReloadButton.jsx";



export default function Page() {
  const data = useData<Data>()
  return (
    <>
      <h1 className="page-title">DNS Servers</h1>
      {
        !data.data || data.data.length == 0 ?
          <p className="help-message">You don't have connected nodes. Visit <a href="/help">/help</a>.</p>
          : null
      }
      <div className="grid">

        {
          data.data?.map(each => {
            return <div key={each.name} className="connection-box">
              <h3>{each.name}</h3>
              <p>{each.remoteAddress}</p>
              <p>Unbound</p>
              <div className="bottom-bar">
                <ControlsReloadButton nodeId={each.remoteAddress} updateIfItChanges={data} />
                <Link aria-label="Domain Blocks" data-balloon-pos="down" className="button" href={`/connections/${each.remoteAddress}/blocks`}>
                  <Ico>block</Ico>
                </Link>
                <Link aria-label="Domain Redirects" data-balloon-pos="down" className="button" href={`/connections/${each.remoteAddress}/redirects`}>
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
