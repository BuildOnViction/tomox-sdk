---
title: Swagger Document for TomoDEX
language_tabs:
  - shell: cURL
  - node: request
  - go: GO
  - ruby: Ruby
  - python: Python
  - java: Java
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<h1 id="swagger-document-for-tomodex">Swagger Document for TomoDEX v1.0.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

TomoDEX API Document

<h1 id="swagger-document-for-tomodex-accounts">accounts</h1>

Account endpoints

## Find account by user address

<a id="opIdhandleGetAccount"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /account/{userAddress} \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/account/{userAddress}", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/account/{userAddress}',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/account/{userAddress}', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/account/{userAddress}");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /account/{userAddress}`

Returns a single account

<h3 id="find-account-by-user-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userAddress|path|string|true|Address of user to return|

> Example responses

> 200 Response

```json
{
  "address": "0xF7349C253FF7747Df661296E0859c44e974fb52E"
}
```

<h3 id="find-account-by-user-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Account](#schemaaccount)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid Address|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Account not found|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
None
</aside>

## [Deprecated] Find account's token balance by user address and token address

<a id="opIdhandleGetAccountTokenBalance"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /account/{userAddress}/{tokenAddress} \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/account/{userAddress}/{tokenAddress}", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/account/{userAddress}/{tokenAddress}',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/account/{userAddress}/{tokenAddress}', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/account/{userAddress}/{tokenAddress}");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /account/{userAddress}/{tokenAddress}`

Returns an object contains token balance of user

<h3 id="[deprecated]-find-account's-token-balance-by-user-address-and-token-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userAddress|path|string|true|Address of user to find token balance|
|tokenAddress|path|string|true|Address of token|

> Example responses

> 200 Response

```json
{
  "address": "string",
  "symbol": "string",
  "balance": "string",
  "availableBalance": "string",
  "inOrderBalance": "string"
}
```

<h3 id="[deprecated]-find-account's-token-balance-by-user-address-and-token-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[TokenBalance](#schematokenbalance)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid Address|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Account not found|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
None
</aside>

## Add a new account by user address

<a id="opIdhandleCreateAccount"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /account/create \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/account/create", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.post '/account/create',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.post('/account/create', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/account/create");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /account/create`

Returns newly created account

<h3 id="add-a-new-account-by-user-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|newAddress|path|string|true|Address of user|

> Example responses

> 201 Response

```json
{
  "address": "0xF7349C253FF7747Df661296E0859c44e974fb52E"
}
```

<h3 id="add-a-new-account-by-user-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Account already exists|None|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Account created|[Account](#schemaaccount)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Invalid Address|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-tokens">tokens</h1>

Token endpoints

## Finds all tokens

<a id="opIdHandleGetTokens"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /tokens \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/tokens", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/tokens',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/tokens', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/tokens");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /tokens`

