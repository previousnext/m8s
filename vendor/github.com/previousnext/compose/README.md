# Compose

A go library for parsing Docker Compose yaml files into go structs.

## Usage

```go
import "github.com/previousnext/compose"

file := "path/to.yml"
dc := compose.Load(file)

// Get the image for service "app"
fmt.Println(dc.Services["app"].Image)
```
