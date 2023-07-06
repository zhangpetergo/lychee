





## 程序的优雅关闭

[graceful-shutdown](https://blog-zhangpetergo.vercel.app/posts/graceful-shutdown)





完成 debug pprof



完成 logger 的构建



## Viper Unmarshal

viper Unmarshal绑定的结构体的tag

mapstructure

```go
type Config struct {
	Server struct {
		APIHost   string `mapstructure:"api_host"`
		DebugHost string `mapstructure:"debug_host"`
		Version   string `mapstructure:"version"`
	}
}

```



结构体 tag 使用`mapstructure` 的原因是

> Viper 是 Go 中流行的配置管理库，它支持从各种来源（例如文件、环境变量和命令行参数）读取配置值。当使用 Viper 将配置值绑定到结构体时，库需要一种方法将配置文件中的键映射到结构体中的字段。
>
> 选择用作`mapstructure`结构字段的标记是特定于 Viper 的实现的。该`mapstructure`标签不是内置的 Go 标签，而是 Viper 使用的约定，用于指示配置键应如何映射到结构字段。
>
> Viper 使用该标签的原因`mapstructure`是利用该`github.com/mitchellh/mapstructure`库，该库提供了一种强大而灵活的方法来解码和映射 Go 中的任意数据结构。通过使用这个库，Viper 可以处理各种映射场景，包括嵌套结构、重命名字段和处理默认值。
>
> 通过在结构体字段上指定`mapstructure`标签，您可以为每个字段定义自定义映射规则，例如在配置文件中指定键名或处理特定数据类型。
>
> 总的来说，在 Viper 中使用`mapstructure`标签可以提供更灵活和可定制的方法来将配置值映射到结构体字段，从而为处理 Go 应用程序中的复杂配置场景提供了强大的工具集。