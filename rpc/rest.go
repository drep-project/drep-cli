package rpc

import (
    "flag"
    "fmt"
    "github.com/astaxie/beego"
    "reflect"
)


var mappingMethodMap = map[string] string {
  //  "/GetAllBlocks": "*:GetAllBlocks",
}

type RestDescription struct {
    Api reflect.Value
}

type Request struct {
    Method string `json:"method"`
    Params string `json:"params"`
}

type Response struct {
    Success bool `json:"success"`
    ErrorMsg string `json:"errMsg"`
    Data interface{} `json:"body"`
}


type RestController struct {
    api reflect.Value
}

func (controller *RestController) Stop() {
    controller.api.MethodByName("StopRun").Call(nil)
}

func StartRest(apiDesc RestDescription) *RestController{
    port := flag.String("port", "55550", "port:default is 55551")
    fmt.Println("http server is ready for listen port:", port)
    switch controller := apiDesc.Api.Interface().(type) {
        case beego.ControllerInterface:
            beego.Router("/",controller)
            t := apiDesc.Api.Type()
            for i:=0;i<t.NumMethod();i++{
                mName := t.Method(i).Name
                beego.Router("/"+mName, controller, "*:"+mName)
            }
    }
    beego.Run(":" + *port)
    return &RestController{
        api:apiDesc.Api,
    }
}

