#!/home/liuyue01/local/expect/bin/expect -f

## Global Conf ###
set timeout 10
set host_conf_file "~/bin/host.conf"
##################

## Global Var ####
array set host_map [list] 
##################

###  function defination ###

proc parse_host_conf {} {
    global host_conf_file host_map
    if {![file exist $host_conf_file]} {
        puts stderr "Host conf file doesn't exist..."
        exit 1
    }
    set fp [open "$host_conf_file" r]
    while {[gets $fp line] >= 0} {
        set line [string trim "$line"]
        set length [string length $line]
        if {$length == 0} {
            continue
        }
        set first_tag [string index $line 0]
        if {$first_tag == "#"} {
            continue
        }
        set host_info [split $line ","]
        foreach {key host username password} $host_info break
        set host_map($key) [list $host $username $password]
    }
}

proc get_ssh_info {host_alias} {    
    global host_map
    set count [array size host_map]
    if {$count == 0} {
        puts "The host list is empty..."
        exit 0
    }
    set flag 0
    foreach {k v} [array get host_map *] {
        if {$k == $host_alias} {
            set flag 1    
        }
    }
    if {$flag != 1} {
        puts "The host key $host_alias is not exist..."
        exit 0
    }
    set host_info $host_map($host_alias)
    set host [lindex $host_info 0]
    set username [lindex $host_info 1]
    set passwd [lindex $host_info 2]
    set ssh_host "$username@$host"
    return [list $ssh_host $passwd]
}

proc get_host_list {} {
    global host_map
    set count [array size host_map]
    if {$count == 0} {
        puts "The host list is empty..."
        exit 0
    }
    foreach {k v} [array get host_map *] {
        puts "$k\t\t[lindex $v 0]" 
    }
}

############################

### Main() ###

set host_alias [lindex $argv 0]

if {$argc < 1} { 
    puts stderr "Usage: $argv0 {host_alias}"
    exit 1 
}

parse_host_conf

switch $host_alias {

    "list" {
        get_host_list
        exit 0
    }

    default {
        set ssh_info [get_ssh_info $host_alias]
    }
}

set ssh_host [lindex $ssh_info 0]
set passwd [lindex $ssh_info 1]

spawn ssh $ssh_host
expect {
    "*yes/no" { send "yes\r"; exp_continue}
    "*assword:" { send "$passwd\r" }
}
interact
