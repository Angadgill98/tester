



export interface Application_f {
    AppId: number;
    AppName: string;
    ParentAppId: number;
    EnvironmentId: number | null | {
        Int64: number;
        Valid: boolean;
    };
    Actions: number[];
    Sequences: number[];
    Children: number[];
}

export interface ChildApplicationInfo {
    Id: number;
    Name: string;
}


export interface Action_f {
    Id: number;
    Name: string;
    Type: string;
}

export interface Sequence_f {
    Id: number;
    Name: string;
}

export interface Environment_f {
    Id: number;
    Name: string;
}