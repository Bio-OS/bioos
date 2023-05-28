### Folder structure

```sh

|-- public            # 静态资源，可通过 `/` 访问
|-- src               # 源代码
    |-- api           # 接口
    |-- assets        # 样式文件，svg文件
    |   |-- styles
    |   |-- svg       # 使用前使用svgo压缩，以 `icon-名称` 命名，可访问 `localhost:8901/icons` 查看已有icon，避免重复添加
    |-- components    # 内部组件
    |-- lib           # 外部组件
    |-- helpers       # 工具库
    |-- pages         # 页面，只写页面文件，组件请放到 `components/` 下
    |-- typings       # 自定义类型声明

```

### Naming convention

| 类别        | 规则        | 示例               |
| ----------- | ----------- | ------------------ |
| 文件夹名    | 小写中划线  | field-item         |
| 文件名      | 大驼峰      | WorkspaceList.tsx  |
| less 文件名 | 大驼峰.less | WorkspaceList.less |
| 函数名      | 小驼峰      | validateFile.ts    |
| 常量        | 大写下划线  | FILE_SIZE          |
| className   | 小驼峰      | .flexAlignCenter   |

### Installation

```sh
npm install
```

### Develop

```sh
npm run dev
```

### Deploy

```sh
npm run build
```

### Tips

- 图标预览请访问 localhost:8901/icons
