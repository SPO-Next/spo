import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';

import { PurchaseService } from '../../../services/purchase.service';
import { MdDialog, MdDialogConfig, MdDialogRef } from '@angular/material';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';
import { config } from '../../../app.config';
import { QrCodeComponent } from '../../layout/qr-code/qr-code.component';


@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.css']
})
export class BuyComponent {

  orders = [];
  otcEnabled: boolean;
  scanning = false;
  supportedTokens = [];
  tellerNotice = {};

  

  constructor(
    public purchaseService: PurchaseService,
    private dialog: MdDialog,
  ) {
    this.otcEnabled = config.otcEnabled;
    this.load();
  }

  addDepositAddress(token) {
    //console.log(token)
    const config = new MdDialogConfig();
    config.width = '500px';
    let dialogRef:MdDialogRef<AddDepositAddressComponent> = this.dialog.open(AddDepositAddressComponent, config);
    dialogRef.componentInstance.tokenType = token;
  }

  searchDepositAddress(address: string,tokenType:string) {
    this.scanning = true;
    this.purchaseService.scan(address,tokenType).subscribe(() => {
      this.disableScanning();
    }, error => {
      this.disableScanning();
    });
  }

  load() {
    //this.purchaseService.getSupportedTokens();
    this.getSupportedTokens();
    this.getTellerNotice();
  }
  
  getSupportedTokens(){
    this.purchaseService.getSupportedTokens().subscribe(result => {
      //console.log(result);
      if (result.code != 0) {
        alert(result.errmsg);
        return;
      }
      this.supportedTokens = result.data;        
    })    
  }

  getTellerNotice() {
    this.purchaseService.getTellerNotice().subscribe(result=> {
      if (result.code != 0) {
        alert(result.errmsg);
        return;
      }
      this.tellerNotice = result.data;  
    })
  }
  
  private disableScanning()
  {
    setTimeout(() => this.scanning = false, 1000);
  }
  showQr(address) {
    const config = new MdDialogConfig();
    config.data =  {address: address};
    this.dialog.open(QrCodeComponent, config);
  }
  
 

  ngOnInit(): void {
    
  }
}
