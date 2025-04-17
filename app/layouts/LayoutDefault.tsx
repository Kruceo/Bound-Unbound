import "./Layout.less";

import React, { useEffect, useRef, useState } from "react";
import SideBar from "./SideBar";
import { NotificationProvider } from "./NotificationContext";
import { onAuthToken } from "./LayoutDefault.telefunc";
import { navigate } from "vike/client/router";

export default function LayoutDefault({ children }: { children: React.ReactNode }) {
  const [logged, setLogged] = useState<boolean>(true)
  const singleExecuted = useRef(false)
  useEffect(() => {
    // this is because in dev mode react mounts 2 times all components
    // but this 'onAuthStatus' method have a brute force protection
    if (singleExecuted.current) return console.log("skipped");
    singleExecuted.current = true

    onAuthToken().then(d => {
      if (!d.ok) navigate("/auth/signin")
      else setLogged(true)
      if (d.cookies) {
        for (const cookie of d.cookies) {
          document.cookie = cookie
        }
      }
    }).catch(() => navigate("/auth/signin"))

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
          <NotificationProvider>
            <Content>{children}</Content>
          </NotificationProvider>
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
