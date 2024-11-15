module github.com/Eandalf/expressgo/examples/helloworld

go 1.23.2

require github.com/Eandalf/expressgo v0.0.0-00010101000000-000000000000

require golang.org/x/text v0.20.0 // indirect

replace github.com/Eandalf/expressgo => ../..

replace github.com/Eandalf/expressgo/bodyparser => ../../bodyparser

replace github.com/Eandalf/expressgo/cors => ../../cors
