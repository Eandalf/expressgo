Write-Host "expressgo: format"
go fmt

Write-Host "expressgo: install"
go install

Write-Host "goto: expressgo/examples/helloworld"
Push-Location ".\examples\helloworld"

Write-Host "expressgo/examples/helloworld: format"
go fmt

Write-Host "expressgo/examples/helloworld: build"
go build

Write-Host "goto: expressgo"
Pop-Location
