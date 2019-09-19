# Hyperledger Chaincode

## Evaluation

### 거래의 평가

거래는 판매자와 구매자 모두 *공통으로 평가*할 수 있는 **상대방의 태도**에 대해서 평가할 수 있도록 `평가항목`이 구성되어 있다.
`평가항목`은 **3**개이다.

#### 평가항목

```
1. 거래는 명시된대로 잘 진행 되었습니까? (상품에 문제가 없는지 또는 서비스에 문제는 없었는지, 기간은 잘 지켜졌는지 등)
2. 거래를 진행하는데 상대방의 태도가 우수했습니까? (반말, 욕설등을 사용하지 않았는지)
3. 필요하다면 상대방과 다시 거래하실 의향이 있으십니까? (재구매 또는 재방문 의향)
```

#### 평가방법

아래와 같이 점수로 환산하며, 점수를 거래정보에 저장할때에는 배열로 평가항목 1,2,3의 순서대로 저정한다. 평가를 하지 않은 경우 모든 항목에 0점 처리하며, 평가점수를 집계할때 제외한다.

```
(1점) 매우 아니다. -강한 부정
(2점) 아니다. -부정
(3점) 보통이다.
(4점) 그렇다. -긍정
(5점) 매우 그렇다. -강한 긍정
```

#### 평가의 활용

평가점수를 집계한 데이터는 사용자 data-set에 저장한다. 사용자 data-set에서는 평가점수의 합을 저장하며, 평균 점수는 총 거래량으로 나누어 계산한다.  
판매 또는 구매 어느 한쪽의 거래량이 많은 사용자에 대해서는 총 평균 점수보다 판매시 평균점수, 구매시 평균점수를 구분하는 것이 신뢰도를 가늠하는데 효율적일 수 있기 때문에
구분하여 제공한다.

사용자가 상대방의 점수를 등록하지 않은 경우, 상대방의 점수는 전부 0점처리되며 집계에 합산하지 않는다. SellEx와 BuyEx는 판매, 구매시 상대방으로부터 평가 받지 않은 거래의 양이다. 
따라서 구매자의 경우 BuyEx, 판매자의 경우 SellEx가 +1 처리된다. 이를 감안하여 평균을 계산할때 다음과 같은 식으로 계산한다.

---

### Data 저장

`TradeId`는 무작위로 32자리 해시를 생성하여 저장한다. 해시임을 표현하기 위하여 `0x` 를 붙여 저장한다.  
임시 평가점수는 **AES256 GCM**모드로 암호화하여 저장한다.
암호화시 필요한 키는 길이에 상관없는 문자열로 받으며, 해당 문자열을 **SHA256**으로 해시하여 32자리의 AES256용 키를 생성한다.  
GCM 모드에서 사용되는 `Nonce` 값은 24자리의 문자열로 받으며, go에서 사용하는 time.Time()의 시간포맷을 이용하여 yyyymmdd**A**HHMMSSNNNNNNNNN 의 24자리로 생성한다.
중간의 **A** 는 24자리를 맞춰주기 위해 넣어준 고정값이다.  
생성한 문자열을 decoding 하여 12자리의 16진수 문자열로 변환하여 실제 **AES256 GCM**모드에서 사용할 수 있는 `Nonce`로 변경한다.
 
---



###### 예시
```
판매 평균점수 = 판매 평가점수 합 / 평가 유효 판매량
판매 평균점수 = 판매 평가점수 합 / (총 판매량 - 평가받지 않은 판매량)

구매 평균점수 = 구매 평가점수 합 / 평가 유효 구매량
구매 평균점수 = 구매 평가점수 합 / (총 구매량 - 평가받지 않은 구매량)
```  

> Sell EvalAvg01 = Sellsum.EvalSum01 / (SellAmt - SellEx)  
> Buy TotAvg = BuySum.TotSum / (BuyAmt - BuyEx)  
> Trade TotAvg = TradeSum.TotSum / ((SellAmt + BuyAmt) - (SellEx + BuyEx))

---

