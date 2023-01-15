<h1 align="center">GameASM - Game Assets Manager</h1>

GameASM aka game assets manager, is a simple implementation for Golang to create and recover incremental update for games.

## Download

下载 GameASM 的最新 Github 版本，请点击：

[![Go](https://github.com/xiaoyu08/gameasm/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/xiaoyu08/gameasm/actions/workflows/go.yml)

在打开的页面中选择最新的成功构建，页面向下拉找到 Artifacts，然后点击文件 `artifact` 下载压缩包即可。

## Purpose

To deal with the incredible game assets size and cover the time and storage it costs,
GameASM can create incremental update packages that is much smaller,
and also easier to manage and check its integrity.

## Features

 - Create incremental update packages with only one command
 - Creating, deleting and modifying files are supported to sync
 - Easy to recover

## Getting Started

### 服务端

作为服务端，你可能会创建和发布增量安装包。下列命令将 **创建一个版本字符串为 1.0 的安装包**.

```
./gameasm.exe server push 1.0
```

由于我们是第一次运行，GameASM 询问我们游戏的原始文件目录：

```
Game Assets Manager [v1]
当前工作目录: ...

打扰一下, 无法找到配置文件, 请回答下列问题:
          游戏资源原始目录 :
```

填写好游戏的原始目录后，将在工作目录下生成 `config.dat` 配置文件。按回车键，之后 GameASM 将正式开始创建：

```
Game Assets Manager [v1]
当前工作目录: ...

打扰一下, 无法找到配置文件, 请回答下列问题:
          游戏资源原始目录 : [您游戏的构建成品目录]
谢谢，已创建配置文件，正在继续下一步...

打包根目录: [您游戏的构建成品目录]
打包缓存文件目录: [工作目录]\gameasm\cache\1_0

INFO 这是第一次发布版本，安装包可能较大...
INFO 新增文件: /test1
INFO 新增文件: /test2
INFO 写入资源索引文件...
INFO 打包为 zip 文件...
INFO 请不要删除 _latest_index.json 文件，这是为了下一次增量更新准备的。
```

接下来，程序将自动打开生成的安装包所在的文件夹。您可以按照需要随意复制和删除其中的 `zip` 文件。

**特别注意：** 在 `artifacts` 文件夹内，将自动生成 `_latest_index.json` 文件，该文件将用于下一次更新时判断文件的增删改情况，请不要删除。删除这个文件后，程序会重新生成包含所有文件的初始安装包，无法进行增量更新。

### 客户端

作为客户端，你可能会对下载好的安装包进行解压和恢复操作。现在我们假设 gameasm.exe 存放在 E:/Game/ 文件夹内。


下列命令将 **解压 `E:/Game/1_0.zip` 文件并将其恢复到 `E:/Game/gameasm/game/` 文件夹内**。

```
./gameasm.exe client install 1_0.zip
```

你也可以指定相对目录（暂时不支持绝对路径）。下列命令将解压 `E:/Game/downloads/1_0.zip` 文件。

```
./gameasm.exe client install downloads/1_0.zip
```

输出如下图：

```
Game Assets Manager [v1]
当前工作目录: ...

INFO 拷贝新文件: ee0de434ebf3edc92188830fd113b759 => xxx
INFO 拷贝新文件: f0391d06bbe7924c71403a718b7cefb1 => xxx
INFO 拷贝新文件: 1d015c1bb1ac2ce88fd464b2ff164ac3 => xxx
INFO 拷贝新文件: 41e5de2ced0c34f535b307c62ef1ca1f => xxx
INFO 拷贝新文件: 7c2d1d67ec12f69cf0eff586bd9be596 => xxx
INFO 拷贝新文件: c8f072cce9f7276221c3a6d99f80fc19 => xxx
INFO Done. 6 files recovered.
```

现在，安装包内的内容将会被同步到 `E:/Game/gameasm/game/` 文件夹内。
