import Content from "../../components/Content";
import Loader from "../../components/Loader";
import SideBar from "../../components/SideBar";
import { NotificationProvider } from "../../components/NotificationContext";

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