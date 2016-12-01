package tool

//模拟try
func Try(success func(), err func(e error)) {
	defer func() {
		if e := recover(); e != nil {
			err(e.(error))
		}
	}()
	success()
}
