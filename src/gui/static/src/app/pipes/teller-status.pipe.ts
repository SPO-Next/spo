import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tellerStatus'
})
export class TellerStatusPipe implements PipeTransform {

  transform(value: any): any {
    console.log(value);
    switch (value) {
      case 'waiting_deposit':
        return 'Waiting for deposit or confirmation';
      case 'waiting_send':
        return 'Waiting to send Spo';
      case 'waiting_confirm':
        return 'Waiting for confirmation';
      case 'done':
        return 'Completed';
      default:
        return 'Unknown';
    }
  }
}
