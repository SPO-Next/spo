import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { PurchaseService } from '../../../../services/purchase.service';
import { MdDialogRef } from '@angular/material';

declare var QRCode: any;

@Component({
  selector: 'app-add-deposit-address',
  templateUrl: './add-deposit-address.component.html',
  styleUrls: ['./add-deposit-address.component.css']
})
export class AddDepositAddressComponent implements OnInit {
  
  @ViewChild('DepositQR') DepositQR: ElementRef;
  
  form: FormGroup;
  tokenType: string;
  loading: Boolean = false;
  tokenAddress: string;
  subscribeRef: any;
  
  constructor(
    public walletService: WalletService,
    private dialogRef: MdDialogRef<AddDepositAddressComponent>,
    private formBuilder: FormBuilder,
    private purchaseService: PurchaseService,
    
  ) {}

  ngOnInit() {
    this.initForm();
  }
  
  ngOnDestroy() {
    if(this.subscribeRef) {
      this.subscribeRef.unsubscribe();
    }
  }

  generate(tokenType) {
    if (this.form.value.address === "") {
      alert("Please choose an address");
      return;
    }
    this.loading = true;
    console.log("tokenType:"+tokenType+":"+this.form.value.address);
    //return;
    this.subscribeRef = this.purchaseService.generate(this.form.value.address, tokenType).subscribe((e: any) => {
      console.log(e);
      if(e.code !== 0) {
        alert(e.errmsg);
        return;
      }
      this.showBindData(e.data);
      //this.dialogRef.close();
    }, (err) => {
      alert(err);  
    }, () => {
      this.loading = false;
    });
  }
  showBindData(aData: any) {
    this.tokenAddress = aData.tokenAddress;
    this.showQRCode(this.tokenAddress);
  }
  
  showQRCode(addr) {
    console.log(this.DepositQR)
    new QRCode(this.DepositQR.nativeElement, {
      text: addr,
      width: 200,
      height: 200,
      colorDark: '#000000',
      colorLight: '#ffffff',
      useSVG: false,
      correctLevel: QRCode.CorrectLevel['M']
    });    
  }
  
  private initForm() {
    this.form = this.formBuilder.group({
      address: ['', Validators.required],
    });
  }
  
  closeDialog() {
    this.dialogRef.close();
  }
}
