
export default function Content({ children }: { children: React.ReactNode }) {
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