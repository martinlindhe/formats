status:
	echo "| What          | Status |" > STATUS.md
	echo "| ------------- | ------:|" >> STATUS.md
	grep -r STATUS parse/* | sed 's/parse\//| /' | sed 's/\/\/ STATUS:/ |/' | sed 's/%/% |/' | sed 's/://' | sort >> STATUS.md