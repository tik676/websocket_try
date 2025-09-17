kafka-topic --create \
 --if-not-exists \
 --bootstrap-server kafka1:29092 \
 --topic user_events \
 --partitions 3 \
 --replication-factor 2