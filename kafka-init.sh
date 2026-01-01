#!/usr/bin/env bash
set -euo pipefail

BOOTSTRAP="${BOOTSTRAP:-kafka1:9092}"
TOPIC_TASKS="${TOPIC_TASKS:-factorial.tasks}"
TOPIC_RESULTS="${TOPIC_RESULTS:-factorial.results}"
PARTITIONS="${PARTITIONS:-8}"
REPLICATION="${REPLICATION:-3}"
MIN_ISR="${MIN_ISR:-2}"
WAIT_SECONDS="${WAIT_SECONDS:-120}"

echo "[kafka-init] waiting for kafka (${BOOTSTRAP})..."
for i in $(seq 1 "${WAIT_SECONDS}"); do
  if /opt/kafka/bin/kafka-topics.sh --bootstrap-server "${BOOTSTRAP}" --list >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

/opt/kafka/bin/kafka-topics.sh --bootstrap-server "${BOOTSTRAP}" --list >/dev/null

echo "[kafka-init] creating topics (partitions=${PARTITIONS}, rf=${REPLICATION}, min.insync.replicas=${MIN_ISR})"
/opt/kafka/bin/kafka-topics.sh --bootstrap-server "${BOOTSTRAP}" --create --if-not-exists \
  --topic "${TOPIC_TASKS}" --partitions "${PARTITIONS}" --replication-factor "${REPLICATION}" \
  --config "min.insync.replicas=${MIN_ISR}"

/opt/kafka/bin/kafka-topics.sh --bootstrap-server "${BOOTSTRAP}" --create --if-not-exists \
  --topic "${TOPIC_RESULTS}" --partitions "${PARTITIONS}" --replication-factor "${REPLICATION}" \
  --config "min.insync.replicas=${MIN_ISR}"

echo "[kafka-init] done"

