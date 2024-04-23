target1:
	@touch target1test
	@touch target2test

target1Clean:
	@rm -rf target*test

failTarget:
	@echo "FAIL"

.PHONY: failTarget target1Clean