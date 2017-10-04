import mysql.connector
import string
import random
import hashlib


class Initializer:
    """ provide dummy user list insert to users table"""

    def __init__(self):
        conn = mysql.connector.connect(
            host='localhost',
            port=3306,
            user='isucon',
            database='isucon',
        )
        conn.ping(reconnect=True)
        if not conn.is_connected():
            raise Exception('mysql', 'failed to connection')
        conn.cursor().execute('TRUNCATE TABLE user')
        conn.commit()
        self.conn = conn

    def gen_salt(self, size):
        src = string.ascii_letters + string.digits
        return ''.join([random.choice(src) for i in range(size)])

    def gen_hash(self, seed):
        h = hashlib.sha256()
        h.update(seed.encode('utf-8'))
        return h.hexdigest()

    def register_users(self):
        with open('mail_list', 'r') as f:
            for line in f:
                email = line.strip()
                name = email.split('@')[0]
                salt = self.gen_salt(16)
                passhash = self.gen_hash(salt + name)

                cur = self.conn.cursor()
                cur.execute('select id from user order by id desc limit 1;')
                latest = cur.fetchone()
                insert_id = 1
                if latest is not None:
                    insert_id = latest[0] + 1

                cur.execute(
                    'insert into user (id, name, email, salt, passhash) values (%s, %s, %s, %s, %s)',
                    [insert_id, name, email, salt, passhash])
                self.conn.commit()

        f.closed


init = Initializer()
init.register_users()
