package component

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/astaxie/beego"

	rpcTypes "github.com/drep-project/drepcli/rpc/types"
)

var mappingMethodMap = map[string]string{
	//  "/GetAllBlocks": "*:GetAllBlocks",
}

type RestController struct {
	api reflect.Value
}

func (controller *RestController) Stop() {
	controller.api.MethodByName("StopRun").Call(nil)
}

func StartRest(apiDesc rpcTypes.RestDescription) *RestController {
	port := flag.String("port", "55550", "port:default is 55551")
	fmt.Println("http server is ready for listen port:", port)
	switch controller := apiDesc.Api.Interface().(type) {
	case beego.ControllerInterface:
		beego.Router("/", controller)
		t := apiDesc.Api.Type()
		for i := 0; i < t.NumMethod(); i++ {
			mName := t.Method(i).Name
			beego.Router("/"+mName, controller, "*:"+mName)
		}
	}
	beego.Run(":" + *port)
	return &RestController{
		api: apiDesc.Api,
	}
}
