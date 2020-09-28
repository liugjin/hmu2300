#!/bin/sh

# ================================
#  配置区域开始
# ================================
# ========= upgrade url ============
# url1="https://upgrade-hmu.huayuan-iot.com/hmu/upgrade/app" # 应使用https，但目前条件未充足
url1="http://upgrade-hmu.huayuan-iot.com/hmu/upgrade/app"
url2=$url1
url3=$url1
appMgrName=hmuupgrade # 安装的文件名
# ================================
#  配置区域结束
# ================================


# ========= global define ============
basepath=$(cd `dirname $0`; pwd)
hmuUpDir=/tmp/${appMgrName} # update directory
mkdir -p $hmuUpDir
rm -rf $hmuUpDir/* #**!!需要特别注意不能删根目录!!***#
dlTarget=$hmuUpDir/${appMgrName}.tar.gz

hmuTarDir=$basepath/upgrade
mkdir -p $hmuTarDir

# 导入工具包
. ${basepath}/lib/base.sh

ExitHmuUpgrade(){
    UnlockFile $appMgrName

    # 执行每一个之程序的升级调用
    for app in `ls "${basepath}/upgrade/"|egrep ".sh"`; do
        if [ "$app" = "publish.sh" ]; then
            # ignore publish.sh
            continue
        fi
        echo "Run "${basepath}/upgrade/$app
	    ${basepath}/upgrade/$app
    done
}

# =============== main ================
# =============== get lock ===========
LockFile $appMgrName
if [ $? -eq 1 ]; then
    # in locking
    return;
fi
Log "Clear log for upgrading">/tmp/${appMgrName}_log.txt
# ============file lock end===========

curVer=""
if [ -f ${basepath}/upgrade/ver.txt ]; then
    curVer=$(cat ${basepath}/upgrade/ver.txt)
fi
Log "curVer:"$curVer

# Fetch upgrade information
downloadDone=0
for val in $url1 $url2 $url3; do
    # echo "download:"$?
    DownloadApp $hmuUpDir $dlTarget $val $curVer
    case $? in
        0)
            # not date to upgrade
            ExitHmuUpgrade
            return;
            ;;
        1)
            # download done
            downloadDone=1
            break;
            ;;
        2)
            # upgrade failed
            continue
            ;;
        *)
            Log "Unknow status:"$?
            ;;
    esac
done
if [ $downloadDone -eq 0 ]; then
    Log "No url for download."
    ExitHmuUpgrade
    return;
fi

# Decompress
Log "Installing new version"
rm -rf ${hmuTarDir}/* # 安装前先清空
tar -xzf $dlTarget -C $hmuTarDir
if [ $? -ne 0 ]; then
    Log "failed decompress:"$?
    Log "dlTarget:"$dlTarget
    Log "hmuTarDir:"$hmuTarDir
    ExitHmuUpgrade
    return
fi

Log "update success, clean cached data"
Log "clean:"$dlTarget
rm -rf $dlTarget
ExitHmuUpgrade
