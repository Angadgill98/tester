
export interface Application_f {
    AppId: string;
    AppName: string;
    ParentAppId: string;
    EnvironmentId: string | null;
    Actions: string[];
    Sequences: string[];
    Children: string[];
}

export interface ChildApplicationInfo {
    Id: string;
    Name: string;
}

export interface Action_f {
    Id: string;
    Name: string;
    Type: string;
}

export interface Sequence_f {
    Id: string;
    Name: string;
}

export interface Environment_f {
    Id: string;
    Name: string;
}

export interface Tab {
    AppId: string;
    TabType: string;
    TypeId: string;
    Name: string;
}

export interface HttpAction {
    ActionName: string;
    AppId: string;
    ActionId: string;
    Opts: HttpActionOpts;
    EnvironmentId: string;
}

export interface HttpActionOpts {
    Method: string;
    Url: string;
    Headers: Record<string, string[]>;
    Query: Record<string, string[]>;
    Body: unknown;
    Cookies: Record<string, string[]>;
}

export interface WSOpts {
    URL: string;
    Headers: Record<string, string[]>;
    QueryParams: Record<string, string[]>;
    Subprotocols: string[];
    Messages: WSMessage[];
    HandshakeTimeout: number;
    WriteTimeout: number;
}

export interface Response<T> {
    success: boolean;
    data: T;
    error: string;
}

export interface WSMessage {
    MessageID: string;
    Sequence: number;
    MessageType: string;
    Payload: unknown;
}

export interface WSAction {
    ActionID: string;
    AppID: string;
    ActionName: string;
    Opts: WSOpts;
    EnvironmentID: string;
}

export interface Sequence {
    SequenceID: string;
    AppID: string;
    SequenceName: string;
    ActionIDs: string[];
}
