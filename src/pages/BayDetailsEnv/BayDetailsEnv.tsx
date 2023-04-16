import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import {
    Env,
    EnvVariable,
    getInstance,
    Instance,
    saveInstanceEnv,
} from "../../backend/backend";
import { useParams } from "react-router-dom";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";
import Button from "../../components/Button/Button";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";

type Props = {};

export default function BayDetailsEnv(props: Props) {
    const { uuid } = useParams();

    const [env, setEnv] = useState<{ env: EnvVariable; value: any }[]>();

    const [instance, setInstance] = useState<Instance>();

    const [uploading, setUploading] = useState(false);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);

    useEffect(() => {
        setEnv(
            instance?.environment.map((e) => ({
                env: e,
                value: instance?.env[e.name] ?? e.default ?? "",
            }))
        );
    }, [instance]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
        setSaved(false);
    };

    useEffect(() => {
        getInstance(uuid).then((i: Instance) => setInstance(i));
    }, [uuid]);

    const save = () => {
        const _env: Env = {};
        env.forEach((e) => {
            _env[e.env.name] = e.value;
        });
        setUploading(true);
        saveInstanceEnv(uuid, _env)
            .then(console.log)
            .catch(console.error)
            .finally(() => {
                setUploading(false);
                setSaved(true);
            });
    };

    return (
        <Fragment>
            <Title>Environment</Title>
            {env?.map((env, i) => (
                <EnvVariableInput
                    env={env.env}
                    value={env.value}
                    onChange={(v) => onChange(i, v)}
                    disabled={uploading}
                />
            ))}
            <Button
                primary
                large
                onClick={save}
                rightSymbol="save"
                loading={uploading}
                disabled={saved || saved === undefined}
            >
                Save
            </Button>
            {saved && (
                <Horizontal alignItems="center" gap={4}>
                    <Symbol name="check" />
                    Saved!
                </Horizontal>
            )}
        </Fragment>
    );
}