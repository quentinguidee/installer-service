import { Fragment, useEffect, useState } from "react";
import {
    Button,
    MaterialIcon,
    TextField,
    Title,
} from "@vertex-center/components";
import { api } from "../../../backend/api/backend";
import { Horizontal } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Content from "../../../components/Content/Content";
import { useSettings } from "../hooks/useSettings";

export default function SettingsNotifications() {
    const [webhook, setWebhook] = useState<string>();
    const [changed, setChanged] = useState(false);
    const [saving, setSaving] = useState(false);

    const { settings, errorSettings, isLoadingSettings } = useSettings();

    useEffect(() => {
        setWebhook(settings?.webhook);
    }, [settings]);

    const onWebhookChange = (e: any) => {
        setWebhook(e.target.value);
        setChanged(true);
    };

    const onSave = () => {
        setSaving(true);
        api.settings
            .patch({ webhook })
            .then(() => setChanged(false))
            .catch(console.error)
            .finally(() => setSaving(false));
    };

    const error = errorSettings;
    const isLoading = isLoadingSettings;

    return (
        <Content>
            <Title variant="h2">Notifications</Title>
            <ProgressOverlay show={isLoading || saving} />
            <APIError error={error} />
            {!error && (
                <Fragment>
                    <TextField
                        id="webhook"
                        label="Webhook"
                        value={webhook}
                        onChange={onWebhookChange}
                        disabled={isLoading}
                        placeholder={isLoading && "Loading..."}
                    />
                    <Horizontal
                        gap={20}
                        justifyContent="flex-end"
                        style={{ marginTop: 15 }}
                    >
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="save" />}
                            onClick={onSave}
                            disabled={!changed || saving}
                        >
                            Save
                        </Button>
                    </Horizontal>
                </Fragment>
            )}
        </Content>
    );
}
