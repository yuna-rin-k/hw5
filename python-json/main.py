#!/usr/bin/env python
# -*- coding: utf-8 -*-

import webapp2
from google.appengine.api import urlfetch
import json
import jinja2

networkJson = urlfetch.fetch("https://tokyo.fantasy-transit.appspot.com/net?format=json").content  # ウェブサイトから電車の線路情報をJSON形式でダウンロードする
network = json.loads(networkJson.decode('utf-8'))  # JSONとしてパースする（stringからdictのlistに変換する）

tmpl = jinja2.Template(  # Jinjaのテンプレートエンジンを使ってHTMLを作ります
    u'''
こういう線路があります：
<ul>
{% for line in network %}
  <li> {{line["Name"]}}
    <ul>
    {% for station in line["Stations"] %}
      <li> {{station}} </li>
    {% endfor %}
  </li>
  </ul>
{% endfor %}
</ul>
''')

class MainPage(webapp2.RequestHandler):
    def get(self):
        self.response.headers['Content-Type'] = 'text/html; charset=UTF-8'
        self.response.write(tmpl.render(network=network))

app = webapp2.WSGIApplication([
    ('/', MainPage),
], debug=True)
