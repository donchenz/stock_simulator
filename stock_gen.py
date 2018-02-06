# -*- coding: UTF-8 -*-  


import random, threading, time, json
from redis import Redis

#股票数量
STOCK_NUM = 5
#刷新间隔
REFRESH_INTERVAL = 1
#
REDIS_HOST = '127.0.0.1'
#
REDIS_PORT = 6379
#redis 队列名字
REDIS_Q_NAME = 'stock_q'



redis = Redis(host=REDIS_HOST, port=REDIS_PORT)

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
		redis.lpush(REDIS_Q_NAME, stock_json)
		print stock_json

		time.sleep(REFRESH_INTERVAL)




if __name__ == '__main__':
	try:
		generate_data()
	except KeyboardInterrupt:
		print 'program exit.'




