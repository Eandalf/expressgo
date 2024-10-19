Write-Host "task: prepare"

Write-Host "expressgo: tidy mod"
go mod tidy -v

Write-Host "goto: expressgo/examples/helloworld"
Push-Location ".\examples\helloworld"

Write-Host "expressgo/examples/helloworld: tidy mod"
go mod tidy -v

Write-Host "expressgo/examples/helloworld: clean"
go clean -x

Write-Host "goto: expressgo"
Pop-Location

Write-Host "task: build"

Write-Host "expressgo: format"
go fmt

Write-Host "expressgo: install"
go install -v

Write-Host "goto: expressgo/examples/helloworld"
Push-Location ".\examples\helloworld"

Write-Host "expressgo/examples/helloworld: format"
go fmt

Write-Host "expressgo/examples/helloworld: build"
go build -v

Write-Host "goto: expressgo"
Pop-Location
