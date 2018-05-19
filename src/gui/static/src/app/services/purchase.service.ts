import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Http } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/timeout';

@Injectable()
export class PurchaseService {

  private purchaseOrders: Subject<any[]> = new BehaviorSubject<any[]>([]);
  private purchaseUrl: string = 'https://teller.spo.network/api/';
  // private purchaseUrl: string = '/teller/';

  private purchaseTokenTypes: Subject<TokenModel[]> = new BehaviorSubject<TokenModel[]>([]);

  //
  private tellerNotice: Subject<any[]> = new BehaviorSubject<any[]>([]);
  
  //buy types
  //private purchaseTokenTypes :Subject<any[]> = new BehaviorSubject<any[]>([]);

  constructor(
    private http: Http,
  ) {
    this.retrievePurchaseOrders();
    this.retrievePurchaseTokenTypes();
  }

  all() {
    return this.purchaseOrders.asObservable();
  }


  generate(address: string, tokenType: string) {
    console.log("address:"+address+", tokenType:"+tokenType);
    return this.post('bind', { address: address, tokenType: tokenType })
      .do(response => {
     
        console.log(response);
        if(response.code == 0) {      
          //一个token可以有多个地址
            this.purchaseOrders.first().subscribe(orders => {
              let index = orders.findIndex(order => order.address === address && order.tokenType === tokenType);
              if (index === -1) {
                orders.push({address: address, tokenType:tokenType, addresses: []});
                index = orders.length - 1;
              }
              const timestamp = Math.floor(Date.now() / 1000);
              orders[index].addresses.unshift({
                tokenType: response.data.tokenType,
                tokenAddress: response.data.tokenAddress,
                status: 'waiting_deposit',
                created: timestamp,
                updated: timestamp,
              });
              this.updatePurchaseOrders(orders)
            });
        }
      });
      
  }
  
  getCoinType(tokenType:string) {
    console.log(tokenType);
    if(tokenType == "skycoin") {
      return "sky";
    } else if(tokenType == "bitcoin") {
      return "btc";
    } else if(tokenType == "ethcoin") {
      return "eth";
    } else {
      return "unknow";
    }
  }
  scan(address: string,tokenType: string) {
    return this.get('status?address=' + address+"&tokenType="+tokenType).do(response => {
      if(response.code != 0) {
        console.log(response);        
        return;
      }
      this.purchaseOrders.first().subscribe(orders => {
       // for(var idx in orders) {
       //   console.log(orders[idx]);
       // }
        let coinType = this.getCoinType(tokenType);

      //  console.log(address+":"+coinType);//
        //find the items
        let index = orders.findIndex(order => order.address === address && order.tokenType === coinType);
      //  console.log("index:"+index);
        //orders[index].addresses一个spo针对同一个币种，可能有多个地址

        // Sort addresses ascending by creation date to match teller status response
        orders[index].addresses.sort((a, b) =>  b.created - a.created);
       // console.log(orders[index]);
        
       // {"errmsg":"","code":0,"data":
       //{"statuses":[
      //{"seq":0,"update_at":1509881371,"address":"277BPWQYRVgUccUPZ3iCsJ9JrBwc4mEJ76","tokenType":"skycoin","status":"done"},
     // {"seq":1,"update_at":1509881380,"address":"277BPWQYRVgUccUPZ3iCsJ9JrBwc4mEJ76","tokenType":"skycoin","status":"done"},
     //{"seq":2,"update_at":1509881389,"address":"277BPWQYRVgUccUPZ3iCsJ9JrBwc4mEJ76","tokenType":"skycoin","status":"done"},
     //{"seq":3,"update_at":1509893810,"address":"ZaFfCQno5frpksPfd2CY1HWMyhpea7S3q4","tokenType":"skycoin","status":"waiting_deposit"},
     //{"seq":4,"update_at":1509893810,"address":"2YG74USfHaDzcZ5xbhVoht8DPMocxhMLaNn","tokenType":"skycoin","status":"waiting_deposit"}]}}

       // const status = response.data.statuses.splice(-1,1)[0];
        //console.log(status+":"+status);
        //console.log("orders[index].addresses");
        //console.log(orders[index].addresses);
        for(var idx in orders[index].addresses) {
          var status;
          var tmpAddress = orders[index].addresses[idx].tokenAddress;
          var seq = 0;
          for(var j in response.data.statuses) {//status里面有该地址该币种的所有状态
            var tmpStatus = response.data.statuses[j];
            if(tmpStatus.address == tmpAddress) {
                if(tmpStatus.seq >= seq) {
                  status = tmpStatus;
                }
            }
          }
          
          orders[index].addresses[idx].status = status.status
          orders[index].addresses[idx].updated =status.update_at;
        //  console.log( orders[index].addresses[idx]);
        }
       

        this.updatePurchaseOrders(orders)
      });
    });
  }
  
  alltokens(): Observable<TokenModel[]> {
    return this.purchaseTokenTypes.asObservable();
  }

  getSupportedTokens() {
    return this.get('tokens').do(response=>{
        this.purchaseTokenTypes.first().subscribe(tokens=> {
            this.updateTokenTypes(tokens)
        });
    });    
  }
    
  getTellerNotice() {
    return this.get('teller').do(response=>{
        this.tellerNotice.first().subscribe(notice=> {
            this.tellerNotice.next(notice);
        });
    });    
  }

  private get(url) {
    return this.http.get(this.purchaseUrl + url).timeout(15000)
      .map((res: any) => res.json())
  }

  private post(url, parameters = {}) {
    return this.http.post(this.purchaseUrl + url, parameters).timeout(15000)
      .map((res: any) => res.json())
  }

  private retrievePurchaseOrders() {
    const orders = JSON.parse(window.localStorage.getItem('purchaseOrders'));
    if (orders) {
      this.purchaseOrders.next(orders);
    }
  }

  private retrievePurchaseTokenTypes() {
    const tokenTypes = JSON.parse(window.localStorage.getItem('purchaseTokenTypes'));
    if(tokenTypes) {
      this.purchaseTokenTypes.next(tokenTypes);
    }

  }

  private updatePurchaseOrders(collection: any[]) {
    this.purchaseOrders.next(collection);
    console.log(collection);
    window.localStorage.setItem('purchaseOrders', JSON.stringify(collection));
  }

  private updateTokenTypes(collection: any[]) {
    this.purchaseTokenTypes.next(collection);
    window.localStorage.setItem('purchaseTokenTypes', JSON.stringify(collection));
  }
}
