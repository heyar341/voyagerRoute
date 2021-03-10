# グーグる〜と

<https://googroutes.com>

## 概要

* グーグる〜とは、Google Maps API を使用し、一つの画面で複数ルートを検索できるアプリケーションです。

## 機能
* まとめ検索

　　最大10個までルートの検索、マップ上への表示ができます。
  
　　旅行の１日の移動スケジュールを組み立てる場合などに便利な機能です。

* 同時検索

　　１つの出発地から、最大9個の目的地までの距離、所要時間を同時に検索できる機能です。

　　いくつか行きたい場所があり、どこに1番早く行けるか知りたい場合などに便利な機能です。

(注意)APIの仕様上、日本国内の公共交通機関の乗り換え案内はご利用いただけません。
## 使用技術

### 言語

* Go-1.15.6
* Javascript
* HTML5 
* CSS3

### DB

* MongoDB-4.4.2

### 開発環境

* Docker, docker-compose

### インフラ
* ### GCP
    * App Engine
    * Cloud Build
    * Cloud Storage
    
* ### DB
    * MongoDB Atlas

### 使用API

#### Google Maps API

* Javascript API
* Directions API
* TimeZone API
* Places API

### 導入検討中技術
* Kubernetes

## 今後の追加予定機能

* Map上の地点をクリックして、出発地や目的地に設定できる機能
* Map上に表示されたルートのラインをクリックして、ルートを選択 できる機能
* 経由地を設定して、ルートを検索できる機能