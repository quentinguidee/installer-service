import { Fragment } from "react";
import { Title } from "../../../../components/Text/Text";
import Icon from "../../../../components/Icon/Icon";

import styles from "./ContainerHome.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import classNames from "classnames";
import useContainer from "../../../hooks/useContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { APIError } from "../../../../components/Error/APIError";

export default function ContainerHome() {
    const { uuid } = useParams();

    const { container, isLoading, error } = useContainer(uuid);

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>URLs</Title>
            <APIError error={error} />
            <nav className={styles.nav}>
                {container?.service?.urls &&
                    container?.service?.urls
                        .filter((u) => u.kind === "client")
                        .map((u) => {
                            const portEnvName = container?.service?.environment
                                ?.filter((e) => e.type === "port")
                                ?.find((e) => e.default === u.port).name;

                            const port =
                                container?.environment[portEnvName] ?? u.port;
                            const disabled = container.status !== "running";

                            // @ts-ignore
                            let url = new URL(window.apiURL);
                            url.port = port;
                            url.pathname = u.home ?? "";

                            return (
                                <a
                                    key={u.port}
                                    href={url.href}
                                    target="_blank"
                                    rel="noreferrer"
                                    className={classNames({
                                        [styles.navItem]: true,
                                        [styles.navItemDisabled]: disabled,
                                    })}
                                >
                                    <Horizontal>
                                        <Icon name="public" />
                                        <Spacer />
                                        <Icon name="open_in_new" />
                                    </Horizontal>
                                    <div>{url.href}</div>
                                </a>
                            );
                        })}
            </nav>
        </Fragment>
    );
}