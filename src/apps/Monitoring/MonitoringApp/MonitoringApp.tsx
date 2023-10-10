import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { useFetch } from "../../../hooks/useFetch";
import { api } from "../../../backend/backend";
import { Instances } from "../../../models/instance";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";

export default function MonitoringApp() {
    const {
        data: instances,
        loading,
        reload: reloadInstances,
    } = useFetch<Instances>(api.vxInstances.instances.all);

    const ledPrometheus = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-prometheus-collector")
    );

    const ledGrafana = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-grafana-visualizer")
    );

    useServerEvent("/app/vx-instances/instances/events", {
        change: (e) => {
            console.log(e);
            reloadInstances().finally();
        },
    });

    const sidebar = (
        <Sidebar root="/app/vx-monitoring">
            <SidebarGroup title="Overview">
                <SidebarItem
                    icon="rule"
                    to="/app/vx-monitoring/metrics"
                    name="Metrics"
                />
            </SidebarGroup>
            <SidebarGroup title="Collectors">
                <SidebarItem
                    icon={<SiPrometheus size={20} />}
                    to="/app/vx-monitoring/prometheus"
                    name="Prometheus"
                    led={ledPrometheus && { status: ledPrometheus?.status }}
                />
            </SidebarGroup>
            <SidebarGroup title="Visualizations">
                <SidebarItem
                    icon={<SiGrafana size={20} />}
                    to="/app/vx-monitoring/grafana"
                    name="Grafana"
                    led={ledGrafana && { status: ledGrafana?.status }}
                />
            </SidebarGroup>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <PageWithSidebar title="Monitoring" sidebar={sidebar} />
        </Fragment>
    );
}
