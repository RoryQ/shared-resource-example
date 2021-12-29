
PHONY: emulator-up
emulator-up:
	docker run \
		-p 9010:9010 \
		-p 9020:9020 \
		--platform linux/amd64 \
		--env SPANNER_DATABASE_ID=db \
		--env SPANNER_INSTANCE_ID=inst \
		--env SPANNER_PROJECT_ID=proj \
		roryq/spanner-emulator
