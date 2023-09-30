# Check the connection status of FE query port.

trap exit TERM
host=$FE_SVC
port=$FE_QUERY_PORT

while true; do
  nc -zv -w 3 $host $port
  if [ $? -eq 0 ]; then
	  break
  else
	  echo "info: failed to connect to $host:$port, sleep 1 second then retry"
	  sleep 1
  fi
done

echo "info: successfully connected to $host:$port, able to initialize Doris now"

