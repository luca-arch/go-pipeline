#!/bin/sh

set +e

run_test() {
	echo "Running ${1}..."

	EXPECTED=$2
	docker run --rm --volume `pwd`/testdata:/mnt "gopipeline:alpine" /mnt/${1}
	ACTUAL=$?

	if [ "${ACTUAL}" == "${EXPECTED}" ]; then
		echo "SUCCESS: ${1} got ${EXPECTED}!"
		return;
	fi

	echo "ERROR: ${1} expected ${EXPECTED} but got ${ACTUAL}"
	exit 1
}

run_test test-pipeline-001.yaml 0
run_test test-pipeline-002.yaml 67
run_test test-pipeline-003.yaml 125
run_test test-pipeline-004.yaml 125
