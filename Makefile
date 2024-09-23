PKG := github.com/lightninglabs/btclog

GOLIST := go list -deps $(PKG)/... | grep '$(PKG)'| grep -v '/vendor/'
GOTEST := go test -v

XARGS := xargs -L 1
UNIT := $(GOLIST) | $(XARGS) env $(GOTEST) $(TEST_FLAGS)

unit:
	@$(call print, "Running unit tests.")
	$(UNIT)