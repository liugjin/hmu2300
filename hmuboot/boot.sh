#!/bin/sh
# hmubootUrl="https://upgrade-hmu.huayuan-iot.com/hmu/upgrade/boot" # 应使用https，但目前条件未充足
hmubootUrl="http://upgrade-hmu.huayuan-iot.com/hmu/upgrade/boot"

# ========= global define ============
basepath=$(cd `dirname $0`; pwd)
appMgrName=hmuboot
hmuUpDir=/tmp/$appMgrName # update directory
mkdir -p $hmuUpDir
rm -rf ${hmuUpDir}/* #**!!需要特别注意不能删根目录!!***#
dlTarget=${hmuUpDir}/${appMgrName}.tar.gz
tarTarget=$basepath

# import base function
. $basepath/boot/base.sh

# =============== main ================
Log "Run boot.sh"
# =============== get lock ===========
LockFile $appMgrName
if [ $? -eq 1 ]; then
    # in locking
    return;
fi
Log "Clear log for upgrading">/tmp/${appMgrName}_log.txt
# ============file lock end===========

curVer=$(cat ${basepath}/ver.txt)
Log "curVer:"$curVer

# Fetch upgrade information
DownloadApp $hmuUpDir $dlTarget $hmubootUrl $curVer
case $? in
    0)
        # not date to upgrade
        ;;
    1)
        # download done
        # Decompress
        Log "Installing new version"
        tar -xzf $dlTarget -C $tarTarget # 覆盖安装
        if [ $? -ne 0 ]; then
            Log "failed decompress:"$?
            Log "dlTarget:"$dlTarget",tarTarget:"$tarTarget
            UnlockFile $appMgrName
            return
        fi
        Log "update success, clean cached data"
        Log "clean:"$dlTarget
        rm -rf $dlTarget
        ;;
    2)
        # upgrade failed
        ;;
    *)
        Log "Unknow status:"$?
        ;;
esac

UnlockFile $appMgrName

# 如果存在指定的升级文件, 执行升级操作
if [ -f ${tarTarget}/upgrade.sh ]; then
    Log "Run upgrade.sh"
    ${tarTarget}/upgrade.sh
fi

