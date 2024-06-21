#!/busybox/sh

echo
date
echo entered entrypoint.sh

/manager
ec=$?

echo
echo MANAGER TERMINATED!  ERROR=${ec} AT $(date)
echo

echo sleeping
sleep 9999
echo terminating 