### User
#### Data-set
###### 설명
```
RecType : 데이터 셋의 성격을 구분하는 ID (User는 1)
UserTkn : 사용자 토큰. 사용자가 누구인지 확인할 때 사용하는 고유한 ID이며 hash값을 문자열로 저장한다. (data-set key)
SellAmt : 판매한 거래의 양.
BuyAmt : 구매한 거래의 양.
SellEx : 판매시, 평가받지 않은 거래의 양. (0점 처리된 판매량)
BuyEx : 구매시, 평가받지 않은 거래의 양. (0점 처리된 구매량)
Date : 사용자 생성날짜
SellSum : 판매 평가점수의 합
    TotSum : 판매에 대해 받은 평가점수의 합.
    EvalSum01 : 판매에 대해 1번 질문에 대하여 받은 평가점수의 합.
    EvalSum02 : 판매에 대해 2번 질문에 대하여 받은 평가점수의 합.
    EvalSum03 : 판매에 대해 3번 질문에 대하여 받은 평가점수의 합.
BuySum : 구매 평가점수의 합
    TotSum : 구매에 대해 받은 평가점수의 합.
    EvalSum01 : 구매에 대해 1번 질문에 대하여 받은 평가점수의 합.
    EvalSum02 : 구매에 대해 2번 질문에 대하여 받은 평가점수의 합.
    EvalSum03 : 구매에 대해 3번 질문에 대하여 받은 평가점수의 합.
TradeSum : 전체 거래 평가점수의 합
    TotSum : 구매에 대해 받은 평가점수의 합.
    EvalSum01 : 전체 거래에 대해 1번 질문에 대하여 받은 평가점수의 합.
    EvalSum02 : 전체 거래에 대해 2번 질문에 대하여 받은 평가점수의 합.
    EvalSum03 : 전체 거래에 대해 3번 질문에 대하여 받은 평가점수의 합. 
```

###### 예시
```
{
    RecType:1
    UserTkn:"0x03AC674216F3E15C761EE1A5E255F067953623C8B388B4459E13F978D7C846F4",
    SellAmt:15230,
    BuyAmt:24,
    SellEx:2367,
    BuyEx:4,
    Date:"20190904143256"
    SellSum:{
        TotSum:60920
        EvalSum01:20306
        EvalSum02:19250
        EvalSum03:21364
    },
    BuySum:{
        TotSum:108
        EvalSum01:34
        EvalSum02:41
        EvalSum03:33
    },
    TradeSum:{
        TotSum:61028
        EvalSum01:20340
        EvalSum02:19291
        EvalSum03:21397
}
```

#### Test

#### 사용자 생성

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["addUser","[UserTkn]"]}'
```

###### 예시
```
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["addUser","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD"]}'
```
 
#### 사용자 조회

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryUser","[UserTkn]"]}'
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args": ["queryUser","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD"]}';
```

#### **전체 내용 조회**

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args": ["queryUser","TOTAL_USER"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args": ["queryUser","TOTAL_USER"]}';
```

---

### Trade
#### Data-set
###### 설명

```
RecType : 데이터 셋의 성격을 구분하는 ID (Trade는 2)
TradeId : 거래에 대한 ID이며 hash값을 문자열로 저장한다. (data-set key)
ServiceCode : 서비스 코드는 화이트핀 프로젝트의 하위 서비스 제공자들에 대한 코드이다. 거래를 발생시킨 서비스가 어디인지 규명할때 사용한다.  
SellerTkn : 판매자 토큰. 판매자가 누구인지 확인할 때 사용하는 고유한 ID이며 hash값을 문자열로 저장한다.  
BuyerTkn : 구매자 토큰. 구매자가 누구인지 확인할 때 사용하는 고유한 ID이며 hash값을 문자열로 저장한다.  
Date : 거래가 생성된 시점, 즉 구매자가 구매 요청을 한 시각을 의미한다.
Close : 거래를 확정하는 시점.  
	SellDone : 판매자가 거래를 확정했는지 여부 (true : 확정, false : 미확정)  
	BuyDone : 구매자가 거래를 확정했는지 여부 (true : 확정, false : 미확정)
	SellDate : 판매자가 거래를 확정한 시각.
	BuyDate : 구매자가 거래를 확정한 시각.  
Score : 거래에 대한 평가 점수  
	SellScore : 판매자의 평가 점수. (판매자가 받은 평가 점수, 구매자가 판매자를 평가한 점수)
	BuyScore : 구매자의 평가 점수. (구매자가 받은 평가 점수, 판매자가 구매자를 평가한 점수)
