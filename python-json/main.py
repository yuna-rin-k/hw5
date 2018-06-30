#!/usr/bin/env python
# -*- coding: utf-8 -*-

import webapp2
from google.appengine.api import urlfetch
import json

networkJson = urlfetch.fetch("http://tokyo.fantasy-transit.appspot.com/net?format=json").content  # ウェブサイトから電車の線路情報をJSON形式でダウンロードする
network = json.loads(networkJson, 'utf-8')  # JSONとしてパースする（stringからdictのlistに変換する）

class MainPage(webapp2.RequestHandler):
    def get(self):
        self.response.headers['Content-Type'] = 'text/html; charset=UTF-8'
        self.response.write(u'こういう線路があります：<ul>')
        for line in network:
          self.response.write('<li>%s</li>' % line["Name"])
        self.response.write(u'</ul>')

app = webapp2.WSGIApplication([
    ('/', MainPage),
], debug=True)
