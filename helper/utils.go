package helper

import (
	"backEnd/config"
	"fmt"
	"github.com/fatih/color"
	"reflect"
)

func Copy(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
func FlashLogo() {
	fmt.Println(" /$$             /$$$$$$                                                      /$$                              \n|  $$           /$$__  $$                                                    | $$                              \n \\  $$         | $$  \\__/  /$$$$$$  /$$$$$$$  /$$$$$$$   /$$$$$$   /$$$$$$$ /$$$$$$    /$$$$$$   /$$$$$$       \n  \\  $$ /$$$$$$| $$       /$$__  $$| $$__  $$| $$__  $$ /$$__  $$ /$$_____/|_  $$_/   /$$__  $$ /$$__  $$      \n   /$$/|______/| $$      | $$  \\ $$| $$  \\ $$| $$  \\ $$| $$$$$$$$| $$        | $$    | $$  \\ $$| $$  \\__/      \n  /$$/         | $$    $$| $$  | $$| $$  | $$| $$  | $$| $$_____/| $$        | $$ /$$| $$  | $$| $$            \n /$$/          |  $$$$$$/|  $$$$$$/| $$  | $$| $$  | $$|  $$$$$$$|  $$$$$$$  |  $$$$/|  $$$$$$/| $$            \n|__/            \\______/  \\______/ |__/  |__/|__/  |__/ \\_______/ \\_______/   \\___/   \\______/ |__/           ")
	color.Red(config.APP_VERSION)
}
