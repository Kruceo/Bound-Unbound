import "./Layout.less";

import React, { useEffect, useState } from "react";
import SideBar from "./SideBar";
import { onAuthStatus, onAuthToken } from "./LayoutDefault.telefunc";
import { navigate } from "vike/client/router";

export default function LayoutDefault({ children }: { children: React.ReactNode }) {
  const [logged, setLogged] = useState<boolean>(false)

  useEffect(() => {
    onAuthStatus().then(f => {
      console.log(f)
      if (!f.Data.AlreadyRegistered) {
        navigate("/auth/register")
        return
      }
      onAuthToken().then(d => {
        if (!d) navigate("/auth/signin")
        else setLogged(true)
      }).catch(() => navigate("/auth/signin"))
    })

  }, [])

  return (
    <div
      style={{
        display: "flex",
        maxWidth: "100%",
        margin: "auto",
      }}
    >
      {
        logged ? <>
          <SideBar></SideBar>
          <Content>{children}</Content>
        </>
          : ""
      }
      <div className="transition-loader">
      <span className="loader"></span>
      </div>
    </div>
  );
}
export function LayoutWithoutBar({ children }: { children: React.ReactNode }) {
  return (
    <div
      style={{
        display: "flex",
        maxWidth: "100%",
        // background:"red",
        margin: "0px",
      }}
    >
      {/* <SideBar></SideBar> */}
      <Content>{children}</Content>
    </div>
  );
}



function Content({ children }: { children: React.ReactNode }) {
  return (
    <div id="page-container" style={{ width: "100%", boxSizing: "border-box" }}>
      <div
        id="page-content"
        style={{
          boxSizing: "border-box",
          width: "100%",
          padding: 48,
          // paddingTop:55,
          paddingLeft: 48,
          // paddingBottom: 50,
          display: "block",
          minHeight: "calc(100vh)",
        }}
      >
        {children}
      </div>
    </div>
  );
}