Return all tokens in an array

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "name": "string",
    "symbol": "string",
    "address": "string",
    "image": {
      "url": "string",
      "meta": {}
    },
    "contractAddress": "string",
    "decimals": 0,
    "active": true,
    "listed": true,
    "quote": true,
    "makeFee": "string",
    "takeFee": "string",
    "usd": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="finds-all-tokens-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="finds-all-tokens-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Token](#schematoken)]|false|none|none|
|» id|string|false|read-only|none|
|» name|string|false|none|none|
|» symbol|string|false|none|none|
|» address|string|false|none|none|
|» image|[Image](#schemaimage)|false|none|none|
|»» url|string|false|none|none|
|»» meta|object|false|none|none|
|» contractAddress|string|false|none|none|
|» decimals|integer(int32)|false|read-only|none|
|» active|boolean|false|none|none|
|» listed|boolean|false|read-only|none|
|» quote|boolean|false|none|none|
|» makeFee|string|false|none|none|
|» takeFee|string|false|none|none|
|» usd|string|false|read-only|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Create new token

<a id="opIdHandleCreateToken"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /tokens \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/tokens", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/tokens',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post('/tokens', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/tokens");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /tokens`

Returns newly created token

> Body parameter

```json
{
  "name": "string",
  "symbol": "string",
  "address": "string",
  "image": {
    "url": "string",
    "meta": {}
  },
  "contractAddress": "string",
  "active": true,
  "quote": true,
  "makeFee": "string",
  "takeFee": "string"
}
```

<h3 id="create-new-token-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Token](#schematoken)|true|Token object that needs to be added|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "name": "string",
  "symbol": "string",
  "address": "string",
  "image": {
    "url": "string",
    "meta": {}
  },
  "contractAddress": "string",
  "decimals": 0,
  "active": true,
  "listed": true,
  "quote": true,
  "makeFee": "string",
  "takeFee": "string",
  "usd": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="create-new-token-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Token](#schematoken)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid payload|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Finds all base tokens

<a id="opIdHandleGetBaseTokens"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /tokens/base \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/tokens/base", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/tokens/base',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/tokens/base', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/tokens/base");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /tokens/base`

Return all base tokens in an array

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "name": "string",
    "symbol": "string",
    "address": "string",
    "image": {
      "url": "string",
      "meta": {}
    },
    "contractAddress": "string",
    "decimals": 0,
    "active": true,
    "listed": true,
    "quote": true,
    "makeFee": "string",
    "takeFee": "string",
    "usd": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="finds-all-base-tokens-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="finds-all-base-tokens-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Token](#schematoken)]|false|none|none|
|» id|string|false|read-only|none|
|» name|string|false|none|none|
|» symbol|string|false|none|none|
|» address|string|false|none|none|
|» image|[Image](#schemaimage)|false|none|none|
|»» url|string|false|none|none|
|»» meta|object|false|none|none|
|» contractAddress|string|false|none|none|
|» decimals|integer(int32)|false|read-only|none|
|» active|boolean|false|none|none|
|» listed|boolean|false|read-only|none|
|» quote|boolean|false|none|none|
|» makeFee|string|false|none|none|
|» takeFee|string|false|none|none|
|» usd|string|false|read-only|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Finds all quote tokens

<a id="opIdHandleGetQuoteTokens"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /tokens/quote \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/tokens/quote", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/tokens/quote',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/tokens/quote', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/tokens/quote");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /tokens/quote`

Return all quote tokens in an array

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "name": "string",
    "symbol": "string",
    "address": "string",
    "image": {
      "url": "string",
      "meta": {}
    },
    "contractAddress": "string",
    "decimals": 0,
    "active": true,
    "listed": true,
    "quote": true,
    "makeFee": "string",
    "takeFee": "string",
    "usd": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="finds-all-quote-tokens-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="finds-all-quote-tokens-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Token](#schematoken)]|false|none|none|
|» id|string|false|read-only|none|
|» name|string|false|none|none|
|» symbol|string|false|none|none|
|» address|string|false|none|none|
|» image|[Image](#schemaimage)|false|none|none|
|»» url|string|false|none|none|
|»» meta|object|false|none|none|
|» contractAddress|string|false|none|none|
|» decimals|integer(int32)|false|read-only|none|
|» active|boolean|false|none|none|
|» listed|boolean|false|read-only|none|
|» quote|boolean|false|none|none|
|» makeFee|string|false|none|none|
|» takeFee|string|false|none|none|
|» usd|string|false|read-only|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the token information corresponding to an address

<a id="opIdHandleGetToken"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /tokens/{address} \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/tokens/{address}", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/tokens/{address}',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/tokens/{address}', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/tokens/{address}");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /tokens/{address}`

Return token object

<h3 id="retrieve-the-token-information-corresponding-to-an-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|path|string|true|Token address|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "name": "string",
  "symbol": "string",
  "address": "string",
  "image": {
    "url": "string",
    "meta": {}
  },
  "contractAddress": "string",
  "decimals": 0,
  "active": true,
  "listed": true,
  "quote": true,
  "makeFee": "string",
  "takeFee": "string",
  "usd": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="retrieve-the-token-information-corresponding-to-an-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Token](#schematoken)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid Address|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-pairs">pairs</h1>

Pair endpoints

## Finds all pairs

<a id="opIdHandleGetPairs"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /pairs?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/pairs", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/pairs',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/pairs', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/pairs?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /pairs`

Return all pairs in an array

<h3 id="finds-all-pairs-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "baseTokenSymbol": "string",
    "baseTokenAddress": "string",
    "baseTokenDecimals": 0,
    "quoteTokenSymbol": "string",
    "quoteTokenAddress": "string",
    "quoteTokenDecimals": 0,
    "listed": true,
    "active": true,
    "rank": 0,
    "makeFee": "string",
    "takeFee": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="finds-all-pairs-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="finds-all-pairs-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Pair](#schemapair)]|false|none|none|
|» id|string|false|read-only|none|
|» baseTokenSymbol|string|false|none|none|
|» baseTokenAddress|string|false|none|none|
|» baseTokenDecimals|integer(int32)|false|read-only|none|
|» quoteTokenSymbol|string|false|none|none|
|» quoteTokenAddress|string|false|none|none|
|» quoteTokenDecimals|integer(int32)|false|read-only|none|
|» listed|boolean|false|read-only|none|
|» active|boolean|false|none|none|
|» rank|integer(int32)|false|read-only|none|
|» makeFee|string|false|read-only|none|
|» takeFee|string|false|read-only|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Create new pair

<a id="opIdHandleCreatePair"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /pairs \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/pairs", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/pairs',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post('/pairs', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/pairs");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /pairs`

Returns newly created pair

> Body parameter

```json
{
  "baseTokenSymbol": "string",
  "baseTokenAddress": "string",
  "quoteTokenSymbol": "string",
  "quoteTokenAddress": "string",
  "active": true
}
```

<h3 id="create-new-pair-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Pair](#schemapair)|true|Pair object that needs to be added|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "baseTokenSymbol": "string",
  "baseTokenAddress": "string",
  "baseTokenDecimals": 0,
  "quoteTokenSymbol": "string",
  "quoteTokenAddress": "string",
  "quoteTokenDecimals": 0,
  "listed": true,
  "active": true,
  "rank": 0,
  "makeFee": "string",
  "takeFee": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="create-new-pair-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Pair](#schemapair)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|*** Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the pair information corresponding to a baseToken and a quoteToken

<a id="opIdHandleGetPair"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /pair?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/pair", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/pair',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/pair', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/pair?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /pair`

Multiple status values can be provided with comma separated strings

<h3 id="retrieve-the-pair-information-corresponding-to-a-basetoken-and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "baseTokenSymbol": "string",
  "baseTokenAddress": "string",
  "baseTokenDecimals": 0,
  "quoteTokenSymbol": "string",
  "quoteTokenAddress": "string",
  "quoteTokenDecimals": 0,
  "listed": true,
  "active": true,
  "rank": 0,
  "makeFee": "string",
  "takeFee": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="retrieve-the-pair-information-corresponding-to-a-basetoken-and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Pair](#schemapair)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|baseToken Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve pair data corresponding to a baseToken and quoteToken

<a id="opIdHandleGetPairData"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /pair/data?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/pair/data", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/pair/data',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/pair/data', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/pair/data?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /pair/data`

Multiple status values can be provided with comma separated strings

<h3 id="retrieve-pair-data-corresponding-to-a-basetoken-and-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
{
  "pair": {
    "pairName": "string",
    "baseToken": "string",
    "quoteToken": "string"
  },
  "open": "string",
  "high": "string",
  "low": 0,
  "close": "string",
  "volume": "string",
  "count": "string",
  "timestamp": "string",
  "orderVolume": "string",
  "orderCount": "string",
  "averageOrderAmount": "string",
  "averageTradeAmount": "string",
  "askPrice": "string",
  "bidPrice": "string",
  "price": "string",
  "rank": 0
}
```

<h3 id="retrieve-pair-data-corresponding-to-a-basetoken-and-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[PairData](#schemapairdata)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|baseToken Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve all pair data

<a id="opIdHandleGetPairsData"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /pairs/data \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/pairs/data", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/pairs/data',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/pairs/data', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/pairs/data");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /pairs/data`

Multiple status values can be provided with comma separated strings

> Example responses

> 200 Response

```json
{
  "pair": {
    "pairName": "string",
    "baseToken": "string",
    "quoteToken": "string"
  },
  "open": "string",
  "high": "string",
  "low": 0,
  "close": "string",
  "volume": "string",
  "count": "string",
  "timestamp": "string",
  "orderVolume": "string",
  "orderCount": "string",
  "averageOrderAmount": "string",
  "averageTradeAmount": "string",
  "askPrice": "string",
  "bidPrice": "string",
  "price": "string",
  "rank": 0
}
```

<h3 id="retrieve-all-pair-data-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[PairData](#schemapairdata)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|baseToken Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-orders">orders</h1>

Order endpoints

