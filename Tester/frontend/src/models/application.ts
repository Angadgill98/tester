



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



export interface Tab{
    AppId:number
    TabType:string
    TypeId:number
    Name:string
}



export interface HttpAction {
    ActionName: string;
    AppId: number;
    ActionId: number;
    Opts: HttpActionOpts;
    EnvironmentId: number;
}

export interface HttpActionOpts {
    Method: string;
    URL: string;
    Headers: Record<string, string>;
    Query: Record<string, string>;
    Body: unknown;
    Cookies: Record<string, string>;
}



export interface WSOpts {
    URL: string;
    Headers: Record<string, string>;
    QueryParams: Record<string, string>;
    Subprotocols: string[];
    Messages: WSMessage[];
    HandshakeTimeout: number;
    WriteTimeout: number;
}

export interface WSMessage {
    MessageID: number;
    Sequence: number;
    MessageType: string;
    Payload: unknown;
}

export interface WSAction {
    ActionID: number;
    AppID: number;
    ActionName: string;
    Opts: WSOpts;
    EnvironmentID: number;
}

export interface Sequence {
    SequenceID: number;
    AppID: number;
    SequenceName: string;
    ActionIDs: number[];
}