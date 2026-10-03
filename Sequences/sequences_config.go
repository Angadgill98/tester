package sequences

import (
	"fmt"
	actions "tester/Actions/v1"
)

func Ignore(){//wont le me import actiosn wihtout using it 
	var a=actions.ActionFS_path
	fmt.Printf("",a)
}


var Sequence_Config = map[string][]SequenceAction{

}
