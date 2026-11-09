package models

import "database/sql"


type ActionInfo struct {
    Id   int64
    Name string
}

type SequenceInfo struct {
    Id   int64
    Name string
}

type ChildApplicationInfo struct {
    Id   int64
    Name string
}

type ApplicationFrontend struct {
    AppId         int64
    AppName       string
    ParentAppId   int64
    EnvironmentId sql.NullInt64

    Actions   []ActionInfo
    Sequences []SequenceInfo
    Children  []ChildApplicationInfo
}