```

###### 예시
```
{
    RecType:2
    TradeId:"0xA665A45920422F9D417E4867EFDC4FB8A04A1F3FFF1FA07E998E86F7F7A27AE3",
    ServiceCode:"Code11",
    SellerTkn:"0x03AC674216F3E15C761EE1A5E255F067953623C8B388B4459E13F978D7C846F4",
    BuyerTkn:"5994471ABB01112AFCC18159F6CC74B4F511B99806DA59B3CAF5A9C173CACFC5",
    Date:"20190904143256",
    Close:{
        SellDone:"",
        BuyDone:"",
        SellDate:"",
        BuyDate:""
    },
    Score:{
        SellScore:[5, 4, 3],
        BuyScore:[2, 5, 3]
    }
}
```

#### Test

#### 거래 생성

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["createTrade","[TradeId]","[ServiceCode]","[SellerTkn]","[BuyerTkn]"]}'
```

###### 예시
```
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["createTrade","0xF6E0A1E2AC41945A9AA7FF8A8AAA0CEBC12A3BCC981A929AD5CF810A090E11AE","code01","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","0xDF7E70E5021544F4834BBEE64A9E3789FEBC4BE81470DF629CAD6DDB03320A5C"]}'
```

#### 거래 완료

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["closeTrade","[TradeId]","[ServiceCode]","[SellerTkn]","[BuyerTkn]"]}'
```

###### 예시
```
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["closeTrade","0xF6E0A1E2AC41945A9AA7FF8A8AAA0CEBC12A3BCC981A929AD5CF810A090E11AE","code01","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","0xDF7E70E5021544F4834BBEE64A9E3789FEBC4BE81470DF629CAD6DDB03320A5C"]}'
```

#### 거래 점수 등록 (둘 다 평가를 마친 이후에)

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["enrollScore","[TradeId]","[AES-key]"]}'
```

###### 예시
```
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["enrollScore","0xF6E0A1E2AC41945A9AA7FF8A8AAA0CEBC12A3BCC981A929AD5CF810A090E11AE","key1234"]}';
```

#### 거래 조회 (key 사용(TradeID))

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["createTrade","[TradeId]","[ServiceCode]","[SellerTkn]","[BuyerTkn]"]}'
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["createTrade","0xF6E0A1E2AC41945A9AA7FF8A8AAA0CEBC12A3BCC981A929AD5CF810A090E11AE","code01","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","0xDF7E70E5021544F4834BBEE64A9E3789FEBC4BE81470DF629CAD6DDB03320A5C"]}'
```

#### 거래 조회 (query string 사용, 테스트용 메소드)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["createTrade","[TradeId]","[ServiceCode]","[SellerTkn]","[BuyerTkn]"]}'
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["createTrade","0xF6E0A1E2AC41945A9AA7FF8A8AAA0CEBC12A3BCC981A929AD5CF810A090E11AE","service01","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","0xDF7E70E5021544F4834BBEE64A9E3789FEBC4BE81470DF629CAD6DDB03320A5C"]}'
```

#### **내부 query string을 통한 거래 조회시 parameter**
[sell / buy / all] : 판매내역 조회 / 구매내역 조회 / 전체 조회  
[normal / page] : 페이징 없이 전체 조회(1000개) / 페이지로 조회  
PageSize : 한 페이지에 보여줄 데이터의 양  
PageNum : 페이지 번호  
Bookmark : 페이지를 점프하기 위해 필요한 북마크. 알면 쓰면 되지만 거의 대부분의 경우 바로 알 수 없으므로 ""로 받으면 됨.

#### 거래 조회 (사용자로 조회)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryTradeWithUser","[UserTkn]","[sell / buy / all]","[asc / desc]","[normal / page]","[PageSize]","[PageNum]","[Bookmark]"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUser","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","sell","desc","normal","","",""]}';
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUser","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","buy","asc","page","3","2",""]}'
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUser","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","all","","page","3","2",""]}'
* all의 경우 desc, asc 사용 불가능. (sort 불가능)
```

#### 거래 조회 (사용자, 서비스코드로 조회)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryTradeWithUserService","[UserTkn]","[ServiceCode]","[sell / buy / all]","[asc / desc]","[normal / page]","[PageSize]","[PageNum]","[Bookmark]"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUserService","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","service01","sell","desc","normal","","",""]}';
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUserService","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","service01","buy","asc","page","3","2",""]}'
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithUserService","0x559AEAD08264D5795D3909718CDD05ABD49572E84FE55590EEF31A88A08FDFFD","service01","all","","page","3","2",""]}'
* all의 경우 desc, asc 사용 불가능. (sort 불가능)
```

