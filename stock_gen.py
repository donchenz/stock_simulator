# -*- coding: UTF-8 -*-  


import random, threading, time, json, socket

#股票数量
STOCK_NUM = 5
#刷新间隔
REFRESH_INTERVAL = 1
#middleware地址
MCAST_ADDR = '224.0.0.1'
#middleware端口
MCAST_PORT = 9999

ANY = '0.0.0.0'
SENDERPORT=1501 



def get_random_num():
	return round(random.random() * 100, 2)

class Stock:
	def __init__(self):
		self.name = ''
		self.buy_price = 0.0
		self.sell_price = 0.0
		self.volume = 0.0

	@classmethod
	def gen_random(self, name):
		s = Stock()
		s.name = name
		s.buy_price = get_random_num()
		s.sell_price = get_random_num()
		s.volume = get_random_num()
		return s



def generate_data():
	print 'start to generate data...'
	stock_name_list = ['stock_' + str(i) for i in range(1, STOCK_NUM+1)]

	while True:
		refresh_stock_num = random.randint(0, STOCK_NUM)
		tmp_name_list = stock_name_list[:]
		random.shuffle(tmp_name_list)
		stock_list = []
		for i in range(0, refresh_stock_num):
			stock = Stock.gen_random(tmp_name_list[i])
			stock_list.append(stock.__dict__)

		stock_json = json.dumps(stock_list)
		send_data(stock_json)

		time.sleep(REFRESH_INTERVAL)

def send_data(data):
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP) 
	sock.bind((ANY,SENDERPORT)) 
	sock.setsockopt(socket.IPPROTO_IP, socket.IP_MULTICAST_TTL, 255)
	sock.sendto(data, (MCAST_ADDR,MCAST_PORT) )


if __name__ == '__main__':
	try:
		generate_data()
	except KeyboardInterrupt:
		print 'program exit.'




