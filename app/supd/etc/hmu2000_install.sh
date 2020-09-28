#!/bin/sh

writeNormalBoot(){
    cat << EOF > /etc/init.d/supd
#!/bin/sh /etc/rc.common
# Example script
# Copyright (C) 2007 OpenWrt.org
 
START=99
STOP=15
USE_PROCD=1
PROG=/usr/local/clc.hmu/app/supd/supd

error() {
	echo "\${initscript}:" "\$@" 1>&2
}               
 
start_service() {
    procd_open_instance
    procd_set_param command "\$PROG"
    procd_set_param respawn
    procd_append_param env GIN_MODE=release
    procd_append_param env PRJ_ROOT=/usr/local/clc.hmu
    procd_close_instance
    echo \$PROG "started"
}

stop_service() {  
    echo "stopped" \$PROG
} 

restart() {  
    stop  
    start  
}
EOF
chmod +x /etc/init.d/supd # make run able
/etc/init.d/supd enable # Enable service autostart
}

####### main start #########
writeNormalBoot
#######  main end  #########
