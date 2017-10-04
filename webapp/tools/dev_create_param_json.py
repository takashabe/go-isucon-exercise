import mysql.connector
import json


class Generater:
    """ provide generate param.json for the benchmark"""

    def __init__(self):
        conn = mysql.connector.connect(
            host='localhost',
            port=3306,
            user='isucon',
            database='isucon', )
        conn.ping(reconnect=True)
        if not conn.is_connected():
            raise Exception('mysql', 'failed to connection')
        self.conn = conn

    def generate(self):
        cur = self.conn.cursor(dictionary=True)
        cur.execute('select id, name, email from user')
        res = cur.fetchall()
        dic = [{
            'id': x['id'],
            'name': x['name'],
            'email': x['email'],
            'password': x['name']
        } for x in res]
        param = {'parameters': dic}
        print(json.dumps(param, indent=4))


g = Generater()
g.generate()
