import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";

export default function ReverseProxyApp() {
    const sidebar = (
        <Sidebar root="/app/vx-reverse-proxy">
            <SidebarGroup title="Providers">
                <SidebarItem
                    name="Vertex Reverse Proxy"
                    icon="router"
                    to="/app/vx-reverse-proxy/vertex"
                />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="Reverse Proxy" sidebar={sidebar} />;
}