import "./style.less";

import React from "react";
import logoUrl from "../assets/logo.svg";
import { Link } from "../components/Link.jsx";
import SideBar from "./SideBar";

export default function LayoutDefault({ children }: { children: React.ReactNode }) {
  return (
    <div
      style={{
        display: "flex",
        maxWidth: "100%",
        margin: "auto",
      }}
    >
      <SideBar></SideBar>
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
          paddingTop:0,
          paddingLeft: 75 + 48,
          // paddingBottom: 50,
          display:"block",
          minHeight: "100vh",
        }}
      >
        {children}
      </div>
    </div>
  );
}
