import { Fragment } from "react";
import { Title } from "../../../components/Text/Text";
import { useFetch } from "../../../hooks/useFetch";
import { About } from "../../../models/about";
import { api } from "../../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";

import styles from "./SettingsAbout.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/Error";
import { ProgressOverlay } from "../../../components/Progress/Progress";

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(api.about.get);

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <APIError error={error} />
            <Vertical gap={20}>
                <Title className={styles.title}>Vertex</Title>
                <KeyValueGroup>
                    <KeyValueInfo
                        name="Version"
                        type="code"
                        symbol="tag"
                        loading={loading}
                    >
                        {about?.version}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Commit"
                        type="code"
                        symbol="commit"
                        loading={loading}
                    >
                        {about?.commit}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Release date"
                        type="code"
                        symbol="calendar_month"
                        loading={loading}
                    >
                        {about?.date}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Compiled for"
                        type="code"
                        symbol="memory"
                        loading={loading}
                    >
                        {about?.os}
                        {about?.arch && `/${about?.arch}`}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>
        </Fragment>
    );
}