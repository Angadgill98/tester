package app



type Application struct {
	AppId       int64
	AppName     string
	ParentAppId int64
	Actions     []int64
	Sequences   []int64
	Children    []int64
}
