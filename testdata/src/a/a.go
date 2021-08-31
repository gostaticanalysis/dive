package a

func f() {
	if true { // want "too long block"
		_ = 10
		_ = 10
		_ = 10
		_ = 10
		_ = 10
		_ = 10
	}

	if true { // OK
		// comment
		// comment
		// comment
		// comment
		// comment
		_ = 10
	}

	if true { // want "too many returns in the block"
		if true {
			return
		}

		if true {
			return
		}

		if true {
			return
		}
	}

	if true { // want "loop in if block"
		for {
		}
	}

	if true { // want "loop in if block" "too deeply nest"
		if true {
			for {
			}
		}
	}

	if true { // want "too deeply nest"
		if true {
			if true {
				return
			}
		}
	}
}
