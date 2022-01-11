#refer in zshrc
source ~/.scripts/secret

aocinput(){
	curl https://adventofcode.com/"$1"/day/"$2"/input --cookie "session=$AOC_SESSION_ID" > in
}
