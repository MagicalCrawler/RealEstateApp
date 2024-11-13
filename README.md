# RealEstateApp
A telegram bot or CLI to search real estate ads

## Development Setup
 - You need a `.env` file in root of project. Clone `.env.example` as `.env` and define your environment variables.

 ## Test
 - All Tests placed in `test` directory
 - The `main_test.go` contains these methods:
   - `TestMain`: Run once before all tests.
   - `clearData`: This method clear all data, If your test need Data manipulation, call this method before and after of your test, like this:
     ```go
     func TestSomething() {
        clearData()
	    defer clearData()
        // Write your Test here
     }
     ```
 - To run all tests, run this command: `go test ./test/...`