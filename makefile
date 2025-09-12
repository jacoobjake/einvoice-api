.PHONY: bob-gen-pgsql

bob-gen-pgsql:
	go run github.com/stephenafamo/bob/gen/bobgen-psql@latest -c bobgen.yaml