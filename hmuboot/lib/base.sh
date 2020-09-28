
basepath=$(pwd)
if [ ! -z "$hmuUpDir" ]; then
    mkdir -p $hmuUpDir
fi

Log(){
    # TODO:Log it to system log
    echo `date` $*
}

# ============file lock start===========
# if dead lock, reboot the device to clean /tmp
LockFile(){
	lockName=$1
	if [ -f /tmp/${lockName}_lock ];then
            #超时30分钟自动解锁
            lockTime=$(cat /tmp/${lockName}_lock)
            now=$(date +%s)
            dur=$((${now}-${lockTime}))
            if [ $dur -lt 1800 ]; then
    		echo "Inlocking! You can delete /tmp/${lockName}.lock to unlock it.Locked:"$dur
                return 1
            fi
	fi
	echo `date +%s`>/tmp/${lockName}.lock
	return 0

}
UnlockFile() {
	lockName=$1
	rm -rf /tmp/${lockName}.lock
}

GetMac(){
    mac="uuid="$(cat /tmp/dxs/snfile)"&" # for hmu2000
	for eth in `ls /sys/class/net/`; do
	    mac=${mac}${eth}"="$(cat /sys/class/net/${eth}/address)"&"
	done
	if [ -z "$mac" ]; then
	    return 2
	fi
	echo $mac|sed 's/.$//'
}

# Parse json value
GetJsonVal(){
    # Log $1 $2
    ${basepath}/boot/JSON.sh -b <$1 | egrep '\["'$2'"\]'|awk -F \" '{print $4}'
    return 0
}

DownloadApp(){
    hmuUpDir=$1
    dlTarget=$2
    verUrl=$3
    curVer=$4
    
    echo ""
    Log "try "$verUrl
    verResp="";
    params="ver="${curVer}"&"$(GetMac)
    curl --capath /usr/local/clc.hmu/hmuboot/ssl/ -d "${params}" -s "${verUrl}">${hmuUpDir}/ver.txt
    if [ $? -ne 0 ]; then
        # fail connection
        Log "failed connection:"$?
        Log "verUrl:"$verUrl
        Log "curVer:"$curVer
        return 2;
    fi
    
    # cat $hmuUpDir/ver.txt
    lastVer=$(GetJsonVal ${hmuUpDir}/ver.txt ver)
    if [ -z "$lastVer" ]; then
        Log "failed version:"$?
        Log "verUrl:"$verUrl
        Log "curVer:"$curVer
        Log "resp:"`cat ${hmuUpDir}/ver.txt`
        return 2;
    fi
    
    # diffent version
    if [ "$lastVer" = "$curVer" ]; then
        Log "The version is lastest, "$lastVer":"$curVer
        return 0
    fi
    
    dlUrl=$(GetJsonVal ${hmuUpDir}/ver.txt dl_url)
    msg=$(GetJsonVal ${hmuUpDir}/ver.txt msg)
    checksum=$(GetJsonVal ${hmuUpDir}/ver.txt checksum)
    Log "FOUND NEW APP:"`cat ${hmuUpDir}/ver.txt`
    
    # download the data
    curl -s --capath /usr/local/clc.hmu/hmuboot/ssl/ $dlUrl -o $dlTarget
    if [ $? -ne 0 ]; then
        Log "failed download:"$?
        Log "url:"$dlUrl
        Log "target:"$dlTarget
        return 2;
    fi
    
    # TODO:需要校验安装包由华远云联所发, 否则安装包有可能会被篡改过。
    targetMD5=$(/usr/bin/md5sum -b $dlTarget | cut -d ' ' -f1)
    if [ ! "$targetMD5" = "$checksum" ]; then
        Log "failed md5:"$?
        Log "MD5 Not Match:"$targetMD5":"$checksum
        return 2
    fi
    Log "file md5:"$targetMD5    
    
    return 1
}
