import { useState } from "react";
import { Title } from "../../../../components/Text/Text";

import styles from "./InstanceDocker.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import Button from "../../../../components/Button/Button";
import { api } from "../../../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../../components/KeyValueInfo/KeyValueInfo";
import byteSize from "byte-size";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";

export default function InstanceDocker() {
    const { uuid } = useParams();

    const {
        data: info,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["instance_docker", uuid],
        queryFn: api.vxInstances.instance(uuid).docker.get,
    });

    const [recreatingContainer, setRecreatingContainer] = useState(false);
    const [recreatingContainerError, setRecreatingContainerError] = useState();

    const recreateContainer = () => {
        setRecreatingContainer(true);
        setRecreatingContainerError(undefined);

        api.vxInstances
            .instance(uuid)
            .docker.recreate()
            .catch((err) => {
                setRecreatingContainerError(
                    err?.response?.data?.message ?? err?.message
                );
                console.error(err);
            })
            .finally(() => {
                setRecreatingContainer(false);
            });
    };

    if (error) return <APIError error={error} />;

    return (
        <Vertical gap={30}>
            <ProgressOverlay show={isLoading} />
            <APIError error={recreatingContainerError} />

            <Vertical gap={20}>
                <Title className={styles.title}>Container</Title>
                <KeyValueGroup>
                    <KeyValueInfo
                        name="ID"
                        type="code"
                        icon="tag"
                        loading={isLoading}
                    >
                        {info?.container?.id}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Container Name"
                        type="code"
                        icon="badge"
                        loading={isLoading}
                    >
                        {info?.container?.name}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Platform"
                        type="code"
                        icon="computer"
                        loading={isLoading}
                    >
                        {info?.container?.platform}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Image</Title>
                <KeyValueGroup>
                    <KeyValueInfo
                        name="ID"
                        type="code"
                        icon="tag"
                        loading={isLoading}
                    >
                        {info?.image?.id}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Architecture"
                        type="code"
                        icon="memory"
                        loading={isLoading}
                    >
                        {info?.image?.architecture}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="OS"
                        type="code"
                        icon="computer"
                        loading={isLoading}
                    >
                        {info?.image?.os}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Size"
                        type="code"
                        icon="straighten"
                        loading={isLoading}
                    >
                        {byteSize(info?.image?.size).toString()}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Tags"
                        type="code"
                        icon="sell"
                        loading={isLoading}
                    >
                        {info?.image?.tags?.join(", ")}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Actions</Title>
                <Vertical gap={8}>
                    <Horizontal alignItems="center" gap={20}>
                        <Button
                            leftIcon="restart_alt"
                            onClick={recreateContainer}
                            disabled={recreatingContainer ?? isLoading}
                        >
                            Recreate container
                        </Button>
                    </Horizontal>
                </Vertical>
            </Vertical>
        </Vertical>
    );
}
