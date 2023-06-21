challenge-2:
	go build -C ./unique/ -o dist-challenge
	maelstrom test -w unique-ids --bin ./unique/dist-challenge --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

challenge-3:
	go build -C ./broadcast/ -o dist-challenge
	maelstrom test -w broadcast --bin ./broadcast/dist-challenge --node-count 1 --time-limit 20 --rate 10