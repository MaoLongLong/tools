# gsu

> 经常要在公司账号和社区账号间来回切换，于是有了 gsu

一个用来切换 git 用户的工具（git-su）

## 安装

```bash
go install go.chensl.me/tools/gsu@latest
```

## 配置

默认读取的配置文件为 `~/.gsu.yml`

```yml
users:
  mll:
    name: MaoLongLong
    email: chensl.dev@outlook.com
  other-user:
    name: foobar
    email: foo@bar.com
```

## 使用

```bash
# 根据配置文件中的 users.mll 修改 git 配置文件 ~/.gitconfig
gsu -g mll

# 如果只想修改某个项目的配置
cd projects/xxx
gsu mll

# 指定其他配置文件
gsu -c /path/to/config/file ...
```
