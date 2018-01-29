//+build ignore

package main

func (r HelloMessageDef) Render() react.Element {
	return react.Div(nil,
		react.S("Hello "+r.Props().Name),
	)
}
