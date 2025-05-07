import { useEffect, useRef, useState } from "react";
import SideBar from "../../components/SideBar";
import { NotificationProvider } from "../../components/NotificationContext";
import { useData } from "vike-react/useData";
import { redirect } from "vike/abort";
import { navigate } from "vike/client/router";
import Content from "../../components/Content";
import Loader from "../../components/Loader";

export default function LayoutDefault({ children }: { children: React.ReactNode }) {
    return (<>
        <SideBar></SideBar>
        <NotificationProvider>
            <Content>{children}</Content>
        </NotificationProvider>
        <Loader></Loader>
    </>
    );
}
