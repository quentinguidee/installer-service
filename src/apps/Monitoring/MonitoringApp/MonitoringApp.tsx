import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import { useLocation } from "react-router-dom";
import l from "../../../components/NavLink/navlink";
import { ContainerLed } from "../../../components/ContainerLed/ContainerLed";

export default function MonitoringApp() {
    useTitle("Monitoring");

    const { pathname } = useLocation();
    const queryClient = useQueryClient();

    const { containers: prometheusContainers, isLoading: isLoadingPrometheus } =
        useContainers({
            tags: ["Vertex Monitoring - Prometheus Collector"],
        });
    const prometheusContainer = Object.values(prometheusContainers ?? {})[0];

    const { containers: grafanaContainers, isLoading: isLoadingGrafana } =
        useContainers({
            tags: ["Vertex Monitoring - Grafana Visualizer"],
        });
    const grafanaContainer = Object.values(grafanaContainers ?? {})[0];

    useServerEvent("/app/vx-containers/containers/events", {
        change: (e) => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    const sidebar = (
        <Sidebar rootUrl="/app/vx-monitoring" currentUrl={pathname}>
            <Sidebar.Group title="Overview">
                <Sidebar.Item
                    label="Metrics"
                    icon={<MaterialIcon icon="rule" />}
                    link={l("/app/vx-monitoring/metrics")}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Collectors">
                <Sidebar.Item
                    label="Prometheus"
                    icon={<SiPrometheus size={20} />}
                    link={l("/app/vx-monitoring/prometheus")}
                    trailing={
                        prometheusContainer && (
                            <ContainerLed
                                small
                                status={prometheusContainer?.status}
                            />
                        )
                    }
                />
            </Sidebar.Group>
            <Sidebar.Group title="Visualizations">
                <Sidebar.Item
                    label="Grafana"
                    icon={<SiGrafana size={20} />}
                    link={l("/app/vx-monitoring/grafana")}
                    trailing={
                        grafanaContainer && (
                            <ContainerLed
                                small
                                status={grafanaContainer?.status}
                            />
                        )
                    }
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoadingGrafana || isLoadingPrometheus} />
            <PageWithSidebar sidebar={sidebar} />
        </Fragment>
    );
}
