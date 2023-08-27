import { BigTitle, Title } from "../../components/Text/Text";
import Header from "../../components/Header/Header";
import styles from "../ReverseProxy/ReverseProxy.module.sass";
import ProxyRedirect from "../../components/ProxyRedirect/ProxyRedirect";
import { Horizontal } from "../../components/Layouts/Layouts";
import Button from "../../components/Button/Button";
import { useFetch } from "../../hooks/useFetch";
import {
    addProxyRedirect,
    getProxyRedirects,
    removeProxyRedirect,
} from "../../backend/backend";
import Popup from "../../components/Popup/Popup";
import React, { useState } from "react";
import Input from "../../components/Input/Input";
import Spacer from "../../components/Spacer/Spacer";

type Props = {};

export default function ReverseProxy(props: Props) {
    const { data: redirects, reload } =
        useFetch<ProxyRedirects>(getProxyRedirects);

    const [showNewRedirectPopup, setShowNewRedirectPopup] = useState(false);

    const [source, setSource] = useState("");
    const [target, setTarget] = useState("");

    const onSourceChange = (e: any) => setSource(e.target.value);
    const onTargetChange = (e: any) => setTarget(e.target.value);

    const openNewRedirectPopup = () => setShowNewRedirectPopup(true);
    const closeNewRedirectPopup = () => setShowNewRedirectPopup(false);

    const addNewRedirection = () => {
        addProxyRedirect(source, target)
            .then(reload)
            .catch(console.error)
            .finally(closeNewRedirectPopup);
    };

    const onDelete = (uuid: string) => {
        removeProxyRedirect(uuid).then(reload).catch(console.error);
    };

    return (
        <div>
            <Header />
            <div className={styles.title}>
                <BigTitle>Reverse Proxy</BigTitle>
            </div>
            <div className={styles.redirects}>
                {Object.entries(redirects ?? {}).map(([uuid, redirect]) => (
                    <ProxyRedirect
                        enabled={true}
                        source={redirect.source}
                        target={redirect.target}
                        onDelete={() => onDelete(uuid)}
                    />
                ))}
            </div>
            <Horizontal className={styles.addRedirect} gap={10}>
                <Button primary onClick={openNewRedirectPopup} leftSymbol="add">
                    Add redirection
                </Button>
            </Horizontal>
            <Popup
                show={showNewRedirectPopup}
                onDismiss={closeNewRedirectPopup}
            >
                <Title>New redirection</Title>
                <Input
                    label="Source"
                    value={source}
                    onChange={onSourceChange}
                />
                <Input
                    label="Target"
                    value={target}
                    onChange={onTargetChange}
                />
                <Horizontal gap={10}>
                    <Spacer />
                    <Button onClick={closeNewRedirectPopup}>Cancel</Button>
                    <Button primary onClick={addNewRedirection}>
                        Send
                    </Button>
                </Horizontal>
            </Popup>
        </div>
    );
}
