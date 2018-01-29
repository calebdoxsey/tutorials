package main

import (
	"io/ioutil"
	"os"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Get("document").Get("documentElement").Set("innerHTML", `<!DOCTYPE html>
<html>
	<head></head>
	<body>
		<form onsubmit="javascript:event.preventDefault()">
			<input type="text" id="filename" placeholder="/tmp/example.txt">
			<input type="submit" value="Load" onclick="onloadfile()">
			<br>
			<br>
			<textarea id="filecontents" cols="70" rows="5" placeholder="contents"></textarea>
			<br>
			<input type="submit" value="Save" onclick="onsavefile()">
		</form>
	</body>
</html>
	`)

	// START OMIT

	js.Global.Set("onloadfile", func() {
		filename := js.Global.Get("document").Call("getElementById", "filename").Get("value").String()
		go func() {
			bs, err := ioutil.ReadFile(filename)
			if err != os.ErrNotExist && err != nil {
				panic(err)
			}
			js.Global.Get("document").Call("getElementById", "filecontents").Set("value", string(bs))
		}()
	})
	js.Global.Set("onsavefile", func() {
		filename := js.Global.Get("document").Call("getElementById", "filename").Get("value").String()
		filecontents := js.Global.Get("document").Call("getElementById", "filecontents").Get("value").String()
		go func() {
			err := ioutil.WriteFile(filename, []byte(filecontents), 0777)
			if err != nil {
				panic(err)
			}
		}()
	})

	// END OMIT
}
