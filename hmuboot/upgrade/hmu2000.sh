#!/bin/sh

# ================================
#  配置区域开始
# ================================

appMgrName="aggregation" # 程序名称
aliveAddr="http://127.0.0.1:8090/hacheck" # 程序存活检查接口，需要返回200状态码，内容为1
verAddr="http://127.0.0.1:8090/version" # 程序版本号获取接口

hmuInDir=/usr/local/clc.hmu/app # 安装目录
mkdir -p $hmuInDir

# ========= upgrade url ============
#url1="https://upgrade-hmu.huayuan-iot.com/hmu/upgrade/hmu2000" # 应使用https，但目前条件未充足
url1="http://upgrade-hmu.huayuan-iot.com/hmu/upgrade/hmu2000"
url2=$url1
url3=$url1

# ================================
#  配置区域结束
# ================================

# TODO:暂时以此识别是否存在sd卡(测试时请手动创建)
if [ ! -d /mnt/sda1 ]; then
    echo "No sdcard"
    return 0
fi

# Clean and make tmp data
#=========================
basepath=$(cd `dirname $0`; pwd)
hmuUpDir=/mnt/sda1/${appMgrName}_upgrade # update directory
mkdir -p $hmuUpDir
rm -rf ${hmuUpDir}/* #**!!需要特别注意不能删根目录!!***#
dlTarget=${hmuUpDir}/${appMgrName}.tar.gz

hmuTarDir=${hmuInDir}/.tmpcache
mkdir -p $hmuTarDir
rm -rf ${hmuTarDir}/* #**!!需要特别注意不能删根目录!!***#

hmuBakDir=/mnt/sda1/bak/clc.hmu/app # backup directory
mkdir -p $hmuBakDir


# 导入工具包
. ${basepath}/lib/base.sh

# ==========================
LockFile $appMgrName
if [ $? -eq 1 ]; then
    # in locking
    return;
fi
# ============file lock end===========
Log "Clear log for upgrading">/tmp/${appMgrName}_log.txt

# ========= function define ============
checkAlive(){
    for i in $(seq 15);
    do
        curl --capath /usr/local/clc.hmu/hmuboot/ssl/ -s ${aliveAddr}>${hmuUpDir}/alive.txt
        alive=$(cat ${hmuUpDir}/alive.txt)
        if [ $? -eq 0 ];then
            if [ "$alive" = "1" ]; then
    	    Log "App online."
                return 1
            fi
        fi
        Log "alive fail:"$alive"=>"$i":"$?
        sleep 1 # wait 1 second to next
    done
    return 0
}

rollback(){
    rm -rf $hmuTarDir
    if [ -d ${hmuBakDir}/aggregation/ ];then
        /usr/local/clc.hmu/app/supd/supd ctl stop clc.hmu.app.aggregation
        /usr/local/clc.hmu/app/supd/supd ctl start clc.hmu.app.aggregation.bak
        checkAlive
        if [ $? -eq 1 ]; then
            Log "The backup is online."
        fi

        Log "Restore the old version"
        rm -rf ${hmuInDir}/aggregation
        cp -rf ${hmuBakDir}/aggregation ${hmuInDir}
        /usr/local/clc.hmu/app/supd/supd ctl stop clc.hmu.app.aggregation.bak
        /usr/local/clc.hmu/app/supd/supd ctl start clc.hmu.app.aggregation
        checkAlive
        if [ $? -eq 1 ]; then
            Log "rollback success"
            rm -rf ${hmuBakDir}/aggregation
        fi
    fi
}

# =============== main ================
curVer=""
curl -s ${verAddr}>${hmuUpDir}/curver.txt
if [ $? -ne 0 ]; then
    if [ -d ${hmuInDir}/aggregation/ ];then
        Log "The app not running, try restart."
        /usr/local/clc.hmu/app/supd/supd ctl start clc.hmu.app.aggregation
        checkAlive
        if [ $? -eq 1 ]; then
            curl -s ${verAddr}>${hmuUpDir}/curver.txt
            curVer=$(cat ${hmuUpDir}/curver.txt)
            Log "curVer:"$curVer
        else
            Log "The old one is online failed, clean and reinstall."
            cp -rf ${hmuInDir}/aggregation ${hmuBakDir}/aggregation
            rm -rf ${hmuInDir}/aggregation
        fi
    else
        Log "App is not installed, try to init."
    fi
else
    curVer=$(cat ${hmuUpDir}/curver.txt)
    Log "curVer:"$curVer
fi

# Fetch upgrade information
downloadDone=0
for val in $url1 $url2 $url3; do
    DownloadApp $hmuUpDir $dlTarget $val $curVer
    case $? in
        0)
            # not date to upgrade
            UnlockFile $appMgrName
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
    UnlockFile $appMgrName
    return;
fi


# backup the old link
if [ -d ${hmuInDir}/aggregation/ ];then
    Log "Backup previous version"
    cp -rf ${hmuInDir}/aggregation $hmuBakDir
    /usr/local/clc.hmu/app/supd/supd ctl stop clc.hmu.app.aggregation
    /usr/local/clc.hmu/app/supd/supd ctl start clc.hmu.app.aggregation.bak
    checkAlive
    if [ $? -ne 1 ];then
        Log "backup failed"
        /usr/local/clc.hmu/app/supd/supd ctl stop clc.hmu.app.aggregation.bak
        /usr/local/clc.hmu/app/supd/supd ctl restart clc.hmu.app.aggregation
        UnlockFile $appMgrName
        return
    fi

    # back success, delete the old to free space
    rm -rf ${hmuInDir}/aggregation
fi

# Decompress
Log "Installing new version"
tar -xzf $dlTarget -C $hmuTarDir
if [ $? -ne 0 ]; then
    Log "failed decompress:"$?
    Log "dlTarget:"$dlTarget
    Log "hmuTarDir:"$hmuTarDir
    rollback
    UnlockFile $appMgrName
    return
fi
tmpName=$(ls $hmuTarDir)
mv ${hmuTarDir}/$tmpName ${hmuInDir}/aggregation
if [ -f ${hmuBakDir}/aggregation/monitoring-units.json ]; then
    # Not upgrade the resource
    cp -rf ${hmuBakDir}/aggregation/*.json ${hmuInDir}/aggregation
    cp -rf ${hmuBakDir}/aggregation/element-lib/* ${hmuInDir}/aggregation/element-lib/
    cp -rf ${hmuBakDir}/aggregation/etc/* ${hmuInDir}/aggregation/etc/
fi

# execute replace
/usr/local/clc.hmu/app/supd/supd ctl stop clc.hmu.app.aggregation.bak
/usr/local/clc.hmu/app/supd/supd ctl start clc.hmu.app.aggregation
checkAlive
installSuc=$?
if [ ${installSuc} -eq 1 ]; then
    Log "update success, clean cached data"
    Log "clean:"${hmuBakDir}/aggregation
    rm -rf $hmuBakDir/aggregation
    Log "clean:"${dlTarget}
    rm -rf $dlTarget
    UnlockFile $appMgrName
    return;
else
    Log "upgrade failed, do rollback."
    rollback
    UnlockFile $appMgrName
    return;
fi

UnlockFile $appMgrName