## Retrieve the sorted list of orders for an Ethereum address

<a id="opIdhandleGetOrders"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orders \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orders", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orders',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orders', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orders`

Return all orders in an array

<h3 id="retrieve-the-sorted-list-of-orders-for-an-ethereum-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|false|User address|
|pageOffset|query|string|false|Page offset|
|pageSize|query|string|false|Number of items per a page|
|sortBy|query|string|false|Sort for query (time, orderStatus, orderType, orderSide)|
|sortType|query|string|false|asc/dec|
|orderStatus|query|string|false|OPEN/CANCELLED/FILLED/PARTIAL_FILLED|
|orderSide|query|string|false|SELL/BUY|
|orderType|query|string|false|LO/MO|
|baseToken|query|string|false|Base token address|
|quoteToken|query|string|false|Quote token address|
|from|query|string|false|the beginning timestamp (number of seconds from 1970/01/01) from which order data has to be queried|
|to|query|string|false|the ending timestamp ((number of seconds from 1970/01/01)) until which order data has to be queried|

> Example responses

> 200 Response

```json
{
  "total": 0,
  "orders": [
    {
      "id": "string",
      "userAddress": "string",
      "exchangeAddress": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "status": "string",
      "side": "string",
      "type": "string",
      "hash": "string",
      "signature": {
        "V": "string",
        "R": "string",
        "S": "string"
      },
      "pricepoint": "string",
      "amount": "string",
      "filledAmount": "string",
      "nonce": "string",
      "makeFee": "string",
      "takeFee": "string",
      "pairName": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}
```

