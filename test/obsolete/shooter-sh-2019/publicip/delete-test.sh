#!/bin/bash
source ../setup.env

#num=0
#for NAME in "${CONNECT_NAMES[@]}"
#do
#        #ID=`curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/publicip?connection_name=${NAME} |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
#	ID=publicipt${num}-powerkim
#        curl -H "${AUTH}" -sX DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=${NAME} &
#
#	num=`expr $num + 1`
#done

TB_PUBLICIP_IDS=`curl -H "${AUTH}" -sX GET http://$TUMBLEBUG_IP:1323/tumblebug/ns/$NS_ID/resources/publicIp | jq -r '.publicIp[].id'`
#echo $TB_PUBLICIP_IDS | json_pp

if [ -n "$TB_PUBLICIP_IDS" ]
then
        #TB_PUBLICIP_IDS=`curl -H "${AUTH}" -sX GET http://$TUMBLEBUG_IP:1323/tumblebug/ns/$NS_ID/resources/publicIp | jq -r '.publicIp[].id'`
        for TB_PUBLICIP_ID in ${TB_PUBLICIP_IDS}
        do
                echo ....Delete ${TB_PUBLICIP_ID} ...
                curl -H "${AUTH}" -sX DELETE http://$TUMBLEBUG_IP:1323/tumblebug/ns/$NS_ID/resources/publicIp/${TB_PUBLICIP_ID} | json_pp
        done
else
        echo ....no publicIps found
fi
