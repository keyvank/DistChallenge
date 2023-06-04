challenge-2:
	go build -o unique
	maelstrom test -w unique-ids --bin unique --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
