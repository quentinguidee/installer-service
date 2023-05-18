import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { getAbout } from "../../backend/backend";
import Loading from "../../components/Loading/Loading";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./SettingsAbout.module.sass";
import { Error } from "../../components/Error/Error";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";

type Props = {};

export default function SettingsAbout(props: Props) {
    const [version, setVersion] = useState<string>();
    const [commit, setCommit] = useState<string>();
    const [date, setDate] = useState<string>();

    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState();

    useEffect(() => {
        setLoading(true);
        getAbout()
            .then((about) => {
                setVersion(about.version);
                setCommit(about.commit);
                setDate(about.date);
                setLoading(false);
            })
            .catch((err) =>
                setError(err?.response?.data?.message ?? err?.message)
            );
    }, []);

    return (
        <Fragment>
            <Title>About</Title>
            {error && <Error error={error} />}
            {loading && !error && <Loading />}
            {!loading && (
                <Vertical gap={4}>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="tag" />
                        </div>
                        <div className={styles.item}>Version:</div>
                        <div className={styles.item}>{version}</div>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="commit" />
                        </div>
                        <div className={styles.item}>Commit:</div>
                        <div className={styles.item}>{commit}</div>
                    </Horizontal>
                    <Horizontal gap={12} alignItems="center">
                        <div className={styles.item}>
                            <Symbol name="calendar_month" />
                        </div>
                        <div className={styles.item}>Release date:</div>
                        <div className={styles.item}>{date}</div>
                    </Horizontal>
                </Vertical>
            )}
        </Fragment>
    );
}
