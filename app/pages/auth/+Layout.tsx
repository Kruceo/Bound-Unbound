import Content from "../../components/Content";
export const config = {
  layout:true
}
export default function LayoutWithoutBar({ children }: { children: React.ReactNode }) {
    return (
      <div
        style={{
          background:"red",
          display: "flex",
          maxWidth: "100%",
          margin: "0px",
        }}
      >
        {/* <SideBar></SideBar> */}
        <Content>{children}</Content>
      </div>
    );
  }