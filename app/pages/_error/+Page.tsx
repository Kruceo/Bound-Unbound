import React from "react";
import { usePageContext } from "vike-react/usePageContext";
import './Page.less'

export default function Page() {
  const { is404 } = usePageContext();
  if (is404) {
    return (
      <>
        <div className="container">
          <h2>404 Page Not Found</h2>
          <p>This page could not be found.</p>
          <a className="button" href="/">Go Home</a>
        </div>
      </>
    );
  }
  return (
    <>
      <div className="container">
        <h2>This page can't load</h2>
        <p>This page reach some problem to load.</p>
        <a className="button" href="/">Home</a>
      </div>
    </>
  );
}
