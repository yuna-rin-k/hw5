#!/usr/bin/env python3

from flask import Flask, render_template, request
import json
import urlfetch

networkJson = urlfetch.fetch("http://tokyo.fantasy-transit.appspot.com/net?format=json").content  # ウェブサイトから電車の線路情報をJSON形式でダウンロードする
network = json.loads(networkJson.decode('utf-8'))  # JSONとしてパースする（stringからdictのlistに変換する）

app = Flask(__name__)
@app.route('/')
def root():
  return render_template('trains.html', network)

@app.route('/pata')
def pata():
  magic = request.args.get('x', '') + request.args.get('y', '')
  return render_template('pata.html', magic=magic)