<h3 id="retrieve-the-sorted-list-of-orders-for-an-ethereum-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-the-sorted-list-of-orders-for-an-ethereum-address-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» total|integer|false|none|none|
|» orders|[[Order](#schemaorder)]|false|none|none|
|»» id|string|false|read-only|none|
|»» userAddress|string|false|none|none|
|»» exchangeAddress|string|false|none|none|
|»» baseToken|string|false|none|none|
|»» quoteToken|string|false|none|none|
|»» status|string|false|none|none|
|»» side|string|false|none|none|
|»» type|string|false|none|none|
|»» hash|string|false|none|none|
|»» signature|[Signature](#schemasignature)|false|none|none|
|»»» V|string|false|none|none|
|»»» R|string|false|none|none|
|»»» S|string|false|none|none|
|»» pricepoint|string|false|none|none|
|»» amount|string|false|none|none|
|»» filledAmount|string|false|none|none|
|»» nonce|string|false|none|none|
|»» makeFee|string|false|none|none|
|»» takeFee|string|false|none|none|
|»» pairName|string|false|none|none|
|»» createdAt|string(date-time)|false|read-only|none|
|»» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Create new order

<a id="opIdHandleNewOrder"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /orders \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/orders", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/orders',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post('/orders', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /orders`

Returns newly created order

> Body parameter

```json
{
  "userAddress": "0x15e08dE16f534c890828F2a0D935433aF5B3CE0C",
  "exchangeAddress": "0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e",
  "baseToken": "0x4d7eA2cE949216D6b120f3AA10164173615A2b6C",
  "quoteToken": "0x0000000000000000000000000000000000000001",
  "side": "SELL/BUY",
  "type": "LO/MO",
  "status": "NEW/CANCELLED",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  },
  "pricepoint": "21207020000000000000000",
  "amount": "4693386710283129",
  "nonce": "1",
  "makeFee": "1",
  "takeFee": "1"
}
```

<h3 id="create-new-order-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[OrderCreate](#schemaordercreate)|true|Order object that needs to be added|

> Example responses

> 201 Response

```json
{
  "id": "string",
  "userAddress": "string",
  "exchangeAddress": "string",
  "baseToken": "string",
  "quoteToken": "string",
  "status": "string",
  "side": "string",
  "type": "string",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  },
  "pricepoint": "string",
  "amount": "string",
  "filledAmount": "string",
  "nonce": "string",
  "makeFee": "string",
  "takeFee": "string",
  "pairName": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="create-new-order-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|successful operation|[Order](#schemaorder)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid payload|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Account is blocked|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## [Deprecated] Retrieve the list of positions for an Ethereum address. Positions are order that have been sent to the matching engine and that are waiting to be matched

<a id="opIdhandleGetPositions"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orders/positions?address=string&limit=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orders/positions", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orders/positions',
  params: {
  'address' => 'string',
'limit' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orders/positions', params={
  'address': 'string',  'limit': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/positions?address=string&limit=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orders/positions`

Return all orders in an array

<h3 id="[deprecated]-retrieve-the-list-of-positions-for-an-ethereum-address.-positions-are-order-that-have-been-sent-to-the-matching-engine-and-that-are-waiting-to-be-matched-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|
|limit|query|string|true|Number of orders returned in query|

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "userAddress": "string",
    "exchangeAddress": "string",
    "baseToken": "string",
    "quoteToken": "string",
    "status": "string",
    "side": "string",
    "type": "string",
    "hash": "string",
    "signature": {
      "V": "string",
      "R": "string",
      "S": "string"
    },
    "pricepoint": "string",
    "amount": "string",
    "filledAmount": "string",
    "nonce": "string",
    "makeFee": "string",
    "takeFee": "string",
    "pairName": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="[deprecated]-retrieve-the-list-of-positions-for-an-ethereum-address.-positions-are-order-that-have-been-sent-to-the-matching-engine-and-that-are-waiting-to-be-matched-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="[deprecated]-retrieve-the-list-of-positions-for-an-ethereum-address.-positions-are-order-that-have-been-sent-to-the-matching-engine-and-that-are-waiting-to-be-matched-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Order](#schemaorder)]|false|none|none|
|» id|string|false|read-only|none|
|» userAddress|string|false|none|none|
|» exchangeAddress|string|false|none|none|
|» baseToken|string|false|none|none|
|» quoteToken|string|false|none|none|
|» status|string|false|none|none|
|» side|string|false|none|none|
|» type|string|false|none|none|
|» hash|string|false|none|none|
|» signature|[Signature](#schemasignature)|false|none|none|
|»» V|string|false|none|none|
|»» R|string|false|none|none|
|»» S|string|false|none|none|
|» pricepoint|string|false|none|none|
|» amount|string|false|none|none|
|» filledAmount|string|false|none|none|
|» nonce|string|false|none|none|
|» makeFee|string|false|none|none|
|» takeFee|string|false|none|none|
|» pairName|string|false|none|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the list of filled order for an Ethereum address

<a id="opIdhandleGetOrderHistory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orders/history?address=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orders/history", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orders/history',
  params: {
  'address' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orders/history', params={
  'address': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/history?address=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orders/history`

Return all orders in an array

<h3 id="retrieve-the-list-of-filled-order-for-an-ethereum-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|
|pageOffset|query|string|false|Page offset, default 0|
|pageSize|query|string|false|Number of items per a page, defaul 50|
|sortBy|query|string|false|Sort for query (time(default), orderStatus, orderType, orderSide)|
|sortType|query|string|false|asc/dec, default asc|
|orderStatus|query|string|false|OPEN/CANCELLED/FILLED/PARTIAL_FILLED|
|orderSide|query|string|false|SELL/BUY|
|orderType|query|string|false|LO/MO|
|baseToken|query|string|false|Base token address|
|quoteToken|query|string|false|Quote token address|
|from|query|string|false|the beginning timestamp (number of seconds from 1970/01/01) from which order data has to be queried|
|to|query|string|false|the ending timestamp ((number of seconds from 1970/01/01)) until which order data has to be queried|

> Example responses

> 200 Response

```json
{
  "total": 0,
  "orders": [
    {
      "id": "string",
      "userAddress": "string",
      "exchangeAddress": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "status": "string",
      "side": "string",
      "type": "string",
      "hash": "string",
      "signature": {
        "V": "string",
        "R": "string",
        "S": "string"
      },
      "pricepoint": "string",
      "amount": "string",
      "filledAmount": "string",
      "nonce": "string",
      "makeFee": "string",
      "takeFee": "string",
      "pairName": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}
```

<h3 id="retrieve-the-list-of-filled-order-for-an-ethereum-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-the-list-of-filled-order-for-an-ethereum-address-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» total|integer|false|none|none|
|» orders|[[Order](#schemaorder)]|false|none|none|
|»» id|string|false|read-only|none|
|»» userAddress|string|false|none|none|
|»» exchangeAddress|string|false|none|none|
|»» baseToken|string|false|none|none|
|»» quoteToken|string|false|none|none|
|»» status|string|false|none|none|
|»» side|string|false|none|none|
|»» type|string|false|none|none|
|»» hash|string|false|none|none|
|»» signature|[Signature](#schemasignature)|false|none|none|
|»»» V|string|false|none|none|
|»»» R|string|false|none|none|
|»»» S|string|false|none|none|
|»» pricepoint|string|false|none|none|
|»» amount|string|false|none|none|
|»» filledAmount|string|false|none|none|
|»» nonce|string|false|none|none|
|»» makeFee|string|false|none|none|
|»» takeFee|string|false|none|none|
|»» pairName|string|false|none|none|
|»» createdAt|string(date-time)|false|read-only|none|
|»» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the total number of orders for an Ethereum address

<a id="opIdhandleGetCountOrder"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orders/count?address=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orders/count", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orders/count',
  params: {
  'address' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orders/count', params={
  'address': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/count?address=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orders/count`

Return a positive integer

<h3 id="retrieve-the-total-number-of-orders-for-an-ethereum-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|

> Example responses

> 200 Response

```json
0
```

<h3 id="retrieve-the-total-number-of-orders-for-an-ethereum-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|integer|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve order nonce for an Ethereum address

<a id="opIdhandleGetOrderNonce"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orders/nonce?address=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orders/nonce", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orders/nonce',
  params: {
  'address' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orders/nonce', params={
  'address': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/nonce?address=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orders/nonce`

Return a positive integer

<h3 id="retrieve-order-nonce-for-an-ethereum-address-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|

> Example responses

> 200 Response

```json
0
```

<h3 id="retrieve-order-nonce-for-an-ethereum-address-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|integer|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Cancel order

<a id="opIdHandleCancelOrder"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /orders/cancel \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/orders/cancel", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.post '/orders/cancel',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.post('/orders/cancel', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/cancel");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /orders/cancel`

Returns the hash of cancelled order

> Body parameter

```json
{
  "orderHash": "string",
  "nonce": "string",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  }
}
```

<h3 id="cancel-order-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[OrderCancel](#schemaordercancel)|true|Cancel order object|

> Example responses

> 200 Response

```json
"string"
```

<h3 id="cancel-order-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|string|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid payload|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## [Deprecated] Cancel all orders

<a id="opIdhandleCancelAllOrders"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /orders/cancelAll?address=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("POST", "/orders/cancelAll", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.post '/orders/cancelAll',
  params: {
  'address' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.post('/orders/cancelAll', params={
  'address': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orders/cancelAll?address=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("POST");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`POST /orders/cancelAll`

This endpoint should implements signature authentication

<h3 id="[deprecated]-cancel-all-orders-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|

> Example responses

> 200 Response

```json
"string"
```

<h3 id="[deprecated]-cancel-all-orders-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|string|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid payload|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-orderbook">orderbook</h1>

Order book endpoints

## Retrieve the orderbook (amount and pricepoint) corresponding to a a baseToken and a quoteToken

<a id="opIdHandleGetOrderBook"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orderbook?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orderbook", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orderbook',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orderbook', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orderbook?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orderbook`

Multiple status values can be provided with comma separated strings

<h3 id="retrieve-the-orderbook-(amount-and-pricepoint)-corresponding-to-a-a-basetoken-and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
{
  "pairName": "string",
  "asks": [
    {
      "amount": "string",
      "pricepoint": "string"
    }
  ],
  "bids": [
    {
      "amount": "string",
      "pricepoint": "string"
    }
  ]
}
```

<h3 id="retrieve-the-orderbook-(amount-and-pricepoint)-corresponding-to-a-a-basetoken-and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[OrderBook](#schemaorderbook)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|*** Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the orderbook (full raw orders, including fields such as hashes, maker, taker addresses, signatures, etc.)
corresponding to a baseToken and a quoteToken

<a id="opIdHandleGetRawOrderBook"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /orderbook/raw?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/orderbook/raw", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/orderbook/raw',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/orderbook/raw', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/orderbook/raw?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /orderbook/raw`

Multiple status values can be provided with comma separated strings

<h3 id="retrieve-the-orderbook-(full-raw-orders,-including-fields-such-as-hashes,-maker,-taker-addresses,-signatures,-etc.)
corresponding-to-a-basetoken-and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
{
  "pairName": "string",
  "orders": [
    {
      "id": "string",
      "userAddress": "string",
      "exchangeAddress": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "status": "string",
      "side": "string",
      "type": "string",
      "hash": "string",
      "signature": {
        "V": "string",
        "R": "string",
        "S": "string"
      },
      "pricepoint": "string",
      "amount": "string",
      "filledAmount": "string",
      "nonce": "string",
      "makeFee": "string",
      "takeFee": "string",
      "pairName": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}
```

<h3 id="retrieve-the-orderbook-(full-raw-orders,-including-fields-such-as-hashes,-maker,-taker-addresses,-signatures,-etc.)
corresponding-to-a-basetoken-and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[RawOrderBook](#schemaraworderbook)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|*** Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-trades">trades</h1>

Trade endpoints

## Retrieve all trades corresponding to a baseToken or/and a quoteToken

<a id="opIdHandleGetTrades"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /trades \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/trades", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/trades',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/trades', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/trades");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /trades`

Return all trades in an array with total match

<h3 id="retrieve-all-trades-corresponding-to-a-basetoken-or/and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|false|Base token address|
|quoteToken|query|string|false|Quote token address|
|pageOffset|query|string|false|none|
|pageSize|query|string|false|number of trade item per page, default 50|
|sortBy|query|string|false|Sort for query (suported sort by time)|
|sortType|query|string|false|asc/dec|
|from|query|string|false|the beginning timestamp (number of seconds from 1970/01/01) from which order data has to be queried|
|to|query|string|false|the ending timestamp ((number of seconds from 1970/01/01)) until which order data has to be queried|

> Example responses

> 200 Response

```json
{
  "total": 0,
  "trades": [
    {
      "id": "string",
      "taker": "string",
      "maker": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "makerOrderHash": "string",
      "takerOrderHash": "string",
      "hash": "string",
      "txHash": "string",
      "pairName": "string",
      "pricepoint": "string",
      "amount": "string",
      "status": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}
```

<h3 id="retrieve-all-trades-corresponding-to-a-basetoken-or/and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|*** Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-all-trades-corresponding-to-a-basetoken-or/and-a-quotetoken-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» total|integer|false|none|none|
|» trades|[[Trade](#schematrade)]|false|none|none|
|»» id|string|false|read-only|none|
|»» taker|string|false|none|none|
|»» maker|string|false|none|none|
|»» baseToken|string|false|none|none|
|»» quoteToken|string|false|none|none|
|»» makerOrderHash|string|false|none|none|
|»» takerOrderHash|string|false|none|none|
|»» hash|string|false|none|none|
|»» txHash|string|false|none|none|
|»» pairName|string|false|none|none|
|»» pricepoint|string|false|none|none|
|»» amount|string|false|none|none|
|»» status|string|false|none|none|
|»» createdAt|string(date-time)|false|read-only|none|
|»» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve the sorted list of trades for an Ethereum address in which the given address is either maker or taker

<a id="opIdHandleGetTradesHistory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /trades/history?address=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/trades/history", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/trades/history',
  params: {
  'address' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/trades/history', params={
  'address': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/trades/history?address=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /trades/history`

Return trades array

<h3 id="retrieve-the-sorted-list-of-trades-for-an-ethereum-address-in-which-the-given-address-is-either-maker-or-taker-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|address|query|string|true|User address|
|pageOffset|query|string|false|none|
|pageSize|query|string|false|number of trade item per page, default 50|
|sortBy|query|string|false|Sort for query (suported sort by time)|
|sortType|query|string|false|asc/dec|
|baseToken|query|string|false|Base token address|
|quoteToken|query|string|false|Quote token address|
|from|query|string|false|the beginning timestamp (number of seconds from 1970/01/01) from which order data has to be queried|
|to|query|string|false|the ending timestamp ((number of seconds from 1970/01/01)) until which order data has to be queried|

> Example responses

> 200 Response

```json
{
  "total": 0,
  "trades": [
    {
      "id": "string",
      "taker": "string",
      "maker": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "makerOrderHash": "string",
      "takerOrderHash": "string",
      "hash": "string",
      "txHash": "string",
      "pairName": "string",
      "pricepoint": "string",
      "amount": "string",
      "status": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}
```

<h3 id="retrieve-the-sorted-list-of-trades-for-an-ethereum-address-in-which-the-given-address-is-either-maker-or-taker-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|address Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-the-sorted-list-of-trades-for-an-ethereum-address-in-which-the-given-address-is-either-maker-or-taker-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» total|integer|false|none|none|
|» trades|[[Trade](#schematrade)]|false|none|none|
|»» id|string|false|read-only|none|
|»» taker|string|false|none|none|
|»» maker|string|false|none|none|
|»» baseToken|string|false|none|none|
|»» quoteToken|string|false|none|none|
|»» makerOrderHash|string|false|none|none|
|»» takerOrderHash|string|false|none|none|
|»» hash|string|false|none|none|
|»» txHash|string|false|none|none|
|»» pairName|string|false|none|none|
|»» pricepoint|string|false|none|none|
|»» amount|string|false|none|none|
|»» status|string|false|none|none|
|»» createdAt|string(date-time)|false|read-only|none|
|»» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-ohlcv">ohlcv</h1>

OHLCV endpoints

## Retrieve OHLCV data corresponding to a baseToken and a quoteToken

<a id="opIdHandleGetOHLCV"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /ohlcv?baseToken=string&quoteToken=string&timeInterval=string&from=string&to=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/ohlcv", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/ohlcv',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string',
'timeInterval' => 'string',
'from' => 'string',
'to' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/ohlcv', params={
  'baseToken': 'string',  'quoteToken': 'string',  'timeInterval': 'string',  'from': 'string',  'to': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/ohlcv?baseToken=string&quoteToken=string&timeInterval=string&from=string&to=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /ohlcv`

Return all ticks in an array

<h3 id="retrieve-ohlcv-data-corresponding-to-a-basetoken-and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|
|timeInterval|query|string|true|Time interval, candle size. Valid values: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 1w, 1mo (1 month)|
|from|query|string|true|the beginning timestamp (number of seconds from 1970/01/01) from which ohlcv data has to be queried|
|to|query|string|true|the ending timestamp ((number of seconds from 1970/01/01)) until which ohlcv data has to be queried|

> Example responses

> 200 Response

```json
[
  {
    "pair": {
      "pairName": "string",
      "baseToken": "string",
      "quoteToken": "string"
    },
    "open": "string",
    "high": "string",
    "low": 0,
    "close": "string",
    "volume": "string",
    "count": "string",
    "timestamp": "string"
  }
]
```

<h3 id="retrieve-ohlcv-data-corresponding-to-a-basetoken-and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|*** Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-ohlcv-data-corresponding-to-a-basetoken-and-a-quotetoken-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Tick](#schematick)]|false|none|none|
|» pair|[PairID](#schemapairid)|false|none|none|
|»» pairName|string|false|none|none|
|»» baseToken|string|false|none|none|
|»» quoteToken|string|false|none|none|
|» open|string|false|none|none|
|» high|string|false|none|none|
|» low|integer(int32)|false|none|none|
|» close|string|false|none|none|
|» volume|string|false|none|none|
|» count|string|false|none|none|
|» timestamp|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-notifications">notifications</h1>

Notification endpoints

## Retrieve the list of notifications for an address with pagination

<a id="opIdHandleGetNotifications"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /notifications?userAddress=string&page=string&perPage=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/notifications", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/notifications',
  params: {
  'userAddress' => 'string',
'page' => 'string',
'perPage' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/notifications', params={
  'userAddress': 'string',  'page': 'string',  'perPage': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/notifications?userAddress=string&page=string&perPage=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /notifications`

Return notifications in an array

<h3 id="retrieve-the-list-of-notifications-for-an-address-with-pagination-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userAddress|query|string|true|User address|
|page|query|string|true|Page number|
|perPage|query|string|true|the number of records returned per page. Valid values are 10, 20, 30, 40, 50|

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "recipient": "string",
    "message": "string",
    "type": "string",
    "status": "string",
    "createdAt": "2019-11-03T15:35:48Z",
    "updatedAt": "2019-11-03T15:35:48Z"
  }
]
```

<h3 id="retrieve-the-list-of-notifications-for-an-address-with-pagination-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid user address|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<h3 id="retrieve-the-list-of-notifications-for-an-address-with-pagination-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Notification](#schemanotification)]|false|none|none|
|» id|string|false|read-only|none|
|» recipient|string|false|none|none|
|» message|string|false|none|none|
|» type|string|false|none|none|
|» status|string|false|none|none|
|» createdAt|string(date-time)|false|read-only|none|
|» updatedAt|string(date-time)|false|read-only|none|

<aside class="success">
This operation does not require authentication
</aside>

## Update status of a notification from UNREAD to READ

<a id="opIdHandleNewOrder"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /notifications/{id} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Content-Type": []string{"application/json"},
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("PUT", "/notifications/{id}", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Content-Type' => 'application/json',
  'Accept' => 'application/json'
}

result = RestClient.put '/notifications/{id}',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Content-Type': 'application/json',
  'Accept': 'application/json'
}

r = requests.put('/notifications/{id}', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/notifications/{id}");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("PUT");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`PUT /notifications/{id}`

Returns newly updated notification

> Body parameter

```json
{
  "recipient": "string",
  "message": "string",
  "type": "string",
  "status": "string"
}
```

<h3 id="update-status-of-a-notification-from-unread-to-read-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Notification](#schemanotification)|true|Notification object that needs to be updated|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "recipient": "string",
  "message": "string",
  "type": "string",
  "status": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="update-status-of-a-notification-from-unread-to-read-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Notification](#schemanotification)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Invalid payload|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-info">info</h1>

Info endpoints

## get__info

> Code samples

```shell
# You can also use wget
curl -X GET /info \
  -H 'Accept: */*'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"*/*"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/info", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => '*/*'
}

result = RestClient.get '/info',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': '*/*'
}

r = requests.get('/info', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/info");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /info`

> Example responses

> 200 Response

<h3 id="get__info-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal server error|None|

<h3 id="get__info-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» exchangeAddress|string|false|none|none|
|» fee|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## get__info_exchange

> Code samples

```shell
# You can also use wget
curl -X GET /info/exchange \
  -H 'Accept: */*'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"*/*"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/info/exchange", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => '*/*'
}

result = RestClient.get '/info/exchange',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': '*/*'
}

r = requests.get('/info/exchange', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/info/exchange");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /info/exchange`

> Example responses

> 200 Response

<h3 id="get__info_exchange-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|

<h3 id="get__info_exchange-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» exchangeAddress|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## get__info_fees

> Code samples

```shell
# You can also use wget
curl -X GET /info/fees \
  -H 'Accept: */*'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"*/*"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/info/fees", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => '*/*'
}

result = RestClient.get '/info/fees',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': '*/*'
}

r = requests.get('/info/fees', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/info/fees");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /info/fees`

> Example responses

> 200 Response

<h3 id="get__info_fees-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|Inline|

<h3 id="get__info_fees-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» fee|string|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="swagger-document-for-tomodex-market">market</h1>

## Retrieve market stats 24h corresponding to a baseToken and a quoteToken

> Code samples

```shell
# You can also use wget
curl -X GET /market/stats?baseToken=string&quoteToken=string \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/market/stats", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/market/stats',
  params: {
  'baseToken' => 'string',
'quoteToken' => 'string'
}, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/market/stats', params={
  'baseToken': 'string',  'quoteToken': 'string'
}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/market/stats?baseToken=string&quoteToken=string");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /market/stats`

Multiple status values can be provided with comma separated strings

<h3 id="retrieve-market-stats-24h-corresponding-to-a-basetoken-and-a-quotetoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|baseToken|query|string|true|Base token address|
|quoteToken|query|string|true|Quote token address|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "baseTokenSymbol": "string",
  "baseTokenAddress": "string",
  "baseTokenDecimals": 0,
  "quoteTokenSymbol": "string",
  "quoteTokenAddress": "string",
  "quoteTokenDecimals": 0,
  "listed": true,
  "active": true,
  "rank": 0,
  "makeFee": "string",
  "takeFee": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}
```

<h3 id="retrieve-market-stats-24h-corresponding-to-a-basetoken-and-a-quotetoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[Pair](#schemapair)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|baseToken Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## Retrieve all market stats 2

> Code samples

```shell
# You can also use wget
curl -X GET /market/stats/all \
  -H 'Accept: application/json'

```

```go
package main

import (
       "bytes"
       "net/http"
)

func main() {

    headers := map[string][]string{
        "Accept": []string{"application/json"},
        
    }

    data := bytes.NewBuffer([]byte{jsonReq})
    req, err := http.NewRequest("GET", "/market/stats/all", data)
    req.Header = headers

    client := &http.Client{}
    resp, err := client.Do(req)
    // ...
}

```

```ruby
require 'rest-client'
require 'json'

headers = {
  'Accept' => 'application/json'
}

result = RestClient.get '/market/stats/all',
  params: {
  }, headers: headers

p JSON.parse(result)

```

```python
import requests
headers = {
  'Accept': 'application/json'
}

r = requests.get('/market/stats/all', params={

}, headers = headers)

print r.json()

```

```java
URL obj = new URL("/market/stats/all");
HttpURLConnection con = (HttpURLConnection) obj.openConnection();
con.setRequestMethod("GET");
int responseCode = con.getResponseCode();
BufferedReader in = new BufferedReader(
    new InputStreamReader(con.getInputStream()));
String inputLine;
StringBuffer response = new StringBuffer();
while ((inputLine = in.readLine()) != null) {
    response.append(inputLine);
}
in.close();
System.out.println(response.toString());

```

`GET /market/stats/all`

Multiple status values can be provided with comma separated strings

> Example responses

> 200 Response

```json
{
  "pair": {
    "pairName": "string",
    "baseToken": "string",
    "quoteToken": "string"
  },
  "open": "string",
  "high": "string",
  "low": 0,
  "close": "string",
  "volume": "string",
  "count": "string",
  "timestamp": "string",
  "orderVolume": "string",
  "orderCount": "string",
  "averageOrderAmount": "string",
  "averageTradeAmount": "string",
  "askPrice": "string",
  "bidPrice": "string",
  "price": "string",
  "rank": 0
}
```

<h3 id="retrieve-all-market-stats-2-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|successful operation|[PairData](#schemapairdata)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|baseToken Parameter missing|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocSaccount">Account</h2>

<a id="schemaaccount"></a>

```json
{
  "address": "0xF7349C253FF7747Df661296E0859c44e974fb52E"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|address|string|false|none|none|
|tokenBalances|object|false|none|none|
|» address|string|false|none|none|
|» symbol|string|false|none|none|
|» balance|string|false|none|none|
|» availableBalance|string|false|none|none|
|» inOrderBalance|string|false|none|none|
|isBlocked|boolean|false|none|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocStokenbalance">TokenBalance</h2>

<a id="schematokenbalance"></a>

```json
{
  "address": "string",
  "symbol": "string",
  "balance": "string",
  "availableBalance": "string",
  "inOrderBalance": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string|false|none|none|
|symbol|string|false|none|none|
|balance|string|false|none|none|
|availableBalance|string|false|none|none|
|inOrderBalance|string|false|none|none|

<h2 id="tocStoken">Token</h2>

<a id="schematoken"></a>

```json
{
  "id": "string",
  "name": "string",
  "symbol": "string",
  "address": "string",
  "image": {
    "url": "string",
    "meta": {}
  },
  "contractAddress": "string",
  "decimals": 0,
  "active": true,
  "listed": true,
  "quote": true,
  "makeFee": "string",
  "takeFee": "string",
  "usd": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|name|string|false|none|none|
|symbol|string|false|none|none|
|address|string|false|none|none|
|image|[Image](#schemaimage)|false|none|none|
|contractAddress|string|false|none|none|
|decimals|integer(int32)|false|read-only|none|
|active|boolean|false|none|none|
|listed|boolean|false|read-only|none|
|quote|boolean|false|none|none|
|makeFee|string|false|none|none|
|takeFee|string|false|none|none|
|usd|string|false|read-only|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocSpair">Pair</h2>

<a id="schemapair"></a>

```json
{
  "id": "string",
  "baseTokenSymbol": "string",
  "baseTokenAddress": "string",
  "baseTokenDecimals": 0,
  "quoteTokenSymbol": "string",
  "quoteTokenAddress": "string",
  "quoteTokenDecimals": 0,
  "listed": true,
  "active": true,
  "rank": 0,
  "makeFee": "string",
  "takeFee": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|baseTokenSymbol|string|false|none|none|
|baseTokenAddress|string|false|none|none|
|baseTokenDecimals|integer(int32)|false|read-only|none|
|quoteTokenSymbol|string|false|none|none|
|quoteTokenAddress|string|false|none|none|
|quoteTokenDecimals|integer(int32)|false|read-only|none|
|listed|boolean|false|read-only|none|
|active|boolean|false|none|none|
|rank|integer(int32)|false|read-only|none|
|makeFee|string|false|read-only|none|
|takeFee|string|false|read-only|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocSpairid">PairID</h2>

<a id="schemapairid"></a>

```json
{
  "pairName": "string",
  "baseToken": "string",
  "quoteToken": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pairName|string|false|none|none|
|baseToken|string|false|none|none|
|quoteToken|string|false|none|none|

<h2 id="tocSpairdata">PairData</h2>

<a id="schemapairdata"></a>

```json
{
  "pair": {
    "pairName": "string",
    "baseToken": "string",
    "quoteToken": "string"
  },
  "open": "string",
  "high": "string",
  "low": 0,
  "close": "string",
  "volume": "string",
  "count": "string",
  "timestamp": "string",
  "orderVolume": "string",
  "orderCount": "string",
  "averageOrderAmount": "string",
  "averageTradeAmount": "string",
  "askPrice": "string",
  "bidPrice": "string",
  "price": "string",
  "rank": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pair|[PairID](#schemapairid)|false|none|none|
|open|string|false|none|none|
|high|string|false|none|none|
|low|integer(int32)|false|none|none|
|close|string|false|none|none|
|volume|string|false|none|none|
|count|string|false|none|none|
|timestamp|string|false|none|none|
|orderVolume|string|false|none|none|
|orderCount|string|false|none|none|
|averageOrderAmount|string|false|none|none|
|averageTradeAmount|string|false|none|none|
|askPrice|string|false|none|none|
|bidPrice|string|false|none|none|
|price|string|false|none|none|
|rank|integer(int32)|false|none|none|

<h2 id="tocSimage">Image</h2>

<a id="schemaimage"></a>

```json
{
  "url": "string",
  "meta": {}
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|url|string|false|none|none|
|meta|object|false|none|none|

<h2 id="tocSorder">Order</h2>

<a id="schemaorder"></a>

```json
{
  "id": "string",
  "userAddress": "string",
  "exchangeAddress": "string",
  "baseToken": "string",
  "quoteToken": "string",
  "status": "string",
  "side": "string",
  "type": "string",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  },
  "pricepoint": "string",
  "amount": "string",
  "filledAmount": "string",
  "nonce": "string",
  "makeFee": "string",
  "takeFee": "string",
  "pairName": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|userAddress|string|false|none|none|
|exchangeAddress|string|false|none|none|
|baseToken|string|false|none|none|
|quoteToken|string|false|none|none|
|status|string|false|none|none|
|side|string|false|none|none|
|type|string|false|none|none|
|hash|string|false|none|none|
|signature|[Signature](#schemasignature)|false|none|none|
|pricepoint|string|false|none|none|
|amount|string|false|none|none|
|filledAmount|string|false|none|none|
|nonce|string|false|none|none|
|makeFee|string|false|none|none|
|takeFee|string|false|none|none|
|pairName|string|false|none|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocSordercreate">OrderCreate</h2>

<a id="schemaordercreate"></a>

```json
{
  "userAddress": "0x15e08dE16f534c890828F2a0D935433aF5B3CE0C",
  "exchangeAddress": "0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e",
  "baseToken": "0x4d7eA2cE949216D6b120f3AA10164173615A2b6C",
  "quoteToken": "0x0000000000000000000000000000000000000001",
  "side": "SELL/BUY",
  "type": "LO/MO",
  "status": "NEW/CANCELLED",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  },
  "pricepoint": "21207020000000000000000",
  "amount": "4693386710283129",
  "nonce": "1",
  "makeFee": "1",
  "takeFee": "1"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|userAddress|string|false|none|none|
|exchangeAddress|string|false|none|none|
|baseToken|string|false|none|none|
|quoteToken|string|false|none|none|
|side|string|false|none|none|
|type|string|false|none|none|
|status|string|false|none|none|
|hash|string|false|none|none|
|signature|[Signature](#schemasignature)|false|none|none|
|pricepoint|string|false|none|none|
|amount|string|false|none|none|
|nonce|string|false|none|none|
|makeFee|string|false|none|none|
|takeFee|string|false|none|none|

<h2 id="tocSordercancel">OrderCancel</h2>

<a id="schemaordercancel"></a>

```json
{
  "orderHash": "string",
  "nonce": "string",
  "hash": "string",
  "signature": {
    "V": "string",
    "R": "string",
    "S": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|orderHash|string|false|none|none|
|nonce|string|false|none|none|
|hash|string|false|none|none|
|signature|[Signature](#schemasignature)|false|none|none|

<h2 id="tocSsignature">Signature</h2>

<a id="schemasignature"></a>

```json
{
  "V": "string",
  "R": "string",
  "S": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|V|string|false|none|none|
|R|string|false|none|none|
|S|string|false|none|none|

<h2 id="tocSorderbook">OrderBook</h2>

<a id="schemaorderbook"></a>

```json
{
  "pairName": "string",
  "asks": [
    {
      "amount": "string",
      "pricepoint": "string"
    }
  ],
  "bids": [
    {
      "amount": "string",
      "pricepoint": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pairName|string|false|none|none|
|asks|[object]|false|none|none|
|» amount|string|false|none|none|
|» pricepoint|string|false|none|none|
|bids|[object]|false|none|none|
|» amount|string|false|none|none|
|» pricepoint|string|false|none|none|

<h2 id="tocSraworderbook">RawOrderBook</h2>

<a id="schemaraworderbook"></a>

```json
{
  "pairName": "string",
  "orders": [
    {
      "id": "string",
      "userAddress": "string",
      "exchangeAddress": "string",
      "baseToken": "string",
      "quoteToken": "string",
      "status": "string",
      "side": "string",
      "type": "string",
      "hash": "string",
      "signature": {
        "V": "string",
        "R": "string",
        "S": "string"
      },
      "pricepoint": "string",
      "amount": "string",
      "filledAmount": "string",
      "nonce": "string",
      "makeFee": "string",
      "takeFee": "string",
      "pairName": "string",
      "createdAt": "2019-11-03T15:35:48Z",
      "updatedAt": "2019-11-03T15:35:48Z"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pairName|string|false|none|none|
|orders|[[Order](#schemaorder)]|false|none|none|

<h2 id="tocStrade">Trade</h2>

<a id="schematrade"></a>

```json
{
  "id": "string",
  "taker": "string",
  "maker": "string",
  "baseToken": "string",
  "quoteToken": "string",
  "makerOrderHash": "string",
  "takerOrderHash": "string",
  "hash": "string",
  "txHash": "string",
  "pairName": "string",
  "pricepoint": "string",
  "amount": "string",
  "status": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|taker|string|false|none|none|
|maker|string|false|none|none|
|baseToken|string|false|none|none|
|quoteToken|string|false|none|none|
|makerOrderHash|string|false|none|none|
|takerOrderHash|string|false|none|none|
|hash|string|false|none|none|
|txHash|string|false|none|none|
|pairName|string|false|none|none|
|pricepoint|string|false|none|none|
|amount|string|false|none|none|
|status|string|false|none|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocStick">Tick</h2>

<a id="schematick"></a>

```json
{
  "pair": {
    "pairName": "string",
    "baseToken": "string",
    "quoteToken": "string"
  },
  "open": "string",
  "high": "string",
  "low": 0,
  "close": "string",
  "volume": "string",
  "count": "string",
  "timestamp": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pair|[PairID](#schemapairid)|false|none|none|
|open|string|false|none|none|
|high|string|false|none|none|
|low|integer(int32)|false|none|none|
|close|string|false|none|none|
|volume|string|false|none|none|
|count|string|false|none|none|
|timestamp|string|false|none|none|

<h2 id="tocSnotification">Notification</h2>

<a id="schemanotification"></a>

```json
{
  "id": "string",
  "recipient": "string",
  "message": "string",
  "type": "string",
  "status": "string",
  "createdAt": "2019-11-03T15:35:48Z",
  "updatedAt": "2019-11-03T15:35:48Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|read-only|none|
|recipient|string|false|none|none|
|message|string|false|none|none|
|type|string|false|none|none|
|status|string|false|none|none|
|createdAt|string(date-time)|false|read-only|none|
|updatedAt|string(date-time)|false|read-only|none|

<h2 id="tocSapiresponse">ApiResponse</h2>

<a id="schemaapiresponse"></a>

```json
{
  "code": 0,
  "type": "string",
  "message": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|integer(int32)|false|none|none|
|type|string|false|none|none|
|message|string|false|none|none|

