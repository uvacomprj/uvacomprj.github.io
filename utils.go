package main

//CheckError checks for errors
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}