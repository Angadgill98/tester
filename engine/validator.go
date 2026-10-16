package main



var validators= map[string]func([]any) Response{

}


func Validate(name string, args []any) Response {
	fn, exists := validators[name]
	if !exists {
		return Response{
			Status: false,
			Msg:    "Unknown Validator op: " + name,
		}
	}

	return fn(args)
}