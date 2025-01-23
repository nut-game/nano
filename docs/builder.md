Builder
===

Nano offers a [`Builder`](../builder.go) object which can be utilized to define a sort of nano properties.

### PostBuildHooks

Post-build hooks can be used to execute additional actions automatically after the build process. It also allows you to interact with the built nano app.

A common use case is where it becomes necessary to perform configuration steps in both the nano builder and the nano app being built. In such cases, an effective approach is to internalize these configurations, enabling you to handle them collectively in a single operation or process. It simplifies the overall configuration process, reducing the need for separate and potentially repetitive steps.

```go
// main.go
cfg := config.NewDefaultBuilderConfig()
builder := nano.NewDefaultBuilder(isFrontEnd, "my-server-type", nano.Cluster, map[string]string{}, *cfg)

customModule := NewCustomModule(builder)
customModule.ConfigureNano(builder)

app := builder.Build()

// custom_object.go
type CustomObject struct {
	builder *nano.Builder
}

func NewCustomObject(builder *nano.Builder) *CustomObject {
	return &CustomObject{
		builder: builder,
	}
}

func (object *CustomObject) ConfigureNano() {
	object.builder.AddAcceptor(...)
	object.builder.AddPostBuildHook(func (app nano.App) {
		app.Register(...)
	})
}
```

In the above example the `ConfigureNano` method of the `CustomObject` is adding an `Acceptor` to the nano app being built, and also adding a post-build function which will register a handler `Component` that will expose endpoints to receive calls.
