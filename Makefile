challenge_2:
	go build -C ./challenge-2/ -o dist-challenge
	maelstrom test -w unique-ids --bin ./challenge-2/dist-challenge --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

challenge_3a:
	go build -C ./challenge-3a/ -o dist-challenge
	maelstrom test -w broadcast --bin ./challenge-3a/dist-challenge --node-count 1 --time-limit 20 --rate 10

challenge_3b:
	go build -C ./challenge-3b/ -o dist-challenge
	maelstrom test -w broadcast --bin ./challenge-3b/dist-challenge --node-count 5 --time-limit 20 --rate 10

challenge_3c:
	go build -C ./challenge-3c/ -o dist-challenge
	maelstrom test -w broadcast --bin ./challenge-3c/dist-challenge --node-count 5 --time-limit 20 --rate 10 --nemesis partition
