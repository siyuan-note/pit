# pit

思源笔记更新程序，Pit 取自 F1 方程式赛车维修区英文名。

## 启动参数

* 配置目录路径：通过参数 `--conf` 传入，实参为思源笔记配置目录路径
* 工作目录路径：通过参数 `--wd` 传入，实参为思源笔记安装路径下的 resources 文件夹绝对路径
* 更新包路径：通过 `--pack` 传入，实参为思源笔记临时目录下的 update.zip，即 `$os.temp/siyuan/update.zip`

## 工作流程

1. 删除工作目录下的文件
   * kernel.exe / kernel-linux / kernel-darwin
   * app.asar
   * appearance 文件夹
   * stage 文件夹
   * guide 文件夹
2. 解压更新包到工作目录
   * Linux、macOS 使用 `chmod +x` 赋予执行权限
3. 删除更新包
