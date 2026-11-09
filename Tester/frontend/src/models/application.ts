export interface ActionInfo {
    Id: number;
    Name: string;
    Type:string
}

export interface SequenceInfo {
    Id: number;
    Name: string;
}

export interface ChildApplicationInfo {
    Id: number;
    Name: string;
}

export interface Application_f {
    AppId: number;
    AppName: string;
    ParentAppId: number;
    EnvironmentId: number | null;

    Actions: ActionInfo[];
    Sequences: SequenceInfo[];
    Children: ChildApplicationInfo[];
}