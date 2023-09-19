import { Service } from "./service";

export type InstallMethod = "script" | "release" | "docker";

export type InstanceQuery = {
    features?: string[];
};

export type InstanceUpdate = {
    current_version: string;
    latest_version: string;
};

export type Operation = {
    op: string;
    from?: string;
    path: string;
    value?: string;
};

export type ServiceUpdate = {
    patch?: Operation[];
};

export type Instance = {
    service: Service;
    uuid: string;
    status: string;
    environment: { [key: string]: string };
    install_method?: InstallMethod;
    launch_on_startup?: boolean;
    display_name?: string;
    databases?: { [key: string]: string };
    update?: InstanceUpdate;
    service_update?: ServiceUpdate;
};

export type Instances = { [uuid: string]: Instance };