#### 거래 조회 (서비스코드로 조회)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryTradeWithService","[ServiceCode]","[sell / buy / all]","[asc / desc]","[normal / page]","[PageSize]","[PageNum]","[Bookmark]"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithService","code01","desc","normal","","",""]}';
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryTradeWithService","code01","asc","page","3","2",""]}';
```

---

### ScoreTemp
#### Data-set
###### 설명
```
RecType : 데이터 셋의 성격을 구분하는 ID (ScoreTemp는 3)
ScoreKey : 임시 평가점수에 대한 키로 사용되며 무작위 값이다. (data-set key)
TradeId : 거래에 대한 ID이며 hash값을 문자열로 저장한다. (Trade data-set과 동일)
ExpiryDate : 평가를 받기위한 만료 시간. 이 시간이 지나면 전부 0점 처리하여 공개한다.
Score : 거래에 대한 상호 평가 점수.
    SellScore : 판매자의 평가 점수. (구매자가 판매자를 평가한 점수) Trade data-set에 들어갈 점수를 보여주기 전, 점수 전문을 양방향 암호화한 문자열 저장. (공개시점 이전 공개 불가능하도록 하기 위해서) 
    BuyScore : 구매자의 평가 점수. (판매자가 구매자를 평가한 점수) Trade data-set에 들어갈 점수를 보여주기 전, 점수 전문을 양방향 암호화한 문자열 저장. (공개시점 이전 공개 불가능하도록 하기 위해서)
IsExpired : 만료됐는지 여부 (false: 만료 안됨, true: 만료됨)
```

###### 예시
```
{
    RecType:3
    ScoreKey:"0x0x03AC674216F3E15C761EE1A5E255F067953623C8B388B4459E13F978D7C846F4"
    TradeId:"0xA665A45920422F9D417E4867EFDC4FB8A04A1F3FFF1FA07E998E86F7F7A27AE3"
    ExpiryDate:"20190904143256"
    Score:{
        SellScore:"BE736DE7249081C95A41CBAF762A92B95E280DE155CACBFBF480E0059FAF88A6",
        BuyScore:"C78BBE24FCC51F8900F2D50FF4894A2136C6FBCF2CE321FD0D94C78E5F234C68"
    }
    IsExpired:false
}
```

#### Test

#### 임시 평가 점수 등록 (생성은 거래 생성시 자동으로 생성됨. 별도의 메소드 제공 없음)

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["enrollTempScore","[TradeId]","[UserTkn]","[A score, B score, C score]","[AES-key]"]}';
```

###### 예시
```
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["enrollTempScore","AB01","BB","[4,4,4]","key1234"]}';
```

#### 임시 평가 점수 조회 (key, 테스트용 메소드)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryScoreTemp","[ScoreTempKey]"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryScoreTemp","AB01_ScoreTemp"]}';
```

#### 임시 평가 점수 조회 (TradeId)

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["queryScoreTempWithTradeId","[TradeId]"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["queryScoreTempWithTradeId","AB01"]}';
```

---

### Properties
#### Data-set
###### 설명
```
EvaluationLimit : 평가 입력 기다려주는 시간 (default 14일, 1,209,600 = 14 * 24 * 60 * 60) 이시간 이후에는 0점 처리
OpenScoreDuration : 거래 당사자들의 모든 평가 입력 후 공개하기 까지 걸리는 시간 (default 5일, 432,000 = 5 * 24 * 60 * 60)
```

###### 예시 (2분(120초), 30초)
```
{
    EvaluationLimit:120000000000
    OpenScoreDuration:30000000000
}
```


#### Test

#### Property 조회

###### 설명
```
peer chaincode query -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["getProperties"]}';
```

###### 예시
```
peer chaincode query -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["getProperties"]}';
```

#### Property 변경

###### 설명
```
peer chaincode invoke -C [channel_name] --peerAddresses [peer_DNS:port] -n [chaincode_name] -c '{"Args":["setProperties",[evaluation_waiting_duration]","[enroll_score_delay]"]}';
```

###### 예시
```
# 시연용 (2분, 30초)
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["setProperties","120","30"]}';
# 운영용 (14일, 5일)
peer chaincode invoke -C ch1 --peerAddresses peer0.peerorg1.testnet.com:7051 -n whitepin -c '{"Args":["setProperties","1209600","432000"]}';
```

---
