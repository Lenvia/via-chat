package views

import (
	"embed"
	"html/template"
)

// main 函数中定义的全局变量 在其他包中无法调用 ，重新定义新包实现全局变量
// go:embed 指令用于嵌入式文件系统的语法定义
// go:embed 不支持相对路径，只能获取当前目录下的目录或文件
var (
	//go:embed *.html
	embedTmpl embed.FS

	// 自定义的函数必须在调用ParseFiles() ParseFS()之前创建。
	funcMap = template.FuncMap{}
	GoTpl   = template.Must(
		template.New("").
			Funcs(funcMap).               // 模板文件中的自定义函数可以在 funcMap 映射中定义，并被嵌入式文件系统中的文件所使用。
			ParseFS(embedTmpl, "*.html")) // 将 embedTmpl 文件系统和 *.html 文件的匹配模式传递给该函数，从而解析并加载所有符合模式的模板文件
)

// 通过以上实现，其他包可以直接引用 views.GoTpl 全局变量，并使用 "text/html" content type 类型渲染模板，从而返回客户端 HTML 页面。
