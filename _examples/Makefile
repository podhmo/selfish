SELFISH ?= ../bin/selfish


create: $(SELFISH)
	$(SELFISH) data/*.md --debug
update: $(SELFISH)
	$(SELFISH) --alias head --silent data/* --debug
delete: $(SELFISH)
	$(SELFISH) --alias head --delete --debug
.PHONY: create update delete

$(SELFISH):
	$(MAKE) -C ../ build
