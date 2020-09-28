## 此脚本的运行需要在root下安装与运行

## App自动安装与升级需要提供的接口
*)必须开启http监听(由升级脚本来确定端口，例如：127.0.0.1:8083).
*)此GET请求(http://127.0.0.1:8083/version)返回200响应内容为版本号.
*)此GET请求(http://127.0.0.1:8083/hacheck)返回200响应内容为1, 代表存活.
*)修改升有脚本的版本号需要重新发布.


## 结构说明

boot.sh -- 负责upgrade.sh的升级与调用。此为基础程序，不建议在线升级，每分钟由系统定时器调用一次。
boot -- boot.sh依赖的库

ver.txt -- upgrade.sh的版本号
lib -- upgrade.sh依赖的库
upgrade.sh -- 由boot.sh调用，负责upgrade文件夹的更新与调用工作。
publish.sh -- 打包lib, upgrade.sh与ver.txt版本

upgrade -- 除package.sh外的所有*.sh将被upgrade.sh调用, *.sh负责app的安装与升级工作
    ver.txt -- App脚本的发布版本号
    publish.sh -- 打包程序
    *.sh -- App安装与升级使用的脚本
    lib -- *.sh依赖的库

install.sh -- 安装boot.sh
uninstall.sh -- 移除boot.sh

hmu2000 -- hmu2000设备上的App升级脚本开发目